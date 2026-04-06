package controllers

import (
	"hrms_go/dto/request/auth"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type AuthController struct {
	repo repositories.UserRepository
}

func NewAuthController(repo repositories.UserRepository) *AuthController {
	return &AuthController{repo}
}


// Login godoc
// @Summary Login
// @Description Login
// @Tags Login
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Login Request"
// @Success 200 {object} map[string]interface{}
// @Router /api/login [post]
func (c *AuthController) Login(ctx fiber.Ctx) error {
	var req auth.LoginRequest
	if err := ctx.Bind().Body(&req); err != nil {
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
