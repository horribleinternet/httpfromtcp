package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sock, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer sock.Close()
	read := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, err := read.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = sock.Write([]byte(line))
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
