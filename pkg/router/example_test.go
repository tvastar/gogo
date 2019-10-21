package router_test

import (
	"github.com/tvastar/gogo/pkg/code"
	"github.com/tvastar/gogo/pkg/router"

	"bytes"
	"fmt"
	"go/format"
	"go/token"
	"net/http"
)

func Example() {
	r := router.New("example", "ex").Generate(
		router.StatusCode(http.StatusOK),
	)

	var buf bytes.Buffer
	node := r.MarshalNode(code.RootScope())
	if err := format.Node(&buf, &token.FileSet{}, node); err != nil {
		fmt.Println("Unexpected error", err)
	}

	fmt.Println(buf.String())

	// Output:
	// package example
	//
	// import "net/http"
	//
	// func (e ex) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(200)
	// }
}
