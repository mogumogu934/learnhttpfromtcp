package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	address := "localhost:42069"
	endPt, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatalf("unable to resolve udp address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, endPt)
	if err != nil {
		log.Fatalf("unable to dial UDP address %s", address)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("unable to read input: %v", err)
			continue
		}

		if line == "" {
			continue
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Printf("unable to write to connection: %v", err)
			continue
		}
	}
}
