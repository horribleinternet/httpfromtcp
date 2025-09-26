package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var buf [8]byte
	var line string
	for {
		n, err := file.Read(buf[:])
		if err != nil {
			break
		}
		parts := strings.Split(string(buf[:n]), "\n")
		line = line + parts[0]
		for i := 1; i < len(parts); i++ {
			fmt.Printf("read: %s\n", line)
			line = parts[i]
		}
	}
	if len(line) > 0 {
		fmt.Printf("read: %s\n", line)
	}
}
