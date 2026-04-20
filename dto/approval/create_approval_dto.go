package approval

import (
	"hrms_go/constant"

	"github.com/google/uuid"
)

type CreateApprovalDto struct {
	TemplateType  constant.ApprovalTemplateTypeConstant
	RequesterBy   string
	ApprovalDocId uuid.UUID
}