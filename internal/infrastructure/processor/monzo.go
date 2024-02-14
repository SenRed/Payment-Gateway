package processor

import (
	"github.com/payment-gateway/internal/domain/model"
)

// MonzoClient Monzo client
type MonzoClient struct {
	// Url or bank + endpoint or API
	HttpClient AcquiringBankHTTPClient
	URL        string
}

func (b *MonzoClient) Process(_ model.Session) (*model.TransactionResult, *model.DomainError) {

	//TODO Not implemented, just for understanding of purpose, check BNP.go for more complete implementation

	// Call the bank API
	// ...
	// return result
	panic("implement me")
}
