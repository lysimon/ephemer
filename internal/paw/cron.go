package paw

import (
	"log"

	"../../pkg/config"
	"github.com/aws/aws-sdk-go/service/s3"
)

func Cron() {
	log.Print("Refreshing all states")
	// Can be on startup if we trust people not to remove them
	validate_buckets()
	// Refreshing state of cloudformation stack + redis
	log.Print("Finish processing all states")
	log.Print(config.GlobalConfig)
	//log.Fatal("Failing hard")

}

// Loop over all aws account + region to create s3 bucket
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
				buckets, err := svc.ListBuckets(nil)
				if err != nil {
					log.Fatalf("Fail to list bucket: %v", err)
				}
				log.Printf("Found bucket list: %v", buckets)
				config.GlobalConfig.Get_aws_account_id(provider.Name)
				// Fail hard if buckets cannot be describe
			}
		}

	}
	log.Print(config.GlobalConfig.Redis.Port)
}
