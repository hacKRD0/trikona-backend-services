// pkg/pagination/paginator.go
package pagination

// CalculateOffsetLimit converts page and pageSize into SQL OFFSET and LIMIT.
// page is 1-based; pageSize is the number of items per page.
func CalculateOffsetLimit(page, pageSize int) (offset, limit int) {
    if page < 1 {
        page = 1
    }
    if pageSize <= 0 {
        pageSize = 20
    }
    offset = (page - 1) * pageSize
    limit = pageSize
    return
}
