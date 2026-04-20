package approval

import "github.com/google/uuid"

type ApprovalDetailResponseDto struct {
	ApprovalTemplateDetailId uuid.UUID
	ApprovalTemplateHeaderId uuid.UUID
	ApproverBy               uuid.UUID
	SequenceNumber           int

	ApprovalDetailId *uuid.UUID
	ApprovalStatus   *string
	Remark           *string
}