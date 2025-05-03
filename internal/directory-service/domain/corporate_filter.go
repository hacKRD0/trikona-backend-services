package domain

// CorporateFilterParams holds query params for /corporates
type CorporateFilterParams struct {
	// general search across user names
	SearchTerm *string `form:"search"`

	// country and state filters
	Country []string `form:"country"`
	States  []string `form:"states"`

	// sectors, industries and services filters
	Sectors    []string `form:"sectors"`
	Industries []string `form:"industries"`
	Services   []string `form:"services"`

	MinSize *int `form:"minSize"`

	// pagination
	Page     int `form:"page,default=1"`
	PageSize int `form:"pageSize,default=20"`
}
