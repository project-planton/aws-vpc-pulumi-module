package localz

import (
	"fmt"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsvpc"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/apiresource/enums/apiresourcekind"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strconv"
)

type SubnetName string
type SubnetCidr string
type AvailabilityZone string

type Locals struct {
	AwsVpc           *awsvpc.AwsVpc
	AwsTags          map[string]string
	PrivateSubnetMap map[AvailabilityZone]map[SubnetName]SubnetCidr
	PublicSubnetMap  map[AvailabilityZone]map[SubnetName]SubnetCidr
}

func Initialize(ctx *pulumi.Context, stackInput *awsvpc.AwsVpcStackInput) *Locals {
	locals := &Locals{}

	//assign value for the locals variable to make it available across the project
	locals.AwsVpc = stackInput.ApiResource

	locals.AwsTags = map[string]string{
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsVpc.Spec.EnvironmentInfo.OrgId,
		awstagkeys.Environment:  locals.AwsVpc.Spec.EnvironmentInfo.EnvId,
		awstagkeys.ResourceKind: apiresourcekind.ApiResourceKind_aws_vpc.String(),
		awstagkeys.ResourceId:   locals.AwsVpc.Metadata.Id,
	}

	locals.PrivateSubnetMap = GetPrivateSubnetMap(locals.AwsVpc)
	locals.PublicSubnetMap = GetPublicSubnetMap(locals.AwsVpc)

	return locals
}

func GetPrivateSubnetMap(awsVpc *awsvpc.AwsVpc) map[AvailabilityZone]map[SubnetName]SubnetCidr {
	privateSubnetMap := make(map[AvailabilityZone]map[SubnetName]SubnetCidr, 0)

	for azIndex, az := range awsVpc.Spec.AvailabilityZones {
		for subnetIndex := 0; subnetIndex < int(awsVpc.Spec.SubnetsPerAvailabilityZone); subnetIndex++ {
			//build private subnet name
			privateSubnetName := fmt.Sprintf("private-subnet-%s-%d", az, subnetIndex)
			//calculate private subnet cidr
			privateSubnetCidr := fmt.Sprintf("10.0.%d.0/%d", 100+azIndex*10+subnetIndex, awsVpc.Spec.SubnetSize)

			// Initialize the map for this AvailabilityZone if it doesn't exist
			if privateSubnetMap[AvailabilityZone(az)] == nil {
				privateSubnetMap[AvailabilityZone(az)] = make(map[SubnetName]SubnetCidr)
			}

			//add private subnet to the locals map
			privateSubnetMap[AvailabilityZone(az)][SubnetName(privateSubnetName)] = SubnetCidr(privateSubnetCidr)
		}
	}
	return privateSubnetMap
}

func GetPublicSubnetMap(awsVpc *awsvpc.AwsVpc) map[AvailabilityZone]map[SubnetName]SubnetCidr {
	publicSubnetMap := make(map[AvailabilityZone]map[SubnetName]SubnetCidr, 0)

	for azIndex, az := range awsVpc.Spec.AvailabilityZones {
		for subnetIndex := 0; subnetIndex < int(awsVpc.Spec.SubnetsPerAvailabilityZone); subnetIndex++ {
			//build public subnet name
			publicSubnetName := fmt.Sprintf("public-subnet-%s-%d", az, subnetIndex)
			//calculate public subnet cidr
			publicSubnetCidr := fmt.Sprintf("10.0.%d.0/%d", azIndex*10+subnetIndex, awsVpc.Spec.SubnetSize)
			// Initialize the map for this AvailabilityZone if it doesn't exist
			if publicSubnetMap[AvailabilityZone(az)] == nil {
				publicSubnetMap[AvailabilityZone(az)] = make(map[SubnetName]SubnetCidr)
			}
			//add public subnet to the locals map
			publicSubnetMap[AvailabilityZone(az)][SubnetName(publicSubnetName)] = SubnetCidr(publicSubnetCidr)
		}
	}
	return publicSubnetMap
}
