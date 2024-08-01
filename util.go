package pietvm

import (
	"fmt"
	"image"
	"image/color"
	"sort"
)

type DpDir image.Point

var (
	east DpDir = DpDir{1, 0}
	south DpDir = DpDir{0, 1}
	west DpDir = DpDir{-1, 0}
	north DpDir = DpDir{0, -1}
)

type CCDir int

const (
	left CCDir = iota
	right
)

type ColorBlock map[image.Point]struct{}

func (cb ColorBlock) String() string {
	s := make(sort.StringSlice, len(cb))
	for p, _ := range cb {
		s = append(s, p.String())
	}
	sort.Sort(s)
	return fmt.Sprint(s)
}

func (cb ColorBlock) Bounds() (r image.Rectangle) {
	for p, _ := range cb {
		r = image.Rectangle{p, image.Point{p.X + 1, p.Y + 1}}
		break
	}
	for p, _ := range cb {
		r = r.Union(image.Rectangle{p, image.Point{p.X + 1, p.Y + 1}})
	}
	return r
}

func adj(p image.Point) [4]image.Point {
	return [4]image.Point {
		{p.X, p.Y + 1},
		{p.X, p.Y - 1},
		{p.X + 1, p.Y},
		{p.X - 1, p.Y},
	}
}

func SameColors(c1, c2 color.Color) bool {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2
}

func ColorAtPos(img image.Image, pos image.Point) color.Color {
	return img.At(pos.X, pos.Y)
}

func ColorInfo(c color.Color) (hue, lightness int) {
	r, g, b, _ := c.RGBA()
	switch {
	case r == 0xFFFF && g == 0xC0C0 && b == 0xC0C0:
		return 0, 0
	case r == 0xFFFF && g == 0xFFFF && b == 0xC0C0:
		return 1, 0
	case r == 0xC0C0 && g == 0xFFFF && b == 0xC0C0:
		return 2, 0
	case r == 0xC0C0 && g == 0xFFFF && b == 0xFFFF:
		return 3, 0
	case r == 0xC0C0 && g == 0xC0C0 && b == 0xFFFF:
		return 4, 0
	case r == 0xFFFF && g == 0xC0C0 && b == 0xFFFF:
		return 5, 0

	case r == 0xFFFF && g == 0x0000 && b == 0x0000:
		return 0, 1
	case r == 0xFFFF && g == 0xFFFF && b == 0x0000:
		return 1, 1
	case r == 0x0000 && g == 0xFFFF && b == 0x0000:
		return 2, 1
	case r == 0x0000 && g == 0xFFFF && b == 0xFFFF:
		return 3, 1
	case r == 0x0000 && g == 0x0000 && b == 0xFFFF:
		return 4, 1
	case r == 0xFFFF && g == 0x0000 && b == 0xFFFF:
		return 5, 1

	case r == 0xC0C0 && g == 0x0000 && b == 0x0000:
		return 0, 2
	case r == 0xC0C0 && g == 0xC0C0 && b == 0x0000:
		return 1, 2
	case r == 0x0000 && g == 0xC0C0 && b == 0x0000:
		return 2, 2
	case r == 0x0000 && g == 0xC0C0 && b == 0xC0C0:
		return 3, 2
	case r == 0x0000 && g == 0x0000 && b == 0xC0C0:
		return 4, 2
	case r == 0xC0C0 && g == 0x0000 && b == 0xC0C0:
		return 5, 2
	default:
		// fmt.Println("sdfsdfsdfsf")
		// fmt.Println("%v", c)
		// This should only be called for the 18 colors with a well-defined
		// hue and lightness, not for white, black, or any other color.
		panic(c)
	}
}

func GetColorBlock(img image.Image, currentPos image.Point, currentColor color.Color) (block ColorBlock) {
	block = map[image.Point]struct{} {
		currentPos: {},
	}

	done := false
	for !done {
		done = true
		for pos, _ := range block {
			for _, newPos := range adj(pos) {
				if newPos.In(img.Bounds()) {
					_, inBlock := block[newPos]
					if !inBlock && SameColors(img.At(newPos.X, newPos.Y), currentColor) {
						block[newPos] = struct{}{}
						done = false
					}
				}
			}
		}
	}
	return block
}

func CanMove(img image.Image, currentPos image.Point, dp DpDir) bool {
	newPos := currentPos.Add(image.Point(dp))
	return newPos.In(img.Bounds()) && !SameColors(color.Black, img.At(newPos.X, newPos.Y))
}
