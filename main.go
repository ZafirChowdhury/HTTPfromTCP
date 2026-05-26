package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer file.Close()

	buffer := make([]byte, 8)
	var line string
	for {

		n, err := file.Read(buffer)

		if err == io.EOF {
			fmt.Printf("read: %s\n", line)
			return
		}

		if err != nil {
			log.Println(err.Error())
			return
		}

		currentIta := string(buffer[:n])
		i := strings.IndexByte(currentIta, '\n')
		if i == -1 {
			line += currentIta
			continue
		}

		line += currentIta[:i]
		fmt.Printf("read: %s\n", line)
		line = currentIta[i+1:]
	}
}
