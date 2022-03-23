package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an AWS resource (S3 Bucket)
		bucket, err := s3.NewBucket(ctx, "my-bucket", nil)
		if err != nil {
			return err
		}

		// Export the name of the bucket
		ctx.Export("bucketName", bucket.ID())
		return nil
	})
}

func deployIAM(ctx *pulumi.Context) error {
	_, roleErr := iam.NewRole(ctx, "developers", &iam.RoleArgs{
		ManagedPolicyArns: pulumi.StringArray{
			pulumi.String("arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"),
		},
		Tags: pulumi.StringMap{
			"tag-key": pulumi.String("tag-value"),
		},
	})

	if (roleErr != nil) {
		return roleErr
	}
	return nil
}
