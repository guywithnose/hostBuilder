package command

import (
	"bytes"
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"testing"
	"text/template"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func setupBaseConfigFile(t *testing.T) (string, *flag.FlagSet) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)

	configData := &config.HostsConfig{
		Groups: map[string][]string{"foo": {"baz.com", "goo"}},
		Hosts: map[string]config.Host{
			"bar":     {Current: hostIgnore},
			"baz.com": {Current: "baz", Options: map[string]string{"bazz": "10.0.0.7"}},
			"goo":     {Current: "foop", Options: map[string]string{"foop": "10.0.0.8"}},
		},
		GlobalIPs: map[string]string{"baz": "10.0.0.4"},
	}

	err = config.WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)

	set := flag.NewFlagSet("test", 0)
	set.String("config", configFile.Name(), "doc")

	return configFile.Name(), set
}

func setupInvalidConfigFile(t *testing.T) string {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)

	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, set.Parse([]string{"foo", "bart"}))

	configData := &config.HostsConfig{
		Groups: nil,
		Hosts: map[string]config.Host{
			"bar":     {Current: hostIgnore},
			"baz.com": {Current: "baz", Options: map[string]string{"bazz": "10.0.0.7"}},
			"goo":     {Current: "foop", Options: map[string]string{"foop": "10.0.0.8"}},
			"unknown": {Current: "unknown"},
		},
		GlobalIPs: map[string]string{"baz": "10.0.0.4"},
	}

	err = config.WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)
	return configFile.Name()
}

type awsTestUtil struct {
	instances     map[string]string
	loadBalancers map[string]string
	profiles      []string
	throwError    bool
}

// ReadAllLoadBalancers gets the load balancer information for all regions
func (util *awsTestUtil) ReadAllLoadBalancers() (map[string]string, error) {
	if util.throwError {
		return nil, errors.New("error")
	}

	return util.loadBalancers, nil
}

// ReadAllInstances gets the instance information for all regions
func (util *awsTestUtil) ReadAllInstances(templ *template.Template) (map[string]string, error) {
	if util.throwError {
		return nil, errors.New("error")
	}

	return util.instances, nil
}

// ListAllProfiles lists all available aws credential profiles
func (util *awsTestUtil) ListAllProfiles() ([]string, error) {
	if util.throwError {
		return nil, errors.New("error")
	}

	return util.profiles, nil
}

// SetProfile sets the aws credential profile to use
func (util *awsTestUtil) SetProfile(string) {
}

func appWithWriter() (*cli.App, *bytes.Buffer) {
	app := cli.NewApp()
	writer := new(bytes.Buffer)
	app.Writer = writer

	return app, writer
}

func appWithErrWriter() (*cli.App, *bytes.Buffer) {
	app := cli.NewApp()
	errWriter := new(bytes.Buffer)
	app.ErrWriter = errWriter

	return app, errWriter
}

func removeFile(t *testing.T, fileName string) {
	assert.Nil(t, os.Remove(fileName))
}
