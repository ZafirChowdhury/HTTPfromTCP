package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()

	buffer := make([]byte, 8)
	for {

		n, err := file.Read(buffer)

		if err == io.EOF {
			return
		}

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("read: %s\n", string(buffer[:n]))
	}
}
