package pica

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
	_ "github.com/jeremaihloo/pica/statik"
	"github.com/rakyll/statik/fs"
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

func (o *Output) EchoStartRequest(request *ApiRequest, runner *APIRunner) error {
	fmt.Println(o.L("="))
	fmt.Println()
	color.Green("%s %s %s", request.Method, request.Url, request.Name)
	targetUrl, err := getTargetUrl(request, runner)
	if err != nil {
		return err
	}
	color.Blue("\nRequest %s\n\n", targetUrl)
	return nil
}

func (o *Output) ErrorRequest(err error) {
	color.Red("do http request error %s", err.Error())
}

func (o *Output) EchoRequstIng(method string, body []byte) {
	fmt.Printf("%s ...", method)
	color.Yellow("\n%s\n\n", body)
}

func (o *Output) Finished(count int, names string) {
	fmt.Println(o.L("="))
	color.Green("\nFinished. [%d] api requests, [%s] passed", count, names)
	fmt.Println(o.L("="))
}

func (o *Output) CopyRight() {
	fmt.Println(o.L("="))
	sfs, err := fs.New()
	if err != nil {
		panic(err)
	}
	file, err := sfs.Open("/copyright.txt")
	if err != nil {
		panic(err)
	}
	DefaultCopyright, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	color.Yellow(string(DefaultCopyright))
}

func (o *Output) Headers(headers http.Header) {
	for key, _ := range headers {
		fmt.Printf("%s: %s\n", key, headers.Get(key))
	}
	fmt.Println()
}

func (o *Output) RequestBody(req *http.Request, runner *APIRunner) error {
	if req.Method != "GET" && req.Method != "DELETE" {
		body := runner.vm.Lookup(strings.ToLower(req.Method))
		data, err := prettyjson.Marshal(body)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	}
	return nil
}

func (o *Output) ResponseBody(res *http.Response) {

}

func (o *Output) Echo(s string) {
	fmt.Printf(s)
}

func (o *Output) Echoln(s string) {
	fmt.Println(s)
}

func (o *Output) Json(obj interface{}) {
	switch obj := obj.(type) {
	case map[string]interface{}:
		o.Json(&obj)
		break
	case *map[string]interface{}:
		data, err := prettyjson.Marshal(obj)
		if err != nil {
			panic(err)
		}
		fmt.Print(string(data))
		break
	case []byte:
		var newObj map[string]interface{}
		err := json.Unmarshal(obj, &newObj)
		if err != nil {
			panic(err)
		}
		o.Json(newObj)
		break
	default:
		data, err := prettyjson.Marshal(obj)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
		break

	}
}

var (
	DefaultOutput = NewOutput(true, os.Stdout)
)
