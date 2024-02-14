package api

import (
	"github.com/payment-gateway/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMapToTransactionsResult(t *testing.T) {
	// Given
	now := time.Now()
	transactions := []model.Transaction{
		{
			Status:  model.SuccessStatus,
			Code:    model.SuccessfulTransactionCode,
			Message: "Transaction executed successfully",
			Date:    now,
		},
		{
			Status:  model.NotStartedStatus,
			Code:    model.NotStartedCode,
			Message: "Transactions created, not started yet",
			Date:    now,
		},
	}
	expectedTransactionsResult := []model.TransactionResult{
		{
			Status:  model.SuccessStatus,
			Code:    model.SuccessfulTransactionCode,
			Message: "Transaction executed successfully",
			Date:    &now,
		},
		{
			Status:  model.NotStartedStatus,
			Code:    model.NotStartedCode,
			Message: "Transactions created, not started yet",
			Date:    &now,
		},
	}

	// When
	transactionResults := mapToTransactionsResult(transactions)
	// Then
	assert.Equalf(t, expectedTransactionsResult, transactionResults, "error when comparing")
}
