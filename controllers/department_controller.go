package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type DepartmentController struct {
	repo repositories.DepartmentRepository
}

func NewDepartmentController(repo repositories.DepartmentRepository) *DepartmentController {
	return &DepartmentController{repo}
}

// Create Department godoc
// @Summary Create department
// @Description Create new department
// @Tags Department
// @Accept json
// @Produce json
// @Param request body models.Department true "Department data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/department [post]
func (c *DepartmentController) Create(ctx fiber.Ctx) error {
	var data models.Department
	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	data.CreatedBy = ctx.Locals("user_id").(string)
	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

// Get All Department godoc
// @Summary Get all departments
// @Description Get list of departments with pagination
// @Tags Department
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
// @Router /api/department [get]
func (c *DepartmentController) FindAll(ctx fiber.Ctx) error {
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

// Update Department godoc
// @Summary Update department
// @Description Update existing department
// @Tags Department
// @Accept json
// @Produce json
// @Param request body models.Department true "Department data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/department [put]
func (c *DepartmentController) Update(ctx fiber.Ctx) error {
	var data models.Department
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

// Delete Department godoc
// @Summary Delete department
// @Description Delete department by ID
// @Tags Department
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/department/{id} [delete]
func (c *DepartmentController) Delete(ctx fiber.Ctx) error {
	if err := c.repo.Delete(ctx.Params("id")); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, nil)
}
