package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	flag "github.com/spf13/pflag"
)

// initConfig initializes the app's configuration manager
// and loads disk and command line configuration values.
func initConfig() config {
	var cfg = config{}
	var koanf = koanf.New(".")
	// Parse Command Line Flags.
	// --config flag to specify the location of config file on fs
	flagSet := flag.NewFlagSet("config", flag.ContinueOnError)
	flagSet.Usage = func() {
		fmt.Println(flagSet.FlagUsages())
		os.Exit(0)
	}
	// Config Location.
	flagSet.StringSlice("config", []string{}, "Path to a config file to load. This can be specified multiple times and the config files will be merged in order")
	// Process flags.
	failOnReadConfigErr(flagSet.Parse(os.Args[1:]))
	// Read default config file. Won't throw the error yet.
	vErr := koanf.Load(file.Provider("config.toml"), toml.Parser())
	// Load the config files provided in the commandline if there are any.
	cFiles, _ := flagSet.GetStringSlice("config")
	for _, c := range cFiles {
		if err := koanf.Load(file.Provider(c), toml.Parser()); err != nil {
			log.Fatalf("error loading config file: %v", err)
		}
	}
	// If no default config is read and no additional config is supplied, exit.
	if vErr != nil {
		if len(cFiles) == 0 {
			log.Fatalf("no config was read: %v", vErr)
		}
	}
	// Read the configuration and load it to internal struct.
	failOnReadConfigErr(koanf.Unmarshal("", &cfg))
	return cfg
}

func failOnReadConfigErr(err error) {
	if err != nil {
		log.Fatalf("error reading config: %v.", err)
	}
}