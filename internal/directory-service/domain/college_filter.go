package domain

// CollegeFilterParams holds query params for /colleges
type CollegeFilterParams struct {
	// general search across user names
	SearchTerm *string `form:"searchTerm"`

	// country and state filters
	Country []string `form:"country"`
	States  []string `form:"states"`

	// student and faculty count filters
	StudentCountRanges []string `form:"studentCountRanges"`
	FacultyCountRanges []string `form:"facultyCountRanges"`

	// pagination
	Page     int `form:"page,default=1"`
	PageSize int `form:"pageSize,default=20"`
}
