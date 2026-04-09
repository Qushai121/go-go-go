package leave

type PostLeaveDto struct {
	LeaveID            string `json:"leave_id"`
	RequestNumber      string    `json:"request_number"`
	EmployeeName       string    `json:"employee_name"`
	LeaveType          string    `json:"leave_type"`
	StartDate          string `json:"start_date"`
	EndDate            string `json:"end_date"`
	ReqDate            string `json:"req_date"`
	Status             string    `json:"status"`
	Message            string    `json:"message"`
	CancellationReason string    `json:"cancellation_reason"`
	LeaveBalance       float64   `json:"leave_balance"`
}