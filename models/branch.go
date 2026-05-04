package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type Branch struct {
	BranchId      uuid.UUID `gorm:"column:branch_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"branch_id"`
	CompanyCode   string    `gorm:"column:company_code;type:varchar(50);not null" json:"company_code"`
	BranchCode    string    `gorm:"column:branch_code;type:varchar(50);not null" json:"branch_code"`
	BranchName    string    `gorm:"column:branch_name;type:varchar(100);not null" json:"branch_name"`
	BranchAddress string    `gorm:"column:branch_address;type:text" json:"branch_address"`
	ObjectCode    string    `gorm:"column:object_code;type:varchar(10);default:BRANCH" json:"object_code"`
	TimezoneSet   string    `gorm:"column:timezone_set;type:varchar(50);default:SE Asia Standard Time" json:"timezone_set"`

	base.AuditFields
}

func (Branch) TableName() string {
	return "hrms_branch"
}
