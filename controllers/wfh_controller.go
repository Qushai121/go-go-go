package controllers

import (
	"hrms_go/dto"
	"hrms_go/dto/wfh"
	mappers "hrms_go/mapper"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type WFHController struct {
	repo repositories.WFHRepository
}

func NewWFHController(repo repositories.WFHRepository) *WFHController {
	return &WFHController{repo}
}

// Create WFH godoc
// @Summary Create WFH
// @Description Create new WFH
// @Tags WFH
// @Accept json
// @Produce json
// @Param request body wfh.PostWfhDto true "WFH data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/wfh [post]
func (c *WFHController) Create(ctx fiber.Ctx) error {
	body := wfh.PostWfhDto{}

	if err := ctx.Bind().Body(&body); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	
	data,err := mappers.ToWFHModel(body)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	userId := ctx.Locals("user_id").(string)
	data.CreatedBy = userId

	
	id, err := uuid.Parse(userId)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	data.UserId = id

	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Get All WFH godoc
// @Summary Get all WFH
// @Description Get list of WFH with pagination
// @Tags WFH
// @Accept json
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/wfh [get]
func (c *WFHController) FindAll(ctx fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}

	if err := ctx.Bind().Query(&queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	result, err := c.repo.FindAll(nil,&queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, result)
}

// Get WFH By User godoc
// @Summary Get all WFH By User
// @Description Get list of WFH By User with pagination
// @Tags WFH
// @Accept json
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/wfh/me [get]
func (c *WFHController) FindByUser(ctx fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}

	if err := ctx.Bind().Query(&queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	userId := ctx.Locals("user_id").(string)
	result, err := c.repo.FindAll(&userId,&queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, result)
}

// Delete WFH godoc
// @Summary Delete WFH
// @Description Delete WFH by ID
// @Tags WFH
// @Accept json
// @Produce json
// @Param id path string true "WFH ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/wfh/{id} [delete]
func (c *WFHController) Delete(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.repo.Delete(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, nil)
}