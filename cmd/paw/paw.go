package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/robfig/cron"
	//"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"../../internal/paw"
)

func main() {
	// Initialization, creating Create_buckets
	paw.Create_buckets()

	// Start clean up cron job
	c := cron.New()
	c.AddFunc("5 * * * * *", func() { paw.Cron() })
	c.Start()
	defer c.Stop()

	// Start web server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("error")
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
