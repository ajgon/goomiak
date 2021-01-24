package cpu

// every mnemonic doesn't include M1R tstates (opcode docoding), as it's done by CPU earlier on.
// only additional T-states are included (for example if M1R=5 => c.states += 1)

func (c *CPU) ldRR_(r, r_ byte) func() {
	// 40       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,B
	// 41       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,C
	// 42       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,D
	// 43       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,E
	// 44       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,H
	// 45       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,L
	// 47       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,A
	// 48       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,B
	// 49       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,C
	// 4A       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,D
	// 4B       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,E
	// 4C       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,H
	// 4D       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,L
	// 4F       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,A
	// 50       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,B
	// 51       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,C
	// 52       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,D
	// 53       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,E
	// 54       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,H
	// 55       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,L
	// 57       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,A
	// 58       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,B
	// 59       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,C
	// 5A       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,D
	// 5B       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,E
	// 5C       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,H
	// 5D       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,L
	// 5F       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,A
	// 60       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD H,B
	// 61       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD H,C
	// 62       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD H,D
	// 63       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD H,E
	// 64       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD H,H
	// 65       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD H,L
	// 67       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD H,A
	// 68       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD L,B
	// 69       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD L,C
	// 6A       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD L,D
	// 6B       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD L,E
	// 6C       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD L,H
	// 6D       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD L,L
	// 6F       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD L,A
	// 78       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,B
	// 79       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,C
	// 7A       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,D
	// 7B       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,E
	// 7C       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,H
	// 7D       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,L
	// 7F       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,A
	// DD40     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,B
	// DD41     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,C
	// DD42     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,D
	// DD43     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,E
	// DD44     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,XH
	// DD45     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,XL
	// DD47     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,A
	// DD48     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,B
	// DD49     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,C
	// DD4A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,D
	// DD4B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,E
	// DD4C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,XH
	// DD4D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,XL
	// DD4F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,A
	// DD50     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,B
	// DD51     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,C
	// DD52     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,D
	// DD53     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,E
	// DD54     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,XH
	// DD55     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,XL
	// DD57     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,A
	// DD58     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,B
	// DD59     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,C
	// DD5A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,D
	// DD5B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,E
	// DD5C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,XH
	// DD5D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,XL
	// DD5F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,A
	// DD60     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD XH,B
	// DD61     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD XH,C
	// DD62     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD XH,D
	// DD63     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD XH,E
	// DD64     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD XH,XH
	// DD65     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD XH,XL
	// DD67     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD XH,A
	// DD68     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD XL,B
	// DD69     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD XL,C
	// DD6A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD XL,D
	// DD6B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD XL,E
	// DD6C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD XL,H
	// DD6D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD XL,L
	// DD6F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD XL,A
	// DD78     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,B
	// DD79     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,C
	// DD7A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,D
	// DD7B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,E
	// DD7C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,XH
	// DD7D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,XL
	// DD7F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,A
	// FD40     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,B
	// FD41     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,C
	// FD42     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,D
	// FD43     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,E
	// FD44     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,YH
	// FD45     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,YL
	// FD47     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,A
	// FD48     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,B
	// FD49     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,C
	// FD4A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,D
	// FD4B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,E
	// FD4C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,YH
	// FD4D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,YL
	// FD4F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,A
	// FD50     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,B
	// FD51     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,C
	// FD52     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,D
	// FD53     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,E
	// FD54     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,YH
	// FD55     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,YL
	// FD57     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,A
	// FD58     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,B
	// FD59     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,C
	// FD5A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,D
	// FD5B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,E
	// FD5C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,YH
	// FD5D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,YL
	// FD5F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,A
	// FD60     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD YH,B
	// FD61     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD YH,C
	// FD62     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD YH,D
	// FD63     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD YH,E
	// FD64     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD YH,YH
	// FD65     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD YH,YL
	// FD67     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD YH,A
	// FD68     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD YL,B
	// FD69     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD YL,C
	// FD6A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD YL,D
	// FD6B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD YL,E
	// FD6C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD YL,H
	// FD6D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD YL,L
	// FD6F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD YL,A
	// FD78     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,B
	// FD79     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,C
	// FD7A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,D
	// FD7B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,E
	// FD7C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,YH
	// FD7D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,YL
	// FD7F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,A
	return func() {
		var lhigh bool
		var lvalue *uint16

		switch r {
		case 'A':
			lhigh, lvalue = true, &c.AF
		case 'B', 'C':
			lhigh, lvalue = r == 'B', &c.BC
		case 'D', 'E':
			lhigh, lvalue = r == 'D', &c.DE
		case 'H', 'L':
			lhigh, lvalue = r == 'H', &c.HL
		case 'X', 'x':
			lhigh, lvalue = r == 'X', &c.IX
		case 'Y', 'y':
			lhigh, lvalue = r == 'Y', &c.IY
		default:
			panic("Invalid `r` part of the mnemonic")
		}

		right := c.extractRegister(r_)

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(right) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(right)
		}

		c.PC++
		if r == 'X' || r == 'x' || r == 'Y' || r == 'y' || r_ == 'X' || r_ == 'x' || r_ == 'Y' || r_ == 'y' {
			c.PC++
		}
	}
}

func (c *CPU) ldRN(r byte) func() {
	// 06U2     07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,U8
	// 0EU2     07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,U8
	// 16U2     07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,U8
	// 1EU2     07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,U8
	// 26U2     07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD H,U8
	// 2EU2     07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD L,U8
	// 3EU2     07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,U8
	// DD06U2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 LD B,U8
	// DD0EU2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 LD C,U8
	// DD16U2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 LD D,U8
	// DD1EU2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 LD E,U8
	// DD26U2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 LD XH,U8
	// DD2EU2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 LD XL,U8
	// DD3EU2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 LD A,U8
	// FD06U2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 LD B,U8
	// FD0EU2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 LD C,U8
	// FD16U2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 LD D,U8
	// FD1EU2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 LD E,U8
	// FD26U2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 LD YH,U8
	// FD2EU2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 LD YL,U8
	// FD3EU2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 LD A,U8
	return func() {
		var lhigh bool
		var lvalue *uint16

		switch r {
		case 'A':
			lhigh, lvalue = true, &c.AF
			c.PC += 2
		case 'B', 'C':
			lhigh, lvalue = r == 'B', &c.BC
			c.PC += 2
		case 'D', 'E':
			lhigh, lvalue = r == 'D', &c.DE
			c.PC += 2
		case 'H', 'L':
			lhigh, lvalue = r == 'H', &c.HL
			c.PC += 2
		case 'X', 'x':
			lhigh, lvalue = r == 'X', &c.IX
			c.PC += 3
		case 'Y', 'y':
			lhigh, lvalue = r == 'Y', &c.IY
			c.PC += 3
		default:
			panic("Invalid `r` part of the mnemonic")
		}

		rvalue := c.readByte(c.PC-1, 3)

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(rvalue) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(rvalue)
		}
	}
}

func (c *CPU) ldR_Ss_(r byte, ss string) func() {
	if ss == "HL" {
		// 46       07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD B,(HL)
		// 4E       07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD C,(HL)
		// 56       07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD D,(HL)
		// 5E       07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD E,(HL)
		// 66       07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD H,(HL)
		// 6E       07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD L,(HL)
		// 7E       07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,(HL)
		return func() {
			var lhigh bool
			var lvalue *uint16
			var right uint8

			switch r {
			case 'A':
				lhigh, lvalue = true, &c.AF
			case 'B', 'C':
				lhigh, lvalue = r == 'B', &c.BC
			case 'D', 'E':
				lhigh, lvalue = r == 'D', &c.DE
			case 'H', 'L':
				lhigh, lvalue = r == 'H', &c.HL
			default:
				panic("Invalid `r` part of the mnemonic")
			}

			right = c.readByte(c.HL, 3)

			if lhigh {
				*lvalue = (*lvalue & 0x00ff) | (uint16(right) << 8)
			} else {
				*lvalue = (*lvalue & 0xff00) | uint16(right)
			}

			c.PC++
		}
	}

	// DD46S2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 LD B,(IX+S8)
	// DD4ES2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 LD C,(IX+S8)
	// DD56S2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 LD D,(IX+S8)
	// DD5ES2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 LD E,(IX+S8)
	// DD66S2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 LD H,(IX+S8)
	// DD6ES2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 LD L,(IX+S8)
	// DD7ES2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 LD A,(IX+S8)
	// FD46S2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 LD B,(IY+S8)
	// FD4ES2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 LD C,(IY+S8)
	// FD56S2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 LD D,(IY+S8)
	// FD5ES2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 LD E,(IY+S8)
	// FD66S2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 LD H,(IY+S8)
	// FD6ES2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 LD L,(IY+S8)
	// FD7ES2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 LD A,(IY+S8)
	return func() {
		var lhigh bool
		var lvalue *uint16
		var right uint8

		switch r {
		case 'A':
			lhigh, lvalue = true, &c.AF
		case 'B', 'C':
			lhigh, lvalue = r == 'B', &c.BC
		case 'D', 'E':
			lhigh, lvalue = r == 'D', &c.DE
		case 'H', 'L':
			lhigh, lvalue = r == 'H', &c.HL
		default:
			panic("Invalid `r` part of the mnemonic")
		}

		shift := c.readByte(c.PC+2, 3)
		c.tstates += 5 // NON
		right = c.readByte(c.shiftedAddress(c.extractRegisterPair(ss), shift), 3)

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(right) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(right)
		}

		c.PC += 3
	}
}

func (c *CPU) ld_Ss_R(ss string, r byte) func() {
	if ss == "HL" {
		// 70       07 00 M1R 4 MWR 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD (HL),B
		// 71       07 00 M1R 4 MWR 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD (HL),C
		// 72       07 00 M1R 4 MWR 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD (HL),D
		// 73       07 00 M1R 4 MWR 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD (HL),E
		// 74       07 00 M1R 4 MWR 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD (HL),H
		// 75       07 00 M1R 4 MWR 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD (HL),L
		// 77       07 00 M1R 4 MWR 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD (HL),A
		return func() {
			c.writeByte(c.HL, c.extractRegister(r), 3)

			c.PC++
		}
	}

	// DD70S1   19 00 M1R 4 M1R 4 MRD 3 NON 5 MWR 3 ... 0 ... 0 LD (IX+S8),B
	// DD71S1   19 00 M1R 4 M1R 4 MRD 3 NON 5 MWR 3 ... 0 ... 0 LD (IX+S8),C
	// DD72S1   19 00 M1R 4 M1R 4 MRD 3 NON 5 MWR 3 ... 0 ... 0 LD (IX+S8),D
	// DD73S1   19 00 M1R 4 M1R 4 MRD 3 NON 5 MWR 3 ... 0 ... 0 LD (IX+S8),E
	// DD74S1   19 00 M1R 4 M1R 4 MRD 3 NON 5 MWR 3 ... 0 ... 0 LD (IX+S8),H
	// DD75S1   19 00 M1R 4 M1R 4 MRD 3 NON 5 MWR 3 ... 0 ... 0 LD (IX+S8),L
	// DD77S1   19 00 M1R 4 M1R 4 MRD 3 NON 5 MWR 3 ... 0 ... 0 LD (IX+S8),A
	// FD70S1   19 00 M1R 4 M1R 4 MRD 3 NON 5 MWR 3 ... 0 ... 0 LD (IY+S8),B
	// FD71S1   19 00 M1R 4 M1R 4 MRD 3 NON 5 MWR 3 ... 0 ... 0 LD (IY+S8),C
	// FD72S1   19 00 M1R 4 M1R 4 MRD 3 NON 5 MWR 3 ... 0 ... 0 LD (IY+S8),D
	// FD73S1   19 00 M1R 4 M1R 4 MRD 3 NON 5 MWR 3 ... 0 ... 0 LD (IY+S8),E
	// FD74S1   19 00 M1R 4 M1R 4 MRD 3 NON 5 MWR 3 ... 0 ... 0 LD (IY+S8),H
	// FD75S1   19 00 M1R 4 M1R 4 MRD 3 NON 5 MWR 3 ... 0 ... 0 LD (IY+S8),L
	// FD77S1   19 00 M1R 4 M1R 4 MRD 3 NON 5 MWR 3 ... 0 ... 0 LD (IY+S8),A
	return func() {
		shift := c.readByte(c.PC+2, 3)
		c.tstates += 5 // NON

		c.writeByte(c.shiftedAddress(c.extractRegisterPair(ss), shift), c.extractRegister(r), 3)

		c.PC += 3
	}
}

func (c *CPU) ld_Ss_N(ss string) func() {
	if ss == "HL" {
		// 36U2     10 00 M1R 4 MRD 3 MWR 3 ... 0 ... 0 ... 0 ... 0 LD (HL),U8
		return func() {
			c.writeByte(c.HL, c.readByte(c.PC+1, 3), 3)

			c.PC += 2
		}
	}

	// DD36S1U2 23 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 4 MWR 3 ... 0 LD (IX+S8),U8
	// FD36S1U2 23 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 4 MWR 3 ... 0 LD (IY+S8),U8
	return func() {
		shift := c.readByte(c.PC+2, 3)
		c.tstates += 5 // NON
		c.writeByte(c.shiftedAddress(c.extractRegisterPair(ss), shift), c.readByte(c.PC+3, 4), 3)

		c.PC += 4
	}
}

// 0A       07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,(BC)
// DD0A     11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 LD A,(BC)
// FD0A     11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 LD A,(BC)
func (c *CPU) ldA_Bc_() {
	c.setAcc(c.readByte(c.BC, 3))
	c.WZ = c.BC + 1
	c.PC++
}

// 1A       07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,(DE)
// DD1A     11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 LD A,(DE)
// FD1A     11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 LD A,(DE)
func (c *CPU) ldA_De_() {
	c.setAcc(c.readByte(c.DE, 3))
	c.WZ = c.DE + 1
	c.PC++
}

// 3AL2H2   13 00 M1R 4 MRD 3 MRD 3 MRD 3 ... 0 ... 0 ... 0 LD A,(U16)
// DD3AL2H2 17 00 M1R 4 M1R 4 MRD 3 MRD 3 MRD 3 ... 0 ... 0 LD A,(U16)
// FD3AL2H2 17 00 M1R 4 M1R 4 MRD 3 MRD 3 MRD 3 ... 0 ... 0 LD A,(U16)
func (c *CPU) ldA_Nn_() {
	address := c.readWord(c.PC+1, 3, 3)
	c.setAcc(c.readByte(address, 3))
	c.PC += 3
	c.WZ = address + 1
}

// 02       07 00 M1R 4 MWR 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD (BC),A
// DD02     11 00 M1R 4 M1R 4 MWR 3 ... 0 ... 0 ... 0 ... 0 LD (BC),A
// FD02     11 00 M1R 4 M1R 4 MWR 3 ... 0 ... 0 ... 0 ... 0 LD (BC),A
func (c *CPU) ld_Bc_A() {
	c.writeByte(c.BC, c.getAcc(), 3)
	c.WZ = ((c.BC + 1) & 0x00ff) | ((uint16(c.getAcc()) << 8) & 0xff00)
	c.PC++
}

// 12       07 00 M1R 4 MWR 3 ... 0 ... 0 ... 0 ... 0 ... 0 LD (DE),A
// DD12     11 00 M1R 4 M1R 4 MWR 3 ... 0 ... 0 ... 0 ... 0 LD (DE),A
// FD12     11 00 M1R 4 M1R 4 MWR 3 ... 0 ... 0 ... 0 ... 0 LD (DE),A
func (c *CPU) ld_De_A() {
	c.writeByte(c.DE, c.getAcc(), 3)
	c.WZ = ((c.DE + 1) & 0x00ff) | ((uint16(c.getAcc()) << 8) & 0xff00)
	c.PC++
}

// 32L1H1   13 00 M1R 4 MRD 3 MRD 3 MWR 3 ... 0 ... 0 ... 0 LD (U16),A
// DD32L1H1 17 00 M1R 4 M1R 4 MRD 3 MRD 3 MWR 3 ... 0 ... 0 LD (U16),A
// FD32L1H1 17 00 M1R 4 M1R 4 MRD 3 MRD 3 MWR 3 ... 0 ... 0 LD (U16),A
func (c *CPU) ld_Nn_A() {
	address := c.readWord(c.PC+1, 3, 3)
	c.writeByte(address, c.getAcc(), 3)
	c.PC += 3
	c.WZ = ((address + 1) & 0x00ff) | ((uint16(c.getAcc()) << 8) & 0xff00)
}

// ED57     09 00 M1R 4 M1R 5 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,I
func (c *CPU) ldAI() {
	c.tstates += 1

	c.setAcc(c.I)

	c.setS(c.I > 127)
	c.setZ(c.I == 0)
	c.setH(false)
	c.setPV(c.States.IFF2)
	c.setN(false)

	c.PC += 2
}

// ED5F     09 00 M1R 4 M1R 5 ... 0 ... 0 ... 0 ... 0 ... 0 LD A,R
func (c *CPU) ldAR() {
	c.tstates += 1

	c.setAcc(c.R)

	c.setS(c.R > 127)
	c.setZ(c.R == 0)
	c.setH(false)
	c.setPV(c.States.IFF2)
	c.setN(false)

	c.PC += 2
}

// ED47     09 00 M1R 4 M1R 5 ... 0 ... 0 ... 0 ... 0 ... 0 LD I,A
func (c *CPU) ldIA() {
	c.tstates += 1

	c.I = c.getAcc()

	c.PC += 2
}

// ED4F     09 00 M1R 4 M1R 5 ... 0 ... 0 ... 0 ... 0 ... 0 LD R,A
func (c *CPU) ldRA() {
	c.tstates += 1

	c.R = c.getAcc()

	c.PC += 2
}

func (c *CPU) ldSsNn(ss string) func() {
	switch ss {
	case "BC":
		// 01L2H2   10 00 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 LD BC,U16
		// DD01L2H2 14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 LD BC,U16
		// FD01L2H2 14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 LD BC,U16
		return func() {
			c.BC = c.readWord(c.PC+1, 3, 3)
			c.PC += 3
		}
	case "DE":
		// 11L2H2   10 00 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 LD DE,U16
		// DD11L2H2 14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 LD DE,U16
		// FD11L2H2 14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 LD DE,U16
		return func() {
			c.DE = c.readWord(c.PC+1, 3, 3)
			c.PC += 3
		}
	case "HL":
		// 21L2H2   10 00 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 LD HL,U16
		return func() {
			c.HL = c.readWord(c.PC+1, 3, 3)
			c.PC += 3
		}
	case "SP":
		// 31L2H2   10 00 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 LD SP,U16
		// DD31L2H2 14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 LD SP,U16
		// FD31L2H2 14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 LD SP,U16
		return func() {
			c.SP = c.readWord(c.PC+1, 3, 3)
			c.PC += 3
		}
	case "IX":
		// DD21L2H2 14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 LD IX,U16
		return func() {
			c.IX = c.readWord(c.PC+2, 3, 3)
			c.PC += 4
		}
	case "IY":
		// FD21L2H2 14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 LD IY,U16
		return func() {
			c.IY = c.readWord(c.PC+2, 3, 3)
			c.PC += 4
		}
	}

	panic("Invalid `ss` type")
}

func (c *CPU) ldSs_Nn_(ss string) func() {
	switch ss {
	case "HL":
		// 2AL2H2   16 00 M1R 4 MRD 3 MRD 3 MRD 3 MRD 3 ... 0 ... 0 LD HL,(U16)
		return func() {
			address := c.readWord(c.PC+1, 3, 3)
			c.HL = c.readWord(address, 3, 3)
			c.WZ = address + 1
			c.PC += 3
		}
	case "IX":
		// DD2AL2H2 20 00 M1R 4 M1R 4 MRD 3 MRD 3 MRD 3 MRD 3 ... 0 LD IX,(U16)
		return func() {
			address := c.readWord(c.PC+2, 3, 3)
			c.IX = c.readWord(address, 3, 3)
			c.WZ = address + 1
			c.PC += 4
		}
	case "IY":
		// FD2AL2H2 20 00 M1R 4 M1R 4 MRD 3 MRD 3 MRD 3 MRD 3 ... 0 LD IY,(U16)
		return func() {
			address := c.readWord(c.PC+2, 3, 3)
			c.IY = c.readWord(address, 3, 3)
			c.WZ = address + 1
			c.PC += 4
		}
	}

	panic("Invalid `ss` type")
}

func (c *CPU) ldRr_Nn_(rr string) func() {
	// ED4BL2H2 20 00 M1R 4 M1R 4 MRD 3 MRD 3 MRD 3 MRD 3 ... 0 LD BC,(U16)
	// ED5BL2H2 20 00 M1R 4 M1R 4 MRD 3 MRD 3 MRD 3 MRD 3 ... 0 LD DE,(U16)
	// ED6BL2H2 20 00 M1R 4 M1R 4 MRD 3 MRD 3 MRD 3 MRD 3 ... 0 LD HL,(U16)
	// ED7BL2H2 20 00 M1R 4 M1R 4 MRD 3 MRD 3 MRD 3 MRD 3 ... 0 LD SP,(U16)
	return func() {
		var lvalue *uint16

		switch rr {
		case "BC":
			lvalue = &c.BC
		case "DE":
			lvalue = &c.DE
		case "HL":
			lvalue = &c.HL
		case "SP":
			lvalue = &c.SP
		default:
			panic("Invalid `rr` part of the mnemonic")
		}

		address := c.readWord(c.PC+2, 3, 3)
		*lvalue = c.readWord(address, 3, 3)

		c.PC += 4
		c.WZ = address + 1
	}
}

func (c *CPU) ld_Nn_Ss(ss string) func() {
	// 22L1H1   16 00 M1R 4 MRD 3 MRD 3 MWR 3 MWR 3 ... 0 ... 0 LD (U16),HL
	if ss == "HL" {
		return func() {
			address := c.readWord(c.PC+1, 3, 3)
			c.writeWord(address, c.HL, 3, 3)
			c.PC += 3
			c.WZ = address + 1
		}
	}

	// DD22L1H1 20 00 M1R 4 M1R 4 MRD 3 MRD 3 MWR 3 MWR 3 ... 0 LD (U16),IX
	// FD22L1H1 20 00 M1R 4 M1R 4 MRD 3 MRD 3 MWR 3 MWR 3 ... 0 LD (U16),IY
	return func() {
		address := c.readWord(c.PC+2, 3, 3)
		c.writeWord(address, c.extractRegisterPair(ss), 3, 3)
		c.PC += 4
		c.WZ = address + 1
	}
}

func (c *CPU) ld_Nn_Rr(rr string) func() {
	// ED43L1H1 20 00 M1R 4 M1R 4 MRD 3 MRD 3 MWR 3 MWR 3 ... 0 LD (U16),BC
	// ED53L1H1 20 00 M1R 4 M1R 4 MRD 3 MRD 3 MWR 3 MWR 3 ... 0 LD (U16),DE
	// ED63L1H1 20 00 M1R 4 M1R 4 MRD 3 MRD 3 MWR 3 MWR 3 ... 0 LD (U16),HL
	// ED73L1H1 20 00 M1R 4 M1R 4 MRD 3 MRD 3 MWR 3 MWR 3 ... 0 LD (U16),SP
	return func() {
		address := c.readWord(c.PC+2, 3, 3)
		c.writeWord(address, c.extractRegisterPair(rr), 3, 3)
		c.WZ = address + 1

		c.PC += 4
	}
}

func (c *CPU) ldSpSs(ss string) func() {
	// F9       06 00 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 LD SP,HL
	if ss == "HL" {
		return func() {
			c.tstates += 2
			c.SP = c.HL

			c.PC++
		}
	}

	// DDF9     10 00 M1R 4 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 LD SP,IX
	// FDF9     10 00 M1R 4 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 LD SP,IY
	return func() {
		c.tstates += 2
		c.SP = c.extractRegisterPair(ss)

		c.PC += 2
	}
}

func (c *CPU) pushSs(ss string) func() {
	switch ss {
	case "BC":
		// C5       11 00 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 ... 0 PUSH BC
		// DDC5     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 PUSH BC
		// FDC5     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 PUSH BC
		return func() {
			c.tstates += 1
			c.pushStack(c.BC)
			c.PC++
		}
	case "DE":
		// D5       11 00 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 ... 0 PUSH DE
		// DDD5     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 PUSH DE
		// FDD5     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 PUSH DE
		return func() {
			c.tstates += 1
			c.pushStack(c.DE)
			c.PC++
		}
	case "HL":
		// E5       11 00 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 ... 0 PUSH HL
		return func() {
			c.tstates += 1
			c.pushStack(c.HL)
			c.PC++
		}
	case "AF":
		// F5       11 00 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 ... 0 PUSH AF
		// DDF5     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 PUSH AF
		// FDF5     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 PUSH AF
		return func() {
			c.tstates += 1
			c.pushStack(c.AF)
			c.PC++
		}
	}

	// DDE5     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 PUSH IX
	// FDE5     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 PUSH IY
	return func() {
		c.tstates += 1
		c.pushStack(c.extractRegisterPair(ss))
		c.PC += 2
	}
}

func (c *CPU) popSs(ss string) func() {
	switch ss {
	case "BC":
		// C1       10 00 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 POP BC
		// DDC1     14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 POP BC
		// FDC1     14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 POP BC
		return func() {
			c.BC = c.popStack()
			c.PC++
		}
	case "DE":
		// D1       10 00 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 POP DE
		// DDD1     14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 POP DE
		// FDD1     14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 POP DE
		return func() {
			c.DE = c.popStack()
			c.PC++
		}
	case "HL":
		// E1       10 00 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 POP HL
		return func() {
			c.HL = c.popStack()
			c.PC++
		}
	case "AF":
		// F1       10 00 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 POP AF
		// DDF1     14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 POP AF
		// FDF1     14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 POP AF
		return func() {
			c.AF = c.popStack()
			c.PC++
		}
	case "IX":
		// DDE1     14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 POP IX
		return func() {
			c.IX = c.popStack()
			c.PC += 2
		}
	case "IY":
		// FDE1     14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 POP IY
		return func() {
			c.IY = c.popStack()
			c.PC += 2
		}
	}

	panic("Invalid `ss` type")
}

// EB       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 EX DE,HL
// DDEB     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 EX DE,HL
// FDEB     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 EX DE,HL
func (c *CPU) exDeSs(ss string) func() {
	if ss == "HL" {
		return func() {
			c.DE, c.HL = c.HL, c.DE

			c.PC++
		}
	}

	return func() {
		c.DE, c.HL = c.HL, c.DE

		c.PC += 2
	}
}

// 08       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 EX AF,AF'
// DD08     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 EX AF,AF'
// FD08     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 EX AF,AF'
func (c *CPU) exAfAf_() {
	c.AF, c.AF_ = c.AF_, c.AF
	c.PC++
}

// D9       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 EXX
// DDD9     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 EXX
// FDD9     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 EXX
func (c *CPU) exx() {
	c.BC, c.BC_ = c.BC_, c.BC
	c.DE, c.DE_ = c.DE_, c.DE
	c.HL, c.HL_ = c.HL_, c.HL

	c.PC++
}

func (c *CPU) ex_Sp_Ss(ss string) func() {
	switch ss {
	// E3       19 00 M1R 4 MRD 3 MRD 4 MWR 3 MWR 5 ... 0 ... 0 EX (SP),HL
	case "HL":
		return func() {
			value := c.readWord(c.SP, 3, 4)
			c.writeWord(c.SP, c.HL, 3, 5)
			c.HL = value
			c.WZ = value

			c.PC++
		}
	// DDE3     23 00 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 5 ... 0 EX (SP),IX
	case "IX":
		return func() {
			value := c.readWord(c.SP, 3, 4)
			c.writeWord(c.SP, c.IX, 3, 5)
			c.IX = value
			c.WZ = value

			c.PC += 2
		}
	// FDE3     23 00 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 5 ... 0 EX (SP),IY
	case "IY":
		return func() {
			value := c.readWord(c.SP, 3, 4)
			c.writeWord(c.SP, c.IY, 3, 5)
			c.IY = value
			c.WZ = value

			c.PC += 2
		}
	}

	panic("Invalid `ss` type")
}

// EDA0     16 00 M1R 4 M1R 4 MRD 3 MWR 5 ... 0 ... 0 ... 0 LDI
func (c *CPU) ldi() {
	rvalue := c.readByte(c.HL, 3)
	c.writeByte(c.DE, rvalue, 5)
	c.DE++
	c.HL++
	c.BC--

	result := rvalue + c.getAcc()

	c.setH(false)
	c.setPV(c.BC != 0)
	c.setN(false)
	c.setF3(result&0x08 == 0x08)
	c.setF5(result&0x02 == 0x02) // 0x02 not 0x20 - this is intended

	c.PC += 2
}

// EDB0     21 16 M1R 4 M1R 4 MRD 3 MWR 5 NON 5 ... 0 ... 0 LDIR
func (c *CPU) ldir() {
	c.ldi()

	if c.BC == 0 {
		return
	}

	c.PC -= 2
	c.WZ = c.PC + 1
	c.tstates += 5 // NON
}

// EDA8     16 00 M1R 4 M1R 4 MRD 3 MWR 5 ... 0 ... 0 ... 0 LDD
func (c *CPU) ldd() {
	rvalue := c.readByte(c.HL, 3)
	c.writeByte(c.DE, rvalue, 5)
	c.DE--
	c.HL--
	c.BC--

	result := rvalue + c.getAcc()

	c.setH(false)
	c.setPV(c.BC != 0)
	c.setN(false)
	c.setF3(result&0x08 == 0x08)
	c.setF5(result&0x02 == 0x02) // 0x02 not 0x20 - this is intended

	c.PC += 2
}

// EDB8     21 16 M1R 4 M1R 4 MRD 3 MWR 5 NON 5 ... 0 ... 0 LDDR
func (c *CPU) lddr() {
	c.ldd()

	if c.BC == 0 {
		return
	}

	c.PC -= 2
	c.WZ = c.PC + 1
	c.tstates += 5 // NON
}

// EDA1     16 00 M1R 4 M1R 4 MRD 3 NON 5 ... 0 ... 0 ... 0 CPI
func (c *CPU) cpi() {
	acc := c.getAcc()
	flagC := c.getC()
	c.setC(true)
	c.adcValueToAcc(c.readByte(c.HL, 3) ^ 0xff)
	result := c.getAcc()
	c.HL++
	c.BC--

	c.setAcc(acc)
	c.setC(flagC)
	c.setN(true)
	c.setPV(c.BC != 0)
	c.setH(!c.getH())

	if c.getH() {
		c.setF3((result-1)&0x08 == 0x08)
		c.setF5((result-1)&0x02 == 0x02) // 0x02 not 0x20 - this is intended
	} else {
		c.setF3(result&0x08 == 0x08)
		c.setF5(result&0x02 == 0x02) // 0x02 not 0x20 - this is intended
	}

	c.PC += 2
	c.WZ++
	c.tstates += 5 // NON
}

// EDB1     21 16 M1R 4 M1R 4 MRD 3 NON 5 NON 5 ... 0 ... 0 CPIR
func (c *CPU) cpir() {
	acc := c.getAcc()
	flagC := c.getC()
	c.setC(true)
	c.adcValueToAcc(c.readByte(c.HL, 3) ^ 0xff)
	result := c.getAcc()
	c.HL++
	c.BC--

	c.setAcc(acc)
	c.setC(flagC)
	c.setN(true)
	c.setPV(c.BC != 0)
	c.setH(!c.getH())
	c.WZ = c.PC + 1

	if c.getH() {
		c.setF3((result-1)&0x08 == 0x08)
		c.setF5((result-1)&0x02 == 0x02) // 0x02 not 0x20 - this is intended
	} else {
		c.setF3(result&0x08 == 0x08)
		c.setF5(result&0x02 == 0x02) // 0x02 not 0x20 - this is intended
	}

	c.tstates += 5 // NON

	if c.BC == 0 || result == 0 {
		c.WZ++
		c.PC += 2
		return
	}
	c.tstates += 5 // NON
}

// EDA9     16 00 M1R 4 M1R 4 MRD 3 NON 5 ... 0 ... 0 ... 0 CPD
func (c *CPU) cpd() {
	acc := c.getAcc()
	flagC := c.getC()
	c.setC(true)
	c.adcValueToAcc(c.readByte(c.HL, 3) ^ 0xff)
	result := c.getAcc()

	c.HL--
	c.BC--

	c.setAcc(acc)
	c.setC(flagC)
	c.setN(true)
	c.setPV(c.BC != 0)
	c.setH(!c.getH())

	if c.getH() {
		result--
	}
	c.setF3(result&0x08 == 0x08)
	c.setF5(result&0x02 == 0x02) // 0x02 not 0x20 - this is intended

	c.PC += 2
	c.WZ--
	c.tstates += 5 // NON
}

// EDB9     21 16 M1R 4 M1R 4 MRD 3 NON 5 NON 5 ... 0 ... 0 CPDR
func (c *CPU) cpdr() {
	acc := c.getAcc()
	flagC := c.getC()
	c.setC(true)
	c.adcValueToAcc(c.readByte(c.HL, 3) ^ 0xff)
	result := c.getAcc()
	c.HL--
	c.BC--

	c.setAcc(acc)
	c.setC(flagC)
	c.setN(true)
	c.setPV(c.BC != 0)
	c.setH(!c.getH())
	c.WZ = c.PC + 1

	if c.getH() {
		result--
	}
	c.setF3(result&0x08 == 0x08)
	c.setF5(result&0x02 == 0x02) // 0x02 not 0x20 - this is intended

	c.tstates += 5 // NON

	if c.BC == 0 || result == 0 {
		c.WZ--
		c.PC += 2
		return
	}

	c.tstates += 5 // NON
}

// 80       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,B
// 81       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,C
// 82       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,D
// 83       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,E
// 84       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,H
// 85       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,L
// 87       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,A
// DD80     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,B
// DD81     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,C
// DD82     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,D
// DD83     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,E
// DD84     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,XH
// DD85     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,XL
// DD87     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,A
// FD80     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,B
// FD81     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,C
// FD82     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,D
// FD83     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,E
// FD84     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,YH
// FD85     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,YL
// FD87     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,A
func (c *CPU) addAR(r byte) func() {
	return func() {
		c.setC(false)
		c.adcValueToAcc(c.extractRegister(r))

		if r == 'X' || r == 'x' || r == 'Y' || r == 'y' {
			c.PC++
		}
		c.PC++
	}
}

// C6U2     07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,U8
// DDC6U2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ADD A,U8
// FDC6U2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ADD A,U8
func (c *CPU) addAN() {
	c.setC(false)
	c.adcValueToAcc(c.readByte(c.PC+1, 3))

	c.PC += 2
}

// 86       07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 ADD A,(HL)
func (c *CPU) addA_Ss_(ss string) func() {
	if ss == "HL" {
		return func() {
			c.setC(false)
			c.adcValueToAcc(c.readByte(c.HL, 3))
			c.PC++
		}
	}

	// DD86S2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 ADD A,(IX+S8)
	// FD86S2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 ADD A,(IY+S8)
	return func() {
		c.setC(false)
		shift := c.readByte(c.PC+2, 3)
		c.tstates += 5 // NON
		c.adcValueToAcc(c.readByte(c.shiftedAddress(c.extractRegisterPair(ss), shift), 3))
		c.PC += 3
	}
}

// 88       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,B
// 89       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,C
// 8A       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,D
// 8B       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,E
// 8C       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,H
// 8D       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,L
// 8F       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,A
// DD88     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,B
// DD89     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,C
// DD8A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,D
// DD8B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,E
// DD8C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,XH
// DD8D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,XL
// DD8F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,A
// FD88     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,B
// FD89     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,C
// FD8A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,D
// FD8B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,E
// FD8C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,YH
// FD8D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,YL
// FD8F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,A
func (c *CPU) adcAR(r byte) func() {
	return func() {
		c.adcValueToAcc(c.extractRegister(r))

		if r == 'X' || r == 'x' || r == 'Y' || r == 'y' {
			c.PC++
		}
		c.PC++
	}
}

// CEU2     07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,U8
// DDCEU2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ADC A,U8
// FDCEU2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ADC A,U8
func (c *CPU) adcAN() {
	c.adcValueToAcc(c.readByte(c.PC+1, 3))

	c.PC += 2
}

// 8E       07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 ADC A,(HL)
// DD8ES2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 ADC A,(IX+S8)
// FD8ES2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 ADC A,(IY+S8)
func (c *CPU) adcA_Ss_(ss string) func() {
	if ss == "HL" {
		return func() {
			c.adcValueToAcc(c.readByte(c.HL, 3))
			c.PC++
		}
	}

	return func() {
		shift := c.readByte(c.PC+2, 3)
		c.tstates += 5 // NON
		c.adcValueToAcc(c.readByte(c.shiftedAddress(c.extractRegisterPair(ss), shift), 3))
		c.PC += 3
	}
}

// 90       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,B
// 91       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,C
// 92       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,D
// 93       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,E
// 94       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,H
// 95       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,L
// 97       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,A
// DD90     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,B
// DD91     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,C
// DD92     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,D
// DD93     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,E
// DD94     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,XH
// DD95     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,XL
// DD97     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,A
// FD90     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,B
// FD91     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,C
// FD92     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,D
// FD93     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,E
// FD94     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,YH
// FD95     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,YL
// FD97     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,A
func (c *CPU) subR(r byte) func() {
	return func() {
		c.setC(true)
		c.adcValueToAcc(c.extractRegister(r) ^ 0xff)

		if r == 'X' || r == 'x' || r == 'Y' || r == 'y' {
			c.PC++
		}
		c.PC++

		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())
	}
}

// D6U2     07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,U8
// DDD6U2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 SUB A,U8
// FDD6U2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 SUB A,U8
func (c *CPU) subN() {
	c.setC(true)
	c.adcValueToAcc(c.readByte(c.PC+1, 3) ^ 0xff)

	c.PC += 2
	c.setN(true)
	c.setC(!c.getC())
	c.setH(!c.getH())
}

func (c *CPU) sub_Ss_(ss string) func() {
	// 96       07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 SUB A,(HL)
	if ss == "HL" {
		return func() {
			c.setC(true)
			c.adcValueToAcc(c.readByte(c.HL, 3) ^ 0xff)

			c.PC++
			c.setN(true)
			c.setC(!c.getC())
			c.setH(!c.getH())
		}
	}

	// DD96S2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 SUB A,(IX+S8)
	// FD96S2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 SUB A,(IY+S8)
	return func() {
		c.setC(true)
		shift := c.readByte(c.PC+2, 3)
		c.tstates += 5 // NON
		c.adcValueToAcc(c.readByte(c.shiftedAddress(c.extractRegisterPair(ss), shift), 3) ^ 0xff)

		c.PC += 3
		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())
	}
}

// 98       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,B
// 99       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,C
// 9A       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,D
// 9B       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,E
// 9C       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,H
// 9D       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,L
// 9F       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,A
// DD98     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,B
// DD99     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,C
// DD9A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,D
// DD9B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,E
// DD9C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,XH
// DD9D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,XL
// DD9F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,A
// FD98     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,B
// FD99     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,C
// FD9A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,D
// FD9B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,E
// FD9C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,YH
// FD9D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,YL
// FD9F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,A
func (c *CPU) sbcAR(r byte) func() {
	return func() {
		c.setC(!c.getC())
		c.adcValueToAcc(c.extractRegister(r) ^ 0xff)

		if r == 'X' || r == 'x' || r == 'Y' || r == 'y' {
			c.PC++
		}
		c.PC++

		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())
	}
}

// DEU2     07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,U8
// DDDEU2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 SBC A,U8
// FDDEU2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 SBC A,U8
func (c *CPU) sbcAN() {
	c.setC(!c.getC())
	c.adcValueToAcc(c.readByte(c.PC+1, 3) ^ 0xff)

	c.PC += 2
	c.setN(true)
	c.setC(!c.getC())
	c.setH(!c.getH())
}

// 9E       07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 SBC A,(HL)
// DD9ES2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 SBC A,(IX+S8)
// FD9ES2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 SBC A,(IY+S8)
func (c *CPU) sbcA_Ss_(ss string) func() {
	if ss == "HL" {
		return func() {
			c.setC(!c.getC())
			c.adcValueToAcc(c.readByte(c.HL, 3) ^ 0xff)

			c.PC++
			c.setN(true)
			c.setC(!c.getC())
			c.setH(!c.getH())
		}
	}

	return func() {
		c.setC(!c.getC())
		shift := c.readByte(c.PC+2, 3)
		c.tstates += 5 // NON
		c.adcValueToAcc(c.readByte(c.shiftedAddress(c.extractRegisterPair(ss), shift), 3) ^ 0xff)

		c.PC += 3
		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())
	}
}

// A0       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,B
// A1       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,C
// A2       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,D
// A3       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,E
// A4       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,H
// A5       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,L
// A7       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,A
// DDA0     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,B
// DDA1     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,C
// DDA2     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,D
// DDA3     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,E
// DDA4     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,XH
// DDA5     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,XL
// DDA7     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,A
// FDA0     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,B
// FDA1     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,C
// FDA2     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,D
// FDA3     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,E
// FDA4     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,YH
// FDA5     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,YL
// FDA7     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,A
func (c *CPU) andR(r byte) func() {
	return func() {
		var result uint8
		result = c.getAcc() & c.extractRegister(r)

		if r == 'X' || r == 'x' || r == 'Y' || r == 'y' {
			c.PC++
		}
		c.PC++

		c.setAcc(result)
		c.setS(result > 127)
		c.setZ(result == 0)
		c.setH(true)
		c.setPV(parityTable[result])
		c.setN(false)
		c.setC(false)
		c.setF5(result&0x20 == 0x20)
		c.setF3(result&0x08 == 0x08)
	}
}

// E6U2     07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,U8
// DDE6U2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 AND A,U8
// FDE6U2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 AND A,U8
func (c *CPU) andN() {
	result := c.getAcc() & c.readByte(c.PC+1, 3)

	c.PC += 2
	c.setAcc(result)
	c.setS(result > 127)
	c.setZ(result == 0)
	c.setH(true)
	c.setPV(parityTable[result])
	c.setN(false)
	c.setC(false)
	c.setF5(result&0x20 == 0x20)
	c.setF3(result&0x08 == 0x08)
}

func (c *CPU) and_Ss_(ss string) func() {
	if ss == "HL" {
		// A6       07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 AND A,(HL)
		return func() {
			result := c.getAcc() & c.readByte(c.HL, 3)

			c.PC++
			c.setAcc(result)
			c.setS(result > 127)
			c.setZ(result == 0)
			c.setH(true)
			c.setPV(parityTable[result])
			c.setN(false)
			c.setC(false)
			c.setF5(result&0x20 == 0x20)
			c.setF3(result&0x08 == 0x08)
		}
	}

	// DDA6S2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 AND A,(IX+S8)
	// FDA6S2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 AND A,(IY+S8)
	return func() {
		shift := c.readByte(c.PC+2, 3)
		c.tstates += 5 // NON
		result := c.getAcc() & c.readByte(c.shiftedAddress(c.extractRegisterPair(ss), shift), 3)

		c.PC += 3
		c.setAcc(result)
		c.setS(result > 127)
		c.setZ(result == 0)
		c.setH(true)
		c.setPV(parityTable[result])
		c.setN(false)
		c.setC(false)
		c.setF5(result&0x20 == 0x20)
		c.setF3(result&0x08 == 0x08)
	}
}

// B0       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,B
// B1       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,C
// B2       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,D
// B3       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,E
// B4       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,H
// B5       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,L
// B7       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,A
// DDB0     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,B
// DDB1     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,C
// DDB2     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,D
// DDB3     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,E
// DDB4     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,XH
// DDB5     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,XL
// DDB7     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,A
// FDB0     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,B
// FDB1     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,C
// FDB2     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,D
// FDB3     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,E
// FDB4     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,YH
// FDB5     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,YL
// FDB7     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,A
func (c *CPU) orR(r byte) func() {
	return func() {
		var result uint8
		result = c.getAcc() | c.extractRegister(r)

		if r == 'X' || r == 'x' || r == 'Y' || r == 'y' {
			c.PC++
		}
		c.PC++

		c.setAcc(result)
		c.setS(result > 127)
		c.setZ(result == 0)
		c.setH(false)
		c.setPV(parityTable[result])
		c.setN(false)
		c.setC(false)
		c.setF5(result&0x20 == 0x20)
		c.setF3(result&0x08 == 0x08)
	}
}

// F6U2     07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,U8
// DDF6U2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 OR A,U8
// FDF6U2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 OR A,U8
func (c *CPU) orN() {
	result := c.getAcc() | c.readByte(c.PC+1, 3)

	c.PC += 2
	c.setAcc(result)
	c.setS(result > 127)
	c.setZ(result == 0)
	c.setH(false)
	c.setPV(parityTable[result])
	c.setN(false)
	c.setC(false)
	c.setF5(result&0x20 == 0x20)
	c.setF3(result&0x08 == 0x08)
}

func (c *CPU) or_Ss_(ss string) func() {
	if ss == "HL" {
		// B6       07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 OR A,(HL)
		return func() {
			result := c.getAcc() | c.readByte(c.HL, 3)

			c.PC++
			c.setAcc(result)
			c.setS(result > 127)
			c.setZ(result == 0)
			c.setH(false)
			c.setPV(parityTable[result])
			c.setN(false)
			c.setC(false)
			c.setF5(result&0x20 == 0x20)
			c.setF3(result&0x08 == 0x08)
		}
	}

	// DDB6S2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 OR A,(IX+S8)
	// FDB6S2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 OR A,(IY+S8)
	return func() {
		shift := c.readByte(c.PC+2, 3)
		c.tstates += 5 // NON
		result := c.getAcc() | c.readByte(c.shiftedAddress(c.extractRegisterPair(ss), shift), 3)

		c.PC += 3
		c.setAcc(result)
		c.setS(result > 127)
		c.setZ(result == 0)
		c.setH(false)
		c.setPV(parityTable[result])
		c.setN(false)
		c.setC(false)
		c.setF5(result&0x20 == 0x20)
		c.setF3(result&0x08 == 0x08)
	}
}

// A8       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,B
// A9       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,C
// AA       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,D
// AB       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,E
// AC       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,H
// AD       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,L
// AF       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,A
// DDA8     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,B
// DDA9     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,C
// DDAA     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,D
// DDAB     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,E
// DDAC     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,XH
// DDAD     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,XL
// DDAF     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,A
// FDA8     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,B
// FDA9     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,C
// FDAA     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,D
// FDAB     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,E
// FDAC     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,YH
// FDAD     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,YL
// FDAF     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,A
func (c *CPU) xorR(r byte) func() {
	return func() {
		var result uint8
		result = c.getAcc() ^ c.extractRegister(r)

		if r == 'X' || r == 'x' || r == 'Y' || r == 'y' {
			c.PC++
		}
		c.PC++

		c.setAcc(result)
		c.setS(result > 127)
		c.setZ(result == 0)
		c.setH(false)
		c.setPV(parityTable[result])
		c.setN(false)
		c.setC(false)
		c.setF5(result&0x20 == 0x20)
		c.setF3(result&0x08 == 0x08)
	}
}

// EEU2     07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,U8
// DDEEU2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 XOR A,U8
// FDEEU2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 XOR A,U8
func (c *CPU) xorN() {
	result := c.getAcc() ^ c.readByte(c.PC+1, 3)

	c.PC += 2
	c.setAcc(result)
	c.setS(result > 127)
	c.setZ(result == 0)
	c.setH(false)
	c.setPV(parityTable[result])
	c.setN(false)
	c.setC(false)
	c.setF5(result&0x20 == 0x20)
	c.setF3(result&0x08 == 0x08)
}

func (c *CPU) xor_Ss_(ss string) func() {
	if ss == "HL" {
		// AE       07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 XOR A,(HL)
		return func() {
			result := c.getAcc() ^ c.readByte(c.HL, 3)

			c.PC++
			c.setAcc(result)
			c.setS(result > 127)
			c.setZ(result == 0)
			c.setH(false)
			c.setPV(parityTable[result])
			c.setN(false)
			c.setC(false)
			c.setF5(result&0x20 == 0x20)
			c.setF3(result&0x08 == 0x08)
		}
	}

	// DDAES2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 XOR A,(IX+S8)
	// FDAES2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 XOR A,(IY+S8)
	return func() {
		shift := c.readByte(c.PC+2, 3)
		c.tstates += 5 // NON
		result := c.getAcc() ^ c.readByte(c.shiftedAddress(c.extractRegisterPair(ss), shift), 3)

		c.PC += 3
		c.setAcc(result)
		c.setS(result > 127)
		c.setZ(result == 0)
		c.setH(false)
		c.setPV(parityTable[result])
		c.setN(false)
		c.setC(false)
		c.setF5(result&0x20 == 0x20)
		c.setF3(result&0x08 == 0x08)
	}
}

// B8       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,B
// B9       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,C
// BA       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,D
// BB       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,E
// BC       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,H
// BD       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,L
// BF       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,A
// DDB8     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,B
// DDB9     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,C
// DDBA     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,D
// DDBB     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,E
// DDBC     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,XH
// DDBD     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,XL
// DDBF     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,A
// FDB8     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,B
// FDB9     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,C
// FDBA     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,D
// FDBB     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,E
// FDBC     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,YH
// FDBD     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,YL
// FDBF     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,A
func (c *CPU) cpR(r byte) func() {
	return func() {
		acc := c.getAcc()
		register := c.extractRegister(r)
		c.setC(true)
		c.adcValueToAcc(register ^ 0xff)

		if r == 'X' || r == 'x' || r == 'Y' || r == 'y' {
			c.PC++
		}
		c.PC++

		c.setAcc(acc)
		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())
		c.setF5(register&0x20 == 0x20)
		c.setF3(register&0x08 == 0x08)
	}
}

// FEU2     07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,U8
// DDFEU2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 CP A,U8
// FDFEU2   11 00 M1R 4 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 CP A,U8
func (c *CPU) cpN() {
	acc := c.getAcc()
	operand := c.readByte(c.PC+1, 3)
	c.setC(true)
	c.adcValueToAcc(operand ^ 0xff)

	c.PC += 2
	c.setAcc(acc)
	c.setN(true)
	c.setC(!c.getC())
	c.setH(!c.getH())
	c.setF5(operand&0x20 == 0x20)
	c.setF3(operand&0x08 == 0x08)
}

func (c *CPU) cp_Ss_(ss string) func() {
	if ss == "HL" {
		// BE       07 00 M1R 4 MRD 3 ... 0 ... 0 ... 0 ... 0 ... 0 CP A,(HL)
		return func() {
			acc := c.getAcc()
			operand := c.readByte(c.HL, 3)
			c.setC(true)
			c.adcValueToAcc(operand ^ 0xff)

			c.PC++
			c.setAcc(acc)
			c.setN(true)
			c.setC(!c.getC())
			c.setH(!c.getH())
			c.setF5(operand&0x20 == 0x20)
			c.setF3(operand&0x08 == 0x08)
		}
	}

	// DDBES2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 CP A,(IX+S8)
	// FDBES2   19 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 3 ... 0 ... 0 CP A,(IY+S8)
	return func() {
		acc := c.getAcc()
		shift := c.readByte(c.PC+2, 3)
		c.tstates += 5 // NON
		operand := c.readByte(c.shiftedAddress(c.extractRegisterPair(ss), shift), 3)
		c.setC(true)
		c.adcValueToAcc(operand ^ 0xff)

		c.PC += 3
		c.setAcc(acc)
		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())
		c.setF5(operand&0x20 == 0x20)
		c.setF3(operand&0x08 == 0x08)
	}
}

// 04       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 INC B
// 0C       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 INC C
// 14       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 INC D
// 1C       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 INC E
// 24       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 INC H
// 2C       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 INC L
// 3C       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 INC A
// DD04     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 INC B
// DD0C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 INC C
// DD14     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 INC D
// DD1C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 INC E
// DD24     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 INC XH
// DD2C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 INC XL
// DD3C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 INC A
// FD04     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 INC B
// FD0C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 INC C
// FD14     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 INC D
// FD1C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 INC E
// FD24     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 INC YH
// FD2C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 INC YL
// FD3C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 INC A
func (c *CPU) incR(r byte) func() {
	return func() {
		var lhigh bool
		var lvalue *uint16

		switch r {
		case 'A':
			lhigh, lvalue = true, &c.AF
			c.PC++
		case 'B', 'C':
			lhigh, lvalue = r == 'B', &c.BC
			c.PC++
		case 'D', 'E':
			lhigh, lvalue = r == 'D', &c.DE
			c.PC++
		case 'H', 'L':
			lhigh, lvalue = r == 'H', &c.HL
			c.PC++
		case 'X', 'x':
			lhigh, lvalue = r == 'X', &c.IX
			c.PC += 2
		case 'Y', 'y':
			lhigh, lvalue = r == 'Y', &c.IY
			c.PC += 2
		default:
			panic("Invalid `r` part of the mnemonic")
		}

		rvalue := c.extractRegister(r) + 1

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(rvalue) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(rvalue)
		}

		c.setN(false)
		c.setPV(rvalue == 0x80)
		c.setH(rvalue&0x0f == 0)
		c.setZ(rvalue == 0)
		c.setS(rvalue > 127)
		c.setF5(rvalue&0x20 == 0x20)
		c.setF3(rvalue&0x08 == 0x08)
	}
}

func (c *CPU) inc_Ss_(ss string) func() {
	if ss == "HL" {
		// 34       11 00 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 ... 0 INC (HL)
		return func() {
			result := c.readByte(c.HL, 4) + 1
			c.writeByte(c.HL, result, 3)
			c.PC++

			c.setN(false)
			c.setPV(result == 0x80)
			c.setH(result&0x0f == 0)
			c.setZ(result == 0)
			c.setS(result > 127)
			c.setF5(result&0x20 == 0x20)
			c.setF3(result&0x08 == 0x08)
		}
	}

	// DD34S1   23 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 4 MWR 3 ... 0 INC (IX+S8)
	// FD34S1   23 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 4 MWR 3 ... 0 INC (IY+S8)
	return func() {
		shift := c.readByte(c.PC+2, 3)
		c.tstates += 5 // NON
		addr := c.shiftedAddress(c.extractRegisterPair(ss), shift)
		result := c.readByte(addr, 4) + 1
		c.writeByte(addr, result, 3)
		c.PC += 3

		c.setN(false)
		c.setPV(result == 0x80)
		c.setH(result&0x0f == 0)
		c.setZ(result == 0)
		c.setS(result > 127)
		c.setF5(result&0x20 == 0x20)
		c.setF3(result&0x08 == 0x08)
	}
}

// 05       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 DEC B
// 0D       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 DEC C
// 15       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 DEC D
// 1D       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 DEC E
// 25       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 DEC H
// 2D       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 DEC L
// 3D       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 DEC A
// DD05     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 DEC B
// DD0D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 DEC C
// DD15     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 DEC D
// DD1D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 DEC E
// DD25     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 DEC XH
// DD2D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 DEC XL
// DD3D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 DEC A
// FD05     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 DEC B
// FD0D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 DEC C
// FD15     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 DEC D
// FD1D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 DEC E
// FD25     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 DEC YH
// FD2D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 DEC YL
// FD3D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 DEC A
func (c *CPU) decR(r byte) func() {
	return func() {
		var lhigh bool
		var lvalue *uint16

		switch r {
		case 'A':
			lhigh, lvalue = true, &c.AF
			c.PC++
		case 'B', 'C':
			lhigh, lvalue = r == 'B', &c.BC
			c.PC++
		case 'D', 'E':
			lhigh, lvalue = r == 'D', &c.DE
			c.PC++
		case 'H', 'L':
			lhigh, lvalue = r == 'H', &c.HL
			c.PC++
		case 'X', 'x':
			lhigh, lvalue = r == 'X', &c.IX
			c.PC += 2
		case 'Y', 'y':
			lhigh, lvalue = r == 'Y', &c.IY
			c.PC += 2
		default:
			panic("Invalid `r` part of the mnemonic")
		}

		rvalue := c.extractRegister(r) - 1

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(rvalue) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(rvalue)
		}

		c.setN(true)
		c.setPV(rvalue == 0x7f)
		c.setH(rvalue&0x0f == 0x0f)
		c.setZ(rvalue == 0)
		c.setS(rvalue > 127)
		c.setF5(rvalue&0x20 == 0x20)
		c.setF3(rvalue&0x08 == 0x08)
	}
}

func (c *CPU) dec_Ss_(ss string) func() {
	if ss == "HL" {
		// 35       11 00 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 ... 0 DEC (HL)
		return func() {
			result := c.readByte(c.HL, 4) - 1
			c.writeByte(c.HL, result, 3)
			c.PC++

			c.setN(true)
			c.setPV(result == 0x7f)
			c.setH(result&0x0f == 0x0f)
			c.setZ(result == 0)
			c.setS(result > 127)
			c.setF5(result&0x20 == 0x20)
			c.setF3(result&0x08 == 0x08)
		}
	}

	// DD35S1   23 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 4 MWR 3 ... 0 DEC (IX+S8)
	// FD35S1   23 00 M1R 4 M1R 4 MRD 3 NON 5 MRD 4 MWR 3 ... 0 DEC (IY+S8)
	return func() {
		addr := c.shiftedAddress(c.extractRegisterPair(ss), c.readByte(c.PC+2, 3))
		c.tstates += 5 // NON
		result := c.readByte(addr, 4) - 1
		c.writeByte(addr, result, 3)
		c.PC += 3

		c.setN(true)
		c.setPV(result == 0x7f)
		c.setH(result&0x0f == 0x0f)
		c.setZ(result == 0)
		c.setS(result > 127)
		c.setF5(result&0x20 == 0x20)
		c.setF3(result&0x08 == 0x08)
	}
}

// 27       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 DAA
// DD27     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 DAA
// FD27     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 DAA
func (c *CPU) daa() {
	t := 0
	a := c.getAcc()

	if c.getH() || (a&0x0f) > 9 {
		t++
	}

	if c.getC() || (a > 0x99) {
		t += 2
		c.setC(true)
	}

	if c.getN() && !c.getH() {
		c.setH(false)
	} else {
		if c.getN() && c.getH() {
			c.setH(a&0x0f < 6)
		} else {
			c.setH(a&0x0f > 9)
		}
	}

	switch t {
	case 1:
		if c.getN() {
			a += 0xfa
		} else {
			a += 0x06
		}
	case 2:
		if c.getN() {
			a += 0xa0
		} else {
			a += 0x60
		}
	case 3:
		if c.getN() {
			a += 0x9a
		} else {
			a += 0x66
		}
	}

	c.setS(a > 127)
	c.setZ(a == 0)
	c.setPV(parityTable[a])
	c.setF5(a&0x20 == 0x20)
	c.setF3(a&0x08 == 0x08)

	c.setAcc(a)

	c.PC++
}

// 2F       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 CPL
// DD2F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 CPL
// FD2F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 CPL
func (c *CPU) cpl() {
	a := c.getAcc() ^ 0xff
	c.setAcc(a)

	c.PC++
	c.setH(true)
	c.setN(true)
	c.setF5(a&0x20 == 0x20)
	c.setF3(a&0x08 == 0x08)
}

// ED44     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NEG
// ED4C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NEG
// ED54     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NEG
// ED5C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NEG
// ED64     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NEG
// ED6C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NEG
// ED74     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NEG
// ED7C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NEG
func (c *CPU) neg() {
	a := c.getAcc()
	c.setAcc(0)

	c.setC(true)
	c.adcValueToAcc(a ^ 0xff)

	c.PC += 2
	c.setN(true)
	c.setC(!c.getC())
	c.setH(!c.getH())
	c.setF5(c.getAcc()&0x20 == 0x20)
	c.setF3(c.getAcc()&0x08 == 0x08)
}

// 3F       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 CCF
// DD3F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 CCF
// FD3F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 CCF
func (c *CPU) ccf() {
	c.PC++

	a := c.getAcc()

	c.setH(c.getC())
	c.setN(false)
	c.setC(!c.getC())
	c.setF5(a&0x20 == 0x20)
	c.setF3(a&0x08 == 0x08)
}

// 37       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 SCF
// DD37     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SCF
// FD37     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SCF
func (c *CPU) scf() {
	c.PC++

	a := c.getAcc()

	c.setC(true)
	c.setN(false)
	c.setH(false)
	c.setF5(a&0x20 == 0x20)
	c.setF3(a&0x08 == 0x08)
}

// 00       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// DD00     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED00     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED01     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED02     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED03     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED04     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED05     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED06     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED07     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED08     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED09     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED0A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED0B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED0C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED0D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED0E     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED0F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED10     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED11     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED12     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED13     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED14     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED15     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED16     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED17     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED18     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED19     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED1A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED1B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED1C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED1D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED1E     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED1F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED20     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED21     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED22     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED23     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED24     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED25     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED26     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED27     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED28     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED29     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED2A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED2B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED2C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED2D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED2E     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED2F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED30     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED31     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED32     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED33     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED34     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED35     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED36     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED37     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED38     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED39     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED3A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED3B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED3C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED3D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED3E     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED3F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED77     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED7F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED80     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED81     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED82     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED83     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED84     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED85     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED86     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED87     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED88     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED89     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED8A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED8B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED8C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED8D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED8E     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED8F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED90     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED91     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED92     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED93     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED94     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED95     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED96     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED97     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED98     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED99     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED9A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED9B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED9C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED9D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED9E     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// ED9F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDA4     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDA5     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDA6     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDA7     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDAC     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDAD     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDAE     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDAF     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDB4     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDB5     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDB6     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDB7     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDBC     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDBD     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDBE     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDBF     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDC0     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDC1     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDC2     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDC3     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDC4     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDC5     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDC6     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDC7     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDC8     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDC9     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDCA     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDCB     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDCC     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDCD     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDCE     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDCF     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDD0     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDD1     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDD2     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDD3     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDD4     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDD5     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDD6     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDD7     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDD8     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDD9     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDDA     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDDB     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDDC     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDDD     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDDE     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDDF     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDE0     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDE1     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDE2     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDE3     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDE4     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDE5     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDE6     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDE7     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDE8     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDE9     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDEA     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDEB     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDEC     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDED     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDEE     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDEF     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDF0     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDF1     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDF2     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDF3     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDF4     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDF5     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDF6     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDF7     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDF8     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDF9     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDFA     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDFB     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDFC     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDFD     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDFE     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// EDFF     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
// FD00     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 NOP
func (c *CPU) nop() {
	c.PC++
}

// 76       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 HALT
// DD76     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 HALT
// FD76     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 HALT
func (c *CPU) halt() {
	c.PC++
	c.States.Halt = true
}

// F3       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 DI
// DDF3     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 DI
// FDF3     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 DI
func (c *CPU) di() {
	c.disableInterrupts()

	c.PC++
}

// FB       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 EI
// DDFB     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 EI
// FDFB     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 EI
func (c *CPU) ei() {
	c.enableInterrupts()

	c.PC++
}

// ED46     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 IM0
// ED4E     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 IM0
// ED56     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 IM1
// ED5E     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 IM2
// ED66     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 IM0
// ED6E     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 IM0
// ED76     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 IM1
// ED7E     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 IM2
func (c *CPU) im(mode uint8) func() {
	return func() {
		c.States.IM = mode
		c.PC += 2
	}
}

// 09       11 00 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 ... 0 ADD HL,BC
// 19       11 00 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 ... 0 ADD HL,DE
// 29       11 00 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 ... 0 ADD HL,HL
// 39       11 00 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 ... 0 ADD HL,SP
// DD09     15 00 M1R 4 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 ADD IX,BC
// DD19     15 00 M1R 4 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 ADD IX,DE
// DD29     15 00 M1R 4 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 ADD IX,IX
// DD39     15 00 M1R 4 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 ADD IX,SP
// FD09     15 00 M1R 4 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 ADD IY,BC
// FD19     15 00 M1R 4 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 ADD IY,DE
// FD29     15 00 M1R 4 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 ADD IY,IY
// FD39     15 00 M1R 4 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 ADD IY,SP
func (c *CPU) addSsRr(ss, rr string) func() {
	switch ss {
	case "HL":
		return func() {
			c.WZ = c.HL + 1
			c.addRegisters(&c.HL, c.extractRegisterPair(rr))
			c.PC++
			c.tstates += 7 // NON + NON
		}
	case "IX":
		return func() {
			c.WZ = c.IX + 1
			c.addRegisters(&c.IX, c.extractRegisterPair(rr))
			c.PC += 2
			c.tstates += 7 // NON + NON
		}
	case "IY":
		return func() {
			c.WZ = c.IY + 1
			c.addRegisters(&c.IY, c.extractRegisterPair(rr))
			c.PC += 2
			c.tstates += 7 // NON + NON
		}
	}

	panic("Invalid `ss` type")
}

// ED4A     15 00 M1R 4 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 ADC HL,BC
// ED5A     15 00 M1R 4 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 ADC HL,DE
// ED6A     15 00 M1R 4 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 ADC HL,HL
// ED7A     15 00 M1R 4 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 ADC HL,SP
func (c *CPU) adcHlRr(rr string) func() {
	return func() {
		c.WZ = c.HL + 1
		c.HL = c.adc16bit(c.HL, c.extractRegisterPair(rr))

		c.PC += 2
		c.tstates += 7 // NON + NON
	}
}

// ED42     15 00 M1R 4 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 SBC HL,BC
// ED52     15 00 M1R 4 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 SBC HL,DE
// ED62     15 00 M1R 4 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 SBC HL,HL
// ED72     15 00 M1R 4 M1R 4 NON 4 NON 3 ... 0 ... 0 ... 0 SBC HL,SP
func (c *CPU) sbcHlRr(rr string) func() {
	return func() {
		c.setC(!c.getC())
		c.WZ = c.HL + 1
		c.HL = c.adc16bit(c.HL, c.extractRegisterPair(rr)^0xffff)

		c.PC += 2
		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())

		c.tstates += 7
	}
}

// 03       06 00 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 INC BC
// 13       06 00 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 INC DE
// 23       06 00 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 INC HL
// 33       06 00 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 INC SP
// DD03     10 00 M1R 4 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 INC BC
// DD13     10 00 M1R 4 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 INC DE
// DD23     10 00 M1R 4 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 INC IX
// DD33     10 00 M1R 4 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 INC SP
// FD03     10 00 M1R 4 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 INC BC
// FD13     10 00 M1R 4 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 INC DE
// FD23     10 00 M1R 4 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 INC IY
// FD33     10 00 M1R 4 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 INC SP
func (c *CPU) incSs(ss string) func() {
	switch ss {
	case "BC":
		return func() {
			c.BC++
			c.PC++
			c.tstates += 2
		}
	case "DE":
		return func() {
			c.DE++
			c.PC++
			c.tstates += 2
		}
	case "HL":
		return func() {
			c.HL++
			c.PC++
			c.tstates += 2
		}
	case "SP":
		return func() {
			c.SP++
			c.PC++
			c.tstates += 2
		}
	case "IX":
		return func() {
			c.IX++
			c.PC += 2
			c.tstates += 2
		}
	case "IY":
		return func() {
			c.IY++
			c.PC += 2
			c.tstates += 2
		}
	}

	panic("Invalid `ss` type")
}

// 0B       06 00 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 DEC BC
// 1B       06 00 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 DEC DE
// 2B       06 00 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 DEC HL
// 3B       06 00 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 DEC SP
// DD0B     10 00 M1R 4 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 DEC BC
// DD1B     10 00 M1R 4 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 DEC DE
// DD2B     10 00 M1R 4 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 DEC IX
// DD3B     10 00 M1R 4 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 DEC SP
// FD0B     10 00 M1R 4 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 DEC BC
// FD1B     10 00 M1R 4 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 DEC DE
// FD2B     10 00 M1R 4 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 DEC IY
// FD3B     10 00 M1R 4 M1R 6 ... 0 ... 0 ... 0 ... 0 ... 0 DEC SP
func (c *CPU) decSs(ss string) func() {
	switch ss {
	case "BC":
		return func() {
			c.BC--
			c.PC++
			c.tstates += 2
		}
	case "DE":
		return func() {
			c.DE--
			c.PC++
			c.tstates += 2
		}
	case "HL":
		return func() {
			c.HL--
			c.PC++
			c.tstates += 2
		}
	case "SP":
		return func() {
			c.SP--
			c.PC++
			c.tstates += 2
		}
	case "IX":
		return func() {
			c.IX--
			c.PC += 2
			c.tstates += 2
		}
	case "IY":
		return func() {
			c.IY--
			c.PC += 2
			c.tstates += 2
		}
	}

	panic("Invalid `ss` type")
}

// 07       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 RLCA
// CB00     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RLC B
// CB01     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RLC C
// CB02     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RLC D
// CB03     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RLC E
// CB04     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RLC H
// CB05     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RLC L
// CB07     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RLC A
// DD07     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RLCA
// FD07     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RLCA
func (c *CPU) rlcR(r byte) func() {
	return func() {
		var size, rvalue uint8
		var lhigh bool
		var lvalue *uint16

		switch r {
		case 'A', ' ':
			lhigh, lvalue = true, &c.AF
		case 'B', 'C':
			lhigh, lvalue = r == 'B', &c.BC
		case 'D', 'E':
			lhigh, lvalue = r == 'D', &c.DE
		case 'H', 'L':
			lhigh, lvalue = r == 'H', &c.HL
		default:
			panic("Invalid `r` part of the mnemonic")
		}

		if r != ' ' {
			size = 1
			rvalue = c.extractRegister(r)
		} else {
			rvalue = c.getAcc()
		}

		signed := rvalue&128 == 128
		rvalue = rvalue << 1
		c.PC += uint16(1 + size)

		if signed {
			rvalue = rvalue | 0x01
		}

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(rvalue) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(rvalue)
		}

		c.setC(signed)
		c.setN(false)
		c.setH(false)
		c.setF5(rvalue&0x20 == 0x20)
		c.setF3(rvalue&0x08 == 0x08)
		if r != ' ' {
			c.setPV(parityTable[rvalue])
			c.setZ(rvalue == 0)
			c.setS(rvalue > 127)
		}
	}
}

// 17       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 RLA
// CB10     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RL B
// CB11     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RL C
// CB12     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RL D
// CB13     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RL E
// CB14     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RL H
// CB15     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RL L
// CB17     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RL A
// DD17     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RLA
// FD17     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RLA
func (c *CPU) rlR(r byte) func() {
	return func() {
		var size, rvalue uint8
		var lhigh bool
		var lvalue *uint16

		switch r {
		case 'A', ' ':
			lhigh, lvalue = true, &c.AF
		case 'B', 'C':
			lhigh, lvalue = r == 'B', &c.BC
		case 'D', 'E':
			lhigh, lvalue = r == 'D', &c.DE
		case 'H', 'L':
			lhigh, lvalue = r == 'H', &c.HL
		default:
			panic("Invalid `r` part of the mnemonic")
		}

		if r != ' ' {
			size = 1
			rvalue = c.extractRegister(r)
		} else {
			rvalue = c.getAcc()
		}

		signed := rvalue&128 == 128
		rvalue = rvalue << 1
		c.PC += uint16(1 + size)

		if c.getC() {
			rvalue = rvalue | 0b00000001
		} else {
			rvalue = rvalue & 0b11111110
		}

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(rvalue) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(rvalue)
		}

		c.setC(signed)
		c.setN(false)
		c.setH(false)
		c.setF5(rvalue&0x20 == 0x20)
		c.setF3(rvalue&0x08 == 0x08)
		if r != ' ' {
			c.setPV(parityTable[rvalue])
			c.setZ(rvalue == 0)
			c.setS(rvalue > 127)
		}
	}
}

// 0F       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 RRCA
// CB08     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RRC B
// CB09     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RRC C
// CB0A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RRC D
// CB0B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RRC E
// CB0C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RRC H
// CB0D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RRC L
// CB0F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RRC A
// DD0F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RRCA
// FD0F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RRCA
func (c *CPU) rrcR(r byte) func() {
	return func() {
		var size, rvalue uint8
		var lhigh bool
		var lvalue *uint16

		switch r {
		case 'A', ' ':
			lhigh, lvalue = true, &c.AF
		case 'B', 'C':
			lhigh, lvalue = r == 'B', &c.BC
		case 'D', 'E':
			lhigh, lvalue = r == 'D', &c.DE
		case 'H', 'L':
			lhigh, lvalue = r == 'H', &c.HL
		default:
			panic("Invalid `r` part of the mnemonic")
		}

		if r != ' ' {
			size = 1
			rvalue = c.extractRegister(r)
		} else {
			rvalue = c.getAcc()
		}

		signed := rvalue&1 == 1
		rvalue = rvalue >> 1
		c.PC += uint16(1 + size)

		if signed {
			rvalue = rvalue | 0x80
		}

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(rvalue) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(rvalue)
		}

		c.setC(signed)
		c.setN(false)
		c.setH(false)
		c.setF5(rvalue&0x20 == 0x20)
		c.setF3(rvalue&0x08 == 0x08)
		if r != ' ' {
			c.setPV(parityTable[rvalue])
			c.setZ(rvalue == 0)
			c.setS(rvalue > 127)
		}
	}
}

// 1F       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 RRA
// CB18     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RR B
// CB19     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RR C
// CB1A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RR D
// CB1B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RR E
// CB1C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RR H
// CB1D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RR L
// CB1F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RR A
// DD1F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RRA
// FD1F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RRA
func (c *CPU) rrR(r byte) func() {
	return func() {
		var size, rvalue uint8
		var lhigh bool
		var lvalue *uint16

		switch r {
		case 'A', ' ':
			lhigh, lvalue = true, &c.AF
		case 'B', 'C':
			lhigh, lvalue = r == 'B', &c.BC
		case 'D', 'E':
			lhigh, lvalue = r == 'D', &c.DE
		case 'H', 'L':
			lhigh, lvalue = r == 'H', &c.HL
		default:
			panic("Invalid `r` part of the mnemonic")
		}

		if r != ' ' {
			size = 1
			rvalue = c.extractRegister(r)

		} else {
			rvalue = c.getAcc()
		}

		signed := rvalue&1 == 1
		rvalue = rvalue >> 1
		c.PC += uint16(1 + size)

		if c.getC() {
			rvalue = rvalue | 0b10000000
		} else {
			rvalue = rvalue & 0b01111111
		}

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(rvalue) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(rvalue)
		}

		c.setC(signed)
		c.setN(false)
		c.setH(false)
		c.setF5(rvalue&0x20 == 0x20)
		c.setF3(rvalue&0x08 == 0x08)
		if r != ' ' {
			c.setPV(parityTable[rvalue])
			c.setZ(rvalue == 0)
			c.setS(rvalue > 127)
		}
	}
}

func (c *CPU) rlcSs(ss string) func() {
	if ss == "HL" {
		//CB06     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 RLC (HL)
		return func() {
			rvalue := c.readByte(c.HL, 4)

			signed := rvalue&128 == 128
			rvalue = rvalue << 1
			c.PC += 2

			if signed {
				rvalue = rvalue | 0x01
			}

			c.writeByte(c.HL, rvalue, 3)

			c.setC(signed)
			c.setN(false)
			c.setPV(parityTable[rvalue])
			c.setH(false)
			c.setZ(rvalue == 0)
			c.setS(rvalue > 127)
			c.setF5(rvalue&0x20 == 0x20)
			c.setF3(rvalue&0x08 == 0x08)
		}
	}

	//DDCBS106 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RLC (IX+S8)
	//FDCBS106 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RLC (IY+S8)
	return func() {
		address := c.shiftedAddress(c.extractRegisterPair(ss), c.readByte(c.PC+2, 5))
		rvalue := c.readByte(address, 4)

		signed := rvalue&128 == 128
		rvalue = rvalue << 1
		c.PC += 4

		if signed {
			rvalue = rvalue | 0x01
		}

		c.writeByte(address, rvalue, 3)

		c.setC(signed)
		c.setN(false)
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setZ(rvalue == 0)
		c.setS(rvalue > 127)
		c.setF5(rvalue&0x20 == 0x20)
		c.setF3(rvalue&0x08 == 0x08)
	}
}

func (c *CPU) rlSs(ss string) func() {
	if ss == "HL" {
		// CB16     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 RL (HL)
		return func() {
			rvalue := c.readByte(c.HL, 4)

			signed := rvalue&128 == 128
			rvalue = rvalue << 1
			c.PC += 2

			if c.getC() {
				rvalue = rvalue | 0b00000001
			} else {
				rvalue = rvalue & 0b11111110
			}

			c.writeByte(c.HL, rvalue, 3)

			c.setC(signed)
			c.setN(false)
			c.setPV(parityTable[rvalue])
			c.setH(false)
			c.setZ(rvalue == 0)
			c.setS(rvalue > 127)
			c.setF5(rvalue&0x20 == 0x20)
			c.setF3(rvalue&0x08 == 0x08)
		}
	}

	// DDCBS116 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RL (IX+S8)
	// FDCBS116 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RL (IY+S8)
	return func() {
		address := c.shiftedAddress(c.extractRegisterPair(ss), c.readByte(c.PC+2, 5))
		rvalue := c.readByte(address, 4)

		signed := rvalue&128 == 128
		rvalue = rvalue << 1
		c.PC += 4

		if c.getC() {
			rvalue = rvalue | 0b00000001
		} else {
			rvalue = rvalue & 0b11111110
		}

		c.writeByte(address, rvalue, 3)

		c.setC(signed)
		c.setN(false)
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setZ(rvalue == 0)
		c.setS(rvalue > 127)
		c.setF5(rvalue&0x20 == 0x20)
		c.setF3(rvalue&0x08 == 0x08)
	}
}

func (c *CPU) rrcSs(ss string) func() {
	if ss == "HL" {
		// CB0E     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 RRC (HL)
		return func() {
			rvalue := c.readByte(c.HL, 4)

			signed := rvalue&1 == 1
			rvalue = rvalue >> 1
			c.PC += 2

			if signed {
				rvalue = rvalue | 0x80
			}

			c.writeByte(c.HL, rvalue, 3)

			c.setC(signed)
			c.setN(false)
			c.setPV(parityTable[rvalue])
			c.setH(false)
			c.setZ(rvalue == 0)
			c.setS(rvalue > 127)
			c.setF5(rvalue&0x20 == 0x20)
			c.setF3(rvalue&0x08 == 0x08)
		}
	}

	// DDCBS10E 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RRC (IX+S8)
	// FDCBS10E 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RRC (IY+S8)
	return func() {
		address := c.shiftedAddress(c.extractRegisterPair(ss), c.readByte(c.PC+2, 5))
		rvalue := c.readByte(address, 4)

		signed := rvalue&1 == 1
		rvalue = rvalue >> 1
		c.PC += 4

		if signed {
			rvalue = rvalue | 0x80
		}

		c.writeByte(address, rvalue, 3)

		c.setC(signed)
		c.setN(false)
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setZ(rvalue == 0)
		c.setS(rvalue > 127)
		c.setF5(rvalue&0x20 == 0x20)
		c.setF3(rvalue&0x08 == 0x08)
	}
}

func (c *CPU) rrSs(ss string) func() {
	if ss == "HL" {
		// CB1E     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 RR (HL)
		return func() {
			rvalue := c.readByte(c.HL, 4)

			signed := rvalue&1 == 1
			rvalue = rvalue >> 1
			c.PC += 2

			if c.getC() {
				rvalue = rvalue | 0b10000000
			} else {
				rvalue = rvalue & 0b01111111
			}

			c.writeByte(c.HL, rvalue, 3)

			c.setC(signed)
			c.setN(false)
			c.setPV(parityTable[rvalue])
			c.setH(false)
			c.setZ(rvalue == 0)
			c.setS(rvalue > 127)
			c.setF5(rvalue&0x20 == 0x20)
			c.setF3(rvalue&0x08 == 0x08)
		}
	}

	// DDCBS11E 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RR (IX+S8)
	// FDCBS11E 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RR (IY+S8)
	return func() {
		address := c.shiftedAddress(c.extractRegisterPair(ss), c.readByte(c.PC+2, 5))
		rvalue := c.readByte(address, 4)

		signed := rvalue&1 == 1
		rvalue = rvalue >> 1
		c.PC += 4

		if c.getC() {
			rvalue = rvalue | 0b10000000
		} else {
			rvalue = rvalue & 0b01111111
		}

		c.writeByte(address, rvalue, 3)

		c.setC(signed)
		c.setN(false)
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setZ(rvalue == 0)
		c.setS(rvalue > 127)
		c.setF5(rvalue&0x20 == 0x20)
		c.setF3(rvalue&0x08 == 0x08)
	}
}

// CB20     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SLA B
// CB21     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SLA C
// CB22     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SLA D
// CB23     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SLA E
// CB24     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SLA H
// CB25     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SLA L
// CB27     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SLA A
func (c *CPU) slaR(r byte) func() {
	return func() {
		var rvalue uint8
		var lhigh bool
		var lvalue *uint16

		switch r {
		case 'A':
			lhigh, lvalue = true, &c.AF
		case 'B', 'C':
			lhigh, lvalue = r == 'B', &c.BC
		case 'D', 'E':
			lhigh, lvalue = r == 'D', &c.DE
		case 'H', 'L':
			lhigh, lvalue = r == 'H', &c.HL
		default:
			panic("Invalid `r` part of the mnemonic")
		}

		rvalue = c.extractRegister(r)

		c.setC(rvalue&128 == 128)
		rvalue = rvalue << 1
		c.PC += 2

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(rvalue) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(rvalue)
		}

		c.setS(rvalue > 127)
		c.setZ(rvalue == 0)
		c.setN(false)
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setF5(rvalue&0x20 == 0x20)
		c.setF3(rvalue&0x08 == 0x08)
	}
}

func (c *CPU) slaSs(ss string) func() {
	if ss == "HL" {
		// CB26     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 SLA (HL)
		return func() {
			rvalue := c.readByte(c.HL, 4)

			c.setC(rvalue&128 == 128)
			rvalue = rvalue << 1
			c.PC += 2

			c.writeByte(c.HL, rvalue, 3)

			c.setN(false)
			c.setPV(parityTable[rvalue])
			c.setH(false)
			c.setZ(rvalue == 0)
			c.setS(rvalue > 127)
			c.setF5(rvalue&0x20 == 0x20)
			c.setF3(rvalue&0x08 == 0x08)
		}
	}

	// DDCBS126 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SLA (IX+S8)
	// FDCBS126 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SLA (IY+S8)
	return func() {
		address := c.shiftedAddress(c.extractRegisterPair(ss), c.readByte(c.PC+2, 5))
		rvalue := c.readByte(address, 4)

		c.setC(rvalue&128 == 128)
		rvalue = rvalue << 1
		c.PC += 4

		c.writeByte(address, rvalue, 3)

		c.setN(false)
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setZ(rvalue == 0)
		c.setS(rvalue > 127)
		c.setF5(rvalue&0x20 == 0x20)
		c.setF3(rvalue&0x08 == 0x08)
	}
}

// CB28     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SRA B
// CB29     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SRA C
// CB2A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SRA D
// CB2B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SRA E
// CB2C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SRA H
// CB2D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SRA L
// CB2F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SRA A
func (c *CPU) sraR(r byte) func() {
	return func() {
		var rvalue uint8
		var lhigh bool
		var lvalue *uint16

		switch r {
		case 'A':
			lhigh, lvalue = true, &c.AF
		case 'B', 'C':
			lhigh, lvalue = r == 'B', &c.BC
		case 'D', 'E':
			lhigh, lvalue = r == 'D', &c.DE
		case 'H', 'L':
			lhigh, lvalue = r == 'H', &c.HL
		default:
			panic("Invalid `r` part of the mnemonic")
		}

		rvalue = c.extractRegister(r)

		c.setC(rvalue&1 == 1)
		rvalue = rvalue >> 1
		if rvalue&64 == 64 {
			rvalue = rvalue | 0b10000000
		} else {
			rvalue = rvalue & 0b01111111
		}
		c.PC += 2

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(rvalue) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(rvalue)
		}

		c.setS(rvalue > 127)
		c.setZ(rvalue == 0)
		c.setN(false)
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setF5(rvalue&0x20 == 0x20)
		c.setF3(rvalue&0x08 == 0x08)
	}
}

func (c *CPU) sraSs(ss string) func() {
	if ss == "HL" {
		// CB2E     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 SRA (HL)
		return func() {
			rvalue := c.readByte(c.HL, 4)

			c.setC(rvalue&1 == 1)
			rvalue = rvalue >> 1
			if rvalue&64 == 64 {
				rvalue = rvalue | 0b10000000
			} else {
				rvalue = rvalue & 0b01111111
			}
			c.PC += 2

			c.writeByte(c.HL, rvalue, 3)

			c.setN(false)
			c.setPV(parityTable[rvalue])
			c.setH(false)
			c.setZ(rvalue == 0)
			c.setS(rvalue > 127)
			c.setF5(rvalue&0x20 == 0x20)
			c.setF3(rvalue&0x08 == 0x08)
		}
	}

	// DDCBS12E 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SRA (IX+S8)
	// FDCBS12E 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SRA (IY+S8)
	return func() {
		address := c.shiftedAddress(c.extractRegisterPair(ss), c.readByte(c.PC+2, 5))
		rvalue := c.readByte(address, 4)

		c.setC(rvalue&1 == 1)
		rvalue = rvalue >> 1
		if rvalue&64 == 64 {
			rvalue = rvalue | 0b10000000
		} else {
			rvalue = rvalue & 0b01111111
		}
		c.PC += 4

		c.writeByte(address, rvalue, 3)

		c.setN(false)
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setZ(rvalue == 0)
		c.setS(rvalue > 127)
		c.setF5(rvalue&0x20 == 0x20)
		c.setF3(rvalue&0x08 == 0x08)
	}
}

// CB38     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SRL B
// CB39     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SRL C
// CB3A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SRL D
// CB3B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SRL E
// CB3C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SRL H
// CB3D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SRL L
// CB3F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SRL A
func (c *CPU) srlR(r byte) func() {
	return func() {
		var rvalue uint8
		var lhigh bool
		var lvalue *uint16

		switch r {
		case 'A':
			lhigh, lvalue = true, &c.AF
		case 'B', 'C':
			lhigh, lvalue = r == 'B', &c.BC
		case 'D', 'E':
			lhigh, lvalue = r == 'D', &c.DE
		case 'H', 'L':
			lhigh, lvalue = r == 'H', &c.HL
		default:
			panic("Invalid `r` part of the mnemonic")
		}

		rvalue = c.extractRegister(r)

		c.setC(rvalue&1 == 1)
		rvalue = rvalue >> 1
		c.PC += 2

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(rvalue) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(rvalue)
		}

		c.setS(false)
		c.setZ(rvalue == 0)
		c.setN(false)
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setF5(rvalue&0x20 == 0x20)
		c.setF3(rvalue&0x08 == 0x08)
	}
}

func (c *CPU) srlSs(ss string) func() {
	if ss == "HL" {
		// CB3E     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 SRL (HL)
		return func() {
			rvalue := c.readByte(c.HL, 4)

			c.setC(rvalue&1 == 1)
			rvalue = rvalue >> 1
			c.PC += 2

			c.writeByte(c.HL, rvalue, 3)

			c.setN(false)
			c.setPV(parityTable[rvalue])
			c.setH(false)
			c.setZ(rvalue == 0)
			c.setS(false)
			c.setF5(rvalue&0x20 == 0x20)
			c.setF3(rvalue&0x08 == 0x08)
		}
	}

	// DDCBS13E 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SRL (IX+S8)
	// FDCBS13E 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SRL (IY+S8)
	return func() {
		address := c.shiftedAddress(c.extractRegisterPair(ss), c.readByte(c.PC+2, 5))
		rvalue := c.readByte(address, 4)

		c.setC(rvalue&1 == 1)
		rvalue = rvalue >> 1
		c.PC += 4

		c.writeByte(address, rvalue, 3)

		c.setN(false)
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setZ(rvalue == 0)
		c.setS(false)
		c.setF5(rvalue&0x20 == 0x20)
		c.setF3(rvalue&0x08 == 0x08)
	}
}

// CB30     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SLL B
// CB31     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SLL C
// CB32     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SLL D
// CB33     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SLL E
// CB34     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SLL H
// CB35     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SLL L
// CB37     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SLL A
func (c *CPU) sllR(r byte) func() {
	return func() {
		var rvalue uint8
		var lhigh bool
		var lvalue *uint16

		switch r {
		case 'A':
			lhigh, lvalue = true, &c.AF
		case 'B', 'C':
			lhigh, lvalue = r == 'B', &c.BC
		case 'D', 'E':
			lhigh, lvalue = r == 'D', &c.DE
		case 'H', 'L':
			lhigh, lvalue = r == 'H', &c.HL
		default:
			panic("Invalid `r` part of the mnemonic")
		}

		rvalue = c.extractRegister(r)

		c.setC(rvalue&128 == 128)
		rvalue = (rvalue << 1) + 1
		c.PC += 2

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(rvalue) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(rvalue)
		}

		c.setS(rvalue > 127)
		c.setZ(rvalue == 0)
		c.setN(false)
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setF5(rvalue&0x20 == 0x20)
		c.setF3(rvalue&0x08 == 0x08)
	}
}

func (c *CPU) sllSs(ss string) func() {
	if ss == "HL" {
		// CB36     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 SLL (HL)
		return func() {
			rvalue := c.readByte(c.HL, 4)

			c.setC(rvalue&128 == 128)
			rvalue = (rvalue << 1) + 1
			c.PC += 2

			c.writeByte(c.HL, rvalue, 3)

			c.setN(false)
			c.setPV(parityTable[rvalue])
			c.setH(false)
			c.setZ(rvalue == 0)
			c.setS(rvalue > 127)
			c.setF5(rvalue&0x20 == 0x20)
			c.setF3(rvalue&0x08 == 0x08)
		}
	}

	// DDCBS136 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SLL (IX+S8)
	// FDCBS136 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SLL (IY+S8)
	return func() {
		address := c.shiftedAddress(c.extractRegisterPair(ss), c.readByte(c.PC+2, 5))
		rvalue := c.readByte(address, 4)

		c.setC(rvalue&128 == 128)
		rvalue = (rvalue << 1) + 1
		c.PC += 4

		c.writeByte(address, rvalue, 3)

		c.setN(false)
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setZ(rvalue == 0)
		c.setS(rvalue > 127)
		c.setF5(rvalue&0x20 == 0x20)
		c.setF3(rvalue&0x08 == 0x08)
	}
}

// ED6F     18 00 M1R 4 M1R 4 MRD 3 NON 4 MWR 3 ... 0 ... 0 RLD
func (c *CPU) rld() {
	value := c.readByte(c.HL, 3)
	a := c.getAcc()

	c.setAcc((a & 0xf0) | ((value >> 4) & 0x0f))
	c.tstates += 4 // NON
	value = value << 4
	value = (a & 0x0f) | value

	c.writeByte(c.HL, value, 3)
	a = c.getAcc()
	c.WZ = c.HL + 1

	c.setS(a > 127)
	c.setZ(a == 0)
	c.setH(false)
	c.setPV(parityTable[a])
	c.setN(false)
	c.setF5(a&0x20 == 0x20)
	c.setF3(a&0x08 == 0x08)

	c.PC += 2
}

// ED67     18 00 M1R 4 M1R 4 MRD 3 NON 4 MWR 3 ... 0 ... 0 RRD
func (c *CPU) rrd() {
	value := c.readByte(c.HL, 3)
	a := c.getAcc()

	c.setAcc((a & 0xf0) | (value & 0x0f))
	c.tstates += 4 // NON
	value = value >> 4
	value = (a << 4) | value

	c.writeByte(c.HL, value, 3)
	a = c.getAcc()
	c.WZ = c.HL + 1

	c.setS(a > 127)
	c.setZ(a == 0)
	c.setH(false)
	c.setPV(parityTable[a])
	c.setN(false)
	c.setF5(a&0x20 == 0x20)
	c.setF3(a&0x08 == 0x08)

	c.PC += 2
}

// CB40     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 0,B
// CB41     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 0,C
// CB42     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 0,D
// CB43     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 0,E
// CB44     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 0,H
// CB45     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 0,L
// CB47     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 0,A
// CB48     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 1,B
// CB49     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 1,C
// CB4A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 1,D
// CB4B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 1,E
// CB4C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 1,H
// CB4D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 1,L
// CB4F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 1,A
// CB50     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 2,B
// CB51     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 2,C
// CB52     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 2,D
// CB53     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 2,E
// CB54     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 2,H
// CB55     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 2,L
// CB57     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 2,A
// CB58     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 3,B
// CB59     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 3,C
// CB5A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 3,D
// CB5B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 3,E
// CB5C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 3,H
// CB5D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 3,L
// CB5F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 3,A
// CB60     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 4,B
// CB61     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 4,C
// CB62     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 4,D
// CB63     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 4,E
// CB64     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 4,H
// CB65     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 4,L
// CB67     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 4,A
// CB68     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 5,B
// CB69     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 5,C
// CB6A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 5,D
// CB6B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 5,E
// CB6C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 5,H
// CB6D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 5,L
// CB6F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 5,A
// CB70     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 6,B
// CB71     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 6,C
// CB72     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 6,D
// CB73     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 6,E
// CB74     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 6,H
// CB75     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 6,L
// CB77     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 6,A
// CB78     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 7,B
// CB79     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 7,C
// CB7A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 7,D
// CB7B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 7,E
// CB7C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 7,H
// CB7D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 7,L
// CB7F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 BIT 7,A
func (c *CPU) bitBR(b uint8, r byte) func() {
	return func() {
		rvalue := c.extractRegister(r)
		mask := uint8(1 << b)

		c.PC += 2

		c.setS(b == 7 && rvalue > 127)
		c.setZ(rvalue&mask == 0)
		c.setF5(rvalue&0x20 == 0x20)
		c.setH(true)
		c.setF3(rvalue&0x08 == 0x08)
		c.setPV(rvalue&mask == 0)
		c.setN(false)
	}
}

func (c *CPU) bitBSs(b uint8, ss string) func() {
	if ss == "HL" {
		// CB46     12 00 M1R 4 M1R 4 MRD 4 ... 0 ... 0 ... 0 ... 0 BIT 0,(HL)
		// CB4E     12 00 M1R 4 M1R 4 MRD 4 ... 0 ... 0 ... 0 ... 0 BIT 1,(HL)
		// CB56     12 00 M1R 4 M1R 4 MRD 4 ... 0 ... 0 ... 0 ... 0 BIT 2,(HL)
		// CB5E     12 00 M1R 4 M1R 4 MRD 4 ... 0 ... 0 ... 0 ... 0 BIT 3,(HL)
		// CB66     12 00 M1R 4 M1R 4 MRD 4 ... 0 ... 0 ... 0 ... 0 BIT 4,(HL)
		// CB6E     12 00 M1R 4 M1R 4 MRD 4 ... 0 ... 0 ... 0 ... 0 BIT 5,(HL)
		// CB76     12 00 M1R 4 M1R 4 MRD 4 ... 0 ... 0 ... 0 ... 0 BIT 6,(HL)
		// CB7E     12 00 M1R 4 M1R 4 MRD 4 ... 0 ... 0 ... 0 ... 0 BIT 7,(HL)
		return func() {
			rvalue := c.readByte(c.HL, 4)
			mask := uint8(1 << b)

			c.PC += 2

			c.setS(b == 7 && rvalue > 127)
			c.setZ(rvalue&mask == 0)
			c.setF5(c.WZ&0x2000 == 0x2000)
			c.setH(true)
			c.setF3(c.WZ&0x0800 == 0x0800)
			c.setPV(rvalue&mask == 0)
			c.setN(false)
		}
	}

	// DDCBS240 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 0,(IX+S8)
	// DDCBS241 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 0,(IX+S8)
	// DDCBS242 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 0,(IX+S8)
	// DDCBS243 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 0,(IX+S8)
	// DDCBS244 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 0,(IX+S8)
	// DDCBS245 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 0,(IX+S8)
	// DDCBS246 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 0,(IX+S8)
	// DDCBS247 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 0,(IX+S8)
	// DDCBS248 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 1,(IX+S8)
	// DDCBS249 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 1,(IX+S8)
	// DDCBS24A 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 1,(IX+S8)
	// DDCBS24B 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 1,(IX+S8)
	// DDCBS24C 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 1,(IX+S8)
	// DDCBS24D 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 1,(IX+S8)
	// DDCBS24E 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 1,(IX+S8)
	// DDCBS24F 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 1,(IX+S8)
	// DDCBS250 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 2,(IX+S8)
	// DDCBS251 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 2,(IX+S8)
	// DDCBS252 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 2,(IX+S8)
	// DDCBS253 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 2,(IX+S8)
	// DDCBS254 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 2,(IX+S8)
	// DDCBS255 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 2,(IX+S8)
	// DDCBS256 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 2,(IX+S8)
	// DDCBS257 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 2,(IX+S8)
	// DDCBS258 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 3,(IX+S8)
	// DDCBS259 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 3,(IX+S8)
	// DDCBS25A 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 3,(IX+S8)
	// DDCBS25B 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 3,(IX+S8)
	// DDCBS25C 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 3,(IX+S8)
	// DDCBS25D 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 3,(IX+S8)
	// DDCBS25E 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 3,(IX+S8)
	// DDCBS25F 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 3,(IX+S8)
	// DDCBS260 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 4,(IX+S8)
	// DDCBS261 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 4,(IX+S8)
	// DDCBS262 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 4,(IX+S8)
	// DDCBS263 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 4,(IX+S8)
	// DDCBS264 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 4,(IX+S8)
	// DDCBS265 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 4,(IX+S8)
	// DDCBS266 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 4,(IX+S8)
	// DDCBS267 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 4,(IX+S8)
	// DDCBS268 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 5,(IX+S8)
	// DDCBS269 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 5,(IX+S8)
	// DDCBS26A 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 5,(IX+S8)
	// DDCBS26B 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 5,(IX+S8)
	// DDCBS26C 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 5,(IX+S8)
	// DDCBS26D 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 5,(IX+S8)
	// DDCBS26E 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 5,(IX+S8)
	// DDCBS26F 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 5,(IX+S8)
	// DDCBS270 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 6,(IX+S8)
	// DDCBS271 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 6,(IX+S8)
	// DDCBS272 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 6,(IX+S8)
	// DDCBS273 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 6,(IX+S8)
	// DDCBS274 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 6,(IX+S8)
	// DDCBS275 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 6,(IX+S8)
	// DDCBS276 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 6,(IX+S8)
	// DDCBS277 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 6,(IX+S8)
	// DDCBS278 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 7,(IX+S8)
	// DDCBS279 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 7,(IX+S8)
	// DDCBS27A 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 7,(IX+S8)
	// DDCBS27B 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 7,(IX+S8)
	// DDCBS27C 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 7,(IX+S8)
	// DDCBS27D 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 7,(IX+S8)
	// DDCBS27E 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 7,(IX+S8)
	// DDCBS27F 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 7,(IX+S8)
	// FDCBS240 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 0,(IY+S8)
	// FDCBS241 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 0,(IY+S8)
	// FDCBS242 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 0,(IY+S8)
	// FDCBS243 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 0,(IY+S8)
	// FDCBS244 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 0,(IY+S8)
	// FDCBS245 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 0,(IY+S8)
	// FDCBS246 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 0,(IY+S8)
	// FDCBS247 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 0,(IY+S8)
	// FDCBS248 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 1,(IY+S8)
	// FDCBS249 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 1,(IY+S8)
	// FDCBS24A 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 1,(IY+S8)
	// FDCBS24B 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 1,(IY+S8)
	// FDCBS24C 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 1,(IY+S8)
	// FDCBS24D 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 1,(IY+S8)
	// FDCBS24E 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 1,(IY+S8)
	// FDCBS24F 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 1,(IY+S8)
	// FDCBS250 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 2,(IY+S8)
	// FDCBS251 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 2,(IY+S8)
	// FDCBS252 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 2,(IY+S8)
	// FDCBS253 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 2,(IY+S8)
	// FDCBS254 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 2,(IY+S8)
	// FDCBS255 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 2,(IY+S8)
	// FDCBS256 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 2,(IY+S8)
	// FDCBS257 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 2,(IY+S8)
	// FDCBS258 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 3,(IY+S8)
	// FDCBS259 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 3,(IY+S8)
	// FDCBS25A 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 3,(IY+S8)
	// FDCBS25B 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 3,(IY+S8)
	// FDCBS25C 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 3,(IY+S8)
	// FDCBS25D 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 3,(IY+S8)
	// FDCBS25E 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 3,(IY+S8)
	// FDCBS25F 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 3,(IY+S8)
	// FDCBS260 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 4,(IY+S8)
	// FDCBS261 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 4,(IY+S8)
	// FDCBS262 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 4,(IY+S8)
	// FDCBS263 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 4,(IY+S8)
	// FDCBS264 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 4,(IY+S8)
	// FDCBS265 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 4,(IY+S8)
	// FDCBS266 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 4,(IY+S8)
	// FDCBS267 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 4,(IY+S8)
	// FDCBS268 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 5,(IY+S8)
	// FDCBS269 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 5,(IY+S8)
	// FDCBS26A 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 5,(IY+S8)
	// FDCBS26B 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 5,(IY+S8)
	// FDCBS26C 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 5,(IY+S8)
	// FDCBS26D 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 5,(IY+S8)
	// FDCBS26E 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 5,(IY+S8)
	// FDCBS26F 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 5,(IY+S8)
	// FDCBS270 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 6,(IY+S8)
	// FDCBS271 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 6,(IY+S8)
	// FDCBS272 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 6,(IY+S8)
	// FDCBS273 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 6,(IY+S8)
	// FDCBS274 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 6,(IY+S8)
	// FDCBS275 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 6,(IY+S8)
	// FDCBS276 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 6,(IY+S8)
	// FDCBS277 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 6,(IY+S8)
	// FDCBS278 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 7,(IY+S8)
	// FDCBS279 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 7,(IY+S8)
	// FDCBS27A 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 7,(IY+S8)
	// FDCBS27B 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 7,(IY+S8)
	// FDCBS27C 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 7,(IY+S8)
	// FDCBS27D 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 7,(IY+S8)
	// FDCBS27E 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 7,(IY+S8)
	// FDCBS27F 20 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 ... 0 ... 0 BIT 7,(IY+S8)
	return func() {
		address := c.shiftedAddress(c.extractRegisterPair(ss), c.readByte(c.PC+2, 5))
		rvalue := c.readByte(address, 4)
		mask := uint8(1 << b)

		c.PC += 4

		c.setS(b == 7 && rvalue > 127)
		c.setZ(rvalue&mask == 0)
		c.setF5(address&0x2000 == 0x2000)
		c.setH(true)
		c.setF3(address&0x0800 == 0x0800)
		c.setPV(rvalue&mask == 0)
		c.setN(false)
	}
}

// CBC0     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 0,B
// CBC1     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 0,C
// CBC2     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 0,D
// CBC3     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 0,E
// CBC4     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 0,H
// CBC5     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 0,L
// CBC7     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 0,A
// CBC8     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 1,B
// CBC9     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 1,C
// CBCA     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 1,D
// CBCB     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 1,E
// CBCC     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 1,H
// CBCD     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 1,L
// CBCF     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 1,A
// CBD0     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 2,B
// CBD1     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 2,C
// CBD2     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 2,D
// CBD3     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 2,E
// CBD4     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 2,H
// CBD5     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 2,L
// CBD7     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 2,A
// CBD8     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 3,B
// CBD9     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 3,C
// CBDA     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 3,D
// CBDB     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 3,E
// CBDC     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 3,H
// CBDD     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 3,L
// CBDF     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 3,A
// CBE0     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 4,B
// CBE1     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 4,C
// CBE2     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 4,D
// CBE3     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 4,E
// CBE4     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 4,H
// CBE5     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 4,L
// CBE7     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 4,A
// CBE8     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 5,B
// CBE9     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 5,C
// CBEA     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 5,D
// CBEB     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 5,E
// CBEC     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 5,H
// CBED     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 5,L
// CBEF     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 5,A
// CBF0     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 6,B
// CBF1     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 6,C
// CBF2     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 6,D
// CBF3     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 6,E
// CBF4     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 6,H
// CBF5     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 6,L
// CBF7     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 6,A
// CBF8     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 7,B
// CBF9     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 7,C
// CBFA     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 7,D
// CBFB     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 7,E
// CBFC     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 7,H
// CBFD     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 7,L
// CBFF     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 SET 7,A
func (c *CPU) setBR(b uint8, r byte) func() {
	return func() {
		var lhigh bool
		var lvalue *uint16

		switch r {
		case 'A':
			lhigh, lvalue = true, &c.AF
		case 'B', 'C':
			lhigh, lvalue = r == 'B', &c.BC
		case 'D', 'E':
			lhigh, lvalue = r == 'D', &c.DE
		case 'H', 'L':
			lhigh, lvalue = r == 'H', &c.HL
		default:
			panic("Invalid `r` part of the mnemonic")
		}

		rvalue := c.extractRegister(r) | uint8(1<<b)

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(rvalue) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(rvalue)
		}

		c.PC += 2
	}
}

func (c *CPU) setBSs(b uint8, ss string) func() {
	switch ss {
	case "HL":
		// CBC6     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 SET 0,(HL)
		// CBCE     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 SET 1,(HL)
		// CBD6     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 SET 2,(HL)
		// CBDE     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 SET 3,(HL)
		// CBE6     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 SET 4,(HL)
		// CBEE     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 SET 5,(HL)
		// CBF6     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 SET 6,(HL)
		// CBFE     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 SET 7,(HL)
		return func() {
			c.writeByte(c.HL, c.readByte(c.HL, 4)|uint8(1<<b), 3)
			c.PC += 2
		}
	case "IX":
		// DDCBS2C0 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 0,(IX+S8),B
		// DDCBS2C1 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 0,(IX+S8),C
		// DDCBS2C2 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 0,(IX+S8),D
		// DDCBS2C3 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 0,(IX+S8),E
		// DDCBS2C4 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 0,(IX+S8),H
		// DDCBS2C5 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 0,(IX+S8),L
		// DDCBS2C6 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 0,(IX+S8)
		// DDCBS2C7 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 0,(IX+S8),A
		// DDCBS2C8 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 1,(IX+S8),B
		// DDCBS2C9 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 1,(IX+S8),C
		// DDCBS2CA 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 1,(IX+S8),D
		// DDCBS2CB 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 1,(IX+S8),E
		// DDCBS2CC 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 1,(IX+S8),H
		// DDCBS2CD 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 1,(IX+S8),L
		// DDCBS2CE 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 1,(IX+S8)
		// DDCBS2CF 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 1,(IX+S8),A
		// DDCBS2D0 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 2,(IX+S8),B
		// DDCBS2D1 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 2,(IX+S8),C
		// DDCBS2D2 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 2,(IX+S8),D
		// DDCBS2D3 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 2,(IX+S8),E
		// DDCBS2D4 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 2,(IX+S8),H
		// DDCBS2D5 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 2,(IX+S8),L
		// DDCBS2D6 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 2,(IX+S8)
		// DDCBS2D7 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 2,(IX+S8),A
		// DDCBS2D8 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 3,(IX+S8),B
		// DDCBS2D9 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 3,(IX+S8),C
		// DDCBS2DA 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 3,(IX+S8),D
		// DDCBS2DB 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 3,(IX+S8),E
		// DDCBS2DC 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 3,(IX+S8),H
		// DDCBS2DD 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 3,(IX+S8),L
		// DDCBS2DE 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 3,(IX+S8)
		// DDCBS2DF 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 3,(IX+S8),A
		// DDCBS2E0 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 4,(IX+S8),B
		// DDCBS2E1 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 4,(IX+S8),C
		// DDCBS2E2 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 4,(IX+S8),D
		// DDCBS2E3 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 4,(IX+S8),E
		// DDCBS2E4 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 4,(IX+S8),H
		// DDCBS2E5 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 4,(IX+S8),L
		// DDCBS2E6 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 4,(IX+S8)
		// DDCBS2E7 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 4,(IX+S8),A
		// DDCBS2E8 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 5,(IX+S8),B
		// DDCBS2E9 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 5,(IX+S8),C
		// DDCBS2EA 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 5,(IX+S8),D
		// DDCBS2EB 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 5,(IX+S8),E
		// DDCBS2EC 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 5,(IX+S8),H
		// DDCBS2ED 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 5,(IX+S8),L
		// DDCBS2EE 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 5,(IX+S8)
		// DDCBS2EF 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 5,(IX+S8),A
		// DDCBS2F0 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 6,(IX+S8),B
		// DDCBS2F1 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 6,(IX+S8),C
		// DDCBS2F2 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 6,(IX+S8),D
		// DDCBS2F3 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 6,(IX+S8),E
		// DDCBS2F4 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 6,(IX+S8),H
		// DDCBS2F5 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 6,(IX+S8),L
		// DDCBS2F6 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 6,(IX+S8)
		// DDCBS2F7 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 6,(IX+S8),A
		// DDCBS2F8 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 7,(IX+S8),B
		// DDCBS2F9 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 7,(IX+S8),C
		// DDCBS2FA 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 7,(IX+S8),D
		// DDCBS2FB 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 7,(IX+S8),E
		// DDCBS2FC 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 7,(IX+S8),H
		// DDCBS2FD 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 7,(IX+S8),L
		// DDCBS2FE 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 7,(IX+S8)
		// DDCBS2FF 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 7,(IX+S8),A
		return func() {
			address := c.shiftedAddress(c.IX, c.readByte(c.PC+2, 5))
			c.writeByte(address, c.readByte(address, 4)|uint8(1<<b), 3)
			c.PC += 4
		}
	case "IY":
		// FDCBS2C0 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 0,(IY+S8),B
		// FDCBS2C1 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 0,(IY+S8),C
		// FDCBS2C2 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 0,(IY+S8),D
		// FDCBS2C3 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 0,(IY+S8),E
		// FDCBS2C4 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 0,(IY+S8),H
		// FDCBS2C5 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 0,(IY+S8),L
		// FDCBS2C6 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 0,(IY+S8)
		// FDCBS2C7 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 0,(IY+S8),A
		// FDCBS2C8 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 1,(IY+S8),B
		// FDCBS2C9 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 1,(IY+S8),C
		// FDCBS2CA 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 1,(IY+S8),D
		// FDCBS2CB 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 1,(IY+S8),E
		// FDCBS2CC 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 1,(IY+S8),H
		// FDCBS2CD 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 1,(IY+S8),L
		// FDCBS2CE 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 1,(IY+S8)
		// FDCBS2CF 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 1,(IY+S8),A
		// FDCBS2D0 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 2,(IY+S8),B
		// FDCBS2D1 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 2,(IY+S8),C
		// FDCBS2D2 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 2,(IY+S8),D
		// FDCBS2D3 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 2,(IY+S8),E
		// FDCBS2D4 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 2,(IY+S8),H
		// FDCBS2D5 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 2,(IY+S8),L
		// FDCBS2D6 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 2,(IY+S8)
		// FDCBS2D7 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 2,(IY+S8),A
		// FDCBS2D8 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 3,(IY+S8),B
		// FDCBS2D9 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 3,(IY+S8),C
		// FDCBS2DA 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 3,(IY+S8),D
		// FDCBS2DB 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 3,(IY+S8),E
		// FDCBS2DC 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 3,(IY+S8),H
		// FDCBS2DD 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 3,(IY+S8),L
		// FDCBS2DE 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 3,(IY+S8)
		// FDCBS2DF 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 3,(IY+S8),A
		// FDCBS2E0 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 4,(IY+S8),B
		// FDCBS2E1 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 4,(IY+S8),C
		// FDCBS2E2 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 4,(IY+S8),D
		// FDCBS2E3 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 4,(IY+S8),E
		// FDCBS2E4 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 4,(IY+S8),H
		// FDCBS2E5 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 4,(IY+S8),L
		// FDCBS2E6 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 4,(IY+S8)
		// FDCBS2E7 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 4,(IY+S8),A
		// FDCBS2E8 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 5,(IY+S8),B
		// FDCBS2E9 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 5,(IY+S8),C
		// FDCBS2EA 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 5,(IY+S8),D
		// FDCBS2EB 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 5,(IY+S8),E
		// FDCBS2EC 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 5,(IY+S8),H
		// FDCBS2ED 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 5,(IY+S8),L
		// FDCBS2EE 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 5,(IY+S8)
		// FDCBS2EF 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 5,(IY+S8),A
		// FDCBS2F0 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 6,(IY+S8),B
		// FDCBS2F1 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 6,(IY+S8),C
		// FDCBS2F2 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 6,(IY+S8),D
		// FDCBS2F3 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 6,(IY+S8),E
		// FDCBS2F4 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 6,(IY+S8),H
		// FDCBS2F5 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 6,(IY+S8),L
		// FDCBS2F6 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 6,(IY+S8)
		// FDCBS2F7 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 6,(IY+S8),A
		// FDCBS2F8 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 7,(IY+S8),B
		// FDCBS2F9 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 7,(IY+S8),C
		// FDCBS2FA 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 7,(IY+S8),D
		// FDCBS2FB 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 7,(IY+S8),E
		// FDCBS2FC 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 7,(IY+S8),H
		// FDCBS2FD 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 7,(IY+S8),L
		// FDCBS2FE 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 7,(IY+S8)
		// FDCBS2FF 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 SET 7,(IY+S8),A
		return func() {
			address := c.shiftedAddress(c.IY, c.readByte(c.PC+2, 5))
			c.writeByte(address, c.readByte(address, 4)|uint8(1<<b), 3)
			c.PC += 4
		}
	}

	panic("Invalid `ss` type")
}

// CB80     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 0,B
// CB81     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 0,C
// CB82     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 0,D
// CB83     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 0,E
// CB84     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 0,H
// CB85     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 0,L
// CB87     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 0,A
// CB88     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 1,B
// CB89     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 1,C
// CB8A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 1,D
// CB8B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 1,E
// CB8C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 1,H
// CB8D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 1,L
// CB8F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 1,A
// CB90     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 2,B
// CB91     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 2,C
// CB92     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 2,D
// CB93     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 2,E
// CB94     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 2,H
// CB95     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 2,L
// CB97     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 2,A
// CB98     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 3,B
// CB99     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 3,C
// CB9A     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 3,D
// CB9B     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 3,E
// CB9C     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 3,H
// CB9D     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 3,L
// CB9F     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 3,A
// CBA0     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 4,B
// CBA1     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 4,C
// CBA2     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 4,D
// CBA3     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 4,E
// CBA4     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 4,H
// CBA5     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 4,L
// CBA7     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 4,A
// CBA8     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 5,B
// CBA9     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 5,C
// CBAA     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 5,D
// CBAB     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 5,E
// CBAC     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 5,H
// CBAD     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 5,L
// CBAF     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 5,A
// CBB0     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 6,B
// CBB1     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 6,C
// CBB2     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 6,D
// CBB3     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 6,E
// CBB4     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 6,H
// CBB5     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 6,L
// CBB7     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 6,A
// CBB8     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 7,B
// CBB9     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 7,C
// CBBA     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 7,D
// CBBB     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 7,E
// CBBC     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 7,H
// CBBD     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 7,L
// CBBF     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 RES 7,A
func (c *CPU) resBR(b uint8, r byte) func() {
	return func() {
		var lhigh bool
		var lvalue *uint16

		switch r {
		case 'A':
			lhigh, lvalue = true, &c.AF
		case 'B', 'C':
			lhigh, lvalue = r == 'B', &c.BC
		case 'D', 'E':
			lhigh, lvalue = r == 'D', &c.DE
		case 'H', 'L':
			lhigh, lvalue = r == 'H', &c.HL
		default:
			panic("Invalid `r` part of the mnemonic")
		}

		rvalue := c.extractRegister(r) & (uint8(1<<b) ^ 0xff)

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(rvalue) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(rvalue)
		}

		c.PC += 2
	}
}

func (c *CPU) resBSs(b uint8, ss string) func() {
	switch ss {
	case "HL":
		// CB86     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 RES 0,(HL)
		// CB8E     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 RES 1,(HL)
		// CB96     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 RES 2,(HL)
		// CB9E     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 RES 3,(HL)
		// CBA6     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 RES 4,(HL)
		// CBAE     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 RES 5,(HL)
		// CBB6     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 RES 6,(HL)
		// CBBE     15 00 M1R 4 M1R 4 MRD 4 MWR 3 ... 0 ... 0 ... 0 RES 7,(HL)
		return func() {
			c.writeByte(c.HL, c.readByte(c.HL, 4)&(uint8(1<<b)^0xff), 3)
			c.PC += 2
		}
	case "IX":
		// DDCBS280 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 0,(IX+S8),B
		// DDCBS281 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 0,(IX+S8),C
		// DDCBS282 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 0,(IX+S8),D
		// DDCBS283 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 0,(IX+S8),E
		// DDCBS284 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 0,(IX+S8),H
		// DDCBS285 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 0,(IX+S8),L
		// DDCBS286 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 0,(IX+S8)
		// DDCBS287 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 0,(IX+S8),A
		// DDCBS288 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 1,(IX+S8),B
		// DDCBS289 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 1,(IX+S8),C
		// DDCBS28A 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 1,(IX+S8),D
		// DDCBS28B 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 1,(IX+S8),E
		// DDCBS28C 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 1,(IX+S8),H
		// DDCBS28D 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 1,(IX+S8),L
		// DDCBS28E 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 1,(IX+S8)
		// DDCBS28F 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 1,(IX+S8),A
		// DDCBS290 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 2,(IX+S8),B
		// DDCBS291 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 2,(IX+S8),C
		// DDCBS292 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 2,(IX+S8),D
		// DDCBS293 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 2,(IX+S8),E
		// DDCBS294 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 2,(IX+S8),H
		// DDCBS295 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 2,(IX+S8),L
		// DDCBS296 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 2,(IX+S8)
		// DDCBS297 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 2,(IX+S8),A
		// DDCBS298 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 3,(IX+S8),B
		// DDCBS299 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 3,(IX+S8),C
		// DDCBS29A 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 3,(IX+S8),D
		// DDCBS29B 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 3,(IX+S8),E
		// DDCBS29C 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 3,(IX+S8),H
		// DDCBS29D 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 3,(IX+S8),L
		// DDCBS29E 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 3,(IX+S8)
		// DDCBS29F 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 3,(IX+S8),A
		// DDCBS2A0 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 4,(IX+S8),B
		// DDCBS2A1 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 4,(IX+S8),C
		// DDCBS2A2 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 4,(IX+S8),D
		// DDCBS2A3 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 4,(IX+S8),E
		// DDCBS2A4 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 4,(IX+S8),H
		// DDCBS2A5 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 4,(IX+S8),L
		// DDCBS2A6 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 4,(IX+S8)
		// DDCBS2A7 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 4,(IX+S8),A
		// DDCBS2A8 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 5,(IX+S8),B
		// DDCBS2A9 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 5,(IX+S8),C
		// DDCBS2AA 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 5,(IX+S8),D
		// DDCBS2AB 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 5,(IX+S8),E
		// DDCBS2AC 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 5,(IX+S8),H
		// DDCBS2AD 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 5,(IX+S8),L
		// DDCBS2AE 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 5,(IX+S8)
		// DDCBS2AF 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 5,(IX+S8),A
		// DDCBS2B0 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 6,(IX+S8),B
		// DDCBS2B1 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 6,(IX+S8),C
		// DDCBS2B2 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 6,(IX+S8),D
		// DDCBS2B3 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 6,(IX+S8),E
		// DDCBS2B4 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 6,(IX+S8),H
		// DDCBS2B5 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 6,(IX+S8),L
		// DDCBS2B6 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 6,(IX+S8)
		// DDCBS2B7 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 6,(IX+S8),A
		// DDCBS2B8 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 7,(IX+S8),B
		// DDCBS2B9 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 7,(IX+S8),C
		// DDCBS2BA 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 7,(IX+S8),D
		// DDCBS2BB 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 7,(IX+S8),E
		// DDCBS2BC 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 7,(IX+S8),H
		// DDCBS2BD 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 7,(IX+S8),L
		// DDCBS2BE 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 7,(IX+S8)
		// DDCBS2BF 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 7,(IX+S8),A
		return func() {
			address := c.shiftedAddress(c.IX, c.readByte(c.PC+2, 5))
			c.writeByte(address, c.readByte(address, 4)&(uint8(1<<b)^0xff), 3)
			c.PC += 4
		}
	case "IY":
		// FDCBS280 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 0,(IY+S8),B
		// FDCBS281 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 0,(IY+S8),C
		// FDCBS282 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 0,(IY+S8),D
		// FDCBS283 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 0,(IY+S8),E
		// FDCBS284 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 0,(IY+S8),H
		// FDCBS285 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 0,(IY+S8),L
		// FDCBS286 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 0,(IY+S8)
		// FDCBS287 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 0,(IY+S8),A
		// FDCBS288 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 1,(IY+S8),B
		// FDCBS289 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 1,(IY+S8),C
		// FDCBS28A 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 1,(IY+S8),D
		// FDCBS28B 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 1,(IY+S8),E
		// FDCBS28C 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 1,(IY+S8),H
		// FDCBS28D 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 1,(IY+S8),L
		// FDCBS28E 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 1,(IY+S8)
		// FDCBS28F 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 1,(IY+S8),A
		// FDCBS290 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 2,(IY+S8),B
		// FDCBS291 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 2,(IY+S8),C
		// FDCBS292 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 2,(IY+S8),D
		// FDCBS293 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 2,(IY+S8),E
		// FDCBS294 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 2,(IY+S8),H
		// FDCBS295 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 2,(IY+S8),L
		// FDCBS296 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 2,(IY+S8)
		// FDCBS297 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 2,(IY+S8),A
		// FDCBS298 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 3,(IY+S8),B
		// FDCBS299 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 3,(IY+S8),C
		// FDCBS29A 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 3,(IY+S8),D
		// FDCBS29B 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 3,(IY+S8),E
		// FDCBS29C 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 3,(IY+S8),H
		// FDCBS29D 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 3,(IY+S8),L
		// FDCBS29E 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 3,(IY+S8)
		// FDCBS29F 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 3,(IY+S8),A
		// FDCBS2A0 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 4,(IY+S8),B
		// FDCBS2A1 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 4,(IY+S8),C
		// FDCBS2A2 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 4,(IY+S8),D
		// FDCBS2A3 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 4,(IY+S8),E
		// FDCBS2A4 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 4,(IY+S8),H
		// FDCBS2A5 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 4,(IY+S8),L
		// FDCBS2A6 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 4,(IY+S8)
		// FDCBS2A7 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 4,(IY+S8),A
		// FDCBS2A8 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 5,(IY+S8),B
		// FDCBS2A9 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 5,(IY+S8),C
		// FDCBS2AA 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 5,(IY+S8),D
		// FDCBS2AB 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 5,(IY+S8),E
		// FDCBS2AC 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 5,(IY+S8),H
		// FDCBS2AD 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 5,(IY+S8),L
		// FDCBS2AE 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 5,(IY+S8)
		// FDCBS2AF 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 5,(IY+S8),A
		// FDCBS2B0 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 6,(IY+S8),B
		// FDCBS2B1 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 6,(IY+S8),C
		// FDCBS2B2 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 6,(IY+S8),D
		// FDCBS2B3 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 6,(IY+S8),E
		// FDCBS2B4 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 6,(IY+S8),H
		// FDCBS2B5 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 6,(IY+S8),L
		// FDCBS2B6 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 6,(IY+S8)
		// FDCBS2B7 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 6,(IY+S8),A
		// FDCBS2B8 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 7,(IY+S8),B
		// FDCBS2B9 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 7,(IY+S8),C
		// FDCBS2BA 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 7,(IY+S8),D
		// FDCBS2BB 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 7,(IY+S8),E
		// FDCBS2BC 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 7,(IY+S8),H
		// FDCBS2BD 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 7,(IY+S8),L
		// FDCBS2BE 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 7,(IY+S8)
		// FDCBS2BF 23 00 M1R 4 M1R 4 MRD 3 MRD 5 MRD 4 MWR 3 ... 0 RES 7,(IY+S8),A
		return func() {
			address := c.shiftedAddress(c.IY, c.readByte(c.PC+2, 5))
			c.writeByte(address, c.readByte(address, 4)&(uint8(1<<b)^0xff), 3)
			c.PC += 4
		}
	}

	panic("Invalid `ss` type")
}

// C3L1H1   10 00 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 JP U16
// DDC3L1H1 14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 JP U16
// FDC3L1H1 14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 JP U16
func (c *CPU) jpNn() {
	c.PC = c.readWord(c.PC+1, 3, 3)
	c.WZ = c.PC
}

// C2L2H2   10 10 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 JP NZ,U16
// DDC2L2H2 14 14 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 JP NZ,U16
// FDC2L2H2 14 14 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 JP NZ,U16
func (c *CPU) jpNzNn() {
	c.WZ = c.readWord(c.PC+1, 3, 3)

	if c.getZ() {
		c.PC += 3
		return
	}

	c.PC = c.WZ
}

// CAL2H2   10 10 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 JP Z,U16
// DDCAL2H2 14 14 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 JP Z,U16
// FDCAL2H2 14 14 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 JP Z,U16
func (c *CPU) jpZNn() {
	c.WZ = c.readWord(c.PC+1, 3, 3)

	if !c.getZ() {
		c.PC += 3
		return
	}

	c.PC = c.WZ
}

// D2L2H2   10 10 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 JP NC,U16
// DDD2L2H2 14 14 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 JP NC,U16
// FDD2L2H2 14 14 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 JP NC,U16
func (c *CPU) jpNcNn() {
	c.WZ = c.readWord(c.PC+1, 3, 3)

	if c.getC() {
		c.PC += 3
		return
	}

	c.PC = c.WZ
}

// DAL2H2   10 10 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 JP C,U16
// DDDAL2H2 14 14 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 JP C,U16
// FDDAL2H2 14 14 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 JP C,U16
func (c *CPU) jpCNn() {
	c.WZ = c.readWord(c.PC+1, 3, 3)
	if !c.getC() {
		c.PC += 3
		return
	}

	c.PC = c.WZ
}

// E2L2H2   10 10 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 JP PO,U16
// DDE2L2H2 14 14 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 JP PO,U16
// FDE2L2H2 14 14 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 JP PO,U16
func (c *CPU) jpPoNn() {
	c.WZ = c.readWord(c.PC+1, 3, 3)
	if c.getPV() {
		c.PC += 3
		return
	}

	c.PC = c.WZ
}

// EAL2H2   10 10 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 JP PE,U16
// DDEAL2H2 14 14 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 JP PE,U16
// FDEAL2H2 14 14 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 JP PE,U16
func (c *CPU) jpPeNn() {
	c.WZ = c.readWord(c.PC+1, 3, 3)
	if !c.getPV() {
		c.PC += 3
		return
	}

	c.PC = c.WZ
}

// F2L2H2   10 10 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 JP P,U16
// DDF2L2H2 14 14 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 JP P,U16
// FDF2L2H2 14 14 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 JP P,U16
func (c *CPU) jpPNn() {
	c.WZ = c.readWord(c.PC+1, 3, 3)
	if c.getS() {
		c.PC += 3
		return
	}

	c.PC = c.WZ
}

// FAL2H2   10 10 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 JP M,U16
// DDFAL2H2 14 14 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 JP M,U16
// FDFAL2H2 14 14 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 JP M,U16
func (c *CPU) jpMNn() {
	c.WZ = c.readWord(c.PC+1, 3, 3)
	if !c.getS() {
		c.PC += 3
		return
	}

	c.PC = c.WZ
}

// 18S1     12 00 M1R 4 MRD 3 NON 5 ... 0 ... 0 ... 0 ... 0 JR S8
// DD18S1   16 00 M1R 4 M1R 4 MRD 3 NON 5 ... 0 ... 0 ... 0 JR S8
// FD18S1   16 00 M1R 4 M1R 4 MRD 3 NON 5 ... 0 ... 0 ... 0 JR S8
func (c *CPU) jrN() {
	c.WZ = 2 + uint16(int16(c.PC)+int16(int8(c.readByte(c.PC+1, 3))))
	c.PC = c.WZ
	c.tstates += 5
}

// 38S2     12 07 M1R 4 MRD 3 NON 5 ... 0 ... 0 ... 0 ... 0 JR C,S8
// DD38S2   16 11 M1R 4 M1R 4 MRD 3 NON 5 ... 0 ... 0 ... 0 JR C,S8
// FD38S2   16 11 M1R 4 M1R 4 MRD 3 NON 5 ... 0 ... 0 ... 0 JR C,S8
func (c *CPU) jrCN() {
	address := c.readByte(c.PC+1, 3)
	if !c.getC() {
		c.PC += 2
		return
	}

	c.WZ = 2 + uint16(int16(c.PC)+int16(int8(address)))
	c.PC = c.WZ
	c.tstates += 5
}

// 30S2     12 07 M1R 4 MRD 3 NON 5 ... 0 ... 0 ... 0 ... 0 JR NC,S8
// DD30S2   16 11 M1R 4 M1R 4 MRD 3 NON 5 ... 0 ... 0 ... 0 JR NC,S8
// FD30S2   16 11 M1R 4 M1R 4 MRD 3 NON 5 ... 0 ... 0 ... 0 JR NC,S8
func (c *CPU) jrNcN() {
	address := c.readByte(c.PC+1, 3)

	if c.getC() {
		c.PC += 2
		return
	}

	c.WZ = 2 + uint16(int16(c.PC)+int16(int8(address)))
	c.PC = c.WZ
	c.tstates += 5
}

// 28S2     12 07 M1R 4 MRD 3 NON 5 ... 0 ... 0 ... 0 ... 0 JR Z,S8
// DD28S2   16 11 M1R 4 M1R 4 MRD 3 NON 5 ... 0 ... 0 ... 0 JR Z,S8
// FD28S2   16 11 M1R 4 M1R 4 MRD 3 NON 5 ... 0 ... 0 ... 0 JR Z,S8
func (c *CPU) jrZN() {
	address := c.readByte(c.PC+1, 3)

	if !c.getZ() {
		c.PC += 2
		return
	}

	c.WZ = 2 + uint16(int16(c.PC)+int16(int8(address)))
	c.PC = c.WZ
	c.tstates += 5
}

// 20S2     12 07 M1R 4 MRD 3 NON 5 ... 0 ... 0 ... 0 ... 0 JR NZ,S8
// DD20S2   16 11 M1R 4 M1R 4 MRD 3 NON 5 ... 0 ... 0 ... 0 JR NZ,S8
// FD20S2   16 11 M1R 4 M1R 4 MRD 3 NON 5 ... 0 ... 0 ... 0 JR NZ,S8
func (c *CPU) jrNzN() {
	address := c.readByte(c.PC+1, 3)

	if c.getZ() {
		c.PC += 2
		return
	}

	c.WZ = 2 + uint16(int16(c.PC)+int16(int8(address)))
	c.PC = c.WZ
	c.tstates += 5
}

func (c *CPU) jp_Ss_(ss string) func() {
	if ss == "HL" {
		// E9       04 00 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 ... 0 JP HL
		return func() {
			c.PC = c.HL
		}
	}

	// DDE9     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 JP IX
	// FDE9     08 00 M1R 4 M1R 4 ... 0 ... 0 ... 0 ... 0 ... 0 JP IY
	return func() {
		c.PC = c.extractRegisterPair(ss)
	}
}

// 10S1     13 08 M1R 5 MRD 3 NON 5 ... 0 ... 0 ... 0 ... 0 DJNZ S8
// DD10S1   17 12 M1R 4 M1R 5 MRD 3 NON 5 ... 0 ... 0 ... 0 DJNZ S8
// FD10S1   17 12 M1R 4 M1R 5 MRD 3 NON 5 ... 0 ... 0 ... 0 DJNZ S8
func (c *CPU) djnzN() {
	c.tstates += 1
	address := c.readByte(c.PC+1, 3)
	c.BC -= 256

	if c.BC < 256 {
		c.PC += 2
		return
	}

	c.tstates += 5 // NON

	c.WZ = 2 + uint16(int16(c.PC)+int16(int8(address)))
	c.PC = c.WZ
}

// CDL1H1   17 00 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 ... 0 CALL U16
// DDCDL1H1 21 00 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 CALL U16
// FDCDL1H1 21 00 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 CALL U16
func (c *CPU) callNn() {
	c.WZ = c.readWord(c.PC+1, 3, 4)

	c.pushStack(c.PC + 3)
	c.PC = c.WZ
}

// C4L2H2   17 10 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 ... 0 CALL NZ,U16
// DDC4L2H2 21 14 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 CALL NZ,U16
// FDC4L2H2 21 14 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 CALL NZ,U16
func (c *CPU) callNzNn() {
	c.WZ = c.readWord(c.PC+1, 3, 3)

	if c.getZ() {
		c.PC += 3
		return
	}
	c.tstates += 1 // additional tstate
	c.pushStack(c.PC + 3)
	c.PC = c.WZ
}

// CCL2H2   17 10 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 ... 0 CALL Z,U16
// DDCCL2H2 21 14 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 CALL Z,U16
// FDCCL2H2 21 14 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 CALL Z,U16
func (c *CPU) callZNn() {
	c.WZ = c.readWord(c.PC+1, 3, 3)

	if !c.getZ() {
		c.PC += 3
		return
	}
	c.tstates += 1 // additional tstate
	c.pushStack(c.PC + 3)
	c.PC = c.WZ
}

// D4L2H2   17 10 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 ... 0 CALL NC,U16
// DDD4L2H2 21 14 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 CALL NC,U16
// FDD4L2H2 21 14 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 CALL NC,U16
func (c *CPU) callNcNn() {
	c.WZ = c.readWord(c.PC+1, 3, 3)

	if c.getC() {
		c.PC += 3
		return
	}
	c.tstates += 1 // additional tstate
	c.pushStack(c.PC + 3)
	c.PC = c.WZ
}

// DCL2H2   17 10 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 ... 0 CALL C,U16
// DDDCL2H2 21 14 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 CALL C,U16
// FDDCL2H2 21 14 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 CALL C,U16
func (c *CPU) callCNn() {
	c.WZ = c.readWord(c.PC+1, 3, 3)

	if !c.getC() {
		c.PC += 3
		return
	}
	c.tstates += 1 // additional tstate
	c.pushStack(c.PC + 3)
	c.PC = c.WZ
}

// E4L2H2   17 10 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 ... 0 CALL PO,U16
// DDE4L2H2 21 14 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 CALL PO,U16
// FDE4L2H2 21 14 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 CALL PO,U16
func (c *CPU) callPoNn() {
	c.WZ = c.readWord(c.PC+1, 3, 3)

	if c.getPV() {
		c.PC += 3
		return
	}
	c.tstates += 1 // additional tstate
	c.pushStack(c.PC + 3)
	c.PC = c.WZ
}

// ECL2H2   17 10 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 ... 0 CALL PE,U16
// DDECL2H2 21 14 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 CALL PE,U16
// FDECL2H2 21 14 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 CALL PE,U16
func (c *CPU) callPeNn() {
	c.WZ = c.readWord(c.PC+1, 3, 3)

	if !c.getPV() {
		c.PC += 3
		return
	}
	c.tstates += 1 // additional tstate
	c.pushStack(c.PC + 3)
	c.PC = c.WZ
}

// F4L2H2   17 10 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 ... 0 CALL P,U16
// DDF4L2H2 21 14 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 CALL P,U16
// FDF4L2H2 21 14 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 CALL P,U16
func (c *CPU) callPNn() {
	c.WZ = c.readWord(c.PC+1, 3, 3)

	if c.getS() {
		c.PC += 3
		return
	}
	c.tstates += 1 // additional tstate
	c.pushStack(c.PC + 3)
	c.PC = c.WZ
}

// FCL2H2   17 10 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 ... 0 CALL M,U16
// DDFCL2H2 21 14 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 CALL M,U16
// FDFCL2H2 21 14 M1R 4 M1R 4 MRD 3 MRD 4 MWR 3 MWR 3 ... 0 CALL M,U16
func (c *CPU) callMNn() {
	c.WZ = c.readWord(c.PC+1, 3, 3)

	if !c.getS() {
		c.PC += 3
		return
	}
	c.tstates += 1 // additional tstate
	c.pushStack(c.PC + 3)
	c.PC = c.WZ
}

// C9       10 00 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 RET
// DDC9     14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET
// ED55     14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET
// ED5D     14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET
// ED65     14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET
// ED6D     14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET
// ED75     14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET
// ED7D     14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET
// FDC9     14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET
func (c *CPU) ret() {
	c.PC = c.popStack()
	c.WZ = c.PC
}

// C0       11 05 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 RET NZ
// DDC0     15 09 M1R 4 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET NZ
// FDC0     15 09 M1R 4 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET NZ
func (c *CPU) retNz() {
	c.tstates += 1

	if c.getZ() {
		c.PC++
		return
	}

	c.PC = c.popStack()
}

// C8       11 05 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 RET Z
// DDC8     15 09 M1R 4 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET Z
// FDC8     15 09 M1R 4 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET Z
func (c *CPU) retZ() {
	c.tstates += 1

	if !c.getZ() {
		c.PC++
		return
	}

	c.PC = c.popStack()
}

// D0       11 05 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 RET NC
// DDD0     15 09 M1R 4 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET NC
// FDD0     15 09 M1R 4 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET NC
func (c *CPU) retNc() {
	c.tstates += 1

	if c.getC() {
		c.PC++
		return
	}

	c.PC = c.popStack()
}

// D8       11 05 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 RET C
// DDD8     15 09 M1R 4 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET C
// FDD8     15 09 M1R 4 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET C
func (c *CPU) retC() {
	c.tstates += 1

	if !c.getC() {
		c.PC++
		return
	}

	c.PC = c.popStack()
}

// E0       11 05 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 RET PO
// DDE0     15 09 M1R 4 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET PO
// FDE0     15 09 M1R 4 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET PO
func (c *CPU) retPo() {
	c.tstates += 1

	if c.getPV() {
		c.PC++
		return
	}

	c.PC = c.popStack()
}

// E8       11 05 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 RET PE
// DDE8     15 09 M1R 4 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET PE
// FDE8     15 09 M1R 4 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET PE
func (c *CPU) retPe() {
	c.tstates += 1

	if !c.getPV() {
		c.PC++
		return
	}

	c.PC = c.popStack()
}

// F0       11 05 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 RET P
// DDF0     15 09 M1R 4 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET P
// FDF0     15 09 M1R 4 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET P
func (c *CPU) retP() {
	c.tstates += 1

	if c.getS() {
		c.PC++
		return
	}

	c.PC = c.popStack()
}

// F8       11 05 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 ... 0 RET M
// DDF8     15 09 M1R 4 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET M
// FDF8     15 09 M1R 4 M1R 5 MRD 3 MRD 3 ... 0 ... 0 ... 0 RET M
func (c *CPU) retM() {
	c.tstates += 1

	if !c.getS() {
		c.PC++
		return
	}

	c.PC = c.popStack()
}

// ED4D     14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 RETI
func (c *CPU) reti() {
	c.PC = c.popStack()
	c.WZ = c.PC
	c.States.IFF1 = c.States.IFF2
}

// ED45     14 00 M1R 4 M1R 4 MRD 3 MRD 3 ... 0 ... 0 ... 0 RETN
func (c *CPU) retn() {
	c.PC = c.popStack()
	c.States.IFF1 = c.States.IFF2
}

// C7       11 00 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 ... 0 RST 00H
// CF       11 00 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 ... 0 RST 08H
// D7       11 00 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 ... 0 RST 10H
// DF       11 00 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 ... 0 RST 18H
// E7       11 00 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 ... 0 RST 20H
// EF       11 00 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 ... 0 RST 28H
// F7       11 00 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 ... 0 RST 30H
// FF       11 00 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 ... 0 RST 38H
// DDC7     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 RST 00H
// DDCF     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 RST 08H
// DDD7     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 RST 10H
// DDDF     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 RST 18H
// DDE7     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 RST 20H
// DDEF     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 RST 28H
// DDF7     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 RST 30H
// DDFF     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 RST 38H
// FDC7     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 RST 00H
// FDCF     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 RST 08H
// FDD7     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 RST 10H
// FDDF     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 RST 18H
// FDE7     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 RST 20H
// FDEF     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 RST 28H
// FDF7     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 RST 30H
// FDFF     15 00 M1R 4 M1R 5 MWR 3 MWR 3 ... 0 ... 0 ... 0 RST 38H
func (c *CPU) rst(p uint8) func() {
	if p != 0x00 && p != 0x08 && p != 0x10 && p != 0x18 && p != 0x20 && p != 0x28 && p != 0x30 && p != 0x38 {
		panic("Invalid `p` value for RST")
	}

	return func() {
		c.tstates += 1
		c.pushStack(c.PC + 1)
		c.PC = uint16(p)
		c.WZ = c.PC
	}
}

// DBU2     11 00 M1R 4 MRD 3 IOR 4 ... 0 ... 0 ... 0 ... 0 IN A,(U8)
// DDDBU2   15 00 M1R 4 M1R 4 MRD 3 IOR 4 ... 0 ... 0 ... 0 IN A,(U8)
// FDDBU2   15 00 M1R 4 M1R 4 MRD 3 IOR 4 ... 0 ... 0 ... 0 IN A,(U8)
func (c *CPU) inA_N_() {
	portAddress := c.readByte(c.PC+1, 3)
	c.WZ = (uint16(c.getAcc()) << 8) + uint16(portAddress) + 1
	c.setAcc(c.getPort(portAddress, 4))

	c.PC += 2
}

// ED40     12 00 M1R 4 M1R 4 IOR 4 ... 0 ... 0 ... 0 ... 0 IN B,(C)
// ED48     12 00 M1R 4 M1R 4 IOR 4 ... 0 ... 0 ... 0 ... 0 IN C,(C)
// ED50     12 00 M1R 4 M1R 4 IOR 4 ... 0 ... 0 ... 0 ... 0 IN D,(C)
// ED58     12 00 M1R 4 M1R 4 IOR 4 ... 0 ... 0 ... 0 ... 0 IN E,(C)
// ED60     12 00 M1R 4 M1R 4 IOR 4 ... 0 ... 0 ... 0 ... 0 IN H,(C)
// ED68     12 00 M1R 4 M1R 4 IOR 4 ... 0 ... 0 ... 0 ... 0 IN L,(C)
// ED70     12 00 M1R 4 M1R 4 IOR 4 ... 0 ... 0 ... 0 ... 0 IN F,(C)
// ED78     12 00 M1R 4 M1R 4 IOR 4 ... 0 ... 0 ... 0 ... 0 IN A,(C)
func (c *CPU) inR_C_(r byte) func() {
	return func() {
		var lhigh bool
		var lvalue *uint16

		switch r {
		case 'A':
			lhigh, lvalue = true, &c.AF
		case 'B', 'C':
			lhigh, lvalue = r == 'B', &c.BC
		case 'D', 'E':
			lhigh, lvalue = r == 'D', &c.DE
		case 'H', 'L':
			lhigh, lvalue = r == 'H', &c.HL
		case ' ':
		default:
			panic("Invalid `r` part of the mnemonic")
		}

		result := c.getPort(uint8(c.BC), 4)

		if r != ' ' {
			if lhigh {
				*lvalue = (*lvalue & 0x00ff) | (uint16(result) << 8)
			} else {
				*lvalue = (*lvalue & 0xff00) | uint16(result)
			}
		}

		c.WZ = c.BC + 1

		c.setS(result > 127)
		c.setZ(result == 0)
		c.setH(false)
		c.setPV(parityTable[result])
		c.setN(false)

		c.PC += 2
	}
}

// EDA2     16 00 M1R 4 M1R 5 IOR 4 MWR 3 ... 0 ... 0 ... 0 INI
func (c *CPU) ini() {
	c.tstates += 1
	c.writeByte(c.HL, c.getPort(c.extractRegister('C'), 4), 3)
	c.WZ = c.BC + 1
	c.HL++
	c.BC -= 256

	c.setZ(c.BC < 256)
	c.setN(true)
	c.setF5(c.BC&0x2000 == 0x2000)
	c.setF3(c.BC&0x0800 == 0x0800)

	c.PC += 2
}

// EDB2     21 16 M1R 4 M1R 5 IOR 4 MWR 3 NON 5 ... 0 ... 0 INIR
func (c *CPU) inir() {
	c.ini()

	if c.extractRegister('B') == 0 {
		return
	}

	c.tstates += 5
	c.PC -= 2
}

// EDAA     16 00 M1R 4 M1R 5 IOR 4 MWR 3 ... 0 ... 0 ... 0 IND
func (c *CPU) ind() {
	c.tstates += 1
	c.writeByte(c.HL, c.getPort(c.extractRegister('C'), 4), 3)
	c.WZ = c.BC - 1
	c.HL--
	c.BC -= 256

	c.setZ(c.BC < 256)
	c.setN(true)
	c.setF5(c.BC&0x2000 == 0x2000)
	c.setF3(c.BC&0x0800 == 0x0800)

	c.PC += 2
}

// EDBA     21 16 M1R 4 M1R 5 IOR 4 MWR 3 NON 5 ... 0 ... 0 INDR
func (c *CPU) indr() {
	c.ind()

	if c.extractRegister('B') == 0 {
		return
	}

	c.tstates += 5
	c.PC -= 2
}

// D3U1     11 00 M1R 4 MRD 3 IOW 4 ... 0 ... 0 ... 0 ... 0 OUT (U8),A
// DDD3U1   15 00 M1R 4 M1R 4 MRD 3 IOW 4 ... 0 ... 0 ... 0 OUT (U8),A
// FDD3U1   15 00 M1R 4 M1R 4 MRD 3 IOW 4 ... 0 ... 0 ... 0 OUT (U8),A
func (c *CPU) out_N_A() {
	portAddress := c.readByte(c.PC+1, 3)
	c.setPort(portAddress, c.getAcc(), 4)
	c.WZ = uint16(portAddress+1) | (uint16(c.getAcc()) << 8)

	c.PC += 2
}

// ED41     12 00 M1R 4 M1R 4 IOW 4 ... 0 ... 0 ... 0 ... 0 OUT (C),B
// ED49     12 00 M1R 4 M1R 4 IOW 4 ... 0 ... 0 ... 0 ... 0 OUT (C),C
// ED51     12 00 M1R 4 M1R 4 IOW 4 ... 0 ... 0 ... 0 ... 0 OUT (C),D
// ED59     12 00 M1R 4 M1R 4 IOW 4 ... 0 ... 0 ... 0 ... 0 OUT (C),E
// ED61     12 00 M1R 4 M1R 4 IOW 4 ... 0 ... 0 ... 0 ... 0 OUT (C),H
// ED69     12 00 M1R 4 M1R 4 IOW 4 ... 0 ... 0 ... 0 ... 0 OUT (C),L
// ED71     12 00 M1R 4 M1R 4 IOW 4 ... 0 ... 0 ... 0 ... 0 OUT (C),0
// ED79     12 00 M1R 4 M1R 4 IOW 4 ... 0 ... 0 ... 0 ... 0 OUT (C),A
func (c *CPU) out_C_R(r byte) func() {
	return func() {
		var right uint8

		if r == ' ' {
			right = 0
		} else {
			right = c.extractRegister(r)
		}

		if r == 'A' {
			c.WZ = c.BC + 1
		}

		c.setPort(uint8(c.BC), right, 4)

		c.PC += 2
	}
}

// EDA3     16 00 M1R 4 M1R 5 MRD 3 IOW 4 ... 0 ... 0 ... 0 OUTI
func (c *CPU) outi() {
	c.tstates += 1
	c.setPort(c.extractRegister('C'), c.readByte(c.HL, 3), 4)
	c.HL++
	c.BC -= 256
	c.WZ = c.BC + 1

	c.setZ(c.BC < 256)
	c.setN(true)
	c.setF5(c.BC&0x2000 == 0x2000)
	c.setF3(c.BC&0x0800 == 0x0800)

	c.PC += 2
}

// EDB3     21 16 M1R 4 M1R 5 MRD 3 IOW 4 NON 5 ... 0 ... 0 OTIR
func (c *CPU) otir() {
	c.outi()

	if c.extractRegister('B') == 0 {
		return
	}

	c.PC -= 2
	c.tstates += 5
}

// EDAB     16 00 M1R 4 M1R 5 MRD 3 IOW 4 ... 0 ... 0 ... 0 OUTD
func (c *CPU) outd() {
	c.tstates += 1
	c.setPort(c.extractRegister('C'), c.readByte(c.HL, 3), 4)
	c.HL--
	c.BC -= 256
	c.WZ = c.BC - 1

	c.setZ(c.BC < 256)
	c.setN(true)
	c.setF5(c.BC&0x2000 == 0x2000)
	c.setF3(c.BC&0x0800 == 0x0800)

	c.PC += 2
}

// EDBB     21 16 M1R 4 M1R 5 MRD 3 IOW 4 NON 5 ... 0 ... 0 OTDR
func (c *CPU) otdr() {
	c.outd()

	if c.extractRegister('B') == 0 {
		return
	}

	c.PC -= 2
	c.tstates += 5
}
