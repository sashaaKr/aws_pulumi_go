package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		user, userErr := iam.NewUser(ctx, "my-user", &iam.UserArgs{
			Path: pulumi.String("/system/"),
			Tags: pulumi.StringMap{
				"orelly-course": pulumi.String("true"),
			},
		})

		if userErr != nil {
			return userErr
		}

		ctx.Export("userId", user.ID())

		group, errGroup := iam.NewGroup(ctx, "developers", &iam.GroupArgs{
			Path: pulumi.String("/users/"),
		})

		if errGroup != nil {
			return errGroup
		}

		_, membershipErr := iam.NewGroupMembership(ctx, "team", &iam.GroupMembershipArgs{
			Group: group.Name,
			Users: pulumi.StringArray{
				user.Name,
			},
		})

		if membershipErr != nil {
			return membershipErr
		}

    // TODO: add ec2 full access policy to dev group
		
		return nil
	})
}
