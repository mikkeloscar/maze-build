package main

import (
	"bytes"
	"fmt"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
)

// ArchBuild defines the config options passed from .drone.yml.
type ArchBuild struct {
	SignKey  string `json:"sign_key"`
	Packager string `json:"packager"`
	// temp
	// Packages []string `json:"packages"`
}

func main() {
	// configure log
	log.SetFormatter(new(formatter))

	err := run()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

type formatter struct{}

func (f *formatter) Format(entry *log.Entry) ([]byte, error) {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "[%s] %s\n", entry.Level.String(), entry.Message)
	return buf.Bytes(), nil
}

func run() error {
	vargs := ArchBuild{}
	repo := os.Getenv("PLUGIN_REPO")
	workspace := os.Getenv("DRONE_WORKSPACE")

	srcsPath := path.Join(workspace, "drone_pkgbuild", "sources")
	repoPath := path.Join(workspace, "drone_pkgbuild", "repo")

	// repoName, repoUrl, err := splitRepoDef(vargs.Repo)

	// pkgRepo := &Repo{
	// 	name:    repoName,
	// 	url:     repoUrl,
	// 	workdir: repoPath,
	// }

	// configure build
	if vargs.Packager == "" {
		vargs.Packager = "maze-build"
	}

	pkgRepo, err := parseRepo(repo, repoPath)
	if err != nil {
		return err
	}

	err = pkgRepo.local.InitDir()
	if err != nil {
		return err
	}

	builder := &Builder{
		workdir: srcsPath,
		repo:    pkgRepo,
		config:  vargs,
	}

	// aur := &AUR{srcsPath}

	// build, err := parseBuildURLInfo(system.Link, srcsPath)
	// if err != nil {
	// 	return err
	// }

	buildInst, err := parseGitLog(workspace, srcsPath)
	if err != nil {
		return err
	}

	// buildInst := &BuildInst{
	// 	Pkgs: vargs.Packages,
	// 	Src:  &AUR{srcsPath},
	// }

	pkgs, err := builder.BuildNew(buildInst.Pkgs, buildInst.Src)
	if err != nil {
		return err
	}

	err = storeBuiltPkgs(path.Join(workspace, "drone_pkgbuild", "packages.built"), pkgs)
	if err != nil {
		return err
	}

	return nil
}
