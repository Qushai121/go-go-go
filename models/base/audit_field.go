package base

import (
	"time"

	"github.com/google/uuid"
)

type AuditFields struct {
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    CreatedBy uuid.UUID `json:"created_by"`
    UpdatedBy uuid.UUID `json:"updated_by"`
}