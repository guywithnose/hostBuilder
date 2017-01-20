package command

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdGroupShow(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo"}))

	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdGroupShow(c))

	expectedOutput := "foo\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
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
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"food"}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	err := CmdGroupShow(c)
	assert.EqualError(t, err, "Group food does not exist")
	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteGroupShowGroupName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteGroupShow(c)

	expectedOutput := "foo\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteGroupShowComplete(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo"}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteGroupShow(c)

	expectedOutput := ""
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteGroupShowNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteGroupShow(c)

	expectedOutput := ""
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}
}
