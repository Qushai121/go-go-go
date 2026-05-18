package models

import (
	"hrms_go/models/base"
	"time"

	"github.com/google/uuid"
)

type ReceiptHeader struct {
	ReceiptId          uuid.UUID `gorm:"column:receipt_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"receipt_id"`
	EmployeeNik        string    `gorm:"column:employee_nik;type:varchar(100);not null" json:"employee_nik"`
	ReceiptCreateDate  time.Time `gorm:"column:receipt_create_date;not null" json:"receipt_create_date"`
	TotalReceipt       int       `gorm:"column:total_receipt;not null" json:"total_receipt"`
	TotalAmountReceipt float64   `gorm:"column:total_amount_receipt;type:numeric(19,2);not null" json:"total_amount_receipt"`
	ObjectCode         string    `gorm:"column:object_code;type:varchar(50);default:RECEIPT_HEADER" json:"object_code"`

	ReceiptDetails []ReceiptDetail `gorm:"foreignKey:ReceiptId;references:ReceiptId" json:"receipt_details,omitempty"`
	Submission     *Submission     `gorm:"foreignKey:ReceiptId;references:ReceiptId" json:"submission,omitempty"`

	base.AuditFields
}

func (ReceiptHeader) TableName() string {
	return "hrms_receipt_header"
}
