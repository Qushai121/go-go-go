package controllers

import (
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	repo repositories.UserRepository
}

func NewUserController(repo repositories.UserRepository) *UserController {
	return &UserController{repo}
}

func (c *UserController) Create(ctx *fiber.Ctx) error {
	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return utils.Error(ctx, 400, "invalid request")
	}

	hashed, _ := utils.HashPassword(user.Password)
	user.Password = hashed
	user.CreatedBy = ctx.Locals("employee_nik").(string)

	if err := c.repo.Create(&user); err != nil {
		return utils.Error(ctx, 500, "failed create user")
	}

	return utils.Success(ctx, user)
}

func (c *UserController) FindAll(ctx *fiber.Ctx) error {
	users, err := c.repo.FindAll()
	if err != nil {
		return utils.Error(ctx, 500, "failed get users")
	}	
	return utils.Success(ctx, users)
}