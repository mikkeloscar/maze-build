package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var a = AUR{baseDir}

// Test getting deps from AUR.
func TestGetDeps(t *testing.T) {
	deps := make(map[string]struct{})

	err := a.getDeps([]string{"sway-git"}, deps)
	require.NoError(t, err)
	require.Len(t, deps, 3)
	require.Contains(t, deps, "sway-git")
	require.Contains(t, deps, "swaybg-git")
	require.Contains(t, deps, "wlroots-git")

	err = os.RemoveAll(baseDir)
	require.NoError(t, err)
}

// Test Getting sources from AUR.
func TestAURGet(t *testing.T) {
	_, err := a.Get([]string{"sway-git"})
	require.NoError(t, err)

	err = os.RemoveAll(baseDir)
	require.NoError(t, err)
}
