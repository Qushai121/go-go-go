package attandance

type PostAttandanceDto struct {
	AttendanceId string `json:"attendance_id" form:"attendance_id"`
	UserId       string `json:"user_id" form:"user_id"`

	CompanyCode string `json:"company_code" form:"company_code"`
	OfficeCode  string `json:"office_code" form:"office_code"`
	LogTime     string `json:"logtime" form:"logtime"`
	FunctionNo  int    `json:"functionno" form:"functionno"`
	ActivityType string `json:"activity_type" form:"activity_type"`

	Latitude            string `json:"latitude" form:"latitude"`
	Longitude           string `json:"longitude" form:"longitude"`
	PresentaseKemiripan string `json:"presentase_kemiripan" form:"presentase_kemiripan"`
	ImagePath           string `json:"imagepath" form:"imagepath"`
	IsOffline           string `json:"is_offline" form:"is_offline"`
	Distance            string `json:"distance" form:"distance"`
	Platforms           string `json:"platforms" form:"platforms"`
	MaxRadius           string `json:"max_radius" form:"max_radius"`
	ExpandRadius        string `json:"expand_radius" form:"expand_radius"`
	ObjectCode          string `json:"object_code" form:"object_code"`

	CreatedAt string `json:"created_at" form:"created_at"`
	UpdatedAt string `json:"updated_at" form:"updated_at"`
	CreatedBy string `json:"created_by" form:"created_by"`
	UpdatedBy string `json:"updated_by" form:"updated_by"`
}
