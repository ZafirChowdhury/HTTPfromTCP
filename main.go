package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {
	f, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("could not open %s: %s\n", inputFilePath, err)
	}

	fmt.Printf("Reading data from %s\n", inputFilePath)
	fmt.Println("=====================================")

	linesChan := getLinesChannel(f)

	// for loops exists when chanel is closed and there is no data to read
	for line := range linesChan {
		fmt.Println("read:", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	// spin up go rutines
	go func() {
		defer f.Close()
		defer close(lines)

		currentLineContents := ""
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)

			if err != nil {
				// current line has something but we are at EOF
				// send remaning data to the chanel
				if currentLineContents != "" {
					lines <- currentLineContents
				}

				if errors.Is(err, io.EOF) {
					break
				}

				// other error then EFO
				fmt.Printf("error: %s\n", err.Error())
				return
			}

			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			// send all the line one by one
			// skip the last part
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}
			// last element is ether "" or a partial line
			// add it to the currLine tracter then read the next 8 bytes
			currentLineContents += parts[len(parts)-1]
		}
	}()

	return lines
}
