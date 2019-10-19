// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package code_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/tvastar/gogo/pkg/code"

	"bytes"
	"go/format"
	"go/token"
	"testing"
)

func TestExpr(t *testing.T) {
	validate := func(test string, expr string, noder code.NodeMarshaler) {
		var buf bytes.Buffer
		result := noder.MarshalNode(code.RootScope())
		if err := format.Node(&buf, &token.FileSet{}, result); err != nil {
			t.Fatal("format error", err)
		}
		if diff := cmp.Diff(expr, buf.String()); diff != "" {
			t.Error("mismatch", diff)
		}
	}

	validate("ident", "x", code.Ident("x"))
	validate("literal int", "5", code.Literal(5))
	validate("literal int", "-5", code.Literal(-5))
	validate("literal float", "5.2", code.Literal(5.2))
	validate("literal chara", "'x'", code.Rune('x'))
	validate("+", "x + y", code.Ident("x").Op("+", code.Ident("y")))
	validate("unary", "-x", code.Ident("x").Op("-", nil))
	validate("<", "x < y", code.Ident("x").Op("<", code.Ident("y")))
	validate("ident prefix", "x < x2", code.IdentPrefix("x").Op("<", code.IdentPrefix("x")))
	x2 := code.IdentPrefix("x")
	validate("ident prefix", "x < x2 < x2", code.IdentPrefix("x").Op("<", x2).Op("<", x2))
	validate("<=", "x <= y", code.Ident("x").Op("<=", code.Ident("y")))

	validate("()", "(x)", code.Ident("x").Paren())
	validate("call", "x()", code.Ident("x").Call())
	validate("call", "x(y)", code.Ident("x").Call(code.Ident("y")))
	validate("nil", "nil", code.Nil())
	validate("then", "if x {\n\ty\n}", code.If(code.Ident("x")).Then(code.Ident("y")))

	validate("then2", "if x {\n\ty\n\tz\n}",
		code.If(code.Ident("x")).
			Then(code.Ident("y"), code.Ident("z")))

	validate("assign", "x = y", code.Assign("=", code.Ident("x"), code.Ident("y")))

	validate("if2", "if x := n; x < y {\n\tz\n}",
		code.If2(code.Assign(":=", code.Ident("x"), code.Ident("n")),
			code.Ident("x").Op("<", code.Ident("y"))).
			Then(code.Ident("z")))

}
