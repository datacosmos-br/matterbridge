package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/google/gops/agent"
	"github.com/matterbridge-org/matterbridge/bridge/config"
	"github.com/matterbridge-org/matterbridge/gateway"
	"github.com/matterbridge-org/matterbridge/gateway/bridgemap"
	"github.com/matterbridge-org/matterbridge/gateway/dify"
	"github.com/matterbridge-org/matterbridge/version"
	prefixed "github.com/matterbridge/logrus-prefixed-formatter"
	"github.com/sirupsen/logrus"
)

var (
	flagConfig  = flag.String("conf", "matterbridge.toml", "config file")
	flagDebug   = flag.Bool("debug", false, "enable debug")
	flagVersion = flag.Bool("version", false, "show version")
	flagGops    = flag.Bool("gops", false, "enable gops agent")
	flagNoColor = flag.Bool("nocolor", false, "disable colored logs")

	flagIdentityLink = flag.String("identity-link", "",
		"link a chat user to a corporate user in IdentityStorePath and exit (account:user_id=corp_user)")
)

func main() {
	flag.Parse()
	if *flagVersion {
		fmt.Printf("version: %s %s\n", version.Release, version.GitHash)
		return
	}

	rootLogger := setupLogger()
	logger := rootLogger.WithFields(logrus.Fields{"prefix": "main"})

	if *flagGops {
		if err := agent.Listen(agent.Options{}); err != nil {
			logger.Errorf("Failed to start gops agent: %#v", err)
		} else {
			defer agent.Close()
		}
	}

	logger.Printf("Running version %s %s", version.Release, version.GitHash)
	if strings.Contains(version.Release, "-dev") {
		logger.Println("WARNING: THIS IS A DEVELOPMENT VERSION. Things may break.")
	}

	cfg := config.NewConfig(rootLogger, *flagConfig)
	cfg.BridgeValues().General.Debug = *flagDebug

	if *flagIdentityLink != "" {
		err := linkIdentity(cfg.BridgeValues().General.IdentityStorePath, *flagIdentityLink)
		if err != nil {
			logger.Fatalf("Linking identity failed: %s", err)
		}

		logger.Printf("Linked identity %s", *flagIdentityLink)

		return
	}

	// if logging to a file, ensure it is closed when the program terminates
	// nolint:errcheck
	defer func() {
		if f, ok := rootLogger.Out.(*os.File); ok {
			f.Sync()
			f.Close()
		}
	}()

	r, err := gateway.NewRouter(rootLogger, cfg, bridgemap.FullMap)
	if err != nil {
		logger.Fatalf("Starting gateway failed: %s", err)
	}
	if err = r.Start(); err != nil {
		logger.Fatalf("Starting gateway failed: %s", err)
	}
	logger.Printf("Gateway(s) started successfully. Now relaying messages")
	// matterbridge relays until the process is killed and has no shutdown
	// path, so r.Close() is not called here: the sqlite stores are WAL
	// journaled and every write is a committed transaction.
	select {}
}

// errNoIdentityStore reports an identity link without IdentityStorePath.
var errNoIdentityStore = errors.New("IdentityStorePath is not set in the configuration")

// linkIdentity stores the identity link spec ("account:user_id=corp_user")
// in the identity database at path.
func linkIdentity(path, spec string) error {
	if path == "" {
		return errNoIdentityStore
	}

	link, err := dify.ParseLink(spec)
	if err != nil {
		return err
	}

	store, err := dify.OpenIdentityStore(path)
	if err != nil {
		return err
	}

	return errors.Join(store.Link(link), store.Close())
}

func setupLogger() *logrus.Logger {
	logger := &logrus.Logger{
		Out: os.Stdout,
		Formatter: &prefixed.TextFormatter{
			PrefixPadding: 13,
			DisableColors: *flagNoColor,
		},
		Level: logrus.InfoLevel,
	}
	if *flagDebug || os.Getenv("DEBUG") == "1" {
		logger.SetReportCaller(true)
		logger.Formatter = &prefixed.TextFormatter{
			PrefixPadding: 13,
			DisableColors: *flagNoColor,
			FullTimestamp: false,

			CallerFormatter: func(function, file string) string {
				return fmt.Sprintf(" [%s:%s]", function, file)
			},
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				sp := strings.SplitAfter(f.File, "/matterbridge/")
				filename := f.File
				if len(sp) > 1 {
					filename = sp[1]
				}
				s := strings.Split(f.Function, ".")
				funcName := s[len(s)-1]
				return funcName, fmt.Sprintf("%s:%d", filename, f.Line)
			},
		}

		logger.Level = logrus.DebugLevel
		logger.WithFields(logrus.Fields{"prefix": "main"}).Info("Enabling debug logging.")
	}
	return logger
}
