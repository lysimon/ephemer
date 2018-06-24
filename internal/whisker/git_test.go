package whisker

import (
	"log"
	"testing"
)

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
