package pica

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

type Output struct {
	Debug            bool
	DefaultLineCount int
	writer           io.Writer
}

func NewOutput(debug bool, writer io.Writer) *Output {
	return &Output{
		Debug:            debug,
		DefaultLineCount: 60,
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

func (o *Output) RepeatLine(e string, count int) string {
	return "\n" + strings.Repeat(e, count)
}

func (o *Output) EchoStartRequest(request ApiRequest) {
	fmt.Println(o.L("="))
	color.Green("\nStarting request [%s %s %s]\n\n", request.Method, request.Url, request.Name)
}

func (o *Output) ErrorRequest(err error) {
	color.Red("do http request error %s", err.Error())
}

func (o *Output) EchoRequstIng(method string, body []byte) {
	fmt.Printf("%s ...", method)
	color.Blue("\n%s\n\n", body)
}

func (o *Output) Finished(count int, names string) {
	fmt.Println(o.L("="))
	color.Green("\nFinished. [%d] api requests, [%s] passed", count, names)
	fmt.Println(o.L("="))
}

func (o *Output) CopyRight() {
	fmt.Println(o.L("="))
	color.Yellow(DefaultCopyright)
}

var (
	DefaultOutput = NewOutput(true, os.Stdout)
)
