package controllers

import (
	"fmt"
	"strings"
	"time"

	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type OfficePermitController struct {
	repo repositories.OfficePermitRepository
}

func NewOfficePermitController(repo repositories.OfficePermitRepository) *OfficePermitController {
	return &OfficePermitController{repo}
}

type officePermitRequest struct {
	OfficePermitId   string  `json:"office_permit_id" form:"office_permit_id"`
	EmployeeNik      string  `json:"employee_nik" form:"employee_nik"`
	OfficePermitDate string  `json:"office_permit_date" form:"office_permit_date"`
	Remarks          *string `json:"remarks" form:"remarks"`
	Status           string  `json:"status" form:"status"`
	CurrentStep      int     `json:"current_step" form:"current_step"`
	ApprovalHeaderId string  `json:"approvalheader_id" form:"approvalheader_id"`
	ObjectCode       string  `json:"object_code" form:"object_code"`
	CreatedBy        string  `json:"created_by" form:"created_by"`
	UpdatedBy        string  `json:"updated_by" form:"updated_by"`
}

// Create Office Permit godoc
// @Summary Create office permit
// @Description Create new office permit transaction
// @Tags OfficePermit
// @Accept json
// @Produce json
// @Param request body models.OfficePermit true "Office permit data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/office-permit [post]
func (c *OfficePermitController) Create(ctx fiber.Ctx) error {
	var request officePermitRequest
	if err := ctx.Bind().Body(&request); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	data, err := request.toModel(true)
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	if employeeNik := strings.TrimSpace(fmt.Sprint(ctx.Locals("employee_nik"))); employeeNik != "" && employeeNik != "<nil>" {
		if data.EmployeeNik == "" {
			data.EmployeeNik = employeeNik
		}
		if data.CreatedBy == "" || strings.EqualFold(data.CreatedBy, "system") {
			data.CreatedBy = employeeNik
		}
	}

	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Get All Office Permit godoc
// @Summary Get all office permits
// @Description Get office permit list with pagination and search
// @Tags OfficePermit
// @Accept json
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param field_search query string false "Field name for dynamic search"
// @Param employee_nik query string false "Employee NIK"
// @Param office_permit_date query string false "Office permit date YYYY-MM-DD"
// @Param status query string false "Status"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/office-permit [get]
func (c *OfficePermitController) FindAll(ctx fiber.Ctx) error {
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

// Update Office Permit godoc
// @Summary Update office permit
// @Description Update existing office permit transaction
// @Tags OfficePermit
// @Accept json
// @Produce json
// @Param request body models.OfficePermit true "Office permit data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/office-permit [put]
func (c *OfficePermitController) Update(ctx fiber.Ctx) error {
	var request officePermitRequest
	if err := ctx.Bind().Body(&request); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	data, err := request.toModel(false)
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	if data.OfficePermitId == uuid.Nil {
		return utils.Error(ctx, 400, "office_permit_id wajib diisi")
	}

	if employeeNik := strings.TrimSpace(fmt.Sprint(ctx.Locals("employee_nik"))); employeeNik != "" && employeeNik != "<nil>" {
		data.UpdatedBy = &employeeNik
	} else if strings.TrimSpace(request.UpdatedBy) != "" {
		updatedBy := strings.TrimSpace(request.UpdatedBy)
		data.UpdatedBy = &updatedBy
	}

	if err := c.repo.Update(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Delete Office Permit godoc
// @Summary Delete office permit
// @Description Delete office permit by ID
// @Tags OfficePermit
// @Accept json
// @Produce json
// @Param id path string true "Office Permit ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/office-permit/{id} [delete]
func (c *OfficePermitController) Delete(ctx fiber.Ctx) error {
	if err := c.repo.Delete(ctx.Params("id")); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, nil)
}

func (r officePermitRequest) toModel(isCreate bool) (models.OfficePermit, error) {
	var data models.OfficePermit

	if strings.TrimSpace(r.OfficePermitId) != "" {
		id, err := uuid.Parse(strings.TrimSpace(r.OfficePermitId))
		if err != nil {
			return data, fmt.Errorf("office_permit_id tidak valid")
		}
		data.OfficePermitId = id
	} else if isCreate {
		data.OfficePermitId = uuid.New()
	}

	data.EmployeeNik = strings.TrimSpace(r.EmployeeNik)
	if isCreate && data.EmployeeNik == "" {
		return data, fmt.Errorf("employee_nik wajib diisi")
	}

	if strings.TrimSpace(r.OfficePermitDate) != "" {
		permitDate, err := parseOfficePermitDate(r.OfficePermitDate)
		if err != nil {
			return data, err
		}
		data.OfficePermitDate = permitDate
	} else if isCreate {
		return data, fmt.Errorf("office_permit_date wajib diisi")
	}

	data.Remarks = trimOptionalString(r.Remarks)
	data.Status = strings.ToUpper(strings.TrimSpace(r.Status))
	if data.Status == "" && isCreate {
		data.Status = "P"
	}

	data.CurrentStep = r.CurrentStep
	if data.CurrentStep <= 0 && isCreate {
		data.CurrentStep = 1
	}

	if strings.TrimSpace(r.ApprovalHeaderId) != "" {
		approvalHeaderId, err := uuid.Parse(strings.TrimSpace(r.ApprovalHeaderId))
		if err != nil {
			return data, fmt.Errorf("approvalheader_id tidak valid")
		}
		data.ApprovalHeaderId = &approvalHeaderId
	}

	data.ObjectCode = strings.TrimSpace(r.ObjectCode)
	if data.ObjectCode == "" && isCreate {
		data.ObjectCode = "LEAVE_HISTORY"
	}

	data.CreatedBy = strings.TrimSpace(r.CreatedBy)
	if data.CreatedBy == "" && isCreate {
		data.CreatedBy = "System"
	}

	return data, nil
}

func parseOfficePermitDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	layouts := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("format office_permit_date harus YYYY-MM-DD atau YYYY-MM-DD HH:mm:ss")
}

func trimOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
