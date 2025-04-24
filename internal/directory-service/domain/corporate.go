package domain

import "gorm.io/gorm"

// Corporate represents a corporate directory entry
type Corporate struct {
    gorm.Model
    UserID      uint   `gorm:"not null;uniqueIndex" json:"userId"`
    CompanyName string `gorm:"size:255" json:"companyName"`
    Industry    string `gorm:"size:100" json:"industry"`
    Size        int    `json:"size"`        // number of employees
    Headquarters string `gorm:"size:255" json:"headquarters"`
}
