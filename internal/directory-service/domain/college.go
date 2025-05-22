package domain

import (
	"gorm.io/gorm"
)

// College represents a college directory entry
type College struct {
	gorm.Model
	CollegeName  string        `gorm:"size:255" json:"collegeName"`
	Branches     []Branch      `gorm:"foreignKey:CollegeID" json:"branches"`
	Users        []CollegeUser `gorm:"foreignKey:CollegeID" json:"users"`
	StudentCount int           `json:"studentCount"`
	FacultyCount int           `json:"facultyCount"`
}
