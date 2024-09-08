package pkg

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsvpc"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsVpc *awsvpc.AwsVpc
	Labels map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsvpc.AwsVpcStackInput) *Locals {
	locals := &Locals{}

	//assign value for the locals variable to make it available across the project
	locals.AwsVpc = stackInput.ApiResource

	return locals
}
