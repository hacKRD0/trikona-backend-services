// internal/directory-service/domain/experience.go
package domain

import (
	"time"

	um "github.com/hacKRD0/trikona_go/internal/user-management-service/domain"
)

// Experience represents one job or project a user has done

type Experience struct {
	ID             uint        `gorm:"primaryKey" json:"id"`

	// Link back to User
	UserID         uint        `gorm:"not null;index;constraint:OnDelete:CASCADE" json:"userId"`
	User           um.User     `gorm:"foreignKey:UserID;references:ID" json:"-"`

	CompanyID      uint        `gorm:"not null;constraint:OnDelete:CASCADE" json:"companyId"`
	Company        CompanyMaster `gorm:"foreignKey:CompanyID" json:"company"`

	Title          string      `gorm:"size:100;not null" json:"title"`
	StartDate      time.Time   `gorm:"index;not null" json:"startDate"`
	EndDate        time.Time   `gorm:"not null" json:"endDate"`
	DurationMonths int         `gorm:"not null;default:0" json:"durationMonths"`
	IsLatest       bool        `gorm:"not null;default:false" json:"isLatest"`
}