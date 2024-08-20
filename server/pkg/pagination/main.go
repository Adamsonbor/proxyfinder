package pagination

type Options struct {
	PerPage int
	Page    int
}

// Offset calculates the offset based on page and perPage
// If page or perPage are less than 1, they will be set to 1 and 10 respectively
func Offset(page int, perPage int) int {
	if perPage < 1 {
		perPage = 10
	}
	if page < 1 {
		page = 1
	}

	return (page - 1) * perPage
}

// Return perPage or 10 if it is less than 1
func Limit(perPage int) int {
	if perPage < 1 {
		perPage = 10
	}
	return perPage
}

// LimitOffset returns the offset and limit based on page and perPage
// If page or perPage are less than 1, they will be set to 1 and 10 respectively
func LimitOffset(page int, perPage int) (int, int) {
	return Limit(perPage), Offset(page, perPage)
}

// PageCount returns the number of pages based on total and perPage
func PageCount(total int, perPage int) int {
	return (total + perPage - 1) / perPage
}

// ParsePagination returns the offset and limit based on the map values
func ParsePagination(o map[string]interface{}) (int, int) {
	page, _ := o["page"].(int)
	perPage, _ := o["perPage"].(int)
	return LimitOffset(page, perPage)
}
