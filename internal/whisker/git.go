package whisker

// Get basic configuration to show to the user
import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var GitFolder = get_git_folder()

// create git folder used for applications, and remove everything inside
func get_git_folder() string {
	folder := os.Getenv("CONFIGURATION_WHISKER_GIT_FOLDER")
	if folder == "" {
		// Using default
		folder = "/tmp/felicette"
	}
	if !strings.HasPrefix(folder, "/tmp") {
		log.Fatalf("Git folder define by env variable CONFIGURATION_WHISKER_GIT_FOLDER should start with /tmp, got: %v", folder)
	}
	_, err := exec.Command("rm", "-Rf", folder).Output()
	if err != nil {
		log.Fatalf("Unable to delete folder %v: %v", folder, err)
	}
	_, err = exec.Command("mkdir", folder).Output()
	if err != nil {
		log.Fatalf("Unable to create folder %v: %v", folder, err)
	}
	return folder
}

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
	_, err := exec.Command("git", "ls-remote", git_url, ">", "/dev/null").Output()
	if err != nil {
		return false
	} else {
		return true
	}
}

func Is_valid_git_branch(git_url string, git_branch string) bool {
	//args := []string{"ls-remote", git_url, ">", "/dev/null"}
	_, err := exec.Command("git", "ls-remote", "--exit-code", git_url, git_branch, ">", "/dev/null").Output()
	if err != nil {
		return false
	} else {
		return true
	}
}

// Clone a remote git repo locally.
// Needed as we want to access one remote file
func Clone_git_repo(git_url string) bool {
	// Only clone it if it does not exist locally
	if _, err := os.Stat(getLocalGitPath(git_url)); os.IsNotExist(err) {

		cmd := exec.Command("git", "clone", git_url, getMD5Hash(git_url))
		log.Printf("Found md5 of %v : %v ", git_url, getMD5Hash(git_url))

		// Use temporary
		cmd.Dir = GitFolder
		_, err := cmd.Output()
		if err != nil {
			log.Printf("Unable to clone git repository  #%v ", git_url)
			return false
		}
	}
	// Get the latest version

	return true
}

func Fetch_latest(git_url string) bool {
	// Be sure that the latest version is there
	if !Clone_git_repo(git_url) {
		return false
	}
	cmd := exec.Command("git", "fetch", "--all", "--prune")

	// Use temporary
	cmd.Dir = getLocalGitPath(git_url)
	_, err := cmd.Output()
	if err != nil {
		log.Printf("Unable to fetch latest  #%v ", git_url)
		return false
	}
	return true

}

// Retrieve a file from a git_url + branch + filename, nil if not possible
func Get_file(git_url string, branch string, filename string) []byte {
	Fetch_latest(git_url)
	// Get file for this branch
	cmd := exec.Command("git", "show", branch+":"+filename)

	// Use temporary
	cmd.Dir = getLocalGitPath(git_url)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Unable to fetch latest  #%v ", git_url)
		return nil
	}
	return output
}

// Retrieve md5 hashes to have better structure
func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getLocalGitPath(git_url string) string {
	return GitFolder + "/" + getMD5Hash(git_url)
}
