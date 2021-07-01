package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/jerloo/funny"
	"github.com/jerloo/pica"
)

var (
	BuildAt = "unknow"
	COMMIT  = "unknow"
	GOLANG  = "unknow"

	app   = kingpin.New("pica", "A command line for api test and doc generate.")
	debug = app.Flag("debug", "Debug mode.").Bool()

	// runs
	cmdRun         = app.Command("run", "Run api file.")
	runFileName    = cmdRun.Arg("filename", "Api file.").Default("pica.fun").ExistingFile()
	runAPINames    = cmdRun.Arg("apiNames", "Api names to excute").Strings()
	runDelay       = cmdRun.Flag("delay", "Delay after one api request.").Int()
	runOutput      = cmdRun.Flag("output", "Output file.").String()
	runOutputTheme = cmdRun.Flag("theme", "Doc template.").Default("default").String()

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
	cmdVersion = app.Command("vc", "Version controls.")

	versionCommit = cmdVersion.Command("commit", "Commit a new version for api file and doc.")
	commitMsg     = versionCommit.Arg("commit-msg", "Msg for this commit.").Required().String()
	commitFile    = versionCommit.Arg("filename", "Api file to be committed.").Default("pica.fun").String()

	versionReset = cmdVersion.Command("reset", "Reset one commit.")

	versionLog = cmdVersion.Command("log", "Show all logs of one api file.")
	logFile    = versionLog.Arg("filename", "Api file to show logs.").Default("pica.fun").String()

	versionDiff     = cmdVersion.Command("diff", "Diff between two commits")
	diffCommitOlder = versionDiff.Arg("hash-older", "The hash of older one.").String()
	diffCommitNewer = versionDiff.Arg("hash-newer", "The hash of newer one.").Default("HEAD").String()

	cliVersionCommand = app.Command("version", "Command line interface version.")

	combineCode   = app.Command("combine", "Combine code that imported.")
	combineFile   = combineCode.Arg("filename", "File to be combined.").Required().String()
	combineOutput = combineCode.Arg("output", "Output file name for combined code.").String()
	listCommand   = app.Command("list", "List all api names.")
	listAPIFile   = listCommand.Arg("apifile", "Api file to list.").Default("pica.fun").String()
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
		apiRunner := pica.NewAPIRunnerFromFile(*runFileName, *runAPINames, *runDelay)
		err := apiRunner.Run()
		if err != nil {
			panic(err)
		}
		if *runOutput != "" {
			gen := pica.NewMarkdownDocGenerator(apiRunner, *runOutputTheme, *runOutput)
			gen.Get()
		}
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
	case cliVersionCommand.FullCommand():
		fmt.Printf("Commit:   %s\n", COMMIT)
		fmt.Printf("Builds:   %s\n", BuildAt)
		fmt.Printf("Golang:   %s\n", GOLANG)
		break
	case combineCode.FullCommand():
		code, err := funny.CombinedCode(*combineFile)
		if err != nil {
			panic(err)
		}
		fmt.Println(code)
	case listCommand.FullCommand():
		apiRunner := pica.NewAPIRunnerFromFile(*listAPIFile, *runAPINames, *runDelay)
		err := apiRunner.Parse()
		if err != nil {
			panic(err)
		}
		err = apiRunner.ParseAPIItems()
		if err != nil {
			panic(err)
		}
		for _, item := range apiRunner.APIItems {
			fmt.Printf("%s %s\n", item.Request.Name, item.Request.Description)
		}
		break
	default:
		kingpin.Usage()
	}

}
