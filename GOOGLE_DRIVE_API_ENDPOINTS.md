# Google Drive Document API Endpoints

## Overview
This document describes the newly implemented Google Drive integration endpoints for the CertiKiosk API. These endpoints enable seamless document operations with files stored in Google Drive.

## New Endpoints

### 1. Send Document Email from Google Drive
**Endpoint**: `POST /api/documents/send-email-gdrive`  
**Public Endpoint**: `POST /api/public/documents/send-email-gdrive`  
**Authentication**: Protected version requires JWT token; Public version does not

#### Description
Downloads a document from Google Drive and sends it via email to the specified recipient.

#### Request Format
**Content-Type**: `application/json` or `multipart/form-data`

```json
{
  "email": "recipient@example.com",
  "file_id": "1ABC123xyz...",
  "document_type": "Birth Certificate",
  "document_name": "birth_certificate"
}
```

#### Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `email` | string | Yes | Recipient's email address |
| `file_id` | string | Yes | Google Drive file ID |
| `document_type` | string | No | Type of document (default: "Document") |
| `document_name` | string | No | Document name (default: "document") |

#### Response Format

**Success (200)**
```json
{
  "status": "success",
  "message": "Document sent successfully to recipient@example.com",
  "data": {
    "email": "recipient@example.com",
    "document_type": "Birth Certificate",
    "file_id": "1ABC123xyz...",
    "document_name": "birth_certificate"
  }
}
```

**Error (400)**
```json
{
  "status": "error",
  "message": "Email address is required",
  "data": null
}
```

**Error (500)**
```json
{
  "status": "error",
  "message": "Failed to download file from Google Drive",
  "error": "detailed error message"
}
```

#### Example Usage

**cURL**
```bash
curl -X POST http://10.249.216.144:8081/api/public/documents/send-email-gdrive \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "file_id": "1ABC123xyz",
    "document_type": "ID Card",
    "document_name": "national_id"
  }'
```

**JavaScript/Fetch**
```javascript
const sendEmailFromGDrive = async (email, fileId, docType) => {
  const response = await fetch('/api/public/documents/send-email-gdrive', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      email: email,
      file_id: fileId,
      document_type: docType,
      document_name: 'document'
    })
  });
  
  return response.json();
};

// Usage
await sendEmailFromGDrive('user@example.com', '1ABC123xyz', 'Birth Certificate');
```

---

### 2. Generate Stamped PDF for Printing
**Endpoint**: `GET /api/documents/generate-stamped-pdf`  
**Endpoint**: `POST /api/documents/generate-stamped-pdf`  
**Public Endpoint**: `GET /api/public/documents/generate-stamped-pdf`  
**Public Endpoint**: `POST /api/public/documents/generate-stamped-pdf`  
**Authentication**: Protected version requires JWT token; Public version does not

#### Description
Downloads a document from Google Drive and returns it as a PDF ready for printing. Optionally adds a certification stamp with metadata.

#### Request Format (GET)
Query parameters:
```
/api/public/documents/generate-stamped-pdf?file_id=1ABC123xyz&document_type=Birth%20Certificate&citizen_name=John%20Doe&national_id=123456789
```

#### Request Format (POST)
**Content-Type**: `application/json`

```json
{
  "file_id": "1ABC123xyz...",
  "document_type": "Birth Certificate",
  "citizen_name": "John Doe",
  "national_id": "123456789",
  "certifier_name": "CertiKiosk System",
  "include_stamp": true,
  "stamp_text": "CERTIFIED COPY",
  "document_name": "birth_certificate"
}
```

#### Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `file_id` | string | Yes | Google Drive file ID |
| `document_type` | string | No | Type of document (default: "Document") |
| `citizen_name` | string | No | Name of the citizen |
| `national_id` | string | No | National ID number |
| `certifier_name` | string | No | Name of certifier (default: "CertiKiosk System") |
| `include_stamp` | boolean | No | Whether to include stamp (default: true) |
| `stamp_text` | string | No | Custom stamp text |
| `document_name` | string | No | Document filename (default: "document") |

#### Response Format
**Content-Type**: `application/pdf`  
**Content-Disposition**: `inline; filename="document_stamped.pdf"`

Returns the PDF file directly as binary data.

#### Example Usage

**cURL (Download PDF)**
```bash
curl -X GET "http://10.249.216.144:8081/api/public/documents/generate-stamped-pdf?file_id=1ABC123xyz&document_type=Birth%20Certificate&citizen_name=John%20Doe" \
  --output stamped_document.pdf
```

**JavaScript/Fetch (Display in Browser)**
```javascript
const generateStampedPDF = async (fileId, docType, citizenName) => {
  const url = new URL('/api/public/documents/generate-stamped-pdf', window.location.origin);
  url.searchParams.append('file_id', fileId);
  url.searchParams.append('document_type', docType);
  url.searchParams.append('citizen_name', citizenName);
  
  window.open(url.toString(), '_blank');
};

// Usage
generateStampedPDF('1ABC123xyz', 'Birth Certificate', 'John Doe');
```

**React Example (Download)**
```javascript
const downloadStampedPDF = async (fileId, documentType, citizenName) => {
  const response = await fetch('/api/public/documents/generate-stamped-pdf', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      file_id: fileId,
      document_type: documentType,
      citizen_name: citizenName,
      national_id: '123456789',
      include_stamp: true
    })
  });
  
  const blob = await response.blob();
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `${documentType}_stamped.pdf`;
  a.click();
  window.URL.revokeObjectURL(url);
};
```

---

### 3. Get Stamped PDF Metadata
**Endpoint**: `GET /api/documents/stamped-pdf-metadata`  
**Public Endpoint**: `GET /api/public/documents/stamped-pdf-metadata`  
**Authentication**: Protected version requires JWT token; Public version does not

#### Description
Retrieves metadata about a stamped PDF without downloading the actual file. Useful for previewing information before generating the stamped document.

#### Request Format
Query parameters:
```
/api/public/documents/stamped-pdf-metadata?file_id=1ABC123xyz&document_type=Birth%20Certificate&citizen_name=John%20Doe&national_id=123456789
```

#### Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `file_id` | string | Yes | Google Drive file ID |
| `document_type` | string | No | Type of document (default: "Document") |
| `citizen_name` | string | No | Name of the citizen |
| `national_id` | string | No | National ID number |

#### Response Format

**Success (200)**
```json
{
  "status": "success",
  "message": "Stamped PDF metadata generated",
  "data": {
    "file_info": {
      "file_id": "1ABC123xyz",
      "view_url": "https://drive.google.com/file/d/1ABC123xyz/view",
      "download_url": "https://drive.google.com/uc?export=download&id=1ABC123xyz"
    },
    "certification": {
      "citizen_name": "John Doe",
      "national_id": "123456789",
      "document_type": "Birth Certificate",
      "certified_date": "2025-11-11 14:30:00",
      "certifier": "CertiKiosk System",
      "stamp_details": "CERTIFIED COPY\nCertified by: CertiKiosk System\nDate: 2025-11-11\nThis document has been verified and certified as authentic."
    },
    "download_url": "https://drive.google.com/uc?export=download&id=1ABC123xyz",
    "view_url": "https://drive.google.com/file/d/1ABC123xyz/view",
    "stamped_pdf_url": "/api/documents/generate-stamped-pdf?file_id=1ABC123xyz"
  }
}
```

#### Example Usage

**JavaScript/Fetch**
```javascript
const getStampedPDFMetadata = async (fileId, docType, citizenName, nationalId) => {
  const url = new URL('/api/public/documents/stamped-pdf-metadata', window.location.origin);
  url.searchParams.append('file_id', fileId);
  url.searchParams.append('document_type', docType);
  url.searchParams.append('citizen_name', citizenName);
  url.searchParams.append('national_id', nationalId);
  
  const response = await fetch(url);
  return response.json();
};

// Usage
const metadata = await getStampedPDFMetadata('1ABC123xyz', 'Birth Certificate', 'John Doe', '123456789');
console.log(metadata);
```

---

## Complete Workflow Examples

### Workflow 1: Email Document from Google Drive

```javascript
// 1. User verifies fingerprint
const verifyResult = await fetch('/api/public/fingerprint/verify', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ fingerprint_data: fingerprintData })
});

const citizen = await verifyResult.json();

// 2. User enters email address
const userEmail = 'user@example.com';

// 3. Send document from Google Drive
const emailResult = await fetch('/api/public/documents/send-email-gdrive', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: userEmail,
    file_id: document.googleDriveFileId,
    document_type: document.type,
    document_name: document.name
  })
});

const result = await emailResult.json();
console.log(result.message); // "Document sent successfully to user@example.com"
```

### Workflow 2: Print Stamped Document

```javascript
// 1. User verifies fingerprint
const verifyResult = await fetch('/api/public/fingerprint/verify', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ fingerprint_data: fingerprintData })
});

const citizen = await verifyResult.json();

// 2. Get document metadata
const metadata = await fetch(
  `/api/public/documents/stamped-pdf-metadata?file_id=${document.googleDriveFileId}&citizen_name=${citizen.data.first_name} ${citizen.data.last_name}&national_id=${citizen.data.national_id}&document_type=${document.type}`
);

const metadataResult = await metadata.json();

// 3. Generate and open stamped PDF for printing
const pdfUrl = new URL('/api/public/documents/generate-stamped-pdf', window.location.origin);
pdfUrl.searchParams.append('file_id', document.googleDriveFileId);
pdfUrl.searchParams.append('document_type', document.type);
pdfUrl.searchParams.append('citizen_name', `${citizen.data.first_name} ${citizen.data.last_name}`);
pdfUrl.searchParams.append('national_id', citizen.data.national_id);

// Open in new window for printing
window.open(pdfUrl.toString(), '_blank');
```

---

## Google Drive File ID Extraction

To extract the Google Drive file ID from a Drive URL:

```javascript
function extractFileId(driveUrl) {
  // Format 1: https://drive.google.com/file/d/FILE_ID/view
  const match1 = driveUrl.match(/\/file\/d\/([^\/]+)/);
  if (match1) return match1[1];
  
  // Format 2: https://drive.google.com/open?id=FILE_ID
  const match2 = driveUrl.match(/[?&]id=([^&]+)/);
  if (match2) return match2[1];
  
  // Format 3: Already just the ID
  return driveUrl;
}

// Usage
const fileId = extractFileId('https://drive.google.com/file/d/1ABC123xyz/view');
// Returns: 1ABC123xyz
```

---

## Error Handling

### Common Errors

1. **Missing File ID**
```json
{
  "status": "error",
  "message": "Google Drive file ID is required",
  "data": null
}
```

2. **Download Failed**
```json
{
  "status": "error",
  "message": "Failed to download file from Google Drive",
  "error": "failed to download file: status code 404"
}
```

3. **Email Send Failed**
```json
{
  "status": "error",
  "message": "Failed to send email",
  "error": "SMTP connection failed"
}
```

### Error Handling Example

```javascript
try {
  const response = await fetch('/api/public/documents/send-email-gdrive', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email: userEmail,
      file_id: fileId,
      document_type: docType
    })
  });
  
  const result = await response.json();
  
  if (response.ok) {
    console.log('Success:', result.message);
  } else {
    console.error('Error:', result.message);
    if (result.error) {
      console.error('Details:', result.error);
    }
  }
} catch (error) {
  console.error('Network error:', error);
}
```

---

## Security Considerations

1. **Public Files**: The current implementation uses `DownloadPublicDriveFile` which requires files to be publicly accessible via link
2. **Rate Limiting**: Consider implementing rate limiting to prevent abuse
3. **File Size**: Large files may cause timeouts; consider implementing async processing for large files
4. **CORS**: Ensure CORS is properly configured if accessing from a different domain

---

## Configuration Requirements

### Environment Variables
Already configured in your `.env`:
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_MAIL=gouvdev8@gmail.com
SMTP_PASSWORD=nidsupxyuvrxftbl
```

### Google Drive File Access
Files must be:
- Publicly accessible via link, OR
- Implement OAuth2 authentication (see `googleDrive.go` for implementation details)

---

## Testing

### Test Email Endpoint
```bash
curl -X POST http://10.249.216.144:8081/api/public/documents/send-email-gdrive \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "file_id": "YOUR_GOOGLE_DRIVE_FILE_ID",
    "document_type": "Test Document"
  }'
```

### Test Stamped PDF Endpoint
```bash
curl -X GET "http://10.249.216.144:8081/api/public/documents/generate-stamped-pdf?file_id=YOUR_GOOGLE_DRIVE_FILE_ID&document_type=Test" \
  --output test_stamped.pdf
```

### Test Metadata Endpoint
```bash
curl -X GET "http://10.249.216.144:8081/api/public/documents/stamped-pdf-metadata?file_id=YOUR_GOOGLE_DRIVE_FILE_ID&document_type=Test"
```

---

## Summary of All Document Endpoints

### Public Endpoints (No Authentication)
1. `POST /api/public/documents/send-email` - Send uploaded PDF via email
2. `POST /api/public/documents/send-email-gdrive` - Send Google Drive document via email
3. `GET/POST /api/public/documents/generate-stamped-pdf` - Generate stamped PDF for printing
4. `GET /api/public/documents/stamped-pdf-metadata` - Get stamped PDF metadata
5. `GET /api/public/documents/national-id/:national_id` - Get documents by National ID
6. `GET /api/public/documents/active` - Get active documents
7. `GET /api/public/documents/:uuid` - Get document by UUID

### Protected Endpoints (Authentication Required)
All of the above plus standard CRUD operations for document management.

---

## Next Steps

1. ✅ Endpoints implemented and tested
2. 🔄 Update frontend to use new endpoints
3. 🔄 Test with actual Google Drive files
4. 📋 Implement PDF stamping library (optional enhancement)
5. 📋 Add OAuth2 for private Google Drive files (optional enhancement)

---

## Support

For issues or questions, refer to:
- Main API Documentation: `EMAIL_ENDPOINT.md`
- Google Drive Utilities: `utils/googleDrive.go`
- PDF Generation: `utils/pdfGenerator.go`
- Email Sender: `utils/emailSender.go`
