# Quick Reference: Google Drive Document Endpoints

## Base URL
```
http://10.249.216.144:8081
```

## Endpoints Summary

### 1. Send Email from Google Drive
```http
POST /api/public/documents/send-email-gdrive
POST /api/documents/send-email-gdrive (requires auth)
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "file_id": "1ABC123xyz",
  "document_type": "Birth Certificate",
  "document_name": "birth_certificate"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Document sent successfully to user@example.com",
  "data": { ... }
}
```

---

### 2. Generate Stamped PDF
```http
GET /api/public/documents/generate-stamped-pdf?file_id=xxx&document_type=xxx
POST /api/public/documents/generate-stamped-pdf
GET /api/documents/generate-stamped-pdf (requires auth)
POST /api/documents/generate-stamped-pdf (requires auth)
```

**Query Params (GET):**
- `file_id` (required)
- `document_type` (optional)
- `citizen_name` (optional)
- `national_id` (optional)
- `certifier_name` (optional)
- `include_stamp` (optional, default: true)

**Request Body (POST):**
```json
{
  "file_id": "1ABC123xyz",
  "document_type": "Birth Certificate",
  "citizen_name": "John Doe",
  "national_id": "123456789",
  "certifier_name": "CertiKiosk System",
  "include_stamp": true,
  "stamp_text": "CERTIFIED COPY",
  "document_name": "birth_certificate"
}
```

**Response:** PDF binary data

---

### 3. Get Stamped PDF Metadata
```http
GET /api/public/documents/stamped-pdf-metadata?file_id=xxx
GET /api/documents/stamped-pdf-metadata (requires auth)
```

**Query Params:**
- `file_id` (required)
- `document_type` (optional)
- `citizen_name` (optional)
- `national_id` (optional)

**Response:**
```json
{
  "status": "success",
  "message": "Stamped PDF metadata generated",
  "data": {
    "file_info": { ... },
    "certification": { ... },
    "stamped_pdf_url": "/api/documents/generate-stamped-pdf?file_id=xxx"
  }
}
```

---

## Quick Integration Examples

### React/Next.js

```javascript
// Email from Google Drive
const sendEmail = async (email, fileId, docType) => {
  const res = await fetch('/api/public/documents/send-email-gdrive', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, file_id: fileId, document_type: docType })
  });
  return res.json();
};

// Print Stamped PDF
const printPDF = (fileId, citizenName, nationalId) => {
  const url = `/api/public/documents/generate-stamped-pdf?file_id=${fileId}&citizen_name=${citizenName}&national_id=${nationalId}`;
  window.open(url, '_blank');
};

// Get Metadata
const getMetadata = async (fileId) => {
  const res = await fetch(`/api/public/documents/stamped-pdf-metadata?file_id=${fileId}`);
  return res.json();
};
```

---

## Common Use Cases

### 1. Kiosk Email Workflow
1. User verifies fingerprint → get citizen info
2. User enters email
3. Call `send-email-gdrive` with file_id from document record

### 2. Kiosk Print Workflow
1. User verifies fingerprint → get citizen info
2. User selects document
3. Call `generate-stamped-pdf` with citizen info
4. Open PDF in new window for printing

### 3. Preview Before Action
1. Call `stamped-pdf-metadata` to get preview info
2. Show info to user
3. Proceed with email or print based on user choice

---

## Google Drive File ID

Extract from URL:
```javascript
function getFileId(url) {
  const match = url.match(/\/file\/d\/([^\/]+)/);
  return match ? match[1] : url;
}
```

Examples:
- `https://drive.google.com/file/d/1ABC123xyz/view` → `1ABC123xyz`
- `https://drive.google.com/open?id=1ABC123xyz` → `1ABC123xyz`

---

## Error Responses

All errors follow this format:
```json
{
  "status": "error",
  "message": "Human readable error message",
  "error": "Detailed error (optional)",
  "data": null
}
```

Common HTTP status codes:
- `400` - Bad request (missing required fields)
- `404` - Resource not found
- `500` - Server error (download failed, email failed, etc.)

---

## Environment Setup

Required in `.env`:
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_MAIL=gouvdev8@gmail.com
SMTP_PASSWORD=nidsupxyuvrxftbl
```

---

## Testing

### cURL Examples

**Test Email:**
```bash
curl -X POST http://10.249.216.144:8081/api/public/documents/send-email-gdrive \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","file_id":"YOUR_FILE_ID","document_type":"Test"}'
```

**Test Print:**
```bash
curl "http://10.249.216.144:8081/api/public/documents/generate-stamped-pdf?file_id=YOUR_FILE_ID" \
  --output test.pdf
```

**Test Metadata:**
```bash
curl "http://10.249.216.144:8081/api/public/documents/stamped-pdf-metadata?file_id=YOUR_FILE_ID"
```

---

## Notes

- ✅ All endpoints available publicly (no auth) and protected (with auth)
- ✅ Google Drive files must be publicly accessible or OAuth2 configured
- ✅ PDF stamping currently returns original PDF (implement library for actual stamping)
- ✅ Email configuration already set up in your `.env`

---

## Full Documentation

See `GOOGLE_DRIVE_API_ENDPOINTS.md` for complete documentation with all details, examples, and workflows.
