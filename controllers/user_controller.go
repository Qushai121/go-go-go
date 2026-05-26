package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"
	"mime/multipart"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UserController struct {
	repo repositories.UserRepository
}

func NewUserController(repo repositories.UserRepository) *UserController {
	return &UserController{repo}
}

// Create User godoc
// @Summary Create user
// @Description Create new user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body models.User true "User data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/users [post]
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

// Get All Users godoc
// @Summary Get all users
// @Description Get list of users
// @Tags Users
// @Accept json
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param field_search query string false "Field name for dynamic search"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/users [get]
func (c *UserController) FindAll(ctx fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}

	if err := utils.BindPaginationParams(ctx, &queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	users, err := c.repo.FindAll(&queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, users)
}

// Update User Picture godoc
// @Summary Update user profile picture
// @Description Upload and update user profile picture
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Param user_id formData string true "User ID (UUID)"
// @Param profile_picture_url formData file true "Profile picture file"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/users/update-profile-picture [post]
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

	if _, err := utils.VerifyFaceImage(file); err != nil {
		return utils.Error(ctx, 422, err.Error())
	}

	targetUser, err := c.repo.FindByID(userId.String())
	if err != nil {
		return utils.Error(ctx, 404, "user not found")
	}

	fileUrl, err := utils.SaveFileToCustomPath(file, "profile-picture/"+targetUser.EmployeeNIK, ctx)
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	authEmployeeNIK := ctx.Locals("employee_nik").(string)
	realData.UpdatedBy = &authEmployeeNIK
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

// Verify User Face godoc
// @Summary Verify face image
// @Description Verify face image without saving user profile picture or attendance
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Face image file"
// @Param nik formData string false "Employee NIK. If filled, image will be matched against registered face data."
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/users/verify-face [post]
func (c *UserController) VerifyFace(ctx fiber.Ctx) error {
	file, err := firstFormFile(ctx, "image", "profile_picture_url", "photo_url", "imagepath")
	if err != nil {
		return utils.Error(ctx, 400, "image file is required")
	}

	nik := ctx.FormValue("nik")
	var result *utils.FaceServiceResult
	if nik != "" {
		result, err = utils.VerifyFaceByNIK(file, nik)
	} else {
		result, err = utils.VerifyFaceImage(file)
	}

	if err != nil {
		status := 422
		if result != nil && result.Status > 0 {
			status = result.Status
		}
		return utils.Error(ctx, status, err.Error())
	}

	return utils.Success(ctx, result)
}

// Find Me godoc
// @Summary Find Me
// @Description Get My User Information
// @Tags Users
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/users/me [get]
func (c *UserController) Me(ctx fiber.Ctx) error {
	userId := ctx.Locals("user_id").(string)

	data, err := c.repo.Me(userId)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

func firstFormFile(ctx fiber.Ctx, keys ...string) (*multipart.FileHeader, error) {
	var lastErr error
	for _, key := range keys {
		file, err := ctx.FormFile(key)
		if err == nil {
			return file, nil
		}
		lastErr = err
	}
	return nil, lastErr
}
