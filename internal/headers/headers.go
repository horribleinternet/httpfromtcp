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

func (h Headers) Get(fieldName string) (string, error) {
	value, ok := h[strings.ToLower(fieldName)]
	if !ok {
		return "", fmt.Errorf("field name %s not present", fieldName)
	}
	return value, nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	index := bytes.Index(data, lineEndBytes)
	//fmt.Printf("index: %d, %s\n", index, string(data[:index]))
	if index < 0 {
		return 0, false, nil
	} else if index == 0 {
		return lineEndLen, true, nil
	} else {
		key, value, found := strings.Cut(string(data[:index]), ":")
		//fmt.Println(key, value)
		if !found {
			return 0, false, fmt.Errorf("no field-name:field-value pair found")
		}
		keyTrim := strings.TrimRightFunc(key, func(r rune) bool { return unicode.IsSpace(r) })
		if keyTrim != key {
			return 0, false, fmt.Errorf("whitespace after field name '%s'", key)
		}
		key = strings.TrimLeftFunc(key, func(r rune) bool { return unicode.IsSpace(r) })
		if !validFieldName(key) {
			return 0, false, fmt.Errorf("invalid character in field name '%s'", key)
		}
		key = strings.ToLower(key)
		value = strings.TrimFunc(value, func(r rune) bool { return unicode.IsSpace(r) })
		existing, there := h[key]
		if there {
			h[key] = existing + ", " + value
		} else {
			h[key] = value
		}
		return index + lineEndLen, false, nil
	}
}

func validFieldName(fieldName string) bool {
	if len(fieldName) == 0 {
		return false
	}
	for _, r := range fieldName {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			continue
		}
		if _, valid := validRuneSet[r]; !valid {
			return false
		}
	}
	return true
}

const validRunes = "!#$%&'*+-.^_`|~"

var validRuneSet map[rune]struct{}

func init() {
	validRuneSet = make(map[rune]struct{})
	for _, r := range validRunes {
		validRuneSet[r] = struct{}{}
	}
}
