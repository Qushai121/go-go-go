package controllers

import (
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v2"
)

type ShiftController struct {
	repo repositories.ShiftRepository
}


func (c *ShiftController) Create(ctx *fiber.Ctx) error {
	var shift models.Shift
	if err := ctx.BodyParser(&shift); err != nil {
		return utils.Error(ctx, 400, "invalid request")
	}
	if err := c.repo.Create(&shift);err != nil {
		return utils.Error(ctx, 500, "failed to create shift")
	}

	return utils.Success(ctx, shift);
}

func (c *ShiftController) FindAll(ctx *fiber.Ctx) error {
	shift, err := c.repo.FindAll()
	if err != nil {
		return utils.Error(ctx, 500, "failed to fetch shift")
	}
	return utils.Success(ctx, shift)
}

func NewShiftController(repo repositories.ShiftRepository) *ShiftController {
	return &ShiftController{repo}
}