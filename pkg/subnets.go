package pkg

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/project-planton/aws-vpc-pulumi-module/pkg/localz"
	"github.com/project-planton/aws-vpc-pulumi-module/pkg/outputs"
	"github.com/project-planton/pulumi-module-golang-commons/pkg/datatypes/stringmaps"
	"github.com/project-planton/pulumi-module-golang-commons/pkg/datatypes/stringmaps/convertstringmaps"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func subnets(ctx *pulumi.Context, locals *localz.Locals, createdVpc *ec2.Vpc,
	createdPublicRouteTable *ec2.RouteTable) error {
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
				return errors.Wrapf(err, "error creating private subnet %s", subnetName)
			}
			ctx.Export(outputs.SubnetIdOutputKey(subnetName), createdSubnet.ID())
			ctx.Export(outputs.SubnetCidrOutputKey(subnetName), createdSubnet.CidrBlock)

			if locals.AwsVpc.Spec.IsNatGatewayEnabled {
				if err := natGateway(ctx, locals, createdVpc, subnetName, createdSubnet); err != nil {
					return errors.Wrapf(err, "failed to create nat-gateway for %s subnet", subnetName)
				}
			}
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
				return errors.Wrapf(err, "error creating public subnet %s", subnetName)
			}

			ctx.Export(outputs.SubnetIdOutputKey(subnetName), createdSubnet.ID())
			ctx.Export(outputs.SubnetCidrOutputKey(subnetName), createdSubnet.CidrBlock)

			_, err = ec2.NewRouteTableAssociation(ctx,
				fmt.Sprintf("public-route-assoc-%s", subnetName),
				&ec2.RouteTableAssociationArgs{
					RouteTableId: createdPublicRouteTable.ID(),
					SubnetId:     createdSubnet.ID(),
				}, pulumi.Parent(createdPublicRouteTable))
			if err != nil {
				return errors.Wrap(err, "error associating route table with public subnet")
			}
		}
	}
	return nil
}
