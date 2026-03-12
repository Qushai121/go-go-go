package controllers

import (
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	repo repositories.UserRepository
}

func NewAuthController(repo repositories.UserRepository) *AuthController {
	return &AuthController{repo}
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.Error(ctx, 400, "invalid request")
	}

	user, err := c.repo.FindByEmail(req.Email)
	if err != nil {
		return utils.Error(ctx, 401, "email not found")
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		return utils.Error(ctx, 401, "wrong password")
	}

	token, err := utils.GenerateToken(user.UserId.String(), user.EmployeeNIK)
	if err != nil {
		return utils.Error(ctx, 500, "failed generate token")
	}

	response := struct {
		*models.User
		Token string `json:"token"`
	}{
		User:  user,
		Token: token,
	}

	return utils.Success(ctx, response)
}