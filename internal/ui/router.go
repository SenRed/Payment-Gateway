package ui

import (
	"github.com/gin-gonic/gin"
	docs "github.com/payment-gateway/docs"
	"github.com/payment-gateway/internal/domain/api"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CreateRouter(paymentService api.PaymentService) *gin.Engine {
	r := gin.New()

	gin.ForceConsoleColor()
	r.Use(gin.Recovery())

	// health endpoints
	r.GET("/", healthCheck)
	r.GET("/health", healthCheck)

	// swagger endpoints
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// session endpoints
	r.POST("/v1/session/", func(c *gin.Context) {
		createNewSession(paymentService, c)
	})
	// payment endpoints
	r.POST("/v1/payment/:session-id", func(c *gin.Context) {
		startPayment(paymentService, c)
	})
	r.GET("/v1/payment/:session-id", func(c *gin.Context) {
		getPaymentDetails(paymentService, c)
	})
	return r
}
