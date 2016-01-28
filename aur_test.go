package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var a = AUR{"mockfiles"}

// Test getting deps from AUR.
func TestGetDeps(t *testing.T) {
	deps := make(map[string]struct{})

	err := a.getDeps([]string{"sway-git"}, deps)
	assert.NoError(t, err, "should not fail")
	assert.Equal(t, len(deps), 2, "should have 2 elements")
	_, ok := deps["sway-git"]
	assert.True(t, ok, "should be true")
	_, ok = deps["wlc-git"]
	assert.True(t, ok, "should be true")
}

// Test Getting sources from AUR.
func TestAURGet(t *testing.T) {
	pkgs, err := a.Get([]string{"sway-git"})
	assert.NoError(t, err, "should not fail")

	for _, pkg := range pkgs {
		err = os.RemoveAll(pkg.Path)
		assert.NoError(t, err, "should not fail")
	}
}
