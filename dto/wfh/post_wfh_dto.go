package wfh

import "github.com/google/uuid"

type PostWfhDto struct {
	WFHId     uuid.UUID `json:"wfh_id"`
	UserId    uuid.UUID `json:"user_id" form:"user_id"`
	Remarks   string    `json:"remarks"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
}