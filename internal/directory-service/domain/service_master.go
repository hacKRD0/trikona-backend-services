package domain

type ServiceMaster struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:255" json:"name"`
}
