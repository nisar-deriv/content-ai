package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/regentmarkets/ContentAI/config"
	"github.com/regentmarkets/ContentAI/data"
	"github.com/regentmarkets/ContentAI/nlp"
)

var cfg config.Config

func InitHandlers() {
	log.Println("Loading configuration")
	cfg = config.LoadConfig() // Corrected to call without parameters
}

type DetailedPayload struct {
	Progress string `json:"progress"`
	Problems string `json:"problems"`
	Plan     string `json:"plan"`
	Insights string `json:"insights"`
}

func ReportGenerationHandlerAi(w http.ResponseWriter, r *http.Request) {
	err := GenerateWeeklyReportsAi()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating weekly reports: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Weekly reports generated successfully")
}

func GenerateWeeklyReportsAi() error {
	weekFolder := getWeekFolder()
	files, err := os.ReadDir(weekFolder)
	if err != nil {
		return fmt.Errorf("error reading week folder: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			var content map[string]interface{}
			err := data.ReadFromFile(fmt.Sprintf("%s/%s", weekFolder, file.Name()), &content)
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", file.Name(), err)
			}

			enhancedContent, err := enhanceFullContent(content)
			if err != nil {
				return fmt.Errorf("error enhancing content for file %s: %v", file.Name(), err)
			}

			enhancedFilename := fmt.Sprintf("%s/enhanced_%s", weekFolder, file.Name())
			err = data.WriteToFile(enhancedFilename, enhancedContent)
			if err != nil {
				return fmt.Errorf("error writing enhanced content to file %s: %v", enhancedFilename, err)
			}
		}
	}

	// Set up SSH key
	err = setupSSHKey()
	if err != nil {
		return err
	}
	// Clone the GitHub repository and get the clone directory
	cloneDir, err := cloneGitHubRepo(cfg.GitHubRepoURL)
	if err != nil {
		return fmt.Errorf("error cloning repo: %v", err)
	}

	// Read data from the enhanced files
	files, err = os.ReadDir(weekFolder)
	if err != nil {
		return fmt.Errorf("error reading week folder: %v", err)
	}

	// Construct the markdown content
	mdContent, err := constructMarkdownContent(files, weekFolder)
	if err != nil {
		return fmt.Errorf("error constructing markdown content: %v", err)
	}

	// Verify if the directory exists in the cloned repository
	mdDir := filepath.Join(cloneDir, "docs", "updates")
	if _, err := os.Stat(mdDir); os.IsNotExist(err) {
		return fmt.Errorf("directory %s does not exist", mdDir)
	}

	// Write the markdown content to a file in the cloned repository
	mdFilename := filepath.Join(mdDir, fmt.Sprintf("%s.md", getMondayDate()))
	err = data.WriteToFile(mdFilename, mdContent)
	if err != nil {
		return fmt.Errorf("error writing markdown content to file %s: %v", mdFilename, err)
	}

	// Commit and push the changes, then raise a PR
	err = commitAndPushChanges(cloneDir, mdFilename)
	if err != nil {
		return fmt.Errorf("error committing and pushing changes: %v", err)
	}

	return nil
}

func enhanceFullContent(content map[string]interface{}) (map[string]interface{}, error) {
	textContent, ok := content["text"].(string)
	if !ok {
		return nil, fmt.Errorf("content does not have 'text' field")
	}

	var enhancedText string
	var err error
	if cfg.UseOllama {
		enhancedText, err = nlp.EnhanceTextWithOllama(textContent)
	} else {
		enhancedText, err = nlp.EnhanceTextWithOpenAI(textContent, cfg.OpenAIKey)
	}
	if err != nil {
		return nil, err
	}

	content["text"] = enhancedText
	return content, nil
}

func cloneGitHubRepo(repoURL string) (string, error) {
	cloneDir := "/tmp/repo"
	cmd := exec.Command("git", "clone", repoURL, cloneDir)
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return cloneDir, nil
}

func getMondayDate() string {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	monday := now.AddDate(0, 0, offset)
	return monday.Format("2006-01-02")
}

func constructMarkdownContent(files []os.DirEntry, weekFolder string) (string, error) {
	var progress, problems, plan, insights strings.Builder

	for _, file := range files {
		if !file.IsDir() {
			var content map[string]interface{}
			err := data.ReadFromFile(filepath.Join(weekFolder, file.Name()), &content)
			if err != nil {
				return "", fmt.Errorf("error reading file %s: %v", file.Name(), err)
			}

			enhancedContent, err := enhanceFullContent(content)
			if err != nil {
				return "", fmt.Errorf("error enhancing content for file %s: %v", file.Name(), err)
			}

			// Parse the enhanced content and append to the respective sections
			parseContent(enhancedContent, &progress, &problems, &plan, &insights)
		}
	}

	mdContent := fmt.Sprintf(`# Weekly progress updates

---

## Week of Monday %s

Greetings, fellow Derivians! Wrapping up another week of Production Operations activities. Until next time, happy coding!

<div class="grid cards" markdown>

- 🚀 __Progress__

%s

- ⚠️ __Problems__

%s

- 📋 __Plan__

%s

- 💡 __Insights__

%s

</div>

## <img src="/assets/images/blocker.png" alt="blocker" width="24" height="24"> **Challenges**
  🔴 Nothing blocks our way as of now 

---`, getMondayDate(), progress.String(), problems.String(), plan.String(), insights.String())

	return mdContent, nil
}

func parseContent(content map[string]interface{}, progress, problems, plan, insights *strings.Builder) {
	// Implement logic to parse the content and append to the respective sections
}

func commitAndPushChanges(cloneDir, mdFilename string) error {
	// Change to the cloned repository directory
	err := os.Chdir(cloneDir)
	if err != nil {
		return fmt.Errorf("error changing directory to %s: %v", cloneDir, err)
	}

	// Add the markdown file to the repository
	cmd := exec.Command("git", "add", mdFilename)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error adding file %s: %v", mdFilename, err)
	}

	// Commit the changes
	cmd = exec.Command("git", "commit", "-m", "Add weekly report")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error committing changes: %v", err)
	}

	// Push the changes
	cmd = exec.Command("git", "push", "origin", "main")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error pushing changes: %v", err)
	}

	// Raise a pull request (this is a placeholder, you need to use GitHub API or a CLI tool like hub or gh)
	cmd = exec.Command("gh", "pr", "create", "--title", "Weekly Report", "--body", "This PR contains the weekly report.")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error creating pull request: %v", err)
	}

	return nil
}

func setupSSHKey() error {
	deployKey := os.Getenv("GIT_DEPLOY_KEY")
	if deployKey == "" {
		return fmt.Errorf("GIT_DEPLOY_KEY is not set in the environment")
	}

	sshDir := os.Getenv("HOME") + "/.ssh"
	err := os.MkdirAll(sshDir, 0700)
	if err != nil {
		return fmt.Errorf("error creating .ssh directory: %v", err)
	}

	keyPath := sshDir + "/id_rsa"
	err = os.WriteFile(keyPath, []byte(deployKey), 0600)
	if err != nil {
		return fmt.Errorf("error writing deploy key to file: %v", err)
	}

	cmd := exec.Command("ssh-keyscan", "github.com")
	knownHosts, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error scanning github.com for SSH keys: %v", err)
	}

	err = os.WriteFile(sshDir+"/known_hosts", knownHosts, 0600)
	if err != nil {
		return fmt.Errorf("error writing known_hosts file: %v", err)
	}

	return nil
}
