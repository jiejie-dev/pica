package pica

import (
	"encoding/json"
	"os"
	"strings"
	"text/template"

	_ "embed"

	"github.com/jerloo/funny"
	"github.com/rbretecher/go-postman-collection"
)

//go:embed gen_template.funny.txt
var GenFromPostmanTemplate string

func GenerateScriptsByPostman(postmanFile string) string {
	file, err := os.Open(postmanFile)
	defer func() {
		file.Close()
	}()

	if err != nil {
		panic(err)
	}

	c, err := postman.ParseCollection(file)
	if err != nil {
		panic(err)
	}
	items := getItems(c.Items)
	t, err := template.New("postman").Funcs(template.FuncMap{
		"getQueryNames": getQuery,
		"joinQuery":     joinQuery,
		"joinStrArr":    joinStrArr,
		"body2Map":      body2Map,
	}).Parse(GenFromPostmanTemplate)
	if err != nil {
		panic(err)
	}
	builder := strings.Builder{}
	err = t.Execute(&builder, items)
	if err != nil {
		panic(err)
	}

	// for _, item := range items {
	// 	item.Request.Method
	// 	item.Request.URL.Variables
	// 	item.Request.Body.Raw
	// }
	data := builder.String()
	return funny.Format([]byte(data))
}

func body2Map(body string) (result map[string]interface{}) {
	err := json.Unmarshal([]byte(body), &result)
	if err != nil {
		panic(err)
	}
	return
}

func joinStrArr(items []string, sp string) string {
	return strings.Join(items, sp)
}

func getQuery(item *postman.Items) []string {
	var arr []string
	if item.Request.URL.Query != nil {
		for _, arg := range item.Request.URL.Query.([]interface{}) {
			a := arg.(map[string]interface{})
			arr = append(arr, a["key"].(string))
		}
	}
	return arr
}

func joinQuery(item *postman.Items) string {
	var arr []string
	if item.Request.URL.Query != nil {
		for _, arg := range item.Request.URL.Query.([]interface{}) {
			if arg != nil {
				a := arg.(map[string]interface{})
				arr = append(arr, a["key"].(string))
			}
		}
	}
	return strings.Join(arr, ",")
}

func getItems(item []*postman.Items) (results []*postman.Items) {
	for _, child := range item {
		if child.IsGroup() {
			results = append(results, getItems(child.Items)...)
		} else {
			results = append(results, child)
		}
	}
	return results
}
