// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package match_test

import (
	"go/parser"
	"go/token"
	"path"
	"runtime"
	"testing"

	"github.com/tvastar/gogo/pkg/match"
)

func TestMatcherOnItself(t *testing.T) {
	_, self, _, _ := runtime.Caller(0)
	matchx := path.Join(path.Dir(self), "match.go")
	f1, err1 := parser.ParseFile(token.NewFileSet(), matchx, nil, 0)
	if err1 != nil {
		t.Fatal("Could not parse match.go", err1)
	}
	f2, err2 := parser.ParseFile(token.NewFileSet(), matchx, nil, 0)
	if err2 != nil {
		t.Fatal("Could not parse match.go", err2)
	}
	if !match.Match(f1, f2) {
		t.Error("Match did not match itself!")
	}
}
