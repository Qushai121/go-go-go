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
	var data models.Attendance

	if err := ctx.BodyParser(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	file, err := ctx.FormFile("photo_url")
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	fileUrl,err := utils.SaveFileToPath(file,"attandance",ctx)
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	if(fileUrl != nil){
		data.PhotoUrl = *fileUrl
	}

	allowed := map[string]bool{
		"WFH": true,
		"VISIT": true,
		"OFFICE": true,
	}

	if !allowed[data.Activity] {
		return utils.Error(ctx, 400, "activity must be WFH, VISIT, or OFFICE")
	}

	allowedCheckType := map[string]bool{
		"1":true,
		"2":true,
	}

	if !allowedCheckType[data.CheckType]{
		return utils.Error(ctx, 400, "checktype value must be 1 or 2")
	}

	if err := c.repo.Create(&data); err != nil {
		if fileUrl != nil {
			utils.RemoveFileFromPath(*fileUrl);
		}
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
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
	queryParams := dto.PaginateFieldDto{}
	if err := ctx.QueryParser(&queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	userId := ctx.Locals("user_id").(string)
	data, err := c.repo.FindByUser(userId, &queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}