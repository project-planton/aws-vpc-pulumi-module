package pkg

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud/aws-vpc-pulumi-module/pkg/outputs"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func natGateways(ctx *pulumi.Context, awsProvider *aws.Provider, createdVpc *ec2.Vpc, createdPrivateSubnets []*ec2.Subnet) error {
	// create nat gateways for private subnets
	for i, createdPrivateSubnet := range createdPrivateSubnets {
		//create elastic ip for nat gateway
		createdElasticIp, err := ec2.NewEip(ctx,
			fmt.Sprintf("natEip-%d", i),
			&ec2.EipArgs{}, pulumi.Provider(awsProvider))
		if err != nil {
			return errors.Wrap(err, "error creating eip for nat gateway")
		}

		//create nat gateway
		createdNatGateway, err := ec2.NewNatGateway(ctx,
			fmt.Sprintf("natGateway-%d", i),
			&ec2.NatGatewayArgs{
				SubnetId:     createdPrivateSubnet.ID(),
				AllocationId: createdElasticIp.ID(),
			}, pulumi.Parent(createdPrivateSubnet))
		if err != nil {
			return errors.Wrap(err, "error creating nat gateway")
		}

		// Extract and export the 'Name' tag from the subnet using Apply
		createdPrivateSubnet.Tags.ApplyT(func(tags map[string]string) (string, error) {
			if nameTag, ok := tags["Name"]; ok {
				ctx.Export(outputs.NatGatewayIdOutputKey(nameTag), createdNatGateway.ID())
				ctx.Export(outputs.NatGatewayPublicIpOutputKey(nameTag), createdNatGateway.PublicIp)
				ctx.Export(outputs.NatGatewayPrivateIpOutputKey(nameTag), createdNatGateway.PrivateIp)
				return nameTag, nil
			}
			return "No Name Tag", nil
		})

		// private route table to route traffic through nat gateway
		createdPrivateRouteTable, err := ec2.NewRouteTable(ctx,
			fmt.Sprintf("privateRouteTable-%d", i),
			&ec2.RouteTableArgs{
				VpcId: createdVpc.ID(),
				Routes: ec2.RouteTableRouteArray{
					&ec2.RouteTableRouteArgs{
						CidrBlock:    pulumi.String("0.0.0.0/0"),
						NatGatewayId: createdNatGateway.ID(),
					},
				},
			}, pulumi.Parent(createdVpc))
		if err != nil {
			return errors.Wrap(err, "error creating private route table")
		}

		// associate private route table with private subnets
		_, err = ec2.NewRouteTableAssociation(ctx,
			fmt.Sprintf("privateRouteAssoc-%d", i),
			&ec2.RouteTableAssociationArgs{
				SubnetId:     createdPrivateSubnet.ID(),
				RouteTableId: createdPrivateRouteTable.ID(),
			}, pulumi.Parent(createdPrivateRouteTable))
		if err != nil {
			return errors.Wrap(err, "error associating private route table")
		}
	}
	return nil
}
