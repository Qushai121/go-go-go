package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type ParamGroupController struct {
	repo repositories.ParamGroupRepository
}

func NewParamGroupController(repo repositories.ParamGroupRepository) *ParamGroupController {
	return &ParamGroupController{repo}
}

// Create Param Group godoc
// @Summary Create param group
// @Description Create new param group
// @Tags Param Group
// @Accept json
// @Produce json
// @Param request body models.ParamGroup true "Param Group data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/param-group [post]
func (c *ParamGroupController) Create(ctx fiber.Ctx) error {
	data := models.ParamGroup{}

	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	userId := ctx.Locals("user_id").(string)
	data.CreatedBy = userId

	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Get All Param Group godoc
// @Summary Get all param groups
// @Description Get list of param groups with pagination
// @Tags Param Group
// @Accept json
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/param-group [get]
func (c *ParamGroupController) FindAll(ctx fiber.Ctx) error {
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

// Update Param Group godoc
// @Summary Update param group
// @Description Update existing param group
// @Tags Param Group
// @Accept json
// @Produce json
// @Param request body models.ParamGroup true "Param Group data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/param-group [put]
func (c *ParamGroupController) Update(ctx fiber.Ctx) error {
	data := models.ParamGroup{}

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

// Delete Param Group godoc
// @Summary Delete param group
// @Description Delete param group by ID
// @Tags Param Group
// @Accept json
// @Produce json
// @Param id path string true "Param Group ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/param-group/{id} [delete]
func (c *ParamGroupController) Delete(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.repo.Delete(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, nil)
}
