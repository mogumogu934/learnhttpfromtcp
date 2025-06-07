package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/mogumogu934/learnhttpfromtcp/internal/headers"
)

const (
	StatusCodeOK                  = 200
	StatusCodeBadRequest          = 400
	StatusCodeInternalServerError = 500
	CRLF                          = "\r\n"
)

func WriteStatusLine(w io.Writer, statusCode int) error {
	base := fmt.Sprintf("HTTP/1.1 %v ", statusCode)
	switch statusCode {
	case StatusCodeOK:
		_, err := w.Write([]byte(base + "OK" + CRLF))
		if err != nil {
			return fmt.Errorf("unable to write status line: %v", err)
		}
	case StatusCodeBadRequest:
		_, err := w.Write([]byte(base + "Bad Request" + CRLF))
		if err != nil {
			return fmt.Errorf("unable to write status line: %v", err)
		}
	case StatusCodeInternalServerError:
		_, err := w.Write([]byte(base + "Internal Server Error" + CRLF))
		if err != nil {
			return fmt.Errorf("unable to write status line: %v", err)
		}
	default:
		_, err := w.Write([]byte(base + CRLF))
		if err != nil {
			return fmt.Errorf("unable to write status line: %v", err)
		}
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	return headers.Headers{
		"Content-Length": strconv.Itoa(contentLen),
		"Connection":     "close",
		"Content-Type":   "text/plain",
	}
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		h := fmt.Sprintf("%s: %s", k, v) + CRLF
		_, err := w.Write([]byte(h))
		if err != nil {
			return fmt.Errorf("unable to write header (%s): %v", h, err)
		}
	}

	_, err := w.Write([]byte(CRLF))
	if err != nil {
		return fmt.Errorf("unable to writer CRLF to signify end of headers: %v", err)
	}
	return nil
}
