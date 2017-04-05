package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

// HostsConfig defines the structure of the hosts config file
type HostsConfig struct {
	LocalHostnames []string            `json:"localHostnames,omitempty"`
	IPv6Defaults   bool                `json:"ipV6Defaults,omitempty"`
	Hosts          map[string]Host     `json:"hosts,omitempty"`
	GlobalIPs      map[string]string   `json:"globalIPs,omitempty"`
	Groups         map[string][]string `json:"groups,omitempty"`
}

// Host defines the data associated with a hostname
type Host struct {
	Current string            `json:"current,omitempty"`
	Options map[string]string `json:"options,omitempty"`
}

// LoadConfigFromFile loads a HostsConfig from a file
func LoadConfigFromFile(fileName string) (*HostsConfig, error) {
	configJSON, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var configData = new(HostsConfig)
	err = json.Unmarshal(configJSON, configData)
	if err != nil {
		return nil, err
	}

	if configData.LocalHostnames == nil {
		configData.LocalHostnames = []string{}
	}

	if configData.Hosts == nil {
		configData.Hosts = map[string]Host{}
	}

	for index, host := range configData.Hosts {
		if host.Options == nil {
			host.Options = map[string]string{}
			configData.Hosts[index] = host
		}
	}

	if configData.GlobalIPs == nil {
		configData.GlobalIPs = map[string]string{}
	}

	if configData.Groups == nil {
		configData.Groups = map[string][]string{}
	}

	return configData, nil
}

// WriteConfig saves a HostsConfig to a file
func WriteConfig(outputFile string, configData *HostsConfig) error {
	formattedConfig, err := json.MarshalIndent(configData, "", "  ")
	if err != nil {
		// This should never happen
		panic(err)
	}

	err = ioutil.WriteFile(outputFile, formattedConfig, 0644)
	if err != nil {
		return err
	}

	return nil
}

// BuildConfigFromHosts builds a config from a map of hostnames to ips
func BuildConfigFromHosts(hosts map[string][]string) *HostsConfig {
	configData := &HostsConfig{
		LocalHostnames: []string{},
		IPv6Defaults:   false,
		Hosts:          map[string]Host{},
		GlobalIPs:      map[string]string{},
		Groups:         map[string][]string{},
	}

	for hostname, ips := range hosts {
		if len(ips) == 0 {
			continue
		}

		for _, ip := range ips {
			if ip == "127.0.1.1" {
				configData.LocalHostnames = append(configData.LocalHostnames, hostname)
			} else if ip != "127.0.0.1" || !strings.Contains(hostname, "localhost") {
				if _, exists := configData.Hosts[hostname]; exists {
					IPName := fmt.Sprintf("default%d", len(configData.Hosts[hostname].Options)+1)
					configData.Hosts[hostname].Options[IPName] = ip
				} else {
					configData.Hosts[hostname] = Host{Current: "default", Options: map[string]string{"default": ip}}
				}
			}
		}
	}

	return configData
}
