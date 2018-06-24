package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/robfig/cron"
	//"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/lysimon/felicette/internal/paw"
	"github.com/lysimon/felicette/pkg/status"
)

func main() {
	// Initialization, creating s3 buckets with retention policy for the future
	paw.Create_buckets()

	// Start clean up cron job => should be done by tail
	c := cron.New()
	c.AddFunc("5 * * * * *", func() { paw.Cron() })
	c.Start()
	defer c.Stop()

	status.Status()
	// Start web server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("error")
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
