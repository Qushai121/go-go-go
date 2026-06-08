package response

import "time"

type AttendanceMeResponseDto struct {
	AttendanceId        string     `gorm:"column:attendance_id" json:"attendance_id"`
	UserId              string     `gorm:"column:user_id" json:"user_id"`
	CompanyCode         string     `gorm:"column:company_code" json:"company_code"`
	SiteType            string     `gorm:"column:site_type" json:"site_type"`
	SiteCode            string     `gorm:"column:site_code" json:"site_code"`
	SiteName            string     `gorm:"column:site_name" json:"site_name"`
	SitePhone           string     `gorm:"column:site_phone" json:"site_phone"`
	SiteAddress         string     `gorm:"column:site_address" json:"site_address"`
	SiteLatitude        string     `gorm:"column:site_latitude" json:"site_latitude"`
	SiteLongitude       string     `gorm:"column:site_longitude" json:"site_longitude"`
	LogTime             *time.Time `gorm:"column:logtime" json:"logtime"`
	FunctionNo          int        `gorm:"column:functionno" json:"functionno"`
	ActivityType        string     `gorm:"column:activity_type" json:"activity_type"`
	ActionType          string     `gorm:"column:action_type" json:"action_type"`
	Latitude            string     `gorm:"column:latitude" json:"latitude"`
	Longitude           string     `gorm:"column:longitude" json:"longitude"`
	PresentaseKemiripan string     `gorm:"column:presentase_kemiripan" json:"presentase_kemiripan"`
	ImagePath           string     `gorm:"column:imagepath" json:"imagepath"`
	IsOffline           string     `gorm:"column:is_offline" json:"is_offline"`
	Distance            string     `gorm:"column:distance" json:"distance"`
	Platforms           string     `gorm:"column:platforms" json:"platforms"`
	MaxRadius           string     `gorm:"column:max_radius" json:"max_radius"`
	ExpandRadius        string     `gorm:"column:expand_radius" json:"expand_radius"`
	ObjectCode          string     `gorm:"column:object_code" json:"object_code"`
	CreatedAt           *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           *time.Time `gorm:"column:updated_at" json:"updated_at"`
	CreatedBy           string     `gorm:"column:created_by" json:"created_by"`
	UpdatedBy           string     `gorm:"column:updated_by" json:"updated_by"`
}
