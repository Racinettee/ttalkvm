package util

type IntStack struct {
	dat []int
	crt int
}

func NewIntStack(initial int) IntStack {
	return IntStack {
		dat: make([]int, initial),
		crt: 0,
	}
}

func (stack *IntStack) Push(val int) {
	stack.dat[stack.crt] = val
	stack.crt += 1
}

func (stack *IntStack) Pop() (result int) {
	result = stack.Top()
	stack.crt -= 1
	return
}

func (stack IntStack) Top() int {
	return stack.dat[stack.crt-1]
}
