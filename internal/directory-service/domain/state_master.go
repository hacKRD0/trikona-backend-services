package domain

type StateMaster struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	Name      string       `gorm:"size:255" json:"name"`
	CountryID uint         `gorm:"not null" json:"countryId"`
	Country   CountryMaster `gorm:"foreignKey:CountryID" json:"country"`
}
