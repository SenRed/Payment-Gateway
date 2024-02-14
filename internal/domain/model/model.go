package model

import "time"

type TransactionResult struct {
	Status  TransactionStatus `json:"status"`
	Code    TransactionCode   `json:"code,omitempty"`
	Message string            `json:"message"`
	Date    *time.Time        `json:"date,omitempty"`
}

type TransactionStatus string

const (
	NotStartedStatus TransactionStatus = "NOT_STARTED"
	SuccessStatus    TransactionStatus = "SUCCESS"
	FailureStatus    TransactionStatus = "FAILURE"
)

type TransactionCode string

const (
	NotStartedCode            TransactionCode = "not-started"
	InvalidCVCCode            TransactionCode = "invalid-cvc"
	InsufficientFundsCode     TransactionCode = "insufficient-funds"
	SuccessfulTransactionCode TransactionCode = "successful-transaction"
)

type Transaction struct {
	// ID transaction ID
	ID uint
	// Status Status of the transaction
	Status TransactionStatus
	// Result of the request, could be "success" or "failure"
	Code TransactionCode
	// Message A human-readable message for more context
	Message string
	// Date Transaction date
	Date time.Time
}

type Amount struct {
	Currency string `json:"currency" example:"EUR"`
	Value    string `json:"value" example:"100"`
}
type CustomerCardInfo struct {
	CardNumber   string `json:"cardNumber" example:"4917484589897107"`
	ExpiryMonth  string `json:"expiryMonth" example:"02"`
	ExpiryYear   string `json:"expiryYear" example:"25"`
	SecurityCode string `json:"cvv" example:"123"`
}

type Session struct {
	SessionID        string
	MerchantID       string
	Amount           Amount
	CustomerCardInfo CustomerCardInfo
	Status           TransactionStatus
	Transactions     []Transaction
}

func (s *Session) MarkAsSuccessful() {
	s.Status = SuccessStatus
	successfulTransaction := Transaction{
		Status:  SuccessStatus,
		Code:    SuccessfulTransactionCode,
		Message: "Transaction executed successfully",
		Date:    time.Now(),
	}
	s.Transactions = append(s.Transactions, successfulTransaction)
}

func (s *Session) MarkAsFailed(transactionResult TransactionResult) {
	s.Status = FailureStatus
	failedTransaction := Transaction{
		Status:  FailureStatus,
		Code:    transactionResult.Code,
		Message: transactionResult.Message,
		Date:    time.Now(),
	}
	s.Transactions = append(s.Transactions, failedTransaction)
}

func NewSession(sessionID string, merchantID string, amount Amount, customerCardInfo CustomerCardInfo) Session {
	return Session{
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
				Date:    time.Now(),
			},
		},
	}
}
