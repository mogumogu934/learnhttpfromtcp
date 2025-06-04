package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mogumogu934/learnhttpfromtcp/internal/headers"
)

const (
	bufferSize                 = 8
	requestStateInitialized    = 0
	requestStateDone           = 1
	requestStateParsingHeaders = 2
	CRLF                       = "\r\n"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0

	r := Request{
		Headers: map[string]string{},
		state:   requestStateInitialized,
	}

	for r.state != requestStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if r.state != requestStateDone {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", r.state, numBytesRead)
				}
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead

		numBytesParsed, err := r.parse(buf[:readToIndex])
		if err != nil {
			return nil, fmt.Errorf("error parsing request from reader: %v", err)
		}

		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return &r, nil
}

func parseRequestLine(data []byte) (parsedLine *RequestLine, numBytesParsed int, err error) {
	endIndex := bytes.Index(data, []byte(CRLF))
	if endIndex == -1 {
		return nil, 0, nil
	}

	reqLine := string(data[:endIndex])
	parts := strings.Split(reqLine, " ")

	if len(parts) != 3 {
		return nil, 0, errors.New("invalid request line")
	}

	if parts[0] != "GET" &&
		parts[0] != "POST" &&
		parts[0] != "PUT" &&
		parts[0] != "PATCH" &&
		parts[0] != "DELETE" {
		return nil, 0, errors.New("invalid request method")
	}

	if parts[1][0] != '/' {
		return nil, 0, errors.New("invalid request target")
	}

	if !strings.HasPrefix(parts[2], "HTTP/") {
		return nil, 0, errors.New("invalid request line")
	}

	version := strings.Split(parts[2], "/")[1]

	if version != "1.1" {
		return nil, 0, errors.New("unsupported version of http")
	}

	return &RequestLine{
		HttpVersion:   version,
		RequestTarget: parts[1],
		Method:        parts[0],
	}, endIndex + 2, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != requestStateDone {
		numBytesParsed, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += numBytesParsed
		if numBytesParsed == 0 {
			break
		}
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		reqLine, numBytesParsed, err := parseRequestLine(data)
		if err != nil {
			return 0, fmt.Errorf("error parsing request line: %v", err)
		}
		if numBytesParsed == 0 {
			// need more data
			return 0, nil
		}
		r.RequestLine = *reqLine
		r.state = requestStateParsingHeaders
		return numBytesParsed, nil
	case requestStateParsingHeaders:
		numBytesParsed, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = requestStateDone
		}
		return numBytesParsed, nil
	case requestStateDone:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("unknown state")
	}
}
