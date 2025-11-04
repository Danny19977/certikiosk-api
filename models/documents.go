package models

import "time"

type Documents struct {
	UUID string `gorm:"primaryKey;not null;unique" json:"uuid"`

	DocumentType    string    `json:"document_type"`
	DocumentDataUrl string    `json:"document_data"`
	IssueDate       time.Time `json:"issue_date"`
	IsActive        bool      `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
