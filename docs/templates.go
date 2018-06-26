package docs


const (
	DOC = `
# {{.Name}}

> {{.Description}}

Version: {{.Version}}
Author: {{.Author}}

## Init Scope

### Headers:

| --- | --- | -- |
| name| value | description |
{{range $i, $v= .Headers}}

{{end}}

### Vars
| --- | --- | -- |
| name| value | description |
{{range $i, $v= .Headers}}

{{end}}

## API

{{range $i, $v = .Items }}
### {{$item.Request.Method}} {{$item.Request.Url}} {{$item.Request.Name}}

> {{$Description}}

#### Query
| --- | --- | -- |
| name | type | description |

{{if $item.Body}}
#### Body
| --- | --- | -- |
| name | type | description |
{{range $key, $val = $item.Request.Headers}}
| {{$key}} | {{$val}} | {{$ |
{{end}}
`
)



