package repository

import (
	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"gorm.io/gorm"
)

type StudentRepository interface {
	Count(params *domain.StudentFilterParams) (int64, error)
	Find(params *domain.StudentFilterParams, offset, limit int) ([]domain.Student, error)
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) Count(params *domain.StudentFilterParams) (int64, error) {
	var count int64
	q := r.db.
		Model(&domain.Student{})

	q = applyStudentFilters(q, params)

	// Ensure we count distinct students when skills filter is applied
	if len(params.Skills) > 0 {
		q = q.Distinct("students.id")
	}

	if err := q.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *studentRepository) Find(params *domain.StudentFilterParams, offset, limit int) ([]domain.Student, error) {
	var students []domain.Student

	q := r.db.
		Preload("User").
		Preload("User").
 
        // preload latest education and its College
        Preload("Educations", func(db *gorm.DB) *gorm.DB {
            return db.Where("is_latest = ?", true)
        }).
        Preload("Educations.College").
 
        // preload latest experience and its Company
        Preload("Experiences", func(db *gorm.DB) *gorm.DB {
            return db.Where("is_latest = ?", true)
        }).
        Preload("Experiences.Company").
		Preload("Skills").
		Model(&domain.Student{})

	q = applyStudentFilters(q, params)

	if err := q.
		Offset(offset).
		Limit(limit).
		Find(&students).Error; err != nil {
		return nil, err
	}

	return students, nil
}

func applyStudentFilters(db *gorm.DB, params *domain.StudentFilterParams) *gorm.DB {
	// 1) Latest-education filters
	if params.CollegeName != nil || params.Level != nil ||
		params.MinCgpa != nil || params.YearOfStudy != nil {

		db = db.Joins(`
			JOIN educations AS latest_edu
			  ON latest_edu.user_id = students.user_id
			 AND latest_edu.is_latest = TRUE
		`)

		if params.CollegeName != nil && *params.CollegeName != "" {
			db = db.Joins(`
				JOIN college_master AS cm
				  ON cm.id = latest_edu.college_id
			`).Where("cm.name ILIKE ?", "%"+*params.CollegeName+"%")
		}
		if params.Level != nil && *params.Level != "" {
			db = db.Where("latest_edu.degree = ?", *params.Level)
		}
		if params.MinCgpa != nil {
			db = db.Where("latest_edu.cgpa >= ?", *params.MinCgpa)
		}
		if params.YearOfStudy != nil {
			db = db.Where("latest_edu.year_of_study = ?", *params.YearOfStudy)
		}
	}

	// 2) Latest-experience filters
	if params.Company != nil || params.Title != nil ||
		params.MinExperienceYears != nil || params.MaxExperienceYears != nil {

		db = db.Joins(`
			JOIN experiences AS latest_exp
			  ON latest_exp.user_id = students.user_id
			 AND latest_exp.is_latest = TRUE
		`)

		if params.Company != nil && *params.Company != "" {
			db = db.Joins(`
				JOIN company_master AS com
				  ON com.id = latest_exp.company_id
			`).Where("com.name ILIKE ?", "%"+*params.Company+"%")
		}
		if params.Title != nil && *params.Title != "" {
			db = db.Where("latest_exp.title ILIKE ?", "%"+*params.Title+"%")
		}
		if params.MinExperienceYears != nil {
			db = db.Where("students.total_experience_years >= ?", *params.MinExperienceYears)
		}
		if params.MaxExperienceYears != nil {
			db = db.Where("students.total_experience_years <= ?", *params.MaxExperienceYears)
		}
	}

	// 3) Skills filter via many2many join
	if len(params.Skills) > 0 {
		db = db.Joins(`
			JOIN student_skill AS ss
			  ON ss.student_id = students.id
		`).Joins(`
			JOIN skill_masters AS sm
			  ON sm.id = ss.skill_master_id
		`).Where("sm.name IN ?", params.Skills).
			Group("students.id").
			Having("COUNT(DISTINCT sm.id) = ?", len(params.Skills))
	}

	// 4) General name search
	if params.SearchTerm != nil && *params.SearchTerm != "" {
		term := "%" + *params.SearchTerm + "%"
		db = db.Joins("JOIN users ON users.id = students.user_id").
			Where("users.first_name ILIKE ? OR users.last_name ILIKE ?", term, term)
	}

	return db
}
