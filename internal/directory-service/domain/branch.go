package domain

import "gorm.io/gorm"

type Branch struct {
	gorm.Model
	CollegeID    uint          `gorm:"not null" json:"collegeId"`
	College      College       `gorm:"foreignKey:CollegeID" json:"college"`
	Name         string        `gorm:"size:255" json:"name"`
	CountryID    uint          `gorm:"not null" json:"countryId"`
	Country      CountryMaster `gorm:"foreignKey:CountryID" json:"country"`
	StateID      uint          `gorm:"not null" json:"stateId"`
	State        StateMaster   `gorm:"foreignKey:StateID" json:"state"`
	City         string        `gorm:"size:255" json:"city"`
	Address      string        `gorm:"size:255" json:"address"`
	PinCode      string        `gorm:"size:255" json:"pinCode"`
	Phone        string        `gorm:"size:255" json:"phone"`
	StudentCount int           `json:"studentCount"`
	FacultyCount int           `json:"facultyCount"`
}
