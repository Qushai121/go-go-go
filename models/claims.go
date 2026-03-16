package models

import (
	"time"

	"github.com/google/uuid"
)

type Claim struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SubmissionID  uuid.UUID `gorm:"type:uuid" json:"id_submission"`
	RequestNumber string    `json:"request_number"`
	EmployeeName  string    `json:"employee_name"`
	ClaimType     string    `json:"claim_type"`
	ReceiptDate   time.Time `json:"receipt_date"`
	Status        string    `json:"status"`
	Remarks       string    `json:"remarks"`
	Amount        float64   `json:"amount"`
}