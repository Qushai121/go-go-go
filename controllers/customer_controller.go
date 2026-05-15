package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type CustomerController struct {
	repo repositories.CustomerRepository
}

func NewCustomerController(repo repositories.CustomerRepository) *CustomerController {
	return &CustomerController{repo}
}

// Create Customer godoc
// @Summary Create customer
// @Description Create new customer
// @Tags Customer
// @Accept json
// @Produce json
// @Param request body models.Customer true "Customer data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/customer [post]
func (c *CustomerController) Create(ctx fiber.Ctx) error {
	var data models.Customer
	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	data.CreatedBy = ctx.Locals("user_id").(string)

	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Get All Customers godoc
// @Summary Get all customers
// @Description Get list of customers with pagination
// @Tags Customer
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
// @Router /api/customer [get]
func (c *CustomerController) FindAll(ctx fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}

	if err := utils.BindPaginationParams(ctx, &queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	customers, err := c.repo.FindAll(&queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, customers)
}

// Update Customer godoc
// @Summary Update customer
// @Description Update existing customer
// @Tags Customer
// @Accept json
// @Produce json
// @Param request body models.Customer true "Customer data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/customer [put]
func (c *CustomerController) Update(ctx fiber.Ctx) error {
	data := models.Customer{}

	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	userId := ctx.Locals("user_id").(string)
	data.UpdatedBy = &userId

	if err := c.repo.Update(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Delete Customer godoc
// @Summary Delete customer
// @Description Delete customer by ID
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/customer/{id} [delete]
func (c *CustomerController) Delete(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.repo.Delete(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, nil)
}
