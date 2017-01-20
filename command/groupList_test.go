package command

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdGroupList(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdGroupList(c))

	expectedOutput := "foo\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdGroupListUsage(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"abc"}))
	c := cli.NewContext(nil, set, nil)

	assert.EqualError(t, CmdGroupList(c), "Usage: \"hostBuilder group list\"")
}

func TestCmdGroupListNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	assert.EqualError(t, CmdGroupList(c), "You must specify a config file")
}

func TestCmdGroupListBadConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)

	set.String("config", "/doesntexist", "doc")

	c := cli.NewContext(nil, set, nil)
	err := CmdGroupList(c)
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}
