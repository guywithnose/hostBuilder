package command

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdGroupList(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdGroupList(c))
	assert.Equal(t, "foo\n", writer.String())
}

func TestCmdGroupListUsage(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"def"}))
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
