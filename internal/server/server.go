package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	handler  Handler
}

type HandlerError struct {
	Status  response.StatusCode
	Message string
}

func (h HandlerError) isError() bool {
	return h.Status != response.HTTPOk
}

var closed atomic.Bool

type Handler func(w io.Writer, req *request.Request) *HandlerError

func Serve(port int, h Handler) (*Server, error) {
	portStr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", portStr)
	if err != nil {
		return nil, err
	}
	closed.Store(false)
	server := Server{listener: listener, handler: h}
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
	defer conn.Close() //no net.Conn gets out alive
	req, err := request.RequestFromReader(conn)
	writer := response.NewWriter(conn)
	if err != nil {
		hErr := &HandlerError{Status: response.HTTPBadRequest, Message: err.Error()}
		hErr.WriteError(writer)
		return
	}
	var buffer bytes.Buffer
	handErr := s.handler(&buffer, req)
	if handErr != nil {
		handErr.WriteError(writer)
		return
	}
	writer.WriteStatusLine(response.HTTPOk)
	header := response.GetDefaultHeaders(buffer.Len())
	writer.WriteHeaders(header)
	conn.Write(buffer.Bytes())
}

func (e *HandlerError) WriteError(w *response.Writer) {
	w.WriteStatusLine(e.Status)
	w.WriteHeaders(response.GetDefaultHeaders(len(e.Message)))
	w.WriteBody([]byte(e.Message))
}
