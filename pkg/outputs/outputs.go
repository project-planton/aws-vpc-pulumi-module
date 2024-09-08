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
	resp := &awsvpc.AwsVpcStackOutputs{
		VpcId:             autoapistackoutput.GetVal(pulumiOutputs, VpcId),
		InternetGatewayId: autoapistackoutput.GetVal(pulumiOutputs, InternetGatewayId),
	}

	privateSubnetMap := localz.GetPrivateSubnetMap(stackInput.ApiResource)
	publicSubnetMap := localz.GetPrivateSubnetMap(stackInput.ApiResource)

	privateSubnetOutputs := make([]*awsvpc.AwsVpcSubnetStackOutputs, 0)
	natGatewayOutputs := make([]*awsvpc.AwsVpcNatGatewayStackOutputs, 0)
	for _, subnetNameCidrMap := range privateSubnetMap {
		for subnetName, _ := range subnetNameCidrMap {
			privateSubnetOutputs = append(privateSubnetOutputs, &awsvpc.AwsVpcSubnetStackOutputs{
				Name: string(subnetName),
				Id:   autoapistackoutput.GetVal(pulumiOutputs, SubnetIdOutputKey(string(subnetName))),
				Cidr: autoapistackoutput.GetVal(pulumiOutputs, SubnetCidrOutputKey(string(subnetName))),
			})

			natGatewayOutputs = append(natGatewayOutputs, &awsvpc.AwsVpcNatGatewayStackOutputs{
				Id:        autoapistackoutput.GetVal(pulumiOutputs, NatGatewayIdOutputKey(string(subnetName))),
				PrivateIp: autoapistackoutput.GetVal(pulumiOutputs, NatGatewayPrivateIpOutputKey(string(subnetName))),
				PublicIp:  autoapistackoutput.GetVal(pulumiOutputs, NatGatewayPublicIpOutputKey(string(subnetName))),
			})
		}
	}

	publicSubnetOutputs := make([]*awsvpc.AwsVpcSubnetStackOutputs, 0)
	for _, subnetNameCidrMap := range publicSubnetMap {
		for subnetName, _ := range subnetNameCidrMap {
			publicSubnetOutputs = append(publicSubnetOutputs, &awsvpc.AwsVpcSubnetStackOutputs{
				Name: string(subnetName),
				Id:   autoapistackoutput.GetVal(pulumiOutputs, SubnetIdOutputKey(string(subnetName))),
				Cidr: autoapistackoutput.GetVal(pulumiOutputs, SubnetCidrOutputKey(string(subnetName))),
			})
		}
	}

	return resp
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
