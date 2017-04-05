package command

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdHostList(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdHostList(c))
	assert.Equal(t, "bar\nbaz.com\ngoo\n", writer.String())
}

func TestCmdHostListNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	assert.EqualError(t, CmdHostList(c), "You must specify a config file")
}

func TestCmdHostListUsage(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"abc"}))
	c := cli.NewContext(nil, set, nil)

	assert.EqualError(t, CmdHostList(c), "Usage: \"hostBuilder host list\"")
}

func TestCmdHostListBadConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)

	set.String("config", "/doesntexist", "doc")

	c := cli.NewContext(nil, set, nil)
	err := CmdHostList(c)
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}
