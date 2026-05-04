package utils

import (
	"encoding/json"
	"hrms_go/dto"
	"log"
	"math"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetQuery(queryParams *dto.PaginateFieldDto, query *gorm.DB, totalRecord *int64,totalPage *int) *gorm.DB {
	return GetQueryBase(queryParams,query,totalRecord,totalPage,nil);
}

func GetQueryBase(queryParams *dto.PaginateFieldDto, query *gorm.DB, totalRecord *int64,totalPage *int,allowedField *map[string]dto.DynamicSearchDto) *gorm.DB {
	
	if queryParams.DynamicFieldSearch != nil && allowedField != nil {
		var filters []dto.DynamicSearchFieldDto

		err := json.Unmarshal([]byte(*queryParams.DynamicFieldSearch), &filters)
		if err != nil {
			log.Println("JSON parse error:", err)
		} else {
			for _, f := range filters {
				fieldConfig, ok := (*allowedField)[f.Field]
				if !ok {
					continue // skip invalid field
				}
				log.Println("data1",fieldConfig.Field + fieldConfig.Query);
				log.Println("data2",f.Value);

				query.Where(fieldConfig.Field + fieldConfig.Query,f.Value)
			}
		}
	}


	if queryParams.SortBy != nil {
		query = query.Order(clause.OrderByColumn{
			Column: clause.Column{
				Name: queryParams.GetSortByWithDefaultId(queryParams.SortBy),
			},
			Desc: queryParams.GetSortOrderBool(),
		})
	}

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
