package command

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdGlobalIPRemove(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)

	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"abc"}))

	configData := &config.HostsConfig{GlobalIPs: map[string]string{"def": "127.0.0.1", "abc": "10.0.0.2"}}
	err = config.WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)
	set.String("config", configFile.Name(), "doc")
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, CmdGlobalIPRemove(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFile.Name())
	assert.Nil(t, err)
	if len(modifiedConfigData.GlobalIPs) != 1 {
		t.Fatalf("Global IPs length was %d, expected 1", len(modifiedConfigData.GlobalIPs))
	}

	if _, ok := modifiedConfigData.GlobalIPs["abc"]; ok {
		t.Fatal("\"abc\" Global IP was set after removal")
	}

	assert.Nil(t, os.Remove(configFile.Name()))
}

func TestCmdGlobalIPRemoveUsage(t *testing.T) {
	set := flag.NewFlagSet("test", 0)

	c := cli.NewContext(nil, set, nil)
	assert.EqualError(t, CmdGlobalIPRemove(c), "Usage: \"hostBuilder globalIP remove {Name}\"")
}

func TestCmdGlobalIPRemoveNoConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"abc"}))
	c := cli.NewContext(nil, set, nil)

	assert.EqualError(t, CmdGlobalIPRemove(c), "You must specify a config file")
}

func TestCmdGlobalIPRemoveNonExistantName(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)

	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo"}))

	configData := &config.HostsConfig{GlobalIPs: map[string]string{"def": "127.0.0.1", "abc": "10.0.0.2"}}
	err = config.WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)
	set.String("config", configFile.Name(), "doc")
	c := cli.NewContext(nil, set, nil)
	err = CmdGlobalIPRemove(c)
	assert.EqualError(t, err, "GlobalIP foo does not exist")
	assert.Nil(t, os.Remove(configFile.Name()))
}

func TestCompleteGlobalIPRemove(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)

	set := flag.NewFlagSet("test", 0)

	configData := &config.HostsConfig{GlobalIPs: map[string]string{"def": "127.0.0.1", "abc": "10.0.0.2"}}
	err = config.WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)
	set.String("config", configFile.Name(), "doc")
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteGlobalIPRemove(c)
	expectedOutput := "abc:10.0.0.2\ndef:127.0.0.1\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFile.Name()))
}

func TestCompleteGlobalIPRemoveNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)

	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer
	c := cli.NewContext(app, set, nil)
	CompleteGlobalIPRemove(c)
	expectedOutput := ""
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}
}
