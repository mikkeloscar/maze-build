package main

import (
	"io/ioutil"
	"os/exec"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/mikkeloscar/gopkgbuild"
)

type Builder struct {
	workdir string
	repo    *Repo
	config  *ArchBuild
}

// BuildNew checks what packages to build based on related repo and builds
// those that have been updated.
func (b *Builder) BuildNew(pkgs []string, aur *AUR) ([]string, error) {
	// make sure environment is up to date
	err := b.update()
	if err != nil {
		return nil, err
	}

	// get packages that should be built
	srcPkgs, err := b.getBuildPkgs(pkgs, aur)
	if err != nil {
		return nil, err
	}

	return b.buildPkgs(srcPkgs)
}

// Update build environment.
func (b *Builder) update() error {
	log.Printf("Updating packages")
	return runCmd(b.workdir, "sudo", "pacman", "-Syu", "--noconfirm")
}

// Get a sorted list of packages to build.
func (b *Builder) getBuildPkgs(pkgs []string, aur *AUR) ([]*SrcPkg, error) {
	log.Printf("Fetching build sources+dependencies for %s", strings.Join(pkgs, ", "))
	pkgSrcs, err := aur.Get(pkgs)
	if err != nil {
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
		return nil, err
	}

	return b.repo.GetUpdated(pkgSrcs)
}

// update package sources.
func (b *Builder) updatePkgSrcs(pkgs []*SrcPkg) error {
	for _, pkg := range pkgs {
		_, err := b.updatePkgSrc(pkg)
		if err != nil {
			return err
		}
	}

	return nil
}

// Check and update if a newer source exist for the package.
func (b *Builder) updatePkgSrc(pkg *SrcPkg) (*SrcPkg, error) {
	p := pkg.PKGBUILD
	if len(p.Pkgnames) > 1 || p.Pkgnames[0] != p.Pkgbase {
		log.Printf("Checking for new version of %s:(%s)", p.Pkgbase, strings.Join(p.Pkgnames, ", "))
	} else {
		log.Printf("Checking for new version of %s", p.Pkgbase)
	}

	err := runCmd(pkg.Path, "makepkg", "--nobuild", "--nodeps", "--noconfirm")
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
	p := pkg.PKGBUILD
	if len(p.Pkgnames) > 1 || p.Pkgnames[0] != p.Pkgbase {
		log.Printf("Building packages %s:(%s)", p.Pkgbase, strings.Join(p.Pkgnames, ", "))
	} else {
		log.Printf("Building packages %s", p.Pkgbase)
	}

	err := runCmd(pkg.Path, "makepkg", "--install", "--syncdeps", "--noconfirm")
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
