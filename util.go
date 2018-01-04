package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/kr/pty"
	"github.com/mikkeloscar/maze/model"
	"github.com/mikkeloscar/maze/repo"
)

// add repo entry to the pacman.conf file specified by path.
func addRepoEntry(path string, r *Repo) error {
	arch := "/$arch"
	if r.url[len(r.url)-1] == '/' {
		arch = arch[1:]
	}

	// TODO: better handling of SigLevel
	entry := fmt.Sprintf("[%s]\nSigLevel = PackageOptional\nServer = %s%s\n",
		r.local.Name,
		r.url, arch)

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

func splitRepoDef(repo string) (string, string, error) {
	split := strings.Split(repo, "=")
	if len(split) != 2 || split[0] == "" || split[1] == "" {
		return "", "", fmt.Errorf("invalid repo defination: '%s'", repo)
	}

	return split[0], split[1], nil
}

// run command from basedir and print output to stdout.
func runCmd(baseDir string, env []string, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	if baseDir != "" {
		cmd.Dir = baseDir
	}
	cmd.Env = env

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

// handles progress bar output (only shows the last status line)..
func fmtOutput(output string) string {
	sections := strings.Split(output, "\r")
	return sections[len(sections)-1]
}

// Clone git repository from url to dst
func gitClone(src, dst string) error {
	err := runCmd("", nil, "git", "clone", "-q", src, dst)
	if err != nil {
		return err
	}

	return nil
}

var repoPattReg = regexp.MustCompile(`([a-z\d][a-z\d@._+-]*)=http(s)?://.+`)
var repoPatt = regexp.MustCompile(`http(s)?://.+/([a-z\d][a-z\d@._+-]*)/([a-z\d][a-z\d@._+-]*)`)

func parseRepo(uri, basePath string) (*Repo, error) {
	matches := repoPattReg.FindStringSubmatch(uri)
	if len(matches) > 0 {
		return &Repo{
			local: *repo.NewRepo(&model.Repo{
				Name: matches[1],
			}, basePath),
			url: uri,
		}, nil
	}

	matches = repoPatt.FindStringSubmatch(uri)
	if len(matches) == 0 {
		return nil, fmt.Errorf("invalid repo uri: %s", uri)
	}

	return &Repo{
		local: *repo.NewRepo(&model.Repo{
			Owner: matches[2],
			Name:  matches[3],
		}, basePath),
		url: uri,
	}, nil
}
