package command

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCmdBuild(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	outputFile, err := ioutil.TempFile("/tmp", "output")
	assert.Nil(t, err)
	defer removeFile(t, configFile.Name())
	defer removeFile(t, outputFile.Name())
	set := flag.NewFlagSet("test", 0)
	configData := &config.HostsConfig{Hosts: map[string]config.Host{"foo.bar": config.Host{Current: "test", Options: map[string]string{"test": "10.0.0.1"}}}}
	err = config.WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)

	set.String("config", configFile.Name(), "doc")
	set.String("output", outputFile.Name(), "doc")
	c := cli.NewContext(nil, set, nil)
	err = CmdBuild(c)
	assert.Nil(t, err)

	hostsFile, err := ioutil.ReadFile(outputFile.Name())
	assert.Nil(t, err)

	expectedHostsFile := "10.0.0.1 foo.bar\n127.0.0.1 localhost\n127.0.0.1 localhost.localdomain\n127.0.0.1 localhost4\n127.0.0.1 localhost4.localdomain4\n"
	assert.Equal(t, expectedHostsFile, string(hostsFile))
}

func TestCmdBuildInvalidConfigFile(t *testing.T) {
	outputFile, err := ioutil.TempFile("/tmp", "output")
	assert.Nil(t, err)
	defer removeFile(t, outputFile.Name())
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, err)

	set.String("config", "/doesntexist", "doc")
	set.String("output", outputFile.Name(), "doc")
	c := cli.NewContext(nil, set, nil)
	err = CmdBuild(c)
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}

func TestCmdBuildNoConfig(t *testing.T) {
	outputFile, err := ioutil.TempFile("/tmp", "output")
	assert.Nil(t, err)
	defer removeFile(t, outputFile.Name())
	set := flag.NewFlagSet("test", 0)
	set.String("output", outputFile.Name(), "doc")
	c := cli.NewContext(nil, set, nil)
	err = CmdBuild(c)
	assert.EqualError(t, err, "You must specify a config file")
}

func TestCmdBuildNoOutput(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	defer removeFile(t, configFile.Name())
	set := flag.NewFlagSet("test", 0)
	configData := &config.HostsConfig{Hosts: map[string]config.Host{"foo.bar": config.Host{Current: "test", Options: map[string]string{"test": "10.0.0.1"}}}}
	err = config.WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)

	set.String("config", configFile.Name(), "doc")
	c := cli.NewContext(nil, set, nil)
	err = CmdBuild(c)
	assert.EqualError(t, err, "You must specify an output file")
}

func TestCompleteBuild(t *testing.T) {
	app, writer := appWithWriter()
	app.Commands = []cli.Command{
		{
			Name: "build",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "output, o",
				},
				cli.BoolFlag{
					Name: "oneLinePerIP",
				},
			},
		},
	}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)
	CompleteBuild(c)

	assert.Equal(t, "--output\n--oneLinePerIP\n", writer.String())
}

func TestCompleteBuildOuput(t *testing.T) {
	app, writer := appWithWriter()
	set := flag.NewFlagSet("test", 0)
	os.Args = []string{"hostBuilder", "build", "--output", "--completion"}
	c := cli.NewContext(app, set, nil)
	CompleteBuild(c)

	assert.Equal(t, "fileCompletion\n", writer.String())
}
