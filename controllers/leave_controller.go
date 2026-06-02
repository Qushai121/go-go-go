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
	var request []models.LeaveHistory

	body := strings.TrimSpace(string(ctx.Body()))
	if body == "" {
		return utils.Error(ctx, 400, "Request body wajib diisi")
	}

	if strings.HasPrefix(body, "[") {
		if err := json.Unmarshal([]byte(body), &request); err != nil {
			return utils.Error(ctx, 400, "Format JSON tidak valid")
		}
	} else {
		var single models.LeaveHistory
		if err := json.Unmarshal([]byte(body), &single); err != nil {
			return utils.Error(ctx, 400, "Format JSON tidak valid")
		}

		request = append(request, single)
	}

	employeeNik := strings.TrimSpace(fmt.Sprint(ctx.Locals("employee_nik")))
	if employeeNik != "" && employeeNik != "<nil>" {
		for i := range request {
			if strings.TrimSpace(request[i].CreatedBy) == "" {
				request[i].CreatedBy = employeeNik
			}

			if strings.TrimSpace(request[i].EmployeeNik) == "" {
				request[i].EmployeeNik = employeeNik
			}
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
// @Description Get cuti by employee_nik and date range
// @Tags Leave
// @Accept json
// @Produce json
// @Param employee_nik query string false "Employee NIK"
// @Param start_date query string false "Start date YYYY-MM-DD"
// @Param end_date query string false "End date YYYY-MM-DD"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/leave/cuti [get]
func (c *LeaveController) FindCuti(ctx fiber.Ctx) error {
	employeeNik := strings.TrimSpace(ctx.Query("employee_nik"))
	startDate := strings.TrimSpace(ctx.Query("start_date"))
	endDate := strings.TrimSpace(ctx.Query("end_date"))

	data, err := c.repo.FindCuti(employeeNik, startDate, endDate)
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	return utils.Success(ctx, data)
}
