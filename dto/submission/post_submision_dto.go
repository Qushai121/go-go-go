package submission

import (
	"github.com/google/uuid"
)

type PostSubmissionDto struct {
	SubmissionID  uuid.UUID 		   `json:"submission_id"`
	UserId        uuid.UUID            `json:"user_id" validate:"required"`
	RequestNumber string               `json:"request_number" validate:"required"`
	SubmitDate    string           	   `json:"submit_date,omitempty"`
	Status        string               `json:"status" validate:"required"`
	Remarks       *string              `json:"remarks,omitempty"`
	Amount        float64              `json:"amount" validate:"required,gte=0"`
}