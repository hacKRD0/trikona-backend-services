package domain

import (
	"github.com/hacKRD0/trikona_go/internal/user-management-service/domain"
	"gorm.io/gorm"
)

// CollegeUser represents a user associated with a college entity
type CollegeUser struct {
	gorm.Model
	CollegeID uint              `gorm:"not null;index" json:"collegeId"`
	UserID    uint              `gorm:"not null;index" json:"userId"`
	UserRole  domain.UserRole   `gorm:"size:20;not null" json:"userRole"`
	Status    domain.UserStatus `gorm:"size:20;not null;default:'pending'" json:"status"`
	College   College           `gorm:"foreignKey:CollegeID" json:"-"`
	User      domain.User       `gorm:"foreignKey:UserID" json:"-"`
}
