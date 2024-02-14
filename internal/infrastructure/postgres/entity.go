package postgres

import (
	"github.com/payment-gateway/internal/domain/model"
	"time"
)

type TransactionEntity struct {
	// ID Transactions ID
	ID uint `gorm:"primaryKey, autoincrement"`
	// Session entity foreign key
	SessionEntityID string
	Status          string
	Code            string
	Message         string
	Date            time.Time
}

func (e *TransactionEntity) toModel() model.Transaction {
	return model.Transaction{
		ID:      e.ID,
		Status:  model.TransactionStatus(e.Status),
		Code:    model.TransactionCode(e.Code),
		Message: e.Message,
		Date:    e.Date,
	}
}

type SessionEntity struct {
	// ID Session ID
	ID                  string `gorm:"primaryKey"`
	MerchantID          string
	Status              string
	Amount              Amount           `gorm:"embedded;embeddedPrefix:amount_"`
	CustomerCardInfo    CustomerCardInfo `gorm:"embedded;embeddedPrefix:card_"`
	TransactionEntities []TransactionEntity
}

func (e *SessionEntity) toModel() model.Session {
	transactions := make([]model.Transaction, len(e.TransactionEntities))
	for i, transactionEntity := range e.TransactionEntities {
		transactions[i] = transactionEntity.toModel()
	}
	return model.Session{
		SessionID:        e.ID,
		MerchantID:       e.MerchantID,
		Status:           model.TransactionStatus(e.Status),
		Amount:           model.Amount(e.Amount),
		CustomerCardInfo: model.CustomerCardInfo(e.CustomerCardInfo),
		Transactions:     transactions,
	}
}

type Amount struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

type CustomerCardInfo struct {
	CardNumber   string `json:"cardNumber"`
	ExpiryMonth  string `json:"expiryMonth"`
	ExpiryYear   string `json:"expiryYear"`
	SecurityCode string `json:"cvv"`
}
