package pica

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"
	"unicode"

	_ "embed"

	"github.com/jerloo/funny"
	"github.com/mozillazg/go-pinyin"
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
	names := make(map[string]bool)
	items := getItems(c.Items)
	for index, item := range items {

		if item.Description == "" {
			item.Description = item.Name
		}
		hansName := false
		for _, r := range item.Name {
			if unicode.Is(unicode.Han, r) {
				hansName = true
			}
		}
		if hansName {
			ps := pinyin.LazyConvert(item.Name, nil)
			for index, psItem := range ps {
				ps[index] = Capitalize(psItem)
			}
			item.Name = strings.Join(ps, "")
		}
		if _, ok := names[item.Name]; ok {
			item.Name = fmt.Sprintf("%s%d", item.Name, index)
		} else {
			names[item.Name] = true
		}
	}
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
	return funny.Format([]byte(data), "")
}

//字符首字母大写转换
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
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
