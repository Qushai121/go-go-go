package approval

type PostApproveDto struct {
	ApprovalHeaderId string
	ApprovalDetailId string
	ApproverBy       string
	ApprovalStatus   string
	Remark           *string
}