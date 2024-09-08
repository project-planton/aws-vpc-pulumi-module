package outputs

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsvpc"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

func PulumiOutputsToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *awsvpc.AwsVpcStackInput) *awsvpc.AwsVpcStackOutputs {
	return &awsvpc.AwsVpcStackOutputs{}
}
