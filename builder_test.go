package main

import (
	"os"
	"path"
	"testing"

	"github.com/mikkeloscar/maze/model"
	"github.com/mikkeloscar/maze/repo"
	"github.com/stretchr/testify/assert"
)

const (
	baseDir = "build_home_test"
)

var (
	pkgRepo = &Repo{
		local: *repo.NewRepo(&model.Repo{Name: "repo"}, baseDir),
	}
	builder = &Builder{
		repo:    pkgRepo,
		workdir: baseDir + "/sources",
	}
	aurSrc = &AUR{baseDir + "/sources"}
)

func setupRepoDirs(repos []*Repo) error {
	for _, repo := range repos {
		err := repo.local.InitDir()
		if err != nil {
			return err
		}
		repo.url = repo.local.PathDeep("x86_64")
	}

	return nil
}

// create buildirs and return full path to repo base dir and sources base dir.
func setupBuildDirs(base string, repos ...*Repo) (string, string, error) {
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

	err = setupRepoDirs(repos)
	if err != nil {
		return "", "", err
	}

	return sources, repo, nil
}

func TestUpdateBuild(t *testing.T) {
	if os.Getenv("DOCKER_TEST") != "1" {
		return
	}

	_, _, err := setupBuildDirs(baseDir, pkgRepo)
	assert.NoError(t, err, "should not fail")

	err = builder.update()
	assert.NoError(t, err, "should not fail")

	// cleanup
	err = os.RemoveAll(baseDir)
	assert.NoError(t, err, "should not fail")
}

func TestUpdatePkgSrc(t *testing.T) {
	if os.Getenv("DOCKER_TEST") != "1" {
		return
	}

	pkgs, err := aurSrc.Get([]string{"wlc-git"})
	assert.NoError(t, err, "should not fail")

	_, _, err = setupBuildDirs(baseDir, pkgRepo)
	assert.NoError(t, err, "should not fail")

	_, err = builder.updatePkgSrc(pkgs[0])
	assert.NoError(t, err, "should not fail")

	// cleanup
	err = os.RemoveAll(baseDir)
	assert.NoError(t, err, "should not fail")
}

func TestBuildPkg(t *testing.T) {
	if os.Getenv("DOCKER_TEST") != "1" {
		return
	}

	pkgs, err := aurSrc.Get([]string{"imgur"})
	assert.NoError(t, err, "should not fail")

	_, _, err = setupBuildDirs(baseDir, pkgRepo)
	assert.NoError(t, err, "should not fail")

	_, err = builder.buildPkg(pkgs[0])
	assert.NoError(t, err, "should not fail")

	// cleanup
	err = os.RemoveAll(baseDir)
	assert.NoError(t, err, "should not fail")
}

func TestBuildPkgs(t *testing.T) {
	if os.Getenv("DOCKER_TEST") != "1" {
		return
	}

	_, _, err := setupBuildDirs(baseDir, pkgRepo)
	assert.NoError(t, err, "should not fail")

	_, err = builder.BuildNew([]string{"imgur"}, aurSrc)
	assert.NoError(t, err, "should not fail")

	// cleanup
	err = os.RemoveAll(baseDir)
	assert.NoError(t, err, "should not fail")
}
