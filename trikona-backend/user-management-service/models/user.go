package models

import (
	"gorm.io/gorm"
)

const (
	RoleStudent     = "student"
	RoleProfessional = "professional"
	RoleCompany     = "company"
	RoleModerator   = "moderator"
	RoleAdmin       = "admin"
	RoleGuest       = "guest"

	UserStatusPending  = "pending"
	UserStatusActive   = "active"
	UserStatusRejected = "rejected"
)

type User struct {
	gorm.Model

	FirstName string `gorm:"size:100;not null" json:"firstName"`
	LastName  string `gorm:"size:100;not null" json:"lastName"`
	Email     string `gorm:"size:100;unique;not null" json:"email"`
	Password  string `gorm:"size:255;not null" json:"-"` // Do not expose password in JSON responses
	Role      string `gorm:"size:20;not null;default:'guest'" json:"role"`
	Status    string `gorm:"size:20;not null;default:'pending'" json:"status"`
	LinkedInURL string `gorm:"size:255; not null; default:''" json:"linkedInUrl"`
}

