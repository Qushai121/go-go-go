package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type ApprovalTemplateController struct {
	repo repositories.ApprovalTemplateRepository
}

func NewApprovalTemplateController(repo repositories.ApprovalTemplateRepository) *ApprovalTemplateController {
	return &ApprovalTemplateController{repo}
}

// Create Template Header godoc
// @Summary Create approval template header
// @Tags Approval Template
// @Accept json
// @Produce json
// @Param request body models.ApprovalTemplateHeader true "Header"
// @Success 200 {object} map[string]interface{}
// @Router /api/approval_template/header [post]
func (c *ApprovalTemplateController) CreateHeader(ctx fiber.Ctx) error {
	var data models.ApprovalTemplateHeader

	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	if err := c.repo.CreateHeader(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Get All Template Header godoc
// @Summary Get all template headers
// @Tags Approval Template
// @Accept json
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /api/approval_template/header [get]
func (c *ApprovalTemplateController) FindAllHeader(ctx fiber.Ctx) error {
	query := dto.PaginateFieldDto{}

	if err := ctx.Bind().Query(&query); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	data, err := c.repo.FindAllHeader(&query)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Detail Template Header godoc
// @Summary Get template header detail
// @Tags Approval Template
// @Param id path string true "Header ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/approval_template/header/{id} [get]
func (c *ApprovalTemplateController) DetailHeader(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	data, err := c.repo.DetailHeader(id)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Update Template Header godoc
// @Summary Update template header
// @Tags Approval Template
// @Param id path string true "Header ID"
// @Param request body models.ApprovalTemplateHeader true "Header"
// @Success 200 {object} map[string]interface{}
// @Router /api/approval_template/header/{id} [put]
func (c *ApprovalTemplateController) UpdateHeader(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var data models.ApprovalTemplateHeader

	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	if err := c.repo.UpdateHeader(id, &data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, "updated")
}

// Delete Template Header godoc
// @Summary Delete template header
// @Tags Approval Template
// @Param id path string true "Header ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/approval_template/header/{id} [delete]
func (c *ApprovalTemplateController) DeleteHeader(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.repo.DeleteHeader(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, "deleted")
}

// Create Template Detail godoc
// @Summary Create template detail
// @Tags Approval Template
// @Param request body models.ApprovalTemplateDetail true "Detail"
// @Success 200 {object} map[string]interface{}
// @Router /api/approval_template/detail [post]
func (c *ApprovalTemplateController) CreateDetail(ctx fiber.Ctx) error {
	var data models.ApprovalTemplateDetail

	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	if err := c.repo.CreateDetail(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Get Template Detail by Header godoc
// @Summary Get detail by header
// @Tags Approval Template
// @Param header_id path string true "Header ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/approval_template/detail/{header_id} [get]
func (c *ApprovalTemplateController) FindDetailByHeader(ctx fiber.Ctx) error {
	id := ctx.Params("header_id")

	data, err := c.repo.FindByHeader(id)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Update Template Detail godoc
// @Summary Update template detail
// @Tags Approval Template
// @Param id path string true "Detail ID"
// @Param request body models.ApprovalTemplateDetail true "Detail"
// @Success 200 {object} map[string]interface{}
// @Router /api/approval_template/detail/{id} [put]
func (c *ApprovalTemplateController) UpdateDetail(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var data models.ApprovalTemplateDetail

	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	if err := c.repo.UpdateDetail(id, &data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, "updated")
}

// Delete Template Detail godoc
// @Summary Delete template detail
// @Tags Approval Template
// @Param id path string true "Detail ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/approval_template/detail/{id} [delete]
func (c *ApprovalTemplateController) DeleteDetail(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.repo.DeleteDetail(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, "deleted")
}