package pica

import (
	"net/http"
	"os"
	"github.com/jeremaihloo/funny/langs"
	"fmt"
	"strings"
	"io/ioutil"
	"time"
	"encoding/json"
	"bytes"
)

type ApiRequest struct {
	Headers     http.Header
	Method      string
	Url         string
	Name        string
	Description string
	Body        []byte
	lines       langs.Block
}

type ApiResponse struct {
	Headers http.Header
	Body    []byte
	Status  int
	lines   langs.Block
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
}

type Pica struct {
	FileName        string
	Delay           int
	Output          *os.File
	Debug           bool
	IfRun           bool
	IfConvert       bool
	IfDoc           bool
	IfServer        bool
	IfFormat        bool
	DocTempalteFile string

	vm     *langs.Interpreter
	parser *langs.Parser
	Block  langs.Block
	client *HttpClient
}

func NewPica(
	filename string,
	output *os.File,
	delay int,
	ifRun,
	ifFormat,
	ifConvert,
	ifDoc,
	ifServer bool,
	template string) *Pica {
	return &Pica{
		FileName:        filename,
		Output:          output,
		Delay:           delay,
		IfRun:           ifRun,
		IfConvert:       ifConvert,
		IfDoc:           ifDoc,
		IfServer:        ifServer,
		IfFormat:        ifFormat,
		DocTempalteFile: template,

		vm: langs.NewInterpreterWithScope(langs.Scope{}),
	}
}

func (p *Pica) Run() error {
	err := p.Parse()
	if err != nil {
		return err
	}
	ctx, err := p.ParseApiContext()
	if p.IfFormat {
		return p.Format()
	} else if p.IfConvert {
		return p.Convert()
	} else if p.IfDoc {
		return p.Document(ctx)
	} else if p.IfRun {
		return p.RunApiContext(ctx)
	}
	return nil
}

func (p *Pica) Document(ctx *ApiContext) error {
	p.runInitPartOfContext(ctx)

	data, err := ioutil.ReadFile(p.DocTempalteFile)
	if err != nil {
		data = []byte(DEFAULT_TEMPLATE)
	}
	generator := NewGenerator(ctx, string(data))
	results, err := generator.Get()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", results)
	if _, err := p.Output.Stat(); err == nil {
		_, err := p.Output.Write(results)
		if err != nil {
			return err
		}
		p.Output.Close()
	}
	return nil
}

func (p *Pica) ParseApiContext() (*ApiContext, error) {
	headers := VmMap2HttpHeaders(DefaultHeaders)
	ctx := &ApiContext{
		Headers: &headers,
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
				ctx.ApiItems[len(ctx.ApiItems)-1].Request.lines = append(ctx.ApiItems[len(ctx.ApiItems)-1].Request.lines, line)
			} else {
				ctx.InitLines = append(ctx.InitLines, line)
			}
		}
		index++
	}
	return ctx, nil
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

func (p *Pica) Format() error {
	p.parser.Consume("")
	flag := 0
	for {
		item := p.parser.ReadStatement()
		if item == nil {
			break
		}
		switch item.(type) {
		case *langs.NewLine:
			flag += 1
			if flag < 1 {
				continue
			}
			break
		}
		fmt.Printf("%s", item.String())
	}
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

func (p *Pica) RunApiContext(ctx *ApiContext) error {

	p.runInitPartOfContext(ctx)

	for index, item := range ctx.ApiItems {
		err := p.RunSingleApi(item)
		if err != nil {
			return fmt.Errorf("error when execute %d %s %s", index, item.Request.Name, err.Error())
		}
		if p.Delay > 0 {
			time.Sleep(time.Duration(p.Delay))
		}
	}
	fmt.Printf("\n\nFinished. [%d] api requests, [%s] passed", len(ctx.ApiItems), "all")
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
		fmt.Printf("Posting ...\n%s\n\n", item.Request.Body)
		break
	case "PUT":
		item.Request.Body = p.getBody("put\n\n")
		fmt.Printf("Putting ...\n%s", item.Request.Body)
		break
	case "PATCH":
		item.Request.Body = p.getBody("patch")
		fmt.Printf("Patching ...\n%s\n\n", item.Request.Body)
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
	fmt.Printf("Starting request [%s %s %s]\n\n", item.Request.Method, item.Request.Url, item.Request.Name)
	p.vm.Assign("url", item.Request.Url)
	// Eval init scope statements
	for _, line := range item.Request.lines {
		p.vm.EvalStatement(line)
	}

	p.setRequestBody(item)
	p.setRequestHeaderFromVm(item)

	// do request
	res, err := p.client.Do(item.Request)
	if err != nil {
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
		var jResults map[string]interface{}
		err := json.Unmarshal(item.Response.Body, &jResults)
		if err != nil {
			return fmt.Errorf("json binding %s %s", err.Error(), item.Response.Body)
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
