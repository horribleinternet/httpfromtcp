package server

import (
	"fmt"
	"httpfromtcp/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
}

var closed atomic.Bool

func Serve(port int) (*Server, error) {
	portStr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", portStr)
	if err != nil {
		return nil, err
	}
	closed.Store(false)
	server := Server{listener: listener}
	go server.listen()
	return &server, nil
}

func (s *Server) Close() error {
	closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if closed.Load() {
				break
			}
			log.Printf("error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	response.WriteStatusLine(conn, response.HTTPOk)
	response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	conn.Close()
}
