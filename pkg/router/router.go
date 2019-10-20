// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package router

import (
	"strings"

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
	}
}

type Config struct {
	Package, Struct, Receiver, Writer, Request string
}

type Route interface {
	Route(c *Config) code.NodeMarshaler
}

func (c *Config) Generate(routes ...Route) code.NodeMarshaler {
	var stmts []code.NodeMarshaler
	for _, r := range routes {
		stmts = append(stmts, r.Route(c))
	}

	writer := code.Ident("http").Dot("ResponseWriter")
	request := code.Ident("http").Dot("Request").Star()
	fn := code.Func("ServeHTTP").
		WithReceiver(code.Ident(c.Receiver), code.Ident(c.Struct), nil).
		WithParam(code.Ident(c.Writer), writer, nil).
		WithParam(code.Ident(c.Request), request, nil).
		WithBody(stmts...)
	return code.File(c.Package, fn)
}

type StatusCode int

func (s StatusCode) Route(c *Config) code.NodeMarshaler {
	return code.Ident(c.Writer).Dot("WriteHeader").Call(code.Literal(int(s)))
}
