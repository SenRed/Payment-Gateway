package api

import "github.com/payment-gateway/internal/domain/model"

func mapToTransactionsResult(transactions []model.Transaction) []model.TransactionResult {
	res := make([]model.TransactionResult, len(transactions))
	for i, transaction := range transactions {
		res[i] = model.TransactionResult{
			Status:  transaction.Status,
			Code:    transaction.Code,
			Message: transaction.Message,
			Date:    &transaction.Date,
		}
	}
	return res
}
