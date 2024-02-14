package ui

import (
	"github.com/payment-gateway/internal/domain/model"
	"testing"
)

func Test_validateParams(t *testing.T) {
	tests := []struct {
		name    string
		request SessionDTO
		wantErr bool
	}{
		{
			name: "Valid request",
			request: SessionDTO{
				MerchantID: "merchant123",
				SessionID:  "9e7873bc-a58a-48eb-bbf3-3f67b48f5cf1",
				Amount:     model.Amount{Currency: "EUR", Value: "10.50"},
				CustomerCardInfo: model.CustomerCardInfo{
					CardNumber:   "1234567890123456",
					ExpiryMonth:  "12",
					ExpiryYear:   "25",
					SecurityCode: "123"},
			},
			wantErr: false,
		},
		{
			name: "Missing session ID",
			request: SessionDTO{
				MerchantID: "merchant123",
				SessionID:  "",
				Amount:     model.Amount{Currency: "EUR", Value: "10.50"},
				CustomerCardInfo: model.CustomerCardInfo{
					CardNumber:   "1234567890123456",
					ExpiryMonth:  "12",
					ExpiryYear:   "25",
					SecurityCode: "123"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateSessionRequest(tt.request); (err != nil) != tt.wantErr {
				t.Errorf("validateSessionRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateAmount(t *testing.T) {
	tests := []struct {
		name        string
		amount      model.Amount
		expectError bool
	}{
		{
			name:        "Valid amount",
			amount:      model.Amount{Value: "10.50", Currency: "EUR"},
			expectError: false,
		},
		{
			name:        "Missing value",
			amount:      model.Amount{Value: "", Currency: "EUR"},
			expectError: true,
		},
		{
			name:        "Invalid value",
			amount:      model.Amount{Value: "abc", Currency: "EUR"},
			expectError: true,
		},
		{
			name:        "Missing currency",
			amount:      model.Amount{Value: "10.50", Currency: ""},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAmount(tt.amount)
			if (err != nil) != tt.expectError {
				t.Errorf("TestValidateAmount(%s) error = %v, expectError %v", tt.name, err, tt.expectError)
				return
			}
		})
	}
}

func TestValidateCardInfos(t *testing.T) {
	tests := []struct {
		name         string
		customerCard model.CustomerCardInfo
		expectError  bool
	}{
		{
			name:         "Valid cardInfo",
			customerCard: model.CustomerCardInfo{CardNumber: "1234567890123456", ExpiryMonth: "12", ExpiryYear: "25", SecurityCode: "123"},
			expectError:  false,
		},
		{
			name:         "Missing cardNumber",
			customerCard: model.CustomerCardInfo{CardNumber: "", ExpiryMonth: "12", ExpiryYear: "25", SecurityCode: "123"},
			expectError:  true,
		},
		{
			name:         "Missing expiryMonth",
			customerCard: model.CustomerCardInfo{CardNumber: "1234567890123456", ExpiryMonth: "", ExpiryYear: "25", SecurityCode: "123"},
			expectError:  true,
		},
		{
			name:         "Missing expiryYear",
			customerCard: model.CustomerCardInfo{CardNumber: "1234567890123456", ExpiryMonth: "12", ExpiryYear: "", SecurityCode: "123"},
			expectError:  true,
		},
		{
			name:         "Missing securityCode",
			customerCard: model.CustomerCardInfo{CardNumber: "1234567890123456", ExpiryMonth: "12", ExpiryYear: "25", SecurityCode: ""},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCardInfos(tt.customerCard)
			if (err != nil) != tt.expectError {
				t.Errorf("TestValidateCardInfos(%s) error = %v, expectError %v", tt.name, err, tt.expectError)
				return
			}
		})
	}
}
