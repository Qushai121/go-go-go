package dto

type PaginateFieldDto struct {
	SortOrder *string
	SortBy    *string
	Search    *string
	Page      *int
	PerPage   *int
}

func (d *PaginateFieldDto) GetSortOrderBool() bool {
	if d.SortOrder == nil {
		return true
	}
	return *d.SortOrder == "desc"
}

func (d *PaginateFieldDto) GetOffset() *int {
	if d.Page == nil || d.PerPage == nil {
		return nil
	}
	offset := (*d.Page - 1) * *d.PerPage

	return &offset
}

func (d *PaginateFieldDto) GetSortByWithDefaultId(sortBy *string) string {
	defaultSortBy := ""

	if sortBy != nil {
		defaultSortBy = *sortBy
	}

	return defaultSortBy
}