package repository

import (
	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ProfessionalRepository interface {
    Count(params *domain.ProfessionalFilterParams) (int64, error)
    Find(params *domain.ProfessionalFilterParams, offset, limit int) ([]domain.Professional, error)
}

type professionalRepository struct {
    db *gorm.DB
}

func NewProfessionalRepository(db *gorm.DB) ProfessionalRepository {
    return &professionalRepository{db: db}
}

func (r *professionalRepository) Count(params *domain.ProfessionalFilterParams) (int64, error) {
    var count int64
    q := r.db.Model(&domain.Professional{})
    q = applyProfessionalFilters(q, params)
    if err := q.Count(&count).Error; err != nil {
        return 0, err
    }
    return count, nil
}

func (r *professionalRepository) Find(params *domain.ProfessionalFilterParams, offset, limit int) ([]domain.Professional, error) {
    var list []domain.Professional
    q := r.db.Model(&domain.Professional{})
    q = applyProfessionalFilters(q, params)
    if err := q.Offset(offset).Limit(limit).Find(&list).Error; err != nil {
        return nil, err
    }
    return list, nil
}

func applyProfessionalFilters(db *gorm.DB, params *domain.ProfessionalFilterParams) *gorm.DB {
    if params.CurrentTitle != nil {
        db = db.Where("current_title ILIKE ?", "%"+*params.CurrentTitle+"%")
    }
    if params.MinExperience != nil {
        db = db.Where("experience_yrs >= ?", *params.MinExperience)
    }
    if len(params.Skills) > 0 {
        db = db.Where("skills @> ?", pq.Array(params.Skills))
    }
    if len(params.Industries) > 0 {
        db = db.Where("industries @> ?", pq.Array(params.Industries))
    }
    if params.SearchTerm != nil && *params.SearchTerm != "" {
        t := "%" + *params.SearchTerm + "%"
        db = db.Where(
            "first_name ILIKE ? OR last_name ILIKE ? OR current_title ILIKE ?",
            t, t, t,
        )
    }
    return db
}
