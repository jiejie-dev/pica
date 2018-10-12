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