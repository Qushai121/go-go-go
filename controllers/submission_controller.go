package controllers

import (
	"hrms_go/dto"
	"hrms_go/dto/submission"
	mappers "hrms_go/mapper"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type SubmissionController struct {
	repo repositories.SubmissionRepository
}

func NewSubmissionController(repo repositories.SubmissionRepository) *SubmissionController {
	return &SubmissionController{repo}
}

// Create Submission godoc
// @Summary Create submission
// @Description Create new claim submission
// @Tags Submission
// @Accept json
// @Produce json
// @Param request body submission.PostSubmissionDto true "Submission data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/claim_submission [post]
func (c *SubmissionController) Create(ctx fiber.Ctx) error {
	body := submission.PostSubmissionDto{}
	data := models.Submission{}

	if err := ctx.Bind().Body(&body); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	data,err := mappers.ToSubmissionModel(body);
	if err != nil {
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

// Get All Submission godoc
// @Summary Get all submissions
// @Description Get list of submissions with pagination
// @Tags Submission
// @Accept json
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/claim_submission [get]
func (c *SubmissionController) FindAll(ctx fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}

	if err := ctx.Bind().Query(&queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	result, err := c.repo.FindAll(&queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, result)
}

// Get My Submission godoc
// @Summary Get current user submissions
// @Description Get submissions by logged-in user
// @Tags Submission
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/claim_submission/me [get]
func (c *SubmissionController) FindByUser(ctx fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}

	if err := ctx.Bind().Query(&queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	userId := ctx.Locals("user_id").(string)
	result, err := c.repo.FindByUser(userId, &queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, result)
}

// Update Submission godoc
// @Summary Update submission
// @Description Update existing submission
// @Tags Submission
// @Accept json
// @Produce json
// @Param request body models.Submission true "Submission data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/claim_submission [put]
func (c *SubmissionController) Update(ctx fiber.Ctx) error {
	body := submission.PostSubmissionDto{}
	data := models.Submission{}

	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	data,err := mappers.ToSubmissionModel(body);
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	userId := ctx.Locals("user_id").(string)
	_ = userId // optional: set UpdatedBy

	if err := c.repo.Update(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

// Delete Submission godoc
// @Summary Delete submission
// @Description Delete submission by ID
// @Tags Submission
// @Accept json
// @Produce json
// @Param id path string true "Submission ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/claim_submission/{id} [delete]
func (c *SubmissionController) Delete(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.repo.Delete(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, nil)
}
