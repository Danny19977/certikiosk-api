package database

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Danny19977/certikiosk.git/models"
	"github.com/Danny19977/certikiosk.git/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	p := utils.Env("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		log.Fatalf("[error] failed to parse DB_PORT (%s): %v", p, err)
	}

	DNS := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", utils.Env("DB_HOST"), port, utils.Env("DB_USER"), utils.Env("DB_PASSWORD"), utils.Env("DB_NAME"))

	// Log the DSN (without exposing the password in logs)
	safeDNS := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", utils.Env("DB_HOST"), port, utils.Env("DB_USER"), utils.Env("DB_NAME"))
	log.Printf("[info] connecting to database: %s", safeDNS)

	connection, err := gorm.Open(postgres.Open(DNS), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("[error] failed to initialize database, got error %v", err)
	}

	DB = connection

	// Enable UUID extension for PostgreSQL
	if err := connection.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		panic("Failed to create uuid-ossp extension: " + err.Error())
	}

	// Migrate in proper order - parent tables first, then child tables
	connection.AutoMigrate(
		&models.User{},
		&models.UserLogs{},
		&models.Notification{},
		&models.PasswordReset{},
		&models.Citizens{},
		&models.Fingerprint{},
		&models.Documents{},
		&models.Certification{},
	)
}
