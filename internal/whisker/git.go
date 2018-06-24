package whisker

// Get basic configuration to show to the user
import (
	"net/http"
	"os/exec"
)

func Git() {
	http.HandleFunc("/git", check_git_repo)
}

// json structure returning the current configuration
func check_git_repo(w http.ResponseWriter, r *http.Request) {
	str := "ok"
	bytes := []byte(str)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func Is_git_rep(git_url string) bool {
	//args := []string{"ls-remote", git_url, ">", "/dev/null"}
	_, err := exec.Command(get_git_path(), "ls-remote", git_url, ">", "/dev/null").Output()
	if err != nil {
		return false
	} else {
		return true
	}
}

func Is_valid_git_branch(git_url string, git_branch string) bool {
	//args := []string{"ls-remote", git_url, ">", "/dev/null"}
	_, err := exec.Command(get_git_path(), "ls-remote", "--exit-code", git_url, git_branch, ">", "/dev/null").Output()
	if err != nil {
		return false
	} else {
		return true
	}
}

func List_git_branch(git_url string, prefix string) bool {
	//args := []string{"ls-remote", git_url, ">", "/dev/null"}
	_, err := exec.Command(get_git_path(), "ls-remote", "--exit-code", git_url, prefix, ">", "/dev/null").Output()
	if err != nil {
		return false
	} else {
		return true
	}
}

func get_git_path() string {
	git_path, lookErr := exec.LookPath("git")
	if lookErr != nil {
		panic(lookErr)
	}
	return git_path
}
