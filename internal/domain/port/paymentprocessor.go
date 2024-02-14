package port

import "github.com/payment-gateway/internal/domain/model"

type IPaymentProcessor interface {
	Process(session model.Session) (*model.TransactionResult, *model.DomainError)
}
