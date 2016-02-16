package main

import (
	"os"
	"testing"

	"github.com/mikkeloscar/maze/model"
	"github.com/mikkeloscar/maze/repo"
	"github.com/stretchr/testify/assert"
)

var (
	repo1 = &Repo{
		local: *repo.NewRepo(&model.Repo{Name: "repo1"}, baseDir+"/repo"),
	}
	repo2 = &Repo{
		local: *repo.NewRepo(&model.Repo{Name: "repo2"}, baseDir+"/repo"),
	}
)

// Test GetUpdated.
func TestGetUpdated(t *testing.T) {
	aurSrc = &AUR{baseDir + "/sources"}

	_, _, err := setupBuildDirs(baseDir, repo1, repo2)
	assert.NoError(t, err, "should not fail")

	pkgs, err := aurSrc.Get([]string{"wlc-git"})
	assert.NoError(t, err, "should not fail")

	_, err = repo2.GetUpdated(pkgs)
	assert.NoError(t, err, "should not fail")

	// cleanup
	err = os.RemoveAll(baseDir)
	assert.NoError(t, err, "should not fail")
}
