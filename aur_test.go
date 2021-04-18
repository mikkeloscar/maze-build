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
	require.Len(t, deps, 4)
	require.Contains(t, deps, "sway-git")
	require.Contains(t, deps, "swaybg-git")
	require.Contains(t, deps, "wlroots-git")
	require.Contains(t, deps, "seatd")

	deps = make(map[string]struct{})
	err = a.getDeps([]string{"google-cloud-sdk"}, deps)
	require.NoError(t, err)
	require.Len(t, deps, 1)

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
