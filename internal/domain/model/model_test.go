package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMarkAsSuccessful(t *testing.T) {
	// Given
	session := &Session{
		Status:       NotStartedStatus,
		Transactions: []Transaction{{Code: NotStartedCode}},
	}

	// When
	session.MarkAsSuccessful()

	// Then
	assert.Equal(t, SuccessStatus, session.Status)

	assert.Len(t, session.Transactions, 2)

	assert.Equal(t, session.Transactions[1].Status, SuccessStatus)
	assert.Equal(t, session.Transactions[1].Code, SuccessfulTransactionCode)
	assert.Equal(t, session.Transactions[1].Message, "Transaction executed successfully")
	assert.WithinDuration(t, session.Transactions[1].Date, time.Now(), time.Second)
}

func TestMarkAsFailed(t *testing.T) {
	// Given
	session := &Session{
		Status:       NotStartedStatus,
		Transactions: []Transaction{{Code: NotStartedCode}},
	}
	transactionResult := TransactionResult{
		Code:    "ERR001",
		Message: "Insufficient funds",
	}

	// When
	session.MarkAsFailed(transactionResult)

	// Then
	assert.Equal(t, FailureStatus, session.Status)

	assert.Len(t, session.Transactions, 2)

	assert.Equal(t, session.Transactions[1].Status, FailureStatus)
	assert.Equal(t, session.Transactions[1].Code, transactionResult.Code)
	assert.Equal(t, session.Transactions[1].Message, transactionResult.Message)
	assert.WithinDuration(t, session.Transactions[1].Date, time.Now(), time.Second)
}

func TestNewSession(t *testing.T) {
	// Given
	sessionID := "session123"
	merchantID := "exampleMerchant"
	amount := Amount{
		Value:    "100",
		Currency: "EUR",
	}
	customerCardInfo := CustomerCardInfo{
		CardNumber:   "1234567890123456",
		ExpiryMonth:  "12",
		ExpiryYear:   "25",
		SecurityCode: "123",
	}

	// When
	session := NewSession(sessionID, merchantID, amount, customerCardInfo)

	// Then
	expectedSession := Session{
		SessionID:        sessionID,
		MerchantID:       merchantID,
		Amount:           amount,
		CustomerCardInfo: customerCardInfo,
		Status:           NotStartedStatus,
		Transactions: []Transaction{
			{
				Status:  NotStartedStatus,
				Code:    NotStartedCode,
				Message: "Transactions created, not started yet",
			},
		},
	}
	assert.Equal(t, expectedSession.SessionID, session.SessionID)
	assert.Equal(t, expectedSession.MerchantID, session.MerchantID)
	assert.Equal(t, expectedSession.Amount, session.Amount)
	assert.Equal(t, expectedSession.CustomerCardInfo, session.CustomerCardInfo)
	assert.Equal(t, expectedSession.Status, session.Status)

	assert.Len(t, session.Transactions, 1)

	assert.Equal(t, expectedSession.Transactions[0].Status, NotStartedStatus)
	assert.Equal(t, expectedSession.Transactions[0].Code, NotStartedCode)
	assert.Equal(t, expectedSession.Transactions[0].Message, "Transactions created, not started yet")
	assert.WithinDuration(t, time.Now(), session.Transactions[0].Date, time.Second)
}
