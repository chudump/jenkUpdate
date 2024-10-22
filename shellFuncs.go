package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// attaches _modified for new files modified with this
func SedThingy(filePath string, matchStr []string, replStr []string) {
	fileData, _ := os.ReadFile(filePath)
	fileString := string(fileData)
	for i := 0; i < len(matchStr); i++ {
		fileString = strings.ReplaceAll(fileString, matchStr[i], replStr[i])
	}
	fileData = []byte(fileString)
	os.WriteFile(filePath+"_modified", []byte(fileString), 0755)
}

func Shellout(command string) string {
	cmd := exec.Command(SHELL, "-c", command)
	cmd.Run()
	stdout, err := cmd.Output()
	fmt.Printf(W("[%d]: %s\n\tstdout %s\n\terr %s\n\n",Gray), os.Getpid(), command, stdout, err)
	if err == nil {
		return string(stdout)
	} else {
		return fmt.Sprintf(W("%s",Red), err)
	}
}

func W(str string,col string) string{
	return fmt.Sprintf("%s%s%s",col,str,Reset)
}

func cloneGitRepo(url string, branch string, fpath string) {
	fmt.Printf(W("[%d] Cloning repo %s\n",Purple), os.Getpid(), url)
	os.Mkdir(fpath, 0755)
	Shellout(fmt.Sprintf("git clone %s --branch %s %s", url, branch, fpath))
}

func perFormSedOperations(fpath string, fname string, matchStr string, replStr string) {
	out := Shellout(fmt.Sprintf("cd %s && sed -i 's#%s#%s#g' %s", fpath, matchStr, replStr, fname))
	fmt.Printf("[%d] Sed : %s\n", os.Getpid(), out)
}