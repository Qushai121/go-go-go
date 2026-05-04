package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type OfficeController struct {
	repo repositories.OfficeRepository
}

func NewOfficeController(repo repositories.OfficeRepository) *OfficeController {
	return &OfficeController{repo}
}

// Create Office godoc
// @Summary Create office
// @Description Create new office
// @Tags Office
// @Accept json
// @Produce json
// @Param request body models.Office true "Office data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/office [post]
func (c *OfficeController) Create(ctx fiber.Ctx) error {

	var data models.Office

	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	data.CreatedBy = ctx.Locals("user_id").(string)

	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Get All Office godoc
// @Summary Get all offices
// @Description Get list of offices with pagination
// @Tags Office
// @Accept json
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/office [get]
func (c *OfficeController) FindAll(ctx fiber.Ctx) error {

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

// Update Office godoc
// @Summary Update office
// @Description Update existing office
// @Tags Office
// @Accept json
// @Produce json
// @Param request body models.Office true "Office data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/office [put]
func (c *OfficeController) Update(ctx fiber.Ctx) error {

	var data models.Office

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

// Delete Office godoc
// @Summary Delete office
// @Description Delete office by ID
// @Tags Office
// @Accept json
// @Produce json
// @Param id path string true "Office ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/office/{id} [delete]
func (c *OfficeController) Delete(ctx fiber.Ctx) error {

	id := ctx.Params("id")

	if err := c.repo.Delete(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, nil)
}
