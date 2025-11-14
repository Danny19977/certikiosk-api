package documents

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Danny19977/certikiosk.git/database"
	"github.com/Danny19977/certikiosk.git/models"
	"github.com/Danny19977/certikiosk.git/utils"
	"github.com/gofiber/fiber/v2"
)

// Helper functions
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func splitString(s, sep string) []string {
	return strings.Split(s, sep)
}

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
		Email             string `json:"email"`
		DocumentUUID      string `json:"document_uuid"`
		DocumentType      string `json:"document_type"`
		FileID            string `json:"file_id"`              // Google Drive file ID
		GoogleDriveFileID string `json:"google_drive_file_id"` // Alternative parameter name
		FileId            string `json:"fileId"`               // CamelCase variant
	}

	var input EmailInput

	// Parse JSON or form data
	if err := c.BodyParser(&input); err != nil {
		// Try parsing as multipart form
		input.Email = c.FormValue("email")
		input.DocumentUUID = c.FormValue("document_uuid")
		input.DocumentType = c.FormValue("document_type")
		input.FileID = c.FormValue("file_id")
		input.GoogleDriveFileID = c.FormValue("google_drive_file_id")
		input.FileId = c.FormValue("fileId")

		// Also try query parameters
		if input.Email == "" {
			input.Email = c.Query("email")
		}
		if input.FileID == "" {
			input.FileID = c.Query("file_id")
		}
		if input.FileId == "" {
			input.FileId = c.Query("fileId")
		}
		if input.GoogleDriveFileID == "" {
			input.GoogleDriveFileID = c.Query("google_drive_file_id")
		}
		if input.DocumentType == "" {
			input.DocumentType = c.Query("document_type")
		}
		if input.DocumentUUID == "" {
			input.DocumentUUID = c.Query("document_uuid")
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

	// Get the Google Drive file ID from any of the parameter variants
	googleDriveFileID := input.FileID
	if googleDriveFileID == "" {
		googleDriveFileID = input.GoogleDriveFileID
	}
	if googleDriveFileID == "" {
		googleDriveFileID = input.FileId
	}

	// Get the uploaded stamped PDF from frontend - try multiple field names
	// Frontend sends the STAMPED PDF in "document" field
	file, err := c.FormFile("document")
	if err != nil {
		// Try alternative field names
		file, err = c.FormFile("pdf")
		if err != nil {
			file, err = c.FormFile("pdfFile")
		}
	}

	// Note: We DON'T need stamp image anymore because frontend already stamped the PDF
	// Keeping this code for backward compatibility if needed
	stampFile, stampErr := c.FormFile("stamp")
	if stampErr != nil {
		stampFile, stampErr = c.FormFile("stampImage")
	}

	var stampData []byte
	if stampErr == nil && stampFile != nil {
		stampHandle, err := stampFile.Open()
		if err == nil {
			defer stampHandle.Close()
			stampData, _ = io.ReadAll(stampHandle)
		}
	}

	var pdfData []byte

	// PRIORITY 1: Use the uploaded stamped PDF from frontend
	if err == nil && file != nil {
		// Read the uploaded stamped PDF
		fileHandle, err := file.Open()
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to read uploaded file",
				"error":   err.Error(),
			})
		}
		defer fileHandle.Close()

		pdfData, err = io.ReadAll(fileHandle)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to read file content",
				"error":   err.Error(),
			})
		}
	} else if googleDriveFileID != "" {
		// FALLBACK: Only download from Google Drive if NO file was uploaded

		pdfData, err = utils.DownloadFileFromDrive(googleDriveFileID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to download file from Google Drive",
				"error":   err.Error(),
			})
		}

		if len(pdfData) == 0 {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Downloaded file is empty (0 bytes). The file may be private or the link is incorrect.",
				"data":    nil,
			})
		}

		// Set document type if not provided
		if input.DocumentType == "" {
			input.DocumentType = "Document"
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

		// Try to download from Google Drive if DocumentDataUrl contains a Google Drive link
		if document.DocumentDataUrl != "" {
			// Extract file ID from Google Drive URL if present
			// Example: https://drive.google.com/file/d/FILE_ID/view
			// Or: https://drive.google.com/uc?export=download&id=FILE_ID
			var extractedFileID string

			// Try to extract file ID from URL
			if len(document.DocumentDataUrl) > 0 && (contains(document.DocumentDataUrl, "drive.google.com") ||
				contains(document.DocumentDataUrl, "docs.google.com")) {

				// Simple extraction - you may want to use regex for better accuracy
				if contains(document.DocumentDataUrl, "/file/d/") {
					parts := splitString(document.DocumentDataUrl, "/file/d/")
					if len(parts) > 1 {
						idParts := splitString(parts[1], "/")
						if len(idParts) > 0 {
							extractedFileID = idParts[0]
						}
					}
				} else if contains(document.DocumentDataUrl, "id=") {
					parts := splitString(document.DocumentDataUrl, "id=")
					if len(parts) > 1 {
						idParts := splitString(parts[1], "&")
						extractedFileID = idParts[0]
					}
				}
			}

			if extractedFileID != "" {
				pdfData, err = utils.DownloadFileFromDrive(extractedFileID)
			}
		}

		// If we still don't have data, return error
		if len(pdfData) == 0 {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Please provide PDF file, Google Drive file ID (file_id), or ensure the document URL is accessible",
				"data":    nil,
			})
		}
	}

	// Final check: ensure we have PDF data from some source
	if len(pdfData) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Either PDF file, document UUID, or Google Drive file ID (file_id) is required",
			"data":    nil,
		})
	}

	// Set default document type if not provided
	if input.DocumentType == "" {
		input.DocumentType = "Document"
	}

	// Determine the document identifier for logging
	docIdentifier := input.DocumentUUID
	if docIdentifier == "" && googleDriveFileID != "" {
		docIdentifier = googleDriveFileID
	}

	// Detect file type from signature
	var fileExt string
	var mimeType string
	if len(pdfData) >= 4 {
		signature := string(pdfData[:4])

		if signature == "%PDF" {
			fileExt = "pdf"
			mimeType = "application/pdf"
		} else if pdfData[0] == 0x89 && pdfData[1] == 0x50 && pdfData[2] == 0x4E && pdfData[3] == 0x47 {
			fileExt = "png"
			mimeType = "image/png"
		} else if pdfData[0] == 0xFF && pdfData[1] == 0xD8 && pdfData[2] == 0xFF {
			fileExt = "jpg"
			mimeType = "image/jpeg"
		} else {
			fileExt = "pdf"
			mimeType = "application/pdf"
		}
	} else {
		fileExt = "pdf"
		mimeType = "application/pdf"
	}

	// Send email with file attachment (and stamp if provided)
	err = utils.SendDocumentEmailWithStamp(input.Email, input.DocumentType, docIdentifier, pdfData, stampData, fileExt, mimeType)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to send email",
			"error":   err.Error(),
		})
	}

	// Log email sent activity
	utils.LogCreateWithDB(database.DB, c, "document_email", "Document sent to "+input.Email, docIdentifier)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Document sent successfully to " + input.Email,
		"data": fiber.Map{
			"email":         input.Email,
			"document_type": input.DocumentType,
			"document_uuid": input.DocumentUUID,
			"file_id":       googleDriveFileID,
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

// DownloadGoogleDriveFile - Proxy endpoint to download Google Drive files (bypasses CORS)
func DownloadGoogleDriveFile(c *fiber.Ctx) error {
	// Support multiple parameter names for compatibility
	fileID := c.Query("file_id")
	if fileID == "" {
		fileID = c.Query("fileId") // Support camelCase
	}
	if fileID == "" {
		fileID = c.Params("file_id")
	}

	if fileID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Google Drive file ID is required (use ?fileId=YOUR_FILE_ID or ?file_id=YOUR_FILE_ID)",
			"data":    nil,
		})
	}

	// Download file from Google Drive using backend (bypasses CORS)
	fileData, err := utils.DownloadFileFromDrive(fileID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to download file from Google Drive",
			"error":   err.Error(),
		})
	}

	if len(fileData) == 0 {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Downloaded file is empty (0 bytes). The file may be private or the link is incorrect.",
			"data":    nil,
		})
	}

	// Set CORS headers
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Set("Access-Control-Allow-Headers", "Content-Type")

	// Detect content type from file data
	contentType := "application/octet-stream"
	if len(fileData) > 0 {
		// Check for common file types
		if len(fileData) >= 4 && string(fileData[0:4]) == "%PDF" {
			contentType = "application/pdf"
		} else if len(fileData) >= 2 && fileData[0] == 0xFF && fileData[1] == 0xD8 {
			contentType = "image/jpeg"
		} else if len(fileData) >= 8 && string(fileData[0:8]) == "\x89PNG\r\n\x1a\n" {
			contentType = "image/png"
		} else if len(fileData) >= 6 && string(fileData[0:6]) == "GIF87a" || string(fileData[0:6]) == "GIF89a" {
			contentType = "image/gif"
		}
	}

	// Set content type
	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", "inline; filename=\"document\"")
	c.Set("Content-Length", fmt.Sprintf("%d", len(fileData)))

	// Return the file data
	return c.Send(fileData)
}

// GetGoogleDriveFileMetadata - Get metadata for a Google Drive file
func GetGoogleDriveFileMetadata(c *fiber.Ctx) error {
	fileID := c.Query("file_id")
	if fileID == "" {
		fileID = c.Params("file_id")
	}

	if fileID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Google Drive file ID is required",
			"data":    nil,
		})
	}

	// Get file metadata from Google Drive
	metadata, err := utils.GetFileMetadata(fileID)
	if err != nil {
		// If API call fails, return basic info
		metadata = map[string]interface{}{
			"file_id":      fileID,
			"view_url":     utils.GetDriveViewURL(fileID),
			"download_url": utils.GetPublicFileURL(fileID),
			"proxy_url":    "/api/public/documents/gdrive/download/" + fileID,
			"error":        err.Error(),
		}
	} else {
		// Add proxy URL to metadata
		metadata["proxy_url"] = "/api/public/documents/gdrive/download/" + fileID
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "File metadata retrieved",
		"data":    metadata,
	})
}
