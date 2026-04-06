package controllers

import (
	"fmt"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type SetupController struct {
	repo repositories.UserRepository
}

func NewSetupController(repo repositories.UserRepository) *SetupController {
	return &SetupController{repo}
}

// InitAdmin godoc
// @Summary Initialize admin
// @Description create first admin user
// @Tags Setup
// @Accept json
// @Produce json
// @Param request body models.User true "User data"
// @Success 200 {object} map[string]interface{}
// @Router /api/setup [post]
func (c *SetupController) InitAdmin(ctx fiber.Ctx) error {
	count, _ := c.repo.Count()
	if count > 0 {
		return utils.Error(ctx, 403, "setup already completed")
	}

	var user models.User
	if err := ctx.Bind().Body(&user); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	user.Password = strings.TrimSpace(user.Password)

	fmt.Println(user.Password)

	hashed, _ := utils.HashPassword(user.Password)
	user.Password = hashed
	user.Role = "SUPERADMIN"
	user.CreatedBy = "SYSTEM"

	if err := c.repo.Create(&user); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, fiber.Map{
		"message": "admin created",
	})
}
