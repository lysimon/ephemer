package whisker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
)

/*
func TestGit(t *testing.T) {
	result := Is_git_rep("https://github.com/lysimon/hello-go-serverless-webapp")
	if !result {
		t.Errorf("expected true, got false")
	}
	result = Is_git_rep("https://github.com/lysimon/unknownrepo")
	if result {
		t.Errorf("Was able to access git url")
	}
	result = Is_git_rep("trololo")
	if result {
		t.Errorf("Was able to access git url")
	}
	result = Is_valid_git_branch("https://github.com/lysimon/hello-go-serverless-webapp", "master")
	if !result {
		t.Errorf("expected true, got false")
	}
	result = Is_valid_git_branch("https://github.com/lysimon/hello-go-serverless-webapp", "stable")
	if !result {
		t.Errorf("expected true, got false")
	}
	result = Is_valid_git_branch("https://github.com/lysimon/hello-go-serverless-webapp", "masterddd")
	if result {
		t.Errorf("expected false, got true")
	}
	log.Print("TestGit finish")

}

func TestGitFile(t *testing.T) {
	result := Clone_git_repo("https://github.com/lysimon/hello-go-serverless-webapp")
	if !result {
		t.Errorf("expected true, got false")
	}
	result = Clone_git_repo("https://github.com/lysimon/hello-go-serverless-webapp")
	if !result {
		t.Errorf("expected true, got false")
	}

	result = Clone_git_repo("wrongurlexpectfalseresult")
	if result {
		t.Errorf("expected false, got true")
	}
}

func TestGitFetchFile(t *testing.T) {
	result := Fetch_latest("https://github.com/lysimon/hello-go-serverless-webapp")
	if !result {
		t.Errorf("expected true, got false")
	}
}

func TestGitGetFile(t *testing.T) {
	result := Get_file("https://github.com/lysimon/hello-go-serverless-webapp", "origin/master", "felicette.yml")
	if result == nil {
		t.Errorf("Expected content")
	}
	str_result := string(result)

	if str_result != "samplefile" {
		t.Errorf("Wrong content, got %v instead of %v", str_result, "samplefile")
	}
}
*/

func TestGetGitProjectValidUrl(t *testing.T) {
	result, err := GetGitProject("aHR0cHM6Ly9naXRodWIuY29tL2x5c2ltb24vaGVsbG8tZ28tc2VydmVybGVzcy13ZWJhcHA=")
	if err != (jsonapi.ErrorObject{}) {
		t.Errorf("Not expected error during ls-remote operation")
	}
	if result.Base64Url != "aHR0cHM6Ly9naXRodWIuY29tL2x5c2ltb24vaGVsbG8tZ28tc2VydmVybGVzcy13ZWJhcHA=" {
		t.Errorf("Base 64 should be the same value")
	}
	if result.Url != "https://github.com/lysimon/hello-go-serverless-webapp" {
		t.Errorf("Url should have been sent correctly")
	}

}

func TestGetGitProjectInvalidUrl(t *testing.T) {
	_, err := GetGitProject("https://github.com/lysimon/hello-go-serverless-webapp")
	if err == (jsonapi.ErrorObject{}) {
		t.Errorf("Expected error with wrong input")
	}
	_, err = GetGitProject("aHR0cHM6Ly9yYW5kb21mYWxzZS93aHRldmVydXNlbGVzc3RoaW5nLmdpdA==")
	if err == (jsonapi.ErrorObject{}) {
		t.Errorf("Expected error with wrong input")
	}
	_, err = GetGitProject("")
	if err == (jsonapi.ErrorObject{}) {
		t.Errorf("Expected error with wrong input")
	}
	_, err = GetGitProject("  fsd . ds")
	if err == (jsonapi.ErrorObject{}) {
		t.Errorf("Expected error with wrong input")
	}
}

func TestGetGitProjectFromRequestSuccess(t *testing.T) {
	rr := httptest.NewRecorder()
	path := fmt.Sprintf("/git/%s", "aHR0cHM6Ly9naXRodWIuY29tL2x5c2ltb24vaGVsbG8tZ28tc2VydmVybGVzcy13ZWJhcHA=")
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/git/{base64url}", GetGitProjectFromRequest)
	router.ServeHTTP(rr, req)

	if rr.Result().StatusCode != 200 {
		t.Errorf("Expected 200")
	}
}

func TestGetGitProjectFromRequestError(t *testing.T) {
	rr := httptest.NewRecorder()
	path := fmt.Sprintf("/git/%s", "aHR0cHM6Ly9naXRodWIuY29tL2x5c2ltb24vaGVsbG8tZ")
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/git/{base64url}", GetGitProjectFromRequest)
	router.ServeHTTP(rr, req)

	if rr.Result().StatusCode != 400 {
		t.Errorf("Expected 400")
	}
}

func TestGetBranchFromRequestSuccess(t *testing.T) {
	rr := httptest.NewRecorder()
	path := fmt.Sprintf("/git/%s/branch/%s", "aHR0cHM6Ly9naXRodWIuY29tL2x5c2ltb24vaGVsbG8tZ28tc2VydmVybGVzcy13ZWJhcHA=", "0ebc05fd263f18d488438ac9b0afe871f8c3b0f1")
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/git/{base64url}/branch/{commitId}", GetBranchFromRequest)
	router.ServeHTTP(rr, req)

	if rr.Result().StatusCode != 200 {
		t.Errorf("Expected 200")
	}
}

func TestGetBranchFromRequestWrongGitRepository(t *testing.T) {
	rr := httptest.NewRecorder()
	path := fmt.Sprintf("/git/%s/branch/%s", "aHR0cHM6Ly9naXRodWIuY29tL2x5c2ltb24vaGVsbG8tZ28tc2VydmVybGVzcy13ZWJh", "0ebc05fd263f18d488438ac9b0afe871f8c3b0f1")
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/git/{base64url}/branch/{commitId}", GetBranchFromRequest)
	router.ServeHTTP(rr, req)

	if rr.Result().StatusCode != 400 {
		t.Errorf("Expected 400")
	}
}

func TestGetBranchFromRequestWrongCommitId(t *testing.T) {
	rr := httptest.NewRecorder()
	path := fmt.Sprintf("/git/%s/branch/%s", "aHR0cHM6Ly9naXRodWIuY29tL2x5c2ltb24vaGVsbG8tZ28tc2VydmVybGVzcy13ZWJhcHA=", "0ebc05fd263f18d488438ac9b0afe871f8c3b0f5")
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/git/{base64url}/branch/{commitId}", GetBranchFromRequest)
	router.ServeHTTP(rr, req)

	if rr.Result().StatusCode != 400 {
		t.Errorf("Expected 400, got %s", rr.Result().Status)
	}
}
