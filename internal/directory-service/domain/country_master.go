package domain

type CountryMaster struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Name    string `gorm:"size:255" json:"name"`
	ISOCode string `gorm:"size:255" json:"iso_code"`
}
