package main

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/client"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"gopkg.in/src-d/go-git.v4/utils/ioutil"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	. "gopkg.in/src-d/go-git.v4/_examples"
)

func main() {
	CheckArgs("<gid>", "<token>", "<sourceDir>")
	gid, _ := strconv.Atoi(os.Args[1])
	token := os.Args[2]
	sourceDir := os.Args[3]

	// Cleanup of moved/deleted projects
	fmt.Println("#######################################")
	CheckRemoteGitProjectExist(sourceDir)

	// First cleanup of empty repositories in case of previous program interruption
	fmt.Println("#######################################")
	CleanEmptyDir(sourceDir)

	fmt.Println("#######################################")
	log.Println("Cloning of Vistaprint projects starting")
	recursiveClone(gid, token, sourceDir)
	log.Println("Cloning of Vistaprint projects finished")

	fmt.Println("#######################################")
	CleanEmptyDir(sourceDir)
}

func recursiveClone(gid int, token, sourceDir string) {

	gitLab := gitlab.NewClient(nil, token)

	optsub := &gitlab.ListSubgroupsOptions{
		AllAvailable: gitlab.Bool(true),
	}

	subgroups, _, _ := gitLab.Groups.ListSubgroups(gid, optsub)

	if len(subgroups) != 0 {
		for _, subgroup := range subgroups {
			log.Printf("Working on %s subgroup\n", subgroup.FullPath)

			optproj := &gitlab.ListGroupProjectsOptions{
				ListOptions: gitlab.ListOptions{
					PerPage: 100,
				},
			}

			projects, _, _ := gitLab.Groups.ListGroupProjects(subgroup.ID, optproj)

			if len(projects) != 0 {
				for _, project := range projects {
					projectPath := sourceDir + "/" + subgroup.FullPath + "/" + project.Path

					if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
						// Project was previously cloned
						log.Printf("%s project already exists, pulling latest changes from remote respository instead\n", project.Name)
						// Pull changes from remote repository
						pullRemoteRepo(projectPath)
					} else {
						os.MkdirAll(projectPath, os.FileMode(0755))
						// Clone the given project to the given directory
						Info("Cloning %s into %s", project.HTTPURLToRepo, sourceDir+"/"+subgroup.FullPath+"/"+project.Path)
						_, err := git.PlainClone(projectPath, false, &git.CloneOptions{
							URL:      project.HTTPURLToRepo[:8] + "gitlab-ci-token:" + token + "@" + project.HTTPURLToRepo[8:],
							Progress: os.Stdout,
						})

						CheckGitError(err, projectPath)
					}
				}
			}
			fmt.Println("---------------------")
			recursiveClone(subgroup.ID, token, sourceDir)
		}
	}
}

func pullRemoteRepo(sourceDir string) {
	// Instantiate a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(sourceDir)
	CheckGitError(err, sourceDir)

	// Get the working directory for the repository
	w, err := r.Worktree()
	CheckGitError(err, sourceDir)

	// Pull the latest changes from the origin remote and merge into the current branch
	err = w.Pull(&git.PullOptions{RemoteName: "origin", Progress: os.Stdout})
	CheckGitError(err, sourceDir)
}

func CheckRemoteGitProjectExist(sourceDir string) {
	log.Println("Cleaning deleted projects starting")
	err := filepath.Walk(sourceDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			dir := filepath.Base(path)

			if dir == ".git" {
				projectPath := filepath.Dir(path)
				repo, _ := git.PlainOpen(projectPath)
				remotes, _ := repo.Remotes()

				for _, remote := range remotes {
					_, err := lsRemote(remote, nil)
					// Delete the project from Opengrok if it's removed from the original path
					CheckGitError(err, projectPath)

					if err == nil {
						log.Printf("Project %s won't be deleted\n", projectPath)
					}
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	log.Println("Cleaning deleted projects finished")
}

func CleanEmptyDir(sourceDir string) {
	log.Println("Cleaning empty directories starting")
	err := filepath.Walk(sourceDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				emptyDir, _ := IsEmpty(path)
				containsGit := strings.Contains(path, ".git")
				if emptyDir && !containsGit {
					log.Printf("Deleting %s\n", path)
					os.RemoveAll(path)
				}
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	log.Println("Cleaning empty directories finished")
}

func CheckGitError(err error, srcPath string) {
	if err == nil {
		return
	}

	switch err.Error() {
	case "remote repository is empty", "already up-to-date", "repository does not exist", "repository already exists", "worktree contains unstaged changes":
		log.Println(err.Error())
		return
	case "repository not found":
		log.Printf("Deleting %s", srcPath)
		os.RemoveAll(srcPath)
		return
	default:
		log.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
		return
	}
}

func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func lsRemote(remote *git.Remote, auth transport.AuthMethod) (memory.ReferenceStorage, error) {
	url := remote.Config().URLs[0]
	s, err := newUploadPackSession(url, auth)
	if err != nil {
		return nil, err
	}
	defer ioutil.CheckClose(s, &err)

	ar, err := s.AdvertisedReferences()
	if err != nil {
		return nil, err
	}

	return ar.AllReferences()
}

func newUploadPackSession(url string, auth transport.AuthMethod) (transport.UploadPackSession, error) {
	c, ep, err := newClient(url)
	if err != nil {
		return nil, err
	}

	return c.NewUploadPackSession(ep, auth)
}

func newClient(url string) (transport.Transport, *transport.Endpoint, error) {
	ep, err := transport.NewEndpoint(url)
	if err != nil {
		return nil, nil, err
	}

	c, err := client.NewClient(ep)
	if err != nil {
		return nil, nil, err
	}

	return c, ep, err
}