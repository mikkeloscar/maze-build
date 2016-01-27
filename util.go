package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

// add repo entry to the pacman.conf file specified by path.
func addRepoEntry(path string, r *Repo) error {
	entry := fmt.Sprintf("[%s]\nServer = %s\n", r.name, r.url)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	lines := strings.Split(string(content), "\n")
	inserted := false
	for _, line := range lines {
		if line == "# :INSERT_REPO:" && !inserted {
			inserted = true
			_, err = w.WriteString(entry + "\n")
			if err != nil {
				return err
			}
		}

		_, err = w.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	if !inserted {
		return fmt.Errorf("No place to insert repo")
	}

	return nil
}

// add pacman.conf from src.
func addPacmanConf(src string) error {
	return exec.Command("sudo", "mv", src, "/etc/pacman.conf").Run()
}

// Add custom mirror to /etc/pacman.d/mirrorlist.
func addMirror(mirror, tmpFile string) error {
	entry := fmt.Sprintf("Server = %s\n", mirror)

	f, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	_, err = w.WriteString(entry)
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return exec.Command("sudo", "mv", tmpFile, "/etc/pacman.d/mirrorlist").Run()
}

// create buildirs and return full path to repo base dir and sources base dir.
func setupBuildDirs(base string) (string, string, error) {
	repo := path.Join(base, "repo")
	err := os.MkdirAll(repo, 0755)
	if err != nil {
		return "", "", err
	}

	sources := path.Join(base, "sources")
	err = os.MkdirAll(repo, 0755)
	if err != nil {
		return "", "", err
	}

	return sources, repo, nil
}

func splitRepoDef(repo string) (string, string, error) {
	split := strings.Split(repo, "=")
	if len(split) != 2 || split[0] == "" || split[1] == "" {
		return "", "", fmt.Errorf("invalid repo defination: '%s'", repo)
	}

	return split[0], split[1], nil
}

// Get a list of aur packages defined in repo.
// TODO: add more ways to define/find packages.
// func getPkgs(base string) ([]string, error) {
// 	aurFile := path.Join(base, "aur_packages")
// 	content, err := ioutil.ReadFile(aurFile)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// TODO: add more options like custom dependencies etc.
// 	lines := strings.Split(string(content), "\n")
// 	if len(lines) == 0 {
// 		return nil, fmt.Errorf("no packages found")
// 	}

// 	return lines, nil
// }
