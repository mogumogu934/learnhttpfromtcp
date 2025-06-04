package main

import (
	"fmt"
	"log"
	"net"

	"github.com/mogumogu934/learnhttpfromtcp/internal/request"
)

func main() {
	port := ":42069"
	lsn, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("unable to set up listener: %v", err)
	}
	defer lsn.Close()

	for {
		conn, err := lsn.Accept()
		if err != nil {
			log.Fatalf("unable to accept connection to port %s: %v", port[1:], err)
		}
		fmt.Printf("Connection to port %s accepted\n", port[1:])

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Printf("unable to parse request: %v", err)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range req.Headers {
			fmt.Printf(" - %s: %s\n", k, v)
		}
	}
}
