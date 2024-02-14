package ui

import (
	"github.com/payment-gateway/internal/domain/model"
	"strconv"
)

// validateSessionRequest validates the Session data
func validateSessionRequest(request SessionDTO) *model.DomainError {
	// Validate MerchantID
	if request.MerchantID == "" {
		return &model.DomainError{Category: model.FunctionalError, Message: "MerchantID is required"}
	}
	// Validate session ID
	err := validateSessionId(request.SessionID)
	if err != nil {
		return err
	}
	// Validate amount
	err = validateAmount(request.Amount)
	if err != nil {
		return err
	}
	// Validate card information
	err = validateCardInfos(request.CustomerCardInfo)
	if err != nil {
		return err
	}

	return nil
}

func validateSessionId(sessionId string) *model.DomainError {
	// Check if sessionId is empty
	if sessionId == "" {
		return &model.DomainError{Category: model.FunctionalError, Message: "sessionID is required"}
	}
	return nil
}

func validateCardInfos(customerCardInfo model.CustomerCardInfo) *model.DomainError {
	// Check if CustomerCardInfo is empty
	if customerCardInfo.CardNumber == "" {
		return &model.DomainError{Category: model.FunctionalError, Message: "cardNumber is required"}
	}
	if customerCardInfo.ExpiryMonth == "" || customerCardInfo.ExpiryYear == "" {
		return &model.DomainError{Category: model.FunctionalError, Message: "expiryMonth and ExpiryYear are required"}
	}
	if customerCardInfo.SecurityCode == "" {
		return &model.DomainError{Category: model.FunctionalError, Message: "securityCode (CVV) is required"}
	}
	return nil
}

func validateAmount(amount model.Amount) *model.DomainError {
	// Check if Amount is empty or not a valid numeric value
	if amount.Value == "" {
		return &model.DomainError{Category: model.FunctionalError, Message: "amount is required"}
	}
	if _, err := strconv.ParseFloat(amount.Value, 64); err != nil {
		return &model.DomainError{Category: model.FunctionalError, Message: "amount must be a valid numeric value"}
	}

	// Check if Currency is empty
	if amount.Currency == "" {
		return &model.DomainError{Category: model.FunctionalError, Message: "currency is required"}
	}
	return nil
}
