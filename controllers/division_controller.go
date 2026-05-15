package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type DivisionController struct {
	repo repositories.DivisionRepository
}

func NewDivisionController(repo repositories.DivisionRepository) *DivisionController {
	return &DivisionController{repo}
}

// Create Division godoc
// @Summary Create division
// @Description Create new division
// @Tags Division
// @Accept json
// @Produce json
// @Param request body models.Division true "Division data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/division [post]
func (c *DivisionController) Create(ctx fiber.Ctx) error {
	var data models.Division
	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	data.CreatedBy = ctx.Locals("user_id").(string)
	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

// Get All Division godoc
// @Summary Get all divisions
// @Description Get list of divisions with pagination
// @Tags Division
// @Accept json
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param field_search query string false "Field name for dynamic search"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/division [get]
func (c *DivisionController) FindAll(ctx fiber.Ctx) error {
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

// Update Division godoc
// @Summary Update division
// @Description Update existing division
// @Tags Division
// @Accept json
// @Produce json
// @Param request body models.Division true "Division data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/division [put]
func (c *DivisionController) Update(ctx fiber.Ctx) error {
	var data models.Division
	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	userID := ctx.Locals("user_id").(string)
	data.UpdatedBy = &userID
	if err := c.repo.Update(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

// Delete Division godoc
// @Summary Delete division
// @Description Delete division by ID
// @Tags Division
// @Accept json
// @Produce json
// @Param id path string true "Division ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/division/{id} [delete]
func (c *DivisionController) Delete(ctx fiber.Ctx) error {
	if err := c.repo.Delete(ctx.Params("id")); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, nil)
}
