package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v2"
)

type CompaniesController struct {
	repo repositories.CompaniesRepository
}

func (c *CompaniesController) Create(ctx *fiber.Ctx) error {
	data := models.Companies{}
	
	if err := ctx.BodyParser(&data); err != nil {
		return utils.Error(ctx, 400, "invalid request")
	}
	data.CreatedBy = ctx.Locals("user_id").(string)

	if err := c.repo.Create(&data);err != nil {
		return utils.Error(ctx, 500, "failed to create company")
	}

	return utils.Success(ctx, data);
}


func (c *CompaniesController) FindAll(ctx *fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}

	if err := ctx.QueryParser(&queryParams); err != nil {
		return utils.Error(ctx, 400, "invalid request")
	}

	shift, err := c.repo.FindAll(&queryParams);
	if err != nil {
		return utils.Error(ctx, 500, "failed to get list of company")
	}
	return utils.Success(ctx, shift)
}


func (c *CompaniesController) Update(ctx *fiber.Ctx) error {
	data := models.Companies{}
	
	if err := ctx.BodyParser(&data); err != nil {
		return utils.Error(ctx, 400, "invalid request")
	}
	userId := ctx.Locals("user_id").(string)
	data.UpdatedBy = &userId

	if err := c.repo.Update(&data);err != nil {
		return utils.Error(ctx, 500, "failed to update company")
	}

	return utils.Success(ctx, data);
}


func (c *CompaniesController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.repo.Delete(id); err != nil {
		return utils.Error(ctx, 500, "failed to delete company")
	}
	return utils.Success(ctx,nil)
}

func NewCompaniesController(repo repositories.CompaniesRepository) *CompaniesController  {
	return &CompaniesController{repo}
}