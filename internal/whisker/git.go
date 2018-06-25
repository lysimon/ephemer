package whisker

// Get basic configuration to show to the user
import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/google/jsonapi"
)

// Git project structure returned by the api
type GitProject struct {
	Base64Url string   `json:"base64url"`
	Url       string   `json:"url"`
	Branches  []Branch `json:"branches"`
}

type Branch struct {
	Name     string `json:"name"`
	CommitId string `json:"commit_id"`
}

// json structure returning the current configuration
func GetGitProjectFromRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	baseUrl := params["id"]
	project, err := GetGitProject(baseUrl)
	w.Header().Set("Content-Type", "application/json")
	if err != (jsonapi.ErrorObject{}) {
		// return json error with status code
		js, _ := json.Marshal(err)
		http.Error(w, string(js), http.StatusInternalServerError)
	} else {
		js, _ := json.Marshal(project)
		w.Write(js)
	}
	return
}

func GetGitProject(base64url string) (GitProject, jsonapi.ErrorObject) {
	// Decoding base64 url
	urlbytes, err := base64.StdEncoding.DecodeString(base64url)
	if err != nil {
		return GitProject{}, jsonapi.ErrorObject{
			Title:  "Git repository not found",
			Detail: "Unable to find the git repository provided.",
			Status: "400",
		}
	}
	url := string(urlbytes)

	// Listing remote branches
	outputBytes, err := exec.Command("git", "ls-remote", url).Output()
	if err != nil {
		return GitProject{}, jsonapi.ErrorObject{
			Title:  "Validation Error",
			Detail: "Given request body was invalid.",
			Status: "400",
		}
	}
	// Triming space to parse the output correctly
	output := string(bytes.TrimSpace(outputBytes))

	// Reading line by line to create branches
	branches := []Branch{}
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, "\t")

		commitId := split[0]
		ref := split[1]
		fmt.Println(ref)

		// If ref starts with refs/heads/ it is a branch
		if strings.HasPrefix(ref, "refs/heads/") || strings.HasPrefix(ref, "refs/tags/") {
			ref = strings.Replace(ref, "refs/heads/", "origin/", 1)
			ref = strings.Replace(ref, "refs/tags/", "tags/", 1)
			branch := Branch{
				CommitId: commitId,
				Name:     ref,
			}
			fmt.Println(ref)
			branches = append(branches, branch)
		}
	}

	return GitProject{
		Url:       url,
		Base64Url: base64url,
		Branches:  branches,
	}, jsonapi.ErrorObject{}
}

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

func Check_git_branch(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got check_git_repo")

	git_url := r.FormValue("git_url")
	git_branch := r.FormValue("git_branch")

	log.Printf("Got git_url %v", git_url)
	str := strconv.FormatBool(Is_valid_git_branch(git_url, git_branch))
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
	return bytes.TrimSpace(output)
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
