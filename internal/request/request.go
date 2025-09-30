package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type parseState int

const (
	done        parseState = 0
	initialized parseState = 1
	lineEnd                = "\r\n"
	lineEndLen             = len(lineEnd)
	bufferSize             = 8
)

type Request struct {
	RequestLine RequestLine
	state       parseState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	if r.state == done {
		return 0, fmt.Errorf("cannot parse done request")
	}
	out, ReqLine, err := parseRequestLine(data)
	if err != nil {
		return 0, err
	}
	if out > 0 {
		r.RequestLine = ReqLine
		r.state = done
	}
	return out, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	var req Request
	req.state = initialized
	readToIndex := 0
	for req.state != done {
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
		parsed, err := req.parse(buf)
		if req.state != done {
			copy(buf, buf[parsed:])
			readToIndex -= parsed
		}
	}
	return &req, nil
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
