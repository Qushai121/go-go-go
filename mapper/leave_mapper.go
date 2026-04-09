package mappers

import (
	"time"

	"hrms_go/dto/leave"
	"hrms_go/models"

	"github.com/google/uuid"
)

func ToLeaveModel(dto leave.PostLeaveDto) (models.Leave, error) {
	var model models.Leave
	var err error

	// Handle UUID (optional)
	if dto.LeaveID != "" {
		model.LeaveID, err = uuid.Parse(dto.LeaveID)
		if err != nil {
			return model, err
		}
	} else {
		model.LeaveID = uuid.New()
	}

	// Parse dates (adjust format if needed)
	layout := "2006-01-02"

	if dto.StartDate != "" {
		model.StartDate, err = time.Parse(layout, dto.StartDate)
		if err != nil {
			return model, err
		}
	}

	if dto.EndDate != "" {
		model.EndDate, err = time.Parse(layout, dto.EndDate)
		if err != nil {
			return model, err
		}
	}

	if dto.ReqDate != "" {
		model.ReqDate, err = time.Parse(layout, dto.ReqDate)
		if err != nil {
			return model, err
		}
	}

	// Map simple fields
	model.RequestNumber = dto.RequestNumber
	model.EmployeeName = dto.EmployeeName
	model.LeaveType = dto.LeaveType
	model.Status = dto.Status
	model.Message = dto.Message
	model.CancellationReason = dto.CancellationReason
	model.LeaveBalance = dto.LeaveBalance

	return model, nil
}