// internal/directory-service/domain/education.go
package domain

import (
	"time"

	um "github.com/hacKRD0/trikona_go/internal/user-management-service/domain"
)

// Education represents one entry in a user's educational background

type Education struct {
	ID             uint        `gorm:"primaryKey" json:"id"`

	// Link back to User
	UserID         uint        `gorm:"not null;index;constraint:OnDelete:CASCADE" json:"userId"`
	User           um.User     `gorm:"foreignKey:UserID;references:ID" json:"-"`

	CollegeID      uint        `gorm:"not null;constraint:OnDelete:CASCADE" json:"collegeId"`
	College        CollegeMaster `gorm:"foreignKey:CollegeID" json:"college"`

	Degree         string  			`gorm:"string;check:degree in ('Bachelors', 'Masters', 'High School', 'Diploma')" json:"degree"`
	FieldOfStudy   string      `gorm:"size:100" json:"fieldOfStudy"`
	StartDate      time.Time   `gorm:"index;not null" json:"startDate"`
	EndDate        time.Time   `gorm:"not null" json:"endDate"`
	YearOfStudy    int         `gorm:"not null;check:year_of_study>=1" json:"yearOfStudy"`
	CGPA           float32     `gorm:"type:decimal(4,2);not null;check:cgpa>=0 AND cgpa<=10" json:"cgpa"`
	DurationMonths int         `gorm:"not null;default:0" json:"durationMonths"`
	IsLatest       bool        `gorm:"not null;default:false" json:"isLatest"`
}
