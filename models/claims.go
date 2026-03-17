package models

import (
	"time"

	"github.com/google/uuid"
)

type Claim struct {
	ClaimId uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"claim_id"`
	SubmissionID  uuid.UUID `gorm:"type:uuid" json:"id_submission"`
	RequestNumber string    `json:"request_number"`
	EmployeeName  string    `json:"employee_name"`
	ClaimType     string    `json:"claim_type"`
	ReceiptDate   time.Time `json:"receipt_date"`
	Status        string    `json:"status"`
	Remarks       string    `json:"remarks"`
	Amount        float64   `json:"amount"`
}


func (Claim) TableName() string {
	return "hrms_claims"
}