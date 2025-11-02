package main

import (
	"fmt"
)

const _ = iota

const (
	a1 = iota
	a2
)

const (
	b1, b2 = iota, iota
)

var _ = 1
var x = 1
var (
	y1 = 1
	y2 = 1
)

var (
	z1, z2 = 1, 1
)

func sample() {
	v := 1
	var u = 2
	u, v = v, u
}

type T struct{}

func (t T) method() {
	fmt.Println("hello")
	s := []int{}
	for i := 0; i < 5; i++ {
		s = append(s, i)
	}
	for k, v := range s {
		fmt.Println(k, v)
	}

	score := 85
	switch {
	case score >= 90:
		fmt.Println("Grade A")
	case score >= 80:
		fmt.Println("Grade B")
	case score >= 70:
		fmt.Println("Grade C")
	default:
		fmt.Println("Grade F")
	}

	if score = 1; score == 0 {
		fmt.Println("0")
	} else if score = 2; score == 1 {
		fmt.Println("1")
	} else if score = 3; score == 2 {
		fmt.Println("2")
	} else if score = 4; score == 3 {
		fmt.Println("3")
	} else {
		fmt.Println("4")
	}
}
