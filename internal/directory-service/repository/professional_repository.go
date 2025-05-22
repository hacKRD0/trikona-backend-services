package repository

import (
	"strconv"
	"strings"

	"github.com/hacKRD0/trikona_go/internal/directory-service/domain"
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
	var professionals []domain.Professional
	q := r.db.
		Preload("User").
		Preload("Educations", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_latest = ?", true)
		}).
		Preload("Educations.College").
		Preload("Experiences", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_latest = ?", true)
		}).
		Preload("Experiences.Company").
		Preload("Skills").
		Model(&domain.Professional{})
	q = applyProfessionalFilters(q, params)
	if err := q.Offset(offset).Limit(limit).Find(&professionals).Error; err != nil {
		return nil, err
	}
	return professionals, nil
}

func applyProfessionalFilters(db *gorm.DB, params *domain.ProfessionalFilterParams) *gorm.DB {
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
	if params.SearchTerm != nil && *params.SearchTerm != "" {
		t := "%" + *params.SearchTerm + "%"
		db = db.Where(
			"first_name ILIKE ? OR last_name ILIKE ?",
			t, t,
		)
	}
	return db
}
