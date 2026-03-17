package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v2"
)

type OfficeController struct {
	repo repositories.OfficeRepository
}

func NewOfficeController(repo repositories.OfficeRepository) *OfficeController {
	return &OfficeController{repo}
}

func (c *OfficeController) Create(ctx *fiber.Ctx) error {

	var data models.Office

	if err := ctx.BodyParser(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	data.CreatedBy = ctx.Locals("user_id").(string)

	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

func (c *OfficeController) FindAll(ctx *fiber.Ctx) error {

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

func (c *OfficeController) Update(ctx *fiber.Ctx) error {

	var data models.Office

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

func (c *OfficeController) Delete(ctx *fiber.Ctx) error {

	id := ctx.Params("id")

	if err := c.repo.Delete(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, nil)
}