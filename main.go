package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"
const netProtocol = "tcp"

func main() {
	listener, err := net.Listen(netProtocol, port)
	if err != nil {
		log.Fatalf("error while reading from port %s: %s\n", port, err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error while accepting listener :%s\n", err.Error())
			return
		}

		log.Printf("Connection accepted at port: %s\n", port)
		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}

		log.Println("connection has beed closed")
	}
}

func getLinesChannel(rc io.ReadCloser) <-chan string {
	lines := make(chan string)

	// spin up go rutines
	go func() {
		defer rc.Close()
		defer close(lines)

		currentLineContents := ""
		for {
			buffer := make([]byte, 8)
			n, err := rc.Read(buffer)

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
