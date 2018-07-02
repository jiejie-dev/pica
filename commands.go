package pica

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"bytes"
	"path/filepath"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/jeremaihloo/funny/langs"
	survey "gopkg.in/AlecAivazis/survey.v1"
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

type InitInfo struct {
	Name        string
	Description string
	Author      string
	Version     string
	BaseUrl     string
}

func Init(filename, templateName string) error {
	if _, err := os.Stat(filename); err == nil {
		override := false
		var q = &survey.Confirm{
			Message: "Api file already exists. Override ?",
		}

		survey.AskOne(q, &override, nil)
		if !override {
			return nil
		}
	}
	data, err := ioutil.ReadFile(templateName)
	if err != nil {
		data = []byte(DEFAULT_API_FILE_TEMPLATE)
	}
	info := InitInfo{}
	// the questions to ask
	dir, _ := os.Getwd()
	l := filepath.SplitList(dir)
	dir = l[len(l)-1]
	var qs = []*survey.Question{
		{
			Name: "Name",
			Prompt: &survey.Input{
				Message: "What is your api file name?",
				Default: filepath.Base(dir),
			},
			Validate: survey.Required,
		},
		{
			Name: "Description",
			Prompt: &survey.Input{
				Message: "What is then description of your apis ?",
			},
		},
		{
			Name: "Author",
			Prompt: &survey.Input{
				Message: "Who are you ?",
			},
		},
		{
			Name: "Version",
			Prompt: &survey.Input{
				Message: "What version to start ?",
			},
		},
		{
			Name: "BaseUrl",
			Prompt: &survey.Input{
				Message: "What is the baseUrl of your apis ?",
			},
		},
	}
	err = survey.Ask(qs, &info)
	if err != nil {
		panic(err)
	}
	t := template.Must(template.New("api").Parse(string(data)))
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, &info)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, buffer.Bytes(), os.ModePerm)
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
