package commands

import (
	"flag"
	"os"
)

const DefaultEnvironment string = "development"

type flagHandler interface {
	PrepareFlags(set *flag.FlagSet)
	ExtractValues(results map[string]interface{})
}

type defaultCommandLineOptions struct {
	env     string
	dir     string
	verbose bool
}

type argParser struct {
	options            defaultCommandLineOptions
	flagSet            *flag.FlagSet
	args               []string
	additionalHandlers []flagHandler
}

func newArgParser(name string, args []string) *argParser {
	set := flag.NewFlagSet(name, flag.ExitOnError)
	return &argParser{flagSet: set, additionalHandlers: []flagHandler{}, args: args}
}

func (a *argParser) AddHandler(h flagHandler) {
	a.additionalHandlers = append(a.additionalHandlers, h)
}

func (a *argParser) Parse() map[string]interface{} {
	a.setupDefaultArgs()
	a.setupAdditionalArgs()
	a.flagSet.Parse(a.args)

	results := map[string]interface{}{}
	a.extractDefaultValues(results)
	a.extractAddtionalValues(results)
	return results
}

func (a *argParser) setupDefaultArgs() {
	currentDir, _ := os.Getwd()
	a.flagSet.StringVar(&a.options.dir, "dir", currentDir, "directory in which the config.yml file is located")
	a.flagSet.StringVar(&a.options.env, "env", DefaultEnvironment, "environment under which configuration to run as")
	a.flagSet.BoolVar(&a.options.verbose, "verbose", false, "run with verbose logging")
}

func (a *argParser) setupAdditionalArgs() {
	for _, handler := range a.additionalHandlers {
		handler.PrepareFlags(a.flagSet)
	}
}

func (a *argParser) extractDefaultValues(results map[string]interface{}) {
	results["directory"] = a.options.dir
	results["environment"] = a.options.env
	results["verbose"] = a.options.verbose
}

func (a *argParser) extractAddtionalValues(results map[string]interface{}) {
	for _, handler := range a.additionalHandlers {
		handler.ExtractValues(results)
	}
}
