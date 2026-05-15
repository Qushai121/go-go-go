package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type BranchController struct {
	repo repositories.BranchRepository
}

func NewBranchController(repo repositories.BranchRepository) *BranchController {
	return &BranchController{repo}
}

// Create Branch godoc
// @Summary Create branch
// @Description Create new branch
// @Tags Branch
// @Accept json
// @Produce json
// @Param request body models.Branch true "Branch data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/branch [post]
func (c *BranchController) Create(ctx fiber.Ctx) error {
	var data models.Branch
	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	data.CreatedBy = ctx.Locals("user_id").(string)

	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Get All Branch godoc
// @Summary Get all branches
// @Description Get list of branches with pagination
// @Tags Branch
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
// @Router /api/branch [get]
func (c *BranchController) FindAll(ctx fiber.Ctx) error {
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

// Update Branch godoc
// @Summary Update branch
// @Description Update existing branch
// @Tags Branch
// @Accept json
// @Produce json
// @Param request body models.Branch true "Branch data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/branch [put]
func (c *BranchController) Update(ctx fiber.Ctx) error {
	var data models.Branch
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

// Delete Branch godoc
// @Summary Delete branch
// @Description Delete branch by ID
// @Tags Branch
// @Accept json
// @Produce json
// @Param id path string true "Branch ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/branch/{id} [delete]
func (c *BranchController) Delete(ctx fiber.Ctx) error {
	if err := c.repo.Delete(ctx.Params("id")); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, nil)
}
