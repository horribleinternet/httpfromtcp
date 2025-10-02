package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"net"
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
		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Request line:")
			fmt.Printf("- Method: %s\n", req.RequestLine.Method)
			fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
			fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
			fmt.Println("Headers:")
			for k, v := range req.Headers {
				fmt.Printf("- %s: %s\n", k, v)
			}
		}
	}
}

/*
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
*/
