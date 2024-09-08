package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/datatypes/stringmaps"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/datatypes/stringmaps/convertstringmaps"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func subnets(ctx *pulumi.Context, locals *Locals, createdVpc *ec2.Vpc) (privateSubnets,
	publicSubnets []*ec2.Subnet, err error) {

	privateSubnets = make([]*ec2.Subnet, 0)
	publicSubnets = make([]*ec2.Subnet, 0)

	// iterate through azs and create the configured number of public and private subnets per az
	for availabilityZone, subnetNameCidrMap := range locals.PrivateSubnetMap {
		for subnetName, subnetCidr := range subnetNameCidrMap {
			// create private subnet
			createdSubnet, err := ec2.NewSubnet(ctx,
				string(subnetName),
				&ec2.SubnetArgs{
					VpcId:            createdVpc.ID(),
					CidrBlock:        pulumi.String(subnetCidr),
					AvailabilityZone: pulumi.String(availabilityZone),
					Tags: convertstringmaps.ConvertGoStringMapToPulumiStringMap(
						stringmaps.AddEntry(locals.AwsTags, "Name", string(subnetName))),
				}, pulumi.Parent(createdVpc))
			if err != nil {
				return nil, nil,
					errors.Wrapf(err, "error creating private subnet %s", subnetName)
			}
			privateSubnets = append(privateSubnets, createdSubnet)
		}

	}

	for availabilityZone, subnetNameCidrMap := range locals.PublicSubnetMap {
		for subnetName, subnetCidr := range subnetNameCidrMap {
			// create public subnet
			createdSubnet, err := ec2.NewSubnet(ctx,
				string(subnetName),
				&ec2.SubnetArgs{
					VpcId:            createdVpc.ID(),
					CidrBlock:        pulumi.String(subnetCidr),
					AvailabilityZone: pulumi.String(availabilityZone),
					//required for public subnets
					MapPublicIpOnLaunch: pulumi.Bool(true),
					Tags: convertstringmaps.ConvertGoStringMapToPulumiStringMap(
						stringmaps.AddEntry(locals.AwsTags, "Name", string(subnetName))),
				}, pulumi.Parent(createdVpc))
			if err != nil {
				return nil, nil, errors.Wrapf(err,
					"error creating public subnet %s", subnetName)
			}
			publicSubnets = append(publicSubnets, createdSubnet)
		}

	}
	return privateSubnets, publicSubnets, nil
}
