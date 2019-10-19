// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package code

import (
	"go/ast"
	"go/token"
)

// NodeMarshaler
type NodeMarshaler interface {
	MarshalNode(s *Scope) ast.Node

	// Paren adds a paren around the expression
	Paren() NodeMarshaler

	// Call represents a function call expression
	Call(args ...NodeMarshaler) NodeMarshaler

	// Then updates the then statement
	Then(stmts ...NodeMarshaler) NodeMarshaler

	// Op represents a binary operation such as "<"
	Op(op string, o NodeMarshaler) NodeMarshaler

	// Assign represents an assignment op such as ":="
	// Use code.Assign for multiple simultaneous assignment
	Assign(op string, o NodeMarshaler) NodeMarshaler

	// WithReceiver adds a receiver to the fn
	WithReceiver(args ...NodeMarshaler) NodeMarshaler

	// WithParam adds a param to the fn
	WithParam(args ...NodeMarshaler) NodeMarshaler

	// WithResult adds a result to the fn
	WithResult(args ...NodeMarshaler) NodeMarshaler

	// WithBody adds a body to the fn
	WithBody(stmt ...NodeMarshaler) NodeMarshaler
}

type nodef func(s *Scope) ast.Node

func (n nodef) MarshalNode(s *Scope) ast.Node {
	return n(s)
}

func (n nodef) Op(op string, o NodeMarshaler) NodeMarshaler {
	var tok token.Token
	for kk := token.ILLEGAL; kk <= token.VAR; kk++ {
		if kk.String() == op {
			tok = kk
		}
	}
	return nodef(func(s *Scope) ast.Node {
		x := n.MarshalNode(s).(ast.Expr)
		if o == nil {
			return &ast.UnaryExpr{X: x, Op: tok}
		}
		y := o.MarshalNode(s).(ast.Expr)
		return &ast.BinaryExpr{X: x, Op: tok, Y: y}
	})
}

func (n nodef) Assign(op string, o NodeMarshaler) NodeMarshaler {
	return Assign(op, n, o)
}

func (n nodef) Paren() NodeMarshaler {
	return nodef(func(s *Scope) ast.Node {
		return &ast.ParenExpr{X: n.MarshalNode(s).(ast.Expr)}
	})

}

func (n nodef) Call(args ...NodeMarshaler) NodeMarshaler {
	return nodef(func(s *Scope) ast.Node {
		fn := n.MarshalNode(s).(ast.Expr)
		exprs := make([]ast.Expr, len(args))
		for kk, arg := range args {
			exprs[kk] = arg.MarshalNode(s).(ast.Expr)
		}
		if len(args) == 0 {
			exprs = nil
		}
		return &ast.CallExpr{Fun: fn, Args: exprs}
	})

}

func (n nodef) Then(stmts ...NodeMarshaler) NodeMarshaler {
	return nodef(func(s *Scope) ast.Node {
		ifstmt := n.MarshalNode(s).(*ast.IfStmt)
		block := &ast.BlockStmt{}
		for _, stmt := range stmts {
			nn := stmt.MarshalNode(s)
			if x, ok := nn.(ast.Expr); ok {
				block.List = append(block.List, &ast.ExprStmt{X: x})
			} else {
				block.List = append(block.List, nn.(ast.Stmt))
			}
		}
		ifstmt.Body = block
		return ifstmt
	})
}

func (n nodef) WithReceiver(args ...NodeMarshaler) NodeMarshaler {
	return nodef(func(s *Scope) ast.Node {
		fn := n.MarshalNode(s).(*ast.FuncDecl)
		fn.Recv.List = append(fn.Recv.List, field(s, args...))
		return fn
	})
}

func (n nodef) WithParam(args ...NodeMarshaler) NodeMarshaler {
	return nodef(func(s *Scope) ast.Node {
		fn := n.MarshalNode(s).(*ast.FuncDecl)
		fn.Type.Params.List = append(fn.Type.Params.List, field(s, args...))
		return fn
	})
}

func (n nodef) WithResult(args ...NodeMarshaler) NodeMarshaler {
	return nodef(func(s *Scope) ast.Node {
		fn := n.MarshalNode(s).(*ast.FuncDecl)
		fn.Type.Results.List = append(fn.Type.Results.List, field(s, args...))
		return fn
	})
}

func (n nodef) WithBody(stmts ...NodeMarshaler) NodeMarshaler {
	return nodef(func(s *Scope) ast.Node {
		fn := n.MarshalNode(s).(*ast.FuncDecl)
		block := &ast.BlockStmt{}
		for _, stmt := range stmts {
			nn := stmt.MarshalNode(s)
			if x, ok := nn.(ast.Expr); ok {
				block.List = append(block.List, &ast.ExprStmt{X: x})
			} else {
				block.List = append(block.List, nn.(ast.Stmt))
			}
		}
		fn.Body = block
		return fn
	})
}

func field(s *Scope, args ...NodeMarshaler) *ast.Field {
	f := &ast.Field{Names: []*ast.Ident{}}

	for kk, arg := range args {
		if arg == nil {
			continue
		}
		n := arg.MarshalNode(s)
		if n == nil {
			continue
		}

		switch kk {
		case len(args) - 1:
			f.Tag = n.(*ast.BasicLit)
		case len(args) - 2:
			f.Type = n.(ast.Expr)
		default:
			f.Names = append(f.Names, n.(*ast.Ident))
		}
	}
	return f
}
