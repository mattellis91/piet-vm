package pietvm

import (
	"image"
	"io"
)

type DpDir image.Point

var (
	right DpDir = DpDir{1, 0}
	down DpDir = DpDir{0, 1}
	left DpDir = DpDir{-1, 0}
	Up DpDir = DpDir{0, -1}
)

type CCDir int

const (
	
)

type Interpreter struct {
	Img image.Image
	Stack 
	io.Writer
	io.Reader
	Dp DpDir
}