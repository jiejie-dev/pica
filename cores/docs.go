package pica

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	_ "github.com/jeremaihloo/pica/statik"
	"github.com/rakyll/statik/fs"
)

var DEFAULT_DOC_TEMPLATE = ""

func init() {
	sfs, err := fs.New()
	if err != nil {
		panic(err)
	}
	file, err := sfs.Open("/doc_template.md")
	if err != nil {
		panic(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	DEFAULT_DOC_TEMPLATE = strings.Replace(string(content), "'''", "```", -1)
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

type MarkdownDocGenerator struct {
	runner     *APIRunner
	template   *template.Template
	versionCtl *ApiVersionController
}

func NewMarkdownDocGenerator(runner *APIRunner, theme, output string) *MarkdownDocGenerator {
	if theme != "default" {
		file, err := os.Open(filepath.Join(PROFILE_HOME, "doc_template.md"))
		if err != nil {
			panic(err)
		}
		bts, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
		DEFAULT_DOC_TEMPLATE = string(bts)
	}
	fnMap := template.FuncMap{
		"json": TSafeJson,
	}
	t := template.Must(template.New("doc").Funcs(fnMap).Parse(DEFAULT_DOC_TEMPLATE))

	return &MarkdownDocGenerator{
		runner:     runner,
		template:   t,
		versionCtl: NewApiVersionController(output),
	}
}

func (g *MarkdownDocGenerator) Get() ([]byte, error) {
	note, err := g.versionCtl.Notes()
	if err != nil {
		panic(err)
		return nil, err
	}
	buffer := new(bytes.Buffer)
	ctx := PicaContextFromRunner(g.runner)
	ctx.VersionNotes = note
	err = g.template.Execute(buffer, ctx)
	if err != nil {
		panic(err)
		return nil, fmt.Errorf("generate doc %s", err.Error())
	}
	return buffer.Bytes(), nil
}
