package repository

import (
	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CollegeRepository interface {
    Count(params *domain.CollegeFilterParams) (int64, error)
    Find(params *domain.CollegeFilterParams, offset, limit int) ([]domain.College, error)
}

type collegeRepository struct {
    db *gorm.DB
}

func NewCollegeRepository(db *gorm.DB) CollegeRepository {
    return &collegeRepository{db: db}
}

func (r *collegeRepository) Count(params *domain.CollegeFilterParams) (int64, error) {
    var count int64
    q := r.db.Model(&domain.College{})
    q = applyCollegeFilters(q, params)
    if err := q.Count(&count).Error; err != nil {
        return 0, err
    }
    return count, nil
}

func (r *collegeRepository) Find(params *domain.CollegeFilterParams, offset, limit int) ([]domain.College, error) {
    var list []domain.College
    q := r.db.Model(&domain.College{})
    q = applyCollegeFilters(q, params)
    if err := q.Offset(offset).Limit(limit).Find(&list).Error; err != nil {
        return nil, err
    }
    return list, nil
}

func applyCollegeFilters(db *gorm.DB, params *domain.CollegeFilterParams) *gorm.DB {
    if params.CollegeName != nil {
        db = db.Where("college_name ILIKE ?", "%"+*params.CollegeName+"%")
    }
    if params.Location != nil {
        db = db.Where("location ILIKE ?", "%"+*params.Location+"%")
    }
    if params.Accreditation != nil {
        db = db.Where("accreditation ILIKE ?", "%"+*params.Accreditation+"%")
    }
    if len(params.Departments) > 0 {
        db = db.Where("departments @> ?", pq.Array(params.Departments))
    }
    if params.SearchTerm != nil && *params.SearchTerm != "" {
        t := "%" + *params.SearchTerm + "%"
        db = db.Where(
            "college_name ILIKE ? OR location ILIKE ?",
            t, t,
        )
    }
    return db
}
