package command

import (
	"bytes"
	"flag"
	"os"
	"reflect"
	"testing"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdGroupAdd(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo", "bar"}))

	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, CmdGroupAdd(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	expectedGroup := []string{"baz.com", "goo", "bar"}
	if !reflect.DeepEqual(modifiedConfigData.Groups["foo"], expectedGroup) {
		t.Fatalf("Group foo was %v, expected %v", modifiedConfigData.Groups["foo"], expectedGroup)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdGroupAddUsage(t *testing.T) {
	c := cli.NewContext(nil, flag.NewFlagSet("test", 0), nil)
	err := CmdGroupAdd(c)
	assert.EqualError(t, err, "Usage: \"hostBuilder group add {groupName} {hostName}\"")
}

func TestCmdGroupAddNoConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo", "bar"}))

	c := cli.NewContext(nil, set, nil)
	err := CmdGroupAdd(c)
	assert.EqualError(t, err, "You must specify a config file")
}

func TestCmdGroupAddBadConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo", "bar"}))

	set.String("config", "/doesntexist", "doc")
	c := cli.NewContext(nil, set, nil)
	err := CmdGroupAdd(c)
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}

func TestCmdGroupAddBadGroupName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"food", "bar"}))
	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, CmdGroupAdd(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	expectedGroup := []string{"bar"}
	if !reflect.DeepEqual(modifiedConfigData.Groups["food"], expectedGroup) {
		t.Fatalf("Group food was %v, expected %v", modifiedConfigData.Groups["food"], expectedGroup)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdGroupAddBadHostName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo", "bart"}))
	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	err := CmdGroupAdd(c)
	assert.EqualError(t, err, "Hostname bart does not exist")
	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdGroupAddAlreadyExists(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo", "goo"}))

	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.ErrWriter = writer
	c := cli.NewContext(app, set, nil)
	err := CmdGroupAdd(c)
	assert.EqualError(t, err, "Group foo already contains goo")

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteGroupAddGroupName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteGroupAdd(c)

	expectedOutput := "foo\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteGroupAddHostName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo"}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteGroupAdd(c)

	expectedOutput := "bar\nbaz.com\ngoo\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteGroupAddComplete(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo", "bar"}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteGroupAdd(c)

	expectedOutput := "\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteGroupAddNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteGroupAdd(c)

	expectedOutput := ""
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}
}
