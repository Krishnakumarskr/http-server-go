package main

import (
	"fmt"
	"http-server/internal/request"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("Error in listening to port")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error in accepting connection")
		}

		req, err := request.RequestFromReader(conn)

		if err != nil {
			log.Fatalf("Error in reading request %+v", err)
		}

		fmt.Printf("Request Line:\n- Method: %s\n- Target: %s\n- Version: %s", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		fmt.Printf("\nHeaders:\n")
		req.Headers.ForEach(func(name, val string) {
			fmt.Printf("-%s: %s\n", name, val)
		})
		fmt.Printf("Body\n%s", req.Body)
	}
}
