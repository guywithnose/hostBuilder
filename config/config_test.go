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

	configData := getTestingConfig()

	err = WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)

	configBytes, err := ioutil.ReadFile(configFile.Name())
	assert.Nil(t, err)

	if string(configBytes) != getTestingConfigJSONString() {
		t.Fatalf("File was %s, expected %s", configBytes, getTestingConfigJSONString())
	}

	assert.Nil(t, os.Remove(configFile.Name()))
}

func TestWriteConfigInvalidFile(t *testing.T) {
	configData := getTestingConfig()
	err := WriteConfig("/doesntexist", configData)
	assert.EqualError(t, err, "open /doesntexist: permission denied")
}

func TestLoadConfigFromFile(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)

	err = ioutil.WriteFile(configFile.Name(), []byte(getTestingConfigJSONString()), 0644)
	assert.Nil(t, err)

	configData, err := LoadConfigFromFile(configFile.Name())
	assert.Nil(t, err)

	expectedConfigData := getTestingConfig()

	if !reflect.DeepEqual(configData, expectedConfigData) {
		t.Fatalf("File was %v, expected %v", configData, expectedConfigData)
	}

	assert.Nil(t, os.Remove(configFile.Name()))
}

func TestLoadConfigFromFileInvalidFile(t *testing.T) {
	_, err := LoadConfigFromFile("/doesntexist")
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}

func TestLoadConfigFromFileInvalidJSON(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)

	err = ioutil.WriteFile(configFile.Name(), []byte("{"), 0644)
	assert.Nil(t, err)

	_, err = LoadConfigFromFile(configFile.Name())
	assert.EqualError(t, err, "unexpected end of JSON input")
	assert.Nil(t, os.Remove(configFile.Name()))
}

func TestLoadEmptyConfigAndWrite(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)

	err = ioutil.WriteFile(configFile.Name(), []byte("{}"), 0644)
	assert.Nil(t, err)

	configData, err := LoadConfigFromFile(configFile.Name())
	assert.Nil(t, err)

	err = WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)

	configBytes, err := ioutil.ReadFile(configFile.Name())
	assert.Nil(t, err)

	expectedJSONString := `{
  "localHostnames": [],
  "ipV6Defaults": false,
  "hosts": {},
  "globalIPs": {},
  "groups": {}
}`
	if string(configBytes) != expectedJSONString {
		t.Fatalf("File was %s, expected %s", configBytes, expectedJSONString)
	}

	assert.Nil(t, os.Remove(configFile.Name()))
}

func TestLoadEmptyHostConfigAndWrite(t *testing.T) {
	configFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)

	err = ioutil.WriteFile(configFile.Name(), []byte(`{"hosts":{"hostname":{}}}`), 0644)
	assert.Nil(t, err)

	configData, err := LoadConfigFromFile(configFile.Name())
	assert.Nil(t, err)

	err = WriteConfig(configFile.Name(), configData)
	assert.Nil(t, err)

	configBytes, err := ioutil.ReadFile(configFile.Name())
	assert.Nil(t, err)

	expectedJSONString := `{
  "localHostnames": [],
  "ipV6Defaults": false,
  "hosts": {
    "hostname": {
      "current": "",
      "options": {}
    }
  },
  "globalIPs": {},
  "groups": {}
}`
	if string(configBytes) != expectedJSONString {
		t.Fatalf("File was %s, expected %s", configBytes, expectedJSONString)
	}

	assert.Nil(t, os.Remove(configFile.Name()))
}

func TestBuildConfigFromHosts(t *testing.T) {
	hosts := map[string][]string{
		"bing":                    []string{"10.0.0.2"},
		"noips":                   []string{},
		"foo.bar":                 []string{"10.0.0.2", "10.0.0.3", "10.0.0.4"},
		"foo":                     []string{"127.0.1.1"},
		"ip6-allnodes":            []string{"ff02::1"},
		"ip6-allrouters":          []string{"ff02::2"},
		"ip6-localhost":           []string{"::1"},
		"ip6-localnet":            []string{"fe00::0"},
		"ip6-loopback":            []string{"::1"},
		"ip6-mcastprefix":         []string{"ff00::0"},
		"localhost4.localdomain4": []string{"127.0.0.1"},
		"localhost4":              []string{"127.0.0.1"},
		"localhost6.localdomain6": []string{"::1"},
		"localhost6":              []string{"::1"},
		"localhost.localdomain":   []string{"127.0.0.1", "::1"},
		"localhost":               []string{"127.0.0.1", "::1"},
	}

	configData := BuildConfigFromHosts(hosts)

	expectedHosts := map[string]Host{
		"bing":                    Host{Current: "default", Options: map[string]string{"default": "10.0.0.2"}},
		"foo.bar":                 Host{Current: "default", Options: map[string]string{"default": "10.0.0.2", "default2": "10.0.0.3", "default3": "10.0.0.4"}},
		"ip6-allnodes":            Host{Current: "default", Options: map[string]string{"default": "ff02::1"}},
		"ip6-allrouters":          Host{Current: "default", Options: map[string]string{"default": "ff02::2"}},
		"ip6-localhost":           Host{Current: "default", Options: map[string]string{"default": "::1"}},
		"ip6-localnet":            Host{Current: "default", Options: map[string]string{"default": "fe00::0"}},
		"ip6-loopback":            Host{Current: "default", Options: map[string]string{"default": "::1"}},
		"ip6-mcastprefix":         Host{Current: "default", Options: map[string]string{"default": "ff00::0"}},
		"localhost":               Host{Current: "default", Options: map[string]string{"default": "::1"}},
		"localhost6":              Host{Current: "default", Options: map[string]string{"default": "::1"}},
		"localhost6.localdomain6": Host{Current: "default", Options: map[string]string{"default": "::1"}},
		"localhost.localdomain":   Host{Current: "default", Options: map[string]string{"default": "::1"}},
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
		Hosts:          map[string]Host{"foo.bar": Host{Current: "test", Options: map[string]string{"test": "10.0.0.1"}}},
		IPv6Defaults:   true,
		GlobalIPs:      map[string]string{"foo": "bar"},
		Groups:         map[string][]string{"fooGroup": []string{"foo.bar"}},
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
