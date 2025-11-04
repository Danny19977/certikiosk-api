# CertiKiosk Quick Reference Guide

## 🚀 Quick Start

### 1. Prerequisites
- Go 1.24.4 or higher
- PostgreSQL database
- Git

### 2. Installation
```bash
# Clone repository
git clone <your-repo-url>
cd certikiosk

# Install dependencies
go mod download

# Create database
createdb certikiosk

# Create .env file (see Environment Variables section)
```

### 3. Run Application
```bash
go run main.go
```

Server starts at: `http://localhost:8000`

---

## 📋 Environment Variables (.env)

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
JWT_SECRET=your-super-secret-key-change-this

# Optional: Google Drive
GOOGLE_CREDENTIALS_FILE=./credentials.json
GOOGLE_TOKEN_FILE=./token.json

# Optional: AWS S3
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_S3_BUCKET_NAME=certikiosk-docs
```

---

## 🔑 Core API Endpoints

### Authentication (Required for all endpoints)
```bash
# Register
POST /api/auth/register
{
  "fullname": "Admin User",
  "email": "admin@example.com",
  "password": "password123",
  "confirm_password": "password123"
}

# Login
POST /api/auth/login
{
  "identifier": "admin@example.com",
  "password": "password123"
}

# Response includes JWT token
# Use in headers: Authorization: Bearer <token>
```

### Main Workflow Endpoints

#### 1️⃣ Register Citizen
```bash
POST /api/citizens/create
Authorization: Bearer <token>
{
  "national_id": 12345678,
  "first_name": "John",
  "last_name": "Doe",
  "date_of_birth": "1990-01-01",
  "email": "john.doe@example.com"
}
```

#### 2️⃣ Enroll Fingerprint
```bash
POST /api/fingerprint/enroll
Authorization: Bearer <token>
{
  "citizens_uuid": "<citizen-uuid-from-step-1>",
  "fingerprint_data": "<hashed-fingerprint-data>"
}
```

#### 3️⃣ Upload/Link Document
```bash
# Option A: Direct upload
POST /api/documents/create
Authorization: Bearer <token>
{
  "document_type": "Birth Certificate",
  "document_data_url": "https://example.com/doc.pdf",
  "is_active": true
}

# Option B: From Google Drive or AWS S3
POST /api/documents/fetch-external
Authorization: Bearer <token>
{
  "source": "google_drive",
  "document_id": "drive-file-id",
  "document_type": "Birth Certificate"
}
```

#### 4️⃣ Verify Fingerprint (Kiosk Login)
```bash
POST /api/fingerprint/verify
Authorization: Bearer <token>
{
  "fingerprint_data": "<scanned-fingerprint-hash>"
}

# Returns citizen information if match found
```

#### 5️⃣ Certify Document (Main Workflow)
```bash
POST /api/certification/certify
Authorization: Bearer <token>
{
  "citizens_uuid": "<citizen-uuid>",
  "document_uuid": "<document-uuid>",
  "fingerprint_data": "<verified-fingerprint>",
  "stamp_details": "Official Certification",
  "output_format": "pdf"
}

# Returns certified document information
```

#### 6️⃣ Download Certified Document
```bash
GET /api/certification/download/<certification-uuid>
Authorization: Bearer <token>

# Returns download URL
```

#### 7️⃣ Print Certified Document
```bash
GET /api/certification/print/<certification-uuid>
Authorization: Bearer <token>

# Sends to print queue
```

---

## 📊 All API Endpoints Reference

### Citizens
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/citizens/all` | Get all citizens |
| GET | `/api/citizens/all/paginate?page=1&limit=15&search=john` | Paginated list |
| GET | `/api/citizens/get/:uuid` | Get by UUID |
| GET | `/api/citizens/national-id/:national_id` | Get by National ID |
| POST | `/api/citizens/create` | Register citizen |
| PUT | `/api/citizens/update/:uuid` | Update citizen |
| DELETE | `/api/citizens/delete/:uuid` | Delete citizen |

### Fingerprint
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/fingerprint/all/paginate` | List all fingerprints |
| GET | `/api/fingerprint/citizen/:citizen_uuid` | Get by citizen |
| POST | `/api/fingerprint/enroll` | Enroll fingerprint |
| POST | `/api/fingerprint/verify` | **Verify & authenticate** |
| PUT | `/api/fingerprint/update/:citizen_uuid` | Update fingerprint |
| DELETE | `/api/fingerprint/delete/:citizen_uuid` | Delete fingerprint |

### Documents
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/documents/all` | Get all documents |
| GET | `/api/documents/all/paginate` | Paginated list |
| GET | `/api/documents/active` | Active documents only |
| GET | `/api/documents/get/:uuid` | Get by UUID |
| POST | `/api/documents/create` | Create document |
| POST | `/api/documents/fetch-external` | **Fetch from Drive/S3** |
| PUT | `/api/documents/update/:uuid` | Update document |
| PUT | `/api/documents/toggle-status/:uuid` | Toggle active status |
| DELETE | `/api/documents/delete/:uuid` | Delete document |

### Certification
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/certification/all` | Get all certifications |
| GET | `/api/certification/all/paginate` | Paginated list |
| GET | `/api/certification/get/:uuid` | Get by UUID |
| GET | `/api/certification/citizen/:citizen_uuid` | By citizen |
| GET | `/api/certification/document/:document_uuid` | By document |
| GET | `/api/certification/download/:uuid` | **Download PDF** |
| GET | `/api/certification/print/:uuid` | **Print document** |
| POST | `/api/certification/certify` | **Main: Certify document** |
| PUT | `/api/certification/revoke/:uuid` | Revoke certification |
| DELETE | `/api/certification/delete/:uuid` | Delete certification |

---

## 🧪 Testing with cURL

### Complete Test Flow

```bash
# 1. Register admin user
curl -X POST http://localhost:8000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "fullname": "Admin User",
    "email": "admin@test.com",
    "password": "password123",
    "confirm_password": "password123",
    "role": "admin"
  }'

# 2. Login and get token
curl -X POST http://localhost:8000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "admin@test.com",
    "password": "password123"
  }'

# Save the token from response
TOKEN="your-jwt-token-here"

# 3. Create citizen
curl -X POST http://localhost:8000/api/citizens/create \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "national_id": 12345678,
    "first_name": "John",
    "last_name": "Doe",
    "date_of_birth": "1990-01-01",
    "email": "john.doe@test.com"
  }'

# Save citizen UUID
CITIZEN_UUID="uuid-from-response"

# 4. Enroll fingerprint
curl -X POST http://localhost:8000/api/fingerprint/enroll \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "citizens_uuid": "'$CITIZEN_UUID'",
    "fingerprint_data": "abc123fingerprinthash"
  }'

# 5. Create document
curl -X POST http://localhost:8000/api/documents/create \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "document_type": "Birth Certificate",
    "document_data_url": "https://example.com/doc.pdf",
    "is_active": true
  }'

# Save document UUID
DOCUMENT_UUID="uuid-from-response"

# 6. Verify fingerprint
curl -X POST http://localhost:8000/api/fingerprint/verify \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "fingerprint_data": "abc123fingerprinthash"
  }'

# 7. Certify document
curl -X POST http://localhost:8000/api/certification/certify \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "citizens_uuid": "'$CITIZEN_UUID'",
    "document_uuid": "'$DOCUMENT_UUID'",
    "fingerprint_data": "abc123fingerprinthash",
    "stamp_details": "Official Certification",
    "output_format": "pdf"
  }'

# Save certification UUID
CERT_UUID="uuid-from-response"

# 8. Download certified document
curl -X GET http://localhost:8000/api/certification/download/$CERT_UUID \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🛠️ Optional Integrations

### Google Drive Setup

```bash
# 1. Install packages
go get golang.org/x/oauth2
go get google.golang.org/api/drive/v3

# 2. Get credentials from Google Cloud Console
# - Create project
# - Enable Drive API
# - Create OAuth 2.0 credentials
# - Download credentials.json

# 3. Uncomment code in utils/googleDrive.go

# 4. Run first time to authenticate
go run main.go
# Follow the OAuth flow
```

### AWS S3 Setup

```bash
# 1. Install SDK
go get github.com/aws/aws-sdk-go

# 2. Set environment variables
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=your-key
export AWS_SECRET_ACCESS_KEY=your-secret
export AWS_S3_BUCKET_NAME=your-bucket

# 3. Uncomment code in utils/awsS3.go
```

### PDF Generation Setup

```bash
# Choose one library:

# Option 1: gofpdf (simple)
go get github.com/jung-kurt/gofpdf

# Option 2: pdfcpu (advanced)
go get github.com/pdfcpu/pdfcpu

# Option 3: gopdf
go get github.com/signintech/gopdf

# Implement functions in utils/pdfGenerator.go
```

---

## 📁 Project Structure Quick Reference

```
certikiosk/
├── controller/          # Business logic
│   ├── citizens/        ← Citizen CRUD
│   ├── fingerprint/     ← Fingerprint ops
│   ├── documents/       ← Document management
│   └── certification/   ← Main workflow
├── models/              # Data models
├── routes/              # API routes
├── utils/               # Helper functions
│   ├── googleDrive.go   ← Google Drive API
│   ├── awsS3.go         ← AWS S3 API
│   └── pdfGenerator.go  ← PDF operations
├── database/            # DB connection
├── middlewares/         # Auth, CORS, etc.
└── main.go             # Entry point
```

---

## 🐛 Troubleshooting

### Database Connection Error
```bash
# Check PostgreSQL is running
pg_isready

# Create database if missing
createdb certikiosk

# Verify credentials in .env
```

### JWT Authentication Error
```bash
# Ensure JWT_SECRET is set in .env
# Check token format: "Bearer <token>"
# Verify token hasn't expired
```

### Fingerprint Not Recognized
```bash
# Ensure fingerprint_data matches exactly
# Check fingerprint exists for citizen
# Verify citizen UUID is correct
```

---

## 📝 Common Response Formats

### Success Response
```json
{
  "status": "success",
  "message": "Operation completed successfully",
  "data": { ... }
}
```

### Error Response
```json
{
  "status": "error",
  "message": "Error description",
  "error": "Detailed error message"
}
```

### Paginated Response
```json
{
  "status": "success",
  "message": "Data retrieved successfully",
  "data": [...],
  "pagination": {
    "total_records": 100,
    "total_pages": 10,
    "current_page": 1,
    "page_size": 15
  }
}
```

---

## 🔒 Security Checklist

- [ ] Change default JWT_SECRET
- [ ] Use HTTPS in production
- [ ] Encrypt fingerprint data
- [ ] Use strong database passwords
- [ ] Enable rate limiting
- [ ] Implement IP whitelisting for kiosks
- [ ] Regular security audits
- [ ] Keep dependencies updated

---

## 📚 Additional Resources

- **Full Documentation**: `SYSTEM_DOCUMENTATION.md`
- **Implementation Details**: `IMPLEMENTATION_SUMMARY.md`
- **Workflow Diagrams**: `WORKFLOW_DIAGRAMS.md`

---

## ⚡ Performance Tips

1. **Database Indexing**: Add indexes on frequently queried fields
2. **Caching**: Implement Redis for session caching
3. **Connection Pooling**: Configure PostgreSQL connection pool
4. **Load Balancing**: Use nginx for multiple instances
5. **CDN**: Use CDN for certified documents

---

## 🆘 Support

For issues, check:
1. Server logs
2. Database logs
3. API error responses
4. Documentation files

---

**Last Updated**: November 4, 2025  
**Version**: 1.0.0
