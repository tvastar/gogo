// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package match implements matching ASTs against one another
package match

import (
	"go/ast"
	"go/token"
	"log"
)

type Matcher interface {
	Matches(other interface{}) bool
}

func Match(left, right interface{}) bool {
	if left == right {
		return true
	}

	if m, ok := left.(Matcher); ok {
		return m.Matches(right)
	}

	if m, ok := right.(Matcher); ok {
		return m.Matches(left)
	}

	switch l := left.(type) {
	case *ast.FieldList:
		r, ok := right.(*ast.FieldList)
		if !ok || len(l.List) != len(r.List) {
			log.Printf("Match failed %#v %#v\n", left, right)
			return false
		}
		for kk := range l.List {
			if !Match(l.List[kk], r.List[kk]) {
				log.Printf("Match failed %#v %#v\n", left, right)
				return false
			}
		}
		return true
	case *ast.Field:
		if r, ok := right.(*ast.Field); ok {
			return Match(&l.Names, &r.Names) && Match(l.Type, r.Type)
		}
	case *[]*ast.Ident:
		r, ok := right.(*[]*ast.Ident)
		if !ok || len(*l) != len(*r) {
			log.Printf("Match failed %#v %#v\n", left, right)
			return false
		}
		for kk := range *l {
			if !Match((*l)[kk], (*r)[kk]) {
				log.Printf("Match failed %#v %#v\n", left, right)
				return false
			}
		}
		return true
	case *ast.Ident:
		r, ok := right.(*ast.Ident)
		return ok && Match(l.Name, r.Name)
	case *ast.BasicLit:
		r, ok := right.(*ast.BasicLit)
		return ok && Match(l.Value, r.Value)
	case *ast.Ellipsis:
		r, ok := right.(*ast.Ellipsis)
		return ok && Match(l.Elt, r.Elt)
	case *ast.FuncLit:
		r, ok := right.(*ast.FuncLit)
		return ok && Match(l.Type, r.Type) && Match(l.Body, r.Body)
	case *ast.CompositeLit:
		r, ok := right.(*ast.CompositeLit)
		return ok && Match(l.Type, r.Type) && Match(&l.Elts, &r.Elts)
	case *[]ast.Expr:
		r, ok := right.(*[]ast.Expr)
		if !ok || len(*l) != len(*r) {
			log.Printf("Match failed %#v %#v\n", left, right)
			return false
		}
		for kk := range *l {
			if !Match((*l)[kk], (*r)[kk]) {
				return false
			}
		}
		return true
	case *ast.ParenExpr:
		r, ok := right.(*ast.ParenExpr)
		return ok && Match(l.X, r.X)
	case *ast.SelectorExpr:
		r, ok := right.(*ast.SelectorExpr)
		return ok && Match(l.X, r.X) && Match(l.Sel, r.Sel)
	case *ast.IndexExpr:
		r, ok := right.(*ast.IndexExpr)
		return ok && Match(l.X, r.X) && Match(l.Index, r.Index)
	case *ast.SliceExpr:
		r, ok := right.(*ast.SliceExpr)
		return ok && Match(l.X, r.X) && Match(l.Low, r.Low) && Match(l.High, r.High) && Match(l.Max, r.Max)
	case *ast.TypeAssertExpr:
		r, ok := right.(*ast.TypeAssertExpr)
		return ok && Match(l.X, r.X) && Match(l.Type, r.Type)
	case *ast.CallExpr:
		r, ok := right.(*ast.CallExpr)
		return ok && (l.Ellipsis == token.NoPos) == (r.Ellipsis == token.NoPos) && Match(l.Fun, r.Fun) && Match(&l.Args, &r.Args)
	case *ast.StarExpr:
		r, ok := right.(*ast.StarExpr)
		return ok && Match(l.X, r.X)
	case *ast.UnaryExpr:
		r, ok := right.(*ast.UnaryExpr)
		return ok && l.Op == r.Op && Match(l.X, r.X)
	case *ast.BinaryExpr:
		r, ok := right.(*ast.BinaryExpr)
		return ok && l.Op == r.Op && Match(l.X, r.X) && Match(l.Y, r.Y)
	case *ast.KeyValueExpr:
		r, ok := right.(*ast.KeyValueExpr)
		return ok && Match(l.Key, r.Key) && Match(l.Value, r.Value)
	case *ast.ArrayType:
		r, ok := right.(*ast.ArrayType)
		return ok && Match(l.Len, r.Len) && Match(l.Elt, r.Elt)
	case *ast.StructType:
		r, ok := right.(*ast.StructType)
		return ok && Match(l.Fields, r.Fields)
	case *ast.FuncType:
		r, ok := right.(*ast.FuncType)
		return ok && Match(l.Params, r.Params) && Match(l.Results, r.Results)
	case *ast.InterfaceType:
		r, ok := right.(*ast.InterfaceType)
		return ok && Match(l.Methods, r.Methods)
	case *ast.MapType:
		r, ok := right.(*ast.MapType)
		return ok && Match(l.Key, r.Key) && Match(l.Value, r.Value)
	case *ast.ChanType:
		r, ok := right.(*ast.ChanType)
		return ok && Match(l.Value, r.Value)
	case *ast.DeclStmt:
		r, ok := right.(*ast.DeclStmt)
		return ok && Match(l.Decl, r.Decl)
	case *ast.EmptyStmt:
		_, ok := right.(*ast.EmptyStmt)
		return ok
	case *ast.LabeledStmt:
		r, ok := right.(*ast.LabeledStmt)
		return ok && Match(l.Label, r.Label) && Match(l.Stmt, r.Stmt)
	case *ast.ExprStmt:
		r, ok := right.(*ast.ExprStmt)
		return ok && Match(l.X, r.X)
	case *ast.SendStmt:
		r, ok := right.(*ast.SendStmt)
		return ok && Match(l.Chan, r.Chan) && Match(l.Value, r.Value)
	case *ast.IncDecStmt:
		r, ok := right.(*ast.IncDecStmt)
		return ok && l.Tok == r.Tok && Match(l.X, r.X)
	case *ast.AssignStmt:
		r, ok := right.(*ast.AssignStmt)
		return ok && l.Tok == r.Tok && Match(&l.Lhs, &r.Lhs) && Match(&l.Rhs, &r.Rhs)
	case *ast.GoStmt:
		r, ok := right.(*ast.GoStmt)
		return ok && Match(l.Call, r.Call)
	case *ast.DeferStmt:
		r, ok := right.(*ast.DeferStmt)
		return ok && Match(l.Call, r.Call)
	case *ast.ReturnStmt:
		r, ok := right.(*ast.ReturnStmt)
		return ok && Match(&l.Results, &r.Results)
	case *ast.BranchStmt:
		r, ok := right.(*ast.BranchStmt)
		return ok && l.Tok == r.Tok && Match(l.Label, r.Label)
	case *ast.BlockStmt:
		r, ok := right.(*ast.BlockStmt)
		return ok && Match(&l.List, &r.List)
	case *[]ast.Stmt:
		r, ok := right.(*[]ast.Stmt)
		if !ok || len(*l) != len(*r) {
			log.Printf("Match failed %#v %#v\n", left, right)
			return false
		}
		for kk := range *l {
			if !Match((*l)[kk], (*r)[kk]) {
				log.Printf("Match failed %#v %#v\n", left, right)
				return false
			}
		}
		return true
	case *ast.IfStmt:
		r, ok := right.(*ast.IfStmt)
		return ok && Match(l.Init, r.Init) && Match(l.Cond, r.Cond) && Match(l.Body, r.Body) && Match(l.Else, r.Else)
	case *ast.CaseClause:
		r, ok := right.(*ast.CaseClause)
		return ok && Match(&l.List, &r.List) && Match(&l.Body, &r.Body)

	case *ast.SwitchStmt:
		r, ok := right.(*ast.SwitchStmt)
		return ok && Match(l.Init, r.Init) && Match(l.Tag, r.Tag) && Match(l.Body, r.Body)
	case *ast.TypeSwitchStmt:
		r, ok := right.(*ast.TypeSwitchStmt)
		return ok && Match(l.Init, r.Init) && Match(l.Assign, r.Assign) && Match(l.Body, r.Body)

	case *ast.CommClause:
		r, ok := right.(*ast.CommClause)
		return ok && Match(l.Comm, r.Comm) && Match(l.Body, r.Body)

	case *ast.SelectStmt:
		r, ok := right.(*ast.SelectStmt)
		return ok && Match(l.Body, r.Body)

	case *ast.ForStmt:
		r, ok := right.(*ast.ForStmt)
		return ok && Match(l.Init, r.Init) && Match(l.Cond, r.Cond) && Match(l.Post, r.Post) && Match(l.Body, r.Body)
	case *ast.RangeStmt:
		r, ok := right.(*ast.RangeStmt)
		return ok && l.Tok == r.Tok && Match(l.Key, r.Key) &&
			Match(l.Value, r.Value) && Match(l.X, r.X) &&
			Match(l.Body, r.Body)
	case *ast.ImportSpec:
		r, ok := right.(*ast.ImportSpec)
		return ok && Match(l.Name, r.Name) && Match(l.Path, r.Path)

	case *ast.ValueSpec:
		r, ok := right.(*ast.ValueSpec)
		return ok && Match(&l.Names, &r.Names) &&
			Match(&l.Values, &r.Values) && Match(l.Type, r.Type)
	case *ast.TypeSpec:
		r, ok := right.(*ast.TypeSpec)
		return ok && Match(l.Name, r.Name) && Match(l.Type, r.Type)

	case *ast.GenDecl:
		r, ok := right.(*ast.GenDecl)
		return ok && l.Tok == r.Tok && Match(&l.Specs, &r.Specs)
	case *[]ast.Spec:
		r, ok := right.(*[]ast.Spec)
		if !ok || len(*l) != len(*r) {
			log.Printf("Match failed %#v %#v\n", left, right)
			return false
		}
		for kk := range *l {
			if !Match((*l)[kk], (*r)[kk]) {
				log.Printf("Match failed %#v %#v\n", left, right)
				return false
			}
		}
		return true
	case *ast.FuncDecl:
		r, ok := right.(*ast.FuncDecl)
		return ok && Match(l.Name, r.Name) && Match(l.Recv, r.Recv) &&
			Match(l.Type, r.Type) && Match(l.Body, r.Body)
	case *ast.File:
		r, ok := right.(*ast.File)
		return ok && Match(l.Name, r.Name) && Match(&l.Decls, &r.Decls)
	case *[]ast.Decl:
		r, ok := right.(*[]ast.Decl)
		if !ok || len(*l) != len(*r) {
			log.Printf("Match failed %#v %#v\n", left, right)
			return false
		}
		for kk := range *l {
			if !Match((*l)[kk], (*r)[kk]) {
				log.Printf("Match failed %#v %#v\n", left, right)
				return false
			}
		}
		return true
	}

	// unexpected node type
	log.Printf("Match failed %#v %#v\n", left, right)
	return false
}
