package controllers

import (
	"encoding/json"
	"fmt"
	"hrms_go/dto"
	receiptDto "hrms_go/dto/receipt"
	"hrms_go/repositories"
	"hrms_go/utils"
	"mime/multipart"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

type ReceiptController struct {
	repo repositories.ReceiptRepository
}

func NewReceiptController(repo repositories.ReceiptRepository) *ReceiptController {
	return &ReceiptController{repo}
}

// Create Receipt godoc
// @Summary Create receipt with details and submission
// @Description Create receipt header, details, and submission in one transaction
// @Tags Receipt
// @Accept json
// @Accept multipart/form-data
// @Produce json
// @Param request body receipt.CreateReceiptRequestDto true "Receipt payload"
// @Param employee_nik formData string false "Employee NIK. If empty, taken from auth token."
// @Param receipt_create_date formData string false "Receipt create date (RFC3339 or YYYY-MM-DD HH:mm:ss)"
// @Param details formData string false "JSON array of receipt details"
// @Param receipt_image_0 formData file false "Receipt image for detail index 0"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/receipt [post]
func (c *ReceiptController) Create(ctx fiber.Ctx) error {
	req, savedFiles, err := parseCreateReceiptRequest(ctx)
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	employeeNik := getStringLocal(ctx, "employee_nik")
	if strings.TrimSpace(req.EmployeeNik) == "" {
		req.EmployeeNik = employeeNik
	}

	if err := validateCreateReceipt(req); err != nil {
		removeSavedFiles(savedFiles)
		return utils.Error(ctx, 422, err.Error())
	}

	userId := getStringLocal(ctx, "user_id")
	createdBy := getStringLocal(ctx, "employee_nik")
	if createdBy == "" {
		createdBy = req.EmployeeNik
	}

	data, err := c.repo.Create(req, createdBy, userId)
	if err != nil {
		removeSavedFiles(savedFiles)
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Get All Receipt godoc
// @Summary Get all receipts
// @Description Get list of receipts with pagination and search
// @Tags Receipt
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
// @Router /api/receipt [get]
func (c *ReceiptController) FindAll(ctx fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}

	if err := utils.BindPaginationParams(ctx, &queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	data, err := c.repo.FindAll(&queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Get Draft Receipt Details godoc
// @Summary Get draft receipt details
// @Description Get receipt details that have not been submitted into receipt header
// @Tags Receipt
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
// @Router /api/receipt/detail/draft [get]
func (c *ReceiptController) FindDraftDetails(ctx fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}

	if err := utils.BindPaginationParams(ctx, &queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	data, err := c.repo.FindDraftDetails(&queryParams, getStringLocal(ctx, "employee_nik"))
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Submit Receipt godoc
// @Summary Submit selected receipt details
// @Description Create receipt header and submission from selected draft receipt details
// @Tags Receipt
// @Accept json
// @Produce json
// @Param request body receipt.SubmitReceiptRequestDto true "Submit receipt payload"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/receipt/submit [post]
func (c *ReceiptController) Submit(ctx fiber.Ctx) error {
	var req receiptDto.SubmitReceiptRequestDto

	if err := ctx.Bind().Body(&req); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	if strings.TrimSpace(req.EmployeeNik) == "" {
		req.EmployeeNik = getStringLocal(ctx, "employee_nik")
	}

	if err := validateSubmitReceipt(req); err != nil {
		return utils.Error(ctx, 422, err.Error())
	}

	userId := getStringLocal(ctx, "user_id")
	createdBy := getStringLocal(ctx, "employee_nik")
	if createdBy == "" {
		createdBy = req.EmployeeNik
	}

	data, err := c.repo.Submit(req, createdBy, userId)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Detail Receipt godoc
// @Summary Get receipt detail
// @Description Get receipt header with receipt details and submission
// @Tags Receipt
// @Accept json
// @Produce json
// @Param id path string true "Receipt ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/receipt/{id} [get]
func (c *ReceiptController) Detail(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	data, err := c.repo.Detail(id)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Update Receipt Header godoc
// @Summary Update receipt header
// @Description Update receipt header basic data
// @Tags Receipt
// @Accept json
// @Produce json
// @Param id path string true "Receipt ID"
// @Param request body receipt.UpdateReceiptHeaderDto true "Receipt header payload"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/receipt/{id} [put]
func (c *ReceiptController) UpdateHeader(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var req receiptDto.UpdateReceiptHeaderDto

	if err := ctx.Bind().Body(&req); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	employeeNIK := getStringLocal(ctx, "employee_nik")
	data, err := c.repo.UpdateHeader(id, req, employeeNIK)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Update Receipt Submission godoc
// @Summary Update receipt submission workflow
// @Description Update submission status, current step, and approval header
// @Tags Receipt
// @Accept json
// @Produce json
// @Param id path string true "Receipt ID"
// @Param request body receipt.UpdateReceiptSubmissionDto true "Receipt submission payload"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/receipt/{id}/submission [put]
func (c *ReceiptController) UpdateSubmission(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var req receiptDto.UpdateReceiptSubmissionDto

	if err := ctx.Bind().Body(&req); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	if err := validateStatus(req.Status); err != nil {
		return utils.Error(ctx, 422, err.Error())
	}
	if req.CurrentStep <= 0 {
		return utils.Error(ctx, 422, "current_step must be greater than 0")
	}

	data, err := c.repo.UpdateSubmission(id, req, getStringLocal(ctx, "employee_nik"))
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Delete Receipt godoc
// @Summary Delete receipt
// @Description Delete receipt header with details and submission
// @Tags Receipt
// @Accept json
// @Produce json
// @Param id path string true "Receipt ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/receipt/{id} [delete]
func (c *ReceiptController) Delete(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.repo.Delete(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, "deleted")
}

// Create Receipt Detail godoc
// @Summary Create receipt detail
// @Description Add draft receipt detail before creating receipt header
// @Tags Receipt
// @Accept json
// @Accept multipart/form-data
// @Produce json
// @Param request body receipt.CreateReceiptDetailDto true "Receipt detail payload"
// @Param receipt_date formData string false "Receipt date (RFC3339 or YYYY-MM-DD HH:mm:ss)"
// @Param receipt_type formData string false "Receipt type"
// @Param receipt_amount formData number false "Receipt amount"
// @Param receipt_description formData string false "Receipt description"
// @Param employee_nik formData string false "Employee NIK. If empty, taken from auth token."
// @Param receipt_image formData file false "Receipt image"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/receipt/detail [post]
func (c *ReceiptController) CreateDraftDetail(ctx fiber.Ctx) error {
	req, savedFile, err := parseCreateReceiptDetailRequest(ctx)
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	if err := validateReceiptDetail(req); err != nil {
		removeSavedFile(savedFile)
		return utils.Error(ctx, 422, err.Error())
	}

	createdBy := getStringLocal(ctx, "employee_nik")
	if createdBy == "" {
		createdBy = getReceiptEmployeeNik(ctx)
	}

	data, err := c.repo.CreateDraftDetail(req, createdBy)
	if err != nil {
		removeSavedFile(savedFile)
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Create Receipt Detail godoc
// @Summary Create receipt detail into existing receipt
// @Description Add receipt detail and recalculate header totals
// @Tags Receipt
// @Accept json
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Receipt ID"
// @Param request body receipt.CreateReceiptDetailDto true "Receipt detail payload"
// @Param receipt_date formData string false "Receipt date (RFC3339 or YYYY-MM-DD HH:mm:ss)"
// @Param receipt_type formData string false "Receipt type"
// @Param receipt_amount formData number false "Receipt amount"
// @Param receipt_description formData string false "Receipt description"
// @Param employee_nik formData string false "Employee NIK. If empty, taken from auth token."
// @Param receipt_image formData file false "Receipt image"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/receipt/{id}/detail [post]
func (c *ReceiptController) CreateDetail(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	req, savedFile, err := parseCreateReceiptDetailRequest(ctx)
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	if err := validateReceiptDetail(req); err != nil {
		removeSavedFile(savedFile)
		return utils.Error(ctx, 422, err.Error())
	}

	data, err := c.repo.CreateDetail(id, req, getStringLocal(ctx, "employee_nik"))
	if err != nil {
		removeSavedFile(savedFile)
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Update Receipt Detail godoc
// @Summary Update receipt detail
// @Description Update receipt detail and recalculate header totals
// @Tags Receipt
// @Accept json
// @Accept multipart/form-data
// @Produce json
// @Param detail_id path string true "Receipt Detail ID"
// @Param request body receipt.UpdateReceiptDetailDto true "Receipt detail payload"
// @Param receipt_date formData string false "Receipt date (RFC3339 or YYYY-MM-DD HH:mm:ss)"
// @Param receipt_type formData string false "Receipt type"
// @Param receipt_amount formData number false "Receipt amount"
// @Param receipt_description formData string false "Receipt description"
// @Param employee_nik formData string false "Employee NIK. If empty, taken from auth token."
// @Param receipt_image formData file false "Receipt image"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/receipt/detail/{detail_id} [put]
func (c *ReceiptController) UpdateDetail(ctx fiber.Ctx) error {
	id := ctx.Params("detail_id")
	req, savedFile, err := parseUpdateReceiptDetailRequest(ctx)
	if err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	if err := validateUpdateReceiptDetail(req); err != nil {
		removeSavedFile(savedFile)
		return utils.Error(ctx, 422, err.Error())
	}

	data, err := c.repo.UpdateDetail(id, req, getStringLocal(ctx, "employee_nik"))
	if err != nil {
		removeSavedFile(savedFile)
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Delete Receipt Detail godoc
// @Summary Delete receipt detail
// @Description Delete receipt detail and recalculate header totals
// @Tags Receipt
// @Accept json
// @Produce json
// @Param detail_id path string true "Receipt Detail ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/receipt/detail/{detail_id} [delete]
func (c *ReceiptController) DeleteDetail(ctx fiber.Ctx) error {
	id := ctx.Params("detail_id")

	if err := c.repo.DeleteDetail(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, "deleted")
}

func validateCreateReceipt(req receiptDto.CreateReceiptRequestDto) error {
	if strings.TrimSpace(req.EmployeeNik) == "" {
		return fiber.NewError(422, "employee_nik is required")
	}
	if len(req.Details) == 0 {
		return fiber.NewError(422, "details is required")
	}
	for _, detail := range req.Details {
		if err := validateReceiptDetail(detail); err != nil {
			return err
		}
	}
	return nil
}

func validateSubmitReceipt(req receiptDto.SubmitReceiptRequestDto) error {
	if strings.TrimSpace(req.EmployeeNik) == "" {
		return fiber.NewError(422, "employee_nik is required")
	}
	if len(req.ReceiptDetailIds) == 0 {
		return fiber.NewError(422, "receipt_detail_ids is required")
	}
	return nil
}

func parseCreateReceiptRequest(ctx fiber.Ctx) (receiptDto.CreateReceiptRequestDto, []string, error) {
	var req receiptDto.CreateReceiptRequestDto
	if !isMultipartRequest(ctx) {
		if err := ctx.Bind().Body(&req); err != nil {
			return req, nil, err
		}
		return req, nil, nil
	}

	req.EmployeeNik = ctx.FormValue("employee_nik")
	if req.EmployeeNik == "" {
		req.EmployeeNik = getStringLocal(ctx, "employee_nik")
	}

	if value := ctx.FormValue("receipt_create_date"); value != "" {
		parsed, err := parseReceiptTime(value)
		if err != nil {
			return req, nil, fmt.Errorf("invalid receipt_create_date")
		}
		req.ReceiptCreateDate = &parsed
	}

	if err := json.Unmarshal([]byte(ctx.FormValue("details")), &req.Details); err != nil {
		return req, nil, fmt.Errorf("invalid details json")
	}

	savedFiles := make([]string, 0)
	for index := range req.Details {
		file, err := receiptDetailFile(ctx, index)
		if err != nil {
			continue
		}

		fileUrl, err := saveReceiptImage(ctx, file, req.EmployeeNik, req.Details[index].ReceiptDate)
		if err != nil {
			removeSavedFiles(savedFiles)
			return req, nil, err
		}
		req.Details[index].ReceiptImage = fileUrl
		savedFiles = append(savedFiles, *fileUrl)
	}

	return req, savedFiles, nil
}

func parseCreateReceiptDetailRequest(ctx fiber.Ctx) (receiptDto.CreateReceiptDetailDto, *string, error) {
	var req receiptDto.CreateReceiptDetailDto
	if !isMultipartRequest(ctx) {
		if err := ctx.Bind().Body(&req); err != nil {
			return req, nil, err
		}
		return req, nil, nil
	}

	parsed, err := parseReceiptTime(ctx.FormValue("receipt_date"))
	if err != nil {
		return req, nil, fmt.Errorf("invalid receipt_date")
	}

	req = receiptDto.CreateReceiptDetailDto{
		ReceiptDate:        parsed,
		ReceiptType:        ctx.FormValue("receipt_type"),
		ReceiptDescription: ctx.FormValue("receipt_description"),
		ReceiptImage:       stringPtrIfFilled(ctx.FormValue("receipt_image")),
	}

	if amount, err := parseReceiptAmount(ctx.FormValue("receipt_amount")); err == nil {
		req.ReceiptAmount = amount
	} else {
		return req, nil, err
	}

	file, err := ctx.FormFile("receipt_image")
	if err != nil {
		return req, nil, nil
	}

	fileUrl, err := saveReceiptImage(ctx, file, getReceiptEmployeeNik(ctx), req.ReceiptDate)
	if err != nil {
		return req, nil, err
	}
	req.ReceiptImage = fileUrl

	return req, fileUrl, nil
}

func parseUpdateReceiptDetailRequest(ctx fiber.Ctx) (receiptDto.UpdateReceiptDetailDto, *string, error) {
	var req receiptDto.UpdateReceiptDetailDto
	if !isMultipartRequest(ctx) {
		if err := ctx.Bind().Body(&req); err != nil {
			return req, nil, err
		}
		return req, nil, nil
	}

	createReq, fileUrl, err := parseCreateReceiptDetailRequest(ctx)
	if err != nil {
		return req, nil, err
	}

	req = receiptDto.UpdateReceiptDetailDto{
		ReceiptDate:        createReq.ReceiptDate,
		ReceiptType:        createReq.ReceiptType,
		ReceiptAmount:      createReq.ReceiptAmount,
		ReceiptDescription: createReq.ReceiptDescription,
		ReceiptImage:       createReq.ReceiptImage,
	}

	return req, fileUrl, nil
}

func receiptDetailFile(ctx fiber.Ctx, index int) (*multipart.FileHeader, error) {
	if file, err := ctx.FormFile(fmt.Sprintf("receipt_image_%d", index)); err == nil {
		return file, nil
	}
	if file, err := ctx.FormFile(fmt.Sprintf("details[%d].receipt_image", index)); err == nil {
		return file, nil
	}

	return nil, fiber.NewError(404, "file not found")
}

func saveReceiptImage(ctx fiber.Ctx, file *multipart.FileHeader, employeeNik string, receiptDate time.Time) (*string, error) {
	employeeNik = strings.TrimSpace(employeeNik)
	if employeeNik == "" || employeeNik == "<nil>" {
		return nil, fmt.Errorf("employee_nik is required to upload receipt image")
	}

	folderDate := receiptDate.Format("2006-01-02")
	return utils.SaveFileToCustomPath(file, fmt.Sprintf("receipt/%s/%s", folderDate, employeeNik), ctx)
}

func getReceiptEmployeeNik(ctx fiber.Ctx) string {
	if value := strings.TrimSpace(ctx.FormValue("employee_nik")); value != "" {
		return value
	}
	return getStringLocal(ctx, "employee_nik")
}

func isMultipartRequest(ctx fiber.Ctx) bool {
	return strings.Contains(strings.ToLower(ctx.Get("Content-Type")), "multipart/form-data")
}

func parseReceiptTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, fmt.Errorf("date is required")
	}

	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	var lastErr error
	for _, layout := range formats {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, nil
		}
		lastErr = err
	}

	return time.Time{}, lastErr
}

func parseReceiptAmount(value string) (float64, error) {
	var amount float64
	if _, err := fmt.Sscan(strings.TrimSpace(value), &amount); err != nil {
		return 0, fmt.Errorf("invalid receipt_amount")
	}
	return amount, nil
}

func stringPtrIfFilled(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return &value
}

func removeSavedFiles(paths []string) {
	for _, path := range paths {
		_ = utils.RemoveFileFromPath(path)
	}
}

func removeSavedFile(path *string) {
	if path != nil {
		_ = utils.RemoveFileFromPath(*path)
	}
}

func validateReceiptDetail(req receiptDto.CreateReceiptDetailDto) error {
	if req.ReceiptDate.IsZero() {
		return fiber.NewError(422, "receipt_date is required")
	}
	if strings.TrimSpace(req.ReceiptType) == "" {
		return fiber.NewError(422, "receipt_type is required")
	}
	if req.ReceiptAmount <= 0 {
		return fiber.NewError(422, "receipt_amount must be greater than 0")
	}
	if strings.TrimSpace(req.ReceiptDescription) == "" {
		return fiber.NewError(422, "receipt_description is required")
	}
	return nil
}

func validateUpdateReceiptDetail(req receiptDto.UpdateReceiptDetailDto) error {
	return validateReceiptDetail(receiptDto.CreateReceiptDetailDto{
		ReceiptDate:        req.ReceiptDate,
		ReceiptType:        req.ReceiptType,
		ReceiptAmount:      req.ReceiptAmount,
		ReceiptDescription: req.ReceiptDescription,
		ReceiptImage:       req.ReceiptImage,
	})
}

func validateStatus(status string) error {
	switch status {
	case "P", "A", "R", "D":
		return nil
	default:
		return fiber.NewError(422, "invalid status")
	}
}

func getStringLocal(ctx fiber.Ctx, key string) string {
	value, ok := ctx.Locals(key).(string)
	if !ok {
		return ""
	}
	return value
}
