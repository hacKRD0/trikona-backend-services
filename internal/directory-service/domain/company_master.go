// internal/directory-service/domain/company_master.go
package domain

// CompanyMaster holds the canonical list of companies
// referenced by Experience.CompanyID

type CompanyMaster struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:255;unique;not null" json:"name"`
}
