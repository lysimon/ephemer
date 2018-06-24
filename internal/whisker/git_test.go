package whisker

import (
	"testing"
)

func TestGit(t *testing.T) {
	result := Is_git_rep("https://github.com/serverless/serverless")
	if !result {
		t.Errorf("expected true, got false")
	}
	result = Is_git_rep("https://github.com/serverless/serverlessss")
	if result {
		t.Errorf("Was able to access git url")
	}
	result = Is_git_rep("trololo")
	if result {
		t.Errorf("Was able to access git url")
	}
	result = Is_valid_git_branch("https://github.com/serverless/serverless", "master")
	if !result {
		t.Errorf("expected true, got false")
	}
	result = Is_valid_git_branch("https://github.com/serverless/serverless", "v1.21.1")
	if !result {
		t.Errorf("expected true, got false")
	}
	result = Is_valid_git_branch("https://github.com/serverless/serverless", "masterddd")
	if result {
		t.Errorf("expected false, got true")
	}
}
