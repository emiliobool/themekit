package commands

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func currentDir() string {
	dir, _ := os.Getwd()
	return dir
}

type HandlerMock struct {
	Prepare func(set *flag.FlagSet)
	Extract func(results map[string]interface{})
}

func (m HandlerMock) PrepareFlags(set *flag.FlagSet) {
	m.Prepare(set)
}

func (m HandlerMock) ExtractValues(results map[string]interface{}) {
	m.Extract(results)
}

func TestBasicArgumentParsing(t *testing.T) {
	arguments := []string{"--env", "production", "--dir", "/hello/world", "--verbose"}
	argParser := newArgParser("program", arguments)
	results := argParser.Parse()
	assert.Equal(t, true, results["verbose"])
	assert.Equal(t, "/hello/world", results["directory"])
	assert.Equal(t, "production", results["environment"])
}

func TestBasicArugmentParsingDefaults(t *testing.T) {
	arguments := []string{}
	argParser := newArgParser("program", arguments)
	results := argParser.Parse()
	assert.Equal(t, false, results["verbose"])
	assert.Equal(t, currentDir(), results["directory"])
	assert.Equal(t, DefaultEnvironment, results["environment"])
}

func TestParsingAdditionalArguments(t *testing.T) {
	var thing bool
	var include string
	mock := HandlerMock{
		Prepare: func(set *flag.FlagSet) {
			set.StringVar(&include, "include", "invalid", "")
			set.BoolVar(&thing, "thing", false, "")
		},
		Extract: func(results map[string]interface{}) {
			results["thing"] = thing
			results["include"] = include
		},
	}
	arguments := []string{"--thing", "--include", "important thing"}
	argParser := newArgParser("program", arguments)
	argParser.AddHandler(mock)
	results := argParser.Parse()

	assert.Equal(t, true, thing)
	assert.Equal(t, "important thing", include)
	assert.Equal(t, thing, results["thing"])
	assert.Equal(t, include, results["include"])
}

func TestAddedHandlersAreRunInOrder(t *testing.T) {
	var abra, kadabra string
	arguments := []string{"--abra", "magic!", "--kadabra", "poof!!"}
	argParser := newArgParser("program", arguments)
	argParser.AddHandler(HandlerMock{
		Prepare: func(set *flag.FlagSet) {
			set.StringVar(&abra, "abra", "fail 1", "message")
		},
		Extract: func(results map[string]interface{}) {
			results["answer"] = abra
		},
	})
	argParser.AddHandler(HandlerMock{
		Prepare: func(set *flag.FlagSet) {
			set.StringVar(&kadabra, "kadabra", "fail 2", "message")
		},
		Extract: func(results map[string]interface{}) {
			results["answer"] = kadabra
		},
	})
	results := argParser.Parse()
	assert.Equal(t, "poof!!", results["answer"])
}
