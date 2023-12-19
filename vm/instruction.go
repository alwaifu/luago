package vm

const MAXARG_Bx = 1<<18 - 1       // 262143
const MAXARG_sBx = MAXARG_Bx >> 1 // 131071

/* OpMode */
/* basic instruction format */
const (
	IABC  = iota // [  B:9  ][  C:9  ][ A:8  ][OP:6]
	IABx         // [      Bx:18     ][ A:8  ][OP:6]
	IAsBx        // [     sBx:18     ][ A:8  ][OP:6]
	IAx          // [           Ax:26        ][OP:6]
)

/* OpArgMask */
const (
	OpArgN = iota // argument is not used
	OpArgU        // argument is used
	OpArgR        // argument is a register or a jump offset
	OpArgK        // argument is a constant or register/constant
)

/*
31       22       13       5    0

	+-------+^------+-^-----+-^-----
	|b=9bits |c=9bits |a=8bits|op=6|
	+-------+^------+-^-----+-^-----
	|    bx=18bits    |a=8bits|op=6|
	+-------+^------+-^-----+-^-----
	|   sbx=18bits    |a=8bits|op=6|
	+-------+^------+-^-----+-^-----
	|    ax=26bits            |op=6|
	+-------+^------+-^-----+-^-----

31      23      15       7      0
*/
type Instruction uint32

func (inst Instruction) Opcode() int {
	return int(inst & 0x3F)
}

func (inst Instruction) ABC() (a, b, c int) {
	a = int(inst >> 6 & 0xFF)
	c = int(inst >> 14 & 0x1FF)
	b = int(inst >> 23 & 0x1FF)
	return
}

func (inst Instruction) ABx() (a, bx int) {
	a = int(inst >> 6 & 0xFF)
	bx = int(inst >> 14)
	return
}

func (inst Instruction) AsBx() (a, sbx int) {
	a, bx := inst.ABx()
	return a, bx - MAXARG_sBx
}

func (inst Instruction) Ax() int {
	return int(inst >> 6)
}

func (inst Instruction) OpName() string {
	return opcodes[inst.Opcode()].name
}

func (inst Instruction) OpMode() byte {
	return opcodes[inst.Opcode()].opMode
}

func (inst Instruction) BMode() byte {
	return opcodes[inst.Opcode()].argBMode
}

func (inst Instruction) CMode() byte {
	return opcodes[inst.Opcode()].argCMode
}
