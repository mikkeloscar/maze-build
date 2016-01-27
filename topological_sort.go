package main

import (
	"fmt"
)

func TopologicalSort(pkgs []*SrcPkg) ([]*SrcPkg, error) {
	sorter := NewSorter(pkgs)
	return sorter.topologicalSort()
}

type node struct {
	depends map[*node]struct{}
	tmp     bool
	pkg     *SrcPkg
}

// creates a DAG for performing the topilogical sort of pkgs.
func createGraph(pkgs []*SrcPkg) map[*node]struct{} {
	nodes := make(map[*node]struct{}, len(pkgs))
	tmpNodes := make(map[string]*node, len(pkgs))

	for _, pkg := range pkgs {
		node := &node{
			pkg:     pkg,
			depends: make(map[*node]struct{}),
		}

		nodes[node] = struct{}{}

		for _, name := range pkg.PKGBUILD.Pkgnames {
			tmpNodes[name] = node
		}
	}

	for node, _ := range nodes {
		for _, dep := range node.pkg.PKGBUILD.BuildDepends() {
			if n, ok := tmpNodes[dep.Name]; ok {
				// ignore if split package that depend on
				// itself. Otherwise add to depends set
				if n != node {
					n.depends[node] = struct{}{}
				}
			}
		}
	}

	return nodes
}

type sorter struct {
	unmarked map[*node]struct{}
	l        []*SrcPkg
}

func NewSorter(pkgs []*SrcPkg) *sorter {
	return &sorter{
		unmarked: createGraph(pkgs),
		l:        make([]*SrcPkg, 0, len(pkgs)),
	}
}

func (s *sorter) topologicalSort() ([]*SrcPkg, error) {
	for len(s.unmarked) > 0 {
		n := s.getUnmarked()
		err := s.visit(n)
		if err != nil {
			return nil, err
		}
	}

	return s.l, nil
}

func (s *sorter) getUnmarked() *node {
	for n, _ := range s.unmarked {
		return n
	}

	return nil
}

func (s *sorter) visit(n *node) error {
	if n.tmp {
		return fmt.Errorf("not a DAG")
	}

	if _, ok := s.unmarked[n]; !ok {
		return nil
	}

	n.tmp = true
	for m, _ := range n.depends {
		err := s.visit(m)
		if err != nil {
			return err
		}
	}

	// mark by removing from unmarked
	delete(s.unmarked, n)
	n.tmp = false
	// prepend to L
	s.l = append([]*SrcPkg{n.pkg}, s.l...)

	return nil
}
