package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func (h *Headers) Set(name string, value string) {
	name = strings.ToLower(name)

	if val, ok := h.headers[name]; ok {
		h.headers[name] = strings.Join([]string{val, value}, ",")
		// h.headers[name] = fmt.Sprintf("%s,%s", val, value)
	} else {
		h.headers[name] = value
	}
}

func (h *Headers) Replace(name string, value string) {
	name = strings.ToLower(name)
	h.headers[name] = value
}

func (h *Headers) Delete(name string) {
	name = strings.ToLower(name)
	delete(h.headers, name)
}

func (h *Headers) ForEach(cb func(name string, val string)) {
	for key, val := range h.headers {
		cb(key, val)
	}
}

func isToken(b []byte) bool {
	for _, ch := range b {
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
}

var crlf = []byte("\r\n")

// Host:ows value ows
func parseHeader(filedLine []byte) (string, string, error) {
	parts := bytes.SplitN(filedLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed field line")
	}
	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasPrefix(name, []byte(" ")) {
		return "", "", fmt.Errorf("malformed field name")
	}

	return string(name), string(value), nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	read := 0
	done = false

	for {
		n := bytes.Index(data[read:], crlf)
		if n == -1 {
			break
		}

		if n == 0 {
			done = true
			read += len(crlf)
			break
		}

		name, value, err := parseHeader(data[read : read+n])
		if err != nil {
			return 0, false, err
		}

		if !isToken([]byte(name)) {
			return 0, false, fmt.Errorf("malformed header name")
		}
		read += n + len(crlf)
		h.Set(name, value)
	}

	return read, done, nil
}
