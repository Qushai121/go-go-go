package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserId          uuid.UUID  `gorm:"column:user_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"user_id"`
	UserGroupId     *uuid.UUID `gorm:"column:usergroup_id;type:uuid" json:"usergroup_id"`
	CompanyCode     *string    `gorm:"column:company_code;type:varchar(100)" json:"company_code"`
	BranchCode      *string    `gorm:"column:branch_code;type:varchar(100)" json:"branch_code"`
	OfficeCode      *string    `gorm:"column:office_code;type:varchar(100)" json:"office_code"`
	DivisionCode    *string    `gorm:"column:division_code;type:varchar(100)" json:"division_code"`
	DepartmentCode  *string    `gorm:"column:department_code;type:varchar(100)" json:"department_code"`
	TitleCode       *string    `gorm:"column:title_code;type:varchar(100)" json:"title_code"`
	EmployeeNIK     string     `gorm:"column:employee_nik;type:varchar(50);not null;unique" json:"employee_nik"`
	Fullname        string     `gorm:"column:fullname;type:varchar(150);not null" json:"fullname"`
	Email           string     `gorm:"column:email;type:varchar(150);not null;unique" json:"email"`
	Password        string     `gorm:"column:password;type:text;not null" json:"password"`
	IsActive        string     `gorm:"column:is_active;type:varchar(1);default:Y" json:"is_active"`
	IsLocked        string     `gorm:"column:is_locked;type:varchar(1);default:Y" json:"is_locked"`
	LockedDate      *time.Time `gorm:"column:locked_date" json:"locked_date"`
	LockedByUser    *string    `gorm:"column:locked_by_user;type:varchar(50)" json:"locked_by_user"`
	NeedReset       *string    `gorm:"column:need_reset;type:varchar(1)" json:"need_reset"`
	PasswordExpDate *time.Time `gorm:"column:password_exp_date;type:date" json:"password_exp_date"`
	FailedAttempt   *int       `gorm:"column:failed_attempt" json:"failed_attempt"`
	ExpandRadius    int        `gorm:"column:expand_radius;default:0" json:"expand_radius"`
	ObjectCode      string     `gorm:"column:object_code;type:varchar(10);default:USR" json:"object_code"`

	Role string `gorm:"-" json:"role,omitempty"`

	UserCompany      []UserCompany      `gorm:"-" json:"user_company,omitempty"`
	UserOffice       []UserOffice       `gorm:"-" json:"user_office,omitempty"`
	UserCustomer     []UserCustomer     `gorm:"-" json:"user_customer,omitempty"`
	UserShift        []UserShiftMapping `gorm:"-" json:"user_shift"`
	UserLeaveBalance []LeaveBalance     `gorm:"-" json:"user_leave_balance"`

	ProfilePictureUrl string `gorm:"column:profile_picture_url;type:text" json:"profile_picture_url,omitempty"`

	CreatedAt time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	CreatedBy string     `gorm:"column:created_by;type:varchar(50)" json:"created_by"`
	UpdatedBy *string    `gorm:"column:updated_by;type:varchar(50)" json:"updated_by"`
}

func (u *User) getModelSave() {

}

func (User) TableName() string {
	return "hrms_users"
}
