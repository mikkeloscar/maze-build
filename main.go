package main

import (
	"fmt"
	"os"
	"path"

	"github.com/drone/drone-plugin-go/plugin"
)

type ArchBuild struct {
	Repo    string   `json:"repo"`
	AURPkgs []string `json:"aur_pkgs"`
	SignKey string   `json:"sign_key"`
}

func main() {
	err := run()
	if err != nil {
		// TODO: better print error fmt
		fmt.Println(err)
		os.Exit(1)
	}
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

	pkgRepo := &Repo{
		name:    "repo",
		url:     repoPath,
		workdir: repoPath,
	}

	builder := &Builder{
		workdir: srcsPath,
		repo:    pkgRepo,
		config:  &vargs,
	}

	// aur := &AUR{srcsPath}

	// build, err := parseBuildURLInfo(system.Link, srcsPath)
	// if err != nil {
	// 	return err
	// }

	build := &Build{
		Pkgs: []string{"sway-git"},
		Src:  &AUR{srcsPath},
	}

	pkgs, err := builder.BuildNew(build.Pkgs, build.Src)
	if err != nil {
		return err
	}

	fmt.Println(pkgs)

	return nil
}
