package pica

import (
	"io"
	"fmt"
	"os"
	"strings"
)

type Output struct {
	Debug            bool
	DefaultLineCount int
	writer           io.Writer
}

func NewOutput(debug bool, writer io.Writer) *Output {
	return &Output{
		Debug:            debug,
		DefaultLineCount: 10,
		writer:           writer,
	}
}

func (o *Output) W(args ...interface{}) *Output {
	if len(args) > 1 {
		o.writer.Write([]byte(fmt.Sprintf(args[0].(string)+"\n", args[1:]...)))
	} else {
		o.writer.Write([]byte(args[0].(string) + "\n"))
	}
	return o
}

func (o *Output) Wln(args ...interface{}) *Output {
	if len(args) == 0 {
		o.writer.Write([]byte("\n"))
	} else {
		args[0] = args[0].(string) + ""
		return o.W(args)
	}
	return o
}

func (o *Output) L(e string) string {
	return "\n" + strings.Repeat(e, o.DefaultLineCount)
}

func (o *Output) Line(e string, count int) string {
	return "\n" + strings.Repeat(e, count)
}

var (
	DefaultOutput = NewOutput(true, os.Stdout)
)
