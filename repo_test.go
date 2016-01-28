package main

import (
	"os"
	"testing"

	"github.com/mikkeloscar/gopkgbuild"
	"github.com/stretchr/testify/assert"
)

var (
	repo1 = Repo{
		name:    "repo1",
		url:     "mockfiles",
		workdir: "mockfiles",
	}
	repo2 = Repo{
		name:    "repo2",
		url:     "mockfiles",
		workdir: "mockfiles",
	}
)

// Test splitting name and version string.
func TestSplitNameVersion(t *testing.T) {
	name, version := splitNameVersion("ca-certificates-20150402-1/")
	assert.Equal(t, "ca-certificates", name, "should be equal")
	assert.Equal(t, "20150402-1", version, "should be equal")

	name, version = splitNameVersion("ca-certificates-2:20150402-1/")
	assert.Equal(t, "ca-certificates", name, "should be equal")
	assert.Equal(t, "2:20150402-1", version, "should be equal")

	name, version = splitNameVersion("zlib-1.2.8-4/")
	assert.Equal(t, "zlib", name, "should be equal")
	assert.Equal(t, "1.2.8-4", version, "should be equal")
}

// Test IsNew.
func TestIsNew(t *testing.T) {
	pkg := &SrcPkg{
		PKGBUILD: &pkgbuild.PKGBUILD{
			Pkgnames: []string{"ca-certificates"},
			Pkgver:   pkgbuild.Version("20150402"),
			Epoch:    0,
			Pkgrel:   1,
		},
	}

	// Check if existing package is new
	new, err := repo1.IsNew(pkg)
	assert.NoError(t, err, "should not fail")
	assert.False(t, new, "should be false")

	// Check if new package is new
	pkg.PKGBUILD.Pkgrel = 2
	new, err = repo1.IsNew(pkg)
	assert.NoError(t, err, "should not fail")
	assert.True(t, new, "should be true")

	// Check if old package is new
	pkg.PKGBUILD.Pkgrel = 1
	pkg.PKGBUILD.Pkgver = pkgbuild.Version("20150401")
	new, err = repo1.IsNew(pkg)
	assert.NoError(t, err, "should not fail")
	assert.False(t, new, "should be false")

	// Check if existing package is new (repo is empty)
	pkg.PKGBUILD.Pkgver = pkgbuild.Version("20150402")
	new, err = repo2.IsNew(pkg)
	assert.NoError(t, err, "should not fail")
	assert.True(t, new, "should be true")
}

// Test GetUpdated.
func TestGetUpdated(t *testing.T) {
	baseDir := "build_home_test"
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
