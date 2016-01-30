package main

import (
	"bytes"
	"fmt"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/drone/drone-plugin-go/plugin"
	"github.com/mikkeloscar/maze-repo/repo"
)

// ArchBuild defines the vargs passed from .drone.yml.
type ArchBuild struct {
	Repo     string `json:"repo"`
	SignKey  string `json:"sign_key"`
	Packager string `json:"packager"`
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
	// var repo = plugin.Repo{}
	// var build = plugin.Build{}
	var system = plugin.System{}
	var workspace = plugin.Workspace{}
	var vargs = ArchBuild{}

	// plugin.Param("repo", &repo)
	// plugin.Param("build", &build)
	plugin.Param("system", &system)
	plugin.Param("workspace", &workspace)
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

	srcsPath := path.Join(workspace.Path, "drone_pkgbuild", "sources")
	repoPath := path.Join(workspace.Path, "drone_pkgbuild", "repo")

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

	pkgRepo := &Repo{
		local: repo.Repo{
			Name: "repo",
			Path: repoPath,
		},
		url: repoPath,
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

	build := &build{
		Pkgs: []string{"linux-ck"},
		Src:  &AUR{srcsPath},
	}

	pkgs, err := builder.BuildNew(build.Pkgs, build.Src)
	if err != nil {
		return err
	}

	fmt.Println(pkgs)

	return nil
}
