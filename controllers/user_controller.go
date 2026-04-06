package controllers

import (
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UserController struct {
	repo repositories.UserRepository
}

func NewUserController(repo repositories.UserRepository) *UserController {
	return &UserController{repo}
}

func (c *UserController) Create(ctx fiber.Ctx) error {
	var user models.User
	if err := ctx.Bind().Body(&user); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	hashed, _ := utils.HashPassword(user.Password)
	user.Password = hashed
	user.CreatedBy = ctx.Locals("employee_nik").(string)

	if err := c.repo.Create(&user); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, user)
}

func (c *UserController) FindAll(ctx fiber.Ctx) error {
	users, err := c.repo.FindAll()
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, users)
}

func (c *UserController) UpdateUserShift(ctx fiber.Ctx) error {
	var data models.User
	realData := models.User{}
	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	data.UpdatedBy = ctx.Locals("user_id").(string)

	realData.UpdatedBy = ctx.Locals("user_id").(string)
	realData.ShiftId = data.ShiftId
	realData.UserId = data.UserId

	if err := c.repo.Update(&realData); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

func (c *UserController) UpdateUserPicture(ctx fiber.Ctx) error {

	realData := models.User{}

	userIdStr := ctx.FormValue("user_id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return utils.Error(ctx, 400, "invalid user_id")
	}

	file, err := ctx.FormFile("profile_picture_url")
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	fileUrl, err := utils.SaveFileToPath(file, "user", ctx)
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	realData.UpdatedBy = ctx.Locals("user_id").(string)
	realData.UserId = userId

	if fileUrl != nil {
		realData.ProfilePictureUrl = *fileUrl
	}

	if err := c.repo.UpdateProfilePicture(&realData); err != nil {
		if fileUrl != nil {
			utils.RemoveFileFromPath(*fileUrl)
		}
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, realData)
}
