package main

import (
	"os"
	"testing"

	"github.com/mikkeloscar/maze-repo/repo"
	"github.com/stretchr/testify/assert"
)

var (
	repo1 = Repo{
		local: repo.Repo{
			Name: "repo1",
			Path: baseDir,
		},
		url: baseDir,
	}
	repo2 = Repo{
		local: repo.Repo{
			Name: "repo2",
			Path: baseDir,
		},
		url: baseDir,
	}
)

// Test GetUpdated.
func TestGetUpdated(t *testing.T) {
	aurSrc = &AUR{baseDir + "/sources"}

	_, _, err := setupBuildDirs(baseDir)
	assert.NoError(t, err, "should not fail")

	pkgs, err := aurSrc.Get([]string{"wlc-git"})
	assert.NoError(t, err, "should not fail")

	_, err = repo2.GetUpdated(pkgs)
	assert.NoError(t, err, "should not fail")

	// cleanup
	err = os.RemoveAll(baseDir)
	assert.NoError(t, err, "should not fail")
}
