package main

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	sd "github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	um "github.com/hacKRD0/trikona_go/internal/user-management-service/domain"
	"github.com/hacKRD0/trikona_go/pkg/config"
	"github.com/hacKRD0/trikona_go/pkg/database"
)

func SeedProfessionals() error {
	// Load environment and initialize DB
	if err := config.LoadEnv(); err != nil {
		return fmt.Errorf("failed to load .env: %v", err)
	}
	db, err := database.InitDB()
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	// Drop existing tables in reverse order to avoid foreign key constraints
	db.Migrator().DropTable("professional_skill")
	db.Migrator().DropTable(&sd.Professional{})

	// Migrate Professional and related tables
	if err := db.AutoMigrate(&sd.Professional{}); err != nil {
		return fmt.Errorf("error migrating professional table: %w", err)
	}

	rand.NewSource(time.Now().UnixNano())
	defaultPassword := "Password@123"

	// Fetch skills from master table
	var skillRecords []sd.SkillMaster
	if err := db.Find(&skillRecords).Error; err != nil {
		return fmt.Errorf("error fetching skills: %w", err)
	}

	// Fetch colleges from master table
	var collegeRecords []sd.CollegeMaster
	if err := db.Find(&collegeRecords).Error; err != nil {
		return fmt.Errorf("error fetching colleges: %w", err)
	}

	// Fetch companies from master table
	var companyRecords []sd.CompanyMaster
	if err := db.Find(&companyRecords).Error; err != nil {
		return fmt.Errorf("error fetching companies: %w", err)
	}

	// Create maps for quick lookup
	skillMap := make(map[string]*sd.SkillMaster, len(skillRecords))
	for i := range skillRecords {
		skillMap[skillRecords[i].Name] = &skillRecords[i]
	}

	collegeMap := make(map[string]*sd.CollegeMaster, len(collegeRecords))
	for i := range collegeRecords {
		collegeMap[collegeRecords[i].Name] = &collegeRecords[i]
	}

	companyMap := make(map[string]*sd.CompanyMaster, len(companyRecords))
	for i := range companyRecords {
		companyMap[companyRecords[i].Name] = &companyRecords[i]
	}

	titles := []string{"Project Manager", "Senior Project Manager", "Construction Manager", "Planning Engineer", "Project Director"}
	durationMonthsOpt := []int{12, 24, 36, 48, 60}

	// Get skill names for random selection
	skillNames := make([]string, 0, len(skillRecords))
	for _, s := range skillRecords {
		skillNames = append(skillNames, s.Name)
	}

	// Get college names for random selection
	collegeNames := make([]string, 0, len(collegeRecords))
	for _, c := range collegeRecords {
		collegeNames = append(collegeNames, c.Name)
	}

	// Get company names for random selection
	companyNames := make([]string, 0, len(companyRecords))
	for _, c := range companyRecords {
		companyNames = append(companyNames, c.Name)
	}

	// Seed professional profiles
	for i := 1; i <= 1; i++ {
		// Create user
		user := createProfessionalUser(db, i, defaultPassword)

		// Create professional record first
		professional := &sd.Professional{UserID: user.ID}
		if result := db.Create(professional); result.Error != nil {
			return fmt.Errorf("error creating professional: %w", result.Error)
		}

		// Small delay to ensure records are created
		time.Sleep(100 * time.Millisecond)

		// Experiences & update total exp years
		exs := createProfessionalExperiences(db, user.ID, companyNames, titles, durationMonthsOpt, companyMap)
		totalExpYears := sumProfessionalExperienceYears(exs)
		professional.TotalExperienceYears = totalExpYears
		if result := db.Save(professional); result.Error != nil {
			return fmt.Errorf("error updating professional: %w", result.Error)
		}

		// Educations
		edus := createProfessionalEducations(db, user.ID, collegeNames, durationMonthsOpt, collegeMap)

		// Associate random skills via implicit join table
		rand.Shuffle(len(skillNames), func(a, b int) { skillNames[a], skillNames[b] = skillNames[b], skillNames[a] })
		n := rand.Intn(5) + 1 // 1-5 skills
		for _, skillName := range skillNames[:n] {
			if skill, exists := skillMap[skillName]; exists {
				if err := db.Model(professional).Association("Skills").Append(skill); err != nil {
					return fmt.Errorf("error associating skill: %w", err)
				}
			}
		}

		fmt.Printf("Seeded Professional %s: expYears=%d, exps=%d, edus=%d, skills=%d\n",
			user.Email, totalExpYears, len(exs), len(edus), n)
	}

	fmt.Println("Professional seeding complete!")
	return nil
}

func createProfessionalUser(db *gorm.DB, idx int, pwd string) *um.User {
	first := fmt.Sprintf("Professional%02d", idx)
	last := fmt.Sprintf("Last%02d", idx)
	email := fmt.Sprintf("professional%02d@example.com", idx)
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	user := &um.User{
		FirstName: first,
		LastName:  last,
		Email:     email,
		Password:  string(hash),
		Role:      um.RoleProfessional,
		Status:    um.UserStatusActive,
	}
	db.Create(user)
	return user
}

func createProfessionalExperiences(db *gorm.DB, userID uint, companies, titles []string, opts []int, companyMap map[string]*sd.CompanyMaster) []sd.Experience {
	exs := []sd.Experience{}
	count := rand.Intn(3) + 1
	for i := 0; i < count; i++ {
		start := time.Now().AddDate(-rand.Intn(10)-5, 0, 0) // 5-15 years experience
		years := opts[rand.Intn(len(opts))] / 12
		end := start.AddDate(years, 0, 0)
		exs = append(exs, sd.Experience{
			UserID:         userID,
			CompanyID:      companyMap[companies[rand.Intn(len(companies))]].ID,
			Title:          titles[rand.Intn(len(titles))],
			StartDate:      start,
			EndDate:        end,
			DurationMonths: years * 12,
			IsLatest:       false,
		})
	}
	// mark latest
	latestIdx := 0
	latestTime := exs[0].StartDate
	for i, e := range exs {
		if e.StartDate.After(latestTime) {
			latestTime = e.StartDate
			latestIdx = i
		}
	}
	exs[latestIdx].IsLatest = true
	// persist
	for _, e := range exs {
		db.Create(&e)
	}
	return exs
}

func createProfessionalEducations(db *gorm.DB, userID uint, colleges []string, opts []int, collegeMap map[string]*sd.CollegeMaster) []sd.Education {
	edus := []sd.Education{}
	degrees := []string{"Bachelors", "Masters", "High School", "Diploma"}
	fieldsOfStudy := []string{"Civil Engineering", "Structural Engineering", "Geotechnical Engineering", "Traffic Engineering"}

	count := rand.Intn(2) + 1 // 1-2 education records
	for i := 0; i < count; i++ {
		start := time.Now().AddDate(-rand.Intn(15)-5, 0, 0) // 5-20 years ago
		years := opts[rand.Intn(len(opts))] / 12
		end := start.AddDate(years, 0, 0)
		edus = append(edus, sd.Education{
			UserID:         userID,
			CollegeID:      collegeMap[colleges[rand.Intn(len(colleges))]].ID,
			Degree:         degrees[rand.Intn(len(degrees))],
			FieldOfStudy:   fieldsOfStudy[rand.Intn(len(fieldsOfStudy))],
			StartDate:      start,
			EndDate:        end,
			YearOfStudy:    rand.Intn(4) + 1,
			CGPA:           float32(rand.Intn(401)+600) / 100, // 6.0-10.0
			DurationMonths: years * 12,
			IsLatest:       false,
		})
	}
	// mark latest
	latestIdx := 0
	latestTime := edus[0].StartDate
	for i, ed := range edus {
		if ed.StartDate.After(latestTime) {
			latestTime = ed.StartDate
			latestIdx = i
		}
	}
	edus[latestIdx].IsLatest = true
	// persist
	for _, ed := range edus {
		db.Create(&ed)
	}
	return edus
}

func sumProfessionalExperienceYears(exs []sd.Experience) int {
	total := 0
	for _, e := range exs {
		total += e.DurationMonths / 12
	}
	return total
}
