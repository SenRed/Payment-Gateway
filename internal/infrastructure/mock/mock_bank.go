package mock

import (
	"encoding/json"
	"github.com/payment-gateway/internal/domain/model"
	"github.com/payment-gateway/internal/infrastructure/processor"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strings"
)

type AcquiringBankMock struct {
}

func (b AcquiringBankMock) Post(URL string, body processor.ProcessQuery) (*http.Response, error) {
	// I use the card number to identify the right use case to return
	switch body.CustomerCardInfo.CardNumber {
	case string(model.InvalidCVCCode):
		return &http.Response{
			StatusCode: 200,
			Body: createBody(processor.BankResponse{
				Result: string(model.InvalidCVCCode),
			}),
		}, nil
	case string(model.InsufficientFundsCode):
		return &http.Response{
			StatusCode: 200,
			Body: createBody(processor.BankResponse{
				Result: string(model.InsufficientFundsCode),
			}),
		}, nil
	case string(model.SuccessfulTransactionCode):
		return &http.Response{
			StatusCode: 200,
			Body: createBody(processor.BankResponse{
				Result: string(model.SuccessfulTransactionCode),
			}),
		}, nil
	default:
		return &http.Response{
			StatusCode: 200,
			Body: createBody(processor.BankResponse{
				Result: string(model.SuccessfulTransactionCode),
			}),
		}, nil
	}
}

func createBody(data processor.BankResponse) io.ReadCloser {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg("error when converting returned acquiring bank result")
	}

	return io.NopCloser(strings.NewReader(string(jsonData)))
}
func NewAcquiringBankHTTPClientMock() processor.AcquiringBankHTTPClient {
	return &AcquiringBankMock{}
}
