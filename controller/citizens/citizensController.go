package citizens

import (
	"strconv"

	"github.com/Danny19977/certikiosk.git/database"
	"github.com/Danny19977/certikiosk.git/models"
	"github.com/Danny19977/certikiosk.git/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetPaginatedCitizens - Get paginated list of citizens with search
func GetPaginatedCitizens(c *fiber.Ctx) error {
	db := database.DB

	// Parse query parameters for pagination
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit", "15"))
	if err != nil || limit <= 0 {
		limit = 15
	}
	offset := (page - 1) * limit

	// Parse search query
	search := c.Query("search", "")

	var citizens []models.Citizens
	var totalRecords int64

	query := db.Model(&models.Citizens{})
	if search != "" {
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	query.Count(&totalRecords)

	err = query.Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&citizens).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch citizens",
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
		"message":    "Citizens retrieved successfully",
		"data":       citizens,
		"pagination": pagination,
	})
}

// GetAllCitizens - Get all citizens without pagination
func GetAllCitizens(c *fiber.Ctx) error {
	db := database.DB
	var citizens []models.Citizens

	if err := db.Find(&citizens).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch citizens",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All citizens retrieved successfully",
		"data":    citizens,
	})
}

// GetCitizen - Get a single citizen by UUID
func GetCitizen(c *fiber.Ctx) error {
	citizenUUID := c.Params("uuid")
	db := database.DB
	var citizen models.Citizens

	if err := db.Where("uuid = ?", citizenUUID).First(&citizen).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Citizen not found",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Citizen found",
		"data":    citizen,
	})
}

// GetCitizenByNationalID - Get a citizen by National ID
func GetCitizenByNationalID(c *fiber.Ctx) error {
	nationalID := c.Params("national_id")
	db := database.DB
	var citizen models.Citizens

	nationalIDInt, err := strconv.Atoi(nationalID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid National ID format",
			"data":    nil,
		})
	}

	if err := db.Where("national_id = ?", nationalIDInt).First(&citizen).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Citizen not found with this National ID",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Citizen found",
		"data":    citizen,
	})
}

// CreateCitizen - Register a new citizen
func CreateCitizen(c *fiber.Ctx) error {
	type CitizenInput struct {
		NationalID  int    `json:"national_id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		DateOfBirth string `json:"date_of_birth"`
		Email       string `json:"email"`
	}

	var input CitizenInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
	}

	// Validate required fields
	if input.NationalID == 0 || input.FirstName == "" || input.LastName == "" || input.Email == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "All fields are required",
			"data":    nil,
		})
	}

	// Check if citizen with same National ID already exists
	var existingCitizen models.Citizens
	if err := database.DB.Where("national_id = ?", input.NationalID).First(&existingCitizen).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{
			"status":  "error",
			"message": "Citizen with this National ID already exists",
			"data":    nil,
		})
	}

	citizen := models.Citizens{
		UUID:        uuid.New(),
		NationalID:  input.NationalID,
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Phone:       input.Email,
	}

	if err := database.DB.Create(&citizen).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create citizen",
			"error":   err.Error(),
		})
	}

	// Log citizen creation activity
	utils.LogCreateWithDB(database.DB, c, "citizen", citizen.FirstName+" "+citizen.LastName, citizen.UUID.String())

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Citizen registered successfully",
		"data":    citizen,
	})
}

// UpdateCitizen - Update citizen information
func UpdateCitizen(c *fiber.Ctx) error {
	citizenUUID := c.Params("uuid")
	db := database.DB

	type UpdateCitizenInput struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Phone       string `json:"phone"`
	}

	var updateData UpdateCitizenInput

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
	}

	var citizen models.Citizens
	if err := db.Where("uuid = ?", citizenUUID).First(&citizen).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Citizen not found",
			"data":    nil,
		})
	}

	// Update fields
	if updateData.FirstName != "" {
		citizen.FirstName = updateData.FirstName
	}
	if updateData.LastName != "" {
		citizen.LastName = updateData.LastName
	}
	if updateData.Phone != "" {
		citizen.Phone = updateData.Phone
	}

	if err := db.Save(&citizen).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update citizen",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Citizen updated successfully",
		"data":    citizen,
	})
}

// DeleteCitizen - Delete a citizen
func DeleteCitizen(c *fiber.Ctx) error {
	citizenUUID := c.Params("uuid")
	db := database.DB

	var citizen models.Citizens
	if err := db.Where("uuid = ?", citizenUUID).First(&citizen).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Citizen not found",
			"data":    nil,
		})
	}

	if err := db.Delete(&citizen).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete citizen",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Citizen deleted successfully",
		"data":    nil,
	})
}
