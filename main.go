package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func main() {
	s, e := tcell.NewScreen()
	if e != nil {
		panic(e)
	}
	e = s.Init()
	if e != nil {
		panic(e)
	}
	w, h := s.Size()
	e = s.Suspend()
	if e != nil {
		panic(e)
	}
	fmt.Printf("screen is %dx%d\n", w, h)
}
