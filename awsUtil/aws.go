package awsUtil

import (
	"bytes"
	"errors"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/go-ini/ini"
)

// AwsUtil handles connecting to AWS and retrieving host data
type AwsUtil struct {
	profileName string
}

// NewAwsUtil creates a new awsUtil with a given profileName
func NewAwsUtil(profileName string) *AwsUtil {
	util := new(AwsUtil)
	util.profileName = profileName
	return util
}

// ReadAllLoadBalancers gets the load balancer information for all regions
func (util *AwsUtil) ReadAllLoadBalancers() (map[string]string, error) {
	loadBalancerHostNames := make(map[string]string)
	regions, err := util.getRegions()
	if err != nil {
		return nil, err
	}

	for _, region := range regions {
		regionLbs, err := util.readLoadBalancers(region)
		if err != nil {
			return nil, err
		}

		for name, address := range regionLbs {
			loadBalancerHostNames[name] = address
		}
	}

	return loadBalancerHostNames, nil
}

// ReadAllInstances gets the instance information for all regions
func (util *AwsUtil) ReadAllInstances(templ *template.Template) (map[string]string, error) {
	instanceIPs := make(map[string]string)
	regions, err := util.getRegions()
	if err != nil {
		return nil, err
	}

	for _, region := range regions {
		regionInstances, err := util.readInstances(region, templ)
		if err != nil {
			return nil, err
		}

		for name, address := range regionInstances {
			instanceIPs[name] = address
		}
	}

	return instanceIPs, nil
}

// ListAllProfiles lists all available aws credential profiles
func (util *AwsUtil) ListAllProfiles() ([]string, error) {
	scp := credentials.SharedCredentialsProvider{}
	_, err := scp.Retrieve()
	if err != nil {
		return nil, err
	}

	config, err := ini.Load(scp.Filename)
	if err != nil {
		return nil, err
	}

	return config.SectionStrings(), nil
}

func (util *AwsUtil) readLoadBalancers(region string) (map[string]string, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(region)},
		Profile: util.profileName,
	})
	if err != nil {
		return nil, errors.New("Failed to create session")
	}

	svc := elb.New(sess)

	params := &elb.DescribeLoadBalancersInput{}
	resp, err := svc.DescribeLoadBalancers(params)
	if err != nil {
		return nil, err
	}

	loadBalancerHostNames := make(map[string]string)
	for _, loadBalancer := range resp.LoadBalancerDescriptions {
		loadBalancerHostNames[*loadBalancer.LoadBalancerName] = *loadBalancer.DNSName
	}

	return loadBalancerHostNames, nil
}

func (util *AwsUtil) readInstances(region string, templ *template.Template) (map[string]string, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(region)},
		Profile: util.profileName,
	})
	if err != nil {
		return nil, errors.New("Failed to create session")
	}

	svc := ec2.New(sess)

	params := &ec2.DescribeInstancesInput{
		MaxResults: aws.Int64(10),
	}
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		return nil, err
	}

	instanceIPs := parseReservations(resp.Reservations, templ)
	return instanceIPs, nil
}

func parseReservations(reservations []*ec2.Reservation, templ *template.Template) map[string]string {
	instanceIPs := make(map[string]string)
	for _, reservation := range reservations {
		instances := reservation.Instances
		for _, instance := range instances {
			name, ip, err := parseInstance(instance, templ)
			if err == nil {
				instanceIPs[name] = ip
			}
		}
	}

	return instanceIPs
}

func parseInstance(instance *ec2.Instance, templ *template.Template) (string, string, error) {
	if instance.PublicIpAddress == nil {
		return "", "", errors.New("instance has no ip address")
	}

	var buffer bytes.Buffer
	err := templ.Execute(&buffer, instance)
	if err != nil {
		return "", "", err
	}

	return buffer.String(), *instance.PublicIpAddress, nil
}

func (util *AwsUtil) getRegions() ([]string, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String("us-east-1")},
		Profile: util.profileName,
	})
	if err != nil {
		return nil, err
	}

	svc := ec2.New(sess)
	resultRegions, err := svc.DescribeRegions(nil)
	if err != nil {
		return nil, err
	}

	regions := resultRegions.Regions
	regionNames := make([]string, 0, len(regions))
	for _, region := range regions {
		regionNames = append(regionNames, *region.RegionName)
	}

	return regionNames, nil
}

// SetProfile sets the aws credential profile to use
func (util *AwsUtil) SetProfile(profile string) {
	util.profileName = profile
}
