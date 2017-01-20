package awsUtil

import "text/template"

// AwsInterface defines a simple way to interact with AWS
type AwsInterface interface {
	// ReadAllLoadBalancers gets the load balancer information for all regions
	ReadAllLoadBalancers() (map[string]string, error)
	// ReadAllInstances gets the instance information for all regions
	ReadAllInstances(templ *template.Template) (map[string]string, error)
	// ListAllProfiles lists all available aws credential profiles
	ListAllProfiles() ([]string, error)
}
