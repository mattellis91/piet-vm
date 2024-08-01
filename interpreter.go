package pietvm

import (
	"bufio"
	_ "fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"unicode"
)

type Interpreter struct {
	Img image.Image
	Stack 
	io.Writer
	io.Reader
	dp DpDir
	cc CCDir
	pos image.Point
	Logger *log.Logger
}

func (in *Interpreter) Run() {
	for in.move() {
	}
}

func (in *Interpreter) move() bool {
	
	//move to edge of current block
	in.moveWithinBlock()

	if in.canMove() {
		newPos := in.pos.Add(image.Point(in.dp))
		oldColor := in.color()
		blockSize := len(in.getColorBlock())

		in.pos = newPos
		in.colorChange(oldColor, blockSize)
		// fmt.Println(oldColor)
		// fmt.Println("moved")
		return true
	}

	return in.recovery()
}

func New(img image.Image) Interpreter {
	return Interpreter{
		Img:    img,
		Writer: os.Stdout,
		Reader: os.Stdin,
		dp:     east,
		cc:     left,
		pos:    img.Bounds().Min,
		Logger: log.New(ioutil.Discard, "", log.Lshortfile),
	}
}

func (i *Interpreter) recovery() bool {
	i.Logger.Println("entering recovery.")
	originalDp := i.dp
	originalCc := i.cc

	// When true, toggle the CC. When false, rotate the DP.
	cc := true

	for !i.canMove() {
		if cc {
			i.toggleCc()
		} else {
			i.rotateDp()
		}
		if i.dp == originalDp && i.cc == originalCc {
			i.Logger.Println("Failed recovery")
			return false
		}

		i.moveWithinBlock()
		i.Logger.Println(i)
		cc = !cc
	}

	i.Logger.Println("recovered.")
	return true
}

func (i *Interpreter) toggleCc() {
	if i.cc == left {
		i.cc = right
	} else {
		i.cc = left
	}
}

func (i *Interpreter) rotateDp() {
	switch i.dp {
	case east:
		i.dp = south
	case south:
		i.dp = west
	case west:
		i.dp = north
	case north:
		i.dp = east
	}
}

func (i *Interpreter) pointer() {
	count := i.pop()
	if count < 0 {
		panic("negative count not implemented")
	}
	for j := 0; j < count; j++ {
		i.rotateDp()
	}
}

func (i *Interpreter) switchCc() {
	count := i.pop()
	if count < 0 {
		count = -count
	}
	if count%2 == 1 {
		i.toggleCc()
	}
}

func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var b byte
	for advance, b = range data {
		if unicode.IsDigit(rune(b)) {
			token = append(token, b)
		} else {
			// Reached a non-digit. Just return what we have.
			return
		}
	}
	if atEOF {
		return
	}

	// 'data' was entirely digits, and we are not at EOF, so signal
	// the Scanner to keep going.
	return 0, nil, nil
}

func (i *Interpreter) inNum() {
	s := bufio.NewScanner(i)
	s.Split(splitFunc)
	if s.Scan() {
		n, err := strconv.Atoi(string(s.Bytes()))
		if err != nil {
			i.Logger.Println(err)
		}
		i.push(n)
		return
	}
	if err := s.Err(); err != nil {
		i.Logger.Println(err)
	}
}

func (in *Interpreter) color() color.Color {
	return ColorAtPos(in.Img, in.pos)
}

func (in *Interpreter) canMove() bool {
	return CanMove(in.Img, in.pos, in.dp)
}

func (in *Interpreter) getColorBlock() ColorBlock {
	return GetColorBlock(in.Img, in.pos, in.color())
}

func (i *Interpreter) inChar() {
	buf := make([]byte, 1)
	_, err := i.Read(buf)
	if err != nil {
		i.Logger.Println(err)
	}
	i.push(int(buf[0]))
}

func (i *Interpreter) colorChange(prevColor color.Color, blockSize int) {
	if SameColors(color.White, prevColor) || SameColors(color.White, i.color()) {
		i.Logger.Println(i, "Moving to/from white: no command to execute")
		return
	}

	oldHue, oldLightness := ColorInfo(prevColor)
	newHue, newLightness := ColorInfo(i.color())

	hueChange := (newHue - oldHue + 6) % 6
	lightnessChange := (newLightness - oldLightness + 3) % 3

	// i.Logger.Println(i, "ΔH:", hueChange, "ΔL:", lightnessChange)
	switch lightnessChange {
	case 0:
		switch hueChange {
		case 1:
			i.add()
		case 2:
			i.divide()
		case 3:
			i.greater()
		case 4:
			i.duplicate()
		case 5:
			i.inChar()
		}
	case 1:
		switch hueChange {
		case 0:
			i.push(blockSize)
		case 1:
			i.subtract()
		case 2:
			i.mod()
		case 3:
			i.pointer()
		case 4:
			i.roll()
		case 5:
			i.outNum()
		}
	case 2:
		switch hueChange {
		case 0:
			i.pop()
		case 1:
			i.multiply()
		case 2:
			i.not()
		case 3:
			i.switchCc()
		case 4:
			i.inNum()
		case 5:
			i.outChar()
		}
	}
	//i.Logger.Println("  stack:", i.data)
}


func (in *Interpreter) moveWithinBlock() {

	//check within white
	if SameColors(color.White, in.color()) {
		newPos := in.pos.Add(image.Point(in.dp))
		for SameColors(color.White, ColorAtPos(in.Img, newPos)) {
			in.pos = newPos
			newPos = in.pos.Add(image.Point(in.dp))
		}
		return
	}

	var newPos *image.Point
	block := in.getColorBlock()
	bounds := block.Bounds()

	switch in.dp {
		case east:
			for p, _ := range block {
				if p.X == bounds.Max.X-1 {
					if newPos == nil || 
						in.cc == left && p.Y < newPos.Y ||
						in.cc == right && p.Y > newPos.Y {
							newPos = &image.Point{p.X, p.Y}
					}
				}
			}
		case south:
			for p, _ := range block {
				if p.Y == bounds.Max.Y-1 {
					if newPos == nil ||
						in.cc == left && p.X > newPos.X ||
						in.cc == right && p.X < newPos.X {
							newPos = &image.Point{p.X, p.Y}
					}
				}
			}
		case west:
			for p, _ := range block {
				if p.X == bounds.Min.X {
					if newPos == nil ||
						in.cc == left && p.Y > newPos.Y ||
						in.cc == right && p.Y < newPos.Y {
							newPos = &image.Point{p.X, p.Y}
					}
				}
			}
		case north:
			for p, _ := range block {
				if p.Y == bounds.Min.Y {
					if newPos == nil ||
						in.cc == left && p.X < newPos.X ||
						in.cc == right && p.X > newPos.X {
							newPos = &image.Point{p.X, p.Y}
					}
				}
			}
	}

	in.pos = *newPos
}

func (i *Interpreter) outNum() {
	io.WriteString(i, strconv.Itoa(i.pop()))
}

func (i *Interpreter) outChar() {
	i.Write([]byte{byte(i.pop())})
}