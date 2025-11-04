package models

import "time"

type Certification struct {
	UUID              string `gorm:"primaryKey;not null;unique" json:"uuid"`
	CitizensUUID      string `json:"citizens_uuid"`
	DocumentUUID      string `json:"document_uuid"`
	Aprovel           bool   `json:"aprovel"`
	CertifiedDocument string `json:"certified_document"`
	StampDetails      string `json:"stamp_details"`
	OutputFormat      string `json:"output_format"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
