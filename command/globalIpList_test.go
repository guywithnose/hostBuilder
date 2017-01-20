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

func TestCmdGlobalIPList(t *testing.T) {
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
	assert.Nil(t, CmdGlobalIPList(c))

	expectedOutput := "abc 10.0.0.2\ndef 127.0.0.1\n"
	if writer.String() != expectedOutput {
		t.Fatalf("Output was %s, expected %s", writer.String(), expectedOutput)
	}

	assert.Nil(t, os.Remove(configFile.Name()))
}

func TestCmdGlobalIPListUsage(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"abc"}))
	c := cli.NewContext(nil, set, nil)

	assert.EqualError(t, CmdGlobalIPList(c), "Usage: \"hostBuilder globalIP list\"")
}

func TestCmdGlobalIPListNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	assert.EqualError(t, CmdGlobalIPList(c), "You must specify a config file")
}

func TestCmdGlobalIPListBadConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)

	set.String("config", "/doesntexist", "doc")

	c := cli.NewContext(nil, set, nil)
	err := CmdGlobalIPList(c)
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}
