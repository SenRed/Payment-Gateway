package postgres

import (
	"testing"
	"time"

	"github.com/payment-gateway/internal/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestTransactionEntity_toModel(t *testing.T) {
	// Given
	entity := TransactionEntity{
		ID:      1,
		Status:  string(model.SuccessStatus),
		Code:    string(model.SuccessfulTransactionCode),
		Message: "Transaction successful",
		Date:    time.Now(),
	}

	// When
	result := entity.toModel()

	// Then
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, model.SuccessStatus, result.Status)
	assert.Equal(t, model.SuccessfulTransactionCode, result.Code)
	assert.Equal(t, "Transaction successful", result.Message)
	assert.WithinDuration(t, time.Now(), result.Date, time.Second)
}

func TestSessionEntity_toModel(t *testing.T) {
	// Given
	entity := SessionEntity{
		ID:         "session123",
		MerchantID: "merchant123",
		Status:     string(model.NotStartedStatus),
		Amount: Amount{
			Currency: "USD",
			Value:    "100",
		},
		CustomerCardInfo: CustomerCardInfo{
			CardNumber:   "1234567890123456",
			ExpiryMonth:  "12",
			ExpiryYear:   "25",
			SecurityCode: "123",
		},
		TransactionEntities: []TransactionEntity{},
	}

	// When
	result := entity.toModel()

	// Then
	assert.Equal(t, "session123", result.SessionID)
	assert.Equal(t, "merchant123", result.MerchantID)
	assert.Equal(t, model.NotStartedStatus, result.Status)
	assert.Equal(t, "USD", result.Amount.Currency)
	assert.Equal(t, "100", result.Amount.Value)
	assert.Equal(t, "1234567890123456", result.CustomerCardInfo.CardNumber)
	assert.Equal(t, "12", result.CustomerCardInfo.ExpiryMonth)
	assert.Equal(t, "25", result.CustomerCardInfo.ExpiryYear)
	assert.Equal(t, "123", result.CustomerCardInfo.SecurityCode)
}
