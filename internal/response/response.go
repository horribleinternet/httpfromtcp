package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strings"
	"unicode"
)

type StatusCode int

const (
	HTTPOk                     StatusCode = 200
	HTTPBadRequest             StatusCode = 400
	HTTPInternalServerError    StatusCode = 500
	hTTPOkStr                             = "OK"
	hTTPBadRequestStr                     = "Bad Request"
	hTTPInternalServerErrorStr            = "Internal Server Error"
	headerLineEnd                         = "\r\n"
)

func formatStatusLine(statusCode StatusCode) []byte {
	return []byte(fmt.Sprintf("HTTP/1.1 %d %s%s", statusCode, hTTPStatuses[statusCode], headerLineEnd))
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write(formatStatusLine(statusCode))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := make(headers.Headers)
	header["Content-Length"] = fmt.Sprintf("%d", contentLen)
	header["Connection"] = "close"
	header["Content-Type"] = "text/plain"
	return header
}

func formatHeaders(headers headers.Headers) []byte {
	var headerStr string
	for k, v := range headers {
		name := strings.TrimRightFunc(k, func(r rune) bool { return unicode.IsSpace(r) })
		headerStr += fmt.Sprintf("%s: %s%s", name, v, headerLineEnd)
	}
	headerStr += headerLineEnd
	return []byte(headerStr)
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	_, err := w.Write(formatHeaders(headers))
	return err
}

var hTTPStatuses map[StatusCode]string

func init() {
	hTTPStatuses = make(map[StatusCode]string)
	hTTPStatuses[HTTPOk] = hTTPOkStr
	hTTPStatuses[HTTPBadRequest] = hTTPBadRequestStr
	hTTPStatuses[HTTPInternalServerError] = hTTPInternalServerErrorStr
}
