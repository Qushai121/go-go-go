package utils

import (
	"hrms_go/dto"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func BindPaginationParams(ctx fiber.Ctx, queryParams *dto.PaginateFieldDto) error {
	if err := ctx.Bind().Query(queryParams); err != nil {
		return err
	}

	if raw := ctx.Get("pageNo"); raw != "" {
		if value, err := strconv.Atoi(raw); err == nil && value > 0 {
			queryParams.Page = &value
		}
	}
	if raw := ctx.Get("pageSize"); raw != "" {
		if value, err := strconv.Atoi(raw); err == nil && value > 0 {
			queryParams.PerPage = &value
		}
	}

	if queryParams.Page == nil {
		if raw := ctx.Get("page"); raw != "" {
			if value, err := strconv.Atoi(raw); err == nil && value > 0 {
				queryParams.Page = &value
			}
		}
	}

	if queryParams.PerPage == nil {
		if raw := ctx.Get("per_page"); raw != "" {
			if value, err := strconv.Atoi(raw); err == nil && value > 0 {
				queryParams.PerPage = &value
			}
		}
	}

	queryParams.Filters = ctx.Queries()

	return nil
}

func GetQuery(queryParams *dto.PaginateFieldDto, query *gorm.DB, totalRecord *int64, totalPage *int) *gorm.DB {

	if queryParams.SortBy != nil {
		query = query.Order(clause.OrderByColumn{
			Column: clause.Column{
				Name: queryParams.GetSortByWithDefaultId(queryParams.SortBy),
			},
			Desc: queryParams.GetSortOrderBool(),
		})
	}

	query.Count(totalRecord)
	offset := queryParams.GetOffset()

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

func GetQueryBase(queryParams *dto.PaginateFieldDto, query *gorm.DB, totalRecord *int64, totalPage *int, allowedDynamicList *map[string]dto.DynamicSearchDto) *gorm.DB {
	if queryParams.SortBy != nil {
		query = query.Order(clause.OrderByColumn{
			Column: clause.Column{
				Name: queryParams.GetSortByWithDefaultId(queryParams.SortBy),
			},
			Desc: queryParams.GetSortOrderBool(),
		})
	}

	if allowedDynamicList != nil && queryParams.DynamicFieldSearch != nil && *queryParams.DynamicFieldSearch != "" && queryParams.Search != nil && *queryParams.Search != "" {
		if dynamicField, ok := (*allowedDynamicList)[*queryParams.DynamicFieldSearch]; ok {
			query = applyDynamicSearch(query, dynamicField, *queryParams.Search)
		}
	}

	if allowedDynamicList != nil && len(queryParams.Filters) > 0 {
		for fieldName, dynamicField := range *allowedDynamicList {
			if value := queryParams.Filters[fieldName]; value != "" {
				query = applyDynamicSearch(query, dynamicField, value)
			}
		}
	}

	query.Count(totalRecord)
	offset := queryParams.GetOffset()

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

func applyDynamicSearch(query *gorm.DB, dynamicField dto.DynamicSearchDto, value string) *gorm.DB {
	searchValue := value
	if strings.Contains(strings.ToUpper(dynamicField.Query), "LIKE") {
		searchValue = "%" + searchValue + "%"
	}

	return query.Where(dynamicField.Field+dynamicField.Query, searchValue)
}
