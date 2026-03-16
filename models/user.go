package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserId          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"user_id"`
	EmployeeNIK string    `json:"employee_nik"`
	Fullname    string    `json:"fullname"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Role        string    `json:"role"`

	ShiftId uuid.UUID `gorm:"type:uuid" json:"shift_id"`
	Shift   Shift     `gorm:"foreignKey:ShiftId;references:ShiftId"`

	ProfilePictureUrl string `gorm:"column:profile_picture_url" json:"profile_picture_url"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
}

func (User) TableName() string {
	return "hrms_users"
}
