package controllers

import (
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type LeaveTypeController struct {
	repo repositories.LeaveTypeRepository
}

func NewLeaveTypeController(repo repositories.LeaveTypeRepository) *LeaveTypeController {
	return &LeaveTypeController{repo: repo}
}

// FindAll godoc
// @Summary Search leave type
// @Description Search tipe cuti aktif untuk transaksi cuti dynamic
// @Tags Leave Type
// @Accept json
// @Produce json
// @Param company_code query string false "Company Code"
// @Param search query string false "Search code/name"
// @Param deduct_leave_balance query string false "true/false"
// @Param is_active query string false "true/false"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/leave-type [get]
func (c *LeaveTypeController) FindAll(ctx fiber.Ctx) error {
	filter := models.LeaveTypeFilter{
		CompanyCode:        ctx.Query("company_code"),
		Search:             ctx.Query("search"),
		DeductLeaveBalance: ctx.Query("deduct_leave_balance"),
		IsActive:           ctx.Query("is_active"),
	}

	data, err := c.repo.FindAll(filter)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}
