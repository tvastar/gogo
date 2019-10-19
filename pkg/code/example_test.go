// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package code_test

import (
	"github.com/tvastar/gogo/pkg/code"

	"bytes"
	"fmt"
	"go/format"
	"go/token"
)

func Example() {
	x := code.Ident("x")
	y := code.Ident("y")
	z := code.Ident("z")
	n := code.Ident("n")
	expr := code.If2(x.Assign(":=", n), x.Op("<", y)).
		Then(z)

	var buf bytes.Buffer
	node := expr.MarshalNode(code.RootScope())
	if err := format.Node(&buf, &token.FileSet{}, node); err != nil {
		fmt.Println("Unexpected error", err)
	}

	fmt.Println(buf.String())

	// Output:
	// if x := n; x < y {
	// 	z
	// }
}
