package models

import (
	"gorm.io/gorm"
)

const (
	RoleStudent = "student"
	RoleProfessional = "professional"
	RoleCompany = "company"
	RoleModerator = "moderator"
	RoleAdmin = "admin"
	RoleGuest = "guest"
)

type User struct {
	gorm.Model
	
	FirstName string			`gorm:"size:100;not null" json:"first_name"`
	LastName  string      `gorm:"size:100;not null" json:"last_name"`
	Email     string      `gorm:"size:100;unique;not null" json:"email"`
	Password  string      `gorm:"size:255;not null" json:"-"` // Do not expose password in JSON responses
	Role      string      `gorm:"size:20;not null;default:'guest'" json:"role"`
}

