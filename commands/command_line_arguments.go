package commands

import (
	"flag"
	"github.com/Shopify/themekit"
	"os"
)

const DefaultEnvironment string = themekit.DefaultEnvironment
const DefaultBucketSize int = themekit.DefaultBucketSize
const DefaultRefillRate int = themekit.DefaultRefillRate

const (
	directoryKey    string = "directory"
	environmentKey         = "environment"
	environmentsKey        = "environments"
	themeClientKey         = "themeClient"
)

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
	a.loadThemeClient(results)
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

func (a *argParser) loadConfigurationValues(results map[string]interface{}) {
	directory, _ := results[directoryKey].(string)
	environment, _ := results[environmentKey].(string)
	results[environmentsKey] = loadEnvironments(directory)
	results[themeClientKey] = loadThemeClient(directory, environment)
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

type watchCommandHandler struct {
	allEnvs bool
	notify  string
}

func (w *watchCommandHandler) PrepareFlags(set *flag.FlagSet) {
	set.BoolVar(&w.allEnvs, "allenvs", false, "start watchers for all environments")
	set.StringVar(&w.notify, "notify", "", "file to touch when workers have gone idle")
}

func (w *watchCommandHandler) ExtractValues(results map[string]interface{}) {
	results["allenvs"] = w.allEnvs
	results["notify"] = w.notify
}

type configureCommandHandler struct {
	accessToken string
	domain      string
	env         string
	bucketSize  int
	refillRate  int
}

func (c *configureCommandHandler) PrepareFlags(set *flag.FlagSet) {
	set.StringVar(&c.accessToken, "access_token", "", "accessToken (or password) to make successful API calls")
	set.StringVar(&c.domain, "domain", "", "your shopify domain")
	set.StringVar(&c.env, "env", DefaultEnvironment, "environment for this configuration")
	set.IntVar(&c.bucketSize, "bucket_size", DefaultBucketSize, "leaky bucket capacity")
	set.IntVar(&c.refillRate, "refill_rate", DefaultRefillRate, "leaky bucket refill rate / second")
}

func (c *configureCommandHandler) ExtractValues(results map[string]interface{}) {
	results["accessToken"] = c.accessToken
	results["domain"] = c.domain
	results["environment"] = c.env
	results["bucketSize"] = c.bucketSize
	results["refillRate"] = c.refillRate
}

type bootstrapCommandHandler struct {
	prefix  string
	setId   bool
	version string
}

func (b *bootstrapCommandHandler) PrepareFlags(set *flag.FlagSet) {
	set.StringVar(&b.prefix, "prefix", "", "prefix to the Timber theme being created")
	set.StringVar(&b.prefix, "version", "latest", "version of Shopify Timber to use")
}
