package langs

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

// CombinedCode get combined code that using import
func CombinedCode(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	reg := regexp.MustCompile("import '(.*?)'")
	newData := reg.ReplaceAllStringFunc(string(data), func(part string) string {
		find := reg.FindStringSubmatch(part)
		if len(find) == 0 {
			panic(fmt.Sprintf("import error %s", part))
		}
		tmpCode, err := CombinedCode(find[1])
		if err != nil {
			return ""
		}
		return tmpCode
	})
	return newData, nil
}
