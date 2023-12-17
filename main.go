package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

var snow = [...]rune{'*', '✢', '✣', '✤', '✥', '✱', '✲', '✵', '✻', '✼', '❅', '❆', '❉', '❊', '❋'}

func moveSnow(s tcell.Screen) {
	x, y := s.Size()
	for yi := y; yi >= 0; yi-- {
		for xi := 0; xi < x; xi++ {
			p, c, st, _ := s.GetContent(xi, yi)
			s.SetContent(xi, yi+1, p, c, st)
			s.SetContent(xi, yi, ' ', c, st)
		}
	}
}

func makeSnow(s tcell.Screen) {
	x, _ := s.Size()
	for xi := 0; xi < x; xi++ {
		if rand.Int()%10 == 0 {
			_, c, st, _ := s.GetContent(xi, 0)
			s.SetContent(xi, 0, snow[rand.Intn(len(snow))], c, st)
		}
	}
}

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
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
				quit <- struct{}{}
				break loop
			} else if ev.Rune() == 'r' {
				moveSnow(s)
				makeSnow(s)
				s.Show()
			}
		}
	}

	s.Fini()
	os.Exit(0)
}
