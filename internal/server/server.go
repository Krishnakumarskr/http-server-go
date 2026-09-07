package server

import (
	"fmt"
	"http-server/internal/request"
	"http-server/internal/response"
	"io"
	"net"
)

type Server struct {
	closed  bool
	handler Handler
}

type Handler func(w *response.Writer, req *request.Request)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	headers := response.GetDefaultHeaders(0)
	responseWriter := response.NewWriter(conn)
	r, err := request.RequestFromReader(conn)

	if err != nil {
		responseWriter.WriteStatusLine(response.BadRequest)
		responseWriter.WriteHeaders(*headers)
		return
	}
	s.handler(responseWriter, r)
}

func runServer(s *Server, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if s.closed {
			return
		}
		if err != nil {
			return
		}

		go runConnection(s, conn)
	}
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}
	server := &Server{closed: false, handler: handler}

	go runServer(server, listener)

	return server, nil
}

func (s *Server) Close() {
	s.closed = true
}

// func (s *Server) listen()

// func (s *Server) handle(conn net.Conn) {}
