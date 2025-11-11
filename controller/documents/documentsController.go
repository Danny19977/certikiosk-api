package documents

import (
	"strconv"
	"time"

	"github.com/Danny19977/certikiosk.git/database"
	"github.com/Danny19977/certikiosk.git/models"
	"github.com/Danny19977/certikiosk.git/utils"
	"github.com/gofiber/fiber/v2"
)

// GetPaginatedDocuments - Get paginated list of documents
func GetPaginatedDocuments(c *fiber.Ctx) error {
	db := database.DB

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit", "15"))
	if err != nil || limit <= 0 {
		limit = 15
	}
	offset := (page - 1) * limit

	search := c.Query("search", "")

	var documents []models.Documents
	var totalRecords int64

	query := db.Model(&models.Documents{})
	if search != "" {
		query = query.Where("document_type ILIKE ?", "%"+search+"%")
	}
	query.Count(&totalRecords)

	err = query.Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&documents).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch documents",
			"error":   err.Error(),
		})
	}

	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))

	pagination := map[string]interface{}{
		"total_records": totalRecords,
		"total_pages":   totalPages,
		"current_page":  page,
		"page_size":     limit,
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "Documents retrieved successfully",
		"data":       documents,
		"pagination": pagination,
	})
}

// GetAllDocuments - Get all documents
func GetAllDocuments(c *fiber.Ctx) error {
	db := database.DB
	var documents []models.Documents

	if err := db.Find(&documents).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch documents",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All documents retrieved successfully",
		"data":    documents,
	})
}

// GetDocument - Get a single document by UUID
func GetDocument(c *fiber.Ctx) error {
	documentUUID := c.Params("uuid")
	db := database.DB
	var document models.Documents

	if err := db.Where("uuid = ?", documentUUID).First(&document).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Document not found",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Document found",
		"data":    document,
	})
}

// GetDocumentsByNationalID - Get all documents for a specific citizen by National ID
func GetDocumentsByNationalID(c *fiber.Ctx) error {
	nationalID := c.Params("national_id")
	db := database.DB
	var documents []models.Documents

	if err := db.Where("national_id = ?", nationalID).Order("created_at DESC").Find(&documents).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch documents",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Documents retrieved successfully",
		"data":    documents,
	})
}

// GetDocumentsByUserUUID - Get all documents for a specific user by User UUID
func GetDocumentsByUserUUID(c *fiber.Ctx) error {
	userUUID := c.Params("user_uuid")
	db := database.DB
	var documents []models.Documents

	if err := db.Where("user_uuid = ?", userUUID).Order("created_at DESC").Find(&documents).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch documents",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Documents retrieved successfully",
		"data":    documents,
	})
}

// GetActiveDocuments - Get all active documents
func GetActiveDocuments(c *fiber.Ctx) error {
	db := database.DB
	var documents []models.Documents

	if err := db.Where("is_active = ?", true).Find(&documents).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch active documents",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Active documents retrieved successfully",
		"data":    documents,
	})
}

// CreateDocument - Upload/Register a new document
func CreateDocument(c *fiber.Ctx) error {
	type DocumentInput struct {
		NationalID      int    `json:"national_id"`
		UserUUID        string `json:"user_uuid"`
		DocumentType    string `json:"document_type"`
		DocumentDataUrl string `json:"document_data_url"`
		IssueDate       string `json:"issue_date"`
		IsActive        bool   `json:"is_active"`
	}

	var input DocumentInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
	}

	// Validate required fields
	if input.DocumentType == "" || input.DocumentDataUrl == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Document type and data URL are required",
			"data":    nil,
		})
	}

	if input.NationalID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "National ID is required",
			"data":    nil,
		})
	}

	// Parse issue date
	issueDate := time.Now()
	if input.IssueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", input.IssueDate)
		if err == nil {
			issueDate = parsedDate
		}
	}

	document := models.Documents{
		UUID:            utils.GenerateUUID(),
		NationalID:      input.NationalID,
		UserUUID:        input.UserUUID,
		DocumentType:    input.DocumentType,
		DocumentDataUrl: input.DocumentDataUrl,
		IssueDate:       issueDate,
		IsActive:        input.IsActive,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := database.DB.Create(&document).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create document",
			"error":   err.Error(),
		})
	}

	// Log document creation
	utils.LogCreateWithDB(database.DB, c, "document", input.DocumentType, document.UUID)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Document created successfully",
		"data":    document,
	})
}

// FetchDocumentFromExternalSource - Retrieve document from Google Drive or AWS
func FetchDocumentFromExternalSource(c *fiber.Ctx) error {
	type FetchDocumentInput struct {
		NationalID   int    `json:"national_id"`
		UserUUID     string `json:"user_uuid"`
		Source       string `json:"source"`        // "google_drive" or "aws_s3"
		DocumentID   string `json:"document_id"`   // ID/Key in external source
		DocumentType string `json:"document_type"` // Type of document
	}

	var input FetchDocumentInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
	}

	if input.Source == "" || input.DocumentID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Source and document ID are required",
			"data":    nil,
		})
	}

	if input.NationalID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "National ID is required",
			"data":    nil,
		})
	}

	// TODO: Implement actual external source retrieval
	// For now, this is a placeholder structure
	var documentUrl string

	switch input.Source {
	case "google_drive":
		// Call Google Drive API utility function
		documentUrl = "https://drive.google.com/file/d/" + input.DocumentID
	case "aws_s3":
		// Call AWS S3 utility function
		documentUrl = "https://s3.amazonaws.com/bucket/" + input.DocumentID
	default:
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid source. Use 'google_drive' or 'aws_s3'",
			"data":    nil,
		})
	}

	// Create document record
	document := models.Documents{
		UUID:            utils.GenerateUUID(),
		NationalID:      input.NationalID,
		UserUUID:        input.UserUUID,
		DocumentType:    input.DocumentType,
		DocumentDataUrl: documentUrl,
		IssueDate:       time.Now(),
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := database.DB.Create(&document).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to save document",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Document fetched and saved successfully",
		"data":    document,
	})
}

// UpdateDocument - Update document information
func UpdateDocument(c *fiber.Ctx) error {
	documentUUID := c.Params("uuid")
	db := database.DB

	type UpdateDocumentInput struct {
		NationalID      int    `json:"national_id"`
		UserUUID        string `json:"user_uuid"`
		DocumentType    string `json:"document_type"`
		DocumentDataUrl string `json:"document_data_url"`
		IsActive        *bool  `json:"is_active"`
	}

	var updateData UpdateDocumentInput

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
	}

	var document models.Documents
	if err := db.Where("uuid = ?", documentUUID).First(&document).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Document not found",
			"data":    nil,
		})
	}

	// Update fields
	if updateData.NationalID != 0 {
		document.NationalID = updateData.NationalID
	}
	if updateData.UserUUID != "" {
		document.UserUUID = updateData.UserUUID
	}
	if updateData.DocumentType != "" {
		document.DocumentType = updateData.DocumentType
	}
	if updateData.DocumentDataUrl != "" {
		document.DocumentDataUrl = updateData.DocumentDataUrl
	}
	if updateData.IsActive != nil {
		document.IsActive = *updateData.IsActive
	}

	document.UpdatedAt = time.Now()

	if err := db.Save(&document).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update document",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Document updated successfully",
		"data":    document,
	})
}

// DeleteDocument - Delete a document
func DeleteDocument(c *fiber.Ctx) error {
	documentUUID := c.Params("uuid")
	db := database.DB

	var document models.Documents
	if err := db.Where("uuid = ?", documentUUID).First(&document).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Document not found",
			"data":    nil,
		})
	}

	if err := db.Delete(&document).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete document",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Document deleted successfully",
		"data":    nil,
	})
}

// ToggleDocumentStatus - Activate/Deactivate a document
func ToggleDocumentStatus(c *fiber.Ctx) error {
	documentUUID := c.Params("uuid")
	db := database.DB

	var document models.Documents
	if err := db.Where("uuid = ?", documentUUID).First(&document).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Document not found",
			"data":    nil,
		})
	}

	document.IsActive = !document.IsActive
	document.UpdatedAt = time.Now()

	if err := db.Save(&document).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to toggle document status",
			"error":   err.Error(),
		})
	}

	status := "deactivated"
	if document.IsActive {
		status = "activated"
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Document " + status + " successfully",
		"data":    document,
	})
}

// SendDocumentEmail - Send document via email
func SendDocumentEmail(c *fiber.Ctx) error {
	type EmailInput struct {
		Email        string `json:"email"`
		DocumentUUID string `json:"document_uuid"`
		DocumentType string `json:"document_type"`
	}

	var input EmailInput

	// Parse multipart form data
	if err := c.BodyParser(&input); err != nil {
		// Try parsing as multipart form
		input.Email = c.FormValue("email")
		input.DocumentUUID = c.FormValue("document_uuid")
		input.DocumentType = c.FormValue("document_type")

		if input.Email == "" {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid input data",
				"error":   err.Error(),
			})
		}
	}

	// Validate email
	if input.Email == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Email address is required",
			"data":    nil,
		})
	}

	// Get PDF file from form
	file, err := c.FormFile("pdf")
	var pdfData []byte

	if err == nil && file != nil {
		// Read uploaded PDF file
		fileHandle, err := file.Open()
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to read PDF file",
				"error":   err.Error(),
			})
		}
		defer fileHandle.Close()

		pdfData = make([]byte, file.Size)
		_, err = fileHandle.Read(pdfData)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to read PDF content",
				"error":   err.Error(),
			})
		}
	} else if input.DocumentUUID != "" {
		// Fetch document from database if UUID provided
		db := database.DB
		var document models.Documents

		if err := db.Where("uuid = ?", input.DocumentUUID).First(&document).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Document not found",
				"data":    nil,
			})
		}

		input.DocumentType = document.DocumentType

		// TODO: Fetch actual PDF data from DocumentDataUrl
		// For now, return an error indicating PDF needs to be uploaded
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Please upload the PDF file directly",
			"data":    nil,
		})
	} else {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Either PDF file or document UUID is required",
			"data":    nil,
		})
	}

	// Set default document type if not provided
	if input.DocumentType == "" {
		input.DocumentType = "Document"
	}

	// Send email with PDF attachment
	err = utils.SendDocumentEmail(input.Email, input.DocumentType, input.DocumentUUID, pdfData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to send email",
			"error":   err.Error(),
		})
	}

	// Log email sent activity
	utils.LogCreateWithDB(database.DB, c, "document_email", "Document sent to "+input.Email, input.DocumentUUID)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Document sent successfully to " + input.Email,
		"data": fiber.Map{
			"email":         input.Email,
			"document_type": input.DocumentType,
			"document_uuid": input.DocumentUUID,
		},
	})
}

// SendDocumentEmailFromGDrive - Send document via email from Google Drive
func SendDocumentEmailFromGDrive(c *fiber.Ctx) error {
	type EmailGDriveInput struct {
		Email        string `json:"email"`
		FileID       string `json:"file_id"`
		DocumentType string `json:"document_type"`
		DocumentName string `json:"document_name"`
	}

	var input EmailGDriveInput

	// Parse JSON or form data
	if err := c.BodyParser(&input); err != nil {
		// Try parsing as form data
		input.Email = c.FormValue("email")
		input.FileID = c.FormValue("file_id")
		input.DocumentType = c.FormValue("document_type")
		input.DocumentName = c.FormValue("document_name")

		if input.Email == "" {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid input data",
				"error":   err.Error(),
			})
		}
	}

	// Validate required fields
	if input.Email == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Email address is required",
			"data":    nil,
		})
	}

	if input.FileID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Google Drive file ID is required",
			"data":    nil,
		})
	}

	// Set defaults
	if input.DocumentType == "" {
		input.DocumentType = "Document"
	}
	if input.DocumentName == "" {
		input.DocumentName = "document"
	}

	// Download file from Google Drive
	pdfData, err := utils.DownloadPublicDriveFile(input.FileID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to download file from Google Drive",
			"error":   err.Error(),
		})
	}

	// Send email with PDF attachment
	err = utils.SendDocumentEmail(input.Email, input.DocumentType, input.FileID, pdfData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to send email",
			"error":   err.Error(),
		})
	}

	// Log email sent activity
	utils.LogCreateWithDB(database.DB, c, "document_email_gdrive", "Document from GDrive sent to "+input.Email, input.FileID)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Document sent successfully to " + input.Email,
		"data": fiber.Map{
			"email":         input.Email,
			"document_type": input.DocumentType,
			"file_id":       input.FileID,
			"document_name": input.DocumentName,
		},
	})
}

// GenerateStampedPDF - Generate a stamped/certified PDF for printing
func GenerateStampedPDF(c *fiber.Ctx) error {
	type StampedPDFInput struct {
		FileID        string `json:"file_id"`
		DocumentType  string `json:"document_type"`
		CitizenName   string `json:"citizen_name"`
		NationalID    string `json:"national_id"`
		CertifierName string `json:"certifier_name"`
		IncludeStamp  bool   `json:"include_stamp"`
		StampText     string `json:"stamp_text"`
		DocumentName  string `json:"document_name"`
	}

	var input StampedPDFInput

	// Parse JSON or query parameters
	if err := c.BodyParser(&input); err != nil {
		// Try parsing from query params
		input.FileID = c.Query("file_id")
		input.DocumentType = c.Query("document_type")
		input.CitizenName = c.Query("citizen_name")
		input.NationalID = c.Query("national_id")
		input.CertifierName = c.Query("certifier_name", "CertiKiosk System")
		input.IncludeStamp = c.Query("include_stamp", "true") == "true"
		input.StampText = c.Query("stamp_text")
		input.DocumentName = c.Query("document_name", "document")
	}

	// Validate required fields
	if input.FileID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Google Drive file ID is required",
			"data":    nil,
		})
	}

	// Set defaults
	if input.DocumentType == "" {
		input.DocumentType = "Document"
	}
	if input.CertifierName == "" {
		input.CertifierName = "CertiKiosk System"
	}
	if input.DocumentName == "" {
		input.DocumentName = "document"
	}

	// Download original file from Google Drive
	pdfData, err := utils.DownloadPublicDriveFile(input.FileID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to download file from Google Drive",
			"error":   err.Error(),
		})
	}

	// Prepare certification info
	certInfo := utils.CertificationInfo{
		CitizenName:   input.CitizenName,
		NationalID:    input.NationalID,
		DocumentType:  input.DocumentType,
		CertifiedDate: time.Now(),
		CertifierName: input.CertifierName,
		StampDetails:  input.StampText,
	}

	if input.StampText == "" {
		certInfo.StampDetails = utils.GetCertificationStampTemplate(input.CertifierName, time.Now())
	}

	// For now, return the original PDF with metadata
	// In production, you would add the stamp to the PDF here
	// using a PDF library like pdfcpu or gofpdf

	// Log the generation activity
	utils.LogCreateWithDB(database.DB, c, "stamped_pdf_generate", "Stamped PDF generated for "+input.CitizenName, input.FileID)

	// Set response headers for PDF download
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "inline; filename=\""+input.DocumentName+"_stamped.pdf\"")

	// Return the PDF data directly
	// In production, this would be the stamped PDF
	return c.Send(pdfData)
}

// GenerateStampedPDFMetadata - Get metadata for stamped PDF without downloading
func GenerateStampedPDFMetadata(c *fiber.Ctx) error {
	fileID := c.Query("file_id")
	documentType := c.Query("document_type", "Document")
	citizenName := c.Query("citizen_name")
	nationalID := c.Query("national_id")

	if fileID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Google Drive file ID is required",
			"data":    nil,
		})
	}

	// Get file info from Google Drive
	fileInfo := utils.GetDriveFileInfo(fileID)

	// Prepare certification metadata
	certInfo := utils.CertificationInfo{
		CitizenName:   citizenName,
		NationalID:    nationalID,
		DocumentType:  documentType,
		CertifiedDate: time.Now(),
		CertifierName: "CertiKiosk System",
	}

	metadata := utils.GenerateCertificationMetadata(certInfo)

	// Merge file info with certification metadata
	response := fiber.Map{
		"file_info":       fileInfo,
		"certification":   metadata,
		"download_url":    fileInfo["download_url"],
		"view_url":        fileInfo["view_url"],
		"stamped_pdf_url": "/api/documents/generate-stamped-pdf?file_id=" + fileID,
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Stamped PDF metadata generated",
		"data":    response,
	})
}
