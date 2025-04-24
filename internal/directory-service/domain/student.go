// internal/directory-service/domain/student.go
package domain

import (
	um "github.com/hacKRD0/trikona_go/internal/user-management-service/domain"
	"gorm.io/gorm"
)

// Student represents the directory entry for a student
// Skills are many2many via student_skill join table

type Student struct {
	gorm.Model

	// Link back to the User in user-management service
	UserID uint      `gorm:"not null;uniqueIndex;constraint:OnDelete:CASCADE" json:"userId"`
	User   um.User   `gorm:"foreignKey:UserID;references:ID" json:"user"`

	Educations           []Education    `gorm:"foreignKey:UserID" json:"educations"`
	Experiences          []Experience   `gorm:"foreignKey:UserID" json:"experiences"`
	TotalExperienceYears int            `gorm:"not null;default:0" json:"totalExperienceYears"`
  Skills []SkillMaster `gorm:"many2many:student_skill;constraint:OnDelete:CASCADE" json:"skills"`
}