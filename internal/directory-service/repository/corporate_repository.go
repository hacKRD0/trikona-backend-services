package repository

import (
	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
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
    q := r.db.Model(&domain.Corporate{})
    q = applyCorporateFilters(q, params)
    if err := q.Offset(offset).Limit(limit).Find(&list).Error; err != nil {
        return nil, err
    }
    return list, nil
}

func applyCorporateFilters(db *gorm.DB, params *domain.CorporateFilterParams) *gorm.DB {
    if params.CompanyName != nil {
        db = db.Where("company_name ILIKE ?", "%"+*params.CompanyName+"%")
    }
    if params.Industry != nil {
        db = db.Where("industry ILIKE ?", "%"+*params.Industry+"%")
    }
    if params.MinSize != nil {
        db = db.Where("size >= ?", *params.MinSize)
    }
    if params.Headquarters != nil {
        db = db.Where("headquarters ILIKE ?", "%"+*params.Headquarters+"%")
    }
    if params.SearchTerm != nil && *params.SearchTerm != "" {
        t := "%" + *params.SearchTerm + "%"
        db = db.Where(
            "company_name ILIKE ? OR headquarters ILIKE ?",
            t, t,
        )
    }
    return db
}
