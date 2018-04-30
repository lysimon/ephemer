package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/robfig/cron"
	//"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"../../internal/paw"
	"../../internal/pkg/config"
)

func yourTaskFunction() {
	sess := session.New(&aws.Config{Region: aws.String("eu-central-1"), Credentials: credentials.NewStaticCredentials("accesskey", "secretkey", "")})
	s3Client := s3.New(sess)
	result, err := s3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Get error details
			log.Println("Error:", awsErr.Code(), awsErr.Message())

			// Prints out full error message, including original error if there was one.
			log.Println("Error:", awsErr.Error())

			// Get original error
			if origErr := awsErr.OrigErr(); origErr != nil {
				// operate on original error.
			}
		} else {
			fmt.Println(err.Error())
		}
	}

	log.Println("Buckets:")
	for _, bucket := range result.Buckets {
		log.Printf("%s : %s\n", aws.StringValue(bucket.Name), bucket.CreationDate)
	}

	log.Fatal("Cron job failed, finishing process")
}

func main() {
	// Starting cron job in the background
	config.Init("/data/config/configuration.yaml")
	c := cron.New()
	//c.AddFunc("5 * * * * *", func() { yourTaskFunction() })
	c.AddFunc("5 * * * * *", func() { paw.Cron() })
	c.Start()
	defer c.Stop()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("error")
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
