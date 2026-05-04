package base

import "time"

type AuditFields struct {
	CreatedAt time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP;autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy string     `json:"created_by" gorm:"type:varchar(100);default:System"`
	UpdatedBy *string    `json:"updated_by" gorm:"type:varchar(100)"`
}
