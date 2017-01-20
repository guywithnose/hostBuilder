package command

import (
	"flag"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestReadHostsFile(t *testing.T) {
	hostsFile, err := ioutil.TempFile("/tmp", "hosts")
	assert.Nil(t, err)
	hostsLines := `10.0.0.2 bing
10.0.0.2 foo.bar
10.0.0.3 foo.bar
127.0.1.1  foo
fe00::0	ip6-localnet
ff00::0 ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
127.0.0.1 localhost localhost.localdomain localhost4 localhost4.localdomain4
::1 ip6-localhost ip6-loopback localhost localhost.localdomain localhost6 localhost6.localdomain6
`
	err = ioutil.WriteFile(hostsFile.Name(), []byte(hostsLines), 0644)
	assert.Nil(t, err)

	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, err)

	set.String("config", configFile.Name(), "doc")
	set.String("hostsFile", hostsFile.Name(), "doc")
	c := cli.NewContext(nil, set, nil)
	err = CmdCreateConfig(c)
	assert.Nil(t, err)

	configData, err := config.LoadConfigFromFile(configFile.Name())
	assert.Nil(t, err)

	expectedHosts := map[string]config.Host{
		"bing":                    config.Host{Current: "default", Options: map[string]string{"default": "10.0.0.2"}},
		"foo.bar":                 config.Host{Current: "default", Options: map[string]string{"default": "10.0.0.2", "default2": "10.0.0.3"}},
		"ip6-allnodes":            config.Host{Current: "default", Options: map[string]string{"default": "ff02::1"}},
		"ip6-allrouters":          config.Host{Current: "default", Options: map[string]string{"default": "ff02::2"}},
		"ip6-localhost":           config.Host{Current: "default", Options: map[string]string{"default": "::1"}},
		"ip6-localnet":            config.Host{Current: "default", Options: map[string]string{"default": "fe00::0"}},
		"ip6-loopback":            config.Host{Current: "default", Options: map[string]string{"default": "::1"}},
		"ip6-mcastprefix":         config.Host{Current: "default", Options: map[string]string{"default": "ff00::0"}},
		"localhost":               config.Host{Current: "default", Options: map[string]string{"default": "::1"}},
		"localhost6":              config.Host{Current: "default", Options: map[string]string{"default": "::1"}},
		"localhost6.localdomain6": config.Host{Current: "default", Options: map[string]string{"default": "::1"}},
		"localhost.localdomain":   config.Host{Current: "default", Options: map[string]string{"default": "::1"}},
	}
	if !reflect.DeepEqual(configData.Hosts, expectedHosts) {
		t.Fatalf("File was \n%v\n, expected \n%v\n", configData.Hosts, expectedHosts)
	}

	expectedLocalHosts := []string{"foo"}
	if !reflect.DeepEqual(configData.LocalHostnames, expectedLocalHosts) {
		t.Fatalf("File was \n%v\n, expected \n%v\n", configData.LocalHostnames, expectedLocalHosts)
	}

	assert.Nil(t, os.Remove(hostsFile.Name()))
	assert.Nil(t, os.Remove(configFile.Name()))
}

func TestReadHostsFileBadConfigFile(t *testing.T) {
	hostsFile, err := ioutil.TempFile("/tmp", "hosts")
	assert.Nil(t, err)
	hostsLines := `10.0.0.2 bing
10.0.0.2 foo.bar
10.0.0.3 foo.bar
127.0.1.1  foo
fe00::0	ip6-localnet
ff00::0 ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
127.0.0.1 localhost localhost.localdomain localhost4 localhost4.localdomain4
::1 ip6-localhost ip6-loopback localhost localhost.localdomain localhost6 localhost6.localdomain6
`
	err = ioutil.WriteFile(hostsFile.Name(), []byte(hostsLines), 0644)
	assert.Nil(t, err)

	assert.Nil(t, err)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, err)

	set.String("config", "/doesntexist", "doc")
	set.String("hostsFile", hostsFile.Name(), "doc")
	c := cli.NewContext(nil, set, nil)
	err = CmdCreateConfig(c)
	assert.EqualError(t, err, "open /doesntexist: permission denied")

	assert.Nil(t, os.Remove(hostsFile.Name()))
}

func TestReadHostsFileNoConfigFile(t *testing.T) {
	hostsFile, err := ioutil.TempFile("/tmp", "hosts")
	assert.Nil(t, err)
	hostsLines := `10.0.0.2 bing
10.0.0.2 foo.bar
10.0.0.3 foo.bar
127.0.1.1  foo
fe00::0	ip6-localnet
ff00::0 ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
127.0.0.1 localhost localhost.localdomain localhost4 localhost4.localdomain4
::1 ip6-localhost ip6-loopback localhost localhost.localdomain localhost6 localhost6.localdomain6
`
	err = ioutil.WriteFile(hostsFile.Name(), []byte(hostsLines), 0644)
	assert.Nil(t, err)

	assert.Nil(t, err)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, err)

	set.String("hostsFile", hostsFile.Name(), "doc")
	c := cli.NewContext(nil, set, nil)
	err = CmdCreateConfig(c)
	assert.EqualError(t, err, "You must specify a config file")

	assert.Nil(t, os.Remove(hostsFile.Name()))
}

func TestReadHostsFileBadHostsFile(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, err)

	set.String("config", configFile.Name(), "doc")
	set.String("hostsFile", "/doesntexist", "doc")
	c := cli.NewContext(nil, set, nil)
	err = CmdCreateConfig(c)
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
	assert.Nil(t, os.Remove(configFile.Name()))
}
