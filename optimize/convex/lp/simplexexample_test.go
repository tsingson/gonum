// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lp_test

import (
	"fmt"
	"log"

	"github.com/tsingson/gonum/mat"
	"github.com/tsingson/gonum/optimize/convex/lp"
)

func ExampleSimplex() {
	c := []float64{-1, -2, 0, 0}
	A := mat.NewDense(2, 4, []float64{-1, 2, 1, 0, 3, 1, 0, 1})
	b := []float64{4, 9}

	opt, x, err := lp.Simplex(c, A, b, 0, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("opt: %v\n", opt)
	fmt.Printf("x: %v\n", x)
	// Output:
	// opt: -8
	// x: [2 3 0 0]
}
