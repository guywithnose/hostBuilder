package hosts

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/stretchr/testify/assert"
)

func TestOutputHostLinesMultipleLinesPerIP(t *testing.T) {
	hostLines := runOutputHosts(t, false)
	expectedLines := "10.0.0.2 bing\n10.0.0.2 foo.bar\n123.12.34.56 bam\n127.0.0.1 localhost\n127.0.0.1 localhost.localdomain\n127.0.0.1 localhost4\n127.0.0.1 localhost4.localdomain4\n127.0.1.1 bar\n127.0.1.1 foo\n::1 ip6-localhost\n::1 ip6-loopback\n::1 localhost\n::1 localhost.localdomain\n::1 localhost6\n::1 localhost6.localdomain6\nfe00::0 ip6-localnet\nff00::0 ip6-mcastprefix\nff02::1 ip6-allnodes\nff02::2 ip6-allrouters\n"
	assert.Equal(t, expectedLines, hostLines)
}

func TestOutputHostLinesOneLinePerIP(t *testing.T) {
	hostLines := runOutputHosts(t, true)
	expectedLines := "10.0.0.2 bing foo.bar\n123.12.34.56 bam\n127.0.0.1 localhost localhost.localdomain localhost4 localhost4.localdomain4\n127.0.1.1 bar foo\n::1 ip6-localhost ip6-loopback localhost localhost.localdomain localhost6 localhost6.localdomain6\nfe00::0 ip6-localnet\nff00::0 ip6-mcastprefix\nff02::1 ip6-allnodes\nff02::2 ip6-allrouters\n"
	assert.Equal(t, expectedLines, hostLines)
}

func TestReadHostsFile(t *testing.T) {
	outputFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	defer removeFile(t, outputFile.Name())
	hostsLines := `10.0.0.2 bing
10.0.0.2 foo.bar
10.0.0.3 foo.bar
 # 10.0.0.4 foo.bar
123.12.34.56 bam
10.0.0.256 notip
127.0.1.1 bar 
127.0.1.1  foo
	fe00::0	ip6-localnet
ff00::0 ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
127.0.0.1 localhost localhost.localdomain localhost4 localhost4.localdomain4
::1 ip6-localhost ip6-loopback localhost localhost.localdomain localhost6 localhost6.localdomain6
`
	err = ioutil.WriteFile(outputFile.Name(), []byte(hostsLines), 0644)
	assert.Nil(t, err)

	expectedHosts := map[string][]string{
		"bam":                     []string{"123.12.34.56"},
		"bar":                     []string{"127.0.1.1"},
		"bing":                    []string{"10.0.0.2"},
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

	hosts, err := ReadHostsFile(outputFile.Name())
	assert.Nil(t, err)

	assert.Equal(t, expectedHosts, hosts)
}

func TestReadHostsFileInvalidHostsFile(t *testing.T) {
	_, err := ReadHostsFile("/doesntexist")
	assert.EqualError(t, err, "open /doesntexist: no such file or directory")
}

func getTestingConfig() *config.HostsConfig {
	return &config.HostsConfig{
		LocalHostnames: []string{"foo", "bar"},
		Hosts:          map[string]config.Host{"foo.bar": config.Host{Current: "test", Options: map[string]string{"test": "10.0.0.2"}}, "bing": config.Host{Current: "goo"}, "bam": config.Host{Current: "boo"}},
		IPv6Defaults:   true,
		GlobalIPs:      map[string]string{"local": "127.0.0.1", "goo": "10.0.0.2", "boo": "123.12.34.56"},
	}
}

func runOutputHosts(t *testing.T, oneLinePerIP bool) string {
	outputFile, err := ioutil.TempFile("/tmp", "config")
	assert.Nil(t, err)
	defer removeFile(t, outputFile.Name())

	configData := getTestingConfig()

	err = OutputHostLines(outputFile.Name(), configData, oneLinePerIP)
	assert.Nil(t, err)

	hostLines, err := ioutil.ReadFile(outputFile.Name())
	assert.Nil(t, err)

	return string(hostLines)
}

func removeFile(t *testing.T, fileName string) {
	assert.Nil(t, os.Remove(fileName))
}
