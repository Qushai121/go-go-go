package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type ParamController struct {
	repo repositories.ParamRepository
}

func NewParamController(repo repositories.ParamRepository) *ParamController {
	return &ParamController{repo}
}

// Create Param godoc
// @Summary Create param
// @Description Create new parameter
// @Tags Param
// @Accept json
// @Produce json
// @Param request body models.Param true "Param data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/param [post]
func (c *ParamController) Create(ctx fiber.Ctx) error {
	var param models.Param

	if err := ctx.Bind().Body(&param); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	userId := ctx.Locals("user_id").(string)
	param.CreatedBy = userId

	if err := c.repo.Create(&param); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, param)
}

// Get All Param godoc
// @Summary Get all params
// @Description Get list of params with pagination
// @Tags Param
// @Accept json
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/param [get]
func (c *ParamController) FindAll(ctx fiber.Ctx) error {
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

// Update Param godoc
// @Summary Update param
// @Description Update existing param
// @Tags Param
// @Accept json
// @Produce json
// @Param request body models.Param true "Param data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/param [put]
func (c *ParamController) Update(ctx fiber.Ctx) error {
	data := models.Param{}

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

// Delete Param godoc
// @Summary Delete param
// @Description Delete param by ID
// @Tags Param
// @Accept json
// @Produce json
// @Param id path string true "Param ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/param/{id} [delete]
func (c *ParamController) Delete(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.repo.Delete(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, nil)
}
