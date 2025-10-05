package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

type ContentType int

const (
	lineEndStr                       = "\r\n"
	lineEndLen                       = len(lineEndStr)
	sep                              = ":"
	ContentTypeTextPlain ContentType = 0
	ContentTypeTextHTML  ContentType = 1
	ContentTypeVideo     ContentType = 2
	trailerFieldName                 = "Trailer"
)

var lineEndBytes = []byte(lineEndStr)

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) SetContextType(conType ContentType) error {
	v, ok := validContentTypes[conType]
	if !ok {
		return fmt.Errorf("invalid content type %d", conType)
	}
	h[contentTypeStr] = v
	return nil
}

func (h Headers) AddHeader(fieldName, fieldValue string) {
	h[fieldName] = fieldValue
}

func (h Headers) AddTrailers(trailerNames []string) {
	if len(trailerNames) == 0 {
		return
	}
	value := trailerNames[0]
	for i := 1; i < len(trailerNames); i++ {
		value += ", " + trailerNames[i]
	}
	h[trailerFieldName] = value
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
		key, value, found := strings.Cut(string(data[:index]), sep)
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

const (
	validRunes     = "!#$%&'*+-.^_`|~"
	contentTypeStr = "Content-Type"
	textPlainStr   = "text/plain"
	textHTMLStr    = "text/html"
	videoStr       = "video/mp4"
)

var validRuneSet map[rune]struct{}
var validContentTypes map[ContentType]string

func init() {
	validRuneSet = make(map[rune]struct{})
	for _, r := range validRunes {
		validRuneSet[r] = struct{}{}
	}
	validContentTypes = make(map[ContentType]string)
	validContentTypes[ContentTypeTextPlain] = textPlainStr
	validContentTypes[ContentTypeTextHTML] = textHTMLStr
	validContentTypes[ContentTypeVideo] = videoStr
}
