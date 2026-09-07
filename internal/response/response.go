package response

import (
	"fmt"
	"http-server/internal/headers"
	"io"
)

type Response struct {
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

type Writer struct {
	writer io.Writer
}

type StatusCode uint16

const Success StatusCode = 200
const BadRequest StatusCode = 400
const InternalError StatusCode = 500

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()

	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusMsg := []byte{}

	switch statusCode {
	case Success:
		{
			statusMsg = []byte("HTTP/1.1 200 OK\r\n")
		}
	case BadRequest:
		{
			statusMsg = []byte("HTTP/1.1 400 Bad Request\r\n")
		}
	case InternalError:
		{
			statusMsg = []byte("HTTP/1.1 500 Internal Server Error\r\n")
		}
	default:
		{
			return fmt.Errorf("Unsupported StatusCode!")
		}
	}

	_, err := w.writer.Write(statusMsg)

	return err
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	b := []byte{}

	h.ForEach(func(k, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", k, v)
	})
	b = fmt.Appendf(b, "\r\n")
	_, err := w.writer.Write(b)
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.writer.Write(p)

	return n, err
}
