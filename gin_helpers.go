package log

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/mtfelian/error"
)

// Error writes error into log
func (logger *Logger) Error(c *gin.Context, httpCode int, errorCode uint, msg string, requestBody []byte) {
	var requestUrlString string
	var request *http.Request
	if c != nil {
		requestUrlString = c.Request.URL.String()
		request = c.Request
	}

	logger.Errorf("[%d][%d] %s [%s] %s", httpCode, errorCode,
		time.Now().Format("02.01.2006 15:04:05"), requestUrlString, msg)
	if requestBody != nil {
		logger.Errorf("Body: %s", string(requestBody))
	}
	if request != nil {
		logger.Errorf("Request: %v", request)
	}
}

// ReturnError writes error into log and returns an error
func (logger *Logger) ReturnError(c *gin.Context, httpCode int, errorCode uint, msg string, requestBody []byte) {
	logger.Error(c, httpCode, errorCode, msg, requestBody)
	if c != nil {
		c.JSON(httpCode, error.StandardError{errorCode, msg})
	}
}

// Success writes success into log
func (logger *Logger) Success(c *gin.Context, httpCode int, msg string, requestBody []byte) {
	var requestUrlString string
	var request *http.Request
	if c != nil {
		requestUrlString = c.Request.URL.String()
		request = c.Request
	}

	logger.Infof("[%d][%d] %s [%s] %s", httpCode, error.CodeSuccess,
		time.Now().Format("02.01.2006 15:04:05"), requestUrlString, msg)
	if requestBody != nil {
		logger.Infof("Body: %s", string(requestBody))
	}
	if request != nil {
		logger.Infof("Request: %v", request)
	}
}
