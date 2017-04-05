package command

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdHostShow(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	err := set.Parse([]string{"goo"})
	assert.Nil(t, err)

	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdHostShow(c))

	assert.Equal(t, "1 Option:\n*foop => 10.0.0.8*\n", writer.String())
}

func TestCmdHostShowGlobalIP(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	err := set.Parse([]string{"baz.com"})
	assert.Nil(t, err)

	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdHostShow(c))

	assert.Equal(t, "1 Option:\nbazz => 10.0.0.7\nCurrent: Global IP baz => 10.0.0.4\n", writer.String())
}

func TestCmdHostShowGlobalUnknown(t *testing.T) {
	configFileName := setupInvalidConfigFile(t)
	defer removeFile(t, configFileName)
	set := flag.NewFlagSet("test", 0)
	err := set.Parse([]string{"unknown"})
	assert.Nil(t, err)

	set.String("config", configFileName, "doc")
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdHostShow(c))

	expectedOutput := "0 Options:\nCurrent: unknown (Warning: no associated IP please validate your config)\n"
	assert.Equal(t, expectedOutput, writer.String())
}

func TestCmdHostShowUsage(t *testing.T) {
	c := cli.NewContext(nil, flag.NewFlagSet("test", 0), nil)
	err := CmdHostShow(c)
	assert.EqualError(t, err, "Usage: \"hostBuilder host show {hostName}\"")
}

func TestCmdHostShowNoConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"goo"}))

	c := cli.NewContext(nil, set, nil)
	err := CmdHostShow(c)
	assert.EqualError(t, err, "You must specify a config file")
}

func TestCmdHostShowBadConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"goo"}))

	set.String("config", "/doesntexist", "doc")
	c := cli.NewContext(nil, set, nil)
	err := CmdHostShow(c)
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}

func TestCmdHostShowBadHostName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"bart"}))
	c := cli.NewContext(nil, set, nil)
	err := CmdHostShow(c)
	assert.EqualError(t, err, "Hostname bart does not exist")
}

func TestCompleteHostShowHostName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{}))
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteHostShow(c)

	assert.Equal(t, "bar\nbaz.com\ngoo\n", writer.String())
}

func TestCompleteHostShowComplete(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"goo"}))
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteHostShow(c)

	assert.Equal(t, "", writer.String())
}

func TestCompleteHostShowNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteHostShow(c)

	assert.Equal(t, "", writer.String())
}
