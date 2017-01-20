package command

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdHostSet(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"baz.com", "bazz"}))

	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, CmdHostSet(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	if modifiedConfigData.Hosts["baz.com"].Current != "bazz" {
		t.Fatal("baz.com was not set to baz")
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdHostSetUsage(t *testing.T) {
	c := cli.NewContext(nil, flag.NewFlagSet("test", 0), nil)
	err := CmdHostSet(c)
	assert.EqualError(t, err, "Usage: \"hostBuilder host set {hostName} {IPName}\"")
}

func TestCmdHostSetNoConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"baz.com", "bazz"}))

	c := cli.NewContext(nil, set, nil)
	err := CmdHostSet(c)
	assert.EqualError(t, err, "You must specify a config file")
}

func TestCmdHostSetBadConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"baz.com", "bazz"}))

	set.String("config", "/doesntexist", "doc")
	c := cli.NewContext(nil, set, nil)
	err := CmdHostSet(c)
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}

func TestCmdHostSetBadHostName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"food", "baz"}))
	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	err := CmdHostSet(c)
	assert.EqualError(t, err, "HostName food does not exist")
	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdHostSetBadIPName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"baz.com", "bar"}))
	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	err := CmdHostSet(c)
	assert.EqualError(t, err, "IPName bar does not exist")
	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteHostSetHostName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteHostSet(c)

	expectedOutput := "bar\nbaz.com\ngoo\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteHostSetIPName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"baz.com"}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteHostSet(c)

	expectedOutput := "bazz:10.0.0.7\nbaz:10.0.0.4\nignore:\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteHostBadHostname(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"bar.com"}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteHostSet(c)

	expectedOutput := ""
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteHostSetComplete(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"baz.com", "baz"}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteHostSet(c)

	expectedOutput := ""
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteHostSetNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteHostSet(c)

	expectedOutput := ""
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}
}
