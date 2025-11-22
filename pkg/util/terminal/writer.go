package terminal

import (
	"bytes"
	"fmt"
	"io"

	"github.com/lunixbochs/vtclean"
)

type Writer struct {
	Writer io.Writer
	buf    *bytes.Buffer
}

func NewWriter(writer io.Writer) *Writer {
	var buf bytes.Buffer
	return &Writer{
		Writer: writer,
		buf:    &buf,
	}
}

func (t *Writer) Write(p []byte) (int, error) {
	n, err := t.buf.Write(p)
	lines := bytes.Split(t.buf.Bytes(), []byte{'\n'})
	t.buf = bytes.NewBuffer(lines[len(lines)-1])
	for idx, line := range lines[:len(lines)-1] {
		cl := vtclean.Clean(string(line), false)
		if idx != len(lines)-1 {
			cl = fmt.Sprintf("%s\n", cl)
		}
		_, err2 := t.Writer.Write([]byte(cl))
		if err2 != nil {
			return n, err2
		}
	}
	return n, err
}

func (t *Writer) Close() error {
	line := t.buf.Bytes()
	cl := vtclean.Clean(string(line), false)
	_, err := t.Writer.Write([]byte(cl))
	return err
}
