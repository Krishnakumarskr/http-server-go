package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")

	if err != nil {
		log.Fatal("Error resolving udp address")
	}

	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		log.Fatal("Error dialing UDP")
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf(">")
		input, err := reader.ReadString('\n')

		if err != nil {
			log.Fatalf("\nError reading input %v", err)
		}

		res, err := conn.Write([]byte(input))
		if err != nil {
			fmt.Printf("\nError in writing to udp %v", err)
		}

		fmt.Printf("Res: %d\n", res)
	}

}
