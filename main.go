package main

import (
	"encoding/json"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const INTERNET_CIRD = "0.0.0.0/0"

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
			Group:     group.Name,
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

		policyTmpJson, policyTmpJsonErr := json.Marshal(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				map[string]interface{}{
					"Effect":   "Allow",
					"Resource": "*",
					"Action": []string{
						"ec2:CreateTags",
						"ec2:CreateVolume",
						"ec2:RunInstances",
					},
				},
			},
		})

		if policyTmpJsonErr != nil {
			return policyTmpJsonErr
		}

		jsonPolisyStr := string(policyTmpJson)
		customEc2DevPolicy, newPolicyErr := iam.NewPolicy(ctx, "policy", &iam.PolicyArgs{
			Name:        pulumi.String("ec2-developers-access"),
			Path:        pulumi.String("/"),
			Description: pulumi.String("My test policy"),
			Policy:      pulumi.String(jsonPolisyStr),
		})

		if newPolicyErr != nil {
			return newPolicyErr
		}

		_, customGroupPolicyAttachmentErr := iam.NewGroupPolicyAttachment(ctx, "developers-custom-group-policy-attachment", &iam.GroupPolicyAttachmentArgs{
			Group:     group.Name,
			PolicyArn: customEc2DevPolicy.Arn,
		})

		if customGroupPolicyAttachmentErr != nil {
			return customGroupPolicyAttachmentErr
		}

		ctx.Export("userId", user.ID())

		createRoleErr := createRole(ctx)
		if createRoleErr != nil {
			return createRoleErr
		}

		createVPCErr := createVPC(ctx)
		if createVPCErr != nil {
			return createVPCErr
		}

		return nil
	})
}

func createRole(ctx *pulumi.Context) error {
	tmpJSON0, err := json.Marshal(map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			map[string]interface{}{
				"Action": "sts:AssumeRole",
				"Effect": "Allow",
				"Sid":    "",
				"Principal": map[string]interface{}{
					"Service": "ec2.amazonaws.com",
				},
			},
		},
	})

	if err != nil {
		return err
	}

	json0 := string(tmpJSON0)
	_, roleErr := iam.NewRole(ctx, "developers", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(json0),
		ManagedPolicyArns: pulumi.StringArray{
			pulumi.String("arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"),
		},
		Tags: pulumi.StringMap{
			"tag-key": pulumi.String("tag-value"),
		},
	})

	if roleErr != nil {
		return roleErr
	}

	return nil
}

func createVPC(ctx *pulumi.Context) error {
	vpc, err := ec2.NewVpc(ctx, "main", &ec2.VpcArgs{
		Tags: pulumi.StringMap{
			"Name": pulumi.String("custom-test-vpc"),
		},
		CidrBlock: pulumi.String("10.2.0.0/16"),
	})
	if err != nil {
		return err
	}

	publicASubnet, subnetPublicAErr := ec2.NewSubnet(ctx, "publicA", &ec2.SubnetArgs{
		VpcId:            vpc.ID(),
		AvailabilityZone: pulumi.String("us-east-1a"),
		CidrBlock:        pulumi.String("10.2.0.0/24"),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("PublicA"),
		},
	})
	if subnetPublicAErr != nil {
		return subnetPublicAErr
	}

	publicBSubnet, subnetPublicBErr := ec2.NewSubnet(ctx, "publicB", &ec2.SubnetArgs{
		VpcId:            vpc.ID(),
		AvailabilityZone: pulumi.String("us-east-1b"),
		CidrBlock:        pulumi.String("10.2.1.0/24"),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("PublicB"),
		},
	})
	if subnetPublicBErr != nil {
		return subnetPublicAErr
	}

	gw, internetGatewayErr := ec2.NewInternetGateway(ctx, "custom-internet-gw", &ec2.InternetGatewayArgs{
		VpcId: vpc.ID(),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("custom-gw"),
		},
	})
	if internetGatewayErr != nil {
		return err
	}

	routeTable, routeTableErr := ec2.NewRouteTable(ctx, "public-traffic", &ec2.RouteTableArgs{
		VpcId: vpc.ID(),
		Routes: ec2.RouteTableRouteArray{
			&ec2.RouteTableRouteArgs{
				CidrBlock: pulumi.String(INTERNET_CIRD),
				GatewayId: gw.ID(),
			},
		},
		Tags: pulumi.StringMap{
			"Name": pulumi.String("public-traffic"),
		},
	})
	if routeTableErr != nil {
		return routeTableErr
	}

	_, routeTableAssociationAErr := ec2.NewRouteTableAssociation(ctx, "routeTableAssociationA", &ec2.RouteTableAssociationArgs{
		SubnetId:     publicASubnet.ID(),
		RouteTableId: routeTable.ID(),
	})
	if routeTableAssociationAErr != nil {
		return routeTableAssociationAErr
	}

	_, routeTableAssociationBErr := ec2.NewRouteTableAssociation(ctx, "routeTableAssociationB", &ec2.RouteTableAssociationArgs{
		SubnetId:     publicBSubnet.ID(),
		RouteTableId: routeTable.ID(),
	})
	if routeTableAssociationBErr != nil {
		return routeTableAssociationAErr
	}

	_, aclErr := ec2.NewNetworkAcl(ctx, "custom-acl", &ec2.NetworkAclArgs{
		VpcId:     vpc.ID(),
		Tags:      pulumi.StringMap{"Name": pulumi.String("custom-acl")},
		SubnetIds: pulumi.StringArray{publicASubnet.ID(), publicBSubnet.ID()},
		Ingress: ec2.NetworkAclIngressArray{
			&ec2.NetworkAclIngressArgs{
				RuleNo:    pulumi.Int(100),
				Protocol:  pulumi.String("tcp"),
				Action:    pulumi.String("allow"),
				FromPort:  pulumi.Int(80),
				ToPort:    pulumi.Int(80),
				CidrBlock: pulumi.String(INTERNET_CIRD),
			},
			&ec2.NetworkAclIngressArgs{
				RuleNo:    pulumi.Int(101),
				Protocol:  pulumi.String("tcp"),
				Action:    pulumi.String("allow"),
				FromPort:  pulumi.Int(443),
				ToPort:    pulumi.Int(443),
				CidrBlock: pulumi.String(INTERNET_CIRD),
			},
		},
		Egress: ec2.NetworkAclEgressArray{
			&ec2.NetworkAclEgressArgs{
				Protocol:  pulumi.String("tcp"),
				RuleNo:    pulumi.Int(100),
				Action:    pulumi.String("allow"),
				FromPort:  pulumi.Int(1024),
				ToPort:    pulumi.Int(65535),
				CidrBlock: pulumi.String(INTERNET_CIRD),
			},
		},
	})
	if aclErr != nil {
		return aclErr
	}

	return nil
}
