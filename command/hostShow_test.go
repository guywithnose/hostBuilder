package command

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdHostShow(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	err := set.Parse([]string{"goo"})
	assert.Nil(t, err)

	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdHostShow(c))

	expectedOutput := "1 Option:\n*foop => 10.0.0.8*\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdHostShowGlobalIP(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	err := set.Parse([]string{"baz.com"})
	assert.Nil(t, err)

	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdHostShow(c))

	expectedOutput := "1 Option:\nbazz => 10.0.0.7\nCurrent: Global IP baz => 10.0.0.4\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdHostShowGlobalUnknown(t *testing.T) {
	configFileName := setupInvalidConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	err := set.Parse([]string{"unknown"})
	assert.Nil(t, err)

	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdHostShow(c))

	expectedOutput := "0 Options:\nCurrent: unknown (Warning: no associated IP please validate your config)\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
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
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"bart"}))
	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	err := CmdHostShow(c)
	assert.EqualError(t, err, "Hostname bart does not exist")
	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteHostShowHostName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteHostShow(c)

	expectedOutput := "bar\nbaz.com\ngoo\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteHostShowComplete(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"goo"}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteHostShow(c)

	expectedOutput := ""
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteHostShowNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteHostShow(c)

	expectedOutput := ""
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}
}
