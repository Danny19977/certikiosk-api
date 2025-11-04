# CertiKiosk System Workflow Diagrams

## Main Kiosk Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│                    CERTIKIOSK SYSTEM FLOW                        │
└─────────────────────────────────────────────────────────────────┘

Step 1: CITIZEN APPROACHES KIOSK
    │
    ├─► Place finger on scanner
    │
    ├─► POST /api/fingerprint/verify
    │   └─► Request: { "fingerprint_data": "..." }
    │
    ├─► System verifies fingerprint
    │   ├─► Search fingerprint in database
    │   ├─► Match found? 
    │   │   ├─► YES → Continue
    │   │   └─► NO → Show error "Not recognized"
    │   │
    │   └─► Return citizen information
    │
    └─► Display: "Welcome, [Citizen Name]"


Step 2: DOCUMENT SELECTION
    │
    ├─► GET /api/documents/active (for this citizen)
    │
    ├─► Display available documents:
    │   ├─ Birth Certificate
    │   ├─ ID Card
    │   ├─ Diploma
    │   └─ etc.
    │
    └─► User selects document


Step 3: CERTIFICATION PROCESS
    │
    ├─► POST /api/certification/certify
    │   └─► Request: {
    │         "citizens_uuid": "...",
    │         "document_uuid": "...",
    │         "fingerprint_data": "...",
    │         "stamp_details": "Official Stamp",
    │         "output_format": "pdf"
    │       }
    │
    ├─► Backend Process:
    │   │
    │   ├─► 1. Verify citizen exists ✓
    │   │
    │   ├─► 2. Verify fingerprint matches ✓
    │   │
    │   ├─► 3. Verify document exists & is active ✓
    │   │
    │   ├─► 4. Retrieve document from source
    │   │   ├─► Google Drive?
    │   │   └─► AWS S3?
    │   │
    │   ├─► 5. Apply certification stamp
    │   │   ├─► Add official stamp
    │   │   ├─► Add date/time
    │   │   ├─► Add certifier signature
    │   │   └─► Add QR code (optional)
    │   │
    │   ├─► 6. Save certified document
    │   │
    │   └─► 7. Create certification record
    │
    └─► Return certified document info


Step 4: OUTPUT SELECTION
    │
    ├─► Display options:
    │   ├─ [ Download PDF ]
    │   └─ [ Print Now ]
    │
    ├─► Option A: DOWNLOAD
    │   │
    │   ├─► GET /api/certification/download/:uuid
    │   │
    │   └─► Download PDF to local/USB drive
    │
    └─► Option B: PRINT
        │
        ├─► GET /api/certification/print/:uuid
        │
        └─► Send to printer queue


Step 5: COMPLETION
    │
    ├─► Show success message
    │
    ├─► Log activity
    │
    └─► Return to start screen
```

## Data Flow Architecture

```
┌────────────────────────────────────────────────────────────────┐
│                    SYSTEM ARCHITECTURE                          │
└────────────────────────────────────────────────────────────────┘

┌─────────────┐         ┌──────────────┐         ┌──────────────┐
│   KIOSK     │────────►│  API SERVER  │────────►│  DATABASE    │
│  Terminal   │  HTTPS  │   (Fiber)    │   SQL   │ (PostgreSQL) │
└─────────────┘         └──────────────┘         └──────────────┘
      │                        │
      │                        │
      ▼                        ▼
┌─────────────┐         ┌──────────────┐
│ Fingerprint │         │  External    │
│   Scanner   │         │   Storage    │
└─────────────┘         └──────────────┘
                               │
                        ┌──────┴──────┐
                        │             │
                   ┌────▼────┐   ┌───▼────┐
                   │  Google │   │  AWS   │
                   │  Drive  │   │   S3   │
                   └─────────┘   └────────┘
```

## Database Relationships

```
┌──────────────────────────────────────────────────────────────┐
│                  DATABASE SCHEMA RELATIONSHIPS                │
└──────────────────────────────────────────────────────────────┘

┌─────────────────┐
│    CITIZENS     │
│─────────────────│
│ UUID (PK)       │◄──────┐
│ NationalID      │       │
│ FirstName       │       │
│ LastName        │       │
│ Email           │       │
└─────────────────┘       │
                          │
                          │ Foreign Key
                          │
                  ┌───────┴──────────┐
                  │                  │
         ┌────────▼────────┐  ┌──────▼──────────┐
         │  FINGERPRINT    │  │  CERTIFICATION  │
         │─────────────────│  │─────────────────│
         │ UUID (PK)       │  │ UUID (PK)       │
         │ CitizensUUID(FK)│  │ CitizensUUID(FK)│◄─┐
         │ FingerprintData │  │ DocumentUUID(FK)│  │
         └─────────────────┘  │ Aprovel         │  │
                              │ CertifiedDoc    │  │
                              │ StampDetails    │  │
                              │ OutputFormat    │  │
                              └─────────────────┘  │
                                       │           │
                                       │           │
                              ┌────────▼───────────┴──┐
                              │     DOCUMENTS         │
                              │───────────────────────│
                              │ UUID (PK)             │
                              │ DocumentType          │
                              │ DocumentDataUrl       │
                              │ IssueDate             │
                              │ IsActive              │
                              └───────────────────────┘
```

## Authentication Flow

```
┌──────────────────────────────────────────────────────────────┐
│              FINGERPRINT AUTHENTICATION FLOW                  │
└──────────────────────────────────────────────────────────────┘

User Scan
    │
    ├─► Capture fingerprint image
    │
    ├─► Convert to hash/template
    │
    ├─► Send to backend
    │       POST /api/fingerprint/verify
    │       { "fingerprint_data": "hash..." }
    │
    ├─► Backend Processing:
    │   │
    │   ├─► Query database for matching fingerprint
    │   │       SELECT * FROM fingerprint 
    │   │       WHERE fingerprint_data = 'hash...'
    │   │
    │   ├─► Found?
    │   │   │
    │   │   ├─► YES
    │   │   │   │
    │   │   │   ├─► Get CitizensUUID
    │   │   │   │
    │   │   │   ├─► Query citizen details
    │   │   │   │       SELECT * FROM citizens
    │   │   │   │       WHERE uuid = CitizensUUID
    │   │   │   │
    │   │   │   └─► Return citizen data
    │   │   │       Status: 200 OK
    │   │   │       Data: { citizen: {...}, fingerprint: {...} }
    │   │   │
    │   │   └─► NO
    │   │       │
    │   │       └─► Return error
    │   │           Status: 404 Not Found
    │   │           Message: "Fingerprint not recognized"
    │   │
    │   └─► End
    │
    └─► Display result on kiosk
```

## Document Certification Process

```
┌──────────────────────────────────────────────────────────────┐
│           DOCUMENT CERTIFICATION DETAILED FLOW                │
└──────────────────────────────────────────────────────────────┘

Input:
  • Citizen UUID
  • Document UUID
  • Fingerprint Data
  • Stamp Details
  • Output Format

    │
    ├─► STEP 1: VERIFY CITIZEN
    │   │
    │   ├─► SELECT * FROM citizens WHERE uuid = ?
    │   │
    │   ├─► Exists?
    │   │   ├─► YES → Continue
    │   │   └─► NO → Return 404 "Citizen not found"
    │   │
    │   └─► Citizen Data Retrieved ✓
    │
    ├─► STEP 2: VERIFY FINGERPRINT
    │   │
    │   ├─► SELECT * FROM fingerprint 
    │   │   WHERE citizens_uuid = ? 
    │   │   AND fingerprint_data = ?
    │   │
    │   ├─► Match?
    │   │   ├─► YES → Continue
    │   │   └─► NO → Return 401 "Fingerprint verification failed"
    │   │
    │   └─► Identity Confirmed ✓
    │
    ├─► STEP 3: VERIFY DOCUMENT
    │   │
    │   ├─► SELECT * FROM documents WHERE uuid = ?
    │   │
    │   ├─► Exists?
    │   │   ├─► YES → Check if active
    │   │   │   ├─► Active? → Continue
    │   │   │   └─► Not Active? → Return 400 "Document not active"
    │   │   │
    │   │   └─► NO → Return 404 "Document not found"
    │   │
    │   └─► Document Retrieved ✓
    │
    ├─► STEP 4: RETRIEVE DOCUMENT CONTENT
    │   │
    │   ├─► Check DocumentDataUrl source
    │   │   │
    │   │   ├─► Google Drive?
    │   │   │   └─► Call utils.DownloadFileFromDrive(fileID)
    │   │   │
    │   │   ├─► AWS S3?
    │   │   │   └─► Call utils.DownloadFileFromS3(key)
    │   │   │
    │   │   └─► Direct URL?
    │   │       └─► HTTP GET request
    │   │
    │   └─► Document Content Retrieved ✓
    │
    ├─► STEP 5: APPLY CERTIFICATION STAMP
    │   │
    │   ├─► Create stamp configuration
    │   │   ├─ Citizen Name
    │   │   ├─ Date/Time
    │   │   ├─ Certifier Signature
    │   │   ├─ Official Seal
    │   │   └─ QR Code (optional)
    │   │
    │   ├─► Call utils.AddStampToPDF(document, stampConfig)
    │   │
    │   ├─► Generate certified PDF
    │   │
    │   └─► Certified Document Created ✓
    │
    ├─► STEP 6: SAVE CERTIFICATION RECORD
    │   │
    │   ├─► Generate Certification UUID
    │   │
    │   ├─► INSERT INTO certification (
    │   │       uuid, citizens_uuid, document_uuid,
    │   │       aprovel, certified_document,
    │   │       stamp_details, output_format
    │   │   )
    │   │
    │   └─► Certification Saved ✓
    │
    └─► STEP 7: RETURN RESULT
        │
        └─► Return JSON:
            {
              "status": "success",
              "message": "Document certified successfully",
              "data": {
                "certification": {...},
                "citizen": {...},
                "document": {...}
              }
            }
```

## Error Handling Flow

```
┌──────────────────────────────────────────────────────────────┐
│                    ERROR HANDLING                             │
└──────────────────────────────────────────────────────────────┘

Every Request
    │
    ├─► Authentication Check
    │   │
    │   ├─► Valid JWT Token?
    │   │   ├─► YES → Continue
    │   │   └─► NO → Return 401 Unauthorized
    │   │
    │   └─► Token Expired?
    │       ├─► YES → Return 401 "Token expired"
    │       └─► NO → Continue
    │
    ├─► Input Validation
    │   │
    │   ├─► Required fields present?
    │   │   ├─► YES → Continue
    │   │   └─► NO → Return 400 "Missing required fields"
    │   │
    │   ├─► Data format valid?
    │   │   ├─► YES → Continue
    │   │   └─► NO → Return 400 "Invalid data format"
    │   │
    │   └─► All checks passed ✓
    │
    ├─► Business Logic Execution
    │   │
    │   ├─► Database Error?
    │   │   └─► Return 500 "Database error: [details]"
    │   │
    │   ├─► Resource Not Found?
    │   │   └─► Return 404 "Resource not found"
    │   │
    │   ├─► Conflict?
    │   │   └─► Return 409 "Resource already exists"
    │   │
    │   └─► Success ✓
    │
    └─► Return Success Response
        └─► Status 200/201 with data
```

## Complete System Integration

```
┌──────────────────────────────────────────────────────────────────┐
│              COMPLETE SYSTEM INTEGRATION DIAGRAM                  │
└──────────────────────────────────────────────────────────────────┘

                        ┌─────────────┐
                        │   ADMIN     │
                        │   PORTAL    │
                        └──────┬──────┘
                               │
                               │ Manage citizens,
                               │ documents, users
                               │
        ┌──────────────────────▼────────────────────────┐
        │            API SERVER (Port 8000)              │
        │  ┌──────────────────────────────────────────┐ │
        │  │         CONTROLLERS                      │ │
        │  │  ┌────────┐ ┌──────────┐ ┌───────────┐ │ │
        │  │  │Citizens│ │Fingerprint│ │ Documents │ │ │
        │  │  └────────┘ └──────────┘ └───────────┘ │ │
        │  │  ┌──────────────┐ ┌─────────┐          │ │
        │  │  │Certification │ │  Auth   │          │ │
        │  │  └──────────────┘ └─────────┘          │ │
        │  └──────────────────────────────────────────┘ │
        │  ┌──────────────────────────────────────────┐ │
        │  │         MIDDLEWARE                       │ │
        │  │  • Authentication (JWT)                  │ │
        │  │  • CORS                                  │ │
        │  │  • Logger                                │ │
        │  └──────────────────────────────────────────┘ │
        │  ┌──────────────────────────────────────────┐ │
        │  │         UTILITIES                        │ │
        │  │  • PDF Generator                         │ │
        │  │  • Google Drive API                      │ │
        │  │  • AWS S3 API                            │ │
        │  │  • Activity Logger                       │ │
        │  └──────────────────────────────────────────┘ │
        └───────────┬──────────────────────┬────────────┘
                    │                      │
                    │                      │
        ┌───────────▼──────────┐  ┌────────▼───────────┐
        │   PostgreSQL DB      │  │ External Storage   │
        │  ┌─────────────────┐ │  │ ┌────────────────┐ │
        │  │ • citizens      │ │  │ │ Google Drive   │ │
        │  │ • fingerprint   │ │  │ │ AWS S3         │ │
        │  │ • documents     │ │  │ └────────────────┘ │
        │  │ • certification │ │  └────────────────────┘
        │  │ • users         │ │
        │  │ • user_logs     │ │
        │  └─────────────────┘ │
        └──────────────────────┘
                    ▲
                    │
        ┌───────────┴──────────┐
        │    KIOSK TERMINAL    │
        │  ┌─────────────────┐ │
        │  │ Touch Screen    │ │
        │  │ Fingerprint     │ │
        │  │ Scanner         │ │
        │  │ Printer         │ │
        │  └─────────────────┘ │
        └──────────────────────┘
```

This completes the visual representation of the entire CertiKiosk system!
