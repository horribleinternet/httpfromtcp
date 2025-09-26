package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c := getLinesChannel(file)
	for line := range c {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)
	go func() {
		defer f.Close()
		defer close(out)
		var buf [8]byte
		var line string
		for {
			n, err := f.Read(buf[:])
			if err != nil {
				if len(line) > 0 {
					out <- line
				}
				if !errors.Is(err, io.EOF) {
					fmt.Println(err.Error())
				}
				break
			}
			parts := strings.Split(string(buf[:n]), "\n")
			line = line + parts[0]
			for i := 1; i < len(parts); i++ {
				out <- line
				line = parts[i]
			}
		}
	}()
	return out
}
