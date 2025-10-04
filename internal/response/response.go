package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strings"
	"unicode"
)

type StatusCode int

type writerState int

const (
	HTTPOk                     StatusCode  = 200
	HTTPBadRequest             StatusCode  = 400
	HTTPInternalServerError    StatusCode  = 500
	hTTPOkStr                              = "OK"
	hTTPBadRequestStr                      = "Bad Request"
	hTTPInternalServerErrorStr             = "Internal Server Error"
	headerLineEnd                          = "\r\n"
	writerStateStatusLine      writerState = 0
	writerStateHeaders         writerState = 1
	writerStateBody            writerState = 2
	writerStateDone            writerState = 3
)

type Writer struct {
	out   io.Writer
	state writerState
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{out: writer, state: writerStateStatusLine}
}

func formatStatusLine(statusCode StatusCode) []byte {
	return []byte(fmt.Sprintf("HTTP/1.1 %d %s%s", statusCode, hTTPStatuses[statusCode], headerLineEnd))
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != writerStateStatusLine {
		return fmt.Errorf("calling WriteStatusLine more than once")
	}
	_, err := w.out.Write(formatStatusLine(statusCode))
	if err == nil {
		w.state = writerStateHeaders
	}
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

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state < writerStateHeaders {
		return fmt.Errorf("calling WriteHeaders before writing status line")
	} else if w.state > writerStateHeaders {
		return fmt.Errorf("calling WriteHeaders more than once")
	}
	_, err := w.out.Write(formatHeaders(headers))
	if err == nil {
		w.state = writerStateBody
	}
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state < writerStateBody {
		return 0, fmt.Errorf("calling WriteBody before writing preceeding sections")
	} else if w.state > writerStateBody {
		return 0, fmt.Errorf("calling WriteBody more than once")
	}
	n, err := w.out.Write(p)
	if err == nil {
		w.state = writerStateDone
	}
	return n, err
}

var hTTPStatuses map[StatusCode]string

func init() {
	hTTPStatuses = make(map[StatusCode]string)
	hTTPStatuses[HTTPOk] = hTTPOkStr
	hTTPStatuses[HTTPBadRequest] = hTTPBadRequestStr
	hTTPStatuses[HTTPInternalServerError] = hTTPInternalServerErrorStr
}
