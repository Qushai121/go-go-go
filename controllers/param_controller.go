package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v2"
)

type ParamController struct {
	repo repositories.ParamRepository
}

func NewParamController(repo repositories.ParamRepository) *ParamController {
	return &ParamController{repo}
}

func (c *ParamController) Create(ctx *fiber.Ctx) error {
	var param models.Param

	if err := ctx.BodyParser(&param); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	userId := ctx.Locals("user_id").(string)
	param.CreatedBy = userId

	if err := c.repo.Create(&param); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, param)
}

func (c *ParamController) FindAll(ctx *fiber.Ctx) error {
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

func (c *ParamController) Update(ctx *fiber.Ctx) error {
	data := models.Param{}

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

func (c *ParamController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.repo.Delete(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, nil)
}