package receipt

import "time"

type CreateReceiptRequestDto struct {
	EmployeeNik       string                   `json:"employee_nik"`
	ReceiptCreateDate *time.Time               `json:"receipt_create_date"`
	Details           []CreateReceiptDetailDto `json:"details"`
}

type SubmitReceiptRequestDto struct {
	EmployeeNik       string     `json:"employee_nik"`
	ReceiptCreateDate *time.Time `json:"receipt_create_date"`
	ReceiptDetailIds  []string   `json:"receipt_detail_ids"`
}

type CreateReceiptDetailDto struct {
	ReceiptDate        time.Time `json:"receipt_date"`
	ReceiptType        string    `json:"receipt_type"`
	ReceiptAmount      float64   `json:"receipt_amount"`
	ReceiptDescription string    `json:"receipt_description"`
	ReceiptImage       *string   `json:"receipt_image"`
}

type UpdateReceiptHeaderDto struct {
	EmployeeNik       string    `json:"employee_nik"`
	ReceiptCreateDate time.Time `json:"receipt_create_date"`
}

type UpdateReceiptSubmissionDto struct {
	Status           string  `json:"status"`
	CurrentStep      int     `json:"current_step"`
	ApprovalHeaderId *string `json:"approvalheader_id"`
}

type UpdateReceiptDetailDto struct {
	ReceiptDate        time.Time `json:"receipt_date"`
	ReceiptType        string    `json:"receipt_type"`
	ReceiptAmount      float64   `json:"receipt_amount"`
	ReceiptDescription string    `json:"receipt_description"`
	ReceiptImage       *string   `json:"receipt_image"`
}
