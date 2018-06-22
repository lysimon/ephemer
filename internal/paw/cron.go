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

func Cron() {
	log.Print("Refreshing all states")
	// Can be on startup if we trust people not to remove them
	validate_buckets()
	// Refreshing state of cloudformation stack + redis
	log.Print("Finish processing all states")
	//log.Print(config.GlobalConfig)
	//log.Fatal("Failing hard")

}

// Loop over all aws account + region to create the s3 bucket for cloudformation
func validate_buckets() {
	log.Print("validate_buckets")
	for _, provider := range config.GlobalConfig.Providers {
		// element is the element from someSlice for where we are
		if provider.Provider == "aws" {

			// For each region, we create the s3 bucket where cloudformation will be
			for _, region := range provider.Regions {

				// Creating credentials to call aws
				sess := config.GlobalConfig.Get_aws_session(provider.Name, region.Name)
				// Need to retrieve the correct bucket name for region
				svc := s3.New(sess)
				// TODO change api call in case too many bucket
				//buckets, err := svc.ListBuckets(nil)
				//if err != nil {
				//	log.Fatalf("Fail to list bucket: %v", err)
				//}
				//log.Printf("Found bucket list: %v", buckets)
				// Find out the bucket name for cloudformation
				bucket := parser.Get_parsed_value(region.Bucket_name, config.GlobalConfig.Get_aws_account_id(provider.Name), region.Name, "")
				log.Printf("We should verify that this bucket: %v exist in region %v", bucket, region.Name)

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
				// Setting retention policy for the created bucket

				// Fail hard if buckets cannot be describe
			}
		}

	}
	log.Print(config.GlobalConfig.Redis.Port)
}
