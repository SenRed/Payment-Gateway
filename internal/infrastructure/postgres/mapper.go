package postgres

import "github.com/payment-gateway/internal/domain/model"

func mapToSessionEntity(session model.Session) SessionEntity {
	return SessionEntity{
		ID:                  session.SessionID,
		MerchantID:          session.MerchantID,
		Status:              string(session.Status),
		Amount:              Amount(session.Amount),
		CustomerCardInfo:    CustomerCardInfo(session.CustomerCardInfo),
		TransactionEntities: mapToTransactions(session.Transactions),
	}
}

func mapToTransactions(transactions []model.Transaction) []TransactionEntity {
	res := make([]TransactionEntity, len(transactions))
	for i, transaction := range transactions {
		res[i] = TransactionEntity{
			ID:      transaction.ID,
			Status:  string(transaction.Status),
			Code:    string(transaction.Code),
			Message: transaction.Message,
			Date:    transaction.Date,
		}
	}
	return res
}
