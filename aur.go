package main

import (
	"fmt"
	"os/exec"
	"path"

	"github.com/mikkeloscar/aur"
	"github.com/mikkeloscar/gopkgbuild"
)

const aurCloneURL = "https://aur.archlinux.org/%s.git"

type AUR struct {
	workdir string
}

type SrcPkg struct {
	PKGBUILD *pkgbuild.PKGBUILD
	Path     string
}

// Get PKGBUILDs from AUR
func (a *AUR) Get(pkgs []string) ([]*SrcPkg, error) {
	deps := make(map[string]struct{})
	err := a.getDeps(pkgs, deps)
	if err != nil {
		return nil, err
	}

	err = a.getSourceRepos(deps)
	if err != nil {
		return nil, err
	}

	var srcPkg *SrcPkg
	var filePath string

	srcPkgs := make([]*SrcPkg, 0, len(deps))

	// get a list of PKGBUILDs/SrcPkgs
	for d, _ := range deps {
		filePath = path.Join(a.workdir, d, ".SRCINFO")

		pkgb, err := pkgbuild.ParseSRCINFO(filePath)
		if err != nil {
			return nil, err
		}
		srcPkg = &SrcPkg{
			PKGBUILD: pkgb,
			Path:     path.Join(a.workdir, d),
		}
		srcPkgs = append(srcPkgs, srcPkg)
	}

	return srcPkgs, nil
}

// query the AUR for build deps to packages.
func (a AUR) getDeps(pkgs []string, updates map[string]struct{}) error {
	pkgsInfo, err := aur.Multiinfo(pkgs)
	if err != nil {
		return err
	}

	for _, pkg := range pkgsInfo {
		updates[pkg.Name] = struct{}{}

		// TODO: maybe add optdepends
		depends := make([]string, 0, len(pkg.Depends)+len(pkg.MakeDepends))
		depends = append(depends, pkg.Depends...)
		depends = append(depends, pkg.MakeDepends...)
		a.getDeps(depends, updates)
	}

	return nil
}

// get source repos from set of package names
func (a *AUR) getSourceRepos(pkgs map[string]struct{}) error {
	clone := make(chan error)

	for pkg, _ := range pkgs {
		go a.updateRepo(pkg, clone)
	}

	errors := make([]error, 0)

	for range pkgs {
		err := <-clone
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		msg := "errors while fetching sources: from AUR\n"
		for _, err := range errors {
			msg += fmt.Sprintf("%s * %s\n", msg, err.Error())
		}
		return fmt.Errorf(msg)
	}

	return nil
}

// update (clone or pull) AUR package repo
func (a *AUR) updateRepo(pkg string, c chan<- error) {
	url := fmt.Sprintf(aurCloneURL, pkg)

	// TODO implement version that can pull instead of clone
	err := gitClone(url, path.Join(a.workdir, pkg))
	c <- err
}

// Clone git repository at url to dst
// TODO: add output
func gitClone(url, dst string) error {
	cmd := exec.Command("git", "clone", url, dst)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// TODO: correct error message?
		return fmt.Errorf("%s", out)
	}

	return nil
}
