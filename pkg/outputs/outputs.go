package outputs

import (
	"fmt"
	"github.com/plantoncloud/aws-vpc-pulumi-module/pkg/localz"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsvpc"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	VpcId             = "vpc-id"
	InternetGatewayId = "internet-gateway-id"
)

func PulumiOutputsToStackOutputsConverter(pulumiOutputs auto.OutputMap, stackInput *awsvpc.AwsVpcStackInput) *awsvpc.AwsVpcStackOutputs {
	privateAzSubnetMap := localz.GetPrivateAzSubnetMap(stackInput.ApiResource)
	sortedPrivateAzKeys := localz.GetSortedAzKeys(privateAzSubnetMap)
	privateSubnetOutputs := make([]*awsvpc.AwsVpcSubnetStackOutputs, 0)

	for _, availabilityZone := range sortedPrivateAzKeys {
		azSubnetMap := privateAzSubnetMap[localz.AvailabilityZone(availabilityZone)]
		sortedSubnetNames := localz.GetSortedSubnetNameKeys(azSubnetMap)
		for _, subnetName := range sortedSubnetNames {
			subnetStackOutputs := &awsvpc.AwsVpcSubnetStackOutputs{
				Name: subnetName,
				Id:   autoapistackoutput.GetVal(pulumiOutputs, SubnetIdOutputKey(subnetName)),
				Cidr: autoapistackoutput.GetVal(pulumiOutputs, SubnetCidrOutputKey(subnetName)),
			}
			if stackInput.ApiResource.Spec.IsNatGatewayEnabled {
				subnetStackOutputs.NatGateway = &awsvpc.AwsVpcNatGatewayStackOutputs{
					Id:        autoapistackoutput.GetVal(pulumiOutputs, NatGatewayIdOutputKey(subnetName)),
					PrivateIp: autoapistackoutput.GetVal(pulumiOutputs, NatGatewayPrivateIpOutputKey(subnetName)),
					PublicIp:  autoapistackoutput.GetVal(pulumiOutputs, NatGatewayPublicIpOutputKey(subnetName)),
				}
			}
			privateSubnetOutputs = append(privateSubnetOutputs, subnetStackOutputs)
		}
	}

	publicAzSubnetMap := localz.GetPublicAzSubnetMap(stackInput.ApiResource)
	sortedPublicAzKeys := localz.GetSortedAzKeys(publicAzSubnetMap)
	publicSubnetOutputs := make([]*awsvpc.AwsVpcSubnetStackOutputs, 0)
	for _, availabilityZone := range sortedPublicAzKeys {
		azSubnetMap := publicAzSubnetMap[localz.AvailabilityZone(availabilityZone)]
		sortedSubnetNames := localz.GetSortedSubnetNameKeys(azSubnetMap)
		for _, subnetName := range sortedSubnetNames {
			subnetStackOutputs := &awsvpc.AwsVpcSubnetStackOutputs{
				Name: subnetName,
				Id:   autoapistackoutput.GetVal(pulumiOutputs, SubnetIdOutputKey(subnetName)),
				Cidr: autoapistackoutput.GetVal(pulumiOutputs, SubnetCidrOutputKey(subnetName)),
			}
			publicSubnetOutputs = append(publicSubnetOutputs, subnetStackOutputs)
		}
	}

	return &awsvpc.AwsVpcStackOutputs{
		VpcId:             autoapistackoutput.GetVal(pulumiOutputs, VpcId),
		InternetGatewayId: autoapistackoutput.GetVal(pulumiOutputs, InternetGatewayId),
		PrivateSubnets:    privateSubnetOutputs,
		PublicSubnets:     publicSubnetOutputs,
	}
}

func SubnetIdOutputKey(subnetName string) string {
	return fmt.Sprintf("%s-id", subnetName)
}

func SubnetCidrOutputKey(subnetName string) string {
	return fmt.Sprintf("%s-cidr", subnetName)
}

func NatGatewayIdOutputKey(subnetName string) string {
	return fmt.Sprintf("%s-nat-gw-id", subnetName)
}

func NatGatewayPrivateIpOutputKey(subnetName string) string {
	return fmt.Sprintf("%s-nat-gw-private-ip", subnetName)
}

func NatGatewayPublicIpOutputKey(subnetName string) string {
	return fmt.Sprintf("%s-nat-gw-public-ip", subnetName)
}
