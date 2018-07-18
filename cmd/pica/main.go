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

	// runs
	cmdRun            = app.Command("run", "Run api file.")
	runFileName       = cmdRun.Arg("filename", "Api file.").Default("pica.fun").ExistingFile()
	runAPINames       = cmdRun.Arg("apiNames", "Api names to excute").Strings()
	runDelay          = cmdRun.Flag("delay", "Delay after one api request.").Int()
	runOutput         = cmdRun.Flag("output", "Output file.").String()
	runOutputTemplate = cmdRun.Flag("template", "Doc template.").Default("pica.doc.t").String()

	// formats
	cmdFormat      = app.Command("format", "Format api file.")
	formatFileName = cmdFormat.Arg("filename", "Format file.").Default("pica.fun").ExistingFile()
	formatSave     = cmdFormat.Flag("save", "Save formated file.").Default("1").Bool()
	formatPrint    = cmdFormat.Flag("print", "Print to stdout").Default("1").Bool()

	// servers
	cmdServer  = app.Command("serve", "Run a document website.")
	apiDocFile = cmdServer.Flag("file", "Api File.").Default("pica.md").String()
	docPort    = cmdServer.Flag("port", "Port for doc.").Default("9000").Int()

	// inits
	cmdInit         = app.Command("init", "Init a new api file from template.")
	cmdInitFileName = cmdInit.Arg("filename", "The filename to initialize.").Default("pica.fun").String()
	cmdInitTemplate = cmdInit.Arg("template", "Init Template. Support local and remote file uri.").String()

	// configs
	cmdConfig = app.Command("config", "Config pica.")
	username  = cmdConfig.Flag("username", "Username.").String()
	email     = cmdConfig.Flag("email", "Email.").String()

	// versions
	cmdVersion = app.Command("version", "Version controls.")

	versionCommit = cmdVersion.Command("commit", "Commit a new version for api file and doc.")
	commitMsg     = versionCommit.Arg("commit-msg", "Msg for this commit.").Required().String()
	commitFile    = versionCommit.Arg("filename", "Api file to be committed.").Default("pica.fun").String()

	versionReset = cmdVersion.Command("reset", "Reset one commit.")

	versionLog = cmdVersion.Command("log", "Show all logs of one api file.")
	logFile    = versionLog.Arg("filename", "Api file to show logs.").Default("pica.fun").String()

	versionDiff     = cmdVersion.Command("diff", "Diff between two commits")
	diffCommitOlder = versionDiff.Arg("hash-older", "The hash of older one.").String()
	diffCommitNewer = versionDiff.Arg("hash-newer", "The hash of newer one.").Default("HEAD").String()
)

func main() {

	c, err := app.Parse(os.Args[1:])
	if !*debug {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("\nerror: %s\n", err)
			}
		}()
	}
	switch kingpin.MustParse(c, err) {
	case cmdRun.FullCommand():
		pica.Run(*runFileName, *runAPINames, *runDelay, *runOutput, *runOutputTemplate)
		break
	case cmdFormat.FullCommand():
		pica.Format(*formatFileName, *formatSave, *formatPrint)
		break
	case cmdServer.FullCommand():
		err := pica.Serve(*apiDocFile, *docPort)
		if err != nil {
			panic(err)
		}
		break
	case cmdInit.FullCommand():
		pica.Init(*cmdInitFileName, *cmdInitTemplate)
		break
	case versionCommit.FullCommand():
		pica.VersionCommit(*commitFile, *commitMsg)
		break
	case versionLog.FullCommand():
		pica.VersionLog(*logFile)
		break
	default:
		kingpin.Usage()
	}

}
