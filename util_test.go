package main

import (
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test Add repo entry to pacman.conf.template.
func TestAddRepoEntry(t *testing.T) {
	src := "pacman.conf.template"

	err := exec.Command("cp", "etc/pacman.conf.template", src).Run()
	assert.NoError(t, err, "should not fail")

	r := &Repo{
		name:    "test",
		url:     "http://test.repo.com",
		workdir: ".",
	}

	err = addRepoEntry(src, r)
	assert.NoError(t, err, "should not fail")

	err = addRepoEntry(src, r)
	assert.NoError(t, err, "should not fail")

	// src does not exist
	err = addRepoEntry("fake", r)
	assert.Error(t, err, "should fail")

	// invalid src
	invalidSrc := "pacman.conf.template.fake"
	f, err := os.Create(invalidSrc)
	assert.NoError(t, err, "should not fail")
	err = f.Close()
	assert.NoError(t, err, "should not fail")
	err = addRepoEntry(invalidSrc, r)
	assert.Error(t, err, "should fail")

	// cleanup
	err = os.Remove(src)
	assert.NoError(t, err, "should not fail")
	err = os.Remove(invalidSrc)
	assert.NoError(t, err, "should not fail")
}

// Test adding a custom pacman.conf file for the build instance.
func TestAddPacmanConf(t *testing.T) {
	if os.Getenv("DOCKER_TEST") != "1" {
		return
	}

	file := "test_pacman.conf"
	f, err := os.Create(file)
	assert.NoError(t, err, "should not fail")
	assert.NotNil(t, f, "should not be nil")

	_, err = f.WriteString("test")
	assert.NoError(t, err, "should not fail")
	err = f.Close()
	assert.NoError(t, err, "should not fail")

	err = addPacmanConf(file)
	assert.NoError(t, err, "should not fail")

	err = exec.Command("sudo", "rm", "/etc/pacman.conf").Run()
	assert.NoError(t, err, "should not fail")
}

// Test adding custom pacman mirror.
func TestAddMirror(t *testing.T) {
	if os.Getenv("DOCKER_TEST") != "1" {
		return
	}

	tmpFile := "test_pacman_mirrorlist"
	mirror := "http://mirror.one.com/archlinux/$repo/os/$arch/"

	err := addMirror(mirror, tmpFile)
	assert.NoError(t, err, "should not fail")
}

// Test setting up build dirs.
func TestSettingUpBuildDirs(t *testing.T) {
	base := "build_home_test"

	sources, repo, err := setupBuildDirs(base)
	assert.NoError(t, err, "should not fail")
	assert.Equal(t, sources, path.Join(base, "sources"), "should be equal")
	assert.Equal(t, repo, path.Join(base, "repo"), "should be equal")

	// cleanup
	err = os.RemoveAll(base)
	assert.NoError(t, err, "should not fail")
}

// Test splitting a repo definition.
func TestSplitRepoDef(t *testing.T) {
	repoDef := "name=http://path.to.repo"

	name, url, err := splitRepoDef(repoDef)
	assert.NoError(t, err, "should not fail")
	assert.Equal(t, name, "name", "should be equal")
	assert.Equal(t, url, "http://path.to.repo", "should be equal")

	repoDef = "name==https://path.to.repo"
	_, _, err = splitRepoDef(repoDef)
	assert.Error(t, err, "should fail")

	repoDef = "=https://path.to.repo"
	_, _, err = splitRepoDef(repoDef)
	assert.Error(t, err, "should fail")

	repoDef = "name="
	_, _, err = splitRepoDef(repoDef)
	assert.Error(t, err, "should fail")
}
