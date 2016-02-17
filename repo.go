package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/mikkeloscar/maze/repo"
)

// Repo defines a local copy of a remote arch repository.
type Repo struct {
	url      string
	dbCached bool
	local    repo.Repo
}

// Fetch db, downloads it first if needed.
func (r *Repo) fetchDB() error {
	if r.dbCached {
		return nil
	}

	fileName := fmt.Sprintf("%s.db.tar.gz", r.local.Name)

	if strings.HasPrefix(r.url, "http://") || strings.HasPrefix(r.url, "https://") {
		// TODO: handle more db naming
		_, err := r.httpDownload(fileName)
		if err != nil {
			return err
		}

		r.dbCached = true
		return nil
	}

	arch := "x86_64"

	// else local repo
	if r.url != r.local.PathDeep(arch) {
		err := copyFile(r.local.DB(arch), path.Join(r.url, fileName))
		if err != nil {
			return err
		}
	}

	r.dbCached = true
	return nil
}

// copy file
func copyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

// download db file over http.
func (r *Repo) httpDownload(file string) (string, error) {
	filePath := path.Join(r.local.PathDeep("x86_64"), file)
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Get data
	resp, err := http.Get(fmt.Sprintf("%s/%s/%s", r.url, "x86_64", file))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Write the data to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// GetUpdated takes a list of package sources and returns a sorted list of the
// packages that need to be build (because the source is newer than what's in
// the repo).
func (r *Repo) GetUpdated(pkgs []*SrcPkg) ([]*SrcPkg, error) {
	err := r.fetchDB()
	if err != nil {
		return nil, err
	}

	updated := make([]*SrcPkg, 0, len(pkgs))
	for _, pkg := range pkgs {
		new, err := r.local.IsNew(pkg.PKGBUILD.Pkgbase, "any", pkg.PKGBUILD.CompleteVersion())
		if err != nil {
			return nil, err
		}

		if new {
			updated = append(updated, pkg)
		}
	}

	return TopologicalSort(updated)
}
