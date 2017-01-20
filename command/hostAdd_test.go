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

func TestCmdHostAdd(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	hooIP := "10.0.0.2"
	err := set.Parse([]string{"bar", hooIP, "hoo"})
	assert.Nil(t, err)

	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, CmdHostAdd(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	expectedOptions := map[string]string{"hoo": hooIP}
	if !reflect.DeepEqual(modifiedConfigData.Hosts["bar"].Options, expectedOptions) {
		t.Fatalf("File was %v, expected %v", modifiedConfigData.Hosts["bar"].Options, expectedOptions)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdHostAddHostname(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	err := set.Parse([]string{"bar", "localhost4", "hoo"})
	assert.Nil(t, err)

	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, CmdHostAdd(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	expectedOptions := map[string]string{"hoo": "127.0.0.1"}
	if !reflect.DeepEqual(modifiedConfigData.Hosts["bar"].Options, expectedOptions) {
		t.Fatalf("File was %v, expected %v", modifiedConfigData.Hosts["bar"].Options, expectedOptions)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdHostAddOverwriteFails(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	hooIP := "10.0.0.2"
	err := set.Parse([]string{"goo", hooIP, "foop"})
	assert.Nil(t, err)

	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.ErrWriter = writer
	c := cli.NewContext(app, set, nil)
	err = CmdHostAdd(c)
	assert.EqualError(t, err, "IP goo already exists")

	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdHostAddOverwriteForced(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	hooIP := "10.0.0.2"
	err := set.Parse([]string{"goo", hooIP, "foop"})
	assert.Nil(t, err)

	set.String("config", configFileName, "doc")
	set.Bool("force", true, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.ErrWriter = writer
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdHostAdd(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	expectedOutput := "Warning: Overwriting foop (10.0.0.8 => 10.0.0.2)"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	expectedOptions := map[string]string{"foop": hooIP}
	if !reflect.DeepEqual(modifiedConfigData.Hosts["goo"].Options, expectedOptions) {
		t.Fatalf("File was %v, expected %v", modifiedConfigData.Hosts["bar"].Options, expectedOptions)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCmdHostAddBadIP(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	err := set.Parse([]string{"bar", "10.0.0.256", "hoo"})
	assert.Nil(t, err)

	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	err = CmdHostAdd(c)
	assert.EqualError(t, err, "Unable to resolve 10.0.0.256")
}

func TestCmdHostAddUsage(t *testing.T) {
	c := cli.NewContext(nil, flag.NewFlagSet("test", 0), nil)
	err := CmdHostAdd(c)
	assert.EqualError(t, err, "Usage: \"hostBuilder host add {hostName} {address} {IPName}\"")
}

func TestCmdHostAddNoConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"bar", "10.0.0.2", "hoo"}))

	c := cli.NewContext(nil, set, nil)
	err := CmdHostAdd(c)
	assert.EqualError(t, err, "You must specify a config file")
}

func TestCmdHostAddBadConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"bar", "10.0.0.2", "hoo"}))

	set.String("config", "/doesntexist", "doc")
	c := cli.NewContext(nil, set, nil)
	err := CmdHostAdd(c)
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}

func TestCmdHostAddNewHostName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	hooIP := "10.0.0.2"
	assert.Nil(t, set.Parse([]string{"barz", hooIP, "hoo"}))
	set.String("config", configFileName, "doc")
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, CmdHostAdd(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	expectedHost := config.Host{Current: "hoo", Options: map[string]string{"hoo": hooIP}}
	if !reflect.DeepEqual(modifiedConfigData.Hosts["barz"], expectedHost) {
		t.Fatalf("Host barz was %v, expected %v", modifiedConfigData.Hosts["barz"], expectedHost)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteHostAddHostName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	app.Commands = []cli.Command{{Name: "add"}}
	c := cli.NewContext(app, set, nil)
	CompleteHostAdd(c)

	expectedOutput := "bar\nbaz.com\ngoo\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteHostAddOptionName(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"bar"}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	app.Commands = []cli.Command{{Name: "add"}}
	c := cli.NewContext(app, set, nil)
	CompleteHostAdd(c)

	expectedOutput := "10.0.0.4:baz\n10.0.0.7:bazz\n10.0.0.8:foop\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteHostAddIP(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"bar", "127.0.0.1"}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	app.Commands = []cli.Command{{Name: "add"}}
	c := cli.NewContext(app, set, nil)
	CompleteHostAdd(c)

	expectedOutput := "bazz\nfoop\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteHostAddComplete(t *testing.T) {
	configFileName := setupBaseConfigFile(t)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"bar", "127.0.0.1", "hoo"}))
	set.String("config", configFileName, "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	app.Commands = []cli.Command{{Name: "add"}}
	c := cli.NewContext(app, set, nil)
	CompleteHostAdd(c)

	expectedOutput := ""
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFileName))
}

func TestCompleteHostAddNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	app.Commands = []cli.Command{{Name: "add"}}
	c := cli.NewContext(app, set, nil)
	CompleteHostAdd(c)

	expectedOutput := ""
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}
}
