package dto

type PaginateFieldDto struct {
	SortOrder          *string `query:"sort_order"`
	SortBy             *string `query:"sort_by"`
	Search             *string `query:"search"`
	Page               *int    `query:"page"`
	PerPage            *int    `query:"per_page"`
	StartDate          *string `query:"start_date"`
	EndDate            *string `query:"end_date"`
	DynamicFieldSearch *string `query:"field_search"`
}

func (d *PaginateFieldDto) GetSortOrderBool() bool {
	if d.SortOrder == nil {
		return false // default ASC
	}

	switch *d.SortOrder {
	case "desc":
		return true
	case "asc":
		return false
	default:
		return false
	}
}

func (d *PaginateFieldDto) GetOffset() *int {
	if d.Page == nil || d.PerPage == nil {
		return nil
	}
	if *d.PerPage <= 0 {
		return nil
	}
	if *d.Page <= 1 {
		offset := 0
		return &offset
	}

	offset := (*d.Page - 1) * *d.PerPage

	return &offset
}

func (d *PaginateFieldDto) GetSortByWithDefaultId(defaultCol *string) string {
	if d.SortBy != nil && *d.SortBy != "" {
		return *d.SortBy
	}
	if defaultCol != nil && *defaultCol != "" {
		return *defaultCol
	}
	return ""
}
