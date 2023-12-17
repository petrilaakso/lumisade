package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

var snow = []rune{'*', '✢', '✣', '✤', '✥', '✱', '✲', '✵', '✻', '✼', '❅', '❆', '❉', '❊', '❋'}

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
		_, c, st, _ := s.GetContent(xi, 0)
		if rand.Int()%5 == 0 {
			s.SetContent(xi, 0, snow[rand.Intn(len(snow))], c, st)
		}
	}
}

func main() {
	s, e := tcell.NewScreen()
	if e != nil {
		panic(e)
	}
	e = s.Init()
	defer s.Fini()
	if e != nil {
		panic(e)
	}

	quit := make(chan struct{})
	go func(quit chan struct{}) {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-quit:
				return
			case <-ticker.C:
				s.PostEvent(tcell.NewEventKey(tcell.KeyRune, rune('r'), 0))
			}
		}
	}(quit)

	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit <- struct{}{}
				s.Fini()
				os.Exit(0)
			} else if ev.Rune() == 'r' {
				moveSnow(s)
				makeSnow(s)
				s.Show()
			}
		}
	}
}
