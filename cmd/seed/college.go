package main

import (
	"fmt"

	sd "github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	um "github.com/hacKRD0/trikona_go/internal/user-management-service/domain"
	"github.com/hacKRD0/trikona_go/pkg/config"
	"github.com/hacKRD0/trikona_go/pkg/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	// Load environment and initialize DB
	if err := config.LoadEnv(); err != nil {
		panic("failed to load .env: " + err.Error())
	}
	db, err := database.InitDB()
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// Drop existing tables in reverse dependency order
	db.Migrator().DropTable(&sd.Branch{}, &sd.CollegeUser{})
	db.Migrator().DropTable(&sd.College{})

	// Migrate tables
	db.AutoMigrate(&sd.College{}, &sd.Branch{}, &sd.CollegeUser{})

	// Create colleges
	if err := SeedColleges(db); err != nil {
		panic("failed to seed colleges: " + err.Error())
	}
}

func SeedColleges(db *gorm.DB) error {
	var countries []sd.CountryMaster
	var states []sd.StateMaster

	db.Find(&countries)
	db.Find(&states)
	colleges := []sd.College{
		{
			CollegeName: "Delhi University",
			Branches: []sd.Branch{
				{
					Name:         "North Campus",
					City:         "Delhi",
					Address:      "University Road, North Campus",
					PinCode:      "110007",
					Phone:        "+91 11 27667255",
					CountryID:    countries[0].ID,
					StateID:      states[0].ID,
					StudentCount: 3000,
					FacultyCount: 150,
				},
				{
					Name:         "South Campus",
					City:         "New Delhi",
					Address:      "Benito Juarez Marg, South Campus",
					PinCode:      "110021",
					Phone:        "+91 11 24119832",
					CountryID:    countries[0].ID,
					StateID:      states[0].ID,
					StudentCount: 2000,
					FacultyCount: 100,
				},
			},
		},
		{
			CollegeName: "Banaras Hindu University",
			Branches: []sd.Branch{
				{
					Name:         "Main Campus",
					City:         "Varanasi",
					Address:      "Ajagara, BHU",
					PinCode:      "221005",
					Phone:        "+91 542 2368558",
					CountryID:    countries[0].ID,
					StateID:      states[1].ID,
					StudentCount: 25000,
					FacultyCount: 1500,
				},
			},
		},
		{
			CollegeName: "Anna University",
			Branches: []sd.Branch{
				{
					Name:         "Guindy Campus",
					City:         "Chennai",
					Address:      "Sardar Patel Rd, Guindy",
					PinCode:      "600025",
					Phone:        "+91 44 22351723",
					CountryID:    countries[0].ID,
					StateID:      states[2].ID,
					StudentCount: 18000,
					FacultyCount: 1200,
				},
				{
					Name:         "MIT Campus",
					City:         "Chennai",
					Address:      "Chromepet",
					PinCode:      "600044",
					Phone:        "+91 44 22232424",
					CountryID:    countries[0].ID,
					StateID:      states[2].ID,
					StudentCount: 15000,
					FacultyCount: 900,
				},
			},
		},
		{
			CollegeName: "Jadavpur University",
			Branches: []sd.Branch{
				{
					Name:         "Main Campus",
					City:         "Kolkata",
					Address:      "188, Raja S.C. Mallick Rd, Jadavpur",
					PinCode:      "700032",
					Phone:        "+91 33 24146666",
					CountryID:    countries[0].ID,
					StateID:      states[3].ID,
					StudentCount: 12000,
					FacultyCount: 800,
				},
			},
		},
		{
			CollegeName: "Pune University",
			Branches: []sd.Branch{
				{
					Name:         "Main Campus",
					City:         "Pune",
					Address:      "Ganeshkhind, Pune University",
					PinCode:      "411007",
					Phone:        "+91 20 25690000",
					CountryID:    countries[0].ID,
					StateID:      states[4].ID,
					StudentCount: 22000,
					FacultyCount: 1400,
				},
			},
		},
		{
			CollegeName: "Osmania University",
			Branches: []sd.Branch{
				{
					Name:         "Main Campus",
					City:         "Hyderabad",
					Address:      "Osmania University Main Rd, Amberpet",
					PinCode:      "500007",
					Phone:        "+91 40 27098000",
					CountryID:    countries[0].ID,
					StateID:      states[5].ID,
					StudentCount: 20000,
					FacultyCount: 1300,
				},
			},
		},
		{
			CollegeName: "Jamia Millia Islamia",
			Branches: []sd.Branch{
				{
					Name:         "Main Campus",
					City:         "New Delhi",
					Address:      "Jamia Nagar, Okhla",
					PinCode:      "110025",
					Phone:        "+91 11 26981717",
					CountryID:    countries[0].ID,
					StateID:      states[0].ID,
					StudentCount: 15000,
					FacultyCount: 1000,
				},
			},
		},
		{
			CollegeName: "Aligarh Muslim University",
			Branches: []sd.Branch{
				{
					Name:         "North Campus",
					City:         "Delhi",
					Address:      "123 North St",
					PinCode:      "110001",
					Phone:        "+91 11 12345678",
					CountryID:    countries[0].ID,
					StateID:      states[0].ID,
					StudentCount: 3000,
					FacultyCount: 150,
				},
				{
					Name:         "South Campus",
					City:         "Delhi",
					Address:      "456 South St",
					PinCode:      "110002",
					Phone:        "+91 11 87654321",
					CountryID:    countries[0].ID,
					StateID:      states[0].ID,
					StudentCount: 2000,
					FacultyCount: 100,
				},
				{
					Name:         "Main Campus",
					City:         "Aligarh",
					CountryID:    countries[0].ID,
					StateID:      states[6].ID,
					StudentCount: 5000,
					FacultyCount: 300,
				},
			},
		},
		{
			CollegeName: "Manipal University",
			Branches: []sd.Branch{
				{
					Name:         "Main Campus",
					City:         "Manipal",
					Address:      "Madhav Nagar",
					PinCode:      "576104",
					Phone:        "+91 820 2922396",
					CountryID:    countries[0].ID,
					StateID:      states[7].ID,
					StudentCount: 25000,
					FacultyCount: 1200,
				},
			},
		},
		{
			CollegeName: "Lovely Professional University",
			Branches: []sd.Branch{
				{
					Name:         "Main Campus",
					City:         "Phagwara",
					Address:      "Jalandhar - Delhi G.T. Road",
					PinCode:      "144411",
					Phone:        "+91 1824 517001",
					CountryID:    countries[0].ID,
					StateID:      states[8].ID,
					StudentCount: 30000,
					FacultyCount: 1800,
				},
			},
		},
		{
			CollegeName: "IIT Delhi",
			Branches: []sd.Branch{
				{
					Name:         "Main Campus",
					City:         "New Delhi",
					Address:      "Hauz Khas",
					PinCode:      "110016",
					Phone:        "+91 11 26591753",
					CountryID:    countries[0].ID,
					StateID:      states[0].ID,
					StudentCount: 12000,
					FacultyCount: 800,
				},
			},
		},
		{
			CollegeName: "IIT Bombay",
			Branches: []sd.Branch{
				{
					Name:         "Main Campus",
					City:         "Mumbai",
					Address:      "Powai",
					PinCode:      "400076",
					Phone:        "+91 22 25722545",
					CountryID:    countries[0].ID,
					StateID:      states[9].ID,
					StudentCount: 13000,
					FacultyCount: 900,
				},
			},
		},
		{
			CollegeName: "IIT Madras",
			Branches: []sd.Branch{
				{
					Name:         "Main Campus",
					City:         "Chennai",
					Address:      "IIT P.O.",
					PinCode:      "600036",
					Phone:        "+91 44 22578200",
					CountryID:    countries[0].ID,
					StateID:      states[2].ID,
					StudentCount: 11000,
					FacultyCount: 850,
				},
			},
		},
	}

	for i, col := range colleges {
		user := um.User{
			FirstName: fmt.Sprintf("College%2d", i),
			LastName:  fmt.Sprintf("Last%2d", i),
			Email:     fmt.Sprintf("college%02d@example.com", i),
			Password:  "Password@123",
			Role:      "college_admin",
			Status:    "active",
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte("Password@123"), bcrypt.DefaultCost)
		user.Password = string(hash)
		if err := db.Create(&user).Error; err != nil {
			fmt.Printf("Error creating user: %v\n", err)
			continue
		}
		fmt.Printf("Created user with ID: %d\n", user.ID)
		// First, calculate total students and faculty from all branches
		totalStudents := 0
		totalFaculty := 0
		for _, branch := range col.Branches {
			totalStudents += branch.StudentCount
			totalFaculty += branch.FacultyCount
		}

		// Create college with calculated totals
		college := sd.College{
			CollegeName:  col.CollegeName,
			StudentCount: totalStudents,
			FacultyCount: totalFaculty,
		}
		if err := db.Create(&college).Error; err != nil {
			fmt.Printf("Error creating college %s: %v\n", col.CollegeName, err)
			continue
		}

		// Create branches
		for _, branch := range col.Branches {
			// Create a new branch without ID to avoid conflicts
			newBranch := sd.Branch{
				CollegeID:    college.ID,
				Name:         branch.Name,
				City:         branch.City,
				Address:      branch.Address,
				PinCode:      branch.PinCode,
				Phone:        branch.Phone,
				CountryID:    branch.CountryID,
				StateID:      branch.StateID,
				StudentCount: branch.StudentCount,
				FacultyCount: branch.FacultyCount,
			}
			if err := db.Create(&newBranch).Error; err != nil {
				fmt.Printf("Error creating branch %s: %v\n", branch.Name, err)
				continue
			}
		}

		// Create college user record
		collegeUser := sd.CollegeUser{
			CollegeID: college.ID,
			UserID:    user.ID,
			UserRole:  "college_admin",
			Status:    "active",
		}
		if err := db.Create(&collegeUser).Error; err != nil {
			fmt.Printf("Error creating college user: %v\n", err)
			continue
		}
		fmt.Printf("Created college user with ID: %d\n", collegeUser.ID)
	}
	return nil
}
