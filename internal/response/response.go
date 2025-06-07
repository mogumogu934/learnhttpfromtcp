package response

import (
	"fmt"

	"github.com/mogumogu934/learnhttpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusCodeOK                  StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

func getStatusLine(statusCode StatusCode) []byte {
	reason := ""
	switch statusCode {
	case StatusCodeOK:
		reason = "OK"
	case StatusCodeBadRequest:
		reason = "Bad Request"
	case StatusCodeInternalServerError:
		reason = "Internal Server Error"
	}
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reason)
	return []byte(statusLine)
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}
