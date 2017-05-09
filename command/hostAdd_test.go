package command

import (
	"flag"
	"testing"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdHostAdd(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	hooIP := "10.0.0.2"
	err := set.Parse([]string{"bar", hooIP, "hoo"})
	assert.Nil(t, err)

	app, _ := appWithErrWriter()
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdHostAdd(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	assert.Equal(t, map[string]string{"hoo": hooIP}, modifiedConfigData.Hosts["bar"].Options)
}

func TestCmdHostAddHostname(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	err := set.Parse([]string{"bar", "localhost4", "hoo"})
	assert.Nil(t, err)

	app, _ := appWithErrWriter()
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdHostAdd(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	assert.Equal(t, map[string]string{"hoo": "127.0.0.1"}, modifiedConfigData.Hosts["bar"].Options)
}

func TestCmdHostAddOverwriteFails(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	hooIP := "10.0.0.2"
	err := set.Parse([]string{"goo", hooIP, "foop"})
	assert.Nil(t, err)

	app := cli.NewApp()
	c := cli.NewContext(app, set, nil)
	err = CmdHostAdd(c)
	assert.EqualError(t, err, "IP goo already exists")
}

func TestCmdHostAddOverwriteForced(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	hooIP := "10.0.0.2"
	err := set.Parse([]string{"goo", hooIP, "foop"})
	assert.Nil(t, err)

	set.Bool("force", true, "doc")
	app, writer := appWithErrWriter()
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdHostAdd(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	assert.Equal(t, "Warning: Overwriting foop (10.0.0.8 => 10.0.0.2)", writer.String())
	assert.Equal(t, map[string]string{"foop": hooIP}, modifiedConfigData.Hosts["goo"].Options)
}

func TestCmdHostAddBadIP(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	err := set.Parse([]string{"bar", "10.0.0.256", "hoo"})
	assert.Nil(t, err)

	app, _ := appWithErrWriter()
	c := cli.NewContext(app, set, nil)
	err = CmdHostAdd(c)
	assert.EqualError(t, err, "Unable to resolve 10.0.0.256")
}

func TestCmdHostAddUsage(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	app, _ := appWithErrWriter()
	c := cli.NewContext(app, set, nil)
	err := CmdHostAdd(c)
	assert.EqualError(t, err, "Usage: \"hostBuilder host add {hostName} ({address} {IPName}|{globalIpName})\"")
}

func TestCmdHostAddNoConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"bar", "10.0.0.2", "hoo"}))

	app, _ := appWithErrWriter()
	c := cli.NewContext(app, set, nil)
	err := CmdHostAdd(c)
	assert.EqualError(t, err, "You must specify a config file")
}

func TestCmdHostAddBadConfigFile(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"bar", "10.0.0.2", "hoo"}))

	set.String("config", "/doesntexist", "doc")
	app, _ := appWithErrWriter()
	c := cli.NewContext(app, set, nil)
	err := CmdHostAdd(c)
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}

func TestCmdHostAddNewHostName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	hooIP := "10.0.0.2"
	assert.Nil(t, set.Parse([]string{"barz", hooIP, "hoo"}))
	app, _ := appWithErrWriter()
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdHostAdd(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	expectedHost := config.Host{Current: "hoo", Options: map[string]string{"hoo": hooIP}}
	assert.Equal(t, expectedHost, modifiedConfigData.Hosts["barz"])
}

func TestCmdHostAddNewGlobalIpHostName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"barz", "baz"}))
	app, _ := appWithErrWriter()
	c := cli.NewContext(app, set, nil)
	assert.Nil(t, CmdHostAdd(c))

	modifiedConfigData, err := config.LoadConfigFromFile(configFileName)
	assert.Nil(t, err)

	expectedHost := config.Host{Current: "baz", Options: map[string]string{}}
	assert.Equal(t, expectedHost, modifiedConfigData.Hosts["barz"])
}

func TestCmdHostAddGlobalIpOverwriteFails(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	err := set.Parse([]string{"goo", "baz"})
	assert.Nil(t, err)

	app := cli.NewApp()
	c := cli.NewContext(app, set, nil)
	err = CmdHostAdd(c)
	assert.EqualError(t, err, "IP goo already exists")
}

func TestCompleteHostAddHostName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{}))
	app, writer := appWithWriter()
	app.Commands = []cli.Command{{Name: "add"}}
	c := cli.NewContext(app, set, nil)
	CompleteHostAdd(c)

	assert.Equal(t, "bar\nbaz.com\ngoo\n", writer.String())
}

func TestCompleteHostAddOptionName(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"bar"}))
	app, writer := appWithWriter()
	app.Commands = []cli.Command{{Name: "add"}}
	c := cli.NewContext(app, set, nil)
	CompleteHostAdd(c)

	assert.Equal(t, "10.0.0.4:baz\n10.0.0.7:bazz\n10.0.0.8:foop\nbaz\n", writer.String())
}

func TestCompleteHostAddIP(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"bar", "127.0.0.1"}))
	app, writer := appWithWriter()
	app.Commands = []cli.Command{{Name: "add"}}
	c := cli.NewContext(app, set, nil)
	CompleteHostAdd(c)

	assert.Equal(t, "bazz\nfoop\n", writer.String())
}

func TestCompleteHostAddFlags(t *testing.T) {
	configFileName, set := setupBaseConfigFile(t)
	defer removeFile(t, configFileName)
	assert.Nil(t, set.Parse([]string{"bar", "127.0.0.1", "hoo"}))
	app, writer := appWithWriter()
	app.Commands = []cli.Command{
		{
			Name: "add",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "force",
				},
			},
		},
	}
	c := cli.NewContext(app, set, nil)
	CompleteHostAdd(c)

	assert.Equal(t, "--force\n", writer.String())
}

func TestCompleteHostAddNoConfig(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	app, writer := appWithWriter()
	app.Commands = []cli.Command{{Name: "add"}}
	c := cli.NewContext(app, set, nil)
	CompleteHostAdd(c)

	assert.Equal(t, "", writer.String())
}
