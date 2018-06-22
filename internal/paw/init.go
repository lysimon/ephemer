package paw

import (
	"context"
	"fmt"
	"log"
	"os"

	"../../pkg/config"
	"../../pkg/parser"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func Create_buckets() {
	log.Print("Create_buckets start")
	// Can be on startup if we trust people not to remove them
	create_aws_buckets()
	// Refreshing state of cloudformation stack + redis
	log.Print("Create_buckets end")
	//log.Print(config.GlobalConfig)
	//log.Fatal("Failing hard")

}

// Loop over all aws account + region to create the s3 bucket for cloudformation
func create_aws_buckets() {
	log.Print("create_aws_buckets")
	for _, provider := range config.GlobalConfig.Providers {
		// element is the element from someSlice for where we are
		if provider.Provider == "aws" {

			// For each region, we create the s3 bucket where cloudformation will be
			for _, region := range provider.Regions {

				// Creating credentials to call aws
				sess := config.GlobalConfig.Get_aws_session(provider.Name, region.Name)
				// Need to retrieve the correct bucket name for region
				svc := s3.New(sess)

				bucket := parser.Get_parsed_value(region.Bucket_name, config.GlobalConfig.Get_aws_account_id(provider.Name), region.Name, "")
				log.Printf("We should verify that the bucket: %v exist in region %v", bucket, region.Name)

				ctx := context.Background()
				bucket_region, err := s3manager.GetBucketRegion(ctx, sess, bucket, region.Name)
				if err != nil {
					if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
						fmt.Fprintf(os.Stderr, "unable to find bucket %s's region not found\n", bucket)
						// We should create the s3 bucket then, and fail if we cannot create it
						fmt.Fprintf(os.Stderr, "Creating bucket %s in region %s\n", bucket, bucket_region)
						_, err = svc.CreateBucket(&s3.CreateBucketInput{
							Bucket: aws.String(bucket),
						})
						if err != nil {
							log.Fatalf("Unable to create bucket %q, %v", bucket, err)
						}
					} else {
						log.Fatalf("Unknow error happened while trying to list bucket %v: %v", bucket, err)
					}
				} else {

					if bucket_region == region.Name {
						fmt.Printf("Bucket %s is in correct %s region\n", bucket, bucket_region)
					} else {
						log.Fatalf("Bucket %s should be in region %s but is in %v, please, delete it\n", bucket, region.Name, bucket_region)
					}
				}
				// Verifying retention policy of bucket
				_, err = svc.GetBucketLifecycleConfiguration(&s3.GetBucketLifecycleConfigurationInput{
					Bucket: aws.String(bucket),
				})

				if err != nil {
					log.Printf("code: %s", err.(awserr.Error).Code())
					if err.(awserr.Error).Code() == "NoSuchLifecycleConfiguration" {
						status := "Enabled"
						id := "ephemerio"
						prefix := "*"
						// Generate default lifecycle rule
						LifecycleRule := s3.LifecycleRule{
							Status: &status,
							Expiration: &s3.LifecycleExpiration{
								Days: &region.Retention_day},
							ID: &id,
							Filter: &s3.LifecycleRuleFilter{
								Prefix: &prefix,
							},
						}

						_, err = svc.PutBucketLifecycleConfiguration(&s3.PutBucketLifecycleConfigurationInput{
							Bucket:                 aws.String(bucket),
							LifecycleConfiguration: &s3.BucketLifecycleConfiguration{Rules: []*s3.LifecycleRule{&LifecycleRule}},
						})
						if err != nil {
							log.Fatalf("Unable to create lifecycle policy %q, %v", bucket, err.(awserr.Error))

						}

					} else {
						log.Fatalf("Unable to get bucket lifecycle policy %q, %v", bucket, err.(awserr.Error))
					}
				}
				// Setting retention policy for the created bucket

				// Fail hard if buckets cannot be describe
			}
		}

	}
	log.Print(config.GlobalConfig.Redis.Port)
}
