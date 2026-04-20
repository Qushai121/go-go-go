package constant

type ApprovalStatus string

const (
	PENDING  ApprovalStatus = "PENDING"
	APPROVED ApprovalStatus = "APPROVED"
	REJECTED ApprovalStatus = "REJECTED"
)