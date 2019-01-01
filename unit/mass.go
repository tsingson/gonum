// Code generated by "go generate github.com/tsingson/gonum/unit”; DO NOT EDIT.

// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unit

import (
	"errors"
	"fmt"
	"math"
)

// Mass represents a mass in kilograms
type Mass float64

const (
	Yottagram Mass = 1e21
	Zettagram Mass = 1e18
	Exagram   Mass = 1e15
	Petagram  Mass = 1e12
	Teragram  Mass = 1e9
	Gigagram  Mass = 1e6
	Megagram  Mass = 1e3
	Kilogram  Mass = 1.0
	Hectogram Mass = 1e-1
	Decagram  Mass = 1e-2
	Gram      Mass = 1e-3
	Decigram  Mass = 1e-4
	Centigram Mass = 1e-5
	Milligram Mass = 1e-6
	Microgram Mass = 1e-9
	Nanogram  Mass = 1e-12
	Picogram  Mass = 1e-15
	Femtogram Mass = 1e-18
	Attogram  Mass = 1e-21
	Zeptogram Mass = 1e-24
	Yoctogram Mass = 1e-27
)

// Unit converts the Mass to a *Unit
func (m Mass) Unit() *Unit {
	return New(float64(m), Dimensions{
		MassDim: 1,
	})
}

// Mass allows Mass to implement a Masser interface
func (m Mass) Mass() Mass {
	return m
}

// From converts the unit into the receiver. From returns an
// error if there is a mismatch in dimension
func (m *Mass) From(u Uniter) error {
	if !DimensionsMatch(u, Gram) {
		*m = Mass(math.NaN())
		return errors.New("Dimension mismatch")
	}
	*m = Mass(u.Unit().Value())
	return nil
}

func (m Mass) Format(fs fmt.State, c rune) {
	switch c {
	case 'v':
		if fs.Flag('#') {
			fmt.Fprintf(fs, "%T(%v)", m, float64(m))
			return
		}
		fallthrough
	case 'e', 'E', 'f', 'F', 'g', 'G':
		p, pOk := fs.Precision()
		w, wOk := fs.Width()
		switch {
		case pOk && wOk:
			fmt.Fprintf(fs, "%*.*"+string(c), w, p, float64(m))
		case pOk:
			fmt.Fprintf(fs, "%.*"+string(c), p, float64(m))
		case wOk:
			fmt.Fprintf(fs, "%*"+string(c), w, float64(m))
		default:
			fmt.Fprintf(fs, "%"+string(c), float64(m))
		}
		fmt.Fprint(fs, " kg")
	default:
		fmt.Fprintf(fs, "%%!%c(%T=%g kg)", c, m, float64(m))
	}
}
