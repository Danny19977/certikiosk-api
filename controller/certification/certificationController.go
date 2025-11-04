package certification

import (
	"strconv"
	"time"

	"github.com/Danny19977/certikiosk.git/database"
	"github.com/Danny19977/certikiosk.git/models"
	"github.com/Danny19977/certikiosk.git/utils"
	"github.com/gofiber/fiber/v2"
)

// CertifyDocument - Main function to certify a document with stamp
func CertifyDocument(c *fiber.Ctx) error {
	type CertificationInput struct {
		CitizensUUID    string `json:"citizens_uuid"`
		DocumentUUID    string `json:"document_uuid"`
		FingerprintData string `json:"fingerprint_data"`
		StampDetails    string `json:"stamp_details"`
		OutputFormat    string `json:"output_format"` // "pdf" or "print"
	}

	var input CertificationInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
	}

	// Validate required fields
	if input.CitizensUUID == "" || input.DocumentUUID == "" || input.FingerprintData == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Citizens UUID, Document UUID, and Fingerprint data are required",
			"data":    nil,
		})
	}

	// Step 1: Verify citizen exists
	var citizen models.Citizens
	if err := database.DB.Where("uuid = ?", input.CitizensUUID).First(&citizen).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Citizen not found",
			"data":    nil,
		})
	}

	// Step 2: Verify fingerprint
	var fingerprint models.Fingerprint
	if err := database.DB.Where("citizens_uuid = ? AND fingerprint_data = ?", input.CitizensUUID, input.FingerprintData).First(&fingerprint).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Fingerprint verification failed",
			"data":    nil,
		})
	}

	// Step 3: Verify document exists
	var document models.Documents
	if err := database.DB.Where("uuid = ?", input.DocumentUUID).First(&document).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Document not found",
			"data":    nil,
		})
	}

	// Step 4: Check if document is active
	if !document.IsActive {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Document is not active",
			"data":    nil,
		})
	}

	// Step 5: Apply certification stamp (placeholder - will use PDF utility)
	// TODO: Integrate with PDF generation utility to add stamp
	certifiedDocumentUrl := document.DocumentDataUrl + "_certified"

	// Set default output format if not provided
	if input.OutputFormat == "" {
		input.OutputFormat = "pdf"
	}

	// Step 6: Create certification record
	certification := models.Certification{
		UUID:              utils.GenerateUUID(),
		CitizensUUID:      input.CitizensUUID,
		DocumentUUID:      input.DocumentUUID,
		Aprovel:           true,
		CertifiedDocument: certifiedDocumentUrl,
		StampDetails:      input.StampDetails,
		OutputFormat:      input.OutputFormat,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := database.DB.Create(&certification).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create certification record",
			"error":   err.Error(),
		})
	}

	// Log certification activity
	utils.LogCreateWithDB(database.DB, c, "certification", "Document certified for "+citizen.FirstName+" "+citizen.LastName, certification.UUID)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Document certified successfully",
		"data": fiber.Map{
			"certification": certification,
			"citizen":       citizen,
			"document":      document,
		},
	})
}

// GetPaginatedCertifications - Get paginated list of certifications
func GetPaginatedCertifications(c *fiber.Ctx) error {
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

	var certifications []models.Certification
	var totalRecords int64

	db.Model(&models.Certification{}).Count(&totalRecords)

	err = db.Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&certifications).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch certifications",
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
		"message":    "Certifications retrieved successfully",
		"data":       certifications,
		"pagination": pagination,
	})
}

// GetAllCertifications - Get all certifications
func GetAllCertifications(c *fiber.Ctx) error {
	db := database.DB
	var certifications []models.Certification

	if err := db.Find(&certifications).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch certifications",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All certifications retrieved successfully",
		"data":    certifications,
	})
}

// GetCertification - Get a single certification by UUID
func GetCertification(c *fiber.Ctx) error {
	certificationUUID := c.Params("uuid")
	db := database.DB
	var certification models.Certification

	if err := db.Where("uuid = ?", certificationUUID).First(&certification).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Certification not found",
			"data":    nil,
		})
	}

	// Get associated citizen and document
	var citizen models.Citizens
	var document models.Documents

	db.Where("uuid = ?", certification.CitizensUUID).First(&citizen)
	db.Where("uuid = ?", certification.DocumentUUID).First(&document)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Certification found",
		"data": fiber.Map{
			"certification": certification,
			"citizen":       citizen,
			"document":      document,
		},
	})
}

// GetCertificationsByCitizen - Get all certifications for a specific citizen
func GetCertificationsByCitizen(c *fiber.Ctx) error {
	citizenUUID := c.Params("citizen_uuid")
	db := database.DB
	var certifications []models.Certification

	if err := db.Where("citizens_uuid = ?", citizenUUID).Find(&certifications).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch certifications for this citizen",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Certifications retrieved successfully",
		"data":    certifications,
	})
}

// GetCertificationsByDocument - Get all certifications for a specific document
func GetCertificationsByDocument(c *fiber.Ctx) error {
	documentUUID := c.Params("document_uuid")
	db := database.DB
	var certifications []models.Certification

	if err := db.Where("document_uuid = ?", documentUUID).Find(&certifications).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch certifications for this document",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Certifications retrieved successfully",
		"data":    certifications,
	})
}

// DownloadCertifiedDocument - Download the certified document as PDF
func DownloadCertifiedDocument(c *fiber.Ctx) error {
	certificationUUID := c.Params("uuid")
	db := database.DB
	var certification models.Certification

	if err := db.Where("uuid = ?", certificationUUID).First(&certification).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Certification not found",
			"data":    nil,
		})
	}

	// TODO: Implement actual PDF download functionality
	// This should generate/retrieve the certified PDF and send it as download

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Certified document ready for download",
		"data": fiber.Map{
			"download_url": certification.CertifiedDocument,
			"format":       certification.OutputFormat,
		},
	})
}

// PrintCertifiedDocument - Prepare document for printing
func PrintCertifiedDocument(c *fiber.Ctx) error {
	certificationUUID := c.Params("uuid")
	db := database.DB
	var certification models.Certification

	if err := db.Where("uuid = ?", certificationUUID).First(&certification).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Certification not found",
			"data":    nil,
		})
	}

	// TODO: Implement print queue or direct print functionality
	// This could integrate with a print service

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Document sent to print queue",
		"data": fiber.Map{
			"certification_uuid": certification.UUID,
			"print_url":          certification.CertifiedDocument,
		},
	})
}

// RevokeCertification - Revoke a certification (set approval to false)
func RevokeCertification(c *fiber.Ctx) error {
	certificationUUID := c.Params("uuid")
	db := database.DB
	var certification models.Certification

	if err := db.Where("uuid = ?", certificationUUID).First(&certification).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Certification not found",
			"data":    nil,
		})
	}

	certification.Aprovel = false
	certification.UpdatedAt = time.Now()

	if err := db.Save(&certification).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to revoke certification",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Certification revoked successfully",
		"data":    certification,
	})
}

// DeleteCertification - Delete a certification record
func DeleteCertification(c *fiber.Ctx) error {
	certificationUUID := c.Params("uuid")
	db := database.DB

	var certification models.Certification
	if err := db.Where("uuid = ?", certificationUUID).First(&certification).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Certification not found",
			"data":    nil,
		})
	}

	if err := db.Delete(&certification).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete certification",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Certification deleted successfully",
		"data":    nil,
	})
}
