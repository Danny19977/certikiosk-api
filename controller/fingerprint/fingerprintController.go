package fingerprint

import (
	"strconv"
	"time"

	"github.com/Danny19977/certikiosk.git/database"
	"github.com/Danny19977/certikiosk.git/models"
	"github.com/Danny19977/certikiosk.git/utils"
	"github.com/gofiber/fiber/v2"
)

// EnrollFingerprint - Register a fingerprint for a citizen
func EnrollFingerprint(c *fiber.Ctx) error {
	type FingerprintInput struct {
		CitizensUUID    string `json:"citizens_uuid"`
		FingerprintData string `json:"fingerprint_data"`
	}

	var input FingerprintInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
	}

	// Validate required fields
	if input.CitizensUUID == "" || input.FingerprintData == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Citizens UUID and fingerprint data are required",
			"data":    nil,
		})
	}

	// Verify citizen exists and update their fingerprint
	var citizen models.Citizens
	if err := database.DB.Where("uuid = ?", input.CitizensUUID).First(&citizen).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Citizen not found",
			"data":    nil,
		})
	}

	// Check if fingerprint already enrolled
	if citizen.Fingerprint != "" {
		return c.Status(409).JSON(fiber.Map{
			"status":  "error",
			"message": "Fingerprint already enrolled for this citizen",
			"data":    nil,
		})
	}

	// Update citizen with fingerprint data
	citizen.Fingerprint = input.FingerprintData
	citizen.UpdatedAt = time.Now()

	if err := database.DB.Save(&citizen).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to enroll fingerprint",
			"error":   err.Error(),
		})
	}

	// Log fingerprint enrollment
	utils.LogCreateWithDB(database.DB, c, "fingerprint", "Fingerprint enrolled for "+citizen.FirstName+" "+citizen.LastName, citizen.UUID.String())

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Fingerprint enrolled successfully",
		"data":    citizen,
	})
}

// VerifyFingerprint - Verify a fingerprint and return citizen information
func VerifyFingerprint(c *fiber.Ctx) error {
	type FingerprintVerifyInput struct {
		FingerprintData string `json:"fingerprint_data"`
	}

	var input FingerprintVerifyInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
	}

	if input.FingerprintData == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Fingerprint data is required",
			"data":    nil,
		})
	}

	// Find citizen by fingerprint
	var citizen models.Citizens
	if err := database.DB.Where("fingerprint = ?", input.FingerprintData).First(&citizen).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Fingerprint not recognized",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Fingerprint verified successfully",
		"data":    citizen,
	})
}

// GetFingerprintByCitizen - Get fingerprint data for a specific citizen
func GetFingerprintByCitizen(c *fiber.Ctx) error {
	citizenUUID := c.Params("citizen_uuid")

	var citizen models.Citizens
	if err := database.DB.Where("uuid = ?", citizenUUID).First(&citizen).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Citizen not found",
			"data":    nil,
		})
	}

	if citizen.Fingerprint == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No fingerprint found for this citizen",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Fingerprint retrieved successfully",
		"data": fiber.Map{
			"citizen_uuid":     citizen.UUID,
			"fingerprint_data": citizen.Fingerprint,
		},
	})
}

// GetAllFingerprints - Get all fingerprints with pagination
func GetPaginatedFingerprints(c *fiber.Ctx) error {
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

	var citizens []models.Citizens
	var totalRecords int64

	// Only get citizens with fingerprints
	query := db.Model(&models.Citizens{}).Where("fingerprint != ''")
	query.Count(&totalRecords)

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&citizens).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch fingerprints",
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
		"message":    "Fingerprints retrieved successfully",
		"data":       citizens,
		"pagination": pagination,
	})
}

// UpdateFingerprint - Update fingerprint data for a citizen
func UpdateFingerprint(c *fiber.Ctx) error {
	citizenUUID := c.Params("citizen_uuid")

	type UpdateFingerprintInput struct {
		FingerprintData string `json:"fingerprint_data"`
	}

	var input UpdateFingerprintInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
	}

	if input.FingerprintData == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Fingerprint data is required",
			"data":    nil,
		})
	}

	var citizen models.Citizens
	if err := database.DB.Where("uuid = ?", citizenUUID).First(&citizen).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Citizen not found",
			"data":    nil,
		})
	}

	citizen.Fingerprint = input.FingerprintData
	citizen.UpdatedAt = time.Now()

	if err := database.DB.Save(&citizen).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update fingerprint",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Fingerprint updated successfully",
		"data":    citizen,
	})
}

// DeleteFingerprint - Delete a fingerprint record
func DeleteFingerprint(c *fiber.Ctx) error {
	citizenUUID := c.Params("citizen_uuid")

	var citizen models.Citizens
	if err := database.DB.Where("uuid = ?", citizenUUID).First(&citizen).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Citizen not found",
			"data":    nil,
		})
	}

	if citizen.Fingerprint == "" {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No fingerprint found for this citizen",
			"data":    nil,
		})
	}

	// Clear the fingerprint field
	citizen.Fingerprint = ""
	citizen.UpdatedAt = time.Now()

	if err := database.DB.Save(&citizen).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete fingerprint",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Fingerprint deleted successfully",
		"data":    nil,
	})
}
