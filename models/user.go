package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UUID            string `gorm:"primaryKey;not null;unique" json:"uuid"`
	Fullname        string `json:"fullname"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Title           string `json:"title"`
	Password        string `json:"password"`
	ConformPassword string `json:"confirm_password"`
	Role            string `json:"role"`
	Permission      string `json:"permission"`
	Status          bool   `json:"status"`

	Signature string `json:"signature"`
}

type UserResponse struct {
	UUID       string    `json:"uuid"`
	Fullname   string    `json:"fullname"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Title      string    `json:"title"`
	Role       string    `json:"role"`
	Permission string    `json:"permission"`
	Status     bool      `json:"status"`
	Signature  string    `json:"signature"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Helper to convert User to UserResponse
func (u *User) ToUserResponse() UserResponse {
	return UserResponse{
		UUID:       u.UUID,
		Fullname:   u.Fullname,
		Email:      u.Email,
		Phone:      u.Phone,
		Title:      u.Title,
		Role:       u.Role,
		Permission: u.Permission,
		Status:     u.Status,
		Signature:  u.Signature,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

type UserPaginate struct {
	UUID       string    `json:"uuid"`
	Fullname   string    `json:"fullname"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Title      string    `json:"title"`
	Role       string    `json:"role"`
	Permission string    `json:"permission"`
	Status     bool      `json:"status"`
	Signature  string    `json:"signature"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ...existing code...

type Login struct {
	Identifier string `json:"identifier" validate:"required"`
	// Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (u *User) SetPassword(p string) {
	hp, _ := bcrypt.GenerateFromPassword([]byte(p), 14)
	u.Password = string(hp)
}

func (u *User) ComparePassword(p string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	return err
}

func (u *User) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(&User{}).Count(&total)
	return total
}
