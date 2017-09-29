package command

import (
	"bytes"
	"flag"
	"io/ioutil"
	"testing"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdGlobalIPAdd(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	defer removeFile(t, configFile.Name())

	_, err = runGlobalIPAddCommand(configFile.Name(), map[string]string{"def": "127.0.0.1"})
	assert.Nil(t, err)

	modifiedConfigData, err := config.LoadConfigFromFile(configFile.Name())
	assert.Nil(t, err)

	assert.Equal(t, 2, len(modifiedConfigData.GlobalIPs))
	assert.Equal(t, "10.0.0.2", modifiedConfigData.GlobalIPs["abc"])
}

func TestCmdGlobalIPHostName(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	defer removeFile(t, configFile.Name())

	set := flag.NewFlagSet("test", 0)
	err = set.Parse([]string{"abc", "localhost4"})
	assert.Nil(t, err)

	configData := &config.HostsConfig{GlobalIPs: map[string]string{"def": "127.0.0.1"}}
	assert.Nil(t, config.WriteConfig(configFile.Name(), configData))

	set.String("config", configFile.Name(), "doc")
	app := cli.NewApp()
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdGlobalIPAdd(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFile.Name())
	assert.Nil(t, err)

	assert.Equal(t, 2, len(modifiedConfigData.GlobalIPs))
	assert.Equal(t, "127.0.0.1", modifiedConfigData.GlobalIPs["abc"])
}

func TestCmdGlobalIPAddUsage(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"abc"}))
	c := cli.NewContext(nil, set, nil)
	err := CmdGlobalIPAdd(c)
	assert.EqualError(t, err, "Usage: \"hostBuilder globalIP add {Name} {address}\"")
}

func TestCmdGlobalIPAddInvalidIP(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	abcIP := "10.0.0.256"
	assert.Nil(t, set.Parse([]string{"abc", abcIP}))

	c := cli.NewContext(nil, set, nil)
	err := CmdGlobalIPAdd(c)
	assert.EqualError(t, err, "Unable to resolve 10.0.0.256")
}

func TestCmdGlobalIPAddNoConfigFile(t *testing.T) {
	set, err := getValidAddArgSet()
	assert.Nil(t, err)

	c := cli.NewContext(nil, set, nil)
	err = CmdGlobalIPAdd(c)
	assert.EqualError(t, err, "You must specify a config file")
}

func TestCmdGlobalIPAddBadConfigFile(t *testing.T) {
	set, err := getValidAddArgSet()
	assert.Nil(t, err)

	set.String("config", "/doesntexist", "doc")

	c := cli.NewContext(nil, set, nil)
	err = CmdGlobalIPAdd(c)
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}

func TestCmdGlobalIPAddNoOverride(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	defer removeFile(t, configFile.Name())
	_, err = runGlobalIPAddCommand(configFile.Name(), map[string]string{"abc": "127.0.0.1"})
	assert.EqualError(t, err, "Global IP abc already exists")
}

func TestCmdGlobalIPAddOverride(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	defer removeFile(t, configFile.Name())
	set, err := getValidAddArgSet()
	set.Bool("force", true, "")
	assert.Nil(t, err)

	configData := &config.HostsConfig{GlobalIPs: map[string]string{"abc": "127.0.0.1"}}
	err = config.WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)
	set.String("config", configFile.Name(), "doc")
	app, writer := appWithErrWriter()
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdGlobalIPAdd(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFile.Name())
	assert.Nil(t, err)

	assert.Equal(t, 1, len(modifiedConfigData.GlobalIPs))
	assert.Equal(t, "Warning: Overwriting abc (127.0.0.1 => 10.0.0.2)", writer.String())
}

func TestCompleteGlobalIPAdd(t *testing.T) {
	app, writer := appWithWriter()
	app.Commands = []cli.Command{
		{
			Name: "add",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "force",
				},
			},
		},
	}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)
	CompleteGlobalIPAdd(c)

	assert.Equal(t, "--force\n", writer.String())
}

func TestCompleteGlobalIPAddForce(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	defer removeFile(t, configFile.Name())
	configData := &config.HostsConfig{GlobalIPs: map[string]string{"abc": "127.0.0.1"}}
	err = config.WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)
	set := flag.NewFlagSet("test", 0)
	set.String("config", configFile.Name(), "doc")
	set.Bool("force", true, "")
	app, writer := appWithWriter()
	app.Commands = []cli.Command{{Name: "add"}}
	c := cli.NewContext(app, set, nil)
	CompleteGlobalIPAdd(c)

	assert.Equal(t, "abc\n", writer.String())
}

func TestCompleteGlobalIPAddIPs(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	defer removeFile(t, configFile.Name())
	configData := &config.HostsConfig{GlobalIPs: map[string]string{"abc": "127.0.0.1"}}
	err = config.WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)
	set := flag.NewFlagSet("test", 0)
	set.String("config", configFile.Name(), "doc")
	set.Bool("force", true, "")
	err = set.Parse([]string{"abc"})
	assert.Nil(t, err)
	app, writer := appWithWriter()
	app.Commands = []cli.Command{{Name: "add"}}
	c := cli.NewContext(app, set, nil)
	CompleteGlobalIPAdd(c)

	assert.Equal(t, "127.0.0.1\n", writer.String())
}

func getValidAddArgSet() (*flag.FlagSet, error) {
	set := flag.NewFlagSet("test", 0)
	abcIP := "10.0.0.2"
	err := set.Parse([]string{"abc", abcIP})
	return set, err
}

func runGlobalIPAddCommand(configFileName string, globalIPs map[string]string) (*bytes.Buffer, error) {
	set, err := getValidAddArgSet()
	if err != nil {
		return nil, err
	}

	configData := &config.HostsConfig{GlobalIPs: globalIPs}
	err = config.WriteConfig(configFileName, configData)
	if err != nil {
		return nil, err
	}

	set.String("config", configFileName, "doc")
	app, writer := appWithErrWriter()
	c := cli.NewContext(app, set, nil)
	err = CmdGlobalIPAdd(c)
	if err != nil {
		return nil, err
	}

	return writer, nil
}
