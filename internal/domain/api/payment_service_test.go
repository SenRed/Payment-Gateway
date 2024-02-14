package api

import (
	"github.com/golang/mock/gomock"
	"github.com/payment-gateway/internal/domain/model"
	"github.com/payment-gateway/internal/infrastructure/mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPaymentService_CreateNewSession(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockIPaymentRepository(ctrl)
	mockProcessor := mock.NewMockIPaymentProcessor(ctrl)
	service := NewPaymentService(mockRepo, mockProcessor)

	mockRepo.EXPECT().IsSessionUnique("session123").Return(true, nil)
	mockRepo.EXPECT().CreateSession(gomock.Any()).Return(nil)

	cardInfo := model.CustomerCardInfo{
		CardNumber:   "1234567890123456",
		ExpiryMonth:  "12",
		ExpiryYear:   "25",
		SecurityCode: "123",
	}
	amount := model.Amount{
		Value:    "100",
		Currency: "EUR",
	}
	// When
	err := service.CreateNewSession(
		"session123",
		"exampleMerchant",
		amount,
		cardInfo)

	// Then
	assert.Empty(t, err, "Unexpected error")
}

func TestPaymentService_CreateNewSession_DuplicateSessionId(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockIPaymentRepository(ctrl)
	mockProcessor := mock.NewMockIPaymentProcessor(ctrl)
	service := NewPaymentService(mockRepo, mockProcessor)

	mockRepo.EXPECT().IsSessionUnique("session123").Return(false, nil)

	cardInfo := model.CustomerCardInfo{
		CardNumber:   "1234567890123456",
		ExpiryMonth:  "12",
		ExpiryYear:   "25",
		SecurityCode: "123",
	}
	amount := model.Amount{
		Value:    "100",
		Currency: "USD",
	}
	expectedError := &model.DomainError{
		Category: model.FunctionalError,
		Type:     model.DuplicatedSessionId,
		Message:  "A session already exist with this id",
	}
	// When
	err := service.CreateNewSession(
		"session123",
		"exampleMerchant",
		amount,
		cardInfo)

	// Then
	assert.Equal(t, expectedError, err)
}

func TestPaymentService_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockIPaymentRepository(ctrl)
	mockProcessor := mock.NewMockIPaymentProcessor(ctrl)
	service := NewPaymentService(mockRepo, mockProcessor)

	t.Run("SessionNotFound", func(t *testing.T) {
		expectedError := &model.DomainError{
			Category: model.FunctionalError,
			Message:  "Session not found",
		}
		mockRepo.EXPECT().GetSessionById("session123").Return(nil, expectedError)

		result, err := service.Start("session123")

		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})

	t.Run("SuccessStatus", func(t *testing.T) {
		session := &model.Session{Status: model.SuccessStatus}
		mockRepo.EXPECT().GetSessionById("session123").Return(session, nil)

		expectedResult := &model.TransactionResult{
			Status:  model.SuccessStatus,
			Message: "The payment has already started and finished successfully",
		}

		result, err := service.Start("session123")

		assert.Equal(t, expectedResult, result)
		assert.Nil(t, err)
	})

	t.Run("FailureStatus", func(t *testing.T) {
		session := &model.Session{Status: model.FailureStatus}
		mockRepo.EXPECT().GetSessionById("session123").Return(session, nil)

		expectedResult := &model.TransactionResult{
			Status:  model.FailureStatus,
			Message: "The payment has already started and failed",
		}

		result, err := service.Start("session123")

		assert.Equal(t, expectedResult, result)
		assert.Nil(t, err)
	})

}

func TestPaymentService_processSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockIPaymentRepository(ctrl)
	mockProcessor := mock.NewMockIPaymentProcessor(ctrl)
	service := NewPaymentService(mockRepo, mockProcessor)

	t.Run("ProcessorSuccess", func(t *testing.T) {
		session := model.Session{}
		processorResult := &model.TransactionResult{
			Status:  model.SuccessStatus,
			Code:    model.SuccessfulTransactionCode,
			Message: "Success message",
		}
		mockProcessor.EXPECT().Process(session).Return(processorResult, nil)
		mockRepo.EXPECT().UpdateSession(gomock.Any()).Return(nil)

		expectedResult := &model.TransactionResult{
			Status:  model.SuccessStatus,
			Code:    model.SuccessfulTransactionCode,
			Message: "The transaction finished successfully",
		}

		result, err := service.processSession(session)

		assert.Nil(t, err)

		assert.Equal(t, expectedResult.Status, result.Status)
		assert.Equal(t, expectedResult.Code, result.Code)
		assert.Equal(t, expectedResult.Message, result.Message)
	})

	t.Run("ProcessorFailure", func(t *testing.T) {
		session := model.Session{}
		processorResult := &model.TransactionResult{
			Status:  model.FailureStatus,
			Code:    "ErrorCode",
			Message: "Failure message",
		}
		mockProcessor.EXPECT().Process(session).Return(processorResult, nil)
		mockRepo.EXPECT().UpdateSession(gomock.Any()).Return(nil)

		expectedResult := &model.TransactionResult{
			Status:  model.FailureStatus,
			Code:    "ErrorCode",
			Message: "Failure message",
		}

		result, err := service.processSession(session)

		assert.Equal(t, expectedResult, result)
		assert.Nil(t, err)
	})

	t.Run("ProcessorError", func(t *testing.T) {
		session := model.Session{}
		processorError := &model.DomainError{
			Category: model.TechnicalError,
			Message:  "Processor error",
		}
		mockProcessor.EXPECT().Process(session).Return(nil, processorError)

		result, err := service.processSession(session)

		assert.Nil(t, result)
		assert.Equal(t, processorError, err)
	})

	t.Run("UpdateSessionError", func(t *testing.T) {
		session := model.Session{}
		processorResult := &model.TransactionResult{
			Status:  model.SuccessStatus,
			Code:    model.SuccessfulTransactionCode,
			Message: "Success message",
		}
		mockProcessor.EXPECT().Process(gomock.Any()).Return(processorResult, nil)
		mockRepo.EXPECT().UpdateSession(gomock.Any()).Return(&model.DomainError{
			Category: model.TechnicalError,
			Message:  "failed to create a new session",
		})

		result, err := service.processSession(session)

		assert.Nil(t, result)
		assert.NotNil(t, err)
	})
}

func TestPaymentService_GetPaymentDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockIPaymentRepository(ctrl)
	mockProcessor := mock.NewMockIPaymentProcessor(ctrl)
	service := NewPaymentService(mockRepo, mockProcessor)

	t.Run("SessionNotFound", func(t *testing.T) {
		sessionID := "session123"
		expectedError := &model.DomainError{
			Category: model.FunctionalError,
			Type:     model.SessionNotFound,
			Message:  "Session not found",
		}
		mockRepo.EXPECT().GetSessionById(sessionID).Return(nil, expectedError)

		result, err := service.GetPaymentDetails(sessionID)

		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})

	t.Run("Success", func(t *testing.T) {
		sessionID := "session123"
		now := time.Now()
		expectedTransactions := []model.Transaction{
			{
				Status:  model.SuccessStatus,
				Code:    model.SuccessfulTransactionCode,
				Message: "Success message",
				Date:    now,
			},
			{
				Status:  model.FailureStatus,
				Code:    "ErrorCode",
				Message: "Failure message",
				Date:    now,
			},
		}
		expectedSession := &model.Session{
			Transactions: expectedTransactions,
		}
		mockRepo.EXPECT().GetSessionById(sessionID).Return(expectedSession, nil)

		expectedResults := []model.TransactionResult{
			{
				Status:  model.SuccessStatus,
				Code:    model.SuccessfulTransactionCode,
				Message: "Success message",
				Date:    &now,
			},
			{
				Status:  model.FailureStatus,
				Code:    "ErrorCode",
				Message: "Failure message",
				Date:    &now,
			},
		}

		result, err := service.GetPaymentDetails(sessionID)

		assert.Equal(t, expectedResults, result)
		assert.Nil(t, err)
	})
}
