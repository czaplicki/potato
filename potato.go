package potato

type Potato struct {
    X, Y, I      int
    Memory       [][]byte
    Lables       map[byte] int
    Instructions string
    LastExecuted []int
    IO           chan byte
}

func New(memoryWidth, memoryHeight int, instructions string, io chan byte) *Potato {
    memory := make([][]byte, memoryWidth)
    for i  := range memory { memory[i] = make([]byte, memoryHeight) }
    return &Potato{
        0, 0, -1,
        memory,
        make(map[byte] int, 25),
        instructions,
        make([]int, 8),
        io,
    }
}

func (p Potato)Peek() (byte, bool) {
   i := p.I + 1
   if i >= len(p.Instructions) {
       return 0, false
    }
    return p.Instructions[i], true
}

func (p *Potato)Read() (i byte, b bool) {
    if i, b = p.Peek(); b { p.I++ }
    return
}

func (p *Potato)Tick() bool {
    i, c := p.Read()
    if !c { return false}
    p.LastExecuted = p.LastExecuted[:0]
    p.Execute(i)
    return true
}

func (p *Potato)Evaluate() {
    for p.Tick() {}
}
func (p *Potato)Execute(i byte) {
    p.LastExecuted = append(p.LastExecuted, p.I)
    switch i {
            case '^': p.Y = max(0                   , p.Y - 1)
            case 'v': p.Y = min(len(p.Memory[0]) - 1, p.Y + 1)
            case '<': p.X = max(0                   , p.X - 1)
            case '>': p.X = min(len(p.Memory   ) - 1, p.X + 1)
            case '+': p.Memory[p.X][p.Y] += 1
            case '-': p.Memory[p.X][p.Y] -= 1
            case '_': p.Memory[p.X][p.Y] = <-p.IO
            case '~': p.IO <- p.Memory[p.X][p.Y]
            case ' ', '\t', '\n': p.Tick()
            case '?', '!':
                    j, c := p.Read()
                    if ! c { return }
                    m := p.Memory[p.X][p.Y];
                    if (m == 0) != (i == '!') { return }
                    p.Execute(j)
            default : switch {
                case isUpperLetter(i): p.Lables[i + 32] = p.I - 1
                case isLowerLetter(i): p.I = p.Lables[i]
                case isNumber(i):
                    j, c := p.Read()
                    r := p.I
                    if ! c { return }
                    for n := i - '0'; n > 0; n-- {
                        p.I = r
                        p.Execute(j)
                    }
                }
    }
}

// func (rt *Executetime)display() {
//     fmt.Print("\033[0;0H")
//     for y := 0; y < memSize; y++ {
//         for x := 0; x < memSize; x++ {
//             if x == rt.x && y == rt.y {
//                 fmt.Printf("\033[0;34m%2X\033[0m ", rt.mem[x][y])
//             } else {
//                 fmt.Printf("%2X ", rt.mem[x][y])
//             }
//         }
//         fmt.Printf("\n")
//     }
//     fmt.Printf("\n")
//     for i, c := range rt.ins {
//         if i >= rt.l && i <= rt.i - 1 || i == rt.l {
//             fmt.Printf("\033[0;34m%c\033[0m", c)
//         } else {
//             fmt.Printf("%c", c)
//         }
//     }
//     fmt.Printf("\n")
//     rt.l = rt.i
// }

func isNumber     (c byte) bool { return c >= '0' && c <= '9'}
func isUpperLetter(c byte) bool { return c >= 'A' && c <= 'Z' && c != 'V'}
func isLowerLetter(c byte) bool { return c >= 'a' && c <= 'z' && c != 'v'}

func max(x, y int) int {
    if x > y { return x }
    return y
}

func min(x, y int) int {
    if x < y { return x }
    return y
}
