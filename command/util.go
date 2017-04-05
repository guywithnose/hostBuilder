package command

import (
	"errors"
	"fmt"
	"net"
	"sort"

	"github.com/guywithnose/hostBuilder/config"
	"github.com/urfave/cli"
)

func loadConfig(c *cli.Context) (*config.HostsConfig, error) {
	configFile := c.GlobalString("config")
	if configFile == "" {
		return nil, errors.New("You must specify a config file")
	}

	configData, err := config.LoadConfigFromFile(configFile)
	if err != nil {
		return nil, err
	}

	return configData, nil
}

func sortHostNames(configData *config.HostsConfig) []string {
	hostNames := make([]string, 0, len(configData.Hosts))
	for hostName := range configData.Hosts {
		hostNames = append(hostNames, hostName)
	}

	sort.Strings(hostNames)
	return hostNames
}

func sortGroupNames(configData *config.HostsConfig) []string {
	groupNames := make([]string, 0, len(configData.Groups))
	for groupName := range configData.Groups {
		groupNames = append(groupNames, groupName)
	}

	sort.Strings(groupNames)
	return groupNames
}

func groupContains(configData *config.HostsConfig, groupName, hostName string) bool {
	for _, host := range configData.Groups[groupName] {
		if host == hostName {
			return true
		}
	}

	return false
}

func sortGlobalIPNames(configData *config.HostsConfig) []string {
	globalIPNames := make([]string, 0, len(configData.GlobalIPs))
	for globalIPName := range configData.GlobalIPs {
		globalIPNames = append(globalIPNames, globalIPName)
	}

	sort.Strings(globalIPNames)
	return globalIPNames
}

func sortOptions(configData *config.HostsConfig, hostName string) []string {
	options := make([]string, 0, len(configData.Hosts[hostName].Options))
	for option := range configData.Hosts[hostName].Options {
		options = append(options, option)
	}

	sort.Strings(options)
	return options
}

func sortAllOptions(configData *config.HostsConfig) []string {
	options := make([]string, 0, 100)
	for _, host := range configData.Hosts {
		for option := range host.Options {
			options = append(options, option)
		}
	}

	sort.Strings(options)
	return options
}

func sortAllIPs(configData *config.HostsConfig) ([]string, map[string]string) {
	IPs := make([]string, 0, 100)
	IPMap := make(map[string]string)
	for _, host := range configData.Hosts {
		for option, IP := range host.Options {
			IPs = append(IPs, IP)
			IPMap[IP] = option
		}
	}

	for globalIPName, IP := range configData.GlobalIPs {
		IPs = append(IPs, IP)
		IPMap[IP] = globalIPName
	}

	sort.Strings(IPs)
	return IPs, IPMap
}

func resolveAddress(address string) (string, error) {
	if net.ParseIP(address) != nil {
		return address, nil
	}

	IPs, err := net.LookupHost(address)
	if err != nil {
		return "", fmt.Errorf("Unable to resolve %s", address)
	}

	return IPs[0], nil
}
