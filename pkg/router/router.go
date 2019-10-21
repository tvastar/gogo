// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package router

import (
	"strings"

	"go/ast"

	"github.com/tvastar/gogo/pkg/code"
)

func New(pkgName, structName string) *Config {
	if structName == "" {
		structName = "Router"
	}

	recv := strings.ToLower(string(([]rune(structName))[:1]))
	if recv == "r" {
		recv = "rr"
	}
	return &Config{
		Package:  pkgName,
		Struct:   structName,
		Receiver: recv,
		Writer:   "w",
		Request:  "r",
		Routes:   nil,
	}
}

type Config struct {
	Package, Struct, Receiver, Writer, Request string
	Routes                                     []code.NodeMarshaler
}

func (c *Config) WithRoutes(routes ...code.NodeMarshaler) *Config {
	c.Routes = append(c.Routes, routes...)
	return c
}

var cfgKey = "config"

func (c *Config) MarshalNode(s *code.Scope) ast.Node {
	s.Stash[&cfgKey] = c
	writer := code.Import("net/http").Dot("ResponseWriter")
	request := code.Import("net/http").Dot("Request").Star()
	fn := code.Func("ServeHTTP").
		WithReceiver(code.Ident(c.Receiver), code.Ident(c.Struct), nil).
		WithParam(code.Ident(c.Writer), writer, nil).
		WithParam(code.Ident(c.Request), request, nil).
		WithBody(c.Routes...)
	return code.File(c.Package, fn).MarshalNode(s)
}

func FromScope(s *code.Scope) *Config {
	x, _ := s.LookupStash(&cfgKey)
	return x.(*Config)
}

func Writer() code.NodeMarshaler {
	return code.MarshalerFunc(func(s *code.Scope) ast.Node {
		c := FromScope(s)
		return ast.NewIdent(c.Writer)
	})
}

func StatusCode(status int) code.NodeMarshaler {
	return Writer().Dot("WriteHeader").Call(code.Literal(status))
}
