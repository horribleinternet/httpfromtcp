package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Connection accepted")
		c := getLinesChannel(conn)
		for line := range c {
			fmt.Println(line)
		}
		fmt.Println("Connection closed")
	}
}

func getLinesChannel(f net.Conn) <-chan string {
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
