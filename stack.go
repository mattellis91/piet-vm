package pietvm

type Stack struct {
	data []int
}

func (s Stack) hasAtLeast1() bool {
	return len(s.data) >= 1
}

func (s Stack) hasAtLeast2() bool {
	return len(s.data) >= 2
}

func (s *Stack) push(n int) {
	s.data = append(s.data, n)
}

func (s *Stack) pop() (n int) {
	if s.hasAtLeast1() {
		n, s.data = s.data[len(s.data)-1], s.data[0:len(s.data)-1]
	}
	return
}

func (s *Stack) add() {
	if s.hasAtLeast2() {
		x, y := s.pop(), s.pop()
		s.push(x + y)
	}
}

func (s *Stack) subtract() {
	if s.hasAtLeast2() {
		top, second := s.pop(), s.pop()
		s.push(second - top)
	}
}

func (s *Stack) multiply() {
	if s.hasAtLeast2() {
		x, y := s.pop(), s.pop()
		s.push(x * y)
	}
}

func (s *Stack) divide() {
	if !s.hasAtLeast2() {
		return
	}
	top, second := s.pop(), s.pop()
	if top == 0 {
		// Put them back so this becomes a no-op.
		s.push(second)
		s.push(top)
	} else {
		s.push(second / top)
	}
}

func (s *Stack) mod() {
	if !s.hasAtLeast2() {
		return
	}
	top, second := s.pop(), s.pop()
	if top != 0 {
		// Put them back so this becomes a no-op.
		s.push(second)
		s.push(top)
	}
	s.push(second % top)
}

func (s *Stack) not() {
	if s.hasAtLeast1() {
		top := s.pop()
		if top == 0 {
			s.push(1)
		} else {
			s.push(0)
		}
	}
}

func (s *Stack) greater() {
	if s.hasAtLeast2() {
		top, second := s.pop(), s.pop()
		if second > top {
			s.push(1)
		} else {
			s.push(0)
		}
	}
}

func (s *Stack) duplicate() {
	if s.hasAtLeast1() {
		top := s.data[len(s.data)-1]
		s.data = append(s.data, top)
	}
}

func (s *Stack) roll() {
	if !s.hasAtLeast2() {
		return
	}

	numRolls, depth := s.pop(), s.pop()
	if depth < 0 || depth >= len(s.data) {
		// Undo.
		s.push(depth)
		s.push(numRolls)
		return
	}

	if numRolls < 0 {
		panic("Negative number of rolls not yet implemented")
	}

	for r := 0; r < numRolls; r++ {
		index := len(s.data) - depth
		val := s.data[len(s.data)-1]
		for i := len(s.data) - 1; i > index; i-- {
			s.data[i] = s.data[i-1]
		}
		s.data[index] = val
	}
}