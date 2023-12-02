package main

import (
	"fmt"
	myUtils "github.com/SharokhAtaie/utils"
	"github.com/jpillora/go-tld"
	"github.com/projectdiscovery/gologger"
	fUtils "github.com/projectdiscovery/utils/file"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	httpgit "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	accessToken   = "ghp_xxxxxxxxxxxxxxxx"                 // Your Github access token
	UserName      = "admin"                                // Your Github Username
	email         = "admin@gmail.com"                      // Your email (optional)
	repoUrl       = "https://github.com/Username/Repo.git" // Your Github Repo
	repoDir       = "./Repo"                               // Local Directory to save data
	TelegramToken = "123456:xxxxxxxxxxxxx"                 // Your Telegram bot token
	ChatID        = int64(-123456789)                      // Your ChatID of telegram bot
)

func main() {
	repo, err := GitCloneRepo(repoUrl, repoDir, UserName, accessToken)

	file, err := fUtils.ReadFile("urls.txt")
	if err != nil {
		fmt.Println(err)
	}

	for url := range file {
		SaveJs(repoDir, url)
	}

	//Commit and push the changes
	err = commitAndPush(repo, UserName, email, accessToken, TelegramToken, ChatID)
	if err != nil {
		fmt.Printf("Error committing and pushing changes: %s\n", err)
		return
	}
}

func SaveJs(RepoDir, URL string) {
	u, _ := tld.Parse(URL)
	if u.Host == "" {
		return
	}

	u.Path = strings.ReplaceAll(u.Path, "/", "-")
	if strings.HasPrefix(u.Path, "-") {
		u.Path = u.Path[1:]
	}

	response := Request(URL)

	data, err := JsParser(response)
	if err != nil {
		data = response
	}

	exist := pathExists(RepoDir + "/" + u.Host)

	if !exist {
		err = createDirectory(RepoDir + "/" + u.Host)
		if err != nil {
			fmt.Println(err)
		}

		err = saveStringToFile(RepoDir+"/"+u.Host+"/"+u.Path, data)
		if err != nil {
			fmt.Println(err)
		}
		gologger.Info().Msgf("[%s] Directory Created and file saved (%s)", u.Host, u.Path)
	} else {
		err = saveStringToFile(RepoDir+"/"+u.Host+"/"+u.Path, data)
		if err != nil {
			fmt.Println(err)
		}
		gologger.Info().Msgf("[%s] File Override (%s)", u.Host, u.Path)
	}
}

func JsParser(input string) (string, error) {
	options := js.Options{
		WhileToFor: true,
		Inline:     true,
	}

	ast, err := js.Parse(parse.NewInputString(input), options)
	if err != nil {
		return "", err
	}

	return ast.JSString(), nil
}

func saveStringToFile(filePath, content string) error {
	// Open the file with the os.O_TRUNC flag to truncate it if it exists
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func pathExists(directoryPath string) bool {
	_, err := os.Stat(directoryPath)
	if os.IsNotExist(err) {
		return false // Directory does not exist
	}
	return true // Directory exists
}

func createDirectory(directoryPath string) error {
	err := os.Mkdir(directoryPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func Request(URL string) string {
	req, err := http.Get(URL)
	if err != nil {
		fmt.Println(err)
	}
	defer req.Body.Close()

	data, _ := ioutil.ReadAll(req.Body)
	return string(data)
}

func GitCloneRepo(RepoUrl, RepoDir, userName, AccessToken string) (*git.Repository, error) {
	repo, err := git.PlainOpen(RepoDir)
	if err == git.ErrRepositoryNotExists {
		gologger.Info().Msgf("Cloning repo for the first time %s", RepoDir)
		repo, err = git.PlainClone(RepoDir, false, &git.CloneOptions{
			URL:           RepoUrl,
			ReferenceName: plumbing.NewBranchReferenceName("main"),
			Auth: &httpgit.BasicAuth{
				Username: userName,
				Password: AccessToken,
			},
		})
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return repo, nil
}

func commitAndPush(Repo *git.Repository, userName, Email, AccessToken, telegramToken string, chatID int64) error {

	// Create a worktree to perform the operations
	wt, err := Repo.Worktree()
	if err != nil {
		return err
	}

	// Get the status
	status, err := wt.Status()
	if err != nil {
		return err
	}

	if status.String() == "" {
		return nil
	}

	commitMSG := "Updated "

	for s := range status {
		commitMSG += s + " "
	}

	// Add the changes to the repository
	_, err = wt.Add(".")
	if err != nil {
		return err
	}

	commit, err := wt.Commit(commitMSG, &git.CommitOptions{
		Author: &object.Signature{
			Name:  userName,
			Email: Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	// Push the changes to the remote repository
	err = Repo.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{"refs/heads/main:refs/heads/main"},
		Auth: &httpgit.BasicAuth{
			Username: userName,
			Password: AccessToken,
		},
	})
	if err != nil {
		return err
	}

	// Send to Telegram
	_ = myUtils.SendTelegramData(commitMSG, telegramToken, chatID)

	fmt.Printf("Changes pushed successfully! Commit hash: %s\n", commit.String())

	return nil
}
