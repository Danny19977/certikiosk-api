# CertiKiosk - Fingerprint-Based Document Certification System

## Overview

CertiKiosk is a secure document certification system that uses fingerprint authentication to retrieve, verify, and certify documents from external storage sources (Google Drive, AWS S3). The system allows users to authenticate via fingerprint, retrieve their documents, apply certification stamps, and download or print certified documents.

## System Architecture

### Workflow
1. **Citizen Registration** - Register citizens in the system with their personal information
2. **Fingerprint Enrollment** - Enroll fingerprints for each citizen
3. **Document Storage** - Store or link documents from Google Drive or AWS S3
4. **Authentication** - Verify citizen identity using fingerprint
5. **Document Retrieval** - Fetch documents from external sources
6. **Certification** - Apply certification stamp to documents
7. **Output** - Download as PDF or send to printer

## Project Structure

```
certikiosk/
├── controller/
│   ├── auth/                  # Authentication controllers
│   ├── certification/         # Document certification logic
│   ├── citizens/              # Citizen management
│   ├── documents/             # Document management
│   ├── fingerprint/           # Fingerprint operations
│   ├── Notification/          # Notification system
│   ├── user/                  # User management
│   └── userlog/               # Activity logging
├── database/
│   └── connection.go          # Database connection and migrations
├── models/
│   ├── certification.go       # Certification data model
│   ├── citizens.go            # Citizens data model
│   ├── documents.go           # Documents data model
│   ├── fingerprint.go         # Fingerprint data model
│   ├── notification.go        # Notification model
│   ├── user.go                # User model
│   └── user_logs.go           # Activity logs model
├── routes/
│   └── routes.go              # API route definitions
├── utils/
│   ├── activityLogger.go      # Activity logging utilities
│   ├── awsS3.go               # AWS S3 integration
│   ├── config.go              # Configuration management
│   ├── googleDrive.go         # Google Drive integration
│   ├── jwt.go                 # JWT token utilities
│   ├── pdfGenerator.go        # PDF generation and stamping
│   └── validateStruct.go      # Input validation
├── middlewares/
│   └── auth.go                # Authentication middleware
├── main.go                    # Application entry point
└── go.mod                     # Go module dependencies
```

## API Endpoints

### Citizens Management
- `GET /api/citizens/all` - Get all citizens
- `GET /api/citizens/all/paginate` - Get paginated citizens with search
- `GET /api/citizens/get/:uuid` - Get citizen by UUID
- `GET /api/citizens/national-id/:national_id` - Get citizen by National ID
- `POST /api/citizens/create` - Register new citizen
- `PUT /api/citizens/update/:uuid` - Update citizen information
- `DELETE /api/citizens/delete/:uuid` - Delete citizen

### Fingerprint Management
- `GET /api/fingerprint/all/paginate` - Get all fingerprints (paginated)
- `GET /api/fingerprint/citizen/:citizen_uuid` - Get fingerprint for a citizen
- `POST /api/fingerprint/enroll` - Enroll new fingerprint
- `POST /api/fingerprint/verify` - Verify fingerprint and return citizen info
- `PUT /api/fingerprint/update/:citizen_uuid` - Update fingerprint data
- `DELETE /api/fingerprint/delete/:citizen_uuid` - Delete fingerprint

### Documents Management
- `GET /api/documents/all` - Get all documents
- `GET /api/documents/all/paginate` - Get paginated documents
- `GET /api/documents/active` - Get only active documents
- `GET /api/documents/get/:uuid` - Get document by UUID
- `POST /api/documents/create` - Create/upload new document
- `POST /api/documents/fetch-external` - Fetch document from Google Drive or AWS S3
- `PUT /api/documents/update/:uuid` - Update document
- `PUT /api/documents/toggle-status/:uuid` - Activate/deactivate document
- `DELETE /api/documents/delete/:uuid` - Delete document

### Certification Management
- `GET /api/certification/all` - Get all certifications
- `GET /api/certification/all/paginate` - Get paginated certifications
- `GET /api/certification/get/:uuid` - Get certification by UUID
- `GET /api/certification/citizen/:citizen_uuid` - Get certifications for a citizen
- `GET /api/certification/document/:document_uuid` - Get certifications for a document
- `GET /api/certification/download/:uuid` - Download certified document
- `GET /api/certification/print/:uuid` - Prepare document for printing
- `POST /api/certification/certify` - Certify a document (main workflow)
- `PUT /api/certification/revoke/:uuid` - Revoke a certification
- `DELETE /api/certification/delete/:uuid` - Delete certification

## Database Schema

### Citizens Table
- UUID (Primary Key)
- NationalID (Unique, Integer)
- FirstName
- LastName
- DateOfBirth
- Email (Unique)
- CreatedAt, UpdatedAt

### Fingerprint Table
- UUID (Primary Key)
- CitizensUUID (Foreign Key)
- FingerprintData (Encrypted)
- CreatedAt, UpdatedAt

### Documents Table
- UUID (Primary Key)
- DocumentType
- DocumentDataUrl (URL to document in Google Drive/AWS)
- IssueDate
- IsActive (Boolean)
- CreatedAt, UpdatedAt

### Certification Table
- UUID (Primary Key)
- CitizensUUID (Foreign Key)
- DocumentUUID (Foreign Key)
- Aprovel (Boolean)
- CertifiedDocument (URL to certified document)
- StampDetails (JSON/Text)
- OutputFormat ("pdf" or "print")
- CreatedAt, UpdatedAt

## Main Workflow: Document Certification

### Request Example
```json
POST /api/certification/certify
{
  "citizens_uuid": "citizen-uuid-here",
  "document_uuid": "document-uuid-here",
  "fingerprint_data": "fingerprint-hash-data",
  "stamp_details": "Official Certification Stamp",
  "output_format": "pdf"
}
```

### Process Flow
1. **Verify Citizen** - Check if citizen exists in database
2. **Verify Fingerprint** - Match fingerprint data with enrolled fingerprint
3. **Verify Document** - Check if document exists and is active
4. **Apply Certification** - Add certification stamp to document
5. **Create Record** - Store certification record in database
6. **Return Result** - Return certified document URL and information

### Response Example
```json
{
  "status": "success",
  "message": "Document certified successfully",
  "data": {
    "certification": {
      "uuid": "cert-uuid",
      "citizens_uuid": "citizen-uuid",
      "document_uuid": "doc-uuid",
      "aprovel": true,
      "certified_document": "url-to-certified-doc",
      "stamp_details": "Official Certification Stamp",
      "output_format": "pdf",
      "created_at": "2025-11-04T10:30:00Z"
    },
    "citizen": { ... },
    "document": { ... }
  }
}
```

## External Integration Setup

### Google Drive Integration

1. Install required packages:
```bash
go get golang.org/x/oauth2
go get google.golang.org/api/drive/v3
```

2. Set up Google Cloud Project:
   - Create project at https://console.cloud.google.com
   - Enable Google Drive API
   - Create OAuth 2.0 credentials
   - Download credentials.json

3. Environment variables:
```env
GOOGLE_CREDENTIALS_FILE=./credentials.json
GOOGLE_TOKEN_FILE=./token.json
```

4. Uncomment implementation in `utils/googleDrive.go`

### AWS S3 Integration

1. Install AWS SDK:
```bash
go get github.com/aws/aws-sdk-go
```

2. Environment variables:
```env
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_S3_BUCKET_NAME=your-bucket-name
```

3. Uncomment implementation in `utils/awsS3.go`

### PDF Generation

1. Install PDF library (choose one):
```bash
# Option 1: gofpdf (create PDFs from scratch)
go get github.com/jung-kurt/gofpdf

# Option 2: gopdf (more features)
go get github.com/signintech/gopdf

# Option 3: pdfcpu (manipulate existing PDFs)
go get github.com/pdfcpu/pdfcpu
```

2. Uncomment and implement functions in `utils/pdfGenerator.go`

## Environment Variables

Create a `.env` file:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=certikiosk

# Server
PORT=8000

# JWT
JWT_SECRET=your-secret-key

# Google Drive
GOOGLE_CREDENTIALS_FILE=./credentials.json
GOOGLE_TOKEN_FILE=./token.json

# AWS S3
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_S3_BUCKET_NAME=certikiosk-documents
```

## Running the Application

1. Install dependencies:
```bash
go mod download
```

2. Set up database (PostgreSQL):
```bash
createdb certikiosk
```

3. Run the application:
```bash
go run main.go
```

4. The API will be available at:
```
http://localhost:8000/api
```

## Security Considerations

1. **Fingerprint Data** - Should be hashed/encrypted before storage
2. **Authentication** - All routes are protected except auth endpoints
3. **JWT Tokens** - Used for session management
4. **HTTPS** - Use HTTPS in production
5. **Database** - Use strong passwords and restrict access
6. **File Access** - Validate all file paths and URLs
7. **Input Validation** - All inputs are validated before processing

## Future Enhancements

1. **Biometric Integration** - Real fingerprint scanner hardware integration
2. **Blockchain** - Store certification hashes on blockchain for verification
3. **Multi-language** - Support for multiple languages
4. **Mobile App** - Mobile application for citizens
5. **Email Notifications** - Send certified documents via email
6. **Audit Trail** - Enhanced logging and audit capabilities
7. **OCR** - Extract text from documents for indexing
8. **Digital Signatures** - Add digital signatures to certified documents

## Testing

### Manual Testing with Postman/cURL

1. Register a citizen:
```bash
curl -X POST http://localhost:8000/api/citizens/create \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "national_id": 12345678,
    "first_name": "John",
    "last_name": "Doe",
    "date_of_birth": "1990-01-01",
    "email": "john.doe@example.com"
  }'
```

2. Enroll fingerprint:
```bash
curl -X POST http://localhost:8000/api/fingerprint/enroll \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "citizens_uuid": "citizen-uuid-from-step-1",
    "fingerprint_data": "hashed-fingerprint-data"
  }'
```

3. Certify document:
```bash
curl -X POST http://localhost:8000/api/certification/certify \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "citizens_uuid": "citizen-uuid",
    "document_uuid": "document-uuid",
    "fingerprint_data": "hashed-fingerprint-data",
    "stamp_details": "Official Stamp",
    "output_format": "pdf"
  }'
```

## Troubleshooting

### Common Issues

1. **Database connection failed**
   - Check PostgreSQL is running
   - Verify database credentials in .env
   - Ensure database exists

2. **Authentication errors**
   - Check JWT_SECRET is set
   - Verify token is valid and not expired
   - Ensure middleware is properly configured

3. **External service errors**
   - Verify API credentials are correct
   - Check network connectivity
   - Ensure required packages are installed

## License

[Your License Here]

## Contributors

[Your Name/Team]

## Support

For issues and questions, please contact [your-email@example.com]
