package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v2"
)

type SubmissionController struct {
	repo repositories.SubmissionRepository
}

func NewSubmissionController(repo repositories.SubmissionRepository) *SubmissionController {
	return &SubmissionController{repo}
}

func (c *SubmissionController) Create(ctx *fiber.Ctx) error {
	data := models.Submission{}

	if err := ctx.BodyParser(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	userId := ctx.Locals("user_id").(string)
	// you can set CreatedBy if needed
	_ = userId

	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

func (c *SubmissionController) FindAll(ctx *fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}

	if err := ctx.QueryParser(&queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	result, err := c.repo.FindAll(&queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, result)
}

func (c *SubmissionController) FindByUser(ctx *fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}

	if err := ctx.QueryParser(&queryParams);err != nil{
		return utils.Error(ctx, 400, err.Error())
	}

	userId := ctx.Locals("user_id").(string)
	result,err := c.repo.FindByUser(userId,&queryParams)
	if err != nil{
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, result)
}

func (c *SubmissionController) Update(ctx *fiber.Ctx) error {
	data := models.Submission{}
	if err := ctx.BodyParser(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	userId := ctx.Locals("user_id").(string)
	_ = userId // optional: set UpdatedBy

	if err := c.repo.Update(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

func (c *SubmissionController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.repo.Delete(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, nil)
}