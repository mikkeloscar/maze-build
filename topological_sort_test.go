package main

import (
	"testing"

	"github.com/mikkeloscar/gopkgbuild"
	"github.com/stretchr/testify/assert"
)

func TestTopologicalSort(t *testing.T) {
	pkgs := []*SrcPkg{
		&SrcPkg{
			PKGBUILD: &pkgbuild.PKGBUILD{
				Pkgnames: []string{
					"neovim-git",
				},
				Depends: []*pkgbuild.Dependency{
					&pkgbuild.Dependency{
						Name: "libtermkey-bzr",
					},
					&pkgbuild.Dependency{
						Name: "libvterm-bzr",
					},
					&pkgbuild.Dependency{
						Name: "unibilium",
					},
				},
			},
		},
		&SrcPkg{
			PKGBUILD: &pkgbuild.PKGBUILD{
				Pkgnames: []string{
					"libtermkey-bzr",
				},
				Depends: []*pkgbuild.Dependency{
					&pkgbuild.Dependency{
						Name: "unibilium",
					},
				},
			},
		},
		&SrcPkg{
			PKGBUILD: &pkgbuild.PKGBUILD{
				Pkgnames: []string{
					"unibilium",
				},
			},
		},
		&SrcPkg{
			PKGBUILD: &pkgbuild.PKGBUILD{
				Pkgnames: []string{
					"libvterm-bzr",
				},
			},
		},
	}

	s, err := TopologicalSort(pkgs)
	assert.NoError(t, err, "should not fail")
	assert.Equal(t, 4, len(s), "length of sorted list should be 4")
	assert.Equal(t, "neovim-git", s[3].PKGBUILD.Pkgnames[0], "last package should be 'neovim-git'")
}
