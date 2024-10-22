// https://stackoverflow.com/questions/61677396/pushing-to-a-remote-with-basic-authentication 
package main

import (
	"fmt"
	_ "io/ioutil"
	"os/exec"
	_ "path/filepath"
	"time"

	"github.com/apex/log"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var defaultRemoteName string = "origin"

// Commit creates a commit in the current repository
func Commit(filename string, gitDataDirectory string, commitMessage string) {
	repo, err := git.PlainOpen(gitDataDirectory)
	if err != nil {
		log.Warn("The data folder could not be converted into a Git repository. Therefore, the versioning does not work as expected.")
		return
	}

	w, _ := repo.Worktree()

	log.Info("Committing new changes...")
	w.AddGlob("*")

	_, _ = w.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  BB_NAME,
			Email: BB_EMAIL,
			When:  time.Now(),
		},
	})

	repo.Remote(defaultRemoteName)

	auth := &http.BasicAuth{
		Username: BB_USER,
		Password: BB_PASSWORD,
	}

	log.Info("Pushing changes to remote")
	err = repo.Push(&git.PushOptions{
		RemoteName: defaultRemoteName,
		Auth:       auth,
	})

	if err != nil {
		log.WithError(err).Warn("Error during push")
	}
}

// GitOperations performs a series of git operations: creating a branch,
// adding a file, committing changes, and pushing to the remote repository.
func GitOperations(branchName, filePath, commitMessage string) error {
	// Helper function to run commands
	runCommand := func(name string, args ...string) error {
		cmd := exec.Command(name, args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error executing command: %s, output: %s", err, string(output))
		}
		return nil
	}

	// Step 1: Create a new branch
	if err := runCommand("git", "checkout", "-b", branchName); err != nil {
		return fmt.Errorf("creating branch: %w", err)
	}

	// Step 2: Stage the specified file
	if err := runCommand("git", "add ", filePath); err != nil {
		return fmt.Errorf("staging file: %w", err)
	}

	// Step 3: Commit the changes
	if err := runCommand("git", "commit", "-m", commitMessage); err != nil {
		return fmt.Errorf("committing changes: %w", err)
	}

	// Step 4: Push changes to the remote repository
	if err := runCommand("git", "push", "origin", branchName); err != nil {
		return fmt.Errorf("pushing changes: %w", err)
	}

	return nil
}
