package controllers

import "hrms_go/models"

// createBranchDoc godoc
// @Summary Create branch
// @Description Create new branch
// @Tags Branch
// @Accept json
// @Produce json
// @Param request body models.Branch true "Branch data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/branch [post]
func createBranchDoc(models.Branch) {}

// findBranchDoc godoc
// @Summary Get all branches
// @Description Get list of branches with pagination
// @Tags Branch
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/branch [get]
func findBranchDoc() {}

// updateBranchDoc godoc
// @Summary Update branch
// @Description Update existing branch
// @Tags Branch
// @Accept json
// @Produce json
// @Param request body models.Branch true "Branch data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/branch [put]
func updateBranchDoc(models.Branch) {}

// deleteBranchDoc godoc
// @Summary Delete branch
// @Description Delete branch by ID
// @Tags Branch
// @Produce json
// @Param id path string true "Branch ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/branch/{id} [delete]
func deleteBranchDoc() {}

// createDivisionDoc godoc
// @Summary Create division
// @Description Create new division
// @Tags Division
// @Accept json
// @Produce json
// @Param request body models.Division true "Division data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/division [post]
func createDivisionDoc(models.Division) {}

// findDivisionDoc godoc
// @Summary Get all divisions
// @Description Get list of divisions with pagination
// @Tags Division
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/division [get]
func findDivisionDoc() {}

// updateDivisionDoc godoc
// @Summary Update division
// @Description Update existing division
// @Tags Division
// @Accept json
// @Produce json
// @Param request body models.Division true "Division data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/division [put]
func updateDivisionDoc(models.Division) {}

// deleteDivisionDoc godoc
// @Summary Delete division
// @Description Delete division by ID
// @Tags Division
// @Produce json
// @Param id path string true "Division ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/division/{id} [delete]
func deleteDivisionDoc() {}

// createDepartmentDoc godoc
// @Summary Create department
// @Description Create new department
// @Tags Department
// @Accept json
// @Produce json
// @Param request body models.Department true "Department data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/department [post]
func createDepartmentDoc(models.Department) {}

// findDepartmentDoc godoc
// @Summary Get all departments
// @Description Get list of departments with pagination
// @Tags Department
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/department [get]
func findDepartmentDoc() {}

// updateDepartmentDoc godoc
// @Summary Update department
// @Description Update existing department
// @Tags Department
// @Accept json
// @Produce json
// @Param request body models.Department true "Department data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/department [put]
func updateDepartmentDoc(models.Department) {}

// deleteDepartmentDoc godoc
// @Summary Delete department
// @Description Delete department by ID
// @Tags Department
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/department/{id} [delete]
func deleteDepartmentDoc() {}

// createTitleDoc godoc
// @Summary Create title
// @Description Create new title
// @Tags Title
// @Accept json
// @Produce json
// @Param request body models.Title true "Title data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/title [post]
func createTitleDoc(models.Title) {}

// findTitleDoc godoc
// @Summary Get all titles
// @Description Get list of titles with pagination
// @Tags Title
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/title [get]
func findTitleDoc() {}

// updateTitleDoc godoc
// @Summary Update title
// @Description Update existing title
// @Tags Title
// @Accept json
// @Produce json
// @Param request body models.Title true "Title data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/title [put]
func updateTitleDoc(models.Title) {}

// deleteTitleDoc godoc
// @Summary Delete title
// @Description Delete title by ID
// @Tags Title
// @Produce json
// @Param id path string true "Title ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/title/{id} [delete]
func deleteTitleDoc() {}

// createUserGroupDoc godoc
// @Summary Create user group
// @Description Create new user group
// @Tags UserGroup
// @Accept json
// @Produce json
// @Param request body models.UserGroup true "User group data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/usergroup [post]
func createUserGroupDoc(models.UserGroup) {}

// findUserGroupDoc godoc
// @Summary Get all user groups
// @Description Get list of user groups with pagination
// @Tags UserGroup
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/usergroup [get]
func findUserGroupDoc() {}

// updateUserGroupDoc godoc
// @Summary Update user group
// @Description Update existing user group
// @Tags UserGroup
// @Accept json
// @Produce json
// @Param request body models.UserGroup true "User group data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/usergroup [put]
func updateUserGroupDoc(models.UserGroup) {}

// deleteUserGroupDoc godoc
// @Summary Delete user group
// @Description Delete user group by ID
// @Tags UserGroup
// @Produce json
// @Param id path string true "User group ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/usergroup/{id} [delete]
func deleteUserGroupDoc() {}

// createUserCompanyDoc godoc
// @Summary Create user company
// @Description Create new user company mapping
// @Tags UserCompany
// @Accept json
// @Produce json
// @Param request body models.UserCompany true "User company data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-company [post]
func createUserCompanyDoc(models.UserCompany) {}

// findUserCompanyDoc godoc
// @Summary Get all user companies
// @Description Get list of user company mappings with pagination
// @Tags UserCompany
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-company [get]
func findUserCompanyDoc() {}

// updateUserCompanyDoc godoc
// @Summary Update user company
// @Description Update existing user company mapping
// @Tags UserCompany
// @Accept json
// @Produce json
// @Param request body models.UserCompany true "User company data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-company [put]
func updateUserCompanyDoc(models.UserCompany) {}

// deleteUserCompanyDoc godoc
// @Summary Delete user company
// @Description Delete user company mapping by ID
// @Tags UserCompany
// @Produce json
// @Param id path string true "User company ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-company/{id} [delete]
func deleteUserCompanyDoc() {}

// createUserOfficeDoc godoc
// @Summary Create user office
// @Description Create new user office mapping
// @Tags UserOffice
// @Accept json
// @Produce json
// @Param request body models.UserOffice true "User office data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-office [post]
func createUserOfficeDoc(models.UserOffice) {}

// findUserOfficeDoc godoc
// @Summary Get all user offices
// @Description Get list of user office mappings with pagination
// @Tags UserOffice
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-office [get]
func findUserOfficeDoc() {}

// updateUserOfficeDoc godoc
// @Summary Update user office
// @Description Update existing user office mapping
// @Tags UserOffice
// @Accept json
// @Produce json
// @Param request body models.UserOffice true "User office data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-office [put]
func updateUserOfficeDoc(models.UserOffice) {}

// deleteUserOfficeDoc godoc
// @Summary Delete user office
// @Description Delete user office mapping by ID
// @Tags UserOffice
// @Produce json
// @Param id path string true "User office ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-office/{id} [delete]
func deleteUserOfficeDoc() {}

// createUserCustomerDoc godoc
// @Summary Create user customer
// @Description Create new user customer mapping
// @Tags UserCustomer
// @Accept json
// @Produce json
// @Param request body models.UserCustomer true "User customer data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-customer [post]
func createUserCustomerDoc(models.UserCustomer) {}

// findUserCustomerDoc godoc
// @Summary Get all user customers
// @Description Get list of user customer mappings with pagination
// @Tags UserCustomer
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-customer [get]
func findUserCustomerDoc() {}

// updateUserCustomerDoc godoc
// @Summary Update user customer
// @Description Update existing user customer mapping
// @Tags UserCustomer
// @Accept json
// @Produce json
// @Param request body models.UserCustomer true "User customer data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-customer [put]
func updateUserCustomerDoc(models.UserCustomer) {}

// deleteUserCustomerDoc godoc
// @Summary Delete user customer
// @Description Delete user customer mapping by ID
// @Tags UserCustomer
// @Produce json
// @Param id path string true "User customer ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-customer/{id} [delete]
func deleteUserCustomerDoc() {}
