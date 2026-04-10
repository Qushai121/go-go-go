package mappers

import (
	"time"

	"hrms_go/dto/wfh"
	"hrms_go/models"
)

func ToWFHModel(dto wfh.PostWfhDto) (models.WFH, error) {
	var model models.WFH
	var err error

	model.WFHId = dto.WFHId

	model.UserId = dto.UserId
	model.Remarks = dto.Remarks

	layout := "2006-01-02"

	if dto.StartTime != "" {
		model.StartTime, err = time.Parse(layout, dto.StartTime)
		if err != nil {
			return model, err
		}
	}

	if dto.EndTime != "" {
		model.EndTime, err = time.Parse(layout, dto.EndTime)
		if err != nil {
			return model, err
		}
	}

	return model, nil
}