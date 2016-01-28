package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/kr/pty"
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
	err = os.MkdirAll(sources, 0755)
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

type Build struct {
	Pkgs []string
	Src  *AUR
}

func parseBuildURLInfo(uri, srcPath string) (*Build, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}

	build := &Build{}

	if v, ok := m["pkgs"]; ok {
		build.Pkgs = v
	}

	if v, ok := m["src"]; ok {
		switch v[0] {
		case "aur":
			fallthrough
		default:
			build.Src = &AUR{srcPath}
		}
	}

	return build, nil
}

// run command from basedir and print output to stdout.
func runCmd(baseDir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	if baseDir != "" {
		cmd.Dir = baseDir
	}

	// print command being run
	fmt.Println("$", strings.Join(cmd.Args, " "))

	tty, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(tty)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for scanner.Scan() {
			out := scanner.Text()
			// fmt.Printf("output:\n *****\n %#v \n *****\n", out)
			fmt.Printf("%s\n", fmtOutput(out))
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return err
	}

	wg.Wait()

	return nil
}

func fmtOutput(output string) string {
	sections := strings.Split(output, "\r")
	return sections[len(sections)-1]
}

// Clone git repository from url to dst
func gitClone(src, dst string) error {
	err := runCmd("", "git", "clone", "-q", src, dst)
	if err != nil {
		return err
	}

	return nil
}
