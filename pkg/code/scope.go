// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package code

import (
	"go/ast"
	"strconv"
)

// Creates a root scope object
func RootScope() *Scope {
	var s *Scope
	return s.New()
}

// Scope tracks all used variable names
//
// It also allows stashing arbitrary "context"
type Scope struct {
	Stash  map[interface{}]interface{}
	Vars   map[string]ast.Node
	Parent *Scope
}

// New creates a new nested scope
func (s *Scope) New() *Scope {
	return &Scope{
		Stash:  map[interface{}]interface{}{},
		Vars:   map[string]ast.Node{},
		Parent: s,
	}
}

// LookupStash looks up the stash (up the parent chain) for a key
func (s *Scope) LookupStash(key interface{}) (interface{}, bool) {
	if s == nil {
		return nil, false
	}
	if v, ok := s.Stash[key]; ok {
		return v, ok
	}
	return s.Parent.LookupStash(key)
}

// LookupVar looks up the scope chain to find a variable with the given name
func (s *Scope) LookupVar(name string) (ast.Node, bool) {
	if s == nil {
		return nil, false
	}
	if v, ok := s.Vars[name]; ok {
		return v, ok
	}
	return s.Parent.LookupVar(name)
}

// PickName picks a unique name with the given prefix
func (s *Scope) PickName(prefix string) string {
	if prefix == "" {
		prefix = "gogox"
	}

	name, idx := prefix, 2
	for {
		if _, ok := s.LookupVar(name); !ok {
			return name
		}
		name = prefix + strconv.Itoa(idx)
		idx++
	}
}
