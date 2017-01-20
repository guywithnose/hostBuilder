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
	outputFile, err := ioutil.TempFile("/tmp", "output")
	assert.Nil(t, err)
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
	if string(hostsFile) != expectedHostsFile {
		t.Fatalf("Output hostsFile was \n%s, expected \n%s\n", hostsFile, expectedHostsFile)
	}
	assert.Nil(t, os.Remove(configFile.Name()))
	assert.Nil(t, os.Remove(outputFile.Name()))
}

func TestCmdBuildInvalidConfigFile(t *testing.T) {
	outputFile, err := ioutil.TempFile("/tmp", "output")
	assert.Nil(t, err)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, err)

	set.String("config", "/doesntexist", "doc")
	set.String("output", outputFile.Name(), "doc")
	c := cli.NewContext(nil, set, nil)
	err = CmdBuild(c)
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
	assert.Nil(t, os.Remove(outputFile.Name()))
}

func TestCmdBuildNoConfig(t *testing.T) {
	outputFile, err := ioutil.TempFile("/tmp", "output")
	set := flag.NewFlagSet("test", 0)
	set.String("output", outputFile.Name(), "doc")
	c := cli.NewContext(nil, set, nil)
	err = CmdBuild(c)
	if err == nil {
		t.Fatal("Error was not thrown when config file was not given")
	}

	assert.EqualError(t, err, "You must specify a config file")
	assert.Nil(t, os.Remove(outputFile.Name()))
}

func TestCmdBuildNoOutput(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	set := flag.NewFlagSet("test", 0)
	configData := &config.HostsConfig{Hosts: map[string]config.Host{"foo.bar": config.Host{Current: "test", Options: map[string]string{"test": "10.0.0.1"}}}}
	err = config.WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)

	set.String("config", configFile.Name(), "doc")
	c := cli.NewContext(nil, set, nil)
	err = CmdBuild(c)
	if err == nil {
		t.Fatal("Error was not thrown when config file was not given")
	}

	assert.EqualError(t, err, "You must specify an output file")
	assert.Nil(t, os.Remove(configFile.Name()))
}
