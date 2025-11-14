package utils

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// GetEmailConfig retrieves email configuration from environment variables
func GetEmailConfig() *EmailConfig {
	// Try new config first, fall back to old config
	smtpHost := Env("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = Env("EMAIL_HOST")
	}

	smtpPort := Env("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = Env("EMAIL_PORT")
	}

	smtpUsername := Env("SMTP_MAIL")
	if smtpUsername == "" {
		smtpUsername = Env("EMAIL_USERNAME")
	}

	smtpPassword := Env("SMTP_PASSWORD")
	if smtpPassword == "" {
		smtpPassword = Env("EMAIL_PASSWORD")
	}

	fromEmail := Env("SMTP_MAIL")
	if fromEmail == "" {
		fromEmail = Env("EMAIL_FROM")
	}

	return &EmailConfig{
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
		FromEmail:    fromEmail,
		FromName:     "CertiKiosk",
	}
}

// SendEmail sends an email with optional attachment
func SendEmail(to, subject, body string, attachment []byte, attachmentName string) error {
	config := GetEmailConfig()

	// Validate configuration
	if config.SMTPHost == "" || config.SMTPPort == "" {
		return fmt.Errorf("SMTP configuration is missing")
	}

	// Set default from email if not configured
	fromEmail := config.FromEmail
	if fromEmail == "" {
		fromEmail = config.SMTPUsername
	}

	fromName := config.FromName
	if fromName == "" {
		fromName = "CertiKiosk"
	}

	// Build email headers
	from := fmt.Sprintf("%s <%s>", fromName, fromEmail)

	// Setup authentication
	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)

	// Compose the email message
	var message string

	if attachment != nil && attachmentName != "" {
		// Email with attachment (multipart)
		boundary := "boundary123456789"
		message = fmt.Sprintf(
			"From: %s\r\n"+
				"To: %s\r\n"+
				"Subject: %s\r\n"+
				"MIME-Version: 1.0\r\n"+
				"Content-Type: multipart/mixed; boundary=%s\r\n"+
				"\r\n"+
				"--%s\r\n"+
				"Content-Type: text/html; charset=UTF-8\r\n"+
				"\r\n"+
				"%s\r\n"+
				"\r\n"+
				"--%s\r\n"+
				"Content-Type: application/pdf; name=\"%s\"\r\n"+
				"Content-Transfer-Encoding: base64\r\n"+
				"Content-Disposition: attachment; filename=\"%s\"\r\n"+
				"\r\n"+
				"%s\r\n"+
				"--%s--\r\n",
			from, to, subject, boundary,
			boundary,
			body,
			boundary, attachmentName, attachmentName,
			encodeBase64(attachment),
			boundary,
		)
	} else {
		// Simple text/html email
		message = fmt.Sprintf(
			"From: %s\r\n"+
				"To: %s\r\n"+
				"Subject: %s\r\n"+
				"MIME-Version: 1.0\r\n"+
				"Content-Type: text/html; charset=UTF-8\r\n"+
				"\r\n"+
				"%s\r\n",
			from, to, subject, body,
		)
	}

	// Send email
	addr := fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort)
	err := smtp.SendMail(addr, auth, fromEmail, []string{to}, []byte(message))

	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// encodeBase64 encodes bytes to base64 string with line breaks every 76 characters
func encodeBase64(data []byte) string {
	// Use Go's standard base64 encoding
	encoded := base64.StdEncoding.EncodeToString(data)

	// Add line breaks every 76 characters as per RFC 2045
	var result strings.Builder
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		result.WriteString(encoded[i:end])
		if end < len(encoded) {
			result.WriteString("\r\n")
		}
	}

	return result.String()
}

// SendDocumentEmail sends a document via email with a formatted template
func SendDocumentEmail(to, documentType, documentID string, pdfData []byte) error {
	fmt.Printf("📧 SendDocumentEmail called with %d bytes\n", len(pdfData))
	if len(pdfData) > 4 {
		fmt.Printf("📧 PDF signature in SendDocumentEmail: %s\n", string(pdfData[:4]))
	}

	subject := fmt.Sprintf("Your %s Document from CertiKiosk", documentType)

	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
				.content { padding: 20px; background-color: #f9f9f9; }
				.footer { text-align: center; padding: 20px; font-size: 12px; color: #777; }
				.button { background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block; margin: 10px 0; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>CertiKiosk Document Delivery</h1>
				</div>
				<div class="content">
					<h2>Hello,</h2>
					<p>Your requested document is now available. Please find your <strong>%s</strong> document attached to this email.</p>
					<p><strong>Document ID:</strong> %s</p>
					<p>The document is attached as a PDF file. If you have any questions or issues accessing the document, please contact our support team.</p>
					<p>Thank you for using CertiKiosk!</p>
				</div>
				<div class="footer">
					<p>&copy; 2025 CertiKiosk. All rights reserved.</p>
					<p>This is an automated email. Please do not reply to this message.</p>
				</div>
			</div>
		</body>
		</html>
	`, documentType, documentID)

	filename := fmt.Sprintf("%s_%s.pdf", documentType, documentID)

	fmt.Printf("📧 Calling SendEmail with attachment: %s (%d bytes)\n", filename, len(pdfData))
	return SendEmail(to, subject, body, pdfData, filename)
}

// SendDocumentEmailWithStamp sends a document via email with optional stamp image
func SendDocumentEmailWithStamp(to, documentType, documentID string, fileData []byte, stampData []byte, fileExt, mimeType string) error {
	fmt.Printf("📧 SendDocumentEmailWithStamp called with %d bytes, type: %s, stamp: %d bytes\n", len(fileData), fileExt, len(stampData))

	// Convert images to PDF with stamp (either provided stamp or text stamp)
	var finalData []byte
	var finalExt string
	var finalMimeType string

	if fileExt == "png" || fileExt == "jpg" || fileExt == "jpeg" {
		if len(stampData) > 0 {
			fmt.Printf("🔄 Converting %s image to PDF with provided stamp image...\n", fileExt)
			pdfData, err := ConvertImageToPDFWithImageStamp(fileData, stampData, fileExt, documentType)
			if err != nil {
				fmt.Printf("⚠️ Failed to convert with image stamp: %v, trying text stamp\n", err)
				pdfData, err = ConvertImageToPDFWithStamp(fileData, fileExt, documentType)
				if err != nil {
					fmt.Printf("⚠️ Failed to convert: %v, sending as %s\n", err, fileExt)
					finalData = fileData
					finalExt = fileExt
					finalMimeType = mimeType
				} else {
					finalData = pdfData
					finalExt = "pdf"
					finalMimeType = "application/pdf"
				}
			} else {
				fmt.Printf("✅ Successfully converted %s to PDF with stamp: %d bytes\n", fileExt, len(pdfData))
				finalData = pdfData
				finalExt = "pdf"
				finalMimeType = "application/pdf"
			}
		} else {
			fmt.Printf("🔄 Converting %s image to PDF with text stamp...\n", fileExt)
			pdfData, err := ConvertImageToPDFWithStamp(fileData, fileExt, documentType)
			if err != nil {
				fmt.Printf("⚠️ Failed to convert: %v, sending as %s\n", err, fileExt)
				finalData = fileData
				finalExt = fileExt
				finalMimeType = mimeType
			} else {
				fmt.Printf("✅ Successfully converted %s to PDF: %d bytes\n", fileExt, len(pdfData))
				finalData = pdfData
				finalExt = "pdf"
				finalMimeType = "application/pdf"
			}
		}
	} else {
		finalData = fileData
		finalExt = fileExt
		finalMimeType = mimeType
	}

	subject := fmt.Sprintf("Your %s Document from CertiKiosk", documentType)

	// Adjust message based on file type
	fileTypeLabel := "file"
	if finalExt == "pdf" {
		fileTypeLabel = "PDF document"
	} else if finalExt == "png" {
		fileTypeLabel = "PNG image"
	} else if finalExt == "jpg" || finalExt == "jpeg" {
		fileTypeLabel = "JPEG image"
	}

	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
				.content { padding: 20px; background-color: #f9f9f9; }
				.footer { text-align: center; padding: 20px; font-size: 12px; color: #777; }
				.button { background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block; margin: 10px 0; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>CertiKiosk Document Delivery</h1>
				</div>
				<div class="content">
					<h2>Hello,</h2>
					<p>Your requested document is now available. Please find your <strong>%s</strong> document attached to this email.</p>
					<p><strong>Document ID:</strong> %s</p>
					<p>The document is attached as a %s. If you have any questions or issues accessing the document, please contact our support team.</p>
					<p>Thank you for using CertiKiosk!</p>
				</div>
				<div class="footer">
					<p>&copy; 2025 CertiKiosk. All rights reserved.</p>
					<p>This is an automated email. Please do not reply to this message.</p>
				</div>
			</div>
		</body>
		</html>
	`, documentType, documentID, fileTypeLabel)

	filename := fmt.Sprintf("%s_%s.%s", documentType, documentID, finalExt)

	fmt.Printf("📧 Calling SendEmailWithMime with attachment: %s (%d bytes, %s)\n", filename, len(finalData), finalMimeType)
	return SendEmailWithMime(to, subject, body, finalData, filename, finalMimeType)
}

// SendDocumentEmailWithType sends a document via email with auto-detected file type
func SendDocumentEmailWithType(to, documentType, documentID string, fileData []byte, fileExt, mimeType string) error {
	fmt.Printf("📧 SendDocumentEmailWithType called with %d bytes, type: %s\n", len(fileData), fileExt)

	// Convert images to PDF with certification stamp
	var finalData []byte
	var finalExt string
	var finalMimeType string

	if fileExt == "png" || fileExt == "jpg" || fileExt == "jpeg" {
		fmt.Printf("🔄 Converting %s image to PDF with certification stamp...\n", fileExt)
		pdfData, err := ConvertImageToPDFWithStamp(fileData, fileExt, documentType)
		if err != nil {
			fmt.Printf("⚠️ Failed to convert image to PDF: %v, sending as %s\n", err, fileExt)
			finalData = fileData
			finalExt = fileExt
			finalMimeType = mimeType
		} else {
			fmt.Printf("✅ Successfully converted %s to certified PDF: %d bytes\n", fileExt, len(pdfData))
			finalData = pdfData
			finalExt = "pdf"
			finalMimeType = "application/pdf"
		}
	} else {
		finalData = fileData
		finalExt = fileExt
		finalMimeType = mimeType
	}

	subject := fmt.Sprintf("Your %s Document from CertiKiosk", documentType)

	// Adjust message based on file type
	fileTypeLabel := "file"
	if finalExt == "pdf" {
		fileTypeLabel = "PDF document"
	} else if finalExt == "png" {
		fileTypeLabel = "PNG image"
	} else if finalExt == "jpg" || finalExt == "jpeg" {
		fileTypeLabel = "JPEG image"
	}

	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
				.content { padding: 20px; background-color: #f9f9f9; }
				.footer { text-align: center; padding: 20px; font-size: 12px; color: #777; }
				.button { background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block; margin: 10px 0; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>CertiKiosk Document Delivery</h1>
				</div>
				<div class="content">
					<h2>Hello,</h2>
					<p>Your requested document is now available. Please find your <strong>%s</strong> document attached to this email.</p>
					<p><strong>Document ID:</strong> %s</p>
					<p>The document is attached as a %s. If you have any questions or issues accessing the document, please contact our support team.</p>
					<p>Thank you for using CertiKiosk!</p>
				</div>
				<div class="footer">
					<p>&copy; 2025 CertiKiosk. All rights reserved.</p>
					<p>This is an automated email. Please do not reply to this message.</p>
				</div>
			</div>
		</body>
		</html>
	`, documentType, documentID, fileTypeLabel)

	filename := fmt.Sprintf("%s_%s.%s", documentType, documentID, finalExt)

	fmt.Printf("📧 Calling SendEmailWithMime with attachment: %s (%d bytes, %s)\n", filename, len(finalData), finalMimeType)
	return SendEmailWithMime(to, subject, body, finalData, filename, finalMimeType)
}

// SendEmailWithMime sends an email with specific MIME type for attachment
func SendEmailWithMime(to, subject, body string, attachment []byte, attachmentName, mimeType string) error {
	config := GetEmailConfig()

	// Validate configuration
	if config.SMTPHost == "" || config.SMTPPort == "" {
		return fmt.Errorf("SMTP configuration is missing")
	}

	// Set default from email if not configured
	fromEmail := config.FromEmail
	if fromEmail == "" {
		fromEmail = config.SMTPUsername
	}

	fromName := config.FromName
	if fromName == "" {
		fromName = "CertiKiosk"
	}

	// Build email headers
	from := fmt.Sprintf("%s <%s>", fromName, fromEmail)

	// Setup authentication
	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)

	// Compose the email message with attachment
	boundary := "boundary123456789"
	message := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: multipart/mixed; boundary=%s\r\n"+
			"\r\n"+
			"--%s\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s\r\n"+
			"\r\n"+
			"--%s\r\n"+
			"Content-Type: %s; name=\"%s\"\r\n"+
			"Content-Transfer-Encoding: base64\r\n"+
			"Content-Disposition: attachment; filename=\"%s\"\r\n"+
			"\r\n"+
			"%s\r\n"+
			"--%s--\r\n",
		from, to, subject, boundary,
		boundary,
		body,
		boundary, mimeType, attachmentName, attachmentName,
		encodeBase64(attachment),
		boundary,
	)

	// Send email
	addr := fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort)
	err := smtp.SendMail(addr, auth, fromEmail, []string{to}, []byte(message))

	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// ValidateEmailConfig checks if email configuration is properly set
func ValidateEmailConfig() error {
	config := GetEmailConfig()

	if config.SMTPHost == "" {
		return fmt.Errorf("SMTP_HOST is not configured")
	}
	if config.SMTPPort == "" {
		return fmt.Errorf("SMTP_PORT is not configured")
	}
	if config.SMTPUsername == "" {
		return fmt.Errorf("SMTP_USERNAME is not configured")
	}
	if config.SMTPPassword == "" {
		return fmt.Errorf("SMTP_PASSWORD is not configured")
	}

	return nil
}

// Initialize email configuration on startup
func init() {
	// Check if .env file exists and load it
	if _, err := os.Stat(".env"); err == nil {
		// .env exists, configuration should be loaded by main
	}
}
