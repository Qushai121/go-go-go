package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LeaveController struct {
	repo repositories.LeaveRepository
}

func NewLeaveController(repo repositories.LeaveRepository) *LeaveController {
	return &LeaveController{repo}
}

func (c *LeaveController) Create(ctx *fiber.Ctx) error {
	var leave models.Leave
	if err := ctx.BodyParser(&leave); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	if leave.LeaveID == uuid.Nil {
		leave.LeaveID = uuid.New()
	}

	if err := c.repo.Create(&leave); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, leave)
}

func (c *LeaveController) FindAll(ctx *fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}

	if err := ctx.QueryParser(&queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	leaves, err := c.repo.FindAll(&queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, leaves)
}

func (c *LeaveController) Update(ctx *fiber.Ctx) error {
	var leave models.Leave

	if err := ctx.BodyParser(&leave); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	userId := ctx.Locals("user_id").(string)
	leave.UpdatedBy = &userId // optional if you track who updated

	if err := c.repo.Update(&leave); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, leave)
}

func (c *LeaveController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.repo.Delete(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, nil)
}