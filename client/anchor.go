// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"os"
	"path"
)

// Anchor
type Anchor struct {
	slash string
	walk []string // path to this anchor within anchor subtree
	dir *dir // open directory of this anchor
}

func openAnchor(slash string) (a *Anchor, err error) {
	a = &Anchor{slash: slash}
	if a.dir, err = openDir(slash); err != nil {
		return nil, err
	}
	return
}

// Path returns the path of this anchor within the local file system.
func (a *Anchor) Path() string {
	return path.Join(append([]string{a.slash}, a.walk...)...)
}

// UseAnchor
func (a *Anchor) Anchor(walk ...string) (sub *Anchor) {
	if len(walk) == 0 {
		return a
	}
	switch walk[0] {
	case "chan", "proc", "help":
		panic("subanchors not allowed in element directories")
	}
	os.MkdirAll(path.Join(a.Path(), walk[0]), 0777) // TODO: unused directories are gc'd by the circuit daemon
	sub = &Anchor{
		slash: a.slash,
		walk: append(a.walk, walk[0]),
	}
	var err error
	if sub.dir, err = openDir(sub.Path()); err != nil {
		panic(err)
	}
	return sub.Anchor(walk[1:]...)
}

// UseChan
func (a *Anchor) Chan(name string) *Chan {
	local := path.Join(a.Path(), "chan", name)
	os.MkdirAll(local, 0777)
	return openChan(local)
}

// UseProc
func (a *Anchor) Proc(name string) *Proc {
	local := path.Join(a.Path(), "proc", name)
	os.MkdirAll(local, 0777)
	return openProc(local)
}
