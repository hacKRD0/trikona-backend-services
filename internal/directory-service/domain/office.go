package domain

import "gorm.io/gorm"

type Office struct {
	gorm.Model
	CorporateID uint          `gorm:"not null" json:"corporateId"`
	Corporate   Corporate     `gorm:"foreignKey:CorporateID" json:"corporate"`
	Name        string        `gorm:"size:255" json:"name"`
	CountryID   uint          `gorm:"not null" json:"countryId"`
	Country     CountryMaster `gorm:"foreignKey:CountryID" json:"country"`
	StateID     uint          `gorm:"not null" json:"stateId"`
	State       StateMaster   `gorm:"foreignKey:StateID" json:"state"`
	City        string        `gorm:"size:255" json:"city"`
	Address     string        `gorm:"size:255" json:"address"`
	PinCode     string        `gorm:"size:255" json:"pinCode"`
	Phone       string        `gorm:"size:255" json:"phone"`
	HeadCount   int           `json:"headCount"`
}
