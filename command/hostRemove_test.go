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

func TestCmdHostRemove(t *testing.T) {
	configFileName := setupBaseConfigFile(t)

	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"goo", "foop"}))

	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, CmdHostRemove(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	expectedHost := config.Host{Current: "ignore", Options: map[string]string{}}
	if !reflect.DeepEqual(modifiedConfigData.Hosts["goo"], expectedHost) {
		t.Fatalf("Host goo was %v, expected %v", modifiedConfigData.Hosts["goo"], expectedHost)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdHostRemoveUsage(t *testing.T) {
	set := flag.NewFlagSet("test", 0)

	c := cli.NewContext(nil, set, nil)
	assert.EqualError(t, CmdHostRemove(c), "Usage: \"hostBuilder host remove {hostName} {IPName}\"")
}

func TestCmdHostRemoveNoConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"goo", "foop"}))
	c := cli.NewContext(nil, set, nil)

	assert.EqualError(t, CmdHostRemove(c), "You must specify a config file")
}

func TestCmdHostRemoveBadHostName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)

	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"goop", "foop"}))

	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	assert.EqualError(t, CmdHostRemove(c), "Host goop does not exist")
}

func TestCmdHostRemoveBadIPName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)

	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"goo", "foo"}))

	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	assert.EqualError(t, CmdHostRemove(c), "IPName foo does not exist")
}

func TestCompleteHostRemoveHostName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteHostRemove(c)

	expectedOutput := "bar\nbaz.com\ngoo\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteHostRemoveBadHostName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"goop"}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteHostRemove(c)

	expectedOutput := ""
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteHostRemoveIPName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"goo"}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteHostRemove(c)

	expectedOutput := "foop:10.0.0.8\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteHostRemoveComplete(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"goo", "foop"}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteHostRemove(c)

	expectedOutput := ""
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteHostRemoveNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteHostRemove(c)

	expectedOutput := ""
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}
}
