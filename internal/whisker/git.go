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

func Can_access_git_rep(git_url string) bool {
	//args := []string{"ls-remote", git_url, ">", "/dev/null"}
	_, err := exec.Command(get_git_path(), "ls-remote", git_url, ">", "/dev/null").Output()
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
