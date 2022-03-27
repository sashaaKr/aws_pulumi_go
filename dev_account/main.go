package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		deployIamErr := deployIAM(ctx)

		if deployIamErr != nil {
			return deployIamErr
		}

		return nil
	})
}

func deployIAM(ctx *pulumi.Context) error {
	instance_assume_role_policy, err := iam.GetPolicyDocument(ctx, &iam.GetPolicyDocumentArgs{
		Statements: []iam.GetPolicyDocumentStatement{
			iam.GetPolicyDocumentStatement{
				Actions: []string{
					"sts:AssumeRole",
				},
				Principals: []iam.GetPolicyDocumentStatementPrincipal{
					iam.GetPolicyDocumentStatementPrincipal{
						Type: "AWS",
						Identifiers: []string{
							"044141213750",
						},
					},
				},
			},
		},
	}, nil)

	if err != nil {
		return err
	}

	_, err = iam.NewRole(ctx, "cross-account-access", &iam.RoleArgs{
		ManagedPolicyArns: pulumi.StringArray{
			pulumi.String("arn:aws:iam::aws:policy/AdministratorAccess"),
		},
		AssumeRolePolicy: pulumi.String(instance_assume_role_policy.Json),
	})

	if err != nil {
		return err
	}
	
	return nil

}
