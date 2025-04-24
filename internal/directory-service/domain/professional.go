package domain

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Professional represents a professional directory entry
type Professional struct {
    gorm.Model
    UserID        uint           `gorm:"not null;uniqueIndex" json:"userId"`
    CurrentTitle  string         `gorm:"size:100" json:"currentTitle"`
    ExperienceYrs int            `gorm:"not null" json:"experienceYrs"`
    Skills        pq.StringArray `gorm:"type:text[]" json:"skills"`
    Industries    pq.StringArray `gorm:"type:text[]" json:"industries"`
}
