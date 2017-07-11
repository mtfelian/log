package log

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/mtfelian/error"
	"net/http"
	"net/url"
	"time"
)

// Error writes error into log
func (logger *Logger) Error(c *gin.Context, httpCode int, errorCode uint, msg string, requestBody []byte) {
	var requestUrlString string
	var request *http.Request
	if c != nil {
		requestUrlString = c.Request.URL.String()
		unescapedRequestUrlString, _ = url.QueryUnescape(requestUrlString)
		request = c.Request
	}

	logger.Errorf("[%d][%d] %s [%s] %s", httpCode, errorCode,
		time.Now().Format("02.01.2006 15:04:05"), unescapedRequestUrlString, msg)
	if requestBody != nil {
		unescapedBody, _ := url.QueryUnescape(string(requestBody))
		logger.Errorf("Body: %s", unescapedBody)
	}
	if request != nil {
		unescapedRequest, _ := url.QueryUnescape(fmt.Sprintf("%v", request))
		logger.Errorf("Request: %s", unescapedRequest)
	}
}

// ReturnError writes error into log and returns an error
func (logger *Logger) ReturnError(c *gin.Context, httpCode int, errorCode uint, msg string, requestBody []byte) {
	logger.Error(c, httpCode, errorCode, msg, requestBody)
	if c != nil {
		c.JSON(httpCode, NewErrorf(errorCode, msg))
	}
}

// Success writes success into log
func (logger *Logger) Success(c *gin.Context, httpCode int, msg string, requestBody []byte) {
	var requestUrlString string
	var request *http.Request
	if c != nil {
		requestUrlString = c.Request.URL.String()
		unescapedRequestUrlString, _ = url.QueryUnescape(requestUrlString)
		request = c.Request
	}

	logger.Infof("[%d][%d] %s [%s] %s", httpCode, CodeSuccess,
		time.Now().Format("02.01.2006 15:04:05"), unescapedRequestUrlString, msg)
	if requestBody != nil {
		unescapedBody, _ := url.QueryUnescape(string(requestBody))
		logger.Infof("Body: %s", unescapedBody)
	}
	if request != nil {
		unescapedRequest, _ := url.QueryUnescape(fmt.Sprintf("%v", request))
		logger.Infof("Request: %s", unescapedRequest)
	}
}
