package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserId          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"user_id"`
	EmployeeNIK string    `json:"employee_nik"`
	FullName    string    `json:"fullname"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	Role        string    `json:"role"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
}

func (User) TableName() string {
	return "hrms_users"
}
