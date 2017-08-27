package main

import (
	"bytes"
	"fmt"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	defaultPackager = "maze-build"
)

var (
	config struct {
		Origin   string
		Package  string
		Packager string
		Repo     string
		SignKey  string // TODO:
	}
)

func main() {
	kingpin.Flag("origin", "Origin of the package e.g. aur or local.").Required().StringVar(&config.Origin)
	kingpin.Flag("package", "Name of the package to build.").StringVar(&config.Package)
	kingpin.Flag("packager", "Name used for the packager.").Default(defaultPackager).StringVar(&config.Packager)
	kingpin.Flag("repo", "URL of upstream repo.").Envar("PLUGIN_REPO").StringVar(&config.Repo)
	kingpin.Parse()

	// configure log
	log.SetFormatter(new(formatter))

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	ws, err := initWorkspace(wd)
	if err != nil {
		log.Fatal(err)
	}

	pkgRepo, err := parseRepo(config.Repo, ws.RepoPath)
	if err != nil {
		log.Fatal(err)
	}

	err = pkgRepo.local.InitDir()
	if err != nil {
		log.Fatal(err)
	}

	builder := &Builder{
		workdir:  ws.SourcesPath,
		repo:     pkgRepo,
		Packager: config.Packager,
	}

	pkgs, err := builder.BuildNew([]string{config.Package}, &AUR{ws.SourcesPath})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pkgs)

	// err = storeBuiltPkgs(path.Join(workspace, "drone_pkgbuild", "packages.built"), pkgs)
	// if err != nil {
	// 	return err
	// }
}

type formatter struct{}

func (f *formatter) Format(entry *log.Entry) ([]byte, error) {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "[%s] %s\n", entry.Level.String(), entry.Message)
	return buf.Bytes(), nil
}

type Workspace struct {
	SourcesPath string
	RepoPath    string
}

func initWorkspace(workdir string) (*Workspace, error) {
	ws := &Workspace{
		SourcesPath: path.Join(workdir, "sources"),
		RepoPath:    path.Join(workdir, "repo"),
	}

	err := os.MkdirAll(ws.SourcesPath, 0755)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(ws.RepoPath, 0755)
	if err != nil {
		return nil, err
	}

	return ws, nil
}

// func run() error {
// 	vargs := ArchBuild{}
// 	repo := os.Getenv("PLUGIN_REPO")
// 	workspace := os.Getenv("DRONE_WORKSPACE")

// 	// srcsPath := path.Join(workspace, "drone_pkgbuild", "sources")
// 	// repoPath := path.Join(workspace, "drone_pkgbuild", "repo")

// 	// repoName, repoUrl, err := splitRepoDef(vargs.Repo)

// 	// pkgRepo := &Repo{
// 	// 	name:    repoName,
// 	// 	url:     repoUrl,
// 	// 	workdir: repoPath,
// 	// }

// 	// // configure build
// 	// if vargs.Packager == "" {
// 	// 	vargs.Packager = "maze-build"
// 	// }

// 	pkgRepo, err := parseRepo(repo, repoPath)
// 	if err != nil {
// 		return err
// 	}

// 	err = pkgRepo.local.InitDir()
// 	if err != nil {
// 		return err
// 	}

// 	builder := &Builder{
// 		workdir: srcsPath,
// 		repo:    pkgRepo,
// 		config:  vargs,
// 	}

// 	// aur := &AUR{srcsPath}

// 	// build, err := parseBuildURLInfo(system.Link, srcsPath)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	buildInst, err := parseGitLog(workspace, srcsPath)
// 	if err != nil {
// 		return err
// 	}

// 	// buildInst := &BuildInst{
// 	// 	Pkgs: vargs.Packages,
// 	// 	Src:  &AUR{srcsPath},
// 	// }

// 	pkgs, err := builder.BuildNew(buildInst.Pkgs, buildInst.Src)
// 	if err != nil {
// 		return err
// 	}

// 	err = storeBuiltPkgs(path.Join(workspace, "drone_pkgbuild", "packages.built"), pkgs)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
