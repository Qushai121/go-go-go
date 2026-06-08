package models

type CreateLeaveTransactionRequest struct {
	EmployeeNik      string `json:"employee_nik" form:"employee_nik"`
	LeaveTypeId      string `json:"leave_type_id" form:"leave_type_id"`
	LeaveStart       string `json:"leave_start" form:"leave_start"`
	LeaveEnd         string `json:"leave_end" form:"leave_end"`
	Remarks          string `json:"remarks" form:"remarks"`
	Location         string `json:"location" form:"location"`
	CurrentStep      int    `json:"current_step" form:"current_step"`
	ApprovalHeaderId string `json:"approvalheader_id" form:"approvalheader_id"`
	CreatedBy        string `json:"created_by" form:"created_by"`
}

type CreateLeaveTransactionResponse struct {
	EmployeeNik       string         `json:"employee_nik"`
	LeaveTypeId       string         `json:"leave_type_id"`
	LeaveTypeCode     string         `json:"leave_type_code"`
	LeaveTypeName     string         `json:"leave_type_name"`
	LeaveStart        string         `json:"leave_start"`
	LeaveEnd          string         `json:"leave_end"`
	InsertedDays      int            `json:"inserted_days"`
	DeductedDays      int            `json:"deducted_days"`
	InsertedDates     []string       `json:"inserted_dates"`
	LeaveTransactions []LeaveHistory `json:"leave_transactions"`
}

type LeaveBalanceResponse struct {
	EmployeeNik    string         `json:"employee_nik"`
	PeriodStart    string         `json:"period_start"`
	PeriodEnd      string         `json:"period_end"`
	TotalRemaining int            `json:"total_remaining"`
	Balances       []LeaveBalance `json:"balances"`
}

type EstimateLeaveRequest struct {
	EmployeeNik string `json:"employee_nik" form:"employee_nik"`
	LeaveTypeId string `json:"leave_type_id" form:"leave_type_id"`
	LeaveStart  string `json:"leave_start" form:"leave_start"`
	LeaveEnd    string `json:"leave_end" form:"leave_end"`
}

type EstimateLeaveResponse struct {
	EmployeeNik        string   `json:"employee_nik"`
	LeaveTypeId        string   `json:"leave_type_id"`
	LeaveTypeCode      string   `json:"leave_type_code"`
	LeaveTypeName      string   `json:"leave_type_name"`
	LeaveStart         string   `json:"leave_start"`
	LeaveEnd           string   `json:"leave_end"`
	UseWorkingDay      bool     `json:"use_working_day"`
	DeductLeaveBalance bool     `json:"deduct_leave_balance"`
	TotalRangeDays     int      `json:"total_range_days"`
	TotalTakenDays     int      `json:"total_taken_days"`
	DeductedDays       int      `json:"deducted_days"`
	RemainingBefore    int      `json:"remaining_before"`
	RemainingAfter     int      `json:"remaining_after"`
	ShortageDays       int      `json:"shortage_days"`
	IsEnoughBalance    bool     `json:"is_enough_balance"`
	PeriodStart        string   `json:"period_start"`
	PeriodEnd          string   `json:"period_end"`
	EstimatedDates     []string `json:"estimated_dates"`
	Message            string   `json:"message"`
}
