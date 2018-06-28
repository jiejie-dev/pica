package pica

import (
	"text/template"
	"fmt"
	"bytes"
	"encoding/json"
	"strings"
)

var (
	DEFAULT_DOC_TEMPLATE = `
# {{.Name}}

> {{.Description}}

Version: {{.Version}}
Author: {{.Author}}

## Init Scope

### Headers:

| --- | --- | -- |
| name| value | description |
{{range $i, $v := .Headers}}
| {{$i}} | {{$v}} | - |
{{end}}

## API

{{range $i, $item := .ApiItems }}
### {{$item.Request.Method}} {{$item.Request.Url}} {{$item.Request.Name}}
{{if $item.Request.Description}}
> {{$item.Request.Description}}
{{end}}

#### Query
| --- | --- | -- |
| name | type | description |

{{if $item.Request.Body}}
#### Body
| --- | --- | -- |
| name | type | description |
{{range $key, $val := $item.Request.Headers}}
| {{$key}} | {{$val}} | - |
{{end}}
{{end}}

#### Response
Headers:
| --- | --- | -- |
| name| value | description |
{{range $i, $v := $item.Response.Headers}}
| {{$i}} | {{$v}} | - |
{{end}}

{{if $item.Response.Body}}
Body:
'''
{{json $item.Response.Body}}
{{end}}
'''
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
	ctx      *ApiContext
	template *template.Template
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
		ctx:      ctx,
		template: t,
	}
}

func (g *DocGenerator) Get() ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := g.template.Execute(buffer, g.ctx)
	if err != nil {
		return nil, fmt.Errorf("generate doc %s", err.Error())
	}
	return buffer.Bytes(), nil
}
