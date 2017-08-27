package main

import (
	"errors"
	"fmt"
	"path"

	"github.com/mikkeloscar/aur"
	"github.com/mikkeloscar/gopkgbuild"
)

const aurCloneURL = "https://aur.archlinux.org/%s.git"

// AUR is used to fetch git sources from AUR.
type AUR struct {
	workdir string
}

// SrcPkg describes a source package including path to basedir of the package
// and a parsed version of the PKGBUILD.
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
	for d := range deps {
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
	pkgsInfo, err := aur.Info(pkgs)
	if err != nil {
		return err
	}

	for _, pkg := range pkgsInfo {
		updates[pkg.PackageBase] = struct{}{}

		// TODO: maybe add optdepends
		depends := make([]string, 0, len(pkg.Depends)+len(pkg.MakeDepends))
		err := addDeps(&depends, pkg.Depends)
		if err != nil {
			return err
		}
		err = addDeps(&depends, pkg.MakeDepends)
		if err != nil {
			return err
		}
		a.getDeps(depends, updates)
	}

	return nil
}

// parses a string slice of dependencies and adds them to the combinedDepends
// slice.
func addDeps(combinedDepends *[]string, deps []string) error {
	parsedDeps, err := pkgbuild.ParseDeps(deps)
	if err != nil {
		return err
	}

	for _, dep := range parsedDeps {
		*combinedDepends = append(*combinedDepends, dep.Name)
	}

	return nil
}

// get source repos from set of package names
func (a *AUR) getSourceRepos(pkgs map[string]struct{}) error {
	clone := make(chan error)

	for pkg := range pkgs {
		go a.updateRepo(pkg, clone)
	}

	var errs []error

	for range pkgs {
		err := <-clone
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		msg := "errors while fetching sources: from AUR\n"
		for _, err := range errs {
			msg += fmt.Sprintf("%s * %s\n", msg, err.Error())
		}
		return errors.New(msg)
	}

	return nil
}

// update (clone or pull) AUR package repo
func (a *AUR) updateRepo(pkg string, c chan<- error) {
	url := fmt.Sprintf(aurCloneURL, pkg)

	err := gitClone(url, path.Join(a.workdir, pkg))
	c <- err
}
