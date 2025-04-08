package headers

import (
	"bytes"
	"errors"
	"strings"
	"unicode"
)

const CRLF = "\r\n"

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	endIndex := bytes.Index(data, []byte(CRLF))
	if endIndex == -1 {
		return 0, false, nil
	}

	if endIndex == 0 {
		return 2, true, nil
	}

	headerLine := string(data[:endIndex])
	parts := strings.SplitN(strings.TrimSpace(headerLine), ":", 2)
	if len(parts) < 2 {
		return endIndex + 2, false, errors.New("invalid header line: colon not detected")
	}

	if strings.ContainsAny(parts[0], " ") {
		return 0, false, errors.New("invalid header line: contains spaces between field name and colon")
	}

	if !validateFieldName(parts[0]) {
		return 0, false, errors.New("invalid header line: field name contains invalid characters")
	}

	fieldName := strings.ToLower(parts[0])
	fieldValue := strings.TrimSpace(parts[1])
	fieldValue = strings.TrimSuffix(fieldValue, ";")

	if existingValue, exists := h[fieldName]; exists {
		h[fieldName] = existingValue + ", " + fieldValue
	} else {
		h[fieldName] = fieldValue
	}

	return endIndex + 2, false, nil
}

func validateFieldName(fieldName string) bool {
	validChars := "!#$%&'*+-.^_`|~"

	for _, c := range fieldName {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) && !strings.ContainsRune(validChars, c) {
			return false
		}
	}

	return true
}
