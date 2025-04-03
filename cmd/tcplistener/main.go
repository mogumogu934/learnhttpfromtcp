package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
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

		linesChan := getLinesChannel(conn)
		for l := range linesChan {
			fmt.Printf("%s\n", l)
		}
		fmt.Printf("Connection to port %s closed\n", port[1:])
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer close(lines)
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			lines <- line
		}

		if err := scanner.Err(); err != nil {
			log.Printf("error reading file: %v", err)
		}
	}()

	return lines
}
