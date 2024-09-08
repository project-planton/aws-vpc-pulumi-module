package pkg

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func subnets(ctx *pulumi.Context, locals *Locals, createdVpc *ec2.Vpc) (publicSubnets,
	privateSubnets []*ec2.Subnet, err error) {

	privateSubnets = make([]*ec2.Subnet, 0)
	publicSubnets = make([]*ec2.Subnet, 0)

	// iterate through azs and create the configured number of public and private subnets per az
	for azIndex, az := range locals.AwsVpc.Spec.AvailabilityZones {
		for subnetIndex := 0; subnetIndex < int(locals.AwsVpc.Spec.SubnetsPerAvailabilityZone); subnetIndex++ {
			// Calculate the public subnet CIDR block
			// Public Subnet CIDR Calculation
			publicSubnetCidr := fmt.Sprintf("10.0.%d.0/%d", azIndex*10+subnetIndex, locals.AwsVpc.Spec.SubnetSize)

			// Private Subnet CIDR Calculation
			privateSubnetCidr := fmt.Sprintf("10.0.%d.0/%d", 100+azIndex*10+subnetIndex, locals.AwsVpc.Spec.SubnetSize)

			// Public Subnet
			createdPublicSubnet, err := ec2.NewSubnet(ctx,
				fmt.Sprintf("publicSubnet-%d-%d", azIndex, subnetIndex),
				&ec2.SubnetArgs{
					VpcId:               createdVpc.ID(),
					CidrBlock:           pulumi.String(publicSubnetCidr),
					AvailabilityZone:    pulumi.String(az),
					MapPublicIpOnLaunch: pulumi.Bool(true),
					Tags: pulumi.StringMap{
						"Name": pulumi.String(fmt.Sprintf("public-subnet-%d-%d", azIndex, subnetIndex)),
					},
				}, pulumi.Parent(createdVpc))
			if err != nil {
				return nil, nil, errors.Wrap(err, "error creating public subnet")
			}

			publicSubnets = append(publicSubnets, createdPublicSubnet)

			// Private Subnet
			createdPrivateSubnet, err := ec2.NewSubnet(ctx,
				fmt.Sprintf("privateSubnet-%d-%d", azIndex, subnetIndex),
				&ec2.SubnetArgs{
					VpcId:            createdVpc.ID(),
					CidrBlock:        pulumi.String(privateSubnetCidr),
					AvailabilityZone: pulumi.String(az),
					Tags: pulumi.StringMap{
						"Name": pulumi.String(fmt.Sprintf("private-subnet-%d-%d", azIndex, subnetIndex)),
					},
				}, pulumi.Parent(createdVpc))
			if err != nil {
				return nil, nil, errors.Wrap(err, "error creating private subnet")
			}

			privateSubnets = append(privateSubnets, createdPrivateSubnet)
		}
	}
	return publicSubnets, privateSubnets, nil
}
