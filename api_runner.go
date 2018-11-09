package pica

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/jeremaihloo/funny/langs"
	"io/ioutil"
	"net/http"
	"strings"
)

type ApiRunner struct {
	Filename string
	ApiNames []string
	Delay    int

	content []byte
	vm      *langs.Interpreter
	parser  *langs.Parser
	client  *http.Client
	output  *Output

	ApiItems  []*ApiItem
	Block     langs.Block
	InitLines langs.Block
}

func NewApiRunnerFromFile(filename string, apiNames []string, delay int) *ApiRunner {
	return &ApiRunner{
		Filename: filename,
		ApiNames: apiNames,
		Delay:    delay,

		client: http.DefaultClient,
		vm:     langs.NewInterpreterWithScope(DefaultInitScope),
		output: DefaultOutput,
	}
}

func NewApiRunnerFromContent(content []byte) *ApiRunner {
	return &ApiRunner{
		Filename: "",
		ApiNames: []string{},
		Delay:    0,
		content:  content,
		client:   http.DefaultClient,
		vm:       langs.NewInterpreterWithScope(DefaultInitScope),
		output:   DefaultOutput,
	}
}

func (runner *ApiRunner) Run() error {
	runner.vm.RegisterFunction("address", Address)
	runner.vm.RegisterFunction("email", Email)
	runner.vm.RegisterFunction("phone", Phone)
	runner.vm.RegisterFunction("words", Words)
	runner.vm.RegisterFunction("name", FullName)
	runner.vm.RegisterFunction()
	err := runner.Parse()
	if err != nil {
		return err
	}
	// parse api file to ApiRequest
	err = runner.ParseApiItems()
	if err != nil {
		return err
	}

	runner.RunInitLines()

	for i := 0; i < len(runner.ApiItems); i++ {
		item := runner.ApiItems[i]
		if len(runner.ApiNames) == 0 {
			err = runner.RunSingle(item)
			if err != nil {
				return err
			}
		} else {
			for index := 0; index < len(runner.ApiNames); index++ {
				name := runner.ApiNames[index]
				if item.Request.Name == name {
					err = runner.RunSingle(item)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (runner *ApiRunner) Parse() error {
	if runner.Filename != "" {
		buffer, err := langs.CombinedCode(runner.Filename)
		if err != nil {
			return fmt.Errorf("parse error %v", err.Error())
		}
		runner.content = []byte(buffer)
	}
	runner.parser = langs.NewParser(runner.content)
	runner.Block = runner.parser.Parse()
	return nil
}

func (runner *ApiRunner) RunInitLines() {
	for _, line := range runner.InitLines {
		runner.vm.EvalStatement(line)
	}
}

func (runner *ApiRunner) RunSingle(item *ApiItem) error {
	// assign vars

	runner.vm.Assign("url", item.Request.Url)
	// Eval init scope statements
	for _, line := range item.Request.lines {
		runner.vm.EvalStatement(line)
	}

	headers := runner.vm.Lookup("headers").(map[string]langs.Value)

	// send ApiRequest by http client
	res, err := runner.DoApiRequest(item.Request)
	if err != nil {
		runner.output.ErrorRequest(err)
		return fmt.Errorf("do http request error %s", err.Error())
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	item.Response.Headers = res.Header
	item.Response.Status = res.StatusCode
	item.Response.Body = buf.Bytes()

	// collect http response to ApiRequest
	headers = HttpHeaders2VmMap(item.Response.Headers)
	runner.vm.Assign("header", headers)
	runner.vm.Assign("status", item.Response.Status)
	runner.vm.Assign("body", item.Response.Body)

	contentType := item.Response.Headers.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		jResults := make(map[string]langs.Value)
		jun := make(map[string]interface{})
		err := json.Unmarshal(item.Response.Body, &jun)
		if err != nil {
			color.Red(fmt.Sprintf("json binding %s %s", err.Error(), item.Response.Body))
		}
		for k, v := range jun {
			jResults[k] = langs.Value(v)
		}
		runner.vm.Assign("json", jResults)

		runner.output.Json(&jResults)
	} else {
		resData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			color.Red(err.Error())
		}
		fmt.Print(string(resData))
	}

	// Eval item response statement
	for _, line := range item.Response.lines {
		runner.vm.EvalStatement(line)
	}

	return nil
}

func (runner *ApiRunner) DoApiRequest(req *ApiRequest) (*http.Response, error) {

	runner.output.EchoStartRequest(req, runner)

	httpReq, err := CreateHttpRequest(req, runner)
	if err != nil {
		return nil, err
	}

	runner.output.Headers(httpReq.Header)
	runner.output.RequestBody(httpReq, runner)
	res, err := runner.client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	runner.output.Echo("\nResponse ")
	if res.StatusCode == 200 {
		color.Green("Status: %d\n\n", res.StatusCode)
	} else {
		color.Red("Status: %d\n\n", res.StatusCode)
	}
	runner.output.Headers(res.Header)
	return res, nil
}

func (runner *ApiRunner) ParseApiItems() error {
	headers := VmMap2HttpHeaders(DefaultHeaders)
	inited := false
	index := 0
	asserting := false
	for index < len(runner.Block) {
		line := runner.Block[index]
		switch line := line.(type) {
		case *langs.Comment:
			text := strings.Trim(line.Value, " ")
			texts := strings.Split(text, " ")
			if len(texts) < 2 {
				break
			}
			methods := []string{"GET", "POST", "DELETE", "PUT", "PATCH"}
			flag := false
			for _, item := range methods {
				if strings.ToUpper(item) == strings.ToUpper(texts[0]) {
					flag = true
				}
			}
			if flag {
				inited = true
				asserting = false
				req := ApiRequest{
					Method:  texts[0],
					Url:     texts[1],
					Headers: headers,
				}
				if len(texts) > 2 {
					req.Name = texts[2]
				}
				if len(texts) > 3 {
					req.Description = texts[3]
				}
				apiItem := &ApiItem{
					Request:  &req,
					Response: &ApiResponse{},
				}
				runner.ApiItems = append(runner.ApiItems, apiItem)
			}
		case *langs.FunctionCall:
			if line.Name == "assert" {
				asserting = true
			}
			if asserting {
				item := runner.ApiItems[len(runner.ApiItems)-1]
				item.Response.lines = append(item.Response.lines, line)
				break
			}
		default:
			if inited {
				if asserting {
					item := runner.ApiItems[len(runner.ApiItems)-1]
					item.Response.lines = append(item.Response.lines, line)
				} else {
					item := runner.ApiItems[len(runner.ApiItems)-1]
					item.Request.lines = append(item.Request.lines, line)
				}
			} else {
				runner.InitLines = append(runner.InitLines, line)
			}
		}
		index++
	}
	return nil
}
