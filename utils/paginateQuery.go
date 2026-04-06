package utils

import (
	"hrms_go/dto"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetQuery(queryParams *dto.PaginateFieldDto, query *gorm.DB, totalRecord *int64) *gorm.DB {
	query = query.Order(clause.OrderByColumn{
		Column: clause.Column{
			Name: queryParams.GetSortByWithDefaultId(queryParams.SortBy),
		},
		Desc: queryParams.GetSortOrderBool(),
	})

	query.Count(totalRecord)

	if offset := queryParams.GetOffset(); offset != nil {
		query = query.Limit(*queryParams.PerPage).Offset(*offset)
	}

	return query
}
