package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteConfig(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	defer removeFile(t, configFile.Name())

	configData := getTestingConfig()

	err = WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)

	configBytes, err := ioutil.ReadFile(configFile.Name())
	assert.Nil(t, err)

	assert.Equal(t, getTestingConfigJSONString(), string(configBytes))
}

func TestWriteConfigInvalidFile(t *testing.T) {
	configData := getTestingConfig()
	err := WriteConfig("/doesntexist", configData)
	assert.EqualError(t, err, "open /doesntexist: permission denied")
}

func TestLoadConfigFromFile(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	defer removeFile(t, configFile.Name())

	err = ioutil.WriteFile(configFile.Name(), []byte(getTestingConfigJSONString()), 0644)
	assert.Nil(t, err)

	configData, err := LoadConfigFromFile(configFile.Name())
	assert.Nil(t, err)

	assert.Equal(t, getTestingConfig(), configData)
}

func TestLoadConfigFromFileInvalidFile(t *testing.T) {
	_, err := LoadConfigFromFile("/doesntexist")
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}

func TestLoadConfigFromFileInvalidJSON(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	defer removeFile(t, configFile.Name())

	err = ioutil.WriteFile(configFile.Name(), []byte("{"), 0644)
	assert.Nil(t, err)

	_, err = LoadConfigFromFile(configFile.Name())
	assert.EqualError(t, err, "unexpected end of JSON input")
}

func TestLoadEmptyConfigAndWrite(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	defer removeFile(t, configFile.Name())

	err = ioutil.WriteFile(configFile.Name(), []byte("{}"), 0644)
	assert.Nil(t, err)

	configData, err := LoadConfigFromFile(configFile.Name())
	assert.Nil(t, err)

	err = WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)

	configBytes, err := ioutil.ReadFile(configFile.Name())
	assert.Nil(t, err)

	assert.Equal(t, `{}`, string(configBytes))
}

func TestLoadEmptyHostConfigAndWrite(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	defer removeFile(t, configFile.Name())

	err = ioutil.WriteFile(configFile.Name(), []byte(`{"hosts":{"hostname":{}}}`), 0644)
	assert.Nil(t, err)

	configData, err := LoadConfigFromFile(configFile.Name())
	assert.Nil(t, err)

	err = WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)

	configBytes, err := ioutil.ReadFile(configFile.Name())
	assert.Nil(t, err)

	expectedJSONString := `{
  "hosts": {
    "hostname": {}
  }
}`
	assert.Equal(t, expectedJSONString, string(configBytes))
}
func TestBuildConfigFromHosts(t *testing.T) {
	hosts := map[string][]string{
		"bing":                    {"10.0.0.2"},
		"noips":                   {},
		"foo.bar":                 {"10.0.0.2", "10.0.0.3", "10.0.0.4"},
		"foo":                     {"127.0.1.1"},
		"ip6-allnodes":            {"ff02::1"},
		"ip6-allrouters":          {"ff02::2"},
		"ip6-localhost":           {"::1"},
		"ip6-localnet":            {"fe00::0"},
		"ip6-loopback":            {"::1"},
		"ip6-mcastprefix":         {"ff00::0"},
		"localhost4.localdomain4": {"127.0.0.1"},
		"localhost4":              {"127.0.0.1"},
		"localhost6.localdomain6": {"::1"},
		"localhost6":              {"::1"},
		"localhost.localdomain":   {"127.0.0.1", "::1"},
		"localhost":               {"127.0.0.1", "::1"},
	}

	configData := BuildConfigFromHosts(hosts)

	expectedHosts := map[string]Host{
		"bing":                    {Current: "default", Options: map[string]string{"default": "10.0.0.2"}},
		"foo.bar":                 {Current: "default", Options: map[string]string{"default": "10.0.0.2", "default2": "10.0.0.3", "default3": "10.0.0.4"}},
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
	if !reflect.DeepEqual(configData.Hosts, expectedHosts) {
		t.Fatalf("File was \n%v\n, expected \n%v\n", configData.Hosts, expectedHosts)
	}

	expectedLocalHosts := []string{"foo"}
	if !reflect.DeepEqual(configData.LocalHostnames, expectedLocalHosts) {
		t.Fatalf("File was \n%v\n, expected \n%v\n", configData.LocalHostnames, expectedLocalHosts)
	}
}

func getTestingConfig() *HostsConfig {
	return &HostsConfig{
		LocalHostnames: []string{"foo", "bar"},
		Hosts:          map[string]Host{"foo.bar": {Current: "test", Options: map[string]string{"test": "10.0.0.1"}}},
		IPv6Defaults:   true,
		GlobalIPs:      map[string]string{"foo": "bar"},
		Groups:         map[string][]string{"fooGroup": {"foo.bar"}},
	}
}

func getTestingConfigJSONString() string {
	return `{
  "localHostnames": [
    "foo",
    "bar"
  ],
  "ipV6Defaults": true,
  "hosts": {
    "foo.bar": {
      "current": "test",
      "options": {
        "test": "10.0.0.1"
      }
    }
  },
  "globalIPs": {
    "foo": "bar"
  },
  "groups": {
    "fooGroup": [
      "foo.bar"
    ]
  }
}`
}

func removeFile(t *testing.T, fileName string) {
	assert.Nil(t, os.Remove(fileName))
}
