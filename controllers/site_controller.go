package controllers

import (
	"strings"

	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type SiteController struct {
	repo repositories.SiteRepository
}

func NewSiteController(repo repositories.SiteRepository) *SiteController {
	return &SiteController{repo}
}

// Create Site godoc
// @Summary Create site
// @Description Create new site. site_type can be BRANCH, OFFICE, or CUSTOMER.
// @Tags Site
// @Accept json
// @Produce json
// @Param request body models.Site true "Site data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/site [post]
func (c *SiteController) Create(ctx fiber.Ctx) error {
	var data models.Site
	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	data.SiteType = strings.ToUpper(strings.TrimSpace(data.SiteType))
	if !isValidSiteType(data.SiteType) {
		return utils.Error(ctx, 400, "site_type must be BRANCH, OFFICE, or CUSTOMER")
	}
	data.CreatedBy = ctx.Locals("employee_nik").(string)

	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

// Get All Site godoc
// @Summary Get all sites
// @Description Get list of sites with pagination. Use search/site_type/company_code/site_code/site_name query for lookup.
// @Tags Site
// @Accept json
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param field_search query string false "Field name for dynamic search"
// @Param company_code query string false "Filter by company code"
// @Param site_type query string false "Filter by site type"
// @Param site_code query string false "Filter by site code"
// @Param site_name query string false "Filter by site name"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/site [get]
func (c *SiteController) FindAll(ctx fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}
	if err := utils.BindPaginationParams(ctx, &queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	data, err := c.repo.FindAll(&queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

// Search Site godoc
// @Summary Search sites
// @Description Alias endpoint for site lookup/search.
// @Tags Site
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /api/site/search [get]
func (c *SiteController) Search(ctx fiber.Ctx) error {
	return c.FindAll(ctx)
}

// Update Site godoc
// @Summary Update site
// @Description Update existing site
// @Tags Site
// @Accept json
// @Produce json
// @Param request body models.Site true "Site data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/site [put]
func (c *SiteController) Update(ctx fiber.Ctx) error {
	var data models.Site
	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	data.SiteType = strings.ToUpper(strings.TrimSpace(data.SiteType))
	if data.SiteType != "" && !isValidSiteType(data.SiteType) {
		return utils.Error(ctx, 400, "site_type must be BRANCH, OFFICE, or CUSTOMER")
	}
	userID := ctx.Locals("employee_nik").(string)
	data.UpdatedBy = &userID

	if err := c.repo.Update(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

// Delete Site godoc
// @Summary Delete site
// @Description Delete site by ID
// @Tags Site
// @Accept json
// @Produce json
// @Param id path string true "Site ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/site/{id} [delete]
func (c *SiteController) Delete(ctx fiber.Ctx) error {
	if err := c.repo.Delete(ctx.Params("id")); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, nil)
}

func isValidSiteType(siteType string) bool {
	switch strings.ToUpper(strings.TrimSpace(siteType)) {
	case "BRANCH", "OFFICE", "CUSTOMER":
		return true
	default:
		return false
	}
}
