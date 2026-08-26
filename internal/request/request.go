package request

import (
	"bytes"
	"fmt"
	"io"
)

type Request struct {
	RequestLine RequestLine
	state       parseState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var SEPERATOR = []byte("\r\n")
var ERROR_MALFORMED_REQUEST_LINE = fmt.Errorf("malformed request line")

type parseState string

var StateInit parseState = "init"
var StateDone parseState = "done"

func newRequest() *Request {
	return &Request{state: StateInit}
}

func parseRequestLine(s []byte) (*RequestLine, int, error) {
	idx := bytes.Index(s, SEPERATOR)

	if idx == -1 {
		return nil, 0, nil
	}

	startLine := s[:idx]
	read := idx + len(SEPERATOR)

	line := bytes.Split(startLine, []byte(" "))
	if len(line) != 3 {
		return nil, 0, ERROR_MALFORMED_REQUEST_LINE
	}

	httpParts := bytes.Split(line[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[1]) != "1.1" {
		return nil, 0, fmt.Errorf("http version parse error")
	}

	return &RequestLine{
		HttpVersion:   string(httpParts[1]),
		RequestTarget: string(line[1]),
		Method:        string(line[0]),
	}, read, nil

}

func (r *Request) parse(data []byte) (int, error) {

	read := 0
outer:
	for {
		switch r.state {
		case StateInit:
			{
				rl, n, err := parseRequestLine(data[read:])

				if err != nil {
					return read, err
				}

				if n == 0 {
					break outer
				}

				r.RequestLine = *rl
				read += n
				r.state = StateDone
			}
		case StateDone:
			{
				break outer
			}
		}
	}

	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])

		if err != nil {
			return nil, err
		}

		bufLen += n
		readN, err := request.parse(buf[:bufLen])

		if err != nil {
			return nil, err
		}
		copy(buf, buf[readN:bufLen])
		bufLen -= readN

	}

	return request, nil
}
