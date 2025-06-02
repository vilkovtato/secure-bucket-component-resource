package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// SecureBucketArgs contains the arguments for creating a SecureBucket
type SecureBucketArgs struct {
	BucketName pulumi.StringInput            `pulumi:"bucketName,optional"`
	Versioning *bool                         `pulumi:"versioning,optional"`
	Encryption *bool                         `pulumi:"encryption,optional"`
	Tags       map[string]pulumi.StringInput `pulumi:"tags,optional"`
}

// SecureBucket is a component that creates an S3 bucket with security best practices
type SecureBucket struct {
	pulumi.ResourceState

	BucketName pulumi.StringOutput `pulumi:"bucketName"`
}

// NewSecureBucket creates a new SecureBucket component
func NewSecureBucket(ctx *pulumi.Context, name string, args SecureBucketArgs, opts ...pulumi.ResourceOption) (*SecureBucket, error) {
	component := &SecureBucket{}
	err := ctx.RegisterComponentResource("custom:module:SecureBucket", name, component, opts...)
	if err != nil {
		return nil, err
	}

	// Create an S3 bucket with best practices by default
	tags := pulumi.StringMap{
		"ManagedBy": pulumi.String("Pulumi"),
	}
	if args.Tags != nil {
		for k, v := range args.Tags {
			tags[k] = v
		}
	}

	bucket, err := s3.NewBucketV2(ctx, name, &s3.BucketV2Args{
		Bucket: args.BucketName,
		Tags:   tags,
	}, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}

	// Conditionally enable versioning
	if args.Versioning != nil && !*args.Versioning {
		// Skip versioning
	} else {
		_, err = s3.NewBucketVersioningV2(ctx, name+"-versioning", &s3.BucketVersioningV2Args{
			Bucket: bucket.ID(),
			VersioningConfiguration: &s3.BucketVersioningV2VersioningConfigurationArgs{
				Status: pulumi.String("Enabled"),
			},
		}, pulumi.Parent(component))
		if err != nil {
			return nil, err
		}
	}

	// Conditionally enable encryption
	if args.Encryption != nil && !*args.Encryption {
		// Skip encryption
	} else {
		_, err = s3.NewBucketServerSideEncryptionConfigurationV2(ctx, name+"-encryption", &s3.BucketServerSideEncryptionConfigurationV2Args{
			Bucket: bucket.ID(),
			Rules: s3.BucketServerSideEncryptionConfigurationV2RuleArray{
				&s3.BucketServerSideEncryptionConfigurationV2RuleArgs{
					ApplyServerSideEncryptionByDefault: &s3.BucketServerSideEncryptionConfigurationV2RuleApplyServerSideEncryptionByDefaultArgs{
						SseAlgorithm: pulumi.String("AES256"),
					},
				},
			},
		}, pulumi.Parent(component))
		if err != nil {
			return nil, err
		}
	}

	component.BucketName = bucket.ID().ToStringOutput()
	return component, nil
}