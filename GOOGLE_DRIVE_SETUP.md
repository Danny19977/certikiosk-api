# Google Drive API Setup Guide

## Overview
This guide will help you set up Google Drive API integration for the CertiKiosk project to download and manage documents from Google Drive.

## Prerequisites
- Google Cloud Account
- Go 1.24.4 or higher
- CertiKiosk API running

## Step 1: Install Required Packages

Run the following command to install the required Go packages:

```bash
go get golang.org/x/oauth2@latest
go get google.golang.org/api@latest
go mod tidy
```

## Step 2: Google Cloud Console Setup

### 2.1 Create a New Project (if needed)
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Click on the project dropdown at the top
3. Click "New Project"
4. Enter project name: `CertiKiosk`
5. Click "Create"

### 2.2 Enable Google Drive API
1. In the Google Cloud Console, go to **APIs & Services** > **Library**
2. Search for "Google Drive API"
3. Click on "Google Drive API"
4. Click "Enable"

### 2.3 Configure OAuth Consent Screen
1. Go to **APIs & Services** > **OAuth consent screen**
2. Select "External" user type (or "Internal" if using Google Workspace)
3. Click "Create"
4. Fill in the required information:
   - App name: `CertiKiosk`
   - User support email: Your email
   - Developer contact information: Your email
5. Click "Save and Continue"
6. On the Scopes page:
   - Click "Add or Remove Scopes"
   - Add these scopes:
     - `https://www.googleapis.com/auth/drive.readonly`
     - `https://www.googleapis.com/auth/drive.file`
   - Click "Update" and then "Save and Continue"
7. Add test users (your email and any other users who need access)
8. Click "Save and Continue"
9. Review and click "Back to Dashboard"

### 2.4 Create OAuth 2.0 Credentials
1. Go to **APIs & Services** > **Credentials**
2. Click "Create Credentials" > "OAuth client ID"
3. Select "Web application"
4. Name it: `CertiKiosk Web Client`
5. Add authorized redirect URIs:
   - `http://localhost:8081/api/auth/google/callback`
   - `http://localhost:8081/oauth2callback`
6. Click "Create"
7. **IMPORTANT**: Copy the **Client ID** and **Client Secret**

### 2.5 Create API Key (Optional - for read-only public files)
1. Go to **APIs & Services** > **Credentials**
2. Click "Create Credentials" > "API Key"
3. Copy the API key
4. Click "Restrict Key"
5. Under "API restrictions", select "Restrict key"
6. Select "Google Drive API"
7. Click "Save"

## Step 3: Update Environment Variables

Your `.env` file has been updated with the following variables:

```env
# Google Drive API Configuration
GOOGLE_CLIENT_ID=808517840096-u95kk11epet9hq6a4eubsht8skt85udn.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=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
GOOGLE_REDIRECT_URL=http://localhost:8081/api/auth/google/callback
GOOGLE_API_KEY=YOUR_API_KEY_HERE (optional)
GOOGLE_TOKEN_FILE=token.json
```

**Note**: Replace `YOUR_API_KEY_HERE` with the API key you created in Step 2.5 (optional).

## Step 4: Make Google Drive Files Accessible

For files to be downloadable, they need to be either:

### Option A: Make Files Public (Simplest)
1. Open the file in Google Drive
2. Click "Share"
3. Click "Change to anyone with the link"
4. Set permissions to "Viewer"
5. Click "Done"

### Option B: Use Service Account (Recommended for Production)
1. Go to **APIs & Services** > **Credentials**
2. Click "Create Credentials" > "Service Account"
3. Fill in the details and create
4. Click on the service account email
5. Go to "Keys" tab
6. Click "Add Key" > "Create new key"
7. Choose JSON and download the file
8. Save it as `credentials.json` in your project root
9. Share your Google Drive files/folders with the service account email

### Option C: OAuth2 Token (Current Implementation)
This requires user authentication flow. The first time you access a file, you'll need to authenticate.

## Step 5: API Endpoints Available

### Public Endpoints (No Authentication Required)

#### 1. Download Google Drive File (Proxy - Bypasses CORS)
```
GET /api/public/documents/gdrive/download/:file_id
GET /api/public/documents/gdrive/download?file_id={FILE_ID}
```

**Example:**
```
http://localhost:8081/api/public/documents/gdrive/download/1KHBHi5LBNmMlBVpti1WY18H6CJZasyN_
```

#### 2. Get File Metadata
```
GET /api/public/documents/gdrive/metadata/:file_id
GET /api/public/documents/gdrive/metadata?file_id={FILE_ID}
```

**Example:**
```
http://localhost:8081/api/public/documents/gdrive/metadata/1KHBHi5LBNmMlBVpti1WY18H6CJZasyN_
```

**Response:**
```json
{
  "status": "success",
  "message": "File metadata retrieved",
  "data": {
    "id": "1KHBHi5LBNmMlBVpti1WY18H6CJZasyN_",
    "name": "document.pdf",
    "mimeType": "application/pdf",
    "size": 12345,
    "downloadUrl": "https://drive.google.com/uc?export=download&id=...",
    "viewUrl": "https://drive.google.com/file/d/.../view",
    "proxy_url": "/api/public/documents/gdrive/download/..."
  }
}
```

#### 3. Generate Stamped PDF
```
GET /api/public/documents/generate-stamped-pdf?file_id={FILE_ID}&citizen_name={NAME}&national_id={ID}
```

#### 4. Send Document via Email
```
POST /api/public/documents/send-email-gdrive
Content-Type: application/json

{
  "email": "user@example.com",
  "file_id": "1KHBHi5LBNmMlBVpti1WY18H6CJZasyN_",
  "document_type": "Birth Certificate",
  "document_name": "birth_certificate"
}
```

## Step 6: Update Frontend

Update your frontend `DocumentViewer.js` to use the backend proxy instead of direct Google Drive URLs:

```javascript
// Instead of direct Google Drive download:
const directUrl = `https://drive.google.com/uc?export=download&id=${fileId}`;

// Use the backend proxy:
const proxyUrl = `http://localhost:8081/api/public/documents/gdrive/download/${fileId}`;

// Fetch the file through proxy
const response = await fetch(proxyUrl);
const blob = await response.blob();
```

**Example Implementation:**
```javascript
const downloadGoogleDriveDocument = async (fileId) => {
  try {
    console.log('🔄 Downloading via backend proxy...');
    
    const proxyUrl = `http://localhost:8081/api/public/documents/gdrive/download/${fileId}`;
    
    const response = await fetch(proxyUrl);
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const blob = await response.blob();
    console.log('✅ Download successful via proxy');
    
    return blob;
  } catch (error) {
    console.error('❌ Download failed:', error);
    throw error;
  }
};
```

## Step 7: Run the Application

1. Install dependencies:
```bash
go mod tidy
```

2. Run the server:
```bash
go run main.go
```

3. Test the endpoints:
```bash
# Test metadata endpoint
curl http://localhost:8081/api/public/documents/gdrive/metadata/1KHBHi5LBNmMlBVpti1WY18H6CJZasyN_

# Test download endpoint
curl http://localhost:8081/api/public/documents/gdrive/download/1KHBHi5LBNmMlBVpti1WY18H6CJZasyN_ -o test.pdf
```

## Troubleshooting

### Error: "Failed to download file from Google Drive"

**Solution 1: Make the file public**
- Share the file with "Anyone with the link" can view

**Solution 2: Check file ID**
- Verify the file ID is correct
- Extract ID from Google Drive URL: `https://drive.google.com/file/d/{FILE_ID}/view`

**Solution 3: Check API credentials**
- Verify `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET` are correct in `.env`
- Ensure Google Drive API is enabled in Google Cloud Console

### Error: "CORS policy"

**Solution**: Use the backend proxy endpoints instead of direct Google Drive URLs
- ✅ Use: `http://localhost:8081/api/public/documents/gdrive/download/{fileId}`
- ❌ Don't use: `https://drive.google.com/uc?export=download&id={fileId}`

### Error: "403 Forbidden"

**Possible causes:**
1. File is private - make it public or share with service account
2. API key restrictions - check API key settings in Google Cloud Console
3. OAuth consent screen not configured - complete OAuth consent screen setup

### Error: "Token not found"

**Solution**: This happens when using OAuth2 without initial authentication
- Option 1: Use API key instead (set `GOOGLE_API_KEY` in `.env`)
- Option 2: Make files public
- Option 3: Implement service account authentication

## Security Best Practices

1. **Never commit credentials to Git**
   - Add `.env` to `.gitignore`
   - Use environment variables for all secrets

2. **Use HTTPS in production**
   - Update redirect URLs to use `https://`
   - Enable CORS only for your frontend domain

3. **Implement rate limiting**
   - Prevent abuse of download endpoints
   - Use API quotas in Google Cloud Console

4. **Validate file IDs**
   - Sanitize input to prevent injection attacks
   - Validate file ID format before processing

## Next Steps

1. ✅ Install Go packages
2. ✅ Configure Google Cloud Console
3. ✅ Update `.env` file with credentials
4. ✅ Make Google Drive files public or share with service account
5. ✅ Update frontend to use proxy endpoints
6. ✅ Test the implementation
7. 🔄 Deploy to production (update redirect URLs and CORS settings)

## Support

For issues or questions, refer to:
- [Google Drive API Documentation](https://developers.google.com/drive/api/v3/about-sdk)
- [OAuth 2.0 Documentation](https://developers.google.com/identity/protocols/oauth2)
