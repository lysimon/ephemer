package whisker

import (
	"testing"
)

func TestGit(t *testing.T) {
	result := Can_access_git_rep("https://github.com/serverless/serverless")
	if !result {
		t.Errorf("expected true, got false")
	}

	result = Can_access_git_rep("https://github.com/serverless/serverlessss")
	if result {
		t.Errorf("Was able to access git url")
	}
}
