package domain

import (
	"github.com/hacKRD0/trikona_go/internal/user-management-service/domain"
	"gorm.io/gorm"
)

// CorporateUser represents a user associated with a corporate entity
type CorporateUser struct {
	gorm.Model
	CorporateID uint              `gorm:"not null;index" json:"corporateId"`
	UserID      uint              `gorm:"not null;index" json:"userId"`
	UserRole    domain.UserRole   `gorm:"size:20;not null" json:"userRole"`
	Status      domain.UserStatus `gorm:"size:20;not null;default:'pending'" json:"status"`
	Corporate   Corporate         `gorm:"foreignKey:CorporateID" json:"-"`
	User        domain.User       `gorm:"foreignKey:UserID" json:"-"`
}
