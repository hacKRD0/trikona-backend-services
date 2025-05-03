package repository

import (
	"fmt"

	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"github.com/hacKRD0/trikona_go/pkg/utils"
	"gorm.io/gorm"
)

type CorporateRepository interface {
	Count(params *domain.CorporateFilterParams) (int64, error)
	Find(params *domain.CorporateFilterParams, offset, limit int) ([]domain.Corporate, error)
}

type corporateRepository struct {
	db *gorm.DB
}

func NewCorporateRepository(db *gorm.DB) CorporateRepository {
	return &corporateRepository{db: db}
}

func (r *corporateRepository) Count(params *domain.CorporateFilterParams) (int64, error) {
	var count int64
	q := r.db.Model(&domain.Corporate{})
	q = applyCorporateFilters(q, params)
	if err := q.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *corporateRepository) Find(params *domain.CorporateFilterParams, offset, limit int) ([]domain.Corporate, error) {
	var list []domain.Corporate
	fmt.Println("Applying corporate filters:", params)
	q := r.db.
		Preload("Industries").
		Preload("Sectors").
		Preload("Services").
		Preload("Offices").
		Preload("Offices.Country").
		Preload("Offices.State").
		Model(&domain.Corporate{})
	q = applyCorporateFilters(q, params)
	if err := q.Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, err
	}
	fmt.Println(q)
	return list, nil
}

// applyCorporateFilters applies filters to the query using lowercase name comparison
func applyCorporateFilters(db *gorm.DB, params *domain.CorporateFilterParams) *gorm.DB {
	if params.MinSize != nil {
		db = db.Where("headCount >= ?", *params.MinSize)
	}

	// Filter by country names (case-insensitive)
	if len(params.Country) > 0 {
		db = db.Joins(`
			JOIN offices ON offices.corporate_id = corporates.id
		`).Joins(`
			JOIN country_masters ON offices.country_id = country_masters.id
		`).Where("LOWER(country_masters.name) IN ?", utils.LowerCase(params.Country))
	}

	// Filter by state names (case-insensitive)
	if len(params.States) > 0 {
		db = db.Joins(`
			JOIN offices ON offices.corporate_id = corporates.id
		`).Joins(`
			JOIN state_masters ON offices.state_id = state_masters.id
		`).Where("LOWER(state_masters.name) IN ?", utils.LowerCase(params.States))
	}

	// Filter by sector names (case-insensitive)
	if len(params.Sectors) > 0 {
		db = db.Joins(`
			JOIN corporate_sectors AS cs ON cs.corporate_id = corporates.id
			JOIN sector_masters AS sm ON sm.id = cs.sector_master_id
		`).Where("LOWER(sm.name) IN ?", utils.LowerCase(params.Sectors)).Group("corporates.id")
	}

	// Filter by industry names (case-insensitive)
	if len(params.Industries) > 0 {
		db = db.Joins(`
			JOIN corporate_industries AS ci ON ci.corporate_id = corporates.id
		`).Joins(`
			JOIN industry_masters AS im ON im.id = ci.industry_master_id
		`).Where("LOWER(im.name) IN ?", utils.LowerCase(params.Industries)).Group("corporates.id")
	}

	// Filter by service names (case-insensitive)
	if len(params.Services) > 0 {
		db = db.Joins(`
			JOIN corporate_services AS cs ON cs.corporate_id = corporates.id
			JOIN service_masters AS sm ON sm.id = cs.service_master_id
		`).Where("LOWER(sm.name) IN ?", utils.LowerCase(params.Services)).Group("corporates.id")
	}

	if params.SearchTerm != nil && *params.SearchTerm != "" {
		t := "%" + *params.SearchTerm + "%"
		db = db.Where(
			"CompanyName ILIKE ? OR headquarters ILIKE ?",
			t, t,
		)
	}

	return db
}
