package pica

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"bytes"
	"path/filepath"
	"text/template"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/howeyc/fsnotify"
	"github.com/jeremaihloo/funny/langs"
	_ "github.com/jeremaihloo/pica/statik"
	"github.com/rakyll/statik/fs"
	"github.com/shurcooL/github_flavored_markdown"
	"gopkg.in/AlecAivazis/survey.v1"
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
		sfs, err := fs.New()
		if err != nil {
			return err
		}
		file, err := sfs.Open("/api_file_template.fun")
		if err != nil {
			return err
		}
		DEFAULT_API_FILE_TEMPLATE, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
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
	output := github_flavored_markdown.Markdown(input)
	r := gin.Default()
	template := BuildHTML(input)

	r.GET("/", func(c *gin.Context) {
		rs := strings.Replace(template, "[body]", string(output), -1)
		c.Data(200, "text/html", []byte(rs))
	})

	srv := &http.Server{
		Addr:    ":9090",
		Handler: r,
	}

	run := func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}

	go run()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case <-watcher.Event:
				fmt.Println("update")

				input, err = ioutil.ReadFile(filename)
				if err != nil {
					panic(err)
				}
				output = github_flavored_markdown.Markdown(input)
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Hang so program doesn't exit
	<-done

	/* ... do stuff ... */
	watcher.Close()
	return nil
}

func GenDocument(apiRunner *APIRunner, output string) error {
	if !strings.HasSuffix(output, ".md") && !strings.HasSuffix(output, ".html") {
		return errors.New("only .md and .html supported")
	}
	data, err := ioutil.ReadFile("")
	if err != nil {
		data = []byte(DEFAULT_DOC_TEMPLATE)
	}
	generator := NewMarkdownDocGenerator(apiRunner, string(data), output)
	results, err := generator.Get()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", results)
	// if _, err := os.Stat(p.Output); err == nil {
	if strings.HasSuffix(output, ".html") {
		results = []byte(BuildHTML(results))
	}
	err = ioutil.WriteFile(output, results, os.ModePerm)
	if err != nil {
		return err
	}
	// }
	// return errors.New("file already exists")
	return nil
}

func VersionCommit(commitFile, commitMsg string) {
	ctrl := NewApiVersionController(commitFile)
	hash, err := ctrl.Commit(commitMsg)
	if err != nil {
		panic(err)
	}
	fmt.Println(hash)
}

func VersionLog(filename string) {
	ctrl := NewApiVersionController(filename)
	commits, err := ctrl.GetCommits()
	if err != nil {
		panic(err)
	}
	for index, item := range commits {
		color.Green(item.String())
		if index != len(commits)-1 {
			fmt.Println("==========================================================")
		}
	}
}
