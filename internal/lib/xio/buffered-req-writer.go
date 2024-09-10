package xio

import (
	"bufio"
	"io"

	"github.com/sirupsen/logrus"
)

func NewBufferedReqWriter(writer io.Writer) BufferedReqWriter {
	if w, ok := writer.(*bufio.Writer); ok {
		return bufReqWriter{w}
	} else {
		return bufReqWriter{bufio.NewWriter(writer)}
	}
}

// BufferedReqWriter defines a wrapper around Go's standard library bufio.Writer
// type that panics instead of returning error values.
//
// The methods have identical signatures to Go's standard library methods,
// meaning they declare that they could return an error, but that value will
// always be nil and may be ignored on all calls to methods on a
// BufferedReqWriter instance.
type BufferedReqWriter interface {
	io.ByteWriter
	io.Writer
	io.StringWriter

	WriteLineFeed()

	Flush()

	// this exists just so standard library types don't incidentally count as
	// valid BufferedReqWriter instances.
	nope()
}

type bufReqWriter struct{ w *bufio.Writer }

func (bufReqWriter) nope() {}

func (b bufReqWriter) WriteByte(y byte) (err error) {
	if err = b.w.WriteByte(y); err != nil {
		logrus.Fatalf("encountered error while flushing writing single byte: %s", err)
		panic(nil) // unreachable
	}

	return
}

func (b bufReqWriter) Write(p []byte) (n int, err error) {
	if n, err = b.w.Write(p); err != nil {
		logrus.Fatalf("encountered error while writing to buffer: %s", err)
		panic(nil) // unreachable
	}

	return
}

func (b bufReqWriter) WriteString(s string) (n int, err error) {
	if n, err = b.w.WriteString(s); err != nil {
		logrus.Fatalf("encountered error while writing to buffer: %s", err)
		panic(nil) // unreachable
	}

	return
}

func (b bufReqWriter) WriteLineFeed() {
	_ = b.WriteByte('\n')
}

func (b bufReqWriter) Flush() {
	if err := b.w.Flush(); err != nil {
		logrus.Fatalf("encountered error while flushing buffer: %s", err)
	}
}
