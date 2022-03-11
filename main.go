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

		group, errGroup := iam.NewGroup(ctx, "developers", &iam.GroupArgs{
			Path: pulumi.String("/users/"),
		})

		if errGroup != nil {
			return errGroup
		}

    _, groupPolicyAttachmentErr := iam.NewGroupPolicyAttachment(ctx, "developers-group-policy-attachment", &iam.GroupPolicyAttachmentArgs{
      Group: group.Name,
      PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEC2FullAccess"),
    })

    if groupPolicyAttachmentErr != nil {
      return groupPolicyAttachmentErr
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
		
		ctx.Export("userId", user.ID())

		return nil
	})
}
