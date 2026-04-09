package utils

import (
	"hrms_go/dto"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetQuery(queryParams *dto.PaginateFieldDto, query *gorm.DB, totalRecord *int64,totalPage *int) *gorm.DB {
	query = query.Order(clause.OrderByColumn{
		Column: clause.Column{
			Name: queryParams.GetSortByWithDefaultId(queryParams.SortBy),
		},
		Desc: queryParams.GetSortOrderBool(),
	})

	query.Count(totalRecord)
	offset := queryParams.GetOffset();

	if totalRecord != nil && queryParams.PerPage != nil && *queryParams.PerPage != 0 {
		totalPages := int(*totalRecord) / *queryParams.PerPage
		totalPage = &totalPages
	}
	
	if offset != nil {
		query = query.Limit(*queryParams.PerPage).Offset(*offset)
	}

	return query
}
