package domain

import (
	"gorm.io/gorm"
)

// Corporate represents a corporate directory entry
type Corporate struct {
	gorm.Model
	CompanyName string           `gorm:"size:255" json:"companyName"`
	Industries  []IndustryMaster `gorm:"many2many:corporate_industries;constraint:OnDelete:CASCADE" json:"industries"`
	Services    []ServiceMaster  `gorm:"many2many:corporate_services;constraint:OnDelete:CASCADE" json:"services"`
	Sectors     []SectorMaster   `gorm:"many2many:corporate_sectors;constraint:OnDelete:CASCADE" json:"sectors"`
	HeadCount   int              `json:"headCount"` // number of employees
	Offices     []Office         `gorm:"foreignKey:CorporateID" json:"offices"`
	Users       []CorporateUser  `gorm:"foreignKey:CorporateID" json:"users"`
}
