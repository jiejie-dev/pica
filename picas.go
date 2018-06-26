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
	FileName  string
	Delay     int
	Output    *os.File
	Debug     bool
	IfRun     bool
	IfConvert bool
	IfDoc     bool
	IfServer  bool
	IfFormat  bool

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
	ifServer bool) *Pica {
	return &Pica{
		FileName:  filename,
		Output:    output,
		Delay:     delay,
		IfRun:     ifRun,
		IfConvert: ifConvert,
		IfDoc:     ifDoc,
		IfServer:  ifServer,
		IfFormat:  ifFormat,
		vm:        langs.NewInterpreterWithScope(langs.Scope{}),
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
	} else if p.IfRun {
		return p.RunApiContext(ctx)
	}
	return nil
}

func (p *Pica) ParseApiContext() (ApiContext, error) {
	ctx := ApiContext{}
	inited := false
	index := 0
	asserting := false
	for index < len(p.Block)-1 {
		line := p.Block[index]
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
				ctx.ApiItems = append(ctx.ApiItems, apiItem)
			}
		case *langs.FunctionCall:
			if line.Name == "must" {
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

func (p *Pica) RunApiContext(ctx ApiContext) error {
	for _, line := range ctx.InitLines {
		p.vm.EvalStatement(line)
	}
	ctx.BaseUrl = p.vm.Lookup("baseUrl").(string)
	p.client = NewHttpClient(ctx.BaseUrl)

	ctx.Name = p.vm.Lookup("name").(string)
	ctx.Version = p.vm.Lookup("version").(string)
	ctx.Author = p.vm.Lookup("author").(string)
	ctx.Description = p.vm.Lookup("description").(string)

	for index, item := range ctx.ApiItems {
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
