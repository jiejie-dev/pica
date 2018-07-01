package pica

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jeremaihloo/funny/langs"
	"gopkg.in/russross/blackfriday.v2"
)

const (
	DEFAULT_API_FILE_TEMPLATE = `
name = '{{.Name}}'
description = '{{.Description}}'
author = '{{.Author}}'
version = '{{.Version}}'

baseUrl = '{{.BaseUrl}}'

headers = {
    'Content-Type' = 'application/json'
}

// Apis format: [method] [path] [description]

// GET /api/users 获取用户列表
headers['Authorization'] = 'slfjaslkfjlasdjfjas=='

// POST /api/users 新建用户
post = {
    // 用户名
    'name' = 'test'
    // 密码
    'age' = 10
}
`
)

func Init(filename, template string) error {
	data, err := ioutil.ReadFile(template)
	if err != nil {
		data = []byte(DEFAULT_API_FILE_TEMPLATE)
	}
	err = ioutil.WriteFile(filename, data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func Run(filename string, apiNames []string, delay int, outputFile, outputTemplate string) error {
	pica := NewPica(filename, delay, outputFile, outputTemplate)
	return pica.Run()
}

func Format(filename string, save, print bool) (string, error) {
	fw := strings.Builder{}
	output := func(text string) {
		if print {
			fmt.Printf("%s", text)
		}
		if save {
			fw.WriteString(text)
		}
	}
	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("parse error %v", err.Error())
	}
	parser := langs.NewParser(buffer)
	parser.Consume("")
	flag := 0
	for {
		item := parser.ReadStatement()
		if item == nil {
			break
		}
		switch item.(type) {
		case *langs.NewLine:
			flag += 1
			if flag < 1 {
				continue
			}
			break
		}
		output(fmt.Sprintf("%s", item.String()))
	}
	if save {
		ioutil.WriteFile(filename, []byte(fw.String()), os.ModePerm)
	}
	return fw.String(), nil
}

func Serve(filename string, port int) error {
	if !strings.HasSuffix(filename, ".md") {
		return fmt.Errorf("unknow type of doc files support [md]")
	}
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	output := blackfriday.Run(input)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.Data(200, "text/plain", output)
	})
	return r.Run(fmt.Sprintf(":%d", port))
}
