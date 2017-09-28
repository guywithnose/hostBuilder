package command

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdGroupShow(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"foo"}))

	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdGroupShow(c))

	assert.Equal(t, "baz.com\ngoo\n", writer.String())
}

func TestCmdGroupShowUsage(t *testing.T) {
	c := cli.NewContext(nil, flag.NewFlagSet("test", 0), nil)
	err := CmdGroupShow(c)
	assert.EqualError(t, err, "Usage: \"hostBuilder group show {groupName}\"")
}

func TestCmdGroupShowNoConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo"}))

	c := cli.NewContext(nil, set, nil)
	err := CmdGroupShow(c)
	assert.EqualError(t, err, "You must specify a config file")
}

func TestCmdGroupShowBadConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo"}))

	set.String("config", "/doesntexist", "doc")
	c := cli.NewContext(nil, set, nil)
	err := CmdGroupShow(c)
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}

func TestCmdGroupShowBadGroupName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"food"}))
	app := cli.NewApp()
	c := cli.NewContext(app, set, nil)
	err := CmdGroupShow(c)
	assert.EqualError(t, err, "Group food does not exist")
}

func TestCompleteGroupShowGroupName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteGroupShow(c)

	assert.Equal(t, "foo\n", writer.String())
}

func TestCompleteGroupShowComplete(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"foo"}))
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteGroupShow(c)

	assert.Equal(t, "", writer.String())
}

func TestCompleteGroupShowNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteGroupShow(c)

	assert.Equal(t, "", writer.String())
}
