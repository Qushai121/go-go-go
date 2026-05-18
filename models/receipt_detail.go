package models

import (
	"hrms_go/models/base"
	"time"

	"github.com/google/uuid"
)

type ReceiptDetail struct {
	ReceiptDetailId    uuid.UUID `gorm:"column:receipt_detail_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"receipt_detail_id"`
	ReceiptId          uuid.UUID `gorm:"column:receipt_id;type:uuid;not null" json:"receipt_id"`
	ReceiptDate        time.Time `gorm:"column:receipt_date;not null" json:"receipt_date"`
	ReceiptType        string    `gorm:"column:receipt_type;type:varchar(100);not null" json:"receipt_type"`
	ReceiptAmount      float64   `gorm:"column:receipt_amount;type:numeric(19,2);not null" json:"receipt_amount"`
	ReceiptDescription string    `gorm:"column:receipt_description;type:text;not null" json:"receipt_description"`
	ReceiptImage       *string   `gorm:"column:receipt_image;type:varchar(255)" json:"receipt_image"`
	ObjectCode         string    `gorm:"column:object_code;type:varchar(50);default:RECEIPT_DETAIL" json:"object_code"`

	base.AuditFields
}

func (ReceiptDetail) TableName() string {
	return "hrms_receipt_detail"
}
