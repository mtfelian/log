package log

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (logger *Logger) LogError(c *gin.Context, responseCode int, msg string, requestBody []byte) {
	var requestUrlString string
	var request *http.Request
	if c != nil {
		requestUrlString = c.Request.URL.String()
		request = c.Request
	}

	logger.Errorf("[%d] %s [%s] %s", responseCode, time.Now().Format("02.01.2006 15:04:05"), requestUrlString, msg)
	if requestBody != nil {
		logger.Errorf("Body: %s", string(requestBody))
	}
	if request != nil {
		logger.Errorf("Request: %v", request)
	}
}

func (logger *Logger) LogAndReturnError(c *gin.Context, responseCode int, msg string, requestBody []byte) {
	logger.LogError(c, responseCode, msg, requestBody)
	if c != nil {
		c.JSON(responseCode, gin.H{"error": msg})
	}
}

// LogSuc
func (logger *Logger) LogSuccess(c *gin.Context, responseCode int, msg string, requestBody []byte) {
	var requestUrlString string
	var request *http.Request
	if c != nil {
		requestUrlString = c.Request.URL.String()
		request = c.Request
	}

	logger.Infof("[%d] %s [%s] %s", responseCode, time.Now().Format("02.01.2006 15:04:05"), requestUrlString, msg)
	if requestBody != nil {
		logger.Infof("Body: %s", string(requestBody))
	}
	if request != nil {
		logger.Infof("Request: %v", request)
	}
}
