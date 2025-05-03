// cmd/seed/corporate.go
package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/hacKRD0/trikona_go/pkg/config"
	"github.com/hacKRD0/trikona_go/pkg/database"

	sd "github.com/hacKRD0/trikona_go/internal/directory-service/domain"
	um "github.com/hacKRD0/trikona_go/internal/user-management-service/domain"
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
	db.Migrator().DropTable("corporate_industries", "corporate_services", "corporate_sectors")
	db.Migrator().DropTable(&sd.Office{})
	db.Migrator().DropTable(&sd.Corporate{})
	// db.Migrator().DropTable(&um.User{})
	db.Migrator().DropTable(&sd.IndustryMaster{}, &sd.ServiceMaster{}, &sd.SectorMaster{})
	db.Migrator().DropTable(&sd.CountryMaster{}, &sd.StateMaster{})

	// Migrate user management tables
	// db.AutoMigrate(&um.User{})

	// Migrate master tables
	db.AutoMigrate(&sd.IndustryMaster{}, &sd.ServiceMaster{}, &sd.SectorMaster{})
	db.AutoMigrate(&sd.CountryMaster{}, &sd.StateMaster{})

	// Migrate corporate tables
	db.AutoMigrate(&sd.Corporate{}, &sd.Office{}, &sd.CorporateUser{})

	rand.NewSource(time.Now().UnixNano())
	defaultPassword := "Password@123"

	// Seed master data
	countries := seedCountries(db)
	states := seedStates(db)
	industries := seedIndustries(db)
	services := seedServices(db)
	sectors := seedSectors(db)

	// Seed corporate data
	seedCorporates(db, defaultPassword, industries, services, sectors, countries, states)

	fmt.Println("Corporate seed data created successfully!")
}

func seedIndustries(db *gorm.DB) []sd.IndustryMaster {
	industries := []sd.IndustryMaster{
		{Name: "Agriculture"},
		{Name: "Oil and Gas"},
		{Name: "Forestry"},
		{Name: "Chemical"},
		{Name: "Mechanical"},
		{Name: "Metal"},
		{Name: "Commerce"},
		{Name: "Construction"},
		{Name: "Hotels"},
	}

	for i := range industries {
		db.Create(&industries[i])
	}
	return industries
}

func seedServices(db *gorm.DB) []sd.ServiceMaster {
	services := []sd.ServiceMaster{
		{Name: "Earth Work"},
		{Name: "Form Work"},
		{Name: "Rebar"},
		{Name: "Mep"},
		{Name: "Concrete"},
		{Name: "Piling"},
		{Name: "Roadmarking"},
		{Name: "RMC"},
		{Name: "Landscaping"},
		{Name: "Masonry"},
		{Name: "Waterproofing"},
	}

	for i := range services {
		db.Create(&services[i])
	}
	return services
}

func seedSectors(db *gorm.DB) []sd.SectorMaster {
	sectors := []sd.SectorMaster{
		{Name: "Energy"},
		{Name: "Mining"},
		{Name: "Manufacturing"},
		{Name: "Services"},
		{Name: "Technology"},
		{Name: "Finance"},
		{Name: "Healthcare"},
		{Name: "Education"},
		{Name: "Retail"},
		{Name: "Transportation"},
	}

	for i := range sectors {
		db.Create(&sectors[i])
	}
	return sectors
}

func seedCountries(db *gorm.DB) []sd.CountryMaster {
	countries := []sd.CountryMaster{
		{Name: "India", ISOCode: "IN"},
	}

	for i := range countries {
		db.Create(&countries[i])
	}
	return countries
}

func seedStates(db *gorm.DB) []sd.StateMaster {
	states := []sd.StateMaster{
		// India States and Union Territories (CountryID: 1)
		{Name: "Andaman and Nicobar Islands", CountryID: 1},
		{Name: "Andhra Pradesh", CountryID: 1},
		{Name: "Arunachal Pradesh", CountryID: 1},
		{Name: "Assam", CountryID: 1},
		{Name: "Bihar", CountryID: 1},
		{Name: "Chandigarh", CountryID: 1},
		{Name: "Chhattisgarh", CountryID: 1},
		{Name: "Dadra and Nagar Haveli and Daman and Diu", CountryID: 1},
		{Name: "Delhi", CountryID: 1},
		{Name: "Goa", CountryID: 1},
		{Name: "Gujarat", CountryID: 1},
		{Name: "Haryana", CountryID: 1},
		{Name: "Himachal Pradesh", CountryID: 1},
		{Name: "Jammu and Kashmir", CountryID: 1},
		{Name: "Jharkhand", CountryID: 1},
		{Name: "Karnataka", CountryID: 1},
		{Name: "Kerala", CountryID: 1},
		{Name: "Ladakh", CountryID: 1},
		{Name: "Lakshadweep", CountryID: 1},
		{Name: "Madhya Pradesh", CountryID: 1},
		{Name: "Maharashtra", CountryID: 1},
		{Name: "Manipur", CountryID: 1},
		{Name: "Meghalaya", CountryID: 1},
		{Name: "Mizoram", CountryID: 1},
		{Name: "Nagaland", CountryID: 1},
		{Name: "Odisha", CountryID: 1},
		{Name: "Puducherry", CountryID: 1},
		{Name: "Punjab", CountryID: 1},
		{Name: "Rajasthan", CountryID: 1},
		{Name: "Sikkim", CountryID: 1},
		{Name: "Tamil Nadu", CountryID: 1},
		{Name: "Telangana", CountryID: 1},
		{Name: "Tripura", CountryID: 1},
		{Name: "Uttar Pradesh", CountryID: 1},
		{Name: "Uttarakhand", CountryID: 1},
		{Name: "West Bengal", CountryID: 1},
	}

	for i := range states {
		db.Create(&states[i])
	}
	return states
}

func seedCorporates(db *gorm.DB, pwd string, industries []sd.IndustryMaster, services []sd.ServiceMaster, sectors []sd.SectorMaster, countries []sd.CountryMaster, states []sd.StateMaster) {
	companies := []struct {
		Name       string
		HeadCount  int
		Industries []sd.IndustryMaster
		Services   []sd.ServiceMaster
		Sectors    []sd.SectorMaster
		Offices    []sd.Office
	}{
		{
			Name:       "TechCorp Solutions",
			HeadCount:  5000,
			Industries: industries[:2],
			Services:   services[:3],
			Sectors:    sectors[:2],
			Offices: []sd.Office{
				{
					Name:      "Delhi Office",
					Address:   "123 Tech Street",
					City:      "Delhi",
					StateID:   states[8].ID,
					CountryID: countries[0].ID,
					PinCode:   "110001",
					Phone:     "+91 11 1234 5678",
					HeadCount: 2000,
				}, {
					Name:      "Mumbai Office",
					Address:   "456 Innovation Ave",
					City:      "Mumbai",
					StateID:   states[20].ID,
					CountryID: countries[0].ID,
					PinCode:   "400001",
					Phone:     "+91 22 1234 5678",
					HeadCount: 3000,
				},
			},
		},
		{
			Name:       "FinTech Innovations",
			HeadCount:  2500,
			Industries: []sd.IndustryMaster{industries[1], industries[2]},
			Services:   []sd.ServiceMaster{services[1], services[2]},
			Sectors:    []sd.SectorMaster{sectors[1], sectors[2]},
			Offices: []sd.Office{
				{
					Name:      "Hyderabad Office",
					Address:   "789 Finance Lane",
					City:      "Hyderabad",
					StateID:   states[28].ID,
					CountryID: countries[0].ID,
					PinCode:   "500001",
					Phone:     "+91 40 1234 5678",
					HeadCount: 2500,
				},
			},
		},
		{
			Name:       "EduLearn Systems",
			HeadCount:  1800,
			Industries: industries[3:5],
			Services:   services[2:5],
			Sectors:    sectors[3:5],
			Offices: []sd.Office{
				{
					Name:      "Bangalore Campus",
					Address:   "88 Knowledge Park",
					City:      "Bangalore",
					StateID:   states[1].ID,
					CountryID: countries[0].ID,
					PinCode:   "560001",
					Phone:     "+91 80 2345 6789",
					HeadCount: 1800,
				},
			},
		},
		{
			Name:       "GreenEnergy Solutions",
			HeadCount:  2600,
			Industries: []sd.IndustryMaster{industries[4], industries[5]},
			Services:   services[3:6],
			Sectors:    sectors[5:7],
			Offices: []sd.Office{
				{
					Name:      "Ahmedabad HQ",
					Address:   "44 Solar Rd",
					City:      "Ahmedabad",
					StateID:   states[5].ID,
					CountryID: countries[0].ID,
					PinCode:   "380001",
					Phone:     "+91 79 1234 5678",
					HeadCount: 2600,
				},
			},
		},
		{
			Name:       "RetailHub Enterprises",
			HeadCount:  1500,
			Industries: []sd.IndustryMaster{industries[6]},
			Services:   services[5:7],
			Sectors:    sectors[7:9],
			Offices: []sd.Office{
				{
					Name:      "Chennai Retail Center",
					Address:   "12 Bazaar St",
					City:      "Chennai",
					StateID:   states[2].ID,
					CountryID: countries[0].ID,
					PinCode:   "600001",
					Phone:     "+91 44 2345 6789",
					HeadCount: 1500,
				},
			},
		},
		{
			Name:       "BuildRight Constructions",
			HeadCount:  2200,
			Industries: []sd.IndustryMaster{industries[7]},
			Services:   services[:3],
			Sectors:    sectors[0:2],
			Offices: []sd.Office{
				{
					Name:      "Jaipur HQ",
					Address:   "88 Fort Road",
					City:      "Jaipur",
					StateID:   states[6].ID,
					CountryID: countries[0].ID,
					PinCode:   "302001",
					Phone:     "+91 141 1234 5678",
					HeadCount: 2200,
				},
			},
		},
		{
			Name:       "AutoDrive Motors",
			HeadCount:  4000,
			Industries: []sd.IndustryMaster{industries[2], industries[5]},
			Services:   services[2:5],
			Sectors:    sectors[2:4],
			Offices: []sd.Office{
				{
					Name:      "Pune Plant",
					Address:   "77 Driveway Blvd",
					City:      "Pune",
					StateID:   states[20].ID,
					CountryID: countries[0].ID,
					PinCode:   "411001",
					Phone:     "+91 20 1234 5678",
					HeadCount: 4000,
				},
			},
		},
		{
			Name:       "HealthPlus MedTech",
			HeadCount:  1800,
			Industries: []sd.IndustryMaster{industries[3]},
			Services:   services[1:4],
			Sectors:    sectors[6:8],
			Offices: []sd.Office{
				{
					Name:      "Lucknow Unit",
					Address:   "59 Wellness Way",
					City:      "Lucknow",
					StateID:   states[9].ID,
					CountryID: countries[0].ID,
					PinCode:   "226001",
					Phone:     "+91 522 1234 5678",
					HeadCount: 1800,
				},
			},
		},
		{
			Name:       "CloudNet Services",
			HeadCount:  3500,
			Industries: []sd.IndustryMaster{industries[0], industries[6]},
			Services:   services[3:6],
			Sectors:    sectors[4:6],
			Offices: []sd.Office{
				{
					Name:      "Bhubaneswar Cloud Hub",
					Address:   "144 Data Park",
					City:      "Bhubaneswar",
					StateID:   states[12].ID,
					CountryID: countries[0].ID,
					PinCode:   "751001",
					Phone:     "+91 674 1234 5678",
					HeadCount: 3500,
				},
			},
		},
		{
			Name:       "AgroGrow Industries",
			HeadCount:  3000,
			Industries: []sd.IndustryMaster{industries[0]},
			Services:   services[:2],
			Sectors:    sectors[:2],
			Offices: []sd.Office{
				{
					Name:      "Patna Agri Park",
					Address:   "199 Harvest Ave",
					City:      "Patna",
					StateID:   states[11].ID,
					CountryID: countries[0].ID,
					PinCode:   "800001",
					Phone:     "+91 612 1234 5678",
					HeadCount: 3000,
				},
			},
		},
	}

	for i, company := range companies {
		// Create user
		user := um.User{
			FirstName: fmt.Sprintf("Corporate%2d", i),
			LastName:  fmt.Sprintf("Last%2d", i),
			Email:     fmt.Sprintf("%s-admin@%s.com", strings.ToLower(company.Name), strings.ToLower(company.Name)),
			Password:  pwd,
			Role:      "corporate_admin",
			Status:    "active",
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		user.Password = string(hash)
		if err := db.Create(&user).Error; err != nil {
			fmt.Printf("Error creating user: %v\n", err)
			continue
		}
		fmt.Printf("Created user with ID: %d\n", user.ID)

		// Create corporate
		corporate := sd.Corporate{
			CompanyName: company.Name,
			HeadCount:   company.HeadCount,
			Industries:  company.Industries,
			Services:    company.Services,
			Sectors:     company.Sectors,
		}
		db.Create(&corporate)

		// Create offices
		for _, office := range company.Offices {
			off := sd.Office{
				CorporateID: corporate.ID,
				Name:        office.Name,
				Address:     office.Address,
				City:        office.City,
				StateID:     office.StateID,
				CountryID:   office.CountryID,
				PinCode:     office.PinCode,
				Phone:       office.Phone,
				HeadCount:   office.HeadCount,
			}
			db.Create(&off)
		}

		// Create corporate user record
		corporateUser := sd.CorporateUser{
			CorporateID: corporate.ID,
			UserID:      user.ID,
			UserRole:    "corporate_admin",
			Status:      "active",
		}
		if err := db.Create(&corporateUser).Error; err != nil {
			fmt.Printf("Error creating corporate user: %v\n", err)
			continue
		}
		fmt.Printf("Created corporate user with ID: %d\n", corporateUser.ID)
	}
}

// func seedCorporateUsers(db *gorm.DB, pwd string, industries []sd.IndustryMaster, services []sd.ServiceMaster, sectors []sd.SectorMaster, countries []sd.CountryMaster, states []sd.StateMaster) {
// 	// Create corporate users
// 	corporateUsers := []sd.CorporateUser{
// 		{
// 			CorporateID: 1,
// 			User: um.User{
// 				FirstName: "John",
// 				LastName:  "Doe",
// 				Email:     "john.doe@company1.com",
// 				Role:      um.RoleCorporateAdmin,
// 				Status:    um.UserStatusActive,
// 			},
// 		},
// 		{
// 			CorporateID: 1,
// 			User: um.User{
// 				FirstName: "Jane",
// 				LastName:  "Smith",
// 				Email:     "jane.smith@company1.com",
// 				Role:      um.RoleCorporateModerator,
// 				Status:    um.UserStatusActive,
// 			},
// 		},
// 		{
// 			CorporateID: 2,
// 			User: um.User{
// 				FirstName: "Mike",
// 				LastName:  "Johnson",
// 				Email:     "mike.johnson@company2.com",
// 				Role:      um.RoleCorporateAdmin,
// 				Status:    um.UserStatusActive,
// 			},
// 		},
// 	}

// 	// Create users first
// 	for _, user := range corporateUsers {
// 		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
// 		user.User.Password = string(hashedPassword)
// 		db.Create(&user.User)
// 	}

// 	// Now create corporate users with foreign key references
// 	for i, user := range corporateUsers {
// 		var dbUser um.User
// 		db.First(&dbUser, "email = ?", user.User.Email)
// 		corporateUsers[i].UserID = dbUser.ID
// 		db.Create(&corporateUsers[i])
// 	}
// }
