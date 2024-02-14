package port

import (
	"github.com/payment-gateway/internal/domain/model"
)

type IPaymentRepository interface {
	GetSessionById(sessionId string) (*model.Session, *model.DomainError)
	IsSessionUnique(sessionID string) (bool, *model.DomainError)

	CreateSession(session model.Session) *model.DomainError
	UpdateSession(session model.Session) *model.DomainError
}
