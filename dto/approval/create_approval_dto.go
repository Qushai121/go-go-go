package approval

import (
	"hrms_go/constant"

	"github.com/google/uuid"
)

type CreateApprovalDto struct {
	TemplateType  constant.ApprovalTemplateTypeConstant
	RequesterBy   string
	CreatedBy     string
	ApprovalDocId uuid.UUID
}
