package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatalf("unable to open file: %v", err)
	}
	defer file.Close()

	linesChan := getLinesChannel(file)
	for l := range linesChan {
		fmt.Printf("read: %s\n", l)
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
