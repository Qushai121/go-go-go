package controllers

import (
	"hrms_go/dto"
	"hrms_go/dto/approval"
	mappers "hrms_go/mapper"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type ApprovalController struct {
	repo repositories.ApprovalRepository
}

func NewApprovalController(repo repositories.ApprovalRepository) *ApprovalController {
	return &ApprovalController{repo}
}

// Get All Approval godoc
// @Summary Get all approvals
// @Description Get all approval header with pagination & search
// @Tags Approval
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
// @Router /api/approval [get]
func (c *ApprovalController) FindAll(ctx fiber.Ctx) error {
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

// Detail Approval godoc
// @Summary Get approval detail
// @Description Get approval header detail by ID
// @Tags Approval
// @Accept json
// @Produce json
// @Param id path string true "Approval Header ID (UUID)"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/approval/{id} [get]
func (c *ApprovalController) Detail(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	data, err := c.repo.Detail(id)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Delete Approval godoc
// @Summary Delete approval
// @Description Delete approval header by ID
// @Tags Approval
// @Accept json
// @Produce json
// @Param id path string true "Approval Header ID (UUID)"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/approval/{id} [delete]
func (c *ApprovalController) Delete(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.repo.Delete(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, "deleted successfully")
}

// Approve godoc
// @Summary Approve document
// @Description Insert or update approval detail (upsert)
// @Tags Approval
// @Accept json
// @Produce json
// @Param request body approval.PostApproveDto true "Approval Detail Payload"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/approval/approve [post]
func (c *ApprovalController) Approve(ctx fiber.Ctx) error {
	var body approval.PostApproveDto

	if err := ctx.Bind().Body(&body); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	isValidStatus := CheckApprovalStatus(body.ApprovalStatus)
	if isValidStatus {
		return utils.Error(ctx, 422, "invalid approval status")
	}

	data, err := mappers.ToApprovalDetail(body)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	if err := c.repo.Approve(data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, "approval updated")
}

func CheckApprovalStatus(status string) bool {
	switch status {
	case "PENDING", "APPROVED", "REJECTED":
		return true
	default:
		return false
	}
}
