package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/secretmanager"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/sql"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		const region, project = "europe-west1", "conro-sbx"

		// compute.LookupNetwork
		vpcNetwork, err := compute.LookupNetwork(ctx, &compute.LookupNetworkArgs{
			Name:    "vpc-conro-sbx",
			Project: pulumi.StringRef(project),
		}, nil)
		if err != nil {
			return err
		}

		// compute.LookupSubnetwork
		vpcSubnet, err := compute.LookupSubnetwork(ctx, &compute.LookupSubnetworkArgs{
			Name:    pulumi.StringRef("subnet-conro-sbx"),
			Project: pulumi.StringRef(project),
			Region:  pulumi.StringRef(region),
		})
		if err != nil {
			return err
		}

		// compute.NewAddress
		ipAddress, err := compute.NewAddress(ctx, "internal-db-ipaddress", &compute.AddressArgs{
			Name:        pulumi.String("cloud-sql-psc"),
			AddressType: pulumi.String("INTERNAL"),
			Region:      pulumi.String(region),
			Subnetwork:  pulumi.String(vpcSubnet.Id),
		})
		if err != nil {
			return err
		}

		// Private VPC Connection was already implemented, so we've skipped that

		// sql.DatabaseInstance
		dbInstance, err := sql.NewDatabaseInstance(
			ctx,
			"cloudsql-instance",
			&sql.DatabaseInstanceArgs{
				Name:            pulumi.String("cnr-instance-20250703"),
				Region:          pulumi.String(region),
				DatabaseVersion: pulumi.String("POSTGRES_17"),
				Settings: &sql.DatabaseInstanceSettingsArgs{
					Tier: pulumi.String("db-f1-micro"),
					IpConfiguration: &sql.DatabaseInstanceSettingsIpConfigurationArgs{
						Ipv4Enabled: pulumi.Bool(
							false,
						), // Enabling this, creates a Public IpAddress
						PscConfigs: sql.DatabaseInstanceSettingsIpConfigurationPscConfigArray{
							&sql.DatabaseInstanceSettingsIpConfigurationPscConfigArgs{
								PscEnabled: pulumi.Bool(true), // To PSC Data
								AllowedConsumerProjects: pulumi.StringArray{
									pulumi.String(project),
								},
							},
						},
						PrivateNetwork:                          pulumi.String(vpcNetwork.SelfLink),
						EnablePrivatePathForGoogleCloudServices: pulumi.Bool(true),
					},
				},
				DeletionProtection: pulumi.Bool(false),
			},
		)
		if err != nil {
			return err
		}

		// compute.ForwardingRule
		forwardingRule, err := compute.NewForwardingRule(
			ctx,
			"psc-cloud-sql",
			&compute.ForwardingRuleArgs{
				Name:                ipAddress.Name,
				IpAddress:           ipAddress.SelfLink,
				LoadBalancingScheme: pulumi.String(""),
				Project:             pulumi.String(project),
				Region:              pulumi.String(region),
				Network:             pulumi.String("vpc-conro-sbx"),
				Target:              dbInstance.PscServiceAttachmentLink,
			},
		)
		if err != nil {
			return err
		}

		// sql.Database
		_, err = sql.NewDatabase(ctx, "cnr-database", &sql.DatabaseArgs{
			Name:     pulumi.String("cnr-app-db"),
			Instance: dbInstance.Name,
		})
		if err != nil {
			return err
		}

		// Random password
		password, err := random.NewRandomPassword(ctx, "app-db-pass", &random.RandomPasswordArgs{
			Length:  pulumi.Int(16),
			Special: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		// sql.User
		user, err := sql.NewUser(ctx, "app-user", &sql.UserArgs{
			Name:     pulumi.String("appuser"),
			Instance: dbInstance.Name,
			Password: password.Result, // Create random, store as secret
		})
		if err != nil {
			return err
		}

		// secretmanager.NewSecret
		secretName := "cnr-database"
		secret, err := secretmanager.NewSecret(ctx, secretName, &secretmanager.SecretArgs{
			SecretId: pulumi.String("creds-cnr-database"),
			Replication: &secretmanager.SecretReplicationArgs{
				Auto: &secretmanager.SecretReplicationAutoArgs{},
			},
		})
		if err != nil {
			panic(err)
		}

		secretData := pulumi.All(
			user.Name,
			password.Result,
		).ApplyT(func(args []interface{}) string {
			user := args[0].(string)
			password := args[1].(string)

			return fmt.Sprintf(
				"user: %s\npassword: %s",
				user, password,
			)
		}).(pulumi.StringOutput)

		// secretmanager.NewSecretVersion
		_, err = secretmanager.NewSecretVersion(
			ctx,
			secretName+"-version",
			&secretmanager.SecretVersionArgs{
				Enabled:    pulumi.Bool(true),
				Secret:     secret.ID(),
				SecretData: secretData,
			},
		)
		if err != nil {
			panic(err)
		}

		ctx.Export("DB Name:", dbInstance.Name)
		ctx.Export("DB IP:", ipAddress.Address)
		ctx.Export("DB Connection:", dbInstance.ConnectionName)
		ctx.Export("DB Password:", password.Result)
		ctx.Export("ForwardingRule:", forwardingRule.CreationTimestamp)

		return nil
	})
}
