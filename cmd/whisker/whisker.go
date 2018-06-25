package main

import (
	"log"
	"net/http"
	//"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gorilla/mux"
	"github.com/lysimon/felicette/internal/whisker"
	"github.com/lysimon/felicette/pkg/status"
)

func main() {
	// Start web server
	router := mux.NewRouter()
	router.HandleFunc("/git/{id}", whisker.GetGitProjectFromRequest)
	router.HandleFunc("/git/{id}/branch/{branch}", whisker.Check_git_branch)
	router.HandleFunc("/status", status.Status)

	log.Fatal(http.ListenAndServe(":8080", router))

}
