package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

const bufferSize = 8
const stateInitialized = 0
const stateDone = 1
const CRLF = "\r\n"

type Request struct {
	RequestLine RequestLine
	State       int
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
		RequestLine: RequestLine{
			HttpVersion:   "",
			RequestTarget: "",
			Method:        "",
		},
		State: stateInitialized,
	}

	for r.State != stateDone {
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			if readToIndex > 0 {
				numBytesParsed, parseErr := r.parse(buf[:readToIndex])
				if parseErr != nil {
					return nil, fmt.Errorf("unable to parse buffer: %w", err)
				}

				if numBytesParsed > 0 {
					copy(buf, buf[numBytesParsed:readToIndex])
					readToIndex -= numBytesParsed
				}
			}

			if r.State == stateInitialized {
				return nil, errors.New("incomplete request: reached EOF before parsing complete")
			}

			r.State = stateDone
			break
		}

		if err != nil {
			return nil, fmt.Errorf("unable to read from io reader: %w", err)
		}

		readToIndex += numBytesRead

		numBytesParsed, err := r.parse(buf[:readToIndex])
		if err != nil {
			return nil, fmt.Errorf("unable to parse buffer: %w", err)
		}

		if numBytesParsed > 0 {
			copy(buf, buf[numBytesParsed:readToIndex])
			readToIndex -= numBytesParsed
		}
	}

	return &r, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.State == stateDone {
		return 0, errors.New("error: trying to read data in a done state")
	}

	if r.State != stateInitialized {
		return 0, errors.New("error: unknown state")
	}

	reqLine, numBytesParsed, err := parseRequestLine(data)
	if err != nil {
		return 0, fmt.Errorf("unable to parse request line: %w", err)
	}

	if numBytesParsed == 0 {
		return 0, nil
	}

	if reqLine != nil {
		r.RequestLine = *reqLine
		r.State = stateDone
	}

	return numBytesParsed, nil
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
