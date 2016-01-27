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
	var workspace = plugin.Workspace{}
	var vargs = ArchBuild{}

	// plugin.Param("repo", &repo)
	// plugin.Param("build", &build)
	plugin.Param("workspace", &workspace)
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

	buildDir := path.Join(workspace.Path, "drone_pkgbuild")
	repoPath, srcsPath, err := setupBuildDirs(buildDir)
	if err != nil {
		return err
	}

	repoName, repoUrl, err := splitRepoDef(vargs.Repo)

	pkgRepo := &Repo{
		name:    repoName,
		url:     repoUrl,
		workdir: repoPath,
	}

	builder := &Builder{
		repo:    pkgRepo,
		workdir: srcsPath,
	}

	aur := &AUR{srcsPath}

	pkgs, err := builder.BuildNew(vargs.AURPkgs, aur)
	if err != nil {
		return err
	}

	fmt.Println(pkgs)

	return nil
}
