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

	lines langs.Block
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
	BaseUrl     string
	Headers     *http.Header
	InitVars    langs.Scope

	initLines langs.Block
	apiItems  []*ApiItem
}

type Pica struct {
	FileName  string
	Delay     int
	Output    *os.File
	Debug     bool
	IfRun     bool
	IfConvert bool
	IfDoc     bool
	IfServer  bool

	vm     *langs.Interpreter
	parser *langs.Parser
	block  langs.Block

	initVars langs.Scope

	client *HttpClient
}

func NewPica(filename string, output *os.File, delay int, ifRun, ifConvert, ifDoc, ifServer bool) *Pica {
	return &Pica{
		FileName:  filename,
		Output:    output,
		Delay:     delay,
		IfRun:     ifRun,
		IfConvert: ifConvert,
		IfDoc:     ifDoc,
		IfServer:  ifServer,
		vm:        langs.NewInterpreter(langs.Scope{}),
	}
}

func (p *Pica) Run() error {
	err := p.Parse()
	if err != nil {
		return err
	}
	ctx, err := p.parseApiContext()
	if p.IfRun {
		err := p.RunApiContext(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pica) parseApiContext() (ApiContext, error) {
	ctx := ApiContext{}
	inited := false
	index := 0
	asserting := false
	for index < len(p.block)-1 {
		line := p.block[index]
		switch line := line.(type) {
		case *langs.Comment:
			text := strings.Trim(line.Value, " ")
			texts := strings.Split(text, " ")
			if len(texts) < 2 {
				break
			}
			methods := []string{"GET", "POST", "DELETE", "PUT"}
			flag := false
			for _, item := range methods {
				if strings.ToUpper(item) == strings.ToUpper(texts[0]) {
					flag = true
				}
			}
			if flag {
				inited = true
				req := ApiRequest{
					Method: texts[0],
					Url:    texts[1],
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
				ctx.apiItems = append(ctx.apiItems, apiItem)
			}
		case *langs.FunctionCall:
			if line.Name == "must" {
				asserting = true
			}
			if asserting {
				ctx.apiItems[len(ctx.apiItems)-1].Response.lines = append(ctx.apiItems[len(ctx.apiItems)-1].Response.lines, line)
				break
			}
		default:
			if inited {
				ctx.apiItems[len(ctx.apiItems)-1].Request.lines = append(ctx.apiItems[len(ctx.apiItems)-1].Request.lines, line)
			} else {
				ctx.initLines = append(ctx.initLines, line)
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
	p.block = p.parser.Parse()
	return nil
}

func (p *Pica) Convert() error {
	return nil
}

func (p *Pica) RunApiContext(ctx ApiContext) error {
	for _, line := range ctx.initLines {
		p.vm.EvalStatement(line)
	}
	ctx.BaseUrl = p.vm.Lookup("baseUrl").(string)
	p.client = NewHttpClient(ctx.BaseUrl)

	for index, item := range ctx.apiItems {
		err := p.RunSingleApi(item)
		if err != nil {
			return fmt.Errorf("error when execute %d %s %s", index, item.Request.Name, err.Error())
		}
		if p.Delay > 0 {
			time.Sleep(time.Duration(p.Delay))
		}
	}
	return nil
}

func (p *Pica) RunSingleApi(item *ApiItem) error {
	p.vm.Assign("url", item.Request.Url)
	for _, line := range item.Request.lines {
		p.vm.EvalStatement(line)
	}
	switch item.Request.Method {
	case "POST":
		item.Request.Body = []byte(p.vm.Lookup("post").(string))
		break
	case "PUT":
		item.Request.Body = []byte(p.vm.Lookup("put").(string))
		break
	case "PATCH":
		item.Request.Body = []byte(p.vm.Lookup("patch").(string))
		break
	}
	res, err := p.client.Do(item.Request)
	if err != nil {
		return err
	}
	item.Response.Headers = res.Header
	item.Response.Status = res.StatusCode
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	item.Response.Body = buf.Bytes()

	p.vm.Assign("headers", item.Response.Headers)
	p.vm.Assign("status", item.Response.Status)
	p.vm.Assign("body", item.Response.Body)
	if item.Request.Headers.Get("Content-Type") == "application/json" {
		var jResults map[string]interface{}
		err := json.Unmarshal(item.Response.Body, jResults)
		if err != nil {
			return err
		}
		p.vm.Assign("json", jResults)
	}

	for _, line := range item.Response.lines {
		p.vm.EvalStatement(line)
	}
	return nil
}
