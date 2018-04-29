package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

// User options for changing the behavior of the script
var options = struct {
	ShouldRemoveRepoOnComplete bool
}{
	ShouldRemoveRepoOnComplete: true,
}

const (
	//baseRepo             = "git@git.hexad.de:carnet"
	baseRepo             = "git@github.com:harihar"
	releaseVersionPrefix = "release"
)

// List of all repos
//var repos = []string{"profile-completion-frontend",
//	"profile-completion-backend", "profile-integrity-service",
//	"spin-management-service", "user-management-service",
//	"token-refresh-service", "vehicle-user-update-service",
//}

//var repos = []string{"go-release-test"}

var repos = []string{"MyGitRepo", "SpringLearning"}

// TODO collect statistics fo each repo - done, time taken, not done etc
func main() {
	if len(os.Args) < 2 {
		log.Println("No arguments passed. Please provide the release version number.")
		log.Println("Example - go run go-release.go 1.0.0")
		os.Exit(1)
	}
	startTime := time.Now()
	newReleaseVersion := os.Args[1]
	newReleaseBranch := fmt.Sprintf("%s-%s", releaseVersionPrefix, newReleaseVersion)
	// TODO validate the version has semver syntax

	// Create a temp dir for cloning all repos
	tempDirForRepos := os.TempDir() + "carnet-repos"
	err := os.RemoveAll(tempDirForRepos)
	panicIfErr(fmt.Sprintf("Could not delete temp dir - %s", tempDirForRepos), err)

	err = os.Mkdir(tempDirForRepos, os.ModePerm)
	panicIfErr(fmt.Sprintf("Could not create temp dir - %s", tempDirForRepos), err)

	for _, repo := range repos {
		err = os.Chdir(tempDirForRepos)
		panicIfErr(fmt.Sprintf("Could not change directory to - %s", tempDirForRepos), err)

		repoURL := fmt.Sprintf("%s/%s.git", baseRepo, repo)

		gitClone(repoURL)

		// Go to the repo dir
		repoDir := fmt.Sprintf("%s/%s", tempDirForRepos, repo)
		err = os.Chdir(repoDir)
		panicIfErr(fmt.Sprintf("Could not change dir to - %s", repoDir), err)

		// Bump up the version in the version file to 1.0.0, commit it, push upstream
		updateVersionNoInVersionBranch(repo, newReleaseVersion)

		// New release branch should be created of the form release-1.0.0, push upstream
		createNewReleaseBranch(repo, newReleaseBranch)

		// Update the concourse release-source resource, in master, to the new version, commit it, push it
	}
	log.Printf("**************All Good**************")
	log.Printf("Total time taken %.2f seconds", time.Since(startTime).Seconds())
	log.Printf("Please update the source branch to %s in concourse for all release pipelines", newReleaseBranch)
}

func checkoutMasterBranch(repo string) {
	// Go to master branch
	err := exec.Command("git", "checkout", "master").Run()
	panicIfErr(fmt.Sprintf("Could not check out master branch of repo - %s", repo), err)
}

func createNewReleaseBranch(repo, newReleaseBranch string) {
	checkoutMasterBranch(repo)
	log.Printf("Creating release branch %s for repo - %s", newReleaseBranch, repo)
	err := exec.Command("git", "checkout", "-b", newReleaseBranch).Run()
	panicIfErr(fmt.Sprintf("Could not create release branch %s for repo - %s", newReleaseBranch, repo), err)
	err = exec.Command("git", "push", "-u", "origin", newReleaseBranch).Run()
	panicIfErr(fmt.Sprintf("Could not push branch %s upstream for repo - %s", newReleaseBranch, repo), err)
	log.Printf("Created and pushed new release branch %s for repo - %s", newReleaseBranch, repo)
}

func gitClone(repoURL string) {
	// Git clone repo
	log.Printf("Cloning repo - %s", repoURL)
	err := exec.Command("git", "clone", repoURL).Run()
	panicIfErr(fmt.Sprintf("Could not clone repo - %s", repoURL), err)
	log.Printf("Done cloning repo - %s", repoURL)
}

func updateVersionNoInVersionBranch(repo, newReleaseVersion string) {
	// Git checkout version branch
	log.Printf("Checking out branch 'version' of repo - %s", repo)
	err := exec.Command("git", "checkout", "-b", "version", "remotes/origin/version").Run()
	panicIfErr(fmt.Sprintf("Could not checkout branch 'version' of repo - %s", repo), err)

	// Bump up the version number
	log.Printf("Writing the new release version number to version file of repo - %s", repo)
	file, err := os.Create("version")
	panicIfErr(fmt.Sprintf("Error opening version file of repo - %s", repo), err)
	_, err = file.WriteString(newReleaseVersion + "\n")
	panicIfErr(fmt.Sprintf("Could not write %s to version file", newReleaseVersion), err)
	log.Printf("Done writing the new release version number to version file of repo - %s", repo)

	// Git commit
	err = exec.Command("git", "add", "version").Run()
	panicIfErr(fmt.Sprintf("Error while staging version file changes for repo - %s", repo), err)

	commitMessage := fmt.Sprintf("'RELEASE CUT - Updating release version no to %s'", newReleaseVersion)
	err = exec.Command("git", "commit", "-m", commitMessage).Run()
	panicIfErr(fmt.Sprintf("Error in git commit in version branch of repo - %s", repo), err)
	log.Printf("Commited changes to version file in version branch of repo - %s", repo)

	// Git push
	err = exec.Command("git", "push").Run()
	panicIfErr(fmt.Sprintf("Error while doing git push to version branch of repo - %s", repo), err)
	log.Printf("Pushed changes to version branch of repo - %s", repo)
}

func panicIfErr(message string, err error) {
	if err != nil {
		log.Panicf("%s. Error - %+v", message, err)
	}
}
