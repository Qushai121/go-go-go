package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type LeaveController struct {
	repo repositories.LeaveRepository
}

func NewLeaveController(repo repositories.LeaveRepository) *LeaveController {
	return &LeaveController{
		repo: repo,
	}
}

// Add Cuti godoc
// @Summary Add cuti
// @Description Add cuti / leave history
// @Tags Leave
// @Accept json
// @Produce json
// @Param request body []models.LeaveHistory true "Leave data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/leave/cuti [post]
func (c *LeaveController) AddCuti(ctx fiber.Ctx) error {
	var request models.LeaveHistory

	body := strings.TrimSpace(string(ctx.Body()))
	if body == "" {
		return utils.Error(ctx, 400, "Request body wajib diisi")
	}

	if strings.HasPrefix(body, "[") {
		return utils.Error(ctx, 400, "Request body harus object, bukan array")
	}

	if err := json.Unmarshal([]byte(body), &request); err != nil {
		return utils.Error(ctx, 400, "Format JSON tidak valid")
	}

	employeeNik := strings.TrimSpace(fmt.Sprint(ctx.Locals("employee_nik")))
	if employeeNik != "" && employeeNik != "<nil>" {
		if strings.TrimSpace(request.CreatedBy) == "" {
			request.CreatedBy = employeeNik
		}

		if strings.TrimSpace(request.EmployeeNik) == "" {
			request.EmployeeNik = employeeNik
		}
	}

	if err := c.repo.AddCuti(request); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, fiber.Map{
		"message": "Cuti berhasil ditambahkan",
	})
}

// Get Cuti godoc
// @Summary Get cuti
// @Description Get cuti by filter
// @Tags Leave
// @Accept json
// @Produce json
// @Param leave_history_id query string false "Leave History ID"
// @Param employee_nik query string false "Employee NIK"
// @Param leave_type query string false "Leave Type"
// @Param start_date query string false "Start date YYYY-MM-DD untuk range overlap cuti"
// @Param end_date query string false "End date YYYY-MM-DD untuk range overlap cuti"
// @Param leave_start query string false "Leave start exact date YYYY-MM-DD"
// @Param leave_end query string false "Leave end exact date YYYY-MM-DD"
// @Param total_days query int false "Total Days"
// @Param remarks query string false "Remarks"
// @Param location query string false "Location"
// @Param leave_year query int false "Leave Year"
// @Param status query string false "Status"
// @Param current_step query int false "Current Step"
// @Param approvalheader_id query string false "Approval Header ID"
// @Param object_code query string false "Object Code"
// @Param created_at query string false "Created at exact date YYYY-MM-DD"
// @Param updated_at query string false "Updated at exact date YYYY-MM-DD"
// @Param created_by query string false "Created By"
// @Param updated_by query string false "Updated By"
// @Param created_at_start query string false "Created at start YYYY-MM-DD"
// @Param created_at_end query string false "Created at end YYYY-MM-DD"
// @Param updated_at_start query string false "Updated at start YYYY-MM-DD"
// @Param updated_at_end query string false "Updated at end YYYY-MM-DD"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/leave/cuti [get]
func (c *LeaveController) FindCuti(ctx fiber.Ctx) error {
	q := func(key string) string {
		return strings.TrimSpace(ctx.Query(key))
	}

	filter := models.LeaveHistoryFilter{
		LeaveHistoryId: q("leave_history_id"),

		EmployeeNik: q("employee_nik"),
		LeaveType:   q("leave_type"),

		StartDate: q("start_date"),
		EndDate:   q("end_date"),

		LeaveStart: q("leave_start"),
		LeaveEnd:   q("leave_end"),

		TotalDays: q("total_days"),
		Remarks:   q("remarks"),
		Location:  q("location"),

		LeaveYear: q("leave_year"),
		Status:    q("status"),

		CurrentStep:      q("current_step"),
		ApprovalHeaderId: q("approvalheader_id"),

		ObjectCode: q("object_code"),

		CreatedAt: q("created_at"),
		UpdatedAt: q("updated_at"),

		CreatedBy: q("created_by"),
		UpdatedBy: q("updated_by"),

		CreatedAtStart: q("created_at_start"),
		CreatedAtEnd:   q("created_at_end"),
		UpdatedAtStart: q("updated_at_start"),
		UpdatedAtEnd:   q("updated_at_end"),
	}

	data, err := c.repo.FindCuti(filter)
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	return utils.Success(ctx, data)
}
