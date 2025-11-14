package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// GoogleDriveConfig holds configuration for Google Drive API
type GoogleDriveConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// GetGoogleDriveConfig returns the Google Drive OAuth2 configuration
func GetGoogleDriveConfig() *oauth2.Config {
	clientID := Env("GOOGLE_CLIENT_ID")
	clientSecret := Env("GOOGLE_CLIENT_SECRET")
	redirectURL := Env("GOOGLE_REDIRECT_URL")

	if redirectURL == "" {
		redirectURL = "http://localhost:8080/api/auth/google/callback"
	}

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			drive.DriveReadonlyScope,
			drive.DriveFileScope,
		},
		Endpoint: google.Endpoint,
	}
}

// GetGoogleDriveService initializes and returns a Google Drive service with API key
func GetGoogleDriveService() (*drive.Service, error) {
	ctx := context.Background()

	// Try to use API key first (simpler for read-only public files)
	apiKey := Env("GOOGLE_API_KEY")
	if apiKey != "" {
		srv, err := drive.NewService(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			return nil, fmt.Errorf("unable to create Drive service with API key: %v", err)
		}
		return srv, nil
	}

	// Fall back to OAuth2 flow
	config := GetGoogleDriveConfig()

	// Try to load token from file
	tokenFile := Env("GOOGLE_TOKEN_FILE")
	if tokenFile == "" {
		tokenFile = "token.json"
	}

	token, err := tokenFromFile(tokenFile)
	if err != nil {
		return nil, fmt.Errorf("unable to load token: %v. Please authenticate first", err)
	}

	client := config.Client(ctx, token)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	return srv, nil
}

// DownloadFileFromDrive downloads a file from Google Drive using the API
func DownloadFileFromDrive(fileID string) ([]byte, error) {
	srv, err := GetGoogleDriveService()
	if err != nil {
		// Fall back to direct download if API service is not available
		return DownloadPublicDriveFile(fileID)
	}

	// Get file metadata first to check permissions
	file, err := srv.Files.Get(fileID).Fields("id, name, mimeType, permissions").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve file metadata: %v", err)
	}

	// Download the file content
	resp, err := srv.Files.Get(fileID).Download()
	if err != nil {
		return nil, fmt.Errorf("unable to download file: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read file content: %v", err)
	}

	fmt.Printf("Downloaded file: %s (%s)\n", file.Name, file.MimeType)
	return data, nil
}

// GetFileMetadata retrieves metadata for a file from Google Drive
func GetFileMetadata(fileID string) (map[string]interface{}, error) {
	srv, err := GetGoogleDriveService()
	if err != nil {
		return nil, err
	}

	file, err := srv.Files.Get(fileID).
		Fields("id, name, mimeType, size, createdTime, modifiedTime, webViewLink, webContentLink, iconLink, thumbnailLink, owners, permissions").
		Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve file metadata: %v", err)
	}

	metadata := map[string]interface{}{
		"id":             file.Id,
		"name":           file.Name,
		"mimeType":       file.MimeType,
		"size":           file.Size,
		"createdTime":    file.CreatedTime,
		"modifiedTime":   file.ModifiedTime,
		"webViewLink":    file.WebViewLink,
		"webContentLink": file.WebContentLink,
		"iconLink":       file.IconLink,
		"thumbnailLink":  file.ThumbnailLink,
		"downloadUrl":    GetPublicFileURL(fileID),
		"viewUrl":        GetDriveViewURL(fileID),
	}

	if len(file.Owners) > 0 {
		metadata["owner"] = file.Owners[0].DisplayName
	}

	return metadata, nil
}

// GetPublicFileURL generates a public URL for a Google Drive file
func GetPublicFileURL(fileID string) string {
	return fmt.Sprintf("https://drive.google.com/uc?export=download&id=%s", fileID)
}

// GetDriveViewURL generates a view URL for a Google Drive file
func GetDriveViewURL(fileID string) string {
	return fmt.Sprintf("https://drive.google.com/file/d/%s/view", fileID)
}

// DownloadPublicDriveFile - Simple HTTP-based download for public Google Drive files
// This is a fallback method when OAuth2 is not configured
func DownloadPublicDriveFile(fileID string) ([]byte, error) {
	fmt.Printf("🔍 Attempting to download public file: %s\n", fileID)

	// Try multiple download URLs
	urls := []string{
		fmt.Sprintf("https://drive.google.com/uc?export=download&id=%s", fileID),
		fmt.Sprintf("https://drive.google.com/uc?id=%s&export=download", fileID),
		fmt.Sprintf("https://www.googleapis.com/drive/v3/files/%s?alt=media", fileID),
		fmt.Sprintf("https://docs.google.com/uc?export=download&id=%s", fileID),
	}

	var lastErr error
	for i, url := range urls {
		fmt.Printf("📡 Attempt %d/%d: %s\n", i+1, len(urls), url)
		data, err := tryDownloadURL(url)
		if err == nil && len(data) > 0 {
			fmt.Printf("✅ Success! Downloaded %d bytes\n", len(data))
			return data, nil
		}
		if err != nil {
			fmt.Printf("❌ Attempt %d failed: %v\n", i+1, err)
			lastErr = err
		}
	}

	return nil, fmt.Errorf("all download attempts failed. Last error: %v. Make sure the file is public (Anyone with the link can view)", lastErr)
}

// tryDownloadURL attempts to download from a specific URL
func tryDownloadURL(url string) ([]byte, error) {
	// Add API key if available
	apiKey := Env("GOOGLE_API_KEY")
	if apiKey != "" && len(url) > 0 {
		if url[len(url)-1] == '=' || url[len(url)-1] == 'a' {
			url = url + "&key=" + apiKey
		} else {
			url = url + "?key=" + apiKey
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add user agent to avoid some blocks
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)
		if len(bodyStr) > 200 {
			bodyStr = bodyStr[:200] + "..."
		}
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, bodyStr)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check if we got HTML error page instead of file
	if len(data) > 0 && data[0] == '<' {
		dataStr := string(data)
		if len(dataStr) > 500 {
			dataStr = dataStr[:500]
		}
		// Check for common error messages
		if len(dataStr) > 100 && (len(data) > 100 && string(data[:100]) == "<!DOCTYPE html>" ||
			len(data) > 20 && string(data[:20]) == "<html" ||
			len(data) > 15 && string(data[:15]) == "<!doctype html>") {
			return nil, fmt.Errorf("received HTML page instead of file - file may be private or link sharing not enabled")
		}
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("received empty response (0 bytes)")
	}

	return data, nil
}

// tokenFromFile retrieves a Token from a given file path
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tok oauth2.Token
	err = json.NewDecoder(f).Decode(&tok)
	return &tok, err
}

// saveToken saves a token to a file path
func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %v", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

// GetDriveFileInfo returns basic file information without requiring full API setup
func GetDriveFileInfo(fileID string) map[string]string {
	return map[string]string{
		"file_id":      fileID,
		"view_url":     GetDriveViewURL(fileID),
		"download_url": GetPublicFileURL(fileID),
	}
}
