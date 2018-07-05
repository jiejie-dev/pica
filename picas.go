package pica

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fixate/go-qs"
	"github.com/jeremaihloo/funny/langs"
)

type ApiRequest struct {
	Headers     http.Header
	Method      string
	Url         string
	Query       Query
	Name        string
	Description string
	Body        []byte
	lines       langs.Block
}

type Query map[string]interface{}

func NewQuery(m map[string]interface{}) Query {
	query := Query{}
	for k, v := range m {
		query[k] = v
	}
	return query
}

func ParseQuery(queryString string) (Query, error) {
	r, err := qs.Unmarshal(queryString)
	if err != nil {
		return nil, err
	}
	return NewQuery(r), nil
}

func (query Query) String() (string, error) {
	return qs.Marshal(query)
}

type ApiResponse struct {
	Headers http.Header
	Body    []byte
	Status  int
	lines   langs.Block

	saveLines langs.Block
}

type ApiItem struct {
	Request  ApiRequest
	Response ApiResponse
}

type ApiContext struct {
	Name        string
	Description string
	Author      string
	Version     string
	BaseUrl     string
	Headers     *http.Header
	InitVars    langs.Scope
	InitLines   langs.Block
	ApiItems    []*ApiItem

	Pica         *Pica
	VersionNotes *VersionNote
}

type Pica struct {
	FileName        string
	Output          string
	Debug           bool
	DocTempalteFile string
	Delay           int

	vm     *langs.Interpreter
	parser *langs.Parser
	Block  langs.Block
	client *HttpClient

	Ctx *ApiContext

	output *Output
}

func NewPica(filename string, delay int, output, template string) *Pica {
	return &Pica{
		FileName:        filename,
		Output:          output,
		Delay:           delay,
		DocTempalteFile: template,

		vm:     langs.NewInterpreterWithScope(langs.Scope{}),
		output: DefaultOutput,
	}
}

func (p *Pica) Run() error {
	err := p.Parse()
	if err != nil {
		return err
	}
	err = p.ParseApiContext()
	if err != nil {
		return err
	}
	err = p.RunApiContext()
	if err != nil {
		return err
	}
	if p.Output != "" {
		err = p.Document()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pica) Document() error {
	p.runInitPartOfContext(p.Ctx)

	data, err := ioutil.ReadFile(p.DocTempalteFile)
	if err != nil {
		data = []byte(DEFAULT_DOC_TEMPLATE)
	}
	generator := NewGenerator(p.Ctx, string(data))
	results, err := generator.Get()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", results)
	if _, err := os.Stat(p.Output); err == nil {
		err := ioutil.WriteFile(p.Output, results, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pica) ParseApiContext() error {
	headers := VmMap2HttpHeaders(DefaultHeaders)
	ctx := &ApiContext{
		Headers: &headers,
		Pica:    p,
	}
	inited := false
	index := 0
	asserting := false
	for index < len(p.Block) {
		line := p.Block[index]
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
				req := ApiRequest{
					Method:  texts[0],
					Url:     texts[1],
					Headers: http.Header{},
				}
				if len(texts) > 2 {
					req.Name = texts[2]
				}
				if len(texts) > 3 {
					req.Description = texts[3]
				}
				apiItem := &ApiItem{
					Request: req,
				}
				ctx.ApiItems = append(ctx.ApiItems, apiItem)
			}
		case *langs.FunctionCall:
			if line.Name == "assert" {
				asserting = true
			}
			if asserting {
				ctx.ApiItems[len(ctx.ApiItems)-1].Response.lines = append(ctx.ApiItems[len(ctx.ApiItems)-1].Response.lines, line)
				break
			}
		default:
			if inited {
				if asserting {
					ctx.ApiItems[len(ctx.ApiItems)-1].Response.lines = append(ctx.ApiItems[len(ctx.ApiItems)-1].Response.lines, line)
				} else {
					ctx.ApiItems[len(ctx.ApiItems)-1].Request.lines = append(ctx.ApiItems[len(ctx.ApiItems)-1].Request.lines, line)
				}
			} else {
				ctx.InitLines = append(ctx.InitLines, line)
			}
		}
		index++
	}
	p.Ctx = ctx
	return nil
}

func (p *Pica) Parse() error {
	buffer, err := ioutil.ReadFile(p.FileName)
	if err != nil {
		return fmt.Errorf("parse error %v", err.Error())
	}
	p.parser = langs.NewParser(buffer)
	p.Block = p.parser.Parse()
	return nil
}

func (p *Pica) Convert() error {
	return nil
}

func (p *Pica) setApiInfoFromVmIntoCtx(ctx *ApiContext) {
	ctx.Name = p.vm.Lookup("name").(string)
	ctx.Version = p.vm.Lookup("version").(string)
	ctx.Author = p.vm.Lookup("author").(string)
	ctx.Description = p.vm.Lookup("description").(string)
}

func (p *Pica) setCtxHeader(ctx *ApiContext) {
	headersIf := p.vm.Lookup("headers")
	if headersIf != nil {
		headers := p.vm.Lookup("headers").(map[string]langs.Value)
		for key, val := range headers {
			switch v := val.(type) {
			case int:
				ctx.Headers.Set(key, string(v))
			case string:
				ctx.Headers.Set(key, v)
			default:
				panic("header's value part must be [string, int] types")
			}
		}
	}
}

func (p *Pica) runInitPartOfContext(ctx *ApiContext) {
	for _, line := range ctx.InitLines {
		p.vm.EvalStatement(line)
	}
	ctx.BaseUrl = p.vm.Lookup("baseUrl").(string)
	p.client = NewHttpClient(ctx.BaseUrl)

	p.setApiInfoFromVmIntoCtx(ctx)
	p.setCtxHeader(ctx)
}

func (p *Pica) RunApiContext() error {
	p.output.CopyRight()

	p.runInitPartOfContext(p.Ctx)

	for index, item := range p.Ctx.ApiItems {
		err := p.RunSingleApi(item)
		if err != nil {
			return fmt.Errorf("error when execute %d %s %s", index, item.Request.Name, err.Error())
		}
		if p.Delay > 0 {
			time.Sleep(time.Duration(p.Delay))
		}
	}
	p.output.Finished(len(p.Ctx.ApiItems), "all")
	return nil
}

func (p *Pica) getBody(tt string) []byte {
	val := p.vm.Lookup(tt).(map[string]langs.Value)
	data, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		return nil
	}
	return data
}

func (p *Pica) setRequestBody(item *ApiItem) {
	switch item.Request.Method {
	case "POST":
		item.Request.Body = p.getBody("post")
		p.output.EchoRequstIng("Posting", item.Request.Body)
		break
	case "PUT":
		item.Request.Body = p.getBody("put\n\n")
		p.output.EchoRequstIng("Posting", item.Request.Body)
		break
	case "PATCH":
		item.Request.Body = p.getBody("patch")
		p.output.EchoRequstIng("Posting", item.Request.Body)
		break
	}
}

func (p *Pica) resetRequestHeader() {

}

func (p *Pica) setRequestHeaderFromVm(item *ApiItem) {
	headersIf := p.vm.Lookup("headers")
	if headersIf != nil {
		headers := p.vm.Lookup("headers").(map[string]langs.Value)
		for key, val := range headers {
			switch v := val.(type) {
			case int:
				item.Request.Headers.Set(key, string(v))
			case string:
				item.Request.Headers.Set(key, v)
			default:
				panic("header's value part must be [string, int] types")
			}
		}
	}
}

func (p *Pica) RunSingleApi(item *ApiItem) error {
	p.output.EchoStartRequest(item.Request)
	p.vm.Assign("url", item.Request.Url)
	// Eval init scope statements
	for _, line := range item.Request.lines {
		p.vm.EvalStatement(line)
	}

	p.setRequestBody(item)
	p.setRequestHeaderFromVm(item)

	// do request
	res, err := p.client.Do(item.Request, p.vm)
	if err != nil {
		p.output.ErrorRequest(err)
		return fmt.Errorf("do http request error %s", err.Error())
	}
	item.Response.Headers = res.Header
	item.Response.Status = res.StatusCode
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	item.Response.Body = buf.Bytes()

	// Assign new header from response to vm
	headers := HttpHeaders2VmMap(item.Response.Headers)
	p.vm.Assign("headers", headers)
	p.vm.Assign("status", item.Response.Status)
	p.vm.Assign("body", item.Response.Body)

	contentType := item.Request.Headers.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		var jResults map[string]langs.Value
		err := json.Unmarshal(item.Response.Body, &jResults)
		if err != nil {
			panic(fmt.Errorf("json binding %s %s", err.Error(), item.Response.Body))
		}
		p.vm.Assign("json", jResults)

		PrintJson(&jResults)
	}

	// Eval item response statement
	for _, line := range item.Response.lines {
		p.vm.EvalStatement(line)
	}

	return nil
}
