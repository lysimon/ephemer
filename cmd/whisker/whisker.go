package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	//"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"../../internal/whisker"
	"../../pkg/status"
)

func main() {
	// Start web server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("error")
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	whisker.Configuration()

	// Loading status
	status.Status()

	log.Fatal(http.ListenAndServe(":8080", nil))

}
