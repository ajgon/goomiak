package cpu

import (
	"fmt"
	"z80/dma"
)

var parityTable [256]bool = [256]bool{
	/*	      0     1      2     3      4     5      6     7      8     9      A     B      C     D      E     F */
	/* 0 */ true, false, false, true, false, true, true, false, false, true, true, false, true, false, false, true,
	/* 1 */ false, true, true, false, true, false, false, true, true, false, false, true, false, true, true, false,
	/* 2 */ false, true, true, false, true, false, false, true, true, false, false, true, false, true, true, false,
	/* 3 */ true, false, false, true, false, true, true, false, false, true, true, false, true, false, false, true,
	/* 4 */ false, true, true, false, true, false, false, true, true, false, false, true, false, true, true, false,
	/* 5 */ true, false, false, true, false, true, true, false, false, true, true, false, true, false, false, true,
	/* 6 */ true, false, false, true, false, true, true, false, false, true, true, false, true, false, false, true,
	/* 7 */ false, true, true, false, true, false, false, true, true, false, false, true, false, true, true, false,
	/* 8 */ false, true, true, false, true, false, false, true, true, false, false, true, false, true, true, false,
	/* 9 */ true, false, false, true, false, true, true, false, false, true, true, false, true, false, false, true,
	/* A */ true, false, false, true, false, true, true, false, false, true, true, false, true, false, false, true,
	/* B */ false, true, true, false, true, false, false, true, true, false, false, true, false, true, true, false,
	/* C */ true, false, false, true, false, true, true, false, false, true, true, false, true, false, false, true,
	/* D */ false, true, true, false, true, false, false, true, true, false, false, true, false, true, true, false,
	/* E */ false, true, true, false, true, false, false, true, true, false, false, true, false, true, true, false,
	/* F */ true, false, false, true, false, true, true, false, false, true, true, false, true, false, false, true,
}

type MnemonicsDebug struct {
	base   [256]string
	xx80xx [256]string
	xxIXxx [256]string
	xxIYxx [256]string
}

var mnemonicsDebug = MnemonicsDebug{
	base: [256]string{
		"nop", "ld bc,nn", "ld (bc),a", "inc bc", "inc b", "dec b", "ld b,n", "rlca",
		"ex af,af'", "add hl,bc", "ld a,(bc)", "dec bc", "inc c", "dec c", "ld c,n", "rrca",
		"djnz n", "ld de,nn", "ld (de),a", "inc de", "inc d", "dec d", "ld d,n", "rla",
		"jr n", "add hl,de", "ld a,(de)", "dec de", "inc e", "dec e", "ld e,n", "rra",
		"jr nz,n", "ld hl,nn", "ld (nn),hl", "inc hl", "inc h", "dec h", "ld h,n", "daa",
		"jr z,n", "add hl,hl", "ld hl,(nn)", "dec hl", "inc l", "dec l", "ld l,n", "cpl",
		"jr nc,n", "ld sp,nn", "ld (nn),a", "inc sp", "inc (hl)", "dec (hl)", "ld (hl),n", "scf",
		"jr c,n", "add hl,sp", "ld a,(nn)", "dec sp", "inc a", "dec a", "ld a,n", "ccf",
		"ld b,b", "ld b,c", "ld b,d", "ld b,e", "ld b,h", "ld b,l", "ld b,(hl)", "ld b,a",
		"ld c,b", "ld c,c", "ld c,d", "ld c,e", "ld c,h", "ld c,l", "ld c,(hl)", "ld c,a",
		"ld d,b", "ld d,c", "ld d,d", "ld d,e", "ld d,h", "ld d,l", "ld d,(hl)", "ld d,a",
		"ld e,b", "ld e,c", "ld e,d", "ld e,e", "ld e,h", "ld e,l", "ld e,(hl)", "ld e,a",
		"ld h,b", "ld h,c", "ld h,d", "ld h,e", "ld h,h", "ld h,l", "ld h,(hl)", "ld h,a",
		"ld l,b", "ld l,c", "ld l,d", "ld l,e", "ld l,h", "ld l,l", "ld l,(hl)", "ld l,a",
		"ld (hl),b", "ld (hl),c", "ld (hl),d", "ld (hl),e", "ld (hl),h", "ld (hl),l", "halt", "ld (hl),a",
		"ld a,b", "ld a,c", "ld a,d", "ld a,e", "ld a,h", "ld a,l", "ld a,(hl)", "ld a,a",
		"add a,b", "add a,c", "add a,d", "add a,e", "add a,h", "add a,l", "add a,(hl)", "add a,a",
		"adc a,b", "adc a,c", "adc a,d", "adc a,e", "adc a,h", "adc a,l", "adc a,(hl)", "adc a,a",
		"sub b", "sub c", "sub d", "sub e", "sub h", "sub l", "sub (hl)", "sub a",
		"sbc a,b", "sbc a,c", "sbc a,d", "sbc a,e", "sbc a,h", "sbc a,l", "sbc a,(hl)", "sbc a,a",
		"and b", "and c", "and d", "and e", "and h", "and l", "and (hl)", "and a",
		"xor b", "xor c", "xor d", "xor e", "xor h", "xor l", "xor (hl)", "xor a",
		"or b", "or c", "or d", "or e", "or h", "or l", "or (hl)", "or a",
		"cp b", "cp c", "cp d", "cp e", "cp h", "cp l", "cp (hl)", "cp a",
		"ret nz", "pop bc", "jp nz,nn", "jp nn", "call nz,nn", "push bc", "add a,n", "rst 00h",
		"ret z", "ret", "jp z,nn", "nnBITnn", "call z,nn", "call nn", "adc a,n", "rst 08h",
		"ret nc", "pop de", "jp nc,nn", "out (n),a", "call nc,nn", "push de", "sub n", "rst 10h",
		"ret c", "enn", "jp c,nn", "in a,(n)", "call c,nn", "xxIXxx", "sbc a,n", "rst 18h",
		"ret po", "pop hl", "jp po,nn", "ex (sp),hl", "call po,nn", "push hl", "and n", "rst 20h",
		"ret pe", "jp (hl)", "jp pe,nn", "ex de,hl", "call pe,nn", "xx80xx", "xor n", "rst 28h",
		"ret p", "pop af", "jp p,nn", "di", "call p,nn", "push af", "or n", "rst 30h",
		"ret m", "ld sp,hl", "jp m,nn", "ei", "call m,nn", "xxIYxx", "cp n", "rst 38h",
	},
	xx80xx: [256]string{
		"nop", "nop", "nop", "nop", "nop", "nop", "nop", "nop",
		"nop", "nop", "nop", "nop", "nop", "nop", "nop", "nop",
		"nop", "nop", "nop", "nop", "nop", "nop", "nop", "nop",
		"nop", "nop", "nop", "nop", "nop", "nop", "nop", "nop",
		"nop", "nop", "nop", "nop", "nop", "nop", "nop", "nop",
		"nop", "nop", "nop", "nop", "nop", "nop", "nop", "nop",
		"nop", "nop", "nop", "nop", "nop", "nop", "nop", "nop",
		"nop", "nop", "nop", "nop", "nop", "nop", "nop", "nop",
		"in b,(c)", "out (c),b", "sbc hl,bc", "ld (nn),bc", "neg", "retn", "im 0", "ld i,a",
		"in c,(c)", "out (c),c", "adc hl,bc", "ld bc,(nn)", "neg", "reti", "nop", "ld r,a",
		"in d,(c)", "out (c),d", "sbc hl,de", "ld (nn),de", "neg", "retn", "im 1", "ld a,i",
		"in e,(c)", "out (c),e", "adc hl,de", "ld de,(nn)", "neg", "retn", "im 2", "ld a,r",
		"in h,(c)", "out (c),h", "sbc hl,hl", "ld (nn),hl", "neg", "retn", "nop", "ld rrd",
		"in l,(c)", "out (c),l", "adc hl,hl", "ld hl,(nn)", "neg", "retn", "nop", "ld rld",
		"in (c)", "out (c)", "sbc hl,sp", "ld (nn),sp", "neg", "retn", "nop", "nop",
		"in a,(c)", "out (c),a", "adc hl,sp", "ld sp,(nn)", "neg", "reti", "nop", "nop",
		"nop", "nop", "nop", "nop", "nop", "nop", "nop", "nop",
		"nop", "nop", "nop", "nop", "nop", "nop", "nop", "nop",
		"nop", "nop", "nop", "nop", "nop", "nop", "nop", "nop",
		"nop", "nop", "nop", "nop", "nop", "nop", "nop", "nop",
		"ldi", "cpi", "ini", "outi", "nop", "nop", "nop", "nop",
		"ldd", "cpd", "ind", "outd", "nop", "nop", "nop", "nop",
		"ldir", "cpir", "inir", "otir", "nop", "nop", "nop", "nop",
		"lddr", "cpdr", "indr", "otdr", "nop", "nop", "nop", "nop",
	},
	xxIXxx: [256]string{
		"nop", "ld bc,nn", "ld (bc),a", "inc bc", "inc b", "dec b", "ld b,n", "rlca",
		"ex af,af'", "add ix,bc", "ld a,(bc)", "dec bc", "inc c", "dec c", "ld c,n", "rrca",
		"djnz n", "ld de,nn", "ld (de),a", "inc de", "inc d", "dec d", "ld d,n", "rla",
		"jr n", "add ix,de", "ld a,(de)", "dec de", "inc e", "dec e", "ld e,n", "rra",
		"jr nz,n", "ld ix,nn", "ld (nn),ix", "inc ix", "inc ixh", "dec ixh", "ld ixh,n", "daa",
		"jr z,n", "add ix,ix", "ld ix,(nn)", "dec ix", "inc ixl", "dec ixl", "ld ixl,n", "cpl",
		"jr nc,n", "ld sp,nn", "ld (nn),a", "inc sp", "inc (ix+d)", "dec (ix+d)", "ld (ix+d),n", "scf",
		"jr c,n", "add ix,sp", "ld a,(nn)", "dec sp", "inc a", "dec a", "ld a,n", "ccf",
		"ld b,b", "ld b,c", "ld b,d", "ld b,e", "ld b,ixh", "ld b,ixl", "ld b,(ix+d)", "ld b,a",
		"ld c,b", "ld c,c", "ld c,d", "ld c,e", "ld c,ixh", "ld c,ixl", "ld c,(ix+d)", "ld c,a",
		"ld d,b", "ld d,c", "ld d,d", "ld d,e", "ld d,ixh", "ld d,ixl", "ld d,(ix+d)", "ld d,a",
		"ld e,b", "ld e,c", "ld e,d", "ld e,e", "ld e,ixh", "ld e,ixl", "ld e,(ix+d)", "ld e,a",
		"ld ixh,b", "ld ixh,c", "ld ixh,d", "ld ixh,e", "ld ixh,ixh", "ld ixh,ixl", "ld ixh,(ix+d)", "ld ixh,a",
		"ld ixl,b", "ld ixl,c", "ld ixl,d", "ld ixl,e", "ld ixl,ixh", "ld ixl,ixl", "ld ixl,(ix+d)", "ld ixl,a",
		"ld (ix+d),b", "ld (ix+d),c", "ld (ix+d),d", "ld (ix+d),e", "ld (ix+d),ixh", "ld (ix+d),ixl", "halt", "ld (ix+d),a",
		"ld a,b", "ld a,c", "ld a,d", "ld a,e", "ld a,ixh", "ld a,ixl", "ld a,(ix+d)", "ld a,a",
		"add a,b", "add a,c", "add a,d", "add a,e", "add a,ixh", "add a,ixl", "add a,(ix+d)", "add a,a",
		"adc a,b", "adc a,c", "adc a,d", "adc a,e", "adc a,ixh", "adc a,ixl", "adc a,(ix+d)", "adc a,a",
		"sub b", "sub c", "sub d", "sub e", "sub ixh", "sub ixl", "sub (ix+d)", "sub a",
		"sbc a,b", "sbc a,c", "sbc a,d", "sbc a,e", "sbc a,ixh", "sbc a,ixl", "sbc a,(ix+d)", "sbc a,a",
		"and b", "and c", "and d", "and e", "and ixh", "and ixl", "and (ix+d)", "and a",
		"xor b", "xor c", "xor d", "xor e", "xor ixh", "xor ixl", "xor (ix+d)", "xor a",
		"or b", "or c", "or d", "or e", "or ixh", "or ixl", "or (ix+d)", "or a",
		"cp b", "cp c", "cp d", "cp e", "cp ixh", "cp ixl", "cp (ix+d)", "cp a",
		"ret nz", "pop bc", "jp nz,nn", "jp nn", "call nz,nn", "push bc", "add a,n", "rst 00h",
		"ret z", "ret", "jp z,nn", "xxIXBITxx", "call z,nn", "call nn", "adc a,n", "rst 08h",
		"ret nc", "pop de", "jp nc,nn", "out (n),a", "call nc,nn", "push de", "sub n", "rst 10h",
		"ret c", "enn", "jp c,nn", "in a,(n)", "call c,nn", "nop", "sbc a,n", "rst 18h",
		"ret po", "pop ix", "jp po,nn", "ex (sp),ix", "call po,nn", "push ix", "and n", "rst 20h",
		"ret pe", "jp (ix+d)", "jp pe,nn", "ex de,ix", "call pe,nn", "nop", "xor n", "rst 28h",
		"ret p", "pop af", "jp p,nn", "di", "call p,nn", "push af", "or n", "rst 30h",
		"ret m", "ld sp,ix", "jp m,nn", "ei", "call m,nn", "nop", "cp n", "rst 38h",
	},
	xxIYxx: [256]string{
		"nop", "ld bc,nn", "ld (bc),a", "inc bc", "inc b", "dec b", "ld b,n", "rlca",
		"ex af,af'", "add iy,bc", "ld a,(bc)", "dec bc", "inc c", "dec c", "ld c,n", "rrca",
		"djnz n", "ld de,nn", "ld (de),a", "inc de", "inc d", "dec d", "ld d,n", "rla",
		"jr n", "add iy,de", "ld a,(de)", "dec de", "inc e", "dec e", "ld e,n", "rra",
		"jr nz,n", "ld iy,nn", "ld (nn),iy", "inc iy", "inc iyh", "dec iyh", "ld iyh,n", "daa",
		"jr z,n", "add iy,iy", "ld iy,(nn)", "dec iy", "inc iyl", "dec iyl", "ld iyl,n", "cpl",
		"jr nc,n", "ld sp,nn", "ld (nn),a", "inc sp", "inc (iy+d)", "dec (iy+d)", "ld (iy+d),n", "scf",
		"jr c,n", "add iy,sp", "ld a,(nn)", "dec sp", "inc a", "dec a", "ld a,n", "ccf",
		"ld b,b", "ld b,c", "ld b,d", "ld b,e", "ld b,iyh", "ld b,iyl", "ld b,(iy+d)", "ld b,a",
		"ld c,b", "ld c,c", "ld c,d", "ld c,e", "ld c,iyh", "ld c,iyl", "ld c,(iy+d)", "ld c,a",
		"ld d,b", "ld d,c", "ld d,d", "ld d,e", "ld d,iyh", "ld d,iyl", "ld d,(iy+d)", "ld d,a",
		"ld e,b", "ld e,c", "ld e,d", "ld e,e", "ld e,iyh", "ld e,iyl", "ld e,(iy+d)", "ld e,a",
		"ld iyh,b", "ld iyh,c", "ld iyh,d", "ld iyh,e", "ld iyh,iyh", "ld iyh,iyl", "ld iyh,(iy+d)", "ld iyh,a",
		"ld iyl,b", "ld iyl,c", "ld iyl,d", "ld iyl,e", "ld iyl,iyh", "ld iyl,iyl", "ld iyl,(iy+d)", "ld iyl,a",
		"ld (iy+d),b", "ld (iy+d),c", "ld (iy+d),d", "ld (iy+d),e", "ld (iy+d),iyh", "ld (iy+d),iyl", "halt", "ld (iy+d),a",
		"ld a,b", "ld a,c", "ld a,d", "ld a,e", "ld a,iyh", "ld a,iyl", "ld a,(iy+d)", "ld a,a",
		"add a,b", "add a,c", "add a,d", "add a,e", "add a,iyh", "add a,iyl", "add a,(iy+d)", "add a,a",
		"adc a,b", "adc a,c", "adc a,d", "adc a,e", "adc a,iyh", "adc a,iyl", "adc a,(iy+d)", "adc a,a",
		"sub b", "sub c", "sub d", "sub e", "sub iyh", "sub iyl", "sub (iy+d)", "sub a",
		"sbc a,b", "sbc a,c", "sbc a,d", "sbc a,e", "sbc a,iyh", "sbc a,iyl", "sbc a,(iy+d)", "sbc a,a",
		"and b", "and c", "and d", "and e", "and iyh", "and iyl", "and (iy+d)", "and a",
		"xor b", "xor c", "xor d", "xor e", "xor iyh", "xor iyl", "xor (iy+d)", "xor a",
		"or b", "or c", "or d", "or e", "or iyh", "or iyl", "or (iy+d)", "or a",
		"cp b", "cp c", "cp d", "cp e", "cp iyh", "cp iyl", "cp (iy+d)", "cp a",
		"ret nz", "pop bc", "jp nz,nn", "jp nn", "call nz,nn", "push bc", "add a,n", "rst 00h",
		"ret z", "ret", "jp z,nn", "xxIYBITxx", "call z,nn", "call nn", "adc a,n", "rst 08h",
		"ret nc", "pop de", "jp nc,nn", "out (n),a", "call nc,nn", "push de", "sub n", "rst 10h",
		"ret c", "enn", "jp c,nn", "in a,(n)", "call c,nn", "nop", "sbc a,n", "rst 18h",
		"ret po", "pop iy", "jp po,nn", "ex (sp),iy", "call po,nn", "push iy", "and n", "rst 20h",
		"ret pe", "jp (iy+d)", "jp pe,nn", "ex de,iy", "call pe,nn", "nop", "xor n", "rst 28h",
		"ret p", "pop af", "jp p,nn", "di", "call p,nn", "push af", "or n", "rst 30h",
		"ret m", "ld sp,iy", "jp m,nn", "ei", "call m,nn", "nop", "cp n", "rst 38h",
	},
}

type CPUStates struct {
	Halt  bool
	Ports [256]uint8
	IFF1  bool
	IFF2  bool
	IM    uint8
}

type CPUMnemonics struct {
	base   [256]func() uint8
	xx80xx [256]func() uint8
	xxIXxx [256]func() uint8
	xxIYxx [256]func() uint8
}

type CPU struct {
	PC     uint16
	SP     uint16
	AF     uint16
	AF_    uint16
	BC     uint16
	BC_    uint16
	DE     uint16
	DE_    uint16
	HL     uint16
	HL_    uint16
	I      uint8
	R      uint8
	IX     uint16
	IY     uint16
	States CPUStates

	dma       *dma.DMA
	mnemonics CPUMnemonics
}

func (c *CPU) initializeMnemonics() {
	for _, reg := range [3]string{"HL", "IX", "IY"} {
		baseList := [256]func() uint8{
			c.nop, c.ldBcNn, c.ld_Bc_A, c.incBc, c.incB, c.decB, c.ldBN, c.rlcR(' '),
			c.exAfAf_, c.addSsRr(reg, "BC"), c.ldA_Bc_, c.decBc, c.incC, c.decC, c.ldCN, c.rrcR(' '),
			c.djnzN, c.ldDeNn, c.ld_De_A, c.incDe, c.incD, c.decD, c.ldDN, c.rlR(' '),
			c.jrN, c.addSsRr(reg, "DE"), c.ldA_De_, c.decDe, c.incE, c.decE, c.ldEN, c.rrR(' '),
			c.jrNzN, c.ldSsNn(reg), c.ld_Nn_Ss(reg), c.incSs(reg), c.incH, c.decH, c.ldHN, c.daa,
			c.jrZN, c.addSsRr(reg, reg), c.ldSs_Nn_(reg), c.decSs(reg), c.incL, c.decL, c.ldLN, c.cpl,
			c.jrNcN, c.ldSpNn, c.ld_Nn_A, c.incSp, c.inc_Ss_(reg), c.dec_Ss_(reg), c.ld_Ss_N(reg), c.scf,
			c.jrCN, c.addSsRr(reg, "SP"), c.ldA_Nn_, c.decSp, c.incA, c.decA, c.ldAN, c.ccf,
			c.ldRR_('B', 'B'), c.ldRR_('B', 'C'), c.ldRR_('B', 'D'), c.ldRR_('B', 'E'), c.ldRR_('B', 'H'), c.ldRR_('B', 'L'), c.ldR_Ss_('B', reg), c.ldRR_('B', 'A'),
			c.ldRR_('C', 'B'), c.ldRR_('C', 'C'), c.ldRR_('C', 'D'), c.ldRR_('C', 'E'), c.ldRR_('C', 'H'), c.ldRR_('C', 'L'), c.ldR_Ss_('C', reg), c.ldRR_('C', 'A'),
			c.ldRR_('D', 'B'), c.ldRR_('D', 'C'), c.ldRR_('D', 'D'), c.ldRR_('D', 'E'), c.ldRR_('D', 'H'), c.ldRR_('D', 'L'), c.ldR_Ss_('D', reg), c.ldRR_('D', 'A'),
			c.ldRR_('E', 'B'), c.ldRR_('E', 'C'), c.ldRR_('E', 'D'), c.ldRR_('E', 'E'), c.ldRR_('E', 'H'), c.ldRR_('E', 'L'), c.ldR_Ss_('E', reg), c.ldRR_('E', 'A'),
			c.ldRR_('H', 'B'), c.ldRR_('H', 'C'), c.ldRR_('H', 'D'), c.ldRR_('H', 'E'), c.ldRR_('H', 'H'), c.ldRR_('H', 'L'), c.ldR_Ss_('H', reg), c.ldRR_('H', 'A'),
			c.ldRR_('L', 'B'), c.ldRR_('L', 'C'), c.ldRR_('L', 'D'), c.ldRR_('L', 'E'), c.ldRR_('L', 'H'), c.ldRR_('L', 'L'), c.ldR_Ss_('L', reg), c.ldRR_('L', 'A'),
			c.ld_Ss_R(reg, 'B'), c.ld_Ss_R(reg, 'C'), c.ld_Ss_R(reg, 'D'), c.ld_Ss_R(reg, 'E'), c.ld_Ss_R(reg, 'H'), c.ld_Ss_R(reg, 'L'), c.halt, c.ld_Ss_R(reg, 'A'),
			c.ldRR_('A', 'B'), c.ldRR_('A', 'C'), c.ldRR_('A', 'D'), c.ldRR_('A', 'E'), c.ldRR_('A', 'H'), c.ldRR_('A', 'L'), c.ldR_Ss_('A', reg), c.ldRR_('A', 'A'),
			c.addAR('B'), c.addAR('C'), c.addAR('D'), c.addAR('E'), c.addAR('H'), c.addAR('L'), c.addA_Ss_(reg), c.addAR('A'),
			c.adcAR('B'), c.adcAR('C'), c.adcAR('D'), c.adcAR('E'), c.adcAR('H'), c.adcAR('L'), c.adcA_Ss_(reg), c.adcAR('A'),
			c.subR('B'), c.subR('C'), c.subR('D'), c.subR('E'), c.subR('H'), c.subR('L'), c.sub_Ss_(reg), c.subR('A'),
			c.sbcAR('B'), c.sbcAR('C'), c.sbcAR('D'), c.sbcAR('E'), c.sbcAR('H'), c.sbcAR('L'), c.sbcA_Ss_(reg), c.sbcAR('A'),
			c.andR('B'), c.andR('C'), c.andR('D'), c.andR('E'), c.andR('H'), c.andR('L'), c.and_Ss_(reg), c.andR('A'),
			c.xorR('B'), c.xorR('C'), c.xorR('D'), c.xorR('E'), c.xorR('H'), c.xorR('L'), c.xor_Ss_(reg), c.xorR('A'),
			c.orR('B'), c.orR('C'), c.orR('D'), c.orR('E'), c.orR('H'), c.orR('L'), c.or_Ss_(reg), c.orR('A'),
			c.cpR('B'), c.cpR('C'), c.cpR('D'), c.cpR('E'), c.cpR('H'), c.cpR('L'), c.cp_Ss_(reg), c.cpR('A'),
			c.retNz, c.popBc, c.jpNzNn, c.jpNn, c.callNzNn, c.pushBc, c.addAN, c.rst(0x00),
			c.retZ, c.ret, c.jpZNn, c.die, c.callZNn, c.callNn, c.adcAN, c.rst(0x08),
			c.retNc, c.popDe, c.jpNcNn, c.out_N_A, c.callNcNn, c.pushDe, c.subN, c.rst(0x10),
			c.retC, c.exx, c.jpCNn, c.inA_N_, c.callCNn, c.nop, c.sbcAN, c.rst(0x18),
			c.retPo, c.popSs(reg), c.jpPoNn, c.ex_Sp_Ss(reg), c.callPoNn, c.pushSs(reg), c.andN, c.rst(0x20),
			c.retPe, c.jp_Ss_(reg), c.jpPeNn, c.exDeSs(reg), c.callPeNn, c.die, c.xorN, c.rst(0x28),
			c.retP, c.popAf, c.jpPNn, c.di, c.callPNn, c.pushAf, c.orN, c.rst(0x30),
			c.retM, c.ldSpSs(reg), c.jpMNn, c.ei, c.callMNn, c.nop, c.cpN, c.rst(0x38),
		}
		switch reg {
		case "HL":
			c.mnemonics.base = baseList
		case "IX":
			c.mnemonics.xxIXxx = baseList
		case "IY":
			c.mnemonics.xxIYxx = baseList
		}
	}
	c.mnemonics.xx80xx = [256]func() uint8{
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.inR_C_('B'), c.out_C_R('B'), c.sbcHlRr("BC"), c.ld_Nn_Rr("BC"), c.neg, c.retn, c.im(0), c.ldIA,
		c.inR_C_('C'), c.out_C_R('C'), c.adcHlRr("BC"), c.ldRr_Nn_("BC"), c.neg, c.reti, c.nop, c.ldRA,
		c.inR_C_('D'), c.out_C_R('D'), c.sbcHlRr("DE"), c.ld_Nn_Rr("DE"), c.neg, c.retn, c.im(1), c.ldAI,
		c.inR_C_('E'), c.out_C_R('E'), c.adcHlRr("DE"), c.ldRr_Nn_("DE"), c.neg, c.retn, c.im(2), c.ldAR,
		c.inR_C_('H'), c.out_C_R('H'), c.sbcHlRr("HL"), c.ld_Nn_Rr("HL"), c.neg, c.retn, c.nop, c.rrd,
		c.inR_C_('L'), c.out_C_R('L'), c.adcHlRr("HL"), c.ldRr_Nn_("HL"), c.neg, c.retn, c.nop, c.rld,
		c.inR_C_(' '), c.out_C_R(' '), c.sbcHlRr("SP"), c.ld_Nn_Rr("SP"), c.neg, c.retn, c.nop, c.nop,
		c.inR_C_('A'), c.out_C_R('A'), c.adcHlRr("SP"), c.ldRr_Nn_("SP"), c.neg, c.reti, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.ldi, c.cpi, c.ini, c.outi, c.nop, c.nop, c.nop, c.nop,
		c.ldd, c.cpd, c.ind, c.outd, c.nop, c.nop, c.nop, c.nop,
		c.ldir, c.cpir, c.inir, c.otir, c.nop, c.nop, c.nop, c.nop,
		c.lddr, c.cpdr, c.indr, c.otdr, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
	}
}

func (c *CPU) getAcc() uint8 {
	return uint8(c.AF >> 8)
}

func (c *CPU) setAcc(value uint8) {
	c.AF = (c.AF & 0x00ff) | (uint16(value) << 8)
}

func (c *CPU) getS() bool {
	return c.AF&0x0080 == 0x0080
}

func (c *CPU) getZ() bool {
	return c.AF&0x0040 == 0x0040
}

func (c *CPU) getH() bool {
	return c.AF&0x0010 == 0x0010
}

func (c *CPU) getPV() bool {
	return c.AF&0x0004 == 0x0004
}

func (c *CPU) getN() bool {
	return c.AF&0x0002 == 0x0002
}

func (c *CPU) getC() bool {
	return c.AF&0x0001 == 0x0001
}

func (c *CPU) getFlags() uint8 {
	return uint8(c.AF)
}

func (c *CPU) setS(value bool) {

	if value {
		c.AF = c.AF | 0x0080
	} else {
		c.AF = c.AF & 0xff7f
	}
}

func (c *CPU) setZ(value bool) {

	if value {
		c.AF = c.AF | 0x0040
	} else {
		c.AF = c.AF & 0xffbf
	}
}

func (c *CPU) setH(value bool) {

	if value {
		c.AF = c.AF | 0x0010
	} else {
		c.AF = c.AF & 0xffef
	}
}

func (c *CPU) setPV(value bool) {

	if value {
		c.AF = c.AF | 0x0004
	} else {
		c.AF = c.AF & 0xfffb
	}
}

func (c *CPU) setN(value bool) {

	if value {
		c.AF = c.AF | 0x0002
	} else {
		c.AF = c.AF & 0xfffd
	}
}

func (c *CPU) setC(value bool) {

	if value {
		c.AF = c.AF | 0x0001
	} else {
		c.AF = c.AF & 0xfffe
	}
}

func (c *CPU) setFlags(value uint8) {
	c.AF = (c.AF & 0xff00) | uint16(value)
}

func (c *CPU) popStack() (value uint16) {
	value = c.readWord(c.SP)
	c.SP += 2

	return
}

func (c *CPU) pushStack(value uint16) {
	c.SP -= 2
	c.writeWord(c.SP, value)
}

func (c *CPU) getPort(addr uint8) uint8 {
	return c.States.Ports[addr]
}

func (c *CPU) setPort(addr uint8, value uint8) {
	c.States.Ports[addr] = value
}

func (c *CPU) disableInterrupts() {
	c.States.IFF1 = false
	c.States.IFF2 = false
}

func (c *CPU) enableInterrupts() {
	c.States.IFF1 = true
	c.States.IFF2 = true
}

func (c *CPU) checkInterrupts() (bool, bool) {
	return c.States.IFF1, c.States.IFF2
}

// reads word and maintains endianess
// example:
// 0040 34 21
// readWord(0x0040) => 0x1234
func (c *CPU) readWord(address uint16) uint16 {
	return uint16(c.dma.GetMemory(address+1))<<8 | uint16(c.dma.GetMemory(address))
}

// writes word to given address and address+1 and maintains endianess
// example:
// writeWord(0x1234, 0x5678)
// 1234  78 56
func (c *CPU) writeWord(address uint16, value uint16) {
	c.dma.SetMemoryBulk(address, []uint8{uint8(value), uint8(value >> 8)})
}

func (c *CPU) extractRegister(r byte) uint8 {
	switch r {
	case 'A':
		return c.getAcc()
	case 'B':
		return uint8(c.BC >> 8)
	case 'C':
		return uint8(c.BC)
	case 'D':
		return uint8(c.DE >> 8)
	case 'E':
		return uint8(c.DE)
	case 'H':
		return uint8(c.HL >> 8)
	case 'L':
		return uint8(c.HL)
	}

	panic("Invalid `r` part of the mnemonic")
}

func (c *CPU) extractRegisterPair(rr string) (rvalue uint16) {
	switch rr {
	case "AF":
		rvalue = c.AF
	case "BC":
		rvalue = c.BC
	case "DE":
		rvalue = c.DE
	case "HL":
		rvalue = c.HL
	case "SP":
		rvalue = c.SP
	case "IX":
		rvalue = c.IX
	case "IY":
		rvalue = c.IY
	default:
		panic("Invalid `rr` part of the mnemonic")
	}

	return
}

func (c *CPU) increaseRegister(name rune) uint8 {
	var register uint8

	switch name {
	case 'A':
		c.AF += 256
		register = c.getAcc()
	case 'B':
		c.BC += 256
		register = uint8(c.BC >> 8)
	case 'C':
		register = uint8(c.BC) + 1
		c.BC = (c.BC & 0xff00) | uint16(register)
	case 'D':
		c.DE += 256
		register = uint8(c.DE >> 8)
	case 'E':
		register = uint8(c.DE) + 1
		c.DE = (c.DE & 0xff00) | uint16(register)
	case 'H':
		c.HL += 256
		register = uint8(c.HL >> 8)
	case 'L':
		register = uint8(c.HL) + 1
		c.HL = (c.HL & 0xff00) | uint16(register)
	}

	c.setN(false)
	c.setPV(register == 0x80)
	c.setH(register&0x0f == 0)
	c.setZ(register == 0)
	c.setS(register > 127)
	c.PC++

	return 4
}

func (c *CPU) decreaseRegister(name rune) uint8 {
	var register uint8

	switch name {
	case 'A':
		c.AF -= 256
		register = c.getAcc()
	case 'B':
		c.BC -= 256
		register = uint8(c.BC >> 8)
	case 'C':
		register = uint8(c.BC) - 1
		c.BC = (c.BC & 0xff00) | uint16(register)
	case 'D':
		c.DE -= 256
		register = uint8(c.DE >> 8)
	case 'E':
		register = uint8(c.DE) - 1
		c.DE = (c.DE & 0xff00) | uint16(register)
	case 'H':
		c.HL -= 256
		register = uint8(c.HL >> 8)
	case 'L':
		register = uint8(c.HL) - 1
		c.HL = (c.HL & 0xff00) | uint16(register)
	}

	c.setN(true)
	c.setPV(register == 0x7f)
	c.setH(register&0x0f == 0x0f)
	c.setZ(register == 0)
	c.setS(register > 127)

	c.PC++
	return 4
}

// left stores the result
func (c *CPU) addRegisters(left *uint16, right uint16) {
	sum := *left + right

	c.setC(sum < *left || sum < right)
	c.setN(false)
	c.setH((*left^right^sum)&0x1000 == 0x1000)

	*left = sum
}

func (c *CPU) adcValueToAcc(value uint8) {
	var carryIn, carryOut uint8

	if c.getC() {
		carryIn = 1
	}

	a := c.getAcc()
	result := a + value + carryIn
	c.setAcc(result)

	if c.getC() {
		c.setC(a >= 0xff-value)
	} else {
		c.setC(a > 0xff-value)
	}

	c.setN(false)

	if c.getC() {
		carryOut = 1
	}

	c.setPV((((result ^ a ^ value) >> 7) ^ carryOut) == 1)

	c.setH((a^value^result)&0x10 == 0x10)
	c.setZ(result == 0)
	c.setS(result > 127)
}

func (c *CPU) adc16bit(addendLeft, addendRight uint16) (result uint16) {
	var carryIn, carryOut uint16

	if c.getC() {
		carryIn = 1
	}

	result = addendLeft + addendRight + carryIn

	if c.getC() {
		c.setC(addendLeft >= 0xffff-addendRight)
	} else {
		c.setC(addendLeft > 0xffff-addendRight)
	}

	c.setN(false)

	if c.getC() {
		carryOut = 1
	}

	c.setPV((((result ^ addendLeft ^ addendRight) >> 15) ^ carryOut) == 1)

	c.setH((addendLeft^addendRight^result)&0x1000 == 0x1000)
	c.setZ(result == 0)
	c.setS(result > 0x7fff)

	return
}

func (c *CPU) nop() uint8 {
	c.PC++

	return 4
}

func (c *CPU) ldBcNn() uint8 {
	c.BC = c.readWord(c.PC + 1)
	c.PC += 3

	return 10
}

func (c *CPU) ld_Bc_A() uint8 {
	c.dma.SetMemoryByte(c.BC, c.getAcc())
	c.PC++
	return 7
}

func (c *CPU) incBc() uint8 {
	c.BC++
	c.PC++
	return 6
}

func (c *CPU) incB() uint8 {
	return c.increaseRegister('B')
}

func (c *CPU) decB() uint8 {
	return c.decreaseRegister('B')
}

func (c *CPU) ldBN() uint8 {
	c.BC = (c.BC & 0x00ff) | (uint16(c.dma.GetMemory(c.PC+1)) << 8)
	c.PC += 2

	return 7
}

func (c *CPU) rlcR(r byte) func() uint8 {
	return func() uint8 {
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
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setZ(rvalue == 0)
		c.setS(rvalue > 127)

		return 4 + size*4
	}
}

func (c *CPU) exAfAf_() uint8 {
	c.AF, c.AF_ = c.AF_, c.AF
	c.PC++

	return 4
}

func (c *CPU) addSsRr(ss, rr string) func() uint8 {
	switch ss {
	case "HL":
		return func() uint8 {
			c.addRegisters(&c.HL, c.extractRegisterPair(rr))
			c.PC++

			return 11
		}
	case "IX":
		return func() uint8 {
			c.addRegisters(&c.IX, c.extractRegisterPair(rr))
			c.PC += 2

			return 15
		}
	case "IY":
		return func() uint8 {
			c.addRegisters(&c.IY, c.extractRegisterPair(rr))
			c.PC += 2

			return 15
		}
	}

	panic("Invalid `ss` type")
}

func (c *CPU) ldA_Bc_() uint8 {
	value := c.dma.GetMemory(c.BC)
	c.setAcc(value)
	c.PC++

	return 7
}

func (c *CPU) decBc() uint8 {
	c.BC--
	c.PC++

	return 6
}

func (c *CPU) incC() uint8 {
	return c.increaseRegister('C')
}

func (c *CPU) decC() uint8 {
	return c.decreaseRegister('C')
}

func (c *CPU) ldCN() uint8 {
	c.BC = (c.BC & 0xff00) | uint16(c.dma.GetMemory(c.PC+1))
	c.PC += 2

	return 7
}

func (c *CPU) djnzN() uint8 {
	c.BC -= 256
	if c.BC < 256 {
		c.PC += 2
		return 8
	}

	c.PC = 2 + uint16(int16(c.PC)+int16(int8(c.dma.GetMemory(c.PC+1))))
	return 13
}

func (c *CPU) ldDeNn() uint8 {
	c.DE = c.readWord(c.PC + 1)
	c.PC += 3

	return 10
}

func (c *CPU) ld_De_A() uint8 {
	c.dma.SetMemoryByte(c.DE, c.getAcc())
	c.PC++
	return 7
}

func (c *CPU) incDe() uint8 {
	c.DE++
	c.PC++
	return 6
}

func (c *CPU) incD() uint8 {
	return c.increaseRegister('D')
}

func (c *CPU) decD() uint8 {
	return c.decreaseRegister('D')
}

func (c *CPU) ldDN() uint8 {
	c.DE = (c.DE & 0x00ff) | (uint16(c.dma.GetMemory(c.PC+1)) << 8)
	c.PC += 2

	return 7
}

func (c *CPU) rlR(r byte) func() uint8 {
	return func() uint8 {
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
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setZ(rvalue == 0)
		c.setS(rvalue > 127)

		return 4 + size*4
	}
}

func (c *CPU) rlSs(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			rvalue := c.dma.GetMemory(c.HL)

			signed := rvalue&128 == 128
			rvalue = rvalue << 1
			c.PC += 2

			if c.getC() {
				rvalue = rvalue | 0b00000001
			} else {
				rvalue = rvalue & 0b11111110
			}

			c.dma.SetMemoryByte(c.HL, rvalue)

			c.setC(signed)
			c.setN(false)
			c.setPV(parityTable[rvalue])
			c.setH(false)
			c.setZ(rvalue == 0)
			c.setS(rvalue > 127)

			return 15
		}
	}

	return func() uint8 {
		address := c.extractRegisterPair(ss) + uint16(c.dma.GetMemory(c.PC+3))
		rvalue := c.dma.GetMemory(address)

		signed := rvalue&128 == 128
		rvalue = rvalue << 1
		c.PC += 4

		if c.getC() {
			rvalue = rvalue | 0b00000001
		} else {
			rvalue = rvalue & 0b11111110
		}

		c.dma.SetMemoryByte(address, rvalue)

		c.setC(signed)
		c.setN(false)
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setZ(rvalue == 0)
		c.setS(rvalue > 127)

		return 23
	}
}

func (c *CPU) jrN() uint8 {
	c.PC = 2 + uint16(int16(c.PC)+int16(int8(c.dma.GetMemory(c.PC+1))))

	return 12
}

func (c *CPU) ldA_De_() uint8 {
	value := c.dma.GetMemory(c.DE)
	c.setAcc(value)
	c.PC++

	return 7
}

func (c *CPU) decDe() uint8 {
	c.DE--
	c.PC++

	return 6
}

func (c *CPU) incE() uint8 {
	return c.increaseRegister('E')
}

func (c *CPU) decE() uint8 {
	return c.decreaseRegister('E')
}

func (c *CPU) ldEN() uint8 {
	c.DE = (c.DE & 0xff00) | uint16(c.dma.GetMemory(c.PC+1))
	c.PC += 2

	return 7
}

func (c *CPU) rrR(r byte) func() uint8 {
	return func() uint8 {
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
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setZ(rvalue == 0)
		c.setS(rvalue > 127)

		return 4 + size*4
	}
}

func (c *CPU) rrSs(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			rvalue := c.dma.GetMemory(c.HL)

			signed := rvalue&1 == 1
			rvalue = rvalue >> 1
			c.PC += 2

			if c.getC() {
				rvalue = rvalue | 0b10000000
			} else {
				rvalue = rvalue & 0b01111111
			}

			c.dma.SetMemoryByte(c.HL, rvalue)

			c.setC(signed)
			c.setN(false)
			c.setPV(parityTable[rvalue])
			c.setH(false)
			c.setZ(rvalue == 0)
			c.setS(rvalue > 127)

			return 15
		}
	}

	return func() uint8 {
		address := c.extractRegisterPair(ss) + uint16(c.dma.GetMemory(c.PC+3))
		rvalue := c.dma.GetMemory(address)

		signed := rvalue&1 == 1
		rvalue = rvalue >> 1
		c.PC += 4

		if c.getC() {
			rvalue = rvalue | 0b10000000
		} else {
			rvalue = rvalue & 0b01111111
		}

		c.dma.SetMemoryByte(address, rvalue)

		c.setC(signed)
		c.setN(false)
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setZ(rvalue == 0)
		c.setS(rvalue > 127)

		return 23
	}
}

func (c *CPU) jrNzN() uint8 {
	if c.getZ() {
		c.PC += 2
		return 7
	}

	c.PC = 2 + uint16(int16(c.PC)+int16(int8(c.dma.GetMemory(c.PC+1))))
	return 12
}

func (c *CPU) ldSsNn(ss string) func() uint8 {
	switch ss {
	case "HL":
		return func() uint8 {
			c.HL = c.readWord(c.PC + 1)
			c.PC += 3

			return 10
		}
	case "IX":
		return func() uint8 {
			c.IX = c.readWord(c.PC + 2)
			c.PC += 4

			return 14
		}
	case "IY":
		return func() uint8 {
			c.IY = c.readWord(c.PC + 2)
			c.PC += 4

			return 14
		}
	}

	panic("Invalid `ss` type")
}

func (c *CPU) ld_Nn_Ss(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			c.writeWord(c.readWord(c.PC+1), c.HL)
			c.PC += 3
			return 5
		}
	}

	return func() uint8 {
		c.writeWord(c.readWord(c.PC+2), c.extractRegisterPair(ss))
		c.PC += 4
		return 20
	}
}

func (c *CPU) incSs(ss string) func() uint8 {
	switch ss {
	case "HL":
		return func() uint8 {
			c.HL++
			c.PC++
			return 6
		}
	case "IX":
		return func() uint8 {
			c.IX++
			c.PC += 2
			return 10
		}
	case "IY":
		return func() uint8 {
			c.IY++
			c.PC += 2
			return 10
		}
	}

	panic("Invalid `ss` type")
}

func (c *CPU) incH() uint8 {
	return c.increaseRegister('H')
}

func (c *CPU) decH() uint8 {
	return c.decreaseRegister('H')
}

func (c *CPU) ldHN() uint8 {
	c.HL = (c.HL & 0x00ff) | (uint16(c.dma.GetMemory(c.PC+1)) << 8)
	c.PC += 2

	return 7
}

func (c *CPU) daa() uint8 {
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

	c.setAcc(a)

	c.PC++

	return 4
}

func (c *CPU) jrZN() uint8 {
	if !c.getZ() {
		c.PC += 2
		return 7
	}

	c.PC = 2 + uint16(int16(c.PC)+int16(int8(c.dma.GetMemory(c.PC+1))))
	return 12
}

func (c *CPU) ldSs_Nn_(ss string) func() uint8 {
	switch ss {
	case "HL":
		return func() uint8 {
			c.HL = c.readWord(c.readWord(c.PC + 1))
			c.PC += 3

			return 16
		}
	case "IX":
		return func() uint8 {
			c.IX = c.readWord(c.readWord(c.PC + 2))
			c.PC += 4
			return 20
		}
	case "IY":
		return func() uint8 {
			c.IY = c.readWord(c.readWord(c.PC + 2))
			c.PC += 4
			return 20
		}
	}

	panic("Invalid `ss` type")
}

func (c *CPU) decSs(ss string) func() uint8 {
	switch ss {
	case "HL":
		return func() uint8 {
			c.HL--
			c.PC++

			return 6
		}
	case "IX":
		return func() uint8 {
			c.IX--
			c.PC += 2

			return 10
		}
	case "IY":
		return func() uint8 {
			c.IY--
			c.PC += 2

			return 10
		}
	}

	panic("Invalid `ss` type")
}

func (c *CPU) incL() uint8 {
	return c.increaseRegister('L')
}

func (c *CPU) decL() uint8 {
	return c.decreaseRegister('L')
}

func (c *CPU) ldLN() uint8 {
	c.HL = (c.HL & 0xff00) | uint16(c.dma.GetMemory(c.PC+1))
	c.PC += 2

	return 7
}

func (c *CPU) cpl() uint8 {
	c.setAcc(c.getAcc() ^ 0xff)
	c.PC++
	c.setH(true)
	c.setN(true)

	return 4
}

func (c *CPU) jrNcN() uint8 {
	if c.getC() {
		c.PC += 2
		return 7
	}

	c.PC = 2 + uint16(int16(c.PC)+int16(int8(c.dma.GetMemory(c.PC+1))))
	return 12
}

func (c *CPU) ldSpNn() uint8 {
	c.SP = c.readWord(c.PC + 1)
	c.PC += 3

	return 10
}

func (c *CPU) ld_Nn_A() uint8 {
	c.dma.SetMemoryByte(c.readWord(c.PC+1), c.getAcc())
	c.PC += 3
	return 13
}

func (c *CPU) incSp() uint8 {
	c.SP++
	c.PC++
	return 6
}

func (c *CPU) inc_Ss_(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			result := c.dma.GetMemory(c.HL) + 1
			c.dma.SetMemoryByte(c.HL, result)
			c.PC++

			c.setN(false)
			c.setPV(result == 0x80)
			c.setH(result&0x0f == 0)
			c.setZ(result == 0)
			c.setS(result > 127)

			return 11

		}
	}

	return func() uint8 {
		addr := c.extractRegisterPair(ss) + uint16(c.dma.GetMemory(c.PC+2))
		result := c.dma.GetMemory(addr) + 1
		c.dma.SetMemoryByte(addr, result)
		c.PC += 3

		c.setN(false)
		c.setPV(result == 0x80)
		c.setH(result&0x0f == 0)
		c.setZ(result == 0)
		c.setS(result > 127)

		return 23
	}
}

func (c *CPU) dec_Ss_(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			result := c.dma.GetMemory(c.HL) - 1
			c.dma.SetMemoryByte(c.HL, result)
			c.PC++

			c.setN(true)
			c.setPV(result == 0x7f)
			c.setH(result&0x0f == 0x0f)
			c.setZ(result == 0)
			c.setS(result > 127)

			return 11

		}
	}

	return func() uint8 {
		addr := c.extractRegisterPair(ss) + uint16(c.dma.GetMemory(c.PC+2))
		result := c.dma.GetMemory(addr) - 1
		c.dma.SetMemoryByte(addr, result)
		c.PC += 3

		c.setN(true)
		c.setPV(result == 0x7f)
		c.setH(result&0x0f == 0x0f)
		c.setZ(result == 0)
		c.setS(result > 127)

		return 23
	}
}

func (c *CPU) ld_Ss_N(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			c.dma.SetMemoryByte(c.HL, c.dma.GetMemory(c.PC+1))
			c.PC += 2
			return 10
		}
	}

	return func() uint8 {
		c.dma.SetMemoryByte(c.extractRegisterPair(ss)+uint16(c.dma.GetMemory(c.PC+2)), c.dma.GetMemory(c.PC+3))
		c.PC += 4
		return 19
	}
}

func (c *CPU) scf() uint8 {
	c.PC++

	c.setC(true)
	c.setN(false)
	c.setH(false)

	return 4
}

func (c *CPU) jrCN() uint8 {
	if !c.getC() {
		c.PC += 2
		return 7
	}

	c.PC = 2 + uint16(int16(c.PC)+int16(int8(c.dma.GetMemory(c.PC+1))))
	return 12
}

func (c *CPU) ldA_Nn_() uint8 {
	c.setAcc(c.dma.GetMemory(c.readWord(c.PC + 1)))
	c.PC += 3

	return 13
}

func (c *CPU) decSp() uint8 {
	c.SP--
	c.PC++

	return 6
}

func (c *CPU) incA() uint8 {
	return c.increaseRegister('A')
}

func (c *CPU) decA() uint8 {
	return c.decreaseRegister('A')
}

func (c *CPU) ldAN() uint8 {
	c.setAcc(c.dma.GetMemory(c.PC + 1))
	c.PC += 2

	return 7
}

func (c *CPU) ccf() uint8 {
	c.PC++

	c.setH(c.getC())
	c.setN(false)
	c.setC(!c.getC())

	return 4
}

func (c *CPU) ldRR_(r, r_ byte) func() uint8 {
	return func() uint8 {
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

		right := c.extractRegister(r_)

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(right) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(right)
		}

		c.PC++

		return 4
	}
}

func (c *CPU) ldR_Ss_(r byte, ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
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

			right = c.dma.GetMemory(c.HL)

			if lhigh {
				*lvalue = (*lvalue & 0x00ff) | (uint16(right) << 8)
			} else {
				*lvalue = (*lvalue & 0xff00) | uint16(right)
			}

			c.PC++
			return 7
		}
	}

	return func() uint8 {
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

		right = c.dma.GetMemory(c.extractRegisterPair(ss) + uint16(c.dma.GetMemory(c.PC+2)))

		if lhigh {
			*lvalue = (*lvalue & 0x00ff) | (uint16(right) << 8)
		} else {
			*lvalue = (*lvalue & 0xff00) | uint16(right)
		}

		c.PC += 3
		return 19
	}

}

func (c *CPU) ld_Ss_R(ss string, r byte) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			c.dma.SetMemoryByte(c.HL, c.extractRegister(r))
			c.PC++
			return 7

		}
	}

	return func() uint8 {
		c.dma.SetMemoryByte(c.extractRegisterPair(ss)+uint16(c.dma.GetMemory(c.PC+2)), c.extractRegister(r))
		c.PC += 3
		return 19
	}
}

func (c *CPU) halt() uint8 {
	c.PC++
	c.States.Halt = true

	return 4
}

func (c *CPU) addAR(r byte) func() uint8 {
	return func() uint8 {
		c.setC(false)
		c.adcValueToAcc(c.extractRegister(r))

		c.PC++

		return 4
	}
}

func (c *CPU) addA_Ss_(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			c.setC(false)
			c.adcValueToAcc(c.dma.GetMemory(c.HL))
			c.PC++
			return 7
		}
	}

	return func() uint8 {
		c.setC(false)
		c.adcValueToAcc(c.dma.GetMemory(c.extractRegisterPair(ss) + uint16(c.dma.GetMemory(c.PC+2))))
		c.PC += 3
		return 19
	}
}

func (c *CPU) adcAR(r byte) func() uint8 {
	return func() uint8 {
		c.adcValueToAcc(c.extractRegister(r))
		c.PC++

		return 4
	}
}

func (c *CPU) adcA_Ss_(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			c.adcValueToAcc(c.dma.GetMemory(c.HL))
			c.PC++
			return 7
		}
	}

	return func() uint8 {
		c.adcValueToAcc(c.dma.GetMemory(c.extractRegisterPair(ss) + uint16(c.dma.GetMemory(c.PC+2))))
		c.PC += 3
		return 19
	}
}

func (c *CPU) subR(r byte) func() uint8 {
	return func() uint8 {
		c.setC(true)
		c.adcValueToAcc(c.extractRegister(r) ^ 0xff)

		c.PC++
		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())

		return 4
	}
}

func (c *CPU) sub_Ss_(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			c.setC(true)
			c.adcValueToAcc(c.dma.GetMemory(c.HL) ^ 0xff)

			c.PC++
			c.setN(true)
			c.setC(!c.getC())
			c.setH(!c.getH())

			return 7
		}
	}

	return func() uint8 {
		c.setC(true)
		c.adcValueToAcc(c.dma.GetMemory(c.extractRegisterPair(ss)+uint16(c.dma.GetMemory(c.PC+2))) ^ 0xff)

		c.PC += 3
		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())

		return 19
	}
}

func (c *CPU) sbcAR(r byte) func() uint8 {
	return func() uint8 {
		c.setC(!c.getC())
		c.adcValueToAcc(c.extractRegister(r) ^ 0xff)

		c.PC++
		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())

		return 4
	}
}

func (c *CPU) sbcA_Ss_(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			c.setC(!c.getC())
			c.adcValueToAcc(c.dma.GetMemory(c.HL) ^ 0xff)

			c.PC++
			c.setN(true)
			c.setC(!c.getC())
			c.setH(!c.getH())

			return 7
		}
	}

	return func() uint8 {
		c.setC(!c.getC())
		c.adcValueToAcc(c.dma.GetMemory(c.extractRegisterPair(ss)+uint16(c.dma.GetMemory(c.PC+2))) ^ 0xff)

		c.PC += 3
		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())

		return 19
	}
}

func (c *CPU) andR(r byte) func() uint8 {
	return func() uint8 {
		var result uint8
		result = c.getAcc() & c.extractRegister(r)

		c.PC++
		c.setAcc(result)
		c.setS(result > 127)
		c.setZ(result == 0)
		c.setH(true)
		c.setPV(parityTable[result])
		c.setN(false)
		c.setC(false)

		return 4
	}
}

func (c *CPU) and_Ss_(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			result := c.getAcc() & c.dma.GetMemory(c.HL)

			c.PC++
			c.setAcc(result)
			c.setS(result > 127)
			c.setZ(result == 0)
			c.setH(true)
			c.setPV(parityTable[result])
			c.setN(false)
			c.setC(false)

			return 7
		}
	}

	return func() uint8 {
		result := c.getAcc() & c.dma.GetMemory(c.extractRegisterPair(ss)+uint16(c.dma.GetMemory(c.PC+2)))

		c.PC += 3
		c.setAcc(result)
		c.setS(result > 127)
		c.setZ(result == 0)
		c.setH(true)
		c.setPV(parityTable[result])
		c.setN(false)
		c.setC(false)

		return 19
	}
}

func (c *CPU) xorR(r byte) func() uint8 {
	return func() uint8 {
		var result uint8
		result = c.getAcc() ^ c.extractRegister(r)

		c.PC++
		c.setAcc(result)
		c.setS(result > 127)
		c.setZ(result == 0)
		c.setH(false)
		c.setPV(parityTable[result])
		c.setN(false)
		c.setC(false)

		return 4
	}
}

func (c *CPU) xor_Ss_(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			result := c.getAcc() ^ c.dma.GetMemory(c.HL)

			c.PC++
			c.setAcc(result)
			c.setS(result > 127)
			c.setZ(result == 0)
			c.setH(false)
			c.setPV(parityTable[result])
			c.setN(false)
			c.setC(false)

			return 7
		}
	}

	return func() uint8 {
		result := c.getAcc() ^ c.dma.GetMemory(c.extractRegisterPair(ss)+uint16(c.dma.GetMemory(c.PC+2)))

		c.PC += 3
		c.setAcc(result)
		c.setS(result > 127)
		c.setZ(result == 0)
		c.setH(false)
		c.setPV(parityTable[result])
		c.setN(false)
		c.setC(false)

		return 19
	}
}

func (c *CPU) orR(r byte) func() uint8 {
	return func() uint8 {
		var result uint8
		result = c.getAcc() | c.extractRegister(r)

		c.PC++
		c.setAcc(result)
		c.setS(result > 127)
		c.setZ(result == 0)
		c.setH(false)
		c.setPV(parityTable[result])
		c.setN(false)
		c.setC(false)

		return 4
	}
}

func (c *CPU) or_Ss_(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			result := c.getAcc() | c.dma.GetMemory(c.HL)

			c.PC++
			c.setAcc(result)
			c.setS(result > 127)
			c.setZ(result == 0)
			c.setH(false)
			c.setPV(parityTable[result])
			c.setN(false)
			c.setC(false)

			return 7
		}
	}

	return func() uint8 {
		result := c.getAcc() | c.dma.GetMemory(c.extractRegisterPair(ss)+uint16(c.dma.GetMemory(c.PC+2)))

		c.PC += 3
		c.setAcc(result)
		c.setS(result > 127)
		c.setZ(result == 0)
		c.setH(false)
		c.setPV(parityTable[result])
		c.setN(false)
		c.setC(false)

		return 19
	}
}

func (c *CPU) cpR(r byte) func() uint8 {
	return func() uint8 {
		acc := c.getAcc()
		c.setC(true)
		c.adcValueToAcc(c.extractRegister(r) ^ 0xff)

		c.PC++
		c.setAcc(acc)
		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())

		return 4
	}
}

func (c *CPU) cp_Ss_(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			acc := c.getAcc()
			c.setC(true)
			c.adcValueToAcc(c.dma.GetMemory(c.HL) ^ 0xff)

			c.PC++
			c.setAcc(acc)
			c.setN(true)
			c.setC(!c.getC())
			c.setH(!c.getH())

			return 7
		}
	}

	return func() uint8 {
		acc := c.getAcc()
		c.setC(true)
		c.adcValueToAcc(c.dma.GetMemory(c.extractRegisterPair(ss)+uint16(c.dma.GetMemory(c.PC+2))) ^ 0xff)

		c.PC += 3
		c.setAcc(acc)
		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())

		return 19
	}
}

func (c *CPU) retNz() uint8 {
	if c.getZ() {
		c.PC++
		return 5
	}

	c.PC = c.popStack()

	return 11
}

func (c *CPU) popBc() uint8 {
	c.BC = c.popStack()
	c.PC++

	return 10
}

func (c *CPU) jpNzNn() uint8 {
	if c.getZ() {
		c.PC += 3
		return 10
	}

	c.PC = c.readWord(c.PC + 1)
	return 10
}

func (c *CPU) jpNn() uint8 {
	c.PC = c.readWord(c.PC + 1)
	return 10
}

func (c *CPU) callNzNn() uint8 {
	if c.getZ() {
		c.PC += 3
		return 10
	}
	c.pushStack(c.PC + 3)
	c.PC = c.readWord(c.PC + 1)

	return 17
}

func (c *CPU) pushBc() uint8 {
	c.pushStack(c.BC)
	c.PC++

	return 11
}

func (c *CPU) addAN() uint8 {
	c.setC(false)
	c.adcValueToAcc(c.dma.GetMemory(c.PC + 1))

	c.PC += 2

	return 7
}

func (c *CPU) rst(p uint8) func() uint8 {
	if p != 0x00 && p != 0x08 && p != 0x10 && p != 0x18 && p != 0x20 && p != 0x28 && p != 0x30 && p != 0x38 {
		panic("Invalid `p` value for RST")
	}

	return func() uint8 {
		c.pushStack(c.PC + 1)
		c.PC = uint16(p)

		return 11
	}
}

func (c *CPU) retZ() uint8 {
	if !c.getZ() {
		c.PC++
		return 5
	}

	c.PC = c.popStack()

	return 11
}

func (c *CPU) ret() uint8 {
	c.PC = c.popStack()

	return 10
}

func (c *CPU) jpZNn() uint8 {
	if !c.getZ() {
		c.PC += 3
		return 10
	}

	c.PC = c.readWord(c.PC + 1)
	return 10
}

func (c *CPU) callZNn() uint8 {
	if !c.getZ() {
		c.PC += 3
		return 10
	}
	c.pushStack(c.PC + 3)
	c.PC = c.readWord(c.PC + 1)

	return 17
}

func (c *CPU) callNn() uint8 {
	c.pushStack(c.PC + 3)
	c.PC = c.readWord(c.PC + 1)

	return 17
}

func (c *CPU) adcAN() uint8 {
	c.adcValueToAcc(c.dma.GetMemory(c.PC + 1))

	c.PC += 2
	return 7
}

func (c *CPU) retNc() uint8 {
	if c.getC() {
		c.PC++
		return 5
	}

	c.PC = c.popStack()

	return 11
}

func (c *CPU) popDe() uint8 {
	c.DE = c.popStack()
	c.PC++

	return 10
}

func (c *CPU) jpNcNn() uint8 {
	if c.getC() {
		c.PC += 3
		return 10
	}

	c.PC = c.readWord(c.PC + 1)
	return 10
}

func (c *CPU) out_N_A() uint8 {
	c.setPort(c.dma.GetMemory(c.PC+1), c.getAcc())

	c.PC += 2
	return 11
}

func (c *CPU) callNcNn() uint8 {
	if c.getC() {
		c.PC += 3
		return 10
	}
	c.pushStack(c.PC + 3)
	c.PC = c.readWord(c.PC + 1)

	return 17
}

func (c *CPU) pushDe() uint8 {
	c.pushStack(c.DE)
	c.PC++

	return 11
}

func (c *CPU) subN() uint8 {
	c.setC(true)
	c.adcValueToAcc(c.dma.GetMemory(c.PC+1) ^ 0xff)

	c.PC += 2
	c.setN(true)
	c.setC(!c.getC())
	c.setH(!c.getH())

	return 7
}

func (c *CPU) retC() uint8 {
	if !c.getC() {
		c.PC++
		return 5
	}

	c.PC = c.popStack()

	return 11
}

func (c *CPU) exx() uint8 {
	c.BC, c.BC_ = c.BC_, c.BC
	c.DE, c.DE_ = c.DE_, c.DE
	c.HL, c.HL_ = c.HL_, c.HL

	c.PC++
	return 4
}

func (c *CPU) jpCNn() uint8 {
	if !c.getC() {
		c.PC += 3
		return 10
	}

	c.PC = c.readWord(c.PC + 1)
	return 10
}

func (c *CPU) inA_N_() uint8 {
	c.setAcc(c.getPort(c.dma.GetMemory(c.PC + 1)))

	c.PC += 2
	return 11
}

func (c *CPU) callCNn() uint8 {
	if !c.getC() {
		c.PC += 3
		return 10
	}
	c.pushStack(c.PC + 3)
	c.PC = c.readWord(c.PC + 1)

	return 17
}

func (c *CPU) sbcAN() uint8 {
	c.setC(!c.getC())
	c.adcValueToAcc(c.dma.GetMemory(c.PC+1) ^ 0xff)

	c.PC += 2
	c.setN(true)
	c.setC(!c.getC())
	c.setH(!c.getH())

	return 7
}

func (c *CPU) retPo() uint8 {
	if c.getPV() {
		c.PC++
		return 5
	}

	c.PC = c.popStack()

	return 11
}

func (c *CPU) popSs(ss string) func() uint8 {
	switch ss {
	case "HL":
		return func() uint8 {
			c.HL = c.popStack()
			c.PC++

			return 10
		}
	case "IX":
		return func() uint8 {
			c.IX = c.popStack()
			c.PC += 2

			return 14
		}
	case "IY":
		return func() uint8 {
			c.IY = c.popStack()
			c.PC += 2

			return 14
		}
	}

	panic("Invalid `ss` type")
}

func (c *CPU) jpPoNn() uint8 {
	if c.getPV() {
		c.PC += 3
		return 10
	}

	c.PC = c.readWord(c.PC + 1)
	return 10
}

func (c *CPU) ex_Sp_Ss(ss string) func() uint8 {
	switch ss {
	case "HL":
		return func() uint8 {
			value := c.readWord(c.SP)
			c.writeWord(c.SP, c.HL)
			c.HL = value

			c.PC++
			return 19
		}
	case "IX":
		return func() uint8 {
			value := c.readWord(c.SP)
			c.writeWord(c.SP, c.IX)
			c.IX = value

			c.PC += 2
			return 23
		}
	case "IY":
		return func() uint8 {
			value := c.readWord(c.SP)
			c.writeWord(c.SP, c.IY)
			c.IY = value

			c.PC += 2
			return 23
		}
	}

	panic("Invalid `ss` type")
}

func (c *CPU) callPoNn() uint8 {
	if c.getPV() {
		c.PC += 3
		return 10
	}
	c.pushStack(c.PC + 3)
	c.PC = c.readWord(c.PC + 1)

	return 17
}

func (c *CPU) pushSs(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			c.pushStack(c.HL)
			c.PC++

			return 11
		}
	}

	return func() uint8 {
		c.pushStack(c.extractRegisterPair(ss))
		c.PC += 2

		return 15
	}

}

func (c *CPU) andN() uint8 {
	result := c.getAcc() & c.dma.GetMemory(c.PC+1)

	c.PC++
	c.setAcc(result)
	c.setS(result > 127)
	c.setZ(result == 0)
	c.setH(true)
	c.setPV(parityTable[result])
	c.setN(false)
	c.setC(false)

	return 7
}

func (c *CPU) retPe() uint8 {
	if !c.getPV() {
		c.PC++
		return 5
	}

	c.PC = c.popStack()

	return 11
}

func (c *CPU) jp_Ss_(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			c.PC = c.readWord(c.HL)
			return 4
		}
	}

	return func() uint8 {
		c.PC = c.readWord(c.extractRegisterPair(ss))
		return 8
	}
}

func (c *CPU) jpPeNn() uint8 {
	if !c.getPV() {
		c.PC += 3
		return 10
	}

	c.PC = c.readWord(c.PC + 1)
	return 10
}

func (c *CPU) exDeSs(ss string) func() uint8 {
	switch ss {
	case "HL":
		return func() uint8 {
			c.DE, c.HL = c.HL, c.DE

			c.PC++
			return 4
		}
	case "IX":
		return func() uint8 {
			c.DE, c.IX = c.IX, c.DE

			c.PC += 2
			return 8
		}
	case "IY":
		return func() uint8 {
			c.DE, c.IY = c.IY, c.DE

			c.PC += 2
			return 8
		}
	}

	panic("Invalid `ss` type")
}

func (c *CPU) callPeNn() uint8 {
	if !c.getPV() {
		c.PC += 3
		return 10
	}
	c.pushStack(c.PC + 3)
	c.PC = c.readWord(c.PC + 1)

	return 17
}

func (c *CPU) xorN() uint8 {
	result := c.getAcc() ^ c.dma.GetMemory(c.PC+1)

	c.PC += 2
	c.setAcc(result)
	c.setS(result > 127)
	c.setZ(result == 0)
	c.setH(false)
	c.setPV(parityTable[result])
	c.setN(false)
	c.setC(false)

	return 7
}

func (c *CPU) retP() uint8 {
	if c.getS() {
		c.PC++
		return 5
	}

	c.PC = c.popStack()

	return 11
}

func (c *CPU) popAf() uint8 {
	c.AF = c.popStack()
	c.PC++

	return 10
}

func (c *CPU) jpPNn() uint8 {
	if c.getS() {
		c.PC += 3
		return 10
	}

	c.PC = c.readWord(c.PC + 1)
	return 10
}

func (c *CPU) di() uint8 {
	c.disableInterrupts()

	c.PC++
	return 4
}

func (c *CPU) callPNn() uint8 {
	if c.getS() {
		c.PC += 3
		return 10
	}
	c.pushStack(c.PC + 3)
	c.PC = c.readWord(c.PC + 1)

	return 17
}

func (c *CPU) pushAf() uint8 {
	c.pushStack(c.AF)
	c.PC++

	return 11
}

func (c *CPU) orN() uint8 {
	result := c.getAcc() | c.dma.GetMemory(c.PC+1)

	c.PC += 2
	c.setAcc(result)
	c.setS(result > 127)
	c.setZ(result == 0)
	c.setH(false)
	c.setPV(parityTable[result])
	c.setN(false)
	c.setC(false)

	return 7
}

func (c *CPU) retM() uint8 {
	if !c.getS() {
		c.PC++
		return 5
	}

	c.PC = c.popStack()

	return 11
}

func (c *CPU) ldSpSs(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			c.SP = c.HL

			c.PC++
			return 6
		}
	}

	return func() uint8 {
		c.SP = c.extractRegisterPair(ss)

		c.PC += 2
		return 10
	}
}

func (c *CPU) jpMNn() uint8 {
	if !c.getS() {
		c.PC += 3
		return 10
	}

	c.PC = c.readWord(c.PC + 1)
	return 10
}

func (c *CPU) ei() uint8 {
	c.enableInterrupts()

	c.PC++
	return 4
}

func (c *CPU) callMNn() uint8 {
	if !c.getS() {
		c.PC += 3
		return 10
	}
	c.pushStack(c.PC + 3)
	c.PC = c.readWord(c.PC + 1)

	return 17
}

func (c *CPU) cpN() uint8 {
	acc := c.getAcc()
	c.setC(true)
	c.adcValueToAcc(c.dma.GetMemory(c.PC+1) ^ 0xff)

	c.PC += 2
	c.setAcc(acc)
	c.setN(true)
	c.setC(!c.getC())
	c.setH(!c.getH())

	return 7
}

func (c *CPU) inR_C_(r byte) func() uint8 {
	return func() uint8 {
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

		result := c.getPort(uint8(c.BC))

		if r != ' ' {
			if lhigh {
				*lvalue = (*lvalue & 0x00ff) | (uint16(result) << 8)
			} else {
				*lvalue = (*lvalue & 0xff00) | uint16(result)
			}
		}

		c.setS(result > 127)
		c.setZ(result == 0)
		c.setH(false)
		c.setPV(parityTable[result])
		c.setN(false)

		c.PC += 2

		return 12
	}

}

func (c *CPU) out_C_R(r byte) func() uint8 {
	return func() uint8 {
		var right uint8

		if r == ' ' {
			right = 0
		} else {
			right = c.extractRegister(r)
		}

		c.setPort(uint8(c.BC), right)

		c.PC += 2
		return 12
	}
}

func (c *CPU) sbcHlRr(rr string) func() uint8 {
	return func() uint8 {
		c.setC(!c.getC())
		c.HL = c.adc16bit(c.HL, c.extractRegisterPair(rr)^0xffff)

		c.PC += 2
		c.setN(true)
		c.setC(!c.getC())
		c.setH(!c.getH())

		return 15
	}
}

func (c *CPU) adcHlRr(rr string) func() uint8 {
	return func() uint8 {
		c.HL = c.adc16bit(c.HL, c.extractRegisterPair(rr))

		c.PC += 2
		return 15
	}
}

func (c *CPU) ld_Nn_Rr(rr string) func() uint8 {
	return func() uint8 {
		c.writeWord(c.readWord(c.PC+2), c.extractRegisterPair(rr))

		c.PC += 4
		return 20
	}
}

func (c *CPU) neg() uint8 {
	a := c.getAcc()
	c.setAcc(0)

	c.setC(false)
	c.adcValueToAcc(a ^ 0xff)

	c.PC += 2
	c.setN(true)
	c.setC(!c.getC())
	c.setH(!c.getH())

	return 8
}

func (c *CPU) retn() uint8 {
	c.PC = c.popStack()
	c.States.IFF1 = c.States.IFF2

	return 14
}

func (c *CPU) reti() uint8 {
	c.PC = c.popStack()
	c.States.IFF1 = c.States.IFF2

	return 14
}

func (c *CPU) im(mode uint8) func() uint8 {
	return func() uint8 {
		c.States.IM = mode
		c.PC += 2

		return 8
	}
}

func (c *CPU) ldIA() uint8 {
	c.I = c.getAcc()

	c.PC += 2
	return 9
}

func (c *CPU) ldAI() uint8 {
	c.setAcc(c.I)

	c.setS(c.I > 127)
	c.setZ(c.I == 0)
	c.setH(false)
	c.setPV(c.States.IFF2)
	c.setN(false)

	c.PC += 2
	return 9
}

func (c *CPU) ldRr_Nn_(rr string) func() uint8 {
	return func() uint8 {
		var lvalue *uint16

		switch rr {
		case "AF":
			lvalue = &c.AF
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

		*lvalue = c.readWord(c.readWord(c.PC + 2))

		c.PC += 4
		return 20
	}
}

func (c *CPU) ldRA() uint8 {
	c.R = c.getAcc()

	c.PC += 2
	return 9
}

func (c *CPU) ldAR() uint8 {
	c.setAcc(c.R)

	c.setS(c.R > 127)
	c.setZ(c.R == 0)
	c.setH(false)
	c.setPV(c.States.IFF2)
	c.setN(false)

	c.PC += 2
	return 9
}

func (c *CPU) rrd() uint8 {
	value := c.dma.GetMemory(c.HL)
	a := c.getAcc()

	c.setAcc((a & 0xf0) | (value & 0x0f))
	value = value >> 4
	value = (a << 4) | value

	c.dma.SetMemoryByte(c.HL, value)
	a = c.getAcc()

	c.setS(a > 127)
	c.setZ(a == 0)
	c.setH(false)
	c.setPV(parityTable[a])
	c.setN(false)

	c.PC += 2
	return 18
}

func (c *CPU) rld() uint8 {
	value := c.dma.GetMemory(c.HL)
	a := c.getAcc()

	c.setAcc((a & 0xf0) | ((value >> 4) & 0x0f))
	value = value << 4
	value = (a & 0x0f) | value

	c.dma.SetMemoryByte(c.HL, value)
	a = c.getAcc()

	c.setS(a > 127)
	c.setZ(a == 0)
	c.setH(false)
	c.setPV(parityTable[a])
	c.setN(false)

	c.PC += 2
	return 18
}

func (c *CPU) ldi() uint8 {
	c.writeWord(c.DE, c.readWord(c.HL))
	c.DE++
	c.HL++
	c.BC--

	c.setH(false)
	c.setPV(c.BC != 0)
	c.setN(false)

	c.PC += 2
	return 16
}

func (c *CPU) cpi() uint8 {
	acc := c.getAcc()
	flagC := c.getC()
	c.setC(true)
	c.adcValueToAcc(c.dma.GetMemory(c.HL) ^ 0xff)
	c.HL++
	c.BC--

	c.setAcc(acc)
	c.setC(flagC)
	c.setN(true)
	c.setPV(c.BC != 0)
	c.setH(!c.getH())

	c.PC += 2
	return 16
}

func (c *CPU) ini() uint8 {
	c.dma.SetMemoryByte(c.HL, c.getPort(c.extractRegister('C')))
	c.HL++
	c.BC -= 256

	c.setZ(c.BC < 256)
	c.setN(true)

	c.PC += 2
	return 16
}

func (c *CPU) outi() uint8 {
	c.setPort(c.extractRegister('C'), c.dma.GetMemory(c.HL))
	c.HL++
	c.BC -= 256

	c.setZ(c.BC < 256)
	c.setN(true)

	c.PC += 2
	return 16
}

func (c *CPU) ldd() uint8 {
	c.writeWord(c.DE, c.readWord(c.HL))
	c.DE--
	c.HL--
	c.BC--

	c.setH(false)
	c.setPV(c.BC != 0)
	c.setN(false)

	c.PC += 2
	return 16
}

func (c *CPU) cpd() uint8 {
	acc := c.getAcc()
	flagC := c.getC()
	c.setC(true)
	c.adcValueToAcc(c.dma.GetMemory(c.HL) ^ 0xff)
	c.HL--
	c.BC--

	c.setAcc(acc)
	c.setC(flagC)
	c.setN(true)
	c.setPV(c.BC != 0)
	c.setH(!c.getH())

	c.PC += 2
	return 16
}

func (c *CPU) ind() uint8 {
	c.dma.SetMemoryByte(c.HL, c.getPort(c.extractRegister('C')))
	c.HL--
	c.BC -= 256

	c.setZ(c.BC < 256)
	c.setN(true)

	c.PC += 2
	return 16
}

func (c *CPU) outd() uint8 {
	c.setPort(c.extractRegister('C'), c.dma.GetMemory(c.HL))
	c.HL--
	c.BC -= 256

	c.setZ(c.BC < 256)
	c.setN(true)

	c.PC += 2
	return 16
}

func (c *CPU) ldir() uint8 {
	c.ldi()

	if c.BC == 0 {
		return 16
	}

	c.PC -= 2
	return 21
}

func (c *CPU) cpir() uint8 {
	acc := c.getAcc()
	flagC := c.getC()
	c.setC(true)
	c.adcValueToAcc(c.dma.GetMemory(c.HL) ^ 0xff)
	result := c.getAcc()
	c.HL++
	c.BC--

	c.setAcc(acc)
	c.setC(flagC)
	c.setN(true)
	c.setPV(c.BC != 0)
	c.setH(!c.getH())

	if c.BC == 0 || result == 0 {
		c.PC += 2
		return 16
	}

	return 21
}

func (c *CPU) inir() uint8 {
	c.ini()

	if c.extractRegister('B') == 0 {
		return 16
	}

	c.PC -= 2
	return 21
}

func (c *CPU) otir() uint8 {
	c.outi()

	if c.extractRegister('B') == 0 {
		return 16
	}

	c.PC -= 2
	return 21
}

func (c *CPU) lddr() uint8 {
	c.ldd()

	if c.BC == 0 {
		return 16
	}

	c.PC -= 2
	return 21
}

func (c *CPU) cpdr() uint8 {
	acc := c.getAcc()
	flagC := c.getC()
	c.setC(true)
	c.adcValueToAcc(c.dma.GetMemory(c.HL) ^ 0xff)
	result := c.getAcc()
	c.HL--
	c.BC--

	c.setAcc(acc)
	c.setC(flagC)
	c.setN(true)
	c.setPV(c.BC != 0)
	c.setH(!c.getH())

	if c.BC == 0 || result == 0 {
		c.PC += 2
		return 16
	}

	return 21
}

func (c *CPU) indr() uint8 {
	c.ind()

	if c.extractRegister('B') == 0 {
		return 16
	}

	c.PC -= 2
	return 21
}

func (c *CPU) otdr() uint8 {
	c.outd()

	if c.extractRegister('B') == 0 {
		return 16
	}

	c.PC -= 2
	return 21
}

func (c *CPU) rlcSs(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			rvalue := c.dma.GetMemory(c.HL)

			signed := rvalue&128 == 128
			rvalue = rvalue << 1
			c.PC += 2

			if signed {
				rvalue = rvalue | 0x01
			}

			c.dma.SetMemoryByte(c.HL, rvalue)

			c.setC(signed)
			c.setN(false)
			c.setPV(parityTable[rvalue])
			c.setH(false)
			c.setZ(rvalue == 0)
			c.setS(rvalue > 127)

			return 15
		}
	}

	return func() uint8 {
		address := c.extractRegisterPair(ss) + uint16(c.dma.GetMemory(c.PC+3))
		rvalue := c.dma.GetMemory(address)

		signed := rvalue&128 == 128
		rvalue = rvalue << 1
		c.PC += 4

		if signed {
			rvalue = rvalue | 0x01
		}

		c.dma.SetMemoryByte(address, rvalue)

		c.setC(signed)
		c.setN(false)
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setZ(rvalue == 0)
		c.setS(rvalue > 127)

		return 23
	}
}

func (c *CPU) rrcR(r byte) func() uint8 {
	return func() uint8 {
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
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setZ(rvalue == 0)
		c.setS(rvalue > 127)

		return 4 + size*4
	}
}

func (c *CPU) rrcSs(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			rvalue := c.dma.GetMemory(c.HL)

			signed := rvalue&1 == 1
			rvalue = rvalue >> 1
			c.PC += 2

			if signed {
				rvalue = rvalue | 0x80
			}

			c.dma.SetMemoryByte(c.HL, rvalue)

			c.setC(signed)
			c.setN(false)
			c.setPV(parityTable[rvalue])
			c.setH(false)
			c.setZ(rvalue == 0)
			c.setS(rvalue > 127)

			return 15
		}
	}

	return func() uint8 {
		address := c.extractRegisterPair(ss) + uint16(c.dma.GetMemory(c.PC+3))
		rvalue := c.dma.GetMemory(address)

		signed := rvalue&1 == 1
		rvalue = rvalue >> 1
		c.PC += 4

		if signed {
			rvalue = rvalue | 0x80
		}

		c.dma.SetMemoryByte(address, rvalue)

		c.setC(signed)
		c.setN(false)
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setZ(rvalue == 0)
		c.setS(rvalue > 127)

		return 23
	}
}

func (c *CPU) slaR(r byte) func() uint8 {
	return func() uint8 {
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

		return 8
	}
}

func (c *CPU) slaSs(ss string) func() uint8 {
	if ss == "HL" {
		return func() uint8 {
			rvalue := c.dma.GetMemory(c.HL)

			c.setC(rvalue&128 == 128)
			rvalue = rvalue << 1
			c.PC += 2

			c.dma.SetMemoryByte(c.HL, rvalue)

			c.setN(false)
			c.setPV(parityTable[rvalue])
			c.setH(false)
			c.setZ(rvalue == 0)
			c.setS(rvalue > 127)

			return 15
		}
	}

	return func() uint8 {
		address := c.extractRegisterPair(ss) + uint16(c.dma.GetMemory(c.PC+3))
		rvalue := c.dma.GetMemory(address)

		c.setC(rvalue&128 == 128)
		rvalue = rvalue << 1
		c.PC += 4

		c.dma.SetMemoryByte(address, rvalue)

		c.setN(false)
		c.setPV(parityTable[rvalue])
		c.setH(false)
		c.setZ(rvalue == 0)
		c.setS(rvalue > 127)

		return 23
	}
}

func (c *CPU) sraR(r byte) func() uint8 {
	return func() uint8 {
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

		return 8
	}
}

func (c *CPU) die() uint8 {
	panic("unimplemented mnemonic")
}

func (c *CPU) Step() uint8 {
	idx := c.dma.GetMemory(c.PC)
	a, b := c.dma.GetMemory(c.PC+1), c.dma.GetMemory(c.PC+2)

	if idx == 0xdd || idx == 0xed || idx == 0xfd {
		var cycles uint8
		groupIdx := idx
		idx = c.dma.GetMemory(c.PC + 1)
		a, b = c.dma.GetMemory(c.PC+2), c.dma.GetMemory(c.PC+3)
		if groupIdx == 0xdd {
			fmt.Printf("%04x: %s [DD %02x %02x %02x] -> ", c.PC, mnemonicsDebug.xxIXxx[idx], idx, a, b)
			cycles := c.mnemonics.xxIXxx[idx]()
			fmt.Printf("(%d) => A: %02x, F: %08b, BC: %04x, DE: %04x, HL: %04x, SP: %04x\n", cycles, c.getAcc(), c.getFlags(), c.BC, c.DE, c.HL, c.SP)
		}
		if groupIdx == 0xed {
			fmt.Printf("%04x: %s [ED %02x %02x %02x] -> ", c.PC, mnemonicsDebug.xx80xx[idx], idx, a, b)
			cycles := c.mnemonics.xx80xx[idx]()
			fmt.Printf("(%d) => A: %02x, F: %08b, BC: %04x, DE: %04x, HL: %04x, SP: %04x\n", cycles, c.getAcc(), c.getFlags(), c.BC, c.DE, c.HL, c.SP)
		}
		if groupIdx == 0xfd {
			fmt.Printf("%04x: %s [FD %02x %02x %02x] -> ", c.PC, mnemonicsDebug.xxIYxx[idx], idx, a, b)
			cycles := c.mnemonics.xxIYxx[idx]()
			fmt.Printf("(%d) => A: %02x, F: %08b, BC: %04x, DE: %04x, HL: %04x, SP: %04x\n", cycles, c.getAcc(), c.getFlags(), c.BC, c.DE, c.HL, c.SP)
		}
		return cycles
	}
	fmt.Printf("%04x: %s [%02x %02x %02x] -> ", c.PC, mnemonicsDebug.base[idx], idx, a, b)
	cycles := c.mnemonics.base[idx]()
	fmt.Printf("(%d) => A: %02x, F: %08b, BC: %04x, DE: %04x, HL: %04x, SP: %04x\n", cycles, c.getAcc(), c.getFlags(), c.BC, c.DE, c.HL, c.SP)
	return cycles
}

func (c *CPU) Reset() {
	c.PC = 0
	c.SP = 0
	c.AF = 0
	c.AF_ = 0
	c.BC = 0
	c.BC_ = 0
	c.DE = 0
	c.DE_ = 0
	c.HL = 0
	c.HL_ = 0
	c.I = 0
	c.R = 0
	c.IX = 0
	c.IY = 0
	c.States = CPUStates{IFF1: true, IFF2: true}
}

func CPUNew(dma *dma.DMA) *CPU {
	cpu := new(CPU)
	cpu.dma = dma
	cpu.initializeMnemonics()
	cpu.Reset()
	return cpu
}
