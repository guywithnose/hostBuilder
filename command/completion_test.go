package command_test

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/guywithnose/hostBuilder/command"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestRootCompletion(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	app, writer, _ := appWithTestWriters()
	app.Commands = append(command.Commands, cli.Command{Hidden: true, Name: "don't show"})
	app.Flags = command.GlobalFlags
	command.RootCompletion(cli.NewContext(app, set, nil))
	assert.Equal(
		t,
		[]string{
			"createConfig:Create a config file from an existing hosts file",
			"build:Builds your host file",
			"globalIP:Add things to the configuration",
			"host:Modify hosts",
			"group:Modify groups",
			"aws:Add information from AWS to the configuration",
			"--config",
			"",
		},
		strings.Split(writer.String(), "\n"),
	)
}

func TestRootCompletionConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	app, writer, _ := appWithTestWriters()
	os.Args = []string{os.Args[0], "--config", "--completion"}
	command.RootCompletion(cli.NewContext(app, set, nil))
	assert.Equal(t, "fileCompletion\n", writer.String())
}

func appWithTestWriters() (*cli.App, *bytes.Buffer, *bytes.Buffer) {
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	errWriter := new(bytes.Buffer)
	app.Writer = writer
	app.ErrWriter = errWriter
	return app, writer, errWriter
}
