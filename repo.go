package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Repo struct {
	name     string
	url      string
	workdir  string
	dbPath   string
	dbCached bool
}

func (r *Repo) getDB() (string, error) {
	if r.dbCached {
		return path.Join(r.workdir, r.dbPath), nil
	}

	fileName := fmt.Sprintf("%s.db.tar.gz", r.name)

	if strings.HasPrefix(r.url, "file://") {
		if r.url[len(r.url)-1] == '/' {
			r.dbPath = r.url + fileName
		} else {
			r.dbPath = r.url + "/" + fileName
		}
		r.dbCached = true
		return r.dbPath, nil
	}

	if strings.HasPrefix(r.url, "http://") {
		// TODO: handle more db naming
		filePath, err := r.httpDownload(fileName)
		if err != nil {
			return "", err
		}

		r.dbPath = filePath
		r.dbCached = true
		return filePath, nil
	}

	return "", fmt.Errorf("invalid url '%s'", r.url)
}

func (r *Repo) httpDownload(file string) (string, error) {
	filePath := path.Join(r.workdir, file)
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Get data
	resp, err := http.Get(r.url)
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

func (r *Repo) GetUpdated(pkgs []*SrcPkg) ([]*SrcPkg, error) {
	updated := make([]*SrcPkg, 0, len(pkgs))
	for _, pkg := range pkgs {
		new, err := r.IsNew(pkg)
		if err != nil {
			return nil, err
		}

		if new {
			updated = append(updated, pkg)
		}
	}

	return TopologicalSort(updated)
}

// AddLocal adds a list of packages to a repo db, moving the package files to
// the repo db directory if needed.
func (r *Repo) AddLocal(pkgPaths []string) error {
	dbPath, err := r.getDB()
	if err != nil {
		return err
	}

	for i, pkg := range pkgPaths {
		pkgPathDir, pkgPathBase := path.Split(pkg)

		if r.workdir != pkgPathDir {
			// move pkg to repo path.
			newPath := path.Join(r.workdir, pkgPathBase)
			err := os.Rename(pkg, newPath)
			if err != nil {
				return err
			}
			pkgPaths[i] = newPath
		}
	}

	args := []string{"-R", dbPath}
	args = append(args, pkgPaths...)

	cmd := exec.Command("repo-add", args...)
	cmd.Dir = r.workdir

	return cmd.Run()
}

func (r *Repo) movePkgFile(db, pkgPath string) error {

	return nil
}

// IsNew returns true if pkg is a newer version than what's in the repo.
// If the package is not found in the repo, it will be marked as new.
func (r *Repo) IsNew(pkg *SrcPkg) (bool, error) {
	dbPath, err := r.getDB()
	if err != nil {
		return false, err
	}

	f, err := os.Open(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}
	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		return false, err
	}

	tarR := tar.NewReader(gzf)

	for {
		header, err := tarR.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return false, err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			n, v := splitNameVersion(header.Name)
			for _, name := range pkg.PKGBUILD.Pkgnames {
				if n == name {
					version := pkg.PKGBUILD.CompleteVersion()
					if version.Newer(v) {
						return true, nil
					}
					return false, nil
				}
			}
		case tar.TypeReg:
			continue
		}
	}

	return true, nil
}

// turn "zlib-1.2.8-4/" into ("zlib", "1.2.8-4").
func splitNameVersion(str string) (string, string) {
	chars := strings.Split(str[:len(str)-1], "-")
	name := chars[:len(chars)-2]
	version := chars[len(chars)-2:]

	return strings.Join(name, "-"), strings.Join(version, "-")
}
