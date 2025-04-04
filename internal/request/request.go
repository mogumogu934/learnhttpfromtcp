package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, fmt.Errorf("unable to read request: %w", err)
	}

	reqLine := strings.Split(string(req), "\r\n")[0]
	parsedLine, err := parseRequestLine(reqLine)
	if err != nil {
		return &Request{}, fmt.Errorf("unable to parse request line: %w", err)
	}

	return parsedLine, nil
}

func parseRequestLine(reqLine string) (*Request, error) {
	parts := strings.Split(reqLine, " ")

	if len(parts) != 3 {
		return &Request{}, errors.New("invalid request line")
	}

	if parts[0] != "GET" &&
		parts[0] != "POST" &&
		parts[0] != "PUT" &&
		parts[0] != "PATCH" &&
		parts[0] != "DELETE" {
		return &Request{}, errors.New("invalid request line")
	}

	if parts[1][0] != '/' {
		return &Request{}, errors.New("invalid request target")
	}

	if parts[2] != "HTTP/1.1" {
		return &Request{}, errors.New("invalid request line or http version")
	}

	return &Request{
		RequestLine{
			HttpVersion:   strings.Split(parts[2], "/")[1],
			RequestTarget: parts[1],
			Method:        parts[0],
		},
	}, nil
}
