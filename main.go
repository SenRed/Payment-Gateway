package main

import (
	"github.com/payment-gateway/internal/bootstrap"
)

// @title payment-gateway
// @version 1.0
// @description This API aims to provide payment gateway endpoints documentation.

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes: http
// @host localhost:8070
// @BasePath /

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	bootstrap.Init().Run()
}
