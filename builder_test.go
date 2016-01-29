package main

import (
	"os"
	"path"
	"testing"

	"github.com/mikkeloscar/maze-repo/repo"
	"github.com/stretchr/testify/assert"
)

const (
	baseDir = "build_home_test"
)

var (
	pkgRepo = &Repo{
		local: repo.Repo{
			Name: "repo",
			Path: baseDir + "/repo",
		},
		url: baseDir + "/repo",
	}
	builder = &Builder{
		repo:    pkgRepo,
		workdir: baseDir + "/sources",
	}
	aurSrc = &AUR{baseDir + "/sources"}
)

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

func TestUpdateBuild(t *testing.T) {
	if os.Getenv("DOCKER_TEST") != "1" {
		return
	}

	_, _, err := setupBuildDirs(baseDir)
	assert.NoError(t, err, "should not fail")

	err = builder.update()
	assert.NoError(t, err, "should not fail")

	// cleanup
	err = os.RemoveAll(baseDir)
	assert.NoError(t, err, "should not fail")
}

func TestupdatePkgSrc(t *testing.T) {
	if os.Getenv("DOCKER_TEST") != "1" {
		return
	}

	pkgs, err := aurSrc.Get([]string{"wlc-git"})
	assert.NoError(t, err, "should not fail")

	_, _, err = setupBuildDirs(baseDir)
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

	_, _, err = setupBuildDirs(baseDir)
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

	_, _, err := setupBuildDirs(baseDir)
	assert.NoError(t, err, "should not fail")

	_, err = builder.BuildNew([]string{"imgur"}, aurSrc)
	assert.NoError(t, err, "should not fail")

	// cleanup
	err = os.RemoveAll(baseDir)
	assert.NoError(t, err, "should not fail")
}
