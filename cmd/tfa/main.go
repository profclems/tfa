package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/profclems/go-dotenv"
	"github.com/profclems/tfa/internal/build"
	"github.com/profclems/tfa/pkg/cmd/root"
	"github.com/profclems/tfa/utils/iomanip"
)

var configRetries = 0
var _configDir = ""

func main() {
	buildDate := build.Date
	version := build.Version

	iom := iomanip.InitIO()

	cfg, err := initConfig()
	if err != nil {
		fmt.Fprintln(iom.StdErr, err)
		os.Exit(1)
	}

	rootCmd := root.NewRootCmd(iom, cfg, version, buildDate)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(iom.StdErr, err)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() (*dotenv.DotEnv, error) {
	cfg := dotenv.Init(configFile())
	cfg.Separator = ":" // override separator since base32 allows '='
	// If a config file is found, read it in.
	if err := cfg.LoadConfig(); err != nil {
		if os.IsNotExist(err) && configRetries <= 2 {
			err = os.MkdirAll(configDir(), 0755)
			if err != nil {
				return nil, err
			}
			err = ioutil.WriteFile(configFile(), []byte(""), 0600)
			if err != nil {
				return nil, err
			}
			configRetries++
			return initConfig()
		}
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}
	return cfg, nil
}

func configDir() string {
	if _configDir != "" {
		return _configDir
	}
	usrConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if usrConfigHome == "" {
		usrConfigHome = os.Getenv("HOME")
		if usrConfigHome == "" {
			usrConfigHome, _ = homedir.Expand("~/.config")
		} else {
			usrConfigHome = filepath.Join(usrConfigHome, ".config")
		}
	}
	_configDir = filepath.Join(usrConfigHome, "tfa")
	return _configDir
}

func configFile() string {
	return filepath.Join(configDir(), ".tfa")
}
