package models

import (
	"time"
)

type Transaction struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id"`
	Amount    int       `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e *Transaction) TableName() string {
	return "transactions"
}
