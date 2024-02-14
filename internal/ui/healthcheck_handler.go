package ui

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type httpStatus struct {
	Value string `json:"status"`
}

// healthCheck http handler
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, httpStatus{Value: "OK"})
}
