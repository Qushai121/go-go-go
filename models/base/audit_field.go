package base

import "time"

type AuditFields struct {
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy string     `json:"created_by"`
	UpdatedBy *string    `json:"updated_by"`
}
