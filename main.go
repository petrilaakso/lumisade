package main

import (
	_ "embed"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

var snow = [...]rune{'*', '✢', '✣', '✤', '✥', '✱', '✲', '✵', '✻', '✼', '❅', '❆', '❉', '❊', '❋'}

type cell struct {
	primary   rune
	combining []rune
	style     tcell.Style
}

type frame struct {
	c [][]cell
}

func newFrame(x, y int) *frame {
	f := frame{}
	c := make([][]cell, x)
	for i := range c {
		c[i] = make([]cell, y)
	}
	f.c = c
	return &f
}

func (f *frame) size() (x int, y int) {
	x = len(f.c)
	y = len(f.c[0])
	return
}

// Return minimum frame size where gfx fits
func gfxSize(gfx string) (x, y int) {
	xi := 0
	if len(gfx) > 0 {
		y++
	}
	for i := range gfx {
		if gfx[i] == '\n' {
			y++
			if i == len(gfx)-1 { // don't count newline at the end
				y--
			}
			if xi > x {
				x = xi
			}
			xi = 0
		} else {
			xi++
		}
	}
	return
}

func gfxFrame(gfx string) *frame {
	f := newFrame(gfxSize(gfx))
	x, y := 0, 0
	for i := range gfx {
		if gfx[i] == '\n' {
			y++
			x = 0
		} else {
			f.c[x][y].primary = rune(gfx[i])
			x++
		}
	}
	return f
}

func (f *frame) draw(p *frame) {
	f.drawAt(0, 0, p)
}

func (d *frame) drawAt(x, y int, s *frame) {
	maxx, maxy := d.size()
	for sx := range s.c {
		for sy := range s.c[sx] {
			if x+sx < maxx && y+sy < maxy {
				if s.c[sx][sy].primary != 0 {
					d.c[x+sx][y+sy] = s.c[sx][sy]
				}
			}
		}
	}
}

func rain(f *frame) *frame {
	nf := newFrame(f.size())
	_, nfymax := nf.size()
	for x := range f.c {
		for y := range f.c[x] {
			if y+1 >= nfymax {
				break
			}
			nf.c[x][y+1] = f.c[x][y]
		}
	}
	f.c = nf.c

	for x := 0; x < len(f.c); x++ {
		if rand.Int()%10 == 0 {
			f.c[x][0].primary = snow[rand.Intn(len(snow))]
		}
	}

	return f
}

func (f *frame) set(s tcell.Screen) {
	for x := range f.c {
		for y := range f.c[x] {
			s.SetContent(x, y, f.c[x][y].primary, f.c[x][y].combining, f.c[x][y].style)
		}
	}
}

//go:embed lb1.txt
var lb1 string

//go:embed lb2.txt
var lb2 string

func main() {
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	snow := newFrame(s.Size())
	town1 := gfxFrame(lb1)
	town2 := gfxFrame(lb2)
	tx, ty := town1.size()

	// periodic refresh event
	quit := make(chan struct{})
	go func(quit <-chan struct{}) {
		ticker := time.NewTicker(500 * time.Millisecond)
		for {
			select {
			case <-quit:
				ticker.Stop()
				return
			case <-ticker.C:
				s.PostEvent(tcell.NewEventKey(tcell.KeyRune, rune('r'), 0))
			}
		}
	}(quit)

loop:
	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			snow = newFrame(s.Size())
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
				quit <- struct{}{}
				break loop
			} else if ev.Rune() == 'r' || ev.Rune() == ' ' {
				world := newFrame(s.Size())
				snow = rain(snow)
				_, wy := world.size()
				world.drawAt(0, wy-ty, town1)
				world.drawAt(tx*1, wy-ty, town2) // Spam towns
				world.drawAt(tx*2, wy-ty, town2)
				world.drawAt(tx*3, wy-ty, town2)
				world.drawAt(tx*4, wy-ty, town2)
				world.drawAt(tx*5, wy-ty, town2)
				world.draw(snow)
				world.set(s)
				s.Show()
			}
		}
	}

	s.Fini()
	os.Exit(0)
}
