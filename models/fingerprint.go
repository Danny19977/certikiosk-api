package models

import "time"

type Fingerprint struct {
	UUID string `gorm:"primaryKey;not null;unique" json:"uuid"`

	CitizensUUID string    `json:"citizens_uuid"`
	FingerprintData string    `json:"fingerprint_data"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
