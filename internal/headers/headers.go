package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

const (
	lineEndStr = "\r\n"
	lineEndLen = len(lineEndStr)
	sep        = ":"
	sepLen     = len(sep)
)

var lineEndBytes = []byte(lineEndStr)

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	index := bytes.Index(data, lineEndBytes)
	fmt.Printf("index: %d, %s\n", index, string(data[:index]))
	if index < 0 {
		return 0, false, nil
	} else if index == 0 {
		return lineEndLen, true, nil
	} else {
		key, value, found := strings.Cut(string(data[:index]), ":")
		fmt.Println(key, value)
		if !found {
			return 0, false, fmt.Errorf("no key:value pair found")
		}
		keyTrim := strings.TrimRightFunc(key, func(r rune) bool { return unicode.IsSpace(r) })
		if keyTrim != key {
			return 0, false, fmt.Errorf("whitespace after key")
		}
		key = strings.TrimLeftFunc(key, func(r rune) bool { return unicode.IsSpace(r) })
		value = strings.TrimFunc(value, func(r rune) bool { return unicode.IsSpace(r) })
		fmt.Println(key, value)
		h[key] = value
		return index + lineEndLen, false, nil
	}
}
