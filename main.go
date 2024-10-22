package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	// "time"
	"log"
	"gopkg.in/yaml.v2"
)

var allowCommits bool

const (
	SHELL = "bash"
	Red = "\033[31m"
	Reset = "\033[0m"
	Green = "\033[32m"
	Yellow = "\033[33m"
	Blue = "\033[34m"
	Purple = "\033[35m"
	Cyan = "\033[36m"
	Gray = "\033[37m"
	White = "\033[97m"
)

type Config struct {
	GitUrls   []string `yaml:"gitUrls"`
	CreateBackup bool `yaml:"createBackup"`
	FileToUpdate string `yaml:"fileToUpdate"`
	Branch struct {
		PullBranch string `yaml:"pullBranch"`
		PushBranch string `yaml:"pushBranch"`
	} `yaml:"branch"`
	StrChanges []struct {
		Match   string `yaml:"match"`
		Replace string `yaml:"replace"`
	} `yaml:"strChanges"`
}

func main() {
	if runtime.GOOS == "linux" {
		fmt.Println(W("[INFO] Application in alpha (maintainer prashant.nandipati@zigram.tech)",Green))
		fmt.Println("[INFO] Running on linux")
	} else {
		fmt.Println("not linux")
	}
	flag.BoolVar(&allowCommits, "commit", false, "Set this flag to auto commit back with messages")
	flag.Parse()

	file, err := os.Open("bbconfig.yaml")
	if err != nil {
		log.Fatalf(W("error: %v",Red), err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf(W("error: %v",Red), err)
	}

	matchStrings := make([]string, len(config.StrChanges))
	replaceStrings := make([]string, len(config.StrChanges))

	for i, change := range config.StrChanges {
		matchStrings[i] = change.Match
		replaceStrings[i] = change.Replace
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, line := range config.GitUrls {
		if len(line) > 0 {
			wg.Add(1)
			go func(repo string) {
				defer wg.Done()
				var rlink []string = strings.Split(line, "/")
				repoName := rlink[len(rlink)-1]
				fmt.Printf("[%s] Starting in Goroutine!\n", repoName)
				b64 := string(base64.StdEncoding.EncodeToString([]byte(line)))
				fpath := fmt.Sprintf("/tmp/git/bbchanges/%s", b64)
				if runtime.GOOS == "windows" {
					fpath = fmt.Sprintf("%s\\bbchanges\\%s", os.TempDir(), b64)
				}

				cloneGitRepo(line, config.Branch.PullBranch, fpath)
				mutex.Lock()
				SedThingy(fpath+"/"+config.FileToUpdate , matchStrings, replaceStrings)
				Shellout(fmt.Sprintf("mv %s %s", fpath+"/Jenkinsfile", fpath+"/Jenkinsfile_old"))
				Shellout(fmt.Sprintf("mv %s %s", fpath+"/Jenkinsfile_modified", fpath+"/Jenkinsfile"))
				Shellout(fmt.Sprintf("cd %s && git checkout -B %s", fpath, config.Branch.PushBranch))
				fmt.Printf(W("[%d] New branch created!\n",Cyan), os.Getpid())
				if allowCommits {
					// use go-git instead
					Commit(fpath+"/*", fpath, fmt.Sprintf("[%d] Updating jenkinsfile for timeout jenkins changes", os.Getpid()))
				} else {
					fmt.Printf(W("[%d] Commiting has been skipped\n",Yellow), os.Getpid())
				}
				mutex.Unlock()
				fmt.Printf(W("[%s] Done!\n",Green), repoName)
			}(line)
		}
	}
	fmt.Printf("[wg] Spawned %d : Awaiting task completions!\n", runtime.NumGoroutine())
	wg.Wait()
}
