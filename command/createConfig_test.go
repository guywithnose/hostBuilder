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

func TestReadHostsFile(t *testing.T) {
	hostsFile, err := ioutil.TempFile("/tmp", "hosts")
	assert.Nil(t, err)
	defer removeFile(t, hostsFile.Name())
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
	defer removeFile(t, configFile.Name())
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
		"bing":                    {Current: "default", Options: map[string]string{"default": "10.0.0.2"}},
		"foo.bar":                 {Current: "default", Options: map[string]string{"default": "10.0.0.2", "default2": "10.0.0.3"}},
		"ip6-allnodes":            {Current: "default", Options: map[string]string{"default": "ff02::1"}},
		"ip6-allrouters":          {Current: "default", Options: map[string]string{"default": "ff02::2"}},
		"ip6-localhost":           {Current: "default", Options: map[string]string{"default": "::1"}},
		"ip6-localnet":            {Current: "default", Options: map[string]string{"default": "fe00::0"}},
		"ip6-loopback":            {Current: "default", Options: map[string]string{"default": "::1"}},
		"ip6-mcastprefix":         {Current: "default", Options: map[string]string{"default": "ff00::0"}},
		"localhost":               {Current: "default", Options: map[string]string{"default": "::1"}},
		"localhost6":              {Current: "default", Options: map[string]string{"default": "::1"}},
		"localhost6.localdomain6": {Current: "default", Options: map[string]string{"default": "::1"}},
		"localhost.localdomain":   {Current: "default", Options: map[string]string{"default": "::1"}},
	}

	assert.Equal(t, expectedHosts, configData.Hosts)
	assert.Equal(t, []string{"foo"}, configData.LocalHostnames)
}

func TestReadHostsFileBadConfigFile(t *testing.T) {
	hostsFile, err := ioutil.TempFile("/tmp", "hosts")
	assert.Nil(t, err)
	defer removeFile(t, hostsFile.Name())
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
}

func TestReadHostsFileNoConfigFile(t *testing.T) {
	hostsFile, err := ioutil.TempFile("/tmp", "hosts")
	assert.Nil(t, err)
	defer removeFile(t, hostsFile.Name())
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
}

func TestReadHostsFileBadHostsFile(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	defer removeFile(t, configFile.Name())
	set := flag.NewFlagSet("test", 0)
	assert.Nil(t, err)

	set.String("config", configFile.Name(), "doc")
	set.String("hostsFile", "/doesntexist", "doc")
	c := cli.NewContext(nil, set, nil)
	err = CmdCreateConfig(c)
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}

func TestCompleteCreateConfig(t *testing.T) {
	app, writer := appWithWriter()
	app.Commands = []cli.Command{
		{
			Name: "createConfig",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "hostsFile, hosts",
				},
			},
		},
	}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)
	CompleteCreateConfig(c)

	assert.Equal(t, "--hostsFile\n", writer.String())
}

func TestCompleteCreateConfigHostsFile(t *testing.T) {
	app, writer := appWithWriter()
	set := flag.NewFlagSet("test", 0)
	os.Args = []string{"hostBuilder", "createConfig", "--hostsFile", "--completion"}
	c := cli.NewContext(app, set, nil)
	CompleteCreateConfig(c)

	assert.Equal(t, "fileCompletion\n", writer.String())
}
