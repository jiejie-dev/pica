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

type ScriptsGenerator interface {
	Name() string
	Generate(filename string) string
}

func NewScriptsGenerator(name string) ScriptsGenerator {
	return map[string]ScriptsGenerator{
		"Postman":  &PostmanScriptsGenerator{},
		"Swagger2": &Swagger2ScriptsGenerator{},
	}[name]
}

type Swagger2ScriptsGenerator struct {
}

func (generator *Swagger2ScriptsGenerator) Name() string {
	return "Swagger2"
}

func (generator *Swagger2ScriptsGenerator) Generate(filename string) string {
	return ""
}

type PostmanScriptsGenerator struct {
}

func (generator *PostmanScriptsGenerator) Name() string {
	return "Postman"
}

func (generator *PostmanScriptsGenerator) Generate(filename string) string {
	file, err := os.Open(filename)
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
	items := genItems("", c.Items)
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
			item.Name = safePinYin(item.Name)
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

//go:embed gen_template.funny.txt
var GenFromPostmanTemplate string

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

func safePinYin(hanz string) string {
	ps := pinyin.LazyConvert(hanz, nil)
	for index, psItem := range ps {
		ps[index] = Capitalize(psItem)
	}
	return strings.Join(ps, "")
}

func genItems(nameRoot string, item []*postman.Items) (results []*postman.Items) {
	for _, child := range item {
		if child.IsGroup() {
			results = append(results, genItems(fmt.Sprintf("%s_%s", nameRoot, safePinYin(child.Name)), child.Items)...)
		} else {
			if nameRoot != "" {
				child.Name = fmt.Sprintf("%s_%s", nameRoot, child.Name)
			}
			results = append(results, child)
		}
	}
	return results
}
