package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type CompaniesController struct {
	repo repositories.CompaniesRepository
}

// Create Company godoc
// @Summary Create company
// @Description Create new company
// @Tags Companies
// @Accept json
// @Produce json
// @Param request body models.Companies true "Company data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/companies [post]
func (c *CompaniesController) Create(ctx fiber.Ctx) error {
	data := models.Companies{}

	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	data.CreatedBy = ctx.Locals("user_id").(string)

	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Get All Companies godoc
// @Summary Get all companies
// @Description Get list of companies with pagination
// @Tags Companies
// @Accept json
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/companies [get]
func (c *CompaniesController) FindAll(ctx fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}

	if err := utils.BindPaginationParams(ctx, &queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	shift, err := c.repo.FindAll(&queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, shift)
}

// Update Company godoc
// @Summary Update company
// @Description Update existing company
// @Tags Companies
// @Accept json
// @Produce json
// @Param request body models.Companies true "Company data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/companies [put]
func (c *CompaniesController) Update(ctx fiber.Ctx) error {
	data := models.Companies{}

	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	userId := ctx.Locals("user_id").(string)
	data.UpdatedBy = &userId

	if err := c.repo.Update(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

func (c *CompaniesController) Delete(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.repo.Delete(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, nil)
}

func NewCompaniesController(repo repositories.CompaniesRepository) *CompaniesController {
	return &CompaniesController{repo}
}
