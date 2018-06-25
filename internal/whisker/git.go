package whisker

// Get basic configuration to show to the user
import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
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

// json structure returning the current configuration
func GetGitProjectFromRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	baseUrl := params["base64url"]
	project, err := GetGitProject(baseUrl)
	w.Header().Set("Content-Type", "application/json")
	if err != (jsonapi.ErrorObject{}) {
		// return json error with status code
		js, _ := json.Marshal(err)

		statusCode, _ := strconv.Atoi(err.Status)
		http.Error(w, string(js), statusCode)
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

		// If ref starts with refs/heads/ it is a branch
		if strings.HasPrefix(ref, "refs/heads/") || strings.HasPrefix(ref, "refs/tags/") {
			ref = strings.Replace(ref, "refs/heads/", "origin/", 1)
			ref = strings.Replace(ref, "refs/tags/", "tags/", 1)
			branch := Branch{
				CommitId: commitId,
				Name:     ref,
			}
			branches = append(branches, branch)
		}
	}

	return GitProject{
		Url:       url,
		Base64Url: base64url,
		Branches:  branches,
	}, jsonapi.ErrorObject{}
}

func GetBranch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	baseUrl := params["base64url"]
	project, err := GetGitProject(baseUrl)
	if err != (jsonapi.ErrorObject{}) {
		// return json error with status code
		js, _ := json.Marshal(err)
		http.Error(w, string(js), http.StatusInternalServerError)
		return
	}
	// retrieve file from the commitId
	commitId := params["commitId"]
	GetFile(project, commitId, "felicette.yml")
	w.Header().Set("Content-Type", "application/json")
	if err != (jsonapi.ErrorObject{}) {
		// return json error with status code
		js, _ := json.Marshal(err)
		http.Error(w, string(js), http.StatusInternalServerError)
	} else {
		// retrieve felicette.yml file
		js, _ := json.Marshal(project)
		w.Write(js)
	}
	return
}

func getBranch(base64url string, commitId string) (Branch, jsonapi.ErrorObject) {
	project, err := GetGitProject(base64url)
	if err != (jsonapi.ErrorObject{}) {
		return Branch{}, err
	}

	err = cloneGitRepository(project)
	if err != (jsonapi.ErrorObject{}) {
		return Branch{}, err
	}

	return Branch{}, err
}

func cloneGitRepository(gitProject GitProject) jsonapi.ErrorObject {
	// If directory does not exists
	if _, err := os.Stat(getLocalGitPath(gitProject.Base64Url)); os.IsNotExist(err) {

		//Clone it locally
		cmd := exec.Command("git", "clone", gitProject.Url, gitProject.Base64Url)
		cmd.Dir = GitFolder
		_, err := cmd.Output()
		if err != nil {
			return jsonapi.ErrorObject{
				Title:  "Unable to clone git repo",
				Detail: "Unable to find the git repository provided.",
				Status: "400",
			}
		}
	}
	// Get the latest version

	return jsonapi.ErrorObject{}
}

func fetchLatest(gitProject GitProject) jsonapi.ErrorObject {
	// Be sure that the latest version is there
	err := cloneGitRepository(gitProject)
	if err != (jsonapi.ErrorObject{}) {
		return err
	}
	cmd := exec.Command("git", "fetch", "--all", "--prune")
	// Use temporary
	cmd.Dir = getLocalGitPath(gitProject.Base64Url)
	_, erro := cmd.Output()
	if erro != nil {
		return jsonapi.ErrorObject{
			Title:  "Unable tofetch latest commit",
			Detail: "git fetch --all --prune failed.",
			Status: "400",
		}
	}
	return jsonapi.ErrorObject{}
}

// Retrieve a file from a git_url + branch + filename, nil if not possible
func GetFile(gitProject GitProject, commitId string, filename string) ([]byte, jsonapi.ErrorObject) {
	fetchLatest(gitProject)
	// Get file for this branch
	cmd := exec.Command("git", "show", commitId+":"+filename)

	// Use temporary
	cmd.Dir = getLocalGitPath(gitProject.Base64Url)
	output, err := cmd.Output()
	if err != nil {
		return []byte{}, jsonapi.ErrorObject{
			Title:  "Unable tofetch latest commit",
			Detail: "git fetch --all --prune failed, please, check felicette-whisker configuration",
			Status: "400",
		}
	}
	return bytes.TrimSpace(output), jsonapi.ErrorObject{}
}

func getLocalGitPath(folderName string) string {
	return GitFolder + "/" + folderName
}
