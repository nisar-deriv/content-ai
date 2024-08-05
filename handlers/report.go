package handlers

import (
	"encoding/base64"
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
	log.Println("Starting to load configuration")
	cfg = config.LoadConfig()
	log.Println("Configuration loaded successfully")
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
			content, err := data.ReadFromFile(fmt.Sprintf("%s/%s", weekFolder, file.Name()))
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", file.Name(), err)
			}

			enhancedContent, err := enhanceFullContent(content)
			if err != nil {
				return fmt.Errorf("error enhancing content for file %s: %v", file.Name(), err)
			}

			enhancedFilename := fmt.Sprintf("%s/enhanced_%s", weekFolder, file.Name())
			// Check if the enhanced file already exists and delete it
			if _, err := os.Stat(enhancedFilename); err == nil {
				err = os.Remove(enhancedFilename)
				if err != nil {
					return fmt.Errorf("error deleting existing enhanced file %s: %v", enhancedFilename, err)
				}
			}
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

func enhanceFullContent(content string) (string, error) {
	if cfg.UseOllama {
		return nlp.EnhanceTextWithOllama(content)
	}
	return nlp.EnhanceTextWithOpenAI(content, cfg.OpenAIKey)
}

func cloneGitHubRepo(repoURL string) (string, error) {
	log.Printf("Cloning GitHub repository from %s", repoURL)
	cloneDir := "/tmp/repo"
	cmd := exec.Command("git", "clone", repoURL, cloneDir)
	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to clone repository: %v", err)
		return "", err
	}
	log.Println("Repository cloned successfully")
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
	log.Println("Starting to construct markdown content")

	for _, file := range files {
		if !file.IsDir() {
			log.Printf("Processing file: %s\n", file.Name())

			content, err := data.ReadFromFile(filepath.Join(weekFolder, file.Name()))
			if err != nil {
				log.Printf("Error reading file %s: %v\n", file.Name(), err)
				return "", fmt.Errorf("error reading file %s: %v", file.Name(), err)
			}

			log.Printf("Successfully read file: %s\n", file.Name())

			enhancedContent, err := enhanceFullContent(content)
			if err != nil {
				log.Printf("Error enhancing content for file %s: %v\n", file.Name(), err)
				return "", fmt.Errorf("error enhancing content for file %s: %v", file.Name(), err)
			}

			log.Printf("Successfully enhanced content for file: %s\n", file.Name())

			// Parse the enhanced content and append to the respective sections
			parseContent(enhancedContent, &progress, &problems, &plan, &insights)
			log.Printf("Parsed content for file: %s\n", file.Name())
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

func parseContent(content string, progress, problems, plan, insights *strings.Builder) {
	// Implement logic to parse the content and append to the respective sections
}

func commitAndPushChanges(cloneDir, mdFilename string) error {
	log.Printf("Committing changes in directory %s", cloneDir)
	cmd := exec.Command("git", "add", mdFilename)
	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to add file %s: %v", mdFilename, err)
		return err
	}
	log.Println("File added successfully")

	cmd = exec.Command("git", "commit", "-m", "Add weekly report")
	err = cmd.Run()
	if err != nil {
		log.Printf("Failed to commit changes: %v", err)
		return err
	}
	log.Println("Changes committed successfully")

	cmd = exec.Command("git", "push", "origin", "main")
	err = cmd.Run()
	if err != nil {
		log.Printf("Failed to push changes: %v", err)
		return err
	}
	log.Println("Changes pushed successfully")

	return nil
}

func setupSSHKey() error {
	deployKey := os.Getenv("GIT_DEPLOY_KEY")
	if deployKey == "" {
		return fmt.Errorf("GIT_DEPLOY_KEY is not set in the environment")
	}

	// Decode the base64 encoded deploy key
	decodedKey, err := base64.StdEncoding.DecodeString(deployKey)
	if err != nil {
		return fmt.Errorf("error decoding base64 deploy key: %v", err)
	}

	sshDir := os.Getenv("HOME") + "/.ssh"
	err = os.MkdirAll(sshDir, 0700)
	if err != nil {
		return fmt.Errorf("error creating .ssh directory: %v", err)
	}

	keyPath := sshDir + "/id_rsa"
	err = os.WriteFile(keyPath, decodedKey, 0600)
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
