package processor

import (
	"net/http"
)

type AcquiringBankHTTPClient interface {
	Post(URL string, body ProcessQuery) (*http.Response, error)
}
