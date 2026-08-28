package main

import (
	"bytes"
	"fmt"
	"http-server/internal/request"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer close(lines)
		defer f.Close()
		str := ""

		for {
			data := make([]byte, 8)
			n, err := f.Read(data)

			if err != nil {
				break
			}

			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				data = data[i+1:]
				lines <- str
				str = ""
			}
			str += string(data)
		}

		if len(str) != 0 {
			lines <- str
		}
	}()

	return lines
}

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
	}
}
