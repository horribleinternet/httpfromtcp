package main

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, testHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func testHandler(w *response.Writer, req *request.Request) *server.HandlerError {

	var buffer bytes.Buffer
	var err error
	var n int
	status := response.HTTPOk
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		status = response.HTTPBadRequest
		n, err = buffer.Write([]byte(badRequest))
	case "/myproblem":
		status = response.HTTPInternalServerError
		n, err = buffer.Write([]byte(internalError))
	default:
		n, err = buffer.Write([]byte(okRequest))
	}
	if err != nil {
		return &server.HandlerError{Status: response.HTTPInternalServerError, Message: err.Error()}
	}
	err = w.WriteStatusLine(status)
	if err != nil {
		fmt.Printf("Unable to write status line for target %s: %s\n", req.RequestLine.RequestTarget, err.Error())
		return &server.HandlerError{Status: response.HTTPInternalServerError, Message: err.Error()}
	}
	header := response.GetDefaultHeaders(n)
	header.SetContextType(headers.ContentTypeTextHTML)
	err = w.WriteHeaders(header)
	if err != nil {
		fmt.Printf("Unable to write header for target %s: %s\n", req.RequestLine.RequestTarget, err.Error())
		return &server.HandlerError{Status: response.HTTPInternalServerError, Message: err.Error()}
	}
	_, err = w.WriteBody(buffer.Bytes())
	if err != nil {
		fmt.Printf("Unable to write body for target %s: %s\n", req.RequestLine.RequestTarget, err.Error())
		return &server.HandlerError{Status: response.HTTPInternalServerError, Message: err.Error()}
	}
	return nil
}

const (
	badRequest    = "<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>"
	internalError = "<html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body></html>"
	okRequest     = "<html><head><title>200 OK</title></head><body><h1>Success!</h1><p>Your request was an absolute banger.</p></body></html>"
)
