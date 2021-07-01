package pica

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/jerloo/funny"
)

// APIRunner the runner
type APIRunner struct {
	Filename string
	APINames []string
	Delay    int

	content []byte
	vm      *funny.Interpreter
	parser  *funny.Parser
	client  *http.Client
	output  *Output

	APIItems  []*ApiItem
	Block     funny.Block
	InitLines funny.Block
}

// NewAPIRunnerFromFile create a runner from a pica file
func NewAPIRunnerFromFile(filename string, apiNames []string, delay int) *APIRunner {
	return &APIRunner{
		Filename: filename,
		APINames: apiNames,
		Delay:    delay,

		client: http.DefaultClient,
		vm:     funny.NewInterpreterWithScope(DefaultInitScope),
		output: DefaultOutput,
	}
}

// NewAPIRunnerFromContent create a runner from a pica content
func NewAPIRunnerFromContent(content []byte) *APIRunner {
	return &APIRunner{
		Filename: "",
		APINames: []string{},
		Delay:    0,
		content:  content,
		client:   http.DefaultClient,
		vm:       funny.NewInterpreterWithScope(DefaultInitScope),
		output:   DefaultOutput,
	}
}

// Run run the task
func (runner *APIRunner) Run() error {
	runner.vm.RegisterFunction("address", Address)
	runner.vm.RegisterFunction("email", Email)
	runner.vm.RegisterFunction("phone", Phone)
	runner.vm.RegisterFunction("words", Words)
	runner.vm.RegisterFunction("name", FullName)
	err := runner.Parse()
	if err != nil {
		return err
	}
	// parse api file to ApiRequest
	err = runner.ParseAPIItems()
	if err != nil {
		return err
	}

	runner.RunInitLines()

	for i := 0; i < len(runner.APIItems); i++ {
		item := runner.APIItems[i]
		if len(runner.APINames) == 0 {
			err = runner.RunSingle(item)
			if err != nil {
				return err
			}
		} else {
			for index := 0; index < len(runner.APINames); index++ {
				name := runner.APINames[index]
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

// Parse parse pica file
func (runner *APIRunner) Parse() error {
	if runner.Filename != "" {
		buffer, err := funny.CombinedCode("", runner.Filename)
		if err != nil {
			return fmt.Errorf("parse error %v", err.Error())
		}
		runner.content = []byte(buffer)
	}
	runner.parser = funny.NewParser(runner.content)
	runner.Block = runner.parser.Parse()
	return nil
}

// RunInitLines run the code of initialization
func (runner *APIRunner) RunInitLines() {
	for _, line := range runner.InitLines {
		runner.vm.EvalStatement(line)
	}
}

// RunSingle run the single api item
func (runner *APIRunner) RunSingle(item *ApiItem) error {
	// assign vars

	runner.vm.Assign("url", item.Request.Url)
	// Eval init scope statements
	for _, line := range item.Request.lines {
		runner.vm.EvalStatement(line)
	}

	headers := runner.vm.Lookup("headers").(map[string]funny.Value)

	// send ApiRequest by http client
	res, err := runner.DoAPIRequest(item.Request)
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
		jResults := make(map[string]funny.Value)
		jun := make(map[string]interface{})
		err := json.Unmarshal(item.Response.Body, &jun)
		if err != nil {
			color.Red(fmt.Sprintf("json binding %s %s", err.Error(), item.Response.Body))
		}
		for k, v := range jun {
			jResults[k] = funny.Value(v)
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

// DoAPIRequest run the api request
func (runner *APIRunner) DoAPIRequest(req *ApiRequest) (*http.Response, error) {

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

// ParseAPIItems parse ap items from pica code
func (runner *APIRunner) ParseAPIItems() error {
	headers := VmMap2HttpHeaders(DefaultHeaders)
	inited := false
	index := 0
	asserting := false
	for index < len(runner.Block) {
		line := runner.Block[index]
		switch line := line.(type) {
		case *funny.Comment:
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
				runner.APIItems = append(runner.APIItems, apiItem)
			}
		case *funny.FunctionCall:
			if line.Name == "assert" {
				asserting = true
			}
			if asserting {
				item := runner.APIItems[len(runner.APIItems)-1]
				item.Response.lines = append(item.Response.lines, line)
				break
			}
		default:
			if inited {
				if asserting {
					item := runner.APIItems[len(runner.APIItems)-1]
					item.Response.lines = append(item.Response.lines, line)
				} else {
					item := runner.APIItems[len(runner.APIItems)-1]
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
