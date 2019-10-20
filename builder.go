package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	pkgbuild "github.com/mikkeloscar/gopkgbuild"
	log "github.com/sirupsen/logrus"
)

var (
	pkgPatt    = regexp.MustCompile(`([a-z\d@._+]+[a-z\d@._+-]+)-((\d+:)?([\da-z\._+]+-\d+))-(i686|x86_64|any).pkg.tar.xz`)
	pkgSigPatt = regexp.MustCompile(`([a-z\d@._+]+[a-z\d@._+-]+)-((\d+:)?([\da-z\._+]+-\d+))-(i686|x86_64|any).pkg.tar.xz.sig`)
)

// BuiltPkg defines a built package and optional signature file.
type BuiltPkg struct {
	Package   string `json:"package"`
	Signature string `json:"signature"`
}

func (b *BuiltPkg) String() string {
	if b.Signature != "" {
		return fmt.Sprintf("%s (%s)", path.Base(b.Package), path.Base(b.Signature))
	}

	return path.Base(b.Package)
}

// Builder is used to build arch packages.
type Builder struct {
	workdir  string
	repo     *Repo
	Packager string
}

// BuildNew checks what packages to build based on related repo and builds
// those that have been updated.
func (b *Builder) BuildNew(pkgs []string, aur *AUR) ([]*BuiltPkg, error) {
	// initialize pacman.conf with upstream repo
	err := b.setup()
	if err != nil {
		return nil, err
	}

	// make sure environment is up to date
	err = b.update()
	if err != nil {
		return nil, err
	}

	// get packages that should be built
	srcPkgs, err := b.getBuildPkgs(pkgs, aur)
	if err != nil {
		return nil, err
	}

	if len(srcPkgs) == 0 {
		log.Print("All packages up to date, nothing to build")
		return nil, nil
	}

	buildPkgs, err := b.buildPkgs(srcPkgs)
	if err != nil {
		return nil, err
	}

	successLog(buildPkgs)
	return buildPkgs, nil
}

// Write packages built to the log.
func successLog(pkgs []*BuiltPkg) {
	var buf bytes.Buffer
	buf.WriteString("Built packages:")
	for _, pkg := range pkgs {
		buf.WriteString("\n * ")
		buf.WriteString(pkg.String())
	}

	log.Print(buf.String())
}

// setup build environment.
func (b *Builder) setup() error {
	err := addRepoEntry("/etc/pacman.conf.template", b.repo)
	if err != nil {
		return err
	}

	return addPacmanConf("/etc/pacman.conf.template")
}

// Update build environment.
func (b *Builder) update() error {
	log.Printf("Updating packages")
	return runCmd(b.workdir, nil, "sudo", "pacman", "--sync", "--refresh", "--sysupgrade", "--noconfirm")
}

func getBuildPkgsLog(msg string, pkgs []string) {
	var buf bytes.Buffer
	buf.WriteString(msg)
	for _, pkg := range pkgs {
		buf.WriteString("\n * ")
		buf.WriteString(pkg)
	}

	log.Print(buf.String())
}

// Get a sorted list of packages to build.
func (b *Builder) getBuildPkgs(pkgs []string, aur *AUR) ([]*SrcPkg, error) {
	getBuildPkgsLog("Fetching build sources+dependencies for packages:", pkgs)
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

	err := runCmd(pkg.Path, nil, "makepkg", "--nobuild", "--nodeps", "--noprepare", "--noconfirm")
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("makepkg", "--printsrcinfo")
	cmd.Dir = pkg.Path
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	pkgb, err := pkgbuild.ParseSRCINFOContent(out)
	if err != nil {
		return nil, err
	}

	pkg.PKGBUILD = pkgb

	return pkg, nil
}

// Build a list of packages.
func (b *Builder) buildPkgs(pkgs []*SrcPkg) ([]*BuiltPkg, error) {
	buildPkgs := make([]*BuiltPkg, 0, len(pkgs))

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
func (b *Builder) buildPkg(pkg *SrcPkg) ([]*BuiltPkg, error) {
	p := pkg.PKGBUILD
	if len(p.Pkgnames) > 1 || p.Pkgnames[0] != p.Pkgbase {
		log.Printf("Building package %s:(%s)", p.Pkgbase, strings.Join(p.Pkgnames, ", "))
	} else {
		log.Printf("Building package %s", p.Pkgbase)
	}

	var env []string
	if b.Packager != "" {
		env = os.Environ()
		env = append(env, fmt.Sprintf("PACKAGER=%s", b.Packager))
	}

	err := runCmd(pkg.Path, env, "makepkg", "--install", "--syncdeps", "--noconfirm")
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(pkg.Path)
	if err != nil {
		return nil, err
	}

	pkgs := make([]*BuiltPkg, 0, 1)

	for _, f := range files {
		if pkgPatt.MatchString(f.Name()) {
			builtPkg := &BuiltPkg{
				Package: path.Join(pkg.Path, f.Name()),
			}
			pkgs = append(pkgs, builtPkg)
		}
	}

	for _, f := range files {
		if pkgSigPatt.MatchString(f.Name()) {
			for _, p := range pkgs {
				if path.Base(p.Package) == f.Name()[:len(f.Name())-4] {
					p.Signature = path.Join(pkg.Path, f.Name())
				}
			}
		}
	}

	return pkgs, nil
}
