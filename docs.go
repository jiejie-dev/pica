package pica

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
)

const (
	DefaultCopyright = `
// =========================================== //
// -------------- Pica Api File -------------- //
// ----https://github.com/jeremaihloo/pica---- //
// =========================================== //
`
)

var (
	DEFAULT_DOC_TEMPLATE = `

# {{.Name}}

> {{.Description}}

Version: {{.Version}}

Author: {{.Author}}

## Init Scope

### Headers

| name| value | description |
| --- | --- | --- |
{{range $i, $v := .Headers}}| {{$i}} | {{$v}} | - |
{{end}}

## API

{{range $i, $item := .ApiItems }}
### {{$item.Request.Method}} {{$item.Request.Url}} {{$item.Request.Name}}
{{if $item.Request.Description}}
> {{$item.Request.Description}}
{{end}}
{{if $item.Request.Query}}
#### Query
| --- | --- | --- |
| name | type | description |
{{end}}
{{if $item.Request.Body}}
#### Body
| name | type | description |
| --- | --- | --- |
{{range $key, $val := $item.Request.Headers}}| {{$key}} | {{$val}} | - |
{{end}}
{{end}}

#### Response

Headers:

| name| value | description |
| --- | --- | --- |
{{range $i, $v := $item.Response.Headers}}| {{$i}} | {{$v}} | - |
{{end}}

{{if $item.Response.Body}}
Body:
'''
{{json $item.Response.Body}}
{{end}}
'''
{{end}}

{{if ne (len .VersionNotes.Changes) 0}}
## Relase Notes

{{range $index,$item := .VersionNotes.Changes}}
### {{$item.Commit.Hash}}

{{$item.Commit.Message}}
{{end}}
{{end}}
`
)

func init() {
	DEFAULT_DOC_TEMPLATE = strings.Replace(DEFAULT_DOC_TEMPLATE, "'''", "```", -1)
}

func TSafeJson(obj interface{}) string {
	switch obj := obj.(type) {
	case map[string]interface{}:
		return TSafeJson(&obj)
	case *map[string]interface{}:
		data, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf("\nJson:\n%s\n\n", data)
	case []byte:
		var newObj map[string]interface{}
		err := json.Unmarshal(obj, &newObj)
		if err != nil {
			panic(err)
		}
		return TSafeJson(newObj)
	default:
		return "unknow type object to serialize"
	}
}

type DocGenerator struct {
	ctx        *ApiContext
	template   *template.Template
	versionCtl *ApiVersionController
}

func NewDefaultGenerator(ctx *ApiContext) *DocGenerator {
	return NewGenerator(ctx, DEFAULT_DOC_TEMPLATE)
}

func NewGenerator(ctx *ApiContext, tStr string) *DocGenerator {
	fnMap := template.FuncMap{
		"json": TSafeJson,
	}
	t := template.Must(template.New("doc").Funcs(fnMap).Parse(tStr))

	return &DocGenerator{
		ctx:        ctx,
		template:   t,
		versionCtl: NewApiVersionController(ctx.Pica.FileName),
	}
}

func (g *DocGenerator) Get() ([]byte, error) {
	note, err := g.versionCtl.Notes()
	if err != nil {
		panic(err)
		return nil, err
	}
	g.ctx.VersionNotes = note
	buffer := new(bytes.Buffer)
	err = g.template.Execute(buffer, g.ctx)
	if err != nil {
		panic(err)
		return nil, fmt.Errorf("generate doc %s", err.Error())
	}
	return buffer.Bytes(), nil
}
