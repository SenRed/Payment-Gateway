package ui

import (
	"github.com/gin-gonic/gin"
	"github.com/payment-gateway/internal/domain/api"
	"github.com/payment-gateway/internal/domain/model"
	"github.com/rs/zerolog/log"
	"net/http"
)

type ResultError struct {
	ErrorReason string `json:"reason"`
}

type SessionDTO struct {
	SessionID        string                 `json:"sessionId" example:"d6f49736-3922-4520-8c74-b2fee3e0d113" comment:"A unique ID, we recommend using a UUID"`
	MerchantID       string                 `json:"merchantID" example:"amazonID"`
	Amount           model.Amount           `json:"amount"`
	CustomerCardInfo model.CustomerCardInfo `json:"customerCardInfo"`
}

// CreateNewSession
// @ID create-session
// @Summary Create a session
// @Description Create a new session for a transaction by sending an idempotency key, card information and amount/currency.
// @Param 	contentBody  body 	SessionDTO 	true 	"Request body"
// @Success 200 {object} nil  "OK - Session created. No data returned in the response body."
// @Failure 400 {object} ResultError "Bad Request - Request parameters are invalid."
// @Failure 500 {object} ResultError "Internal Server Error - An error occurred while processing the request on the server side."
// @Router /v1/session/ [post]
func createNewSession(paymentService api.PaymentService, c *gin.Context) {
	// Step 0 - Read the JSON sent by the client
	var sessionDTO SessionDTO
	if err := c.ShouldBind(&sessionDTO); err != nil {
		log.Debug().Err(err).Interface("sessionDTO", sessionDTO).Msg("Error when parsing sessionDTO")
		c.JSON(http.StatusBadRequest, ResultError{ErrorReason: "Error when parsing sessionDTO"})
		return
	}
	// Step 1 - Validate the inputs
	if err := validateSessionRequest(sessionDTO); err != nil {
		log.Debug().Interface("sessionDTO", sessionDTO).Msg(err.Message)
		c.JSON(http.StatusBadRequest, ResultError{ErrorReason: err.Message})
		return
	}
	// Step 2 - Create the session
	err := paymentService.CreateNewSession(sessionDTO.SessionID, sessionDTO.MerchantID, sessionDTO.Amount, sessionDTO.CustomerCardInfo)
	// Step 3 - According to the result, return the right status code and content
	handleResult(sessionDTO.SessionID, nil, err, c)
}

// StartPayment
// @ID start-transaction
// @Summary Start a transaction
// @Description Start a transaction by sending the session ID.
// @Param session-id path string true "Session ID"
// @Success 200 {object} model.TransactionResult "OK - The request was executed as expected. The statusMessage field contains information about the transaction status."
// @Failure 400 {object} ResultError "Bad Request - Request parameters are invalid."
// @Failure 500 {object} ResultError "Internal Server Error - An error occurred while processing the request on the server side."
// @Produce json
// @Router /v1/payment/{session-id} [post]
func startPayment(paymentService api.PaymentService, c *gin.Context) {
	// Step 1 - Read session ID
	sessionId := c.Param("session-id")
	if err := validateSessionId(sessionId); err != nil {
		// pkg.LogErrorAndReturn(c, fmt.Sprintf("Invalid query parameters: %s", err))
		return
	}

	// Step 2 - Process the payment
	res, err := paymentService.Start(sessionId)

	// Step 3 - According to the result, return the right status code and content
	handleResult(sessionId, res, err, c)
}

// getPaymentDetails retrieves payment details for a session.
// @Summary Get payment details
// @Description Retrieve payment details for a session by providing the session ID.
// @Param session-id path string true "Session ID"
// @Success 200 {object} []model.TransactionResult "All transaction of the session."
// @Failure 400 {object} ResultError "Bad Request - Request parameters are invalid."
// @Failure 500 {object} ResultError "Internal Server Error - An error occurred while processing the request on the server side."
// @Produce json
// @Router /v1/payment/{session-id} [get]
func getPaymentDetails(paymentService api.PaymentService, c *gin.Context) {
	// Step 1 -  Validate session ID
	sessionId := c.Param("session-id")
	if err := validateSessionId(sessionId); err != nil {
		// pkg.LogErrorAndReturn(c, fmt.Sprintf("Invalid query parameters: %s", err))
		return
	}

	// Step 2 - Get session
	res, err := paymentService.GetPaymentDetails(sessionId)

	// Step 3 - According to the result, return the right status code and content
	handleResult(sessionId, res, err, c)
}

func handleResult(sessionId string, res any, err *model.DomainError, c *gin.Context) {
	switch {
	case err == nil:
		// Case no error
		if res != nil {
			c.JSON(http.StatusOK, res)
		} else {
			c.Status(http.StatusOK)
		}
	case err.Type == model.SessionNotFound || err.Type == model.DuplicatedSessionId:
		// Handle session not found or Duplicated session ID
		log.Debug().Interface("sessionId", sessionId).Msg(err.Message)
		c.JSON(http.StatusBadRequest, ResultError{ErrorReason: err.Message})
	default:
		// Other case
		log.Error().Err(err.RootCause).Interface("sessionID", sessionId).Msg(err.Message)
		c.JSON(http.StatusInternalServerError, ResultError{ErrorReason: "An error occurred when processing the creation request"})
	}
}
