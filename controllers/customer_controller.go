package controllers

import (
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v2"
)

type CustomerController struct {
	repo repositories.CustomerRepository
}

func NewCustomerController(repo repositories.CustomerRepository) *CustomerController {
	return &CustomerController{repo}
}

func (c *CustomerController) Create(ctx *fiber.Ctx) error {
	var customer models.Customer
	if err := ctx.BodyParser(&customer); err != nil {
		return utils.Error(ctx, 400, "invalid request")
	}
	if err := c.repo.Create(&customer);err != nil {
		return utils.Error(ctx, 500, "failed to create customer")
	}

	return utils.Success(ctx, customer);
}

func (c *CustomerController) FindAll(ctx *fiber.Ctx) error {
	customers, err := c.repo.FindAll()
	if err != nil {
		return utils.Error(ctx, 500, "failed to fetch customers")
	}
	return utils.Success(ctx, customers)
}