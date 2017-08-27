package main

import (
	"bytes"
	"fmt"
	"net/url"
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
		Repo     *url.URL
		Token    string
		Upload   bool
		SignKey  string // TODO:
	}
)

func main() {
	kingpin.Flag("origin", "Origin of the package e.g. aur or local.").Required().StringVar(&config.Origin)
	kingpin.Flag("package", "Name of the package to build.").StringVar(&config.Package)
	kingpin.Flag("packager", "Name used for the packager.").Default(defaultPackager).StringVar(&config.Packager)
	kingpin.Flag("repo", "URL of upstream repo.").Envar("PLUGIN_REPO").URLVar(&config.Repo)
	kingpin.Flag("token", "Token used when authenticating with upstream repo.").Envar("TOKEN").StringVar(&config.Token)
	kingpin.Flag("upload", "Specify whether to upload packages or not.").Default("false").BoolVar(&config.Upload)
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

	pkgRepo, err := parseRepo(config.Repo.String(), ws.RepoPath)
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

	if config.Upload {
		config.Repo.Path = ""

		uploader := Uploader{
			client: NewClientToken(config.Repo.String(), config.Token),
			owner:  pkgRepo.local.Owner,
			name:   pkgRepo.local.Name,
		}

		err = uploader.Do(pkgs)
		if err != nil {
			log.Fatal(err)
		}
	}

	// fmt.Println(pkgs)
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
