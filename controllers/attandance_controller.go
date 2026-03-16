package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v2"
)

type AttendanceController struct {
	repo repositories.AttendanceRepository
}

func NewAttendanceController(repo repositories.AttendanceRepository) *AttendanceController {
	return &AttendanceController{repo}
}

func (c *AttendanceController) Create(ctx *fiber.Ctx) error {
	var attendance models.Attendance

	if err := ctx.BodyParser(&attendance); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	allowed := map[string]bool{
		"WFH": true,
		"VISIT": true,
		"OFFICE": true,
	}

	if !allowed[attendance.Activity] {
		return utils.Error(ctx, 400, "activity must be WFH, VISIT, or OFFICE")
	}

	if err := c.repo.Create(&attendance); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, attendance)
}

func (c *AttendanceController) FindAll(ctx *fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}

	if err := ctx.QueryParser(&queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	data, err := c.repo.FindAll(&queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

func (c *AttendanceController) FindByUser(ctx *fiber.Ctx) error {
	userId := ctx.Params("user_id")

	queryParams := dto.PaginateFieldDto{}
	if err := ctx.QueryParser(&queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	data, err := c.repo.FindByUser(userId, &queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}