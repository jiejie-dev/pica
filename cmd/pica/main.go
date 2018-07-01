package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/jeremaihloo/pica"
)

var (
	app   = kingpin.New("pica", "A command line for api test and doc generate.")
	debug = app.Flag("debug", "Debug mode.").Bool()

	cmdRun            = app.Command("run", "Run api file.")
	runFileName       = cmdRun.Arg("filename", "Api file.").Default("pica.fun").ExistingFile()
	runApiNames       = cmdRun.Arg("apiNames", "Api names to excute").Strings()
	runDelay          = cmdRun.Flag("delay", "Delay after one api request.").Int()
	runOutput         = cmdRun.Flag("output", "Output file.").String()
	runOutputTemplate = cmdRun.Flag("template", "Doc template.").Default("pica.doc.t").String()

	cmdFormat      = app.Command("format", "Format api file.")
	formatFileName = cmdFormat.Arg("filename", "Format file.").Default("pica.fun").ExistingFile()
	formatSave     = cmdFormat.Flag("save", "Save formated file.").Default("1").Bool()
	formatPrint    = cmdFormat.Flag("print", "Print to stdout").Default("1").Bool()

	cmdServer  = app.Command("serve", "Run a document website.")
	apiDocFile = cmdServer.Flag("file", "Api File.").Default("pica.md").String()
	docPort    = cmdServer.Flag("port", "Port for doc.").Default("9000").Int()

	cmdInit         = app.Command("init", "Init a new api file from template.")
	cmdInitFileName = cmdInit.Arg("filename", "The filename to initialize.").Default("pica.fun").String()
	cmdInitTemplate = cmdInit.Arg("template", "Init Template. Support local and remote file uri.").String()

	cmdConfig = app.Command("config", "Config pica.")
	username  = cmdConfig.Flag("username", "Username.").String()
	email     = cmdConfig.Flag("email", "Email.").String()
)

func main() {
	if !*debug {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("\nerror: %s\n", err)
			}
		}()
	}
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case cmdRun.FullCommand():
		pica.Run(*runFileName, *runApiNames, *runDelay, *runOutput, *runOutputTemplate)
		break
	case cmdFormat.FullCommand():
		pica.Format(*formatFileName, *formatSave, *formatPrint)
		break
	case cmdServer.FullCommand():
		pica.Serve(*apiDocFile, *docPort)
		break
	case cmdInit.FullCommand():
		pica.Init(*cmdInitFileName, *cmdInitTemplate)
		break
	default:
		kingpin.Usage()
	}

}
