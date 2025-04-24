// internal/directory-service/domain/college_master.go
package domain

// CollegeMaster holds the canonical list of colleges
// +-----------+      +--------------------+      +---------------+
// | students  |<---->| student_education  |<---->| college_master |
// +-----------+      +--------------------+      +---------------+

// CollegeMaster is the lookup for valid colleges
// referenced by Education.CollegeID

type CollegeMaster struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:255;unique;not null" json:"name"`
}
