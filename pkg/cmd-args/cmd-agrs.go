package cmdArgs

import (
	"flag"
	"fmt"
	"os"
	stringsHelpers "string-helpers"

	lv "latest-version"

	"github.com/oleksii-honchar/blablo"
	c "github.com/oleksii-honchar/coteco"
)

var f = fmt.Sprintf
var logger *blablo.Logger

var ConfigPath string = "./nrp.yaml"
var LogLevel string = string(blablo.LevelInfo)

func isValidLogLevel(level string) bool {
	validLevels := []string{
		string(blablo.LevelDebug),
		string(blablo.LevelInfo),
		string(blablo.LevelWarn),
		string(blablo.LevelError),
	}

	return stringsHelpers.CheckIfStrInArray(level, validLevels)
}

func Init() bool {

	flag.StringVar(&ConfigPath, "config", "./nrp.yaml", "Specify 'config' path value")
	flag.StringVar(&LogLevel, "log-level", string(blablo.LevelInfo), "Specify 'log-level' value")
	showVersion1 := flag.Bool("v", false, "Show current version")
	showVersion2 := flag.Bool("version", false, "Show current version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println("Visit https://github.com/oleksii-honchar/nrp-cli for more details")
	}

	flag.Parse()

	if !isValidLogLevel(LogLevel) {
		fmt.Println("Invalid 'log-level' value. Please provide a valid log level.")
		return false
	}

	if *showVersion1 || *showVersion2 {
		fmt.Println(lv.LatestVersion)
		return false
	}

	logger = blablo.NewLogger("cmd-args", string(LogLevel))

	logger.Debug(c.WithGray247(f("cmd arg 'config' = %s", ConfigPath)))
	logger.Debug(c.WithGray247(f("cmd arg 'log-level' = %s", LogLevel)))

	return true
}
