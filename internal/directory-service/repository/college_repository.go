package repository

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"github.com/hacKRD0/trikona_go/pkg/utils"
	"gorm.io/gorm"
)

type CollegeRepository interface {
	Count(params *domain.CollegeFilterParams) (int64, error)
	Find(params *domain.CollegeFilterParams, offset, limit int) ([]domain.College, error)
	GetByID(id uint) (*domain.College, error)
	Create(college *domain.College) error
	Update(college *domain.College) error
	Delete(id uint) error
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
	q := r.db.
		Preload("Branches").
		Preload("Branches.State").
		Model(&domain.College{})
	q = applyCollegeFilters(q, params)
	if err := q.Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *collegeRepository) GetByID(id uint) (*domain.College, error) {
	var college domain.College
	result := r.db.First(&college, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &college, nil
}

func (r *collegeRepository) Create(college *domain.College) error {
	return r.db.Create(college).Error
}

func (r *collegeRepository) Update(college *domain.College) error {
	return r.db.Save(college).Error
}

func (r *collegeRepository) Delete(id uint) error {
	return r.db.Delete(&domain.College{}, id).Error
}

func applyCollegeFilters(db *gorm.DB, params *domain.CollegeFilterParams) *gorm.DB {
	fmt.Println("Applying college filters:", params)
	if params.StudentCountRanges != nil {
		var conditions []string
		var args []interface{}
		ranges := strings.Split(params.StudentCountRanges[0], ",")
		for _, r := range ranges {
			r = strings.TrimSpace(r)
			if r == "500+" {
				conditions = append(conditions, "student_count > ?")
				args = append(args, 500)
				continue
			}
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				continue
			}
			min, err1 := strconv.ParseFloat(parts[0], 32)
			max, err2 := strconv.ParseFloat(parts[1], 32)
			if err1 != nil || err2 != nil {
				continue
			}
			conditions = append(conditions, "student_count BETWEEN ? AND ?")
			args = append(args, min, max)
		}
		if len(conditions) > 0 {
			whereClause := strings.Join(conditions, " OR ")
			db = db.Where(whereClause, args...)
		}
	}
	if params.FacultyCountRanges != nil {
		var conditions []string
		var args []interface{}
		ranges := strings.Split(params.FacultyCountRanges[0], ",")
		for _, r := range ranges {
			r = strings.TrimSpace(r)
			if r == "500+" {
				conditions = append(conditions, "faculty_count > ?")
				args = append(args, 500)
				continue
			}
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				continue
			}
			min, err1 := strconv.ParseFloat(parts[0], 32)
			max, err2 := strconv.ParseFloat(parts[1], 32)
			if err1 != nil || err2 != nil {
				continue
			}
			conditions = append(conditions, "faculty_count BETWEEN ? AND ?")
			args = append(args, min, max)
		}
		if len(conditions) > 0 {
			whereClause := strings.Join(conditions, " OR ")
			db = db.Where(whereClause, args...)
		}
	}
	// Join the branches and related tables once
	if len(params.Country) > 0 || len(params.States) > 0 {
		db = db.Joins(`
			JOIN branches ON branches.college_id = colleges.id
		`)

		if len(params.Country) > 0 {
			db = db.Joins(`
				JOIN country_masters ON branches.country_id = country_masters.id
			`).Where("LOWER(country_masters.name) IN ?", utils.LowerCase(params.Country))
		}

		if len(params.States) > 0 {
			db = db.Joins(`
				JOIN state_masters ON branches.state_id = state_masters.id
			`).Where("LOWER(state_masters.name) IN ?", utils.LowerCase(params.States))
		}
	}
	if params.SearchTerm != nil && *params.SearchTerm != "" {
		t := "%" + *params.SearchTerm + "%"
		db = db.Where(
			"college_name ILIKE ?",
			t,
		)
	}
	return db
}
