package pica

import (
	"text/template"
	"fmt"
	"bytes"
)

const (
	DEFAULT_TEMPLATE = `
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
{{range $i, $v := $item.Response.Body}}
| {{$i}} | {{$v}} | - |
{{end}}

{{if $item.Response.Body}}
Body:
{{$item.Response.Body}}
{{end}}

{{end}}
`
)

type DocGenerator struct {
	ctx      *ApiContext
	template *template.Template
}

func NewDefaultGenerator(ctx *ApiContext) *DocGenerator {
	return NewGenerator(ctx, DEFAULT_TEMPLATE)
}

func NewGenerator(ctx *ApiContext, tStr string) *DocGenerator {
	t := template.Must(template.New("doc").Parse(tStr))

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
