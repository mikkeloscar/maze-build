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
	assert.NoError(t, err, "should not fail")
	assert.Equal(t, 2, len(deps), "should have 2 elements")
	_, ok := deps["sway-git"]
	assert.True(t, ok, "should be true")
	_, ok = deps["wlc-git"]
	assert.True(t, ok, "should be true")

	deps = make(map[string]struct{})
	err = a.getDeps([]string{"virtualbox-host-modules-mainline"}, deps)
	assert.NoError(t, err, "should not fail")
	assert.Equal(t, 2, len(deps), "should have 2 elements")

	err = os.RemoveAll(baseDir)
	assert.NoError(t, err, "should not fail")
}

// Test Getting sources from AUR.
func TestAURGet(t *testing.T) {
	_, err := a.Get([]string{"sway-git"})
	assert.NoError(t, err, "should not fail")

	err = os.RemoveAll(baseDir)
	assert.NoError(t, err, "should not fail")
}
