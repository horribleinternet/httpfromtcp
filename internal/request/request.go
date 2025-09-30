package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	header, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	requestline, _, err := parseRequestLine(header)
	if err != nil {
		return nil, err
	}
	fullRequest := Request{RequestLine: requestline}
	return &fullRequest, nil
}

func parseRequestLine(header []byte) (RequestLine, []byte, error) {
	req, rest, found := strings.Cut(string(header), "\r\n")
	if !found {
		return RequestLine{}, header, fmt.Errorf("request line not found")
	}
	var parsed RequestLine
	method, remain, found := strings.Cut(req, " ")
	parsed.Method = method
	if !found || !allUpper(parsed.Method) {
		return RequestLine{}, header, fmt.Errorf("invalid method %s", parsed.Method)
	}
	middle, verStr, found := cutLast(remain, " ")
	if !found {
		return RequestLine{}, header, fmt.Errorf("no version string")
	}
	verParts := strings.Split(verStr, "/")
	if len(verParts) != 2 || verParts[0] != "HTTP" {
		return RequestLine{}, header, fmt.Errorf("invalid version string %s", verStr)
	}
	if verParts[1] != "1.1" {
		return RequestLine{}, header, fmt.Errorf("unsupported version %s", verParts[1])
	}
	if strings.Contains(middle, " ") {
		return RequestLine{}, header, fmt.Errorf("invalid request target %s", middle)
	}
	parsed.HttpVersion = verParts[1]
	parsed.RequestTarget = middle

	return parsed, []byte(rest), nil
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
