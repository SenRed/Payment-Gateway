package api

import (
	"fmt"
	"github.com/payment-gateway/internal/domain/model"
	"github.com/payment-gateway/internal/domain/port"
)

type PaymentService struct {
	paymentRepository port.IPaymentRepository
	paymentProcessor  port.IPaymentProcessor
}

func NewPaymentService(paymentRepository port.IPaymentRepository, paymentProcessor port.IPaymentProcessor) PaymentService {
	return PaymentService{paymentRepository: paymentRepository, paymentProcessor: paymentProcessor}
}

func (service *PaymentService) CreateNewSession(sessionID string, merchantID string, amount model.Amount, cardInfo model.CustomerCardInfo) *model.DomainError {
	// Verify no session already exist
	isUnique, err := service.paymentRepository.IsSessionUnique(sessionID)
	if err != nil {
		return err
	}
	// If session is not unique then return duplicated session error
	if !isUnique {
		return &model.DomainError{
			Category: model.FunctionalError,
			Type:     model.DuplicatedSessionId,
			Message:  "A session already exist with this id",
		}
	}
	// If everything is good, then create the session by setting the initial status
	session := model.NewSession(sessionID, merchantID, amount, cardInfo)
	// CreateSession the session
	return service.paymentRepository.CreateSession(session)
}

func (service *PaymentService) Start(sessionsId string) (*model.TransactionResult, *model.DomainError) {
	session, err := service.paymentRepository.GetSessionById(sessionsId)
	if err != nil {
		return nil, err
	}
	switch session.Status {
	case model.FailureStatus:
		return &model.TransactionResult{
			Status:  model.FailureStatus,
			Message: "The payment has already started and failed",
		}, nil
	case model.SuccessStatus:
		return &model.TransactionResult{
			Status:  model.SuccessStatus,
			Message: "The payment has already started and finished successfully",
		}, nil
	case model.NotStartedStatus:
		// Process the payment
		return service.processSession(*session)
	default:
		return nil, &model.DomainError{Category: model.FunctionalError, Message: fmt.Sprintf("Unhandled session status: %s", session.Status)}
	}
}

func (service *PaymentService) processSession(session model.Session) (*model.TransactionResult, *model.DomainError) {
	// Step 1 - Call the acquiring bank processor
	processorResult, err := service.paymentProcessor.Process(session)
	if err != nil {
		return nil, err
	}
	// Step 2 - Handle transaction result
	var paymentResult model.TransactionResult

	switch processorResult.Status {
	case model.SuccessStatus:
		// Mark as successful
		session.MarkAsSuccessful()
		paymentResult = model.TransactionResult{
			Status:  model.SuccessStatus,
			Code:    model.SuccessfulTransactionCode,
			Message: "The transaction finished successfully",
		}
	case model.FailureStatus:
		// Mark as failed
		session.MarkAsFailed(*processorResult)
		paymentResult = model.TransactionResult{
			Status:  model.FailureStatus,
			Code:    processorResult.Code,
			Message: processorResult.Message,
		}
	default:
		return nil, &model.DomainError{Category: model.TechnicalError, Message: "Unhandled processor result status"}
	}
	// Step 3 - Update & return the result
	err = service.paymentRepository.UpdateSession(session)
	if err != nil {
		return nil, err
	}
	return &paymentResult, nil
}

func (service *PaymentService) GetPaymentDetails(sessionId string) ([]model.TransactionResult, *model.DomainError) {
	session, err := service.paymentRepository.GetSessionById(sessionId)
	if err != nil {
		return nil, err
	}
	return mapToTransactionsResult(session.Transactions), nil
}
