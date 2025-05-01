package repository

import (
	"strconv"
	"strings"
	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	"gorm.io/gorm"
)

type StudentRepository interface {
	Count(params *domain.StudentFilterParams) (int64, error)
	Find(params *domain.StudentFilterParams, offset, limit int) ([]domain.Student, error)
	GetByID(id uint) (*domain.Student, error)
	Create(student *domain.Student) error
	Update(student *domain.Student) error
	Delete(id uint) error
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

// GetByID retrieves a student by ID
func (r *studentRepository) GetByID(id uint) (*domain.Student, error) {
	var student domain.Student
	result := r.db.Preload("User").Preload("Skills").Preload("Educations").Preload("Experiences").First(&student, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &student, nil
}

// Create adds a new student to the database
func (r *studentRepository) Create(student *domain.Student) error {
	return r.db.Create(student).Error
}

// Update modifies an existing student
func (r *studentRepository) Update(student *domain.Student) error {
	return r.db.Save(student).Error
}

// Delete removes a student by ID
func (r *studentRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Student{}, id).Error
}

func applyStudentFilters(db *gorm.DB, params *domain.StudentFilterParams) *gorm.DB {
	// 1) Latest-education filters
	if params.CollegeName != nil || params.Level != nil ||
		params.CgpaRanges != "" || params.YearOfStudy != nil || params.FieldOfStudy != nil {

		db = db.Joins(`
			JOIN educations AS latest_edu
			  ON latest_edu.user_id = students.user_id
			 AND latest_edu.is_latest = TRUE
		`)

		if params.CollegeName != nil && *params.CollegeName != "" {
			db = db.Joins(`
				JOIN college_masters AS cm
				  ON cm.id = latest_edu.college_id
			`).Where("cm.name ILIKE ?", "%"+*params.CollegeName+"%")
		}
		if params.Level != nil && *params.Level != "" {
			db = db.Where("latest_edu.degree = ?", *params.Level)
		}
		if params.CgpaRanges != "" {
            // Split the ranges string by comma
            ranges := strings.Split(params.CgpaRanges, ",")
            
            // Build CGPA range conditions
            var cgpaConditions []string
            var cgpaValues []interface{}
            
            for _, r := range ranges {
                // Split each range by hyphen
                parts := strings.Split(strings.TrimSpace(r), "-")
                if len(parts) != 2 {
                    continue // Skip invalid ranges
                }
                
                // Parse min and max values
                min, err1 := strconv.ParseFloat(parts[0], 32)
                max, err2 := strconv.ParseFloat(parts[1], 32)
                if err1 != nil || err2 != nil {
                    continue // Skip invalid numbers
                }
                
                condition := "(latest_edu.cgpa >= ? AND latest_edu.cgpa <= ?)"
                cgpaConditions = append(cgpaConditions, condition)
                cgpaValues = append(cgpaValues, float32(min), float32(max))
            }
            
            if len(cgpaConditions) > 0 {
                // Combine conditions with OR
                db = db.Where(strings.Join(cgpaConditions, " OR "), cgpaValues...)
            }
		}
		if params.YearOfStudy != nil {
			db = db.Where("latest_edu.year_of_study = ?", *params.YearOfStudy)
		}
		if params.FieldOfStudy != nil && *params.FieldOfStudy != "" {
			fieldsOfStudy := strings.Split(*params.FieldOfStudy, ",")
			db = db.Where("latest_edu.field_of_study IN ?", fieldsOfStudy)
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
				JOIN company_masters AS com
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
