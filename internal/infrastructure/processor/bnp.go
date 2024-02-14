package processor

import (
	"encoding/json"
	"github.com/payment-gateway/internal/domain/model"
	"github.com/payment-gateway/internal/domain/port"
	"net/http"
)

// BNPClient BNP Paribas bank client
type BNPClient struct {
	HttpClient AcquiringBankHTTPClient
}

type ProcessQuery struct {
	MerchantID       string                 `json:"merchantID"`
	Amount           model.Amount           `json:"amount"`
	CustomerCardInfo model.CustomerCardInfo `json:"customerCardInfo"`
}

type BankResponse struct {
	Result string `json:"result"`
}

func (b *BNPClient) Process(session model.Session) (*model.TransactionResult, *model.DomainError) {
	// Call the bank API
	processBody := ProcessQuery{
		MerchantID:       session.MerchantID,
		Amount:           session.Amount,
		CustomerCardInfo: session.CustomerCardInfo,
	}

	response, err := b.HttpClient.Post("/transaction", processBody)
	if err != nil {
		return nil, &model.DomainError{Category: model.TechnicalError, Message: "Error when calling BNP bank client"}
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, &model.DomainError{Category: model.TechnicalError, RootCause: err, Message: "Error when calling BNP bank client"}
	}

	var bankResponse BankResponse
	if err := json.NewDecoder(response.Body).Decode(&bankResponse); err != nil {
		return nil, &model.DomainError{Category: model.TechnicalError, RootCause: err, Message: "Error when calling BNP bank client"}
	}

	return mapToTransactionResult(bankResponse)
}

func mapToTransactionResult(bankResponse BankResponse) (*model.TransactionResult, *model.DomainError) {
	switch bankResponse.Result {
	case string(model.InsufficientFundsCode):
		return &model.TransactionResult{Status: model.FailureStatus, Code: model.InsufficientFundsCode, Message: "Insufficient funds"}, nil
	case string(model.InvalidCVCCode):
		return &model.TransactionResult{Status: model.FailureStatus, Code: model.InvalidCVCCode, Message: "Invalid CVC"}, nil
	case string(model.SuccessfulTransactionCode):
		return &model.TransactionResult{Status: model.SuccessStatus, Code: model.SuccessfulTransactionCode, Message: "The transaction is successful"}, nil
	default:
		return nil, &model.DomainError{Category: model.TechnicalError, Message: "Unknown bank result code " + bankResponse.Result}
	}
}

func NewBankProcessor(httpClient AcquiringBankHTTPClient) port.IPaymentProcessor {
	return &BNPClient{HttpClient: httpClient}
}
