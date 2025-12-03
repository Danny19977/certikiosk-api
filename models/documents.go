package models

import "time"

type Documents struct {
	UUID string `gorm:"primaryKey;not null;unique" json:"uuid"`

	NationalID      int64     `json:"national_id"`
	DocumentType    string    `json:"document_type"`
	DocumentDataUrl string    `json:"document_data" gorm:"column:document_data_url"`
	UserUUID        string    `json:"user_uuid"`
	IssueDate       time.Time `json:"issue_date"`
	IsActive        bool      `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
