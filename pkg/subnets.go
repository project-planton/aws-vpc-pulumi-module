package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/aws-vpc-pulumi-module/pkg/localz"
	"github.com/plantoncloud/aws-vpc-pulumi-module/pkg/outputs"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/datatypes/stringmaps"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/datatypes/stringmaps/convertstringmaps"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func subnets(ctx *pulumi.Context, locals *localz.Locals, createdVpc *ec2.Vpc) (privateSubnetMap,
	publicSubnetMap map[localz.SubnetName]*ec2.Subnet, err error) {

	privateSubnetMap = make(map[localz.SubnetName]*ec2.Subnet, 0)
	publicSubnetMap = make(map[localz.SubnetName]*ec2.Subnet, 0)

	// iterate through azs and create the configured number of public and private subnets per az
	sortedPrivateAzKeys := localz.GetSortedAzKeys(locals.PrivateAzSubnetMap)
	// create private subnets
	for _, availabilityZone := range sortedPrivateAzKeys {
		azSubnetMap := locals.PrivateAzSubnetMap[localz.AvailabilityZone(availabilityZone)]
		sortedSubnetNames := localz.GetSortedSubnetNameKeys(azSubnetMap)
		for _, subnetName := range sortedSubnetNames {
			// create private subnet in az
			createdSubnet, err := ec2.NewSubnet(ctx,
				subnetName,
				&ec2.SubnetArgs{
					VpcId:            createdVpc.ID(),
					CidrBlock:        pulumi.String(azSubnetMap[localz.SubnetName(subnetName)]),
					AvailabilityZone: pulumi.String(availabilityZone),
					Tags: convertstringmaps.ConvertGoStringMapToPulumiStringMap(
						stringmaps.AddEntry(locals.AwsTags, "Name", subnetName)),
				}, pulumi.Parent(createdVpc))
			if err != nil {
				return nil, nil,
					errors.Wrapf(err, "error creating private subnet %s", subnetName)
			}
			ctx.Export(outputs.SubnetIdOutputKey(subnetName), createdSubnet.ID())
			ctx.Export(outputs.SubnetCidrOutputKey(subnetName), createdSubnet.CidrBlock)
			privateSubnetMap[localz.SubnetName(subnetName)] = createdSubnet
		}
	}

	sortedPublicAzKeys := localz.GetSortedAzKeys(locals.PublicAzSubnetMap)
	// create public subnets
	for _, availabilityZone := range sortedPublicAzKeys {
		azSubnetMap := locals.PublicAzSubnetMap[localz.AvailabilityZone(availabilityZone)]
		sortedSubnetNames := localz.GetSortedSubnetNameKeys(azSubnetMap)
		for _, subnetName := range sortedSubnetNames {
			// create public subnet in az
			createdSubnet, err := ec2.NewSubnet(ctx,
				subnetName,
				&ec2.SubnetArgs{
					VpcId:            createdVpc.ID(),
					CidrBlock:        pulumi.String(azSubnetMap[localz.SubnetName(subnetName)]),
					AvailabilityZone: pulumi.String(availabilityZone),
					//required for public subnets
					MapPublicIpOnLaunch: pulumi.Bool(true),
					Tags: convertstringmaps.ConvertGoStringMapToPulumiStringMap(
						stringmaps.AddEntry(locals.AwsTags, "Name", subnetName)),
				}, pulumi.Parent(createdVpc))
			if err != nil {
				return nil, nil,
					errors.Wrapf(err, "error creating public subnet %s", subnetName)
			}
			ctx.Export(outputs.SubnetIdOutputKey(subnetName), createdSubnet.ID())
			ctx.Export(outputs.SubnetCidrOutputKey(subnetName), createdSubnet.CidrBlock)
			publicSubnetMap[localz.SubnetName(subnetName)] = createdSubnet
		}

	}
	return privateSubnetMap, publicSubnetMap, nil
}
