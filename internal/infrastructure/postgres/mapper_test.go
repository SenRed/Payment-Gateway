package postgres

import (
	"github.com/payment-gateway/internal/domain/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMapToSessionEntity(t *testing.T) {
	// Given
	now := time.Now()
	session := model.Session{
		SessionID:  "session123",
		MerchantID: "merchant123",
		Amount: model.Amount{
			Currency: "USD",
			Value:    "100",
		},
		CustomerCardInfo: model.CustomerCardInfo{
			CardNumber:   "1234567890123456",
			ExpiryMonth:  "12",
			ExpiryYear:   "2025",
			SecurityCode: "123",
		},
		Transactions: []model.Transaction{
			{
				Status:  model.NotStartedStatus,
				Code:    model.NotStartedCode,
				Message: "Transactions not started yet",
				Date:    now,
			},
			{
				Status:  model.FailureStatus,
				Code:    model.InvalidCVCCode,
				Message: "Invalid CVC",
				Date:    now,
			},
		},
	}

	expectedSessionEntity := SessionEntity{
		ID:         "session123",
		MerchantID: "merchant123",
		Amount: Amount{
			Currency: "USD",
			Value:    "100",
		},
		CustomerCardInfo: CustomerCardInfo{
			CardNumber:   "1234567890123456",
			ExpiryMonth:  "12",
			ExpiryYear:   "2025",
			SecurityCode: "123",
		},
		TransactionEntities: []TransactionEntity{
			{
				Status:  string(model.NotStartedStatus),
				Code:    string(model.NotStartedCode),
				Message: "Transactions not started yet",
				Date:    now,
			},
			{
				Status:  string(model.FailureStatus),
				Code:    string(model.InvalidCVCCode),
				Message: "Invalid CVC",
				Date:    now,
			},
		},
	}
	// When
	actualSessionEntity := mapToSessionEntity(session)
	// Then
	assert.Equal(t, expectedSessionEntity, actualSessionEntity)
}
