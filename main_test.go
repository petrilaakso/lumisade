package main

import (
	"testing"
)

func TestNewFrame(t *testing.T) {
	x, y := 80, 24
	f := newFrame(x, y)
	fx, fy := f.size()
	if x != fx {
		t.Fatalf("frame reports wrong x size")
	}
	if y != fy {
		t.Fatalf("frame reports wrong y size")
	}
}

func TestFrameAccess(t *testing.T) {
	x, y := 80, 24
	f := newFrame(x, y)
	_ = f.c[x-1][y-1]
}

func TestGfxSize(t *testing.T) {
	s := `111
222
333`
	x, y := gfxSize(s)
	if x != 3 {
		t.Fatalf("gfx wrong x size")
	}
	if y != 3 {
		t.Fatalf("gfx wrong y size")
	}
}

func TestGfxZeroSize(t *testing.T) {
	s := ""
	x, y := gfxSize(s)
	if x != 0 {
		t.Fatalf("gfx wrong x size")
	}
	if y != 0 {
		t.Fatalf("gfx wrong y size")
	}
}

func TestGfxNewlineSize(t *testing.T) {
	s := `111
222
333
`
	x, y := gfxSize(s)
	if x != 3 {
		t.Fatalf("gfx wrong x size")
	}
	if y != 3 {
		t.Fatalf("gfx wrong y size")
	}
}

func TestGfxSizeLongline(t *testing.T) {
	s := `111
222 spammy spam !
333`
	x, y := gfxSize(s)
	if x != 17 {
		t.Fatalf("gfx wrong x size")
	}
	if y != 3 {
		t.Fatalf("gfx wrong y size")
	}
}

func TestNewGfxFrame(t *testing.T) {
	s := `a11
222
33b`
	f := gfxFrame(s)
	if f.c[0][0].primary != 'a' {
		t.Fatalf("invalid value")
	}
	if f.c[2][2].primary != 'b' {
		t.Fatalf("invalid value")
	}
}

func TestNewGfxFrameNewline(t *testing.T) {
	s := `a11
222
33b
`
	f := gfxFrame(s)
	if f.c[0][0].primary != 'a' {
		t.Fatalf("invalid value")
	}
	if f.c[2][2].primary != 'b' {
		t.Fatalf("invalid value")
	}
}

func TestNewGfxFrameLongline(t *testing.T) {
	s := `a11
222 spammy spam !
33b`
	f := gfxFrame(s)
	if f.c[16][1].primary != '!' {
		t.Fatalf("invalid value")
	}
}

func TestDrawAt(t *testing.T) {
	gfx := `axx
xxx
xxb`
	s := gfxFrame(gfx)
	d := newFrame(80, 24)
	d.drawAt(0, 0, s)
	if d.c[0][0].primary != 'a' {
		t.Fatalf("invalid value")
	}
	if d.c[2][2].primary != 'b' {
		t.Fatalf("invalid value")
	}
}

func TestDrawAtBottom(t *testing.T) {
	gfx := `axx
xxx
xxb`
	s := gfxFrame(gfx)
	d := newFrame(10, 10)
	d.drawAt(0, 7, s)
	if d.c[0][7].primary != 'a' {
		t.Fatalf("invalid value")
	}
	if d.c[2][9].primary != 'b' {
		t.Fatalf("invalid value")
	}
}
