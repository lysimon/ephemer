package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Providers []struct {
		Provider    string
		Name        string
		Credentials struct {
			Access_key string
			Secret_key string
		}
		Regions []struct {
			Name          string
			Bucket_name   string
			Retention_day int64
		}
	}

	Redis struct {
		Host string
		Port int16
	}

	Paw struct {
		Host string
		Port int16
	}
}

func (c Config) To_json() []byte {
	js, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("Unable to retrieve configuration: %v", err)
	}
	return js
}

func (c Config) Get_redis_host() int16 {
	return c.Redis.Port
}

// Get a session for the aws account defined by this name and region
// TODO add cache here for 30 min
func (c Config) Get_aws_session(name string, region string) *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: c.Get_aws_credentials(name),
	})
	// Fail hard on err
	if err != nil {
		log.Fatalf("Fail to create aws session: %v", err)
	}

	return sess
}

func (c Config) Get_aws_credentials(name string) *credentials.Credentials {
	for _, provider := range c.Providers {
		// element is the element from someSlice for where we are
		print(provider.Name)
		if provider.Name == name {
			// Only static credentials managed for now
			return credentials.NewStaticCredentials(provider.Credentials.Access_key, provider.Credentials.Secret_key, "")
		}
	}
	return nil
}

// Retrieve the account id, can be used for creating unique bucket
func (c Config) Get_aws_account_id(name string) string {
	// using default region us-east-1
	sess := c.Get_aws_session(name, "us-east-1")
	svc := sts.New(sess)
	req, resp := svc.GetCallerIdentityRequest(nil)
	err := req.Send()
	if err != nil {
		log.Fatalf("Fail to retrieve identity: %v", err)
	}
	log.Printf("Got successful account answer: %v", resp)

	log.Printf("Got successful account: %v", *resp.Account)
	log.Printf("Got successful arn: %v", *resp.Arn)

	return *resp.Account
}

var GlobalConfig = GetConfiguration("/data/config/configuration.yaml")

// Initialize the configuration
func GetConfiguration(file_path string) Config {

	var c Config

	yamlFile, err := ioutil.ReadFile(file_path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	fmt.Printf("--- m:\n%v\n\n", string(yamlFile))

	fmt.Printf("Got final configuration")
	fmt.Printf("--- m:\n%v\n\n", c)
	fmt.Printf("End final configuration")
	return c
}
