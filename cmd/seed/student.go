// cmd/seed/main.go
package main

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/hacKRD0/trikona_go/pkg/config"
	"github.com/hacKRD0/trikona_go/pkg/database"

	sd "github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	um "github.com/hacKRD0/trikona_go/internal/user-management-service/domain"
)

func notmain() {
	// Load environment and initialize DB
	if err := config.LoadEnv(); err != nil {
		panic("failed to load .env: " + err.Error())
	}
	db, err := database.InitDB()
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// Drop old implicit join table
	db.Migrator().DropTable("student_skill")

	// Drop existing tables
	db.Migrator().DropTable(&um.User{})
	db.Migrator().DropTable(&sd.Student{}, &sd.Experience{}, &sd.Education{})
	db.Migrator().DropTable(&sd.SkillMaster{}, &sd.CollegeMaster{}, &sd.CompanyMaster{})

	// Migrate user management tables
	db.AutoMigrate(&um.User{})

	// Migrate master tables
	db.AutoMigrate(&sd.SkillMaster{}, &sd.CollegeMaster{}, &sd.CompanyMaster{})

	// Migrate directory tables; Student (with Skills many2many creates join table) and Experience
	db.AutoMigrate(&sd.Student{}, &sd.Experience{})

	// Migrate Education
	db.AutoMigrate(&sd.Education{})

	rand.NewSource(time.Now().UnixNano())
	defaultPassword := "Password@123"

	// Data pools
	skills := []string{"Structural Analysis", "AutoCAD", "Concrete Design", "Steel Design", "Soil Mechanics", "Surveying", "Hydraulics", "Project Management", "Estimating", "Construction Materials", "Geotechnical Engineering", "Traffic Engineering"}
	colleges := []string{"Example University", "State College", "Tech Institute", "Liberal Arts College"}
	companies := []string{"Acme Corp", "BuildIt LLC", "ConstructCo", "InfraWorks", "UrbanPlan Inc"}
	titles := []string{"Intern", "Junior Engineer", "Engineer", "Senior Engineer", "Project Manager"}
	durationMonthsOpt := []int{6, 12, 24, 36, 48}

	// Seed master tables
	for _, name := range skills {
		db.Create(&sd.SkillMaster{Name: name})
	}
	for _, name := range colleges {
		db.Create(&sd.CollegeMaster{Name: name})
	}
	for _, name := range companies {
		db.Create(&sd.CompanyMaster{Name: name})
	}

	// Load master records into maps
	skillMap := make(map[string]*sd.SkillMaster)
	var skillRecs []sd.SkillMaster
	db.Find(&skillRecs)
	for i := range skillRecs {
		m := &skillRecs[i]
		skillMap[m.Name] = m
	}

	collegeMap := make(map[string]*sd.CollegeMaster)
	var collegeRecs []sd.CollegeMaster
	db.Find(&collegeRecs)
	for i := range collegeRecs {
		m := &collegeRecs[i]
		collegeMap[m.Name] = m
	}

	companyMap := make(map[string]*sd.CompanyMaster)
	var companyRecs []sd.CompanyMaster
	db.Find(&companyRecs)
	for i := range companyRecs {
		m := &companyRecs[i]
		companyMap[m.Name] = m
	}

	// Seed student profiles
	for i := 1; i <= 20; i++ {
		// Create user
		user := createUser(db, i, defaultPassword)

		// Create student record first
		student := &sd.Student{UserID: user.ID}
		db.Create(student)

		// Wait 1s to make sure student is created
		time.Sleep(time.Second)

		// Experiences & update total exp years
		exs := createExperiences(db, user.ID, companies, titles, durationMonthsOpt, companyMap)
		totalExpYears := sumExperienceYears(exs)
		student.TotalExperienceYears = totalExpYears
		db.Save(student)

		// Educations
		eds := createEducations(db, user.ID, colleges, durationMonthsOpt, collegeMap)

		// Associate random skills via implicit join table
		rand.Shuffle(len(skills), func(a, b int) { skills[a], skills[b] = skills[b], skills[a] })
		n := rand.Intn(4) + 2
		for _, name := range skills[:n] {
			if sm := skillMap[name]; sm != nil {
				db.Model(student).Association("Skills").Append(sm)
			}
		}

		fmt.Printf("Seeded %s: expYears=%d, exps=%d, edus=%d, skills=%d\n", user.Email, totalExpYears, len(exs), len(eds), n)
	}

	fmt.Println("Seeding complete!")
}

func createUser(db *gorm.DB, idx int, pwd string) *um.User {
	first := fmt.Sprintf("Student%02d", idx)
	last := fmt.Sprintf("Last%02d", idx)
	email := fmt.Sprintf("student%02d@example.com", idx)
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	user := &um.User{FirstName: first, LastName: last, Email: email, Password: string(hash), Role: um.RoleStudent, Status: um.UserStatusActive}
	db.Create(user)
	return user
}

func createExperiences(db *gorm.DB, userID uint, companies, titles []string, opts []int, companyMap map[string]*sd.CompanyMaster) []sd.Experience {
	exs := []sd.Experience{}
	count := rand.Intn(3) + 1
	for i := 0; i < count; i++ {
		start := time.Now().AddDate(-rand.Intn(5)-1, 0, 0)
		years := opts[rand.Intn(len(opts))] / 12
		end := start.AddDate(years, 0, 0)
		exs = append(exs, sd.Experience{UserID: userID, CompanyID: companyMap[companies[rand.Intn(len(companies))]].ID, Title: titles[rand.Intn(len(titles))], StartDate: start, EndDate: end, DurationMonths: years * 12, IsLatest: false})
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

func createEducations(db *gorm.DB, userID uint, colleges []string, opts []int, collegeMap map[string]*sd.CollegeMaster) []sd.Education {
	eds := []sd.Education{}
	degrees := []string{"Bachelors", "Masters", "High School", "Diploma"}
	fieldsOfStudy := []string{"Civil Engineering", "Structural Engineering", "Geotechnical Engineering", "Traffic Engineering"}
	count := rand.Intn(3) + 1
	for i := 0; i < count; i++ {
		start := time.Now().AddDate(-rand.Intn(5)-1, 0, 0)
		years := opts[rand.Intn(len(opts))] / 12
		end := start.AddDate(years, 0, 0)
		eds = append(eds, sd.Education{UserID: userID, CollegeID: collegeMap[colleges[rand.Intn(len(colleges))]].ID, Degree: degrees[rand.Intn(len(colleges))], FieldOfStudy: fieldsOfStudy[rand.Intn(len(fieldsOfStudy))], StartDate: start, EndDate: end, YearOfStudy: rand.Intn(4) + 1, CGPA: float32(rand.Intn(1001)) / 100, DurationMonths: years * 12, IsLatest: false})
	}
	// mark latest
	latestIdx := 0
	latestTime := eds[0].StartDate
	for i, ed := range eds {
		if ed.StartDate.After(latestTime) {
			latestTime = ed.StartDate
			latestIdx = i
		}
	}
	eds[latestIdx].IsLatest = true
	// persist
	for _, ed := range eds {
		db.Create(&ed)
	}
	return eds
}

func sumExperienceYears(exs []sd.Experience) int {
	total := 0
	for _, e := range exs {
		total += e.DurationMonths / 12
	}
	return total
}
