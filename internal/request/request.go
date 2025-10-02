package request

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type parseState int

const (
	requestStateDone           parseState = 0
	requestStateInitialized    parseState = 1
	requestStateParsingHeaders parseState = 2
	requestStateParsingBody    parseState = 3
	lineEnd                               = "\r\n"
	lineEndLen                            = len(lineEnd)
	contentLengthFieldName     string     = "content-length"
	bufferSize                            = 8
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       parseState
}

func newRequest() *Request {
	var request Request
	request.state = requestStateInitialized
	request.Headers = headers.NewHeaders()
	return &request
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	if r.state == requestStateDone {
		return 0, fmt.Errorf("cannot parse done request")
	}
	parsedBytes := 0
	for r.state != requestStateDone {
		n, err := r.parseNext(data[parsedBytes:])
		parsedBytes += n
		if err != nil {
			return parsedBytes, fmt.Errorf("'%v' parsing from byte %d", err, parsedBytes)
		}
		if n == 0 {
			break
		}
	}
	return parsedBytes, nil
}

func (r *Request) parseNext(data []byte) (int, error) {
	switch r.state {
	case requestStateDone:
		return 0, fmt.Errorf("cannot parse done request")
	case requestStateInitialized:
		n, line, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n > 0 {
			r.RequestLine = line
			r.state = requestStateParsingHeaders
		}
		return n, nil
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = requestStateParsingBody
		}
		return n, nil
	case requestStateParsingBody:
		lengthStr, err := r.Headers.Get(contentLengthFieldName)
		if err != nil {
			r.state = requestStateDone
			return 0, nil
		}
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return 0, err
		}
		n := len(data)
		r.Body = append(r.Body, data...)
		if len(r.Body) > length {
			return n, fmt.Errorf("content length header value was %d but body is %d bytes long", length, len(r.Body))
		} else if len(r.Body) == length {
			r.state = requestStateDone
		}
		return n, nil
	default:
		return 0, fmt.Errorf("invalid request state %d", r.state)
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	req := newRequest()
	readToIndex := 0
	for req.state != requestStateDone {
		if readToIndex == len(buf) {
			newbuf := make([]byte, 2*len(buf))
			copy(newbuf, buf)
			buf = newbuf
		}
		read, err := reader.Read(buf[readToIndex:])
		if err != nil {
			return nil, err
		}
		readToIndex += read
		parsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		if req.state != requestStateDone {
			copied := copy(buf, buf[parsed:])
			clear(buf[copied:])
			readToIndex -= parsed
		}
	}
	return req, nil
}

func parseRequestLine(header []byte) (int, RequestLine, error) {
	dataStr := string(header)
	index := strings.Index(dataStr, lineEnd)
	if index < 0 {
		return 0, RequestLine{}, nil
	}
	req, _, found := strings.Cut(string(header), lineEnd)
	if !found {
		return 0, RequestLine{}, fmt.Errorf("request line not found")
	}
	var parsed RequestLine
	method, remain, found := strings.Cut(req, " ")
	parsed.Method = method
	if !found || !allUpper(parsed.Method) {
		return 0, RequestLine{}, fmt.Errorf("invalid method %s", parsed.Method)
	}
	middle, verStr, found := cutLast(remain, " ")
	if !found {
		return 0, RequestLine{}, fmt.Errorf("no version string")
	}
	verParts := strings.Split(verStr, "/")
	if len(verParts) != 2 || verParts[0] != "HTTP" {
		return 0, RequestLine{}, fmt.Errorf("invalid version string %s", verStr)
	}
	if verParts[1] != "1.1" {
		return 0, RequestLine{}, fmt.Errorf("unsupported version %s", verParts[1])
	}
	if strings.Contains(middle, " ") {
		return 0, RequestLine{}, fmt.Errorf("invalid request target %s", middle)
	}
	parsed.HttpVersion = verParts[1]
	parsed.RequestTarget = middle

	return index + lineEndLen, parsed, nil
}

func allUpper(str string) bool {
	for _, r := range str {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func cutLast(str string, sep string) (before, after string, found bool) {
	if i := strings.LastIndex(str, sep); i >= 0 {
		return str[:i], str[i+len(sep):], true
	}
	return str, "", false
}
