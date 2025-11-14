# 🔗 Google Drive Download Endpoints - Complete Reference

## ✅ Server Status
```
✅ Running on: http://localhost:8081
✅ Handlers: 143 (includes new Google Drive endpoints)
✅ Database: Connected
```

---

## 📍 Available Endpoints

Your backend now supports **MULTIPLE endpoint formats** for maximum compatibility:

### **Option 1: Frontend-Friendly Endpoint** ⭐ RECOMMENDED

```
GET http://localhost:8081/api/documents/download-google-drive?fileId={FILE_ID}
GET http://localhost:8081/api/public/documents/download-google-drive?fileId={FILE_ID}
```

**Parameters:**
- `fileId` (query param) - Google Drive file ID
- Also supports: `file_id` (snake_case)

**Example:**
```javascript
const fileId = '1-oabwbqpdDytWxBBFpOU2NSygkIPA2wh';
const url = `http://localhost:8081/api/documents/download-google-drive?fileId=${fileId}`;

const response = await fetch(url);
const blob = await response.blob();
```

**Browser Test:**
```
http://localhost:8081/api/documents/download-google-drive?fileId=1-oabwbqpdDytWxBBFpOU2NSygkIPA2wh
```

---

### **Option 2: RESTful Path Parameter**

```
GET http://localhost:8081/api/public/documents/gdrive/download/{FILE_ID}
```

**Example:**
```
http://localhost:8081/api/public/documents/gdrive/download/1-oabwbqpdDytWxBBFpOU2NSygkIPA2wh
```

---

### **Option 3: Query Parameter (snake_case)**

```
GET http://localhost:8081/api/public/documents/gdrive/download?file_id={FILE_ID}
```

**Example:**
```
http://localhost:8081/api/public/documents/gdrive/download?file_id=1-oabwbqpdDytWxBBFpOU2NSygkIPA2wh
```

---

## 🔒 Authentication

### Public Endpoints (No Auth Required):
```
✅ /api/public/documents/download-google-drive
✅ /api/public/documents/gdrive/download
✅ /api/public/documents/gdrive/download/:file_id
```

### Protected Endpoints (Requires Authentication):
```
🔐 /api/documents/download-google-drive
🔐 /api/documents/gdrive/download
🔐 /api/documents/gdrive/download/:file_id
```

**Use public endpoints for kiosk/frontend applications!**

---

## 📊 Response Format

### Success (File Download):
```
HTTP/1.1 200 OK
Content-Type: application/pdf (or image/jpeg, image/png, etc.)
Content-Length: 54321
Access-Control-Allow-Origin: *

[Binary file data]
```

### Error (Empty File):
```json
{
  "status": "error",
  "message": "Downloaded file is empty (0 bytes). The file may be private or the link is incorrect.",
  "data": null
}
```

### Error (Missing File ID):
```json
{
  "status": "error",
  "message": "Google Drive file ID is required (use ?fileId=YOUR_FILE_ID or ?file_id=YOUR_FILE_ID)",
  "data": null
}
```

### Error (Download Failed):
```json
{
  "status": "error",
  "message": "Failed to download file from Google Drive",
  "error": "All download attempts failed. Last error: received HTML page instead of file - file may be private or link sharing not enabled. Make sure the file is public (Anyone with the link can view)"
}
```

---

## 🎯 Frontend Implementation

### React/JavaScript Example:

```javascript
async function downloadGoogleDriveDocument(fileId) {
  try {
    console.log('📥 Downloading file:', fileId);
    
    // Use the frontend-friendly endpoint
    const url = `http://localhost:8081/api/documents/download-google-drive?fileId=${fileId}`;
    
    const response = await fetch(url);
    
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'Download failed');
    }
    
    const blob = await response.blob();
    
    if (blob.size === 0) {
      throw new Error('Downloaded file is empty (0 bytes). File may be private.');
    }
    
    console.log(`✅ Downloaded ${blob.size} bytes`);
    return blob;
    
  } catch (error) {
    console.error('❌ Download failed:', error);
    throw error;
  }
}

// Usage
const fileId = '1-oabwbqpdDytWxBBFpOU2NSygkIPA2wh';
const blob = await downloadGoogleDriveDocument(fileId);

// Display in viewer
const blobUrl = URL.createObjectURL(blob);
document.getElementById('viewer').src = blobUrl;

// Or download to disk
const a = document.createElement('a');
a.href = blobUrl;
a.download = 'document.pdf';
a.click();
```

---

## 🔍 Metadata Endpoints

Get file information without downloading:

```
GET /api/documents/google-drive-metadata?fileId={FILE_ID}
GET /api/public/documents/google-drive-metadata?fileId={FILE_ID}
GET /api/public/documents/gdrive/metadata/{FILE_ID}
```

**Response:**
```json
{
  "status": "success",
  "message": "File metadata retrieved",
  "data": {
    "file_id": "1-oabwbqpdDytWxBBFpOU2NSygkIPA2wh",
    "view_url": "https://drive.google.com/file/d/.../view",
    "download_url": "https://drive.google.com/uc?export=download&id=...",
    "proxy_url": "/api/public/documents/gdrive/download/..."
  }
}
```

---

## 🧪 Testing

### PowerShell Test:
```powershell
# Test with fileId parameter
Invoke-WebRequest -Uri "http://localhost:8081/api/documents/download-google-drive?fileId=1-oabwbqpdDytWxBBFpOU2NSygkIPA2wh" -OutFile "test.pdf"

# Check file size
(Get-Item "test.pdf").Length
```

### Browser Test:
```
http://localhost:8081/api/documents/download-google-drive?fileId=1-oabwbqpdDytWxBBFpOU2NSygkIPA2wh
```

### cURL Test:
```bash
curl "http://localhost:8081/api/documents/download-google-drive?fileId=1-oabwbqpdDytWxBBFpOU2NSygkIPA2wh" -o test.pdf
```

---

## ⚠️ Important: Make Files Public!

**The #1 reason for 0 bytes is private files!**

### How to Make File Public:
1. Open Google Drive
2. Find your file
3. Right-click → **Share**
4. Click **"Change to anyone with the link"**
5. Set role to **"Viewer"**
6. Click **"Done"**

### Verify File is Public:
Try this direct URL in your browser:
```
https://drive.google.com/uc?export=download&id=YOUR_FILE_ID
```

If it downloads, the file is public ✅
If you see a Google login page, the file is private ❌

---

## 📝 Server Logs

When you make a request, watch the terminal for logs:

**Success:**
```
📥 Downloading file from Google Drive: 1-oabwbqpdDytWxBBFpOU2NSygkIPA2wh
🔍 Attempting to download public file: 1-oabwbqpdDytWxBBFpOU2NSygkIPA2wh
📡 Attempt 1/4: https://drive.google.com/uc?export=download&id=...
✅ Success! Downloaded 54321 bytes
✅ Successfully downloaded 54321 bytes for file 1-oabwbqpdDytWxBBFpOU2NSygkIPA2wh
```

**Failure (Private File):**
```
📥 Downloading file from Google Drive: 1-oabwbqpdDytWxBBFpOU2NSygkIPA2wh
🔍 Attempting to download public file: 1-oabwbqpdDytWxBBFpOU2NSygkIPA2wh
📡 Attempt 1/4: https://drive.google.com/uc?export=download&id=...
❌ Attempt 1 failed: received HTML page instead of file - file may be private
📡 Attempt 2/4: https://drive.google.com/uc?id=...&export=download
❌ Attempt 2 failed: received HTML page instead of file - file may be private
...
⚠️ Downloaded 0 bytes for file 1-oabwbqpdDytWxBBFpOU2NSygkIPA2wh
```

---

## 🎯 Complete Endpoint List

### Public (No Authentication):
| Endpoint | Method | Parameters |
|----------|--------|-----------|
| `/api/public/documents/download-google-drive` | GET | `fileId` or `file_id` |
| `/api/public/documents/gdrive/download` | GET | `file_id` |
| `/api/public/documents/gdrive/download/:file_id` | GET | Path param |
| `/api/public/documents/google-drive-metadata` | GET | `fileId` or `file_id` |
| `/api/public/documents/gdrive/metadata` | GET | `file_id` |
| `/api/public/documents/gdrive/metadata/:file_id` | GET | Path param |

### Protected (Requires Authentication):
| Endpoint | Method | Parameters |
|----------|--------|-----------|
| `/api/documents/download-google-drive` | GET | `fileId` or `file_id` |
| `/api/documents/gdrive/download` | GET | `file_id` |
| `/api/documents/gdrive/download/:file_id` | GET | Path param |
| `/api/documents/google-drive-metadata` | GET | `fileId` or `file_id` |
| `/api/documents/gdrive/metadata` | GET | `file_id` |
| `/api/documents/gdrive/metadata/:file_id` | GET | Path param |

---

## ✅ Quick Troubleshooting

### Issue: Still getting 0 bytes
**Solution:** Make the file public! (See "Make Files Public" section above)

### Issue: File downloads but won't display
**Solution:** Check Content-Type header - backend auto-detects PDF, JPEG, PNG, GIF

### Issue: CORS error
**Solution:** Use the `/api/public/...` endpoints, they have CORS headers enabled

### Issue: 404 Not Found
**Solution:** Make sure server is running on port 8081

---

## 🚀 You're All Set!

The backend now supports:
- ✅ Multiple endpoint formats
- ✅ Both `fileId` and `file_id` parameters
- ✅ Public and protected routes
- ✅ Detailed error messages
- ✅ Automatic content-type detection
- ✅ CORS headers for frontend access

**Just make your Google Drive files public and it will work!** 🎉
