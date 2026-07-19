package headers

import (
	"bytes"
	"fmt"
	"strings"
)

func isToken(str []byte) bool {
	for _, ch := range str {
		found := false
		if ch >= 'A' && ch <= 'Z' || ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' {
			found = true
		}
		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			found = true
		}

		if !found {
			return false
		}
	}
	return true
	// Uppercase letters: A-Z
	// Lowercase letters: a-z
	// Digits: 0-9
	// Special characters: !, #, $, %, &, ', *, +, -, ., ^, _, `, |, ~
}

type Headers struct {
	headers map[string]string
}

var rn = []byte("\r\n")

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func parseHeader(fieldLane []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLane, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed header")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if len(bytes.TrimSpace(name)) != len(name) {
		return "", "", fmt.Errorf("malformed field name")
	}

	return string(name), string(value), nil
}

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func (h *Headers) Set(name, value string) {
	h.headers[strings.ToLower(name)] = value
}

func (h *Headers) Parse(data []byte) (int, bool, error) {
	idx := bytes.Index(data, rn)
	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		return len(rn), true, nil
	}

	name, value, err := parseHeader(data[:idx])
	if err != nil {
		return 0, false, err
	}

	if !isToken([]byte(name)) {
		return 0, false, fmt.Errorf("err: header name")
	}

	h.Set(name, value)

	return idx + len(rn), false, nil
}
