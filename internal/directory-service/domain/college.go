package domain

import "github.com/lib/pq"

// College represents a college directory entry
type College struct {
    UserID       uint           `gorm:"primaryKey" json:"userId"`
    CollegeName  string         `gorm:"size:255" json:"collegeName"`
    Location     string         `gorm:"size:255" json:"location"`
    Accreditation string        `gorm:"size:100" json:"accreditation"`
    Departments  pq.StringArray `gorm:"type:text[]" json:"departments"`
}
