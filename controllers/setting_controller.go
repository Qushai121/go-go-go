package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v2"
)

type SettingController struct {
	repo repositories.SettingRepository
}

func NewSettingController(repo repositories.SettingRepository) *SettingController {
	return &SettingController{repo}
}

func (c *SettingController) Create(ctx *fiber.Ctx) error {
	var setting models.Setting

	if err := ctx.BodyParser(&setting); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	userId := ctx.Locals("user_id").(string)
	setting.CreatedBy = userId

	if err := c.repo.Create(&setting); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, setting)
}

func (c *SettingController) FindAll(ctx *fiber.Ctx) error {
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

func (c *SettingController) Update(ctx *fiber.Ctx) error {
	data := models.Setting{}

	if err := ctx.BodyParser(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	userId := ctx.Locals("user_id").(string)
	data.UpdatedBy = &userId

	if err := c.repo.Update(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

func (c *SettingController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.repo.Delete(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, nil)
}