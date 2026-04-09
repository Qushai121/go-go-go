package mappers

import (
	"time"

	"hrms_go/dto/submission"
	"hrms_go/models"
)

func ToSubmissionModel(dto submission.PostSubmissionDto) (models.Submission, error) {
	var submitDate time.Time
	var err error

	// Handle SubmitDate (string → time.Time)
	if dto.SubmitDate != "" {
		// adjust layout to your format (example: "2006-01-02")
		submitDate, err = time.Parse("2006-01-02", dto.SubmitDate)
		if err != nil {
			return models.Submission{}, err
		}
	} else {
		submitDate = time.Now()
	}

	model := models.Submission{
		SubmissionID:  dto.SubmissionID, // optional if DB auto გენ
		UserId:        dto.UserId,
		RequestNumber: dto.RequestNumber,
		SubmitDate:    submitDate,
		Status:        dto.Status,
		Amount:        dto.Amount,
	}

	// Handle optional remarks
	if dto.Remarks != nil {
		model.Remarks = *dto.Remarks
	}

	return model, nil
}