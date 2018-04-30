package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	providers []struct {
		provider    string
		name        string
		credentials struct {
			access_key string
			secret_key string
		}
	}

	redis struct {
		host string
		port int16
	}

	paw struct {
		host string
		port int16
	}
}

func Init(file_path string) {

	log.Print("Initializing configuration")
	conf := make(map[interface{}]interface{})

	reader, _ := os.Open(file_path)

	buff, _ := ioutil.ReadAll(reader)
	s := string(buff)
	fmt.Print(s)

	err := yaml.Unmarshal(buff, &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("--- m:\n%v\n\n", conf)
}
