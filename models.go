package pica

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jerloo/funny"
	"github.com/mitchellh/go-homedir"
)

type Task struct {
	Names []string
}

type Project struct {
	Name      string
	Version   string
	Author    string
	CreatedAt string
	LastRunAt string
	Tasks     map[string][]*Task
}

func NewProject(name, version, author, created, lastRunAt string) *Project {
	return &Project{
		Name:      name,
		Version:   version,
		Author:    author,
		CreatedAt: created,
		LastRunAt: lastRunAt,
	}
}

func (p *Project) Save() error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("pica.json", data, os.ModePerm)
}

//func (p *Project) RunTask(name string) error {
//	for key, val := range p.Tasks{
//		if key == name{
//			apiRunner := NewApiRunnerFromFile()
//		}
//	}
//}

type ApiRequest struct {
	Headers     http.Header
	Method      string
	Url         string
	Query       Query
	Name        string
	Description string
	Body        []byte
	lines       funny.Block
}

type ApiResponse struct {
	Headers http.Header
	Body    []byte
	Status  int
	lines   funny.Block

	saveLines funny.Block
}

type ApiItem struct {
	Request  *ApiRequest
	Response *ApiResponse
}

type PicaContext struct {
	Name        string
	Description string
	Author      string
	Version     string

	PicaVersion      string
	MaxArrayShowRows int

	ApiItems     []*ApiItem
	Headers      map[string]funny.Value
	VersionNotes *VersionNote
}

var DefaultHeaders = map[string]funny.Value{
	"Accept":          "* /*",
	"Accept-Language": "en-US,en;q=0.8",
	"Cache-Control":   "max-age=0",
	"User-Agent":      fmt.Sprintf("Pica Api Test Client/%s https://github.com/jerloo/pica", Version),
	"Connection":      "keep-alive",
	"Referer":         "http://www.baidu.com/",
	"Content-Type":    "application/json",
}

var DefaultInitScope = map[string]funny.Value{
	"headers": DefaultHeaders,
}

var PROFILE_HOME = ""

func init() {
	PROFILE_HOME, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	PROFILE_HOME = filepath.Join(PROFILE_HOME, ".pica")
	_, err = os.Stat(PROFILE_HOME)
	if err != nil {
		err = os.Mkdir(PROFILE_HOME, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}
