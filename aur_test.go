package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var a = AUR{baseDir}

// Test getting deps from AUR.
func TestGetDeps(t *testing.T) {
	deps := make(map[string]struct{})

	err := a.getDeps([]string{"sway-git"}, deps)
	assert.NoError(t, err)
	assert.Len(t, deps, 4)
	_, ok := deps["sway-git"]
	assert.True(t, ok)
	_, ok = deps["wlroots-git"]
	assert.True(t, ok)

	deps = make(map[string]struct{})
	err = a.getDeps([]string{"virtualbox-host-modules-mainline"}, deps)
	assert.NoError(t, err)
	assert.Len(t, deps, 2)

	err = os.RemoveAll(baseDir)
	assert.NoError(t, err)
}

// Test Getting sources from AUR.
func TestAURGet(t *testing.T) {
	_, err := a.Get([]string{"sway-git"})
	assert.NoError(t, err)

	err = os.RemoveAll(baseDir)
	assert.NoError(t, err)
}
