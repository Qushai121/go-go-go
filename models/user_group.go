package models

import (
	"hrms_go/models/base"

	"github.com/google/uuid"
)

type UserGroup struct {
	UserGroupId   uuid.UUID `gorm:"column:usergroup_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"usergroup_id"`
	CompanyCode   string    `gorm:"column:company_code;type:varchar(50);not null" json:"company_code"`
	UserGroupCode string    `gorm:"column:usergroup_code;type:varchar(150);not null" json:"usergroup_code"`
	UserGroupName string    `gorm:"column:usergroup_name;type:varchar(150);not null" json:"usergroup_name"`
	ObjectCode    string    `gorm:"column:object_code;type:varchar(10);default:USERGROUP" json:"object_code"`
	TimezoneSet   string    `gorm:"column:timezone_set;type:varchar(50);default:SE Asia Standard Time" json:"timezone_set"`

	base.AuditFields
}

func (UserGroup) TableName() string {
	return "hrms_usergroup"
}
