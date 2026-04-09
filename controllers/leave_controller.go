package controllers

import (
	"hrms_go/dto"
	"hrms_go/dto/leave"
	mappers "hrms_go/mapper"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type LeaveController struct {
	repo repositories.LeaveRepository
}

func NewLeaveController(repo repositories.LeaveRepository) *LeaveController {
	return &LeaveController{repo}
}

// Create Leave godoc
// @Summary Create leave
// @Description Create new leave request
// @Tags Leave
// @Accept json
// @Produce json
// @Param request body models.Leave true "Leave data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/leave [post]
func (c *LeaveController) Create(ctx fiber.Ctx) error {
	var body leave.PostLeaveDto
	if err := ctx.Bind().Body(&body); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	data,err := mappers.ToLeaveModel(body)
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	// data.CreatedBy = ctx.Locals("user_id").(string)

	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Get All Leave godoc
// @Summary Get all leave
// @Description Get list of leave with pagination
// @Tags Leave
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/leave [get]
func (c *LeaveController) FindAll(ctx fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}

	if err := ctx.Bind().Query(&queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	leaves, err := c.repo.FindAll(&queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, leaves)
}

// Update Leave godoc
// @Summary Update leave
// @Description Update existing leave
// @Tags Leave
// @Accept json
// @Produce json
// @Param request body models.Leave true "Leave data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/leave [put]
func (c *LeaveController) Update(ctx fiber.Ctx) error {
	var body leave.PostLeaveDto

	if err := ctx.Bind().Body(&body); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	
	data,err := mappers.ToLeaveModel(body)
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	
	// userId := ctx.Locals("user_id").(string)
	// data.UpdatedBy = &userId

	if err := c.repo.Update(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Delete Leave godoc
// @Summary Delete leave
// @Description Delete leave by ID
// @Tags Leave
// @Accept json
// @Produce json
// @Param id path string true "Leave ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/leave/{id} [delete]
func (c *LeaveController) Delete(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.repo.Delete(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, nil)
}
