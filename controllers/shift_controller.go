package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type ShiftController struct {
	repo repositories.ShiftRepository
}

// Create Shift godoc
// @Summary Create shift
// @Description Create new shift
// @Tags Shift
// @Accept json
// @Produce json
// @Param request body models.Shift true "Shift data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/shift [post]
func (c *ShiftController) Create(ctx fiber.Ctx) error {
	var data models.Shift
	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	data.CreatedBy = ctx.Locals("user_id").(string)

	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Get All Shift godoc
// @Summary Get all shifts
// @Description Get list of shifts with pagination
// @Tags Shift
// @Accept json
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/shift [get]
func (c *ShiftController) FindAll(ctx fiber.Ctx) error {
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

// Update Shift godoc
// @Summary Update shift
// @Description Update existing shift
// @Tags Shift
// @Accept json
// @Produce json
// @Param request body models.Shift true "Shift data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/shift [put]
func (c *ShiftController) Update(ctx fiber.Ctx) error {
	data := models.Shift{}

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

// Delete Shift godoc
// @Summary Delete shift
// @Description Delete shift by ID
// @Tags Shift
// @Accept json
// @Produce json
// @Param id path string true "Shift ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/shift/{id} [delete]
func (c *ShiftController) Delete(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.repo.Delete(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, nil)
}

func NewShiftController(repo repositories.ShiftRepository) *ShiftController {
	return &ShiftController{repo}
}
