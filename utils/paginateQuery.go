package utils

import (
	"hrms_go/dto"
	"log"
	"math"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetQuery(queryParams *dto.PaginateFieldDto, query *gorm.DB, totalRecord *int64,totalPage *int) *gorm.DB {
	log.Println(queryParams.GetSortByWithDefaultId(queryParams.SortBy))
	log.Println(queryParams.GetSortOrderBool())

	query = query.Order(clause.OrderByColumn{
		Column: clause.Column{
			Name: queryParams.GetSortByWithDefaultId(queryParams.SortBy),
		},
		Desc: queryParams.GetSortOrderBool(),
	})

	query.Count(totalRecord)
	offset := queryParams.GetOffset();

	if totalRecord != nil && totalPage != nil && queryParams.PerPage != nil && *queryParams.PerPage > 0 {
		*totalPage = int(math.Ceil(float64(*totalRecord) / float64(*queryParams.PerPage)))
	}

	if queryParams.StartDate != nil && queryParams.EndDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", *queryParams.StartDate, *queryParams.EndDate)
	} else if queryParams.StartDate != nil {
		query = query.Where("created_at >= ?", *queryParams.StartDate)
	} else if queryParams.EndDate != nil {
		query = query.Where("created_at <= ?", *queryParams.EndDate)
	}
	
	if offset != nil && queryParams.PerPage != nil {
		query = query.Limit(*queryParams.PerPage).Offset(*offset)
	}

	return query
}
