package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/kr/pty"
)

// add repo entry to the pacman.conf file specified by path.
func addRepoEntry(path string, r *Repo) error {
	entry := fmt.Sprintf("[%s]\nServer = %s\n", r.local.Name, r.url)

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

// BuildInst defines a build instruction with package name and sourcer.
type BuildInst struct {
	Pkgs []string
	Src  *AUR
}

// func parseBuildURLInfo(uri, srcPath string) (*build, error) {
// 	u, err := url.Parse(uri)
// 	if err != nil {
// 		return nil, err
// 	}

// 	m, err := url.ParseQuery(u.RawQuery)
// 	if err != nil {
// 		return nil, err
// 	}

// 	build := &build{}

// 	if v, ok := m["pkgs"]; ok {
// 		build.Pkgs = v
// 	}

// 	if v, ok := m["src"]; ok {
// 		switch v[0] {
// 		case "aur":
// 			fallthrough
// 		default:
// 			build.Src = &AUR{srcPath}
// 		}
// 	}

// 	return build, nil
// }

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

// parse build instructions from the git log.
func parseGitLog(dir, srcPath string) (*BuildInst, error) {
	cmd := exec.Command("git", "log", "-1", `--pretty=%B`)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(out)
		return nil, err
	}

	out = bytes.Trim(out, "\n")

	values := strings.Split(string(out), ":")
	if len(values) != 2 {
		return nil, fmt.Errorf("failed to parse log: %s", out)
	}

	buildInst := &BuildInst{
		Pkgs: []string{values[0]},
	}

	switch values[1] {
	case "aur":
		fallthrough
	default:
		buildInst.Src = &AUR{srcPath}
	}

	return buildInst, nil
}

// Store a list of packages built.
func storeBuiltPkgs(file string, pkgs []*BuiltPkg) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(pkgs)
}
