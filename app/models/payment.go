package models

import (
	"encoding/json"
	// "strconv"
	// "strings"
	"time"

	// "github.com/google/uuid"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Payment struct {
	ID                string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	Order             Order
	OrderID           string           `gorm:"size:36;index"`
	Number            string           `gorm:"size:100;index"`
	Amount            decimal.Decimal  `gorm:"type:decimal(16,2)"`
	TransactionID     string           `gorm:"size:100;index"`
	TransactionStatus string           `gorm:"size:100;index"`
	Payload           *json.RawMessage `gorm:"type:json;not null;default:'{}'"`
	PaymentType       string           `gorm:"size:100"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt
}