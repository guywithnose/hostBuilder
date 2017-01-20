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

func TestCmdGroupSet(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo", "baz"}))

	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, CmdGroupSet(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	if modifiedConfigData.Hosts["goo"].Current != "baz" {
		t.Fatal("goo was not set to baz")
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdGroupSetUsage(t *testing.T) {
	c := cli.NewContext(nil, flag.NewFlagSet("test", 0), nil)
	err := CmdGroupSet(c)
	assert.EqualError(t, err, "Usage: \"hostBuilder group set {groupName} {globalIPName}\"")
}

func TestCmdGroupSetNoConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo", "baz"}))

	c := cli.NewContext(nil, set, nil)
	err := CmdGroupSet(c)
	assert.EqualError(t, err, "You must specify a config file")
}

func TestCmdGroupSetBadConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo", "baz"}))

	set.String("config", "/doesntexist", "doc")
	c := cli.NewContext(nil, set, nil)
	err := CmdGroupSet(c)
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}

func TestCmdGroupSetBadGroupName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"food", "baz"}))
	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	err := CmdGroupSet(c)
	assert.EqualError(t, err, "Group food does not exist")
	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdGroupSetBadHostName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo", "barz"}))
	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	err := CmdGroupSet(c)
	assert.EqualError(t, err, "Global IP barz does not exist")
	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteGroupSetGroupName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteGroupSet(c)

	expectedOutput := "foo\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteGroupSetGlobalIPs(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo"}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteGroupSet(c)

	expectedOutput := "baz:10.0.0.4\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteGroupSetComplete(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo", "baz"}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteGroupSet(c)

	expectedOutput := ""
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteGroupSetNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteGroupSet(c)

	expectedOutput := ""
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}
}
