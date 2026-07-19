package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

var rn = []byte("\r\n")

func NewHeaders() Headers {
	return make(Headers)
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

func (h Headers) Parse(data []byte) (int, bool, error) {
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
	h[name] = value

	return idx + len(rn), false, nil
}
