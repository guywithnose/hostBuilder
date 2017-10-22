package hosts

import (
	"io/ioutil"
	"net"
	"regexp"
	"sort"
	"strings"

	"github.com/guywithnose/hostBuilder/config"
)

// OutputHostLines writes the hostLines array to a file
func OutputHostLines(outputFile string, configData *config.HostsConfig, oneLinePerIP bool) error {
	hostLines := buildHostLines(configData)
	output := ""
	ips := make([]string, 0, len(hostLines))
	for ip := range hostLines {
		ips = append(ips, ip)
	}

	sort.Strings(ips)

	for _, ip := range ips {
		if len(hostLines[ip]) > 0 {
			sort.Strings(hostLines[ip])
			if oneLinePerIP {
				output = output + ip + " " + strings.Join(hostLines[ip], " ") + "\n"
			} else {
				for _, hostName := range hostLines[ip] {
					output = output + ip + " " + hostName + "\n"
				}
			}
		}
	}

	return ioutil.WriteFile(outputFile, []byte(output), 0644)
}

func buildHostLines(configData *config.HostsConfig) map[string][]string {
	hostLines := map[string][]string{
		"127.0.0.1": {"localhost", "localhost.localdomain", "localhost4", "localhost4.localdomain4"},
		"127.0.1.1": configData.LocalHostnames,
	}

	if configData.IPv6Defaults {
		hostLines["::1"] = []string{"ip6-localhost", "ip6-loopback", "localhost", "localhost.localdomain", "localhost6", "localhost6.localdomain6"}
		hostLines["fe00::0"] = []string{"ip6-localnet"}
		hostLines["ff00::0"] = []string{"ip6-mcastprefix"}
		hostLines["ff02::1"] = []string{"ip6-allnodes"}
		hostLines["ff02::2"] = []string{"ip6-allrouters"}
	}

	for hostName, data := range configData.Hosts {
		ip, ok := data.Options[data.Current]
		if !ok {
			ip, ok = configData.GlobalIPs[data.Current]
		}

		if ok {
			if existingHosts, ok := hostLines[ip]; ok {
				hostLines[ip] = append(existingHosts, hostName)
			} else {
				hostLines[ip] = []string{hostName}
			}
		}
	}

	return hostLines
}

// ReadHostsFile reads a hosts file and returns the parsed hostnames and ips
func ReadHostsFile(fileName string) (map[string][]string, error) {
	hostsData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	hostsLines := strings.Split(string(hostsData), "\n")

	hosts := make(map[string][]string)
	for _, line := range hostsLines {
		ip, hostnames := parseHostLine(line)
		for _, hostname := range hostnames {
			if _, exists := hosts[hostname]; exists {
				hosts[hostname] = append(hosts[hostname], ip)
			} else {
				hosts[hostname] = []string{ip}
			}
		}
	}

	return hosts, nil
}

func parseHostLine(line string) (string, []string) {
	// Clear out any comments
	line = strings.Replace(line, "\t", " ", -1)
	line = strings.TrimSpace(line)
	fullLineCommentRegex := regexp.MustCompile("^#")
	line = fullLineCommentRegex.ReplaceAllString(line, "")
	commentRegex := regexp.MustCompile("#.*")
	line = commentRegex.ReplaceAllString(line, "")
	line = strings.TrimSpace(line)
	parts := strings.Split(line, " ")
	if len(parts) >= 2 {
		hostsnames := make([]string, 0, len(parts)-1)
		IP := strings.TrimSpace(parts[0])
		if net.ParseIP(IP) == nil {
			return "", nil
		}

		hostsToParse := parts[1:]
		for _, hostname := range hostsToParse {
			hostname = strings.TrimSpace(hostname)
			if hostname == "" {
				continue
			}

			hostsnames = append(hostsnames, hostname)
		}

		return IP, hostsnames
	}

	return "", nil
}
