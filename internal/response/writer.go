package response

import (
	"fmt"
	"io"

	"github.com/mogumogu934/learnhttpfromtcp/internal/headers"
)

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateHeaders
	writerStateBody
	CRLF = "\r\n"
)

type Writer struct {
	writerState writerState
	writer      io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writerState: writerStateStatusLine,
		writer:      w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writerStateStatusLine {
		return fmt.Errorf("unable to write status line in state %d", w.writerState)
	}
	defer func() {
		w.writerState = writerStateHeaders
	}()

	_, err := w.writer.Write(getStatusLine(statusCode))
	return err
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.writerState != writerStateHeaders {
		return fmt.Errorf("unable to write headers in state %d", w.writerState)
	}
	defer func() {
		w.writerState = writerStateBody
	}()

	for k, v := range h {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s", k, v) + CRLF))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte(CRLF))
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("unable to write body in state %d", w.writerState)
	}
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("unable to write body in state %d", w.writerState)
	}

	totalBytes := 0

	lenHex := fmt.Sprintf("%x", len(p))
	n, err := w.writer.Write([]byte(lenHex))
	if err != nil {
		return totalBytes, err
	}
	totalBytes += n

	n, err = w.writer.Write([]byte(CRLF))
	if err != nil {
		return totalBytes, err
	}
	totalBytes += n

	n, err = w.writer.Write(p)
	if err != nil {
		return totalBytes, err
	}
	totalBytes += n

	n, err = w.writer.Write([]byte(CRLF))
	if err != nil {
		return totalBytes, err
	}
	totalBytes += n

	return totalBytes, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("unable to write body in state %d", w.writerState)
	}

	n, err := w.WriteBody([]byte("0" + CRLF + CRLF))
	if err != nil {
		return 0, err
	}

	return n, nil
}
