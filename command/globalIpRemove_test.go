package command

import (
	"flag"
	"io/ioutil"
	"testing"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdGlobalIPRemove(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	defer removeFile(t, configFile.Name())

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
	assert.Equal(t, 1, len(modifiedConfigData.GlobalIPs), "Global IPs length was not 1")
	_, ok := modifiedConfigData.GlobalIPs["abc"]
	assert.Equal(t, false, ok, "\"abc\" Global IP was set after removal")
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
	defer removeFile(t, configFile.Name())

	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo"}))

	configData := &config.HostsConfig{GlobalIPs: map[string]string{"def": "127.0.0.1", "abc": "10.0.0.2"}}
	err = config.WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)
	set.String("config", configFile.Name(), "doc")
	c := cli.NewContext(nil, set, nil)
	err = CmdGlobalIPRemove(c)
	assert.EqualError(t, err, "GlobalIP foo does not exist")
}

func TestCompleteGlobalIPRemove(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	defer removeFile(t, configFile.Name())

	set := flag.NewFlagSet("test", 0)

	configData := &config.HostsConfig{GlobalIPs: map[string]string{"def": "127.0.0.1", "abc": "10.0.0.2"}}
	err = config.WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)
	set.String("config", configFile.Name(), "doc")
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteGlobalIPRemove(c)
	assert.Equal(t, "abc:10.0.0.2\ndef:127.0.0.1\n", writer.String())
}

func TestCompleteGlobalIPRemoveNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	app, writer := appWithWriter()
	c := cli.NewContext(app, set, nil)
	CompleteGlobalIPRemove(c)
	assert.Equal(t, "", writer.String())
}
