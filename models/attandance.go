package models

import (
	"time"

	"github.com/google/uuid"
)

type Attendance struct {
	AttendanceId uuid.UUID `gorm:"column:attendance_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"attendance_id" form:"attendance_id"`

	UserId       uuid.UUID `gorm:"column:user_id;type:uuid;not null;index" json:"user_id" form:"user_id"`
	CompanyCode  string    `gorm:"column:company_code;type:varchar(100);not null" json:"company_code" form:"company_code"`
	SiteType     string    `gorm:"column:site_type;type:varchar(50);not null" json:"site_type" form:"site_type"`
	SiteCode     string    `gorm:"column:site_code;type:varchar(255);not null" json:"site_code" form:"site_code"`
	LogTime      time.Time `gorm:"column:logtime;type:timestamp;not null" json:"logtime" form:"logtime"`
	FunctionNo   int       `gorm:"column:functionno;type:int;not null" json:"functionno" form:"functionno"`
	ActivityType *string   `gorm:"column:activity_type;type:varchar(100)" json:"activity_type" form:"activity_type"`

	Latitude            *string `gorm:"column:latitude;type:varchar(255)" json:"latitude" form:"latitude"`
	Longitude           *string `gorm:"column:longitude;type:varchar(255)" json:"longitude" form:"longitude"`
	PresentaseKemiripan *string `gorm:"column:presentase_kemiripan;type:varchar(100)" json:"presentase_kemiripan" form:"presentase_kemiripan"`
	ImagePath           string  `gorm:"column:imagepath;type:varchar(255);not null" json:"imagepath" form:"imagepath"`
	IsOffline           *string `gorm:"column:is_offline;type:varchar(1)" json:"is_offline" form:"is_offline"`
	Distance            *string `gorm:"column:distance;type:varchar(100)" json:"distance" form:"distance"`
	Platforms           *string `gorm:"column:platforms;type:varchar(100)" json:"platforms" form:"platforms"`
	MaxRadius           *int    `gorm:"column:max_radius;type:int" json:"max_radius" form:"max_radius"`
	ExpandRadius        *int    `gorm:"column:expand_radius;type:int" json:"expand_radius" form:"expand_radius"`
	ObjectCode          string  `gorm:"column:object_code;type:varchar(100);default:'ATTENDANCE'" json:"object_code" form:"object_code"`

	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at" form:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at" form:"updated_at"`
	CreatedBy string     `gorm:"column:created_by;type:varchar(100);default:'system'" json:"created_by" form:"created_by"`
	UpdatedBy *string    `gorm:"column:updated_by;type:varchar(100)" json:"updated_by" form:"updated_by"`

	User User `gorm:"foreignKey:UserId;references:UserId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (Attendance) TableName() string {
	return "hrms_attendance"
}
