package domain

// CollegeFilterParams holds query params for /colleges
type CollegeFilterParams struct {
    CollegeName   *string  `form:"collegeName"`
    Location      *string  `form:"location"`
    Accreditation *string  `form:"accreditation"`
    Departments   []string `form:"departments"` // comma-separated
    SearchTerm    *string  `form:"search"`
    Page          int      `form:"page,default=1"`
    PageSize      int      `form:"pageSize,default=20"`
}
