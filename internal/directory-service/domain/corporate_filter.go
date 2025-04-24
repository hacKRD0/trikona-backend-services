package domain

// CorporateFilterParams holds query params for /corporates
type CorporateFilterParams struct {
    CompanyName   *string `form:"companyName"`
    Industry      *string `form:"industry"`
    MinSize       *int    `form:"minSize"`
    Headquarters  *string `form:"hq"`
    SearchTerm    *string `form:"search"`
    Page          int     `form:"page,default=1"`
    PageSize      int     `form:"pageSize,default=20"`
}
