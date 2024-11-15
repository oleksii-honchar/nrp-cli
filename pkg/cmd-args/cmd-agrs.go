package cmdArgs

import (
	"flag"
	"fmt"
	"os"
	stringsHelpers "string-helpers"

	lv "latest-version"

	cd "config-defaults"

	"github.com/oleksii-honchar/blablo"
	c "github.com/oleksii-honchar/coteco"
)

var f = fmt.Sprintf
var logger *blablo.Logger

var ConfigPath string = "./nrp.yaml"
var LogLevel string = string(blablo.LevelInfo)
var DefaultsMode string = cd.DefaultsProdMode
var CertbotWait bool = false
var CheckAndUpdatePublicIp bool = false
var Force bool = false

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
	flag.StringVar(&ConfigPath, "config", ConfigPath, "Specify 'config' path value: './nrp.yaml'(default)")
	flag.StringVar(&LogLevel, "log-level", LogLevel, "Specify 'log-level' value: info(default)|error|warn|debug")
	flag.StringVar(&DefaultsMode, "defaults-mode", DefaultsMode, "Specify 'defaults-mode' value : prod(default)|dev. When 'dev' defaults used - nginx cmds has a form of docker cmds. When 'prod' default used - nginx cmds assume nrp-cli executed in conatiner and nginx available directly.")
	flag.BoolVar(&CertbotWait, "certbot-wait", false, "When making certbot call, wait 5min right before the request to check nginx")
	flag.BoolVar(&CheckAndUpdatePublicIp, "check-and-update-public-ip", false, "Special cli mode when only public-ip registrant CDK will be executed in order to update public IP and domains A records. Should be used in cron tasks automatically.")
	flag.BoolVar(&Force, "force", false, "forcing 'check-and-update-public-ip' even if public IP updated")

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

	logger = blablo.NewLogger("cmd-args", string(LogLevel), false)

	logger.Debug(c.WithGray247(f("cmd arg 'config' = %s", ConfigPath)))
	logger.Debug(c.WithGray247(f("cmd arg 'log-level' = %s", LogLevel)))
	logger.Debug(c.WithGray247(f("cmd arg 'defaults-mode' = %s", DefaultsMode)))
	logger.Debug(c.WithGray247(f("cmd arg 'certbot-wait' = %v", CertbotWait)))
	logger.Debug(c.WithGray247(f("cmd arg 'check-and-update-public-ip' = %v", CheckAndUpdatePublicIp)))
	logger.Debug(c.WithGray247(f("cmd arg 'force' = %v", Force)))

	return true
}
