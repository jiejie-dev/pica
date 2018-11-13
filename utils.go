//go:generate statik -src=./assets
//go:generate go fmt statik/statik.go

package pica

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"regexp"

	"github.com/fixate/go-qs"
	"github.com/jeremaihloo/funny/langs"
	"github.com/rakyll/statik/fs"
	"github.com/shurcooL/github_flavored_markdown"
)

func HttpHeaders2VmMap(httpHeader http.Header) map[string]langs.Value {
	var headers = map[string]langs.Value{}
	for k, _ := range httpHeader {
		headers[k] = httpHeader.Get(k)
	}
	return headers
}

func VmMap2HttpHeaders(vmMap map[string]langs.Value) http.Header {
	headers := http.Header{}
	for k, v := range vmMap {
		headers.Set(k, v.(string))
	}
	return headers
}

func CompileUrl(url string, vm *langs.Interpreter) (string, Query, error) {
	queryValue := vm.LookupDefault("query", nil)
	query := Query{}
	if queryValue != nil {
		valQuery := queryValue.(map[string]langs.Value)
		for key, val := range valQuery {
			query[key] = val
		}
	}

	reg, err := regexp.Compile("<(.*?)>")
	if err != nil {
		return "", nil, err
	}
	result := reg.ReplaceAllStringFunc(url, func(repl string) string {
		repl = repl[1 : len(repl)-1]
		val := vm.Lookup(repl)
		switch val := val.(type) {
		case int:
			return string(val)
		case string:
			return val
		case float64:
			return strconv.FormatFloat(val, 'f', -1, 64)
		default:
			panic(fmt.Errorf("unsupport type [%s], only support [int][string]", langs.Typing(val)))
		}
	})
	if query == nil || len(query) == 0 {
		return result, nil, nil
	}
	qs, err := query.String()
	if err != nil {
		return "", nil, err
	}
	if qs != "?" {
		return fmt.Sprintf("%s?%s", result, qs), nil, nil
	}
	return result, query, nil
}

func BuildHtml(input []byte) string {
	output := github_flavored_markdown.Markdown(input)
	statikFS, _ := fs.New()
	tFile, err := statikFS.Open("/doc_template.html")
	if err != nil {
		panic(err)
	}
	template, err := ioutil.ReadAll(tFile)
	if err != nil {
		panic(err)
	}

	rs := strings.Replace(string(template), "[body]", string(output), -1)
	return rs
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

func PicaContextFromRunner(runner *APIRunner) *PicaContext {
	return &PicaContext{}
}
