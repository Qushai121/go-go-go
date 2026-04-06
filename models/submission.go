package models

import (
	"time"

	"github.com/google/uuid"
)

type Submission struct {
	SubmissionID  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"submission_id"`
	UserId        uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	RequestNumber string    `json:"request_number"`
	SubmitDate    time.Time `json:"submit_date"`
	Status        string    `json:"status"`
	Remarks       string    `json:"remarks"`
	Amount        float64   `json:"amount"`
	Claims        []Claim   `gorm:"foreignKey:SubmissionID" json:"claims"`
	User          User      `gorm:"foreignKey:UserId;references:UserId"`
}

func (Submission) TableName() string {
	return "hrms_submissions"
}
