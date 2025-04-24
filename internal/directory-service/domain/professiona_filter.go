package domain

// ProfessionalFilterParams holds query params for /professionals
type ProfessionalFilterParams struct {
    CurrentTitle   *string   `form:"title"`
    MinExperience  *int      `form:"minExperience"`
    Skills         []string  `form:"skills"`      // comma-separated
    Industries     []string  `form:"industries"`  // comma-separated
    SearchTerm     *string   `form:"search"`      // searches firstName,lastName,currentTitle
    Page           int       `form:"page,default=1"`
    PageSize       int       `form:"pageSize,default=20"`
}
	