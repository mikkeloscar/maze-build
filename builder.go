package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"

	"github.com/mikkeloscar/gopkgbuild"
)

type Builder struct {
	workdir string
	repo    *Repo
	// PGPKeys []
}

// BuildNew checks what packages to build based on related repo and builds
// those that have been updated.
func (b *Builder) BuildNew(pkgs []string, aur *AUR) ([]string, error) {
	// make sure environment is up to date
	err := b.update()
	if err != nil {
		fmt.Println("ERROR ONE!")
		return nil, err
	}

	// get packages that should be built
	srcPkgs, err := b.getBuildPkgs(pkgs, aur)
	if err != nil {
		fmt.Println("ERROR TWO!")
		return nil, err
	}

	return b.buildPkgs(srcPkgs)
}

// Update build environment.
func (b *Builder) update() error {
	return runCmd(b.workdir, "sudo", "pacman", "-Syu", "--noconfirm")
}

// Get a sorted list of packages to build.
func (b *Builder) getBuildPkgs(pkgs []string, aur *AUR) ([]*SrcPkg, error) {
	pkgSrcs, err := aur.Get(pkgs)
	if err != nil {
		fmt.Println("ERROR THREE!")
		return nil, err
	}

	// Get a list of devel packages (-{bzr,git,svn,hg}) where an extra
	// version check is needed.
	updates := make([]*SrcPkg, 0, len(pkgSrcs))

	for _, pkgSrc := range pkgSrcs {
		if pkgSrc.PKGBUILD.IsDevel() {
			updates = append(updates, pkgSrc)
		}
	}

	err = b.updatePkgSrcs(updates)
	if err != nil {
		fmt.Println("ERROR FOUR!")
		return nil, err
	}

	return b.repo.GetUpdated(pkgSrcs)
}

// update package sources.
func (b *Builder) updatePkgSrcs(pkgs []*SrcPkg) error {
	for _, pkg := range pkgs {
		_, err := b.updatePkgSrc(pkg)
		if err != nil {
			fmt.Println("ERROR FIVE!")
			return err
		}
	}

	return nil
}

// Check and update if a newer source exist for the package.
func (b *Builder) updatePkgSrc(pkg *SrcPkg) (*SrcPkg, error) {
	err := runCmd(pkg.Path, "makepkg", "-od", "--noconfirm", "--skippgpcheck")
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("mksrcinfo")
	cmd.Dir = pkg.Path
	if err != nil {
		return nil, err
	}

	filePath := path.Join(pkg.Path, ".SRCINFO")

	pkgb, err := pkgbuild.ParseSRCINFO(filePath)
	if err != nil {
		return nil, err
	}

	pkg.PKGBUILD = pkgb

	return pkg, nil
}

// Build a list of packages.
func (b *Builder) buildPkgs(pkgs []*SrcPkg) ([]string, error) {
	buildPkgs := make([]string, 0, len(pkgs))

	for _, pkg := range pkgs {
		pkgPaths, err := b.buildPkg(pkg)
		if err != nil {
			return nil, err
		}

		buildPkgs = append(buildPkgs, pkgPaths...)
	}

	return buildPkgs, nil
}

// Build package and return a list of resulting package archives.
func (b *Builder) buildPkg(pkg *SrcPkg) ([]string, error) {
	err := runCmd(pkg.Path, "makepkg", "-is", "--noconfirm", "--skippgpcheck")
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(pkg.Path)
	if err != nil {
		return nil, err
	}

	pkgs := make([]string, 0, 1)

	for _, f := range files {
		if strings.HasSuffix(f.Name(), "pkg.tar.xz") {
			pkgPath := path.Join(pkg.Path, f.Name())
			pkgs = append(pkgs, pkgPath)
		}
	}

	return pkgs, nil
}
