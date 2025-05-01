// internal/directory-service/domain/student_filter_params.go
package domain

// StudentFilterParams holds query params for GET /students
type StudentFilterParams struct {
    // latest-education filters
    CollegeName       *string  `form:"collegeName"`   // latest Education.College
    Level             *string  `form:"level"`         // latest Education.Degree
    CgpaRanges        string  `form:"cgpaRanges"` // comma-separated CGPA ranges (e.g., "6.0-7.0,7.0-8.0")
    YearOfStudy       *int     `form:"yearOfStudy"`   // latest Education.YearOfStudy
    FieldOfStudy      *string  `form:"fieldOfStudy"`  // latest Education.FieldOfStudy

    // latest-experience filters
    Company           *string  `form:"company"`       // latest Experience.Company
    Title             *string  `form:"title"`         // latest Experience.Title
    MinExperienceYears *int    `form:"minExpYears"`   // students.total_experience_years ≥
    MaxExperienceYears *int    `form:"maxExpYears"`   // students.total_experience_years ≤

    // skills filter (many-to-many)
    Skills            []string `form:"skills"`

    // general search across user names
    SearchTerm        *string  `form:"search"`

    // pagination
    Page              int      `form:"page,default=1"`
    PageSize          int      `form:"pageSize,default=20"`
}
