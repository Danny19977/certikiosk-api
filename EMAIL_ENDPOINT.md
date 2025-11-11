# Email Document Endpoint

## Overview
The `/api/documents/send-email` endpoint allows sending documents via email. This endpoint is available both as a public route (for kiosk users) and as a protected route (for authenticated admin users).

## Endpoints

### Public Endpoint
- **URL**: `POST /api/public/documents/send-email`
- **Authentication**: None required (for kiosk users)

### Protected Endpoint
- **URL**: `POST /api/documents/send-email`
- **Authentication**: Required (Bearer token)

## Request Format

The endpoint accepts `multipart/form-data` with the following fields:

### Required Fields:
- `email` (string): Recipient's email address

### Optional Fields:
- `pdf` (file): PDF document to send as attachment
- `document_uuid` (string): UUID of document in database (if not uploading PDF directly)
- `document_type` (string): Type of document (e.g., "Birth Certificate", "ID Card")

## Examples

### Example 1: Send PDF with File Upload
```javascript
const formData = new FormData();
formData.append('email', 'user@example.com');
formData.append('pdf', pdfFile); // File object
formData.append('document_type', 'Birth Certificate');
formData.append('document_uuid', '123e4567-e89b-12d3-a456-426614174000');

fetch('http://10.249.216.144:8081/api/public/documents/send-email', {
  method: 'POST',
  body: formData
})
.then(response => response.json())
.then(data => console.log(data));
```

### Example 2: Send from React/Next.js
```javascript
const sendDocumentEmail = async (email, pdfBlob, documentType) => {
  const formData = new FormData();
  formData.append('email', email);
  formData.append('pdf', pdfBlob, 'document.pdf');
  formData.append('document_type', documentType);

  const response = await fetch('/api/public/documents/send-email', {
    method: 'POST',
    body: formData
  });

  return response.json();
};
```

## Response Format

### Success Response (200)
```json
{
  "status": "success",
  "message": "Document sent successfully to user@example.com",
  "data": {
    "email": "user@example.com",
    "document_type": "Birth Certificate",
    "document_uuid": "123e4567-e89b-12d3-a456-426614174000"
  }
}
```

### Error Responses

#### Missing Email (400)
```json
{
  "status": "error",
  "message": "Email address is required",
  "data": null
}
```

#### Missing PDF (400)
```json
{
  "status": "error",
  "message": "Either PDF file or document UUID is required",
  "data": null
}
```

#### Email Send Failure (500)
```json
{
  "status": "error",
  "message": "Failed to send email",
  "error": "SMTP connection failed"
}
```

## Email Configuration

The endpoint requires the following environment variables to be configured in `.env`:

```env
# SMTP Configuration (already configured in your .env)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_MAIL=gouvdev8@gmail.com
SMTP_PASSWORD=nidsupxyuvrxftbl
```

### Alternative Configuration Variables
The system also supports these alternative variable names for backward compatibility:
- `EMAIL_HOST` (alternative to `SMTP_HOST`)
- `EMAIL_PORT` (alternative to `SMTP_PORT`)
- `EMAIL_USERNAME` (alternative to `SMTP_MAIL`)
- `EMAIL_PASSWORD` (alternative to `SMTP_PASSWORD`)
- `EMAIL_FROM` (alternative to `SMTP_MAIL`)

## Email Template

The email sent includes:
- **Subject**: "Your [Document Type] Document from CertiKiosk"
- **HTML formatted body** with:
  - CertiKiosk branding
  - Document type and ID
  - Professional styling
  - Footer with copyright information
- **PDF attachment** with filename format: `[DocumentType]_[DocumentUUID].pdf`

## Security Considerations

1. **Rate Limiting**: Consider implementing rate limiting to prevent abuse
2. **Email Validation**: The endpoint validates email format before sending
3. **File Size Limits**: Consider adding file size limits for PDF uploads
4. **Spam Prevention**: Monitor usage patterns to prevent spam

## Testing

You can test the endpoint using curl:

```bash
curl -X POST http://10.249.216.144:8081/api/public/documents/send-email \
  -F "email=test@example.com" \
  -F "pdf=@document.pdf" \
  -F "document_type=Test Document" \
  -F "document_uuid=test-uuid-123"
```

## Troubleshooting

### Common Issues

1. **SMTP Authentication Failed**
   - Verify SMTP credentials in `.env`
   - For Gmail, ensure "Less secure app access" is enabled or use App Password

2. **Connection Timeout**
   - Check firewall settings
   - Verify SMTP_HOST and SMTP_PORT are correct

3. **Email Not Received**
   - Check spam folder
   - Verify recipient email is valid
   - Check server logs for errors

## Integration with Frontend

The frontend should already be configured to use this endpoint at:
```
http://10.249.216.144:8081/api/documents/send-email
```

Update to use the public route for kiosk users:
```
http://10.249.216.144:8081/api/public/documents/send-email
```
