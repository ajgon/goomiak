package cpu

import (
	"fmt"
	"z80/dma"
	"z80/loader"
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
	base      [256]string
	xx80xx    [256]string
	xxIXxx    [256]string
	xxIYxx    [256]string
	xxBITxx   [256]string
	xxIXBITxx [256]string
	xxIYBITxx [256]string
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
	xxBITxx: [256]string{
		"rlc b", "rlc c", "rlc d", "rlc e", "rlc h", "rlc l", "rlc (hl)", "rlc a",
		"rrc b", "rrc c", "rrc d", "rrc e", "rrc h", "rrc l", "rrc (hl)", "rrc a",
		"rl b", "rl c", "rl d", "rl e", "rl h", "rl l", "rl (hl)", "rl a",
		"rr b", "rr c", "rr d", "rr e", "rr h", "rr l", "rr (hl)", "rr a",
		"sla b", "sla c", "sla d", "sla e", "sla h", "sla l", "sla (hl)", "sla a",
		"sra b", "sra c", "sra d", "sra e", "sra h", "sra l", "sra (hl)", "sra a",
		"sll b", "sll c", "sll d", "sll e", "sll h", "sll l", "sll (hl)", "sll a",
		"srl b", "srl c", "srl d", "srl e", "srl h", "srl l", "srl (hl)", "srl a",
		"bit 0,b", "bit 0,c", "bit 0,d", "bit 0,e", "bit 0,h", "bit 0,l", "bit 0,(hl)", "bit 0, a",
		"bit 1,b", "bit 1,c", "bit 1,d", "bit 1,e", "bit 1,h", "bit 1,l", "bit 1,(hl)", "bit 1, a",
		"bit 2,b", "bit 2,c", "bit 2,d", "bit 2,e", "bit 2,h", "bit 2,l", "bit 2,(hl)", "bit 2, a",
		"bit 3,b", "bit 3,c", "bit 3,d", "bit 3,e", "bit 3,h", "bit 3,l", "bit 3,(hl)", "bit 3, a",
		"bit 4,b", "bit 4,c", "bit 4,d", "bit 4,e", "bit 4,h", "bit 4,l", "bit 4,(hl)", "bit 4, a",
		"bit 5,b", "bit 5,c", "bit 5,d", "bit 5,e", "bit 5,h", "bit 5,l", "bit 5,(hl)", "bit 5, a",
		"bit 6,b", "bit 6,c", "bit 6,d", "bit 6,e", "bit 6,h", "bit 6,l", "bit 6,(hl)", "bit 6, a",
		"bit 7,b", "bit 7,c", "bit 7,d", "bit 7,e", "bit 7,h", "bit 7,l", "bit 7,(hl)", "bit 7, a",
		"res 0,b", "res 0,c", "res 0,d", "res 0,e", "res 0,h", "res 0,l", "res 0,(hl)", "res 0, a",
		"res 1,b", "res 1,c", "res 1,d", "res 1,e", "res 1,h", "res 1,l", "res 1,(hl)", "res 1, a",
		"res 2,b", "res 2,c", "res 2,d", "res 2,e", "res 2,h", "res 2,l", "res 2,(hl)", "res 2, a",
		"res 3,b", "res 3,c", "res 3,d", "res 3,e", "res 3,h", "res 3,l", "res 3,(hl)", "res 3, a",
		"res 4,b", "res 4,c", "res 4,d", "res 4,e", "res 4,h", "res 4,l", "res 4,(hl)", "res 4, a",
		"res 5,b", "res 5,c", "res 5,d", "res 5,e", "res 5,h", "res 5,l", "res 5,(hl)", "res 5, a",
		"res 6,b", "res 6,c", "res 6,d", "res 6,e", "res 6,h", "res 6,l", "res 6,(hl)", "res 6, a",
		"res 7,b", "res 7,c", "res 7,d", "res 7,e", "res 7,h", "res 7,l", "res 7,(hl)", "res 7, a",
		"set 0,b", "set 0,c", "set 0,d", "set 0,e", "set 0,h", "set 0,l", "set 0,(hl)", "set 0, a",
		"set 1,b", "set 1,c", "set 1,d", "set 1,e", "set 1,h", "set 1,l", "set 1,(hl)", "set 1, a",
		"set 2,b", "set 2,c", "set 2,d", "set 2,e", "set 2,h", "set 2,l", "set 2,(hl)", "set 2, a",
		"set 3,b", "set 3,c", "set 3,d", "set 3,e", "set 3,h", "set 3,l", "set 3,(hl)", "set 3, a",
		"set 4,b", "set 4,c", "set 4,d", "set 4,e", "set 4,h", "set 4,l", "set 4,(hl)", "set 4, a",
		"set 5,b", "set 5,c", "set 5,d", "set 5,e", "set 5,h", "set 5,l", "set 5,(hl)", "set 5, a",
		"set 6,b", "set 6,c", "set 6,d", "set 6,e", "set 6,h", "set 6,l", "set 6,(hl)", "set 6, a",
		"set 7,b", "set 7,c", "set 7,d", "set 7,e", "set 7,h", "set 7,l", "set 7,(hl)", "set 7, a",
	},
	xxIXBITxx: [256]string{
		"rlc b", "rlc c", "rlc d", "rlc e", "rlc h", "rlc l", "rlc (ix)", "rlc a",
		"rrc b", "rrc c", "rrc d", "rrc e", "rrc h", "rrc l", "rrc (ix)", "rrc a",
		"rl b", "rl c", "rl d", "rl e", "rl h", "rl l", "rl (ix)", "rl a",
		"rr b", "rr c", "rr d", "rr e", "rr h", "rr l", "rr (ix)", "rr a",
		"sla b", "sla c", "sla d", "sla e", "sla h", "sla l", "sla (ix)", "sla a",
		"sra b", "sra c", "sra d", "sra e", "sra h", "sra l", "sra (ix)", "sra a",
		"sll b", "sll c", "sll d", "sll e", "sll h", "sll l", "sll (ix)", "sll a",
		"srl b", "srl c", "srl d", "srl e", "srl h", "srl l", "srl (ix)", "srl a",
		"bit 0,b", "bit 0,c", "bit 0,d", "bit 0,e", "bit 0,h", "bit 0,l", "bit 0,(ix)", "bit 0, a",
		"bit 1,b", "bit 1,c", "bit 1,d", "bit 1,e", "bit 1,h", "bit 1,l", "bit 1,(ix)", "bit 1, a",
		"bit 2,b", "bit 2,c", "bit 2,d", "bit 2,e", "bit 2,h", "bit 2,l", "bit 2,(ix)", "bit 2, a",
		"bit 3,b", "bit 3,c", "bit 3,d", "bit 3,e", "bit 3,h", "bit 3,l", "bit 3,(ix)", "bit 3, a",
		"bit 4,b", "bit 4,c", "bit 4,d", "bit 4,e", "bit 4,h", "bit 4,l", "bit 4,(ix)", "bit 4, a",
		"bit 5,b", "bit 5,c", "bit 5,d", "bit 5,e", "bit 5,h", "bit 5,l", "bit 5,(ix)", "bit 5, a",
		"bit 6,b", "bit 6,c", "bit 6,d", "bit 6,e", "bit 6,h", "bit 6,l", "bit 6,(ix)", "bit 6, a",
		"bit 7,b", "bit 7,c", "bit 7,d", "bit 7,e", "bit 7,h", "bit 7,l", "bit 7,(ix)", "bit 7, a",
		"res 0,b", "res 0,c", "res 0,d", "res 0,e", "res 0,h", "res 0,l", "res 0,(ix)", "res 0, a",
		"res 1,b", "res 1,c", "res 1,d", "res 1,e", "res 1,h", "res 1,l", "res 1,(ix)", "res 1, a",
		"res 2,b", "res 2,c", "res 2,d", "res 2,e", "res 2,h", "res 2,l", "res 2,(ix)", "res 2, a",
		"res 3,b", "res 3,c", "res 3,d", "res 3,e", "res 3,h", "res 3,l", "res 3,(ix)", "res 3, a",
		"res 4,b", "res 4,c", "res 4,d", "res 4,e", "res 4,h", "res 4,l", "res 4,(ix)", "res 4, a",
		"res 5,b", "res 5,c", "res 5,d", "res 5,e", "res 5,h", "res 5,l", "res 5,(ix)", "res 5, a",
		"res 6,b", "res 6,c", "res 6,d", "res 6,e", "res 6,h", "res 6,l", "res 6,(ix)", "res 6, a",
		"res 7,b", "res 7,c", "res 7,d", "res 7,e", "res 7,h", "res 7,l", "res 7,(ix)", "res 7, a",
		"set 0,b", "set 0,c", "set 0,d", "set 0,e", "set 0,h", "set 0,l", "set 0,(ix)", "set 0, a",
		"set 1,b", "set 1,c", "set 1,d", "set 1,e", "set 1,h", "set 1,l", "set 1,(ix)", "set 1, a",
		"set 2,b", "set 2,c", "set 2,d", "set 2,e", "set 2,h", "set 2,l", "set 2,(ix)", "set 2, a",
		"set 3,b", "set 3,c", "set 3,d", "set 3,e", "set 3,h", "set 3,l", "set 3,(ix)", "set 3, a",
		"set 4,b", "set 4,c", "set 4,d", "set 4,e", "set 4,h", "set 4,l", "set 4,(ix)", "set 4, a",
		"set 5,b", "set 5,c", "set 5,d", "set 5,e", "set 5,h", "set 5,l", "set 5,(ix)", "set 5, a",
		"set 6,b", "set 6,c", "set 6,d", "set 6,e", "set 6,h", "set 6,l", "set 6,(ix)", "set 6, a",
		"set 7,b", "set 7,c", "set 7,d", "set 7,e", "set 7,h", "set 7,l", "set 7,(ix)", "set 7, a",
	},
	xxIYBITxx: [256]string{
		"rlc b", "rlc c", "rlc d", "rlc e", "rlc h", "rlc l", "rlc (iy)", "rlc a",
		"rrc b", "rrc c", "rrc d", "rrc e", "rrc h", "rrc l", "rrc (iy)", "rrc a",
		"rl b", "rl c", "rl d", "rl e", "rl h", "rl l", "rl (iy)", "rl a",
		"rr b", "rr c", "rr d", "rr e", "rr h", "rr l", "rr (iy)", "rr a",
		"sla b", "sla c", "sla d", "sla e", "sla h", "sla l", "sla (iy)", "sla a",
		"sra b", "sra c", "sra d", "sra e", "sra h", "sra l", "sra (iy)", "sra a",
		"sll b", "sll c", "sll d", "sll e", "sll h", "sll l", "sll (iy)", "sll a",
		"srl b", "srl c", "srl d", "srl e", "srl h", "srl l", "srl (iy)", "srl a",
		"bit 0,b", "bit 0,c", "bit 0,d", "bit 0,e", "bit 0,h", "bit 0,l", "bit 0,(iy)", "bit 0, a",
		"bit 1,b", "bit 1,c", "bit 1,d", "bit 1,e", "bit 1,h", "bit 1,l", "bit 1,(iy)", "bit 1, a",
		"bit 2,b", "bit 2,c", "bit 2,d", "bit 2,e", "bit 2,h", "bit 2,l", "bit 2,(iy)", "bit 2, a",
		"bit 3,b", "bit 3,c", "bit 3,d", "bit 3,e", "bit 3,h", "bit 3,l", "bit 3,(iy)", "bit 3, a",
		"bit 4,b", "bit 4,c", "bit 4,d", "bit 4,e", "bit 4,h", "bit 4,l", "bit 4,(iy)", "bit 4, a",
		"bit 5,b", "bit 5,c", "bit 5,d", "bit 5,e", "bit 5,h", "bit 5,l", "bit 5,(iy)", "bit 5, a",
		"bit 6,b", "bit 6,c", "bit 6,d", "bit 6,e", "bit 6,h", "bit 6,l", "bit 6,(iy)", "bit 6, a",
		"bit 7,b", "bit 7,c", "bit 7,d", "bit 7,e", "bit 7,h", "bit 7,l", "bit 7,(iy)", "bit 7, a",
		"res 0,b", "res 0,c", "res 0,d", "res 0,e", "res 0,h", "res 0,l", "res 0,(iy)", "res 0, a",
		"res 1,b", "res 1,c", "res 1,d", "res 1,e", "res 1,h", "res 1,l", "res 1,(iy)", "res 1, a",
		"res 2,b", "res 2,c", "res 2,d", "res 2,e", "res 2,h", "res 2,l", "res 2,(iy)", "res 2, a",
		"res 3,b", "res 3,c", "res 3,d", "res 3,e", "res 3,h", "res 3,l", "res 3,(iy)", "res 3, a",
		"res 4,b", "res 4,c", "res 4,d", "res 4,e", "res 4,h", "res 4,l", "res 4,(iy)", "res 4, a",
		"res 5,b", "res 5,c", "res 5,d", "res 5,e", "res 5,h", "res 5,l", "res 5,(iy)", "res 5, a",
		"res 6,b", "res 6,c", "res 6,d", "res 6,e", "res 6,h", "res 6,l", "res 6,(iy)", "res 6, a",
		"res 7,b", "res 7,c", "res 7,d", "res 7,e", "res 7,h", "res 7,l", "res 7,(iy)", "res 7, a",
		"set 0,b", "set 0,c", "set 0,d", "set 0,e", "set 0,h", "set 0,l", "set 0,(iy)", "set 0, a",
		"set 1,b", "set 1,c", "set 1,d", "set 1,e", "set 1,h", "set 1,l", "set 1,(iy)", "set 1, a",
		"set 2,b", "set 2,c", "set 2,d", "set 2,e", "set 2,h", "set 2,l", "set 2,(iy)", "set 2, a",
		"set 3,b", "set 3,c", "set 3,d", "set 3,e", "set 3,h", "set 3,l", "set 3,(iy)", "set 3, a",
		"set 4,b", "set 4,c", "set 4,d", "set 4,e", "set 4,h", "set 4,l", "set 4,(iy)", "set 4, a",
		"set 5,b", "set 5,c", "set 5,d", "set 5,e", "set 5,h", "set 5,l", "set 5,(iy)", "set 5, a",
		"set 6,b", "set 6,c", "set 6,d", "set 6,e", "set 6,h", "set 6,l", "set 6,(iy)", "set 6, a",
		"set 7,b", "set 7,c", "set 7,d", "set 7,e", "set 7,h", "set 7,l", "set 7,(iy)", "set 7, a",
	},
}

type CPUConfig struct {
	ContentionDelays []uint8
	FrameLength      uint64
}

type CPUStates struct {
	Halt  bool
	Ports [256]uint8
	IFF1  bool
	IFF2  bool
	IM    uint8
	IRQ   bool
	Tmp   uint64
}

type CPUMnemonics struct {
	base      [256]func()
	xx80xx    [256]func()
	xxIXxx    [256]func()
	xxIYxx    [256]func()
	xxBITxx   [256]func()
	xxIXBITxx [256]func()
	xxIYBITxx [256]func()
}

type CPU struct {
	PC     uint16
	SP     uint16
	WZ     uint16
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

	config    CPUConfig
	dma       *dma.DMA
	mnemonics CPUMnemonics
	tstates   uint64
}

func (c *CPU) initializeMnemonics() {
	for _, reg := range [3]string{"HL", "IX", "IY"} {
		var highReg, lowReg byte
		switch reg {
		case "HL":
			highReg, lowReg = 'H', 'L'
		case "IX":
			highReg, lowReg = 'X', 'x'
		case "IY":
			highReg, lowReg = 'Y', 'y'
		}

		baseList := [256]func(){
			c.nop, c.ldSsNn("BC"), c.ld_Bc_A, c.incSs("BC"), c.incR('B'), c.decR('B'), c.ldRN('B'), c.rlcR(' '),
			c.exAfAf_, c.addSsRr(reg, "BC"), c.ldA_Bc_, c.decSs("BC"), c.incR('C'), c.decR('C'), c.ldRN('C'), c.rrcR(' '),
			c.djnzN, c.ldSsNn("DE"), c.ld_De_A, c.incSs("DE"), c.incR('D'), c.decR('D'), c.ldRN('D'), c.rlR(' '),
			c.jrN, c.addSsRr(reg, "DE"), c.ldA_De_, c.decSs("DE"), c.incR('E'), c.decR('E'), c.ldRN('E'), c.rrR(' '),
			c.jrNzN, c.ldSsNn(reg), c.ld_Nn_Ss(reg), c.incSs(reg), c.incR(highReg), c.decR(highReg), c.ldRN(highReg), c.daa,
			c.jrZN, c.addSsRr(reg, reg), c.ldSs_Nn_(reg), c.decSs(reg), c.incR(lowReg), c.decR(lowReg), c.ldRN(lowReg), c.cpl,
			c.jrNcN, c.ldSsNn("SP"), c.ld_Nn_A, c.incSs("SP"), c.inc_Ss_(reg), c.dec_Ss_(reg), c.ld_Ss_N(reg), c.scf,
			c.jrCN, c.addSsRr(reg, "SP"), c.ldA_Nn_, c.decSs("SP"), c.incR('A'), c.decR('A'), c.ldRN('A'), c.ccf,
			c.ldRR_('B', 'B'), c.ldRR_('B', 'C'), c.ldRR_('B', 'D'), c.ldRR_('B', 'E'), c.ldRR_('B', highReg), c.ldRR_('B', lowReg), c.ldR_Ss_('B', reg), c.ldRR_('B', 'A'),
			c.ldRR_('C', 'B'), c.ldRR_('C', 'C'), c.ldRR_('C', 'D'), c.ldRR_('C', 'E'), c.ldRR_('C', highReg), c.ldRR_('C', lowReg), c.ldR_Ss_('C', reg), c.ldRR_('C', 'A'),
			c.ldRR_('D', 'B'), c.ldRR_('D', 'C'), c.ldRR_('D', 'D'), c.ldRR_('D', 'E'), c.ldRR_('D', highReg), c.ldRR_('D', lowReg), c.ldR_Ss_('D', reg), c.ldRR_('D', 'A'),
			c.ldRR_('E', 'B'), c.ldRR_('E', 'C'), c.ldRR_('E', 'D'), c.ldRR_('E', 'E'), c.ldRR_('E', highReg), c.ldRR_('E', lowReg), c.ldR_Ss_('E', reg), c.ldRR_('E', 'A'),
			c.ldRR_(highReg, 'B'), c.ldRR_(highReg, 'C'), c.ldRR_(highReg, 'D'), c.ldRR_(highReg, 'E'), c.ldRR_(highReg, highReg), c.ldRR_(highReg, lowReg), c.ldR_Ss_('H', reg), c.ldRR_(highReg, 'A'),
			c.ldRR_(lowReg, 'B'), c.ldRR_(lowReg, 'C'), c.ldRR_(lowReg, 'D'), c.ldRR_(lowReg, 'E'), c.ldRR_(lowReg, highReg), c.ldRR_(lowReg, lowReg), c.ldR_Ss_('L', reg), c.ldRR_(lowReg, 'A'),
			c.ld_Ss_R(reg, 'B'), c.ld_Ss_R(reg, 'C'), c.ld_Ss_R(reg, 'D'), c.ld_Ss_R(reg, 'E'), c.ld_Ss_R(reg, 'H'), c.ld_Ss_R(reg, 'L'), c.halt, c.ld_Ss_R(reg, 'A'),
			c.ldRR_('A', 'B'), c.ldRR_('A', 'C'), c.ldRR_('A', 'D'), c.ldRR_('A', 'E'), c.ldRR_('A', highReg), c.ldRR_('A', lowReg), c.ldR_Ss_('A', reg), c.ldRR_('A', 'A'),
			c.addAR('B'), c.addAR('C'), c.addAR('D'), c.addAR('E'), c.addAR(highReg), c.addAR(lowReg), c.addA_Ss_(reg), c.addAR('A'),
			c.adcAR('B'), c.adcAR('C'), c.adcAR('D'), c.adcAR('E'), c.adcAR(highReg), c.adcAR(lowReg), c.adcA_Ss_(reg), c.adcAR('A'),
			c.subR('B'), c.subR('C'), c.subR('D'), c.subR('E'), c.subR(highReg), c.subR(lowReg), c.sub_Ss_(reg), c.subR('A'),
			c.sbcAR('B'), c.sbcAR('C'), c.sbcAR('D'), c.sbcAR('E'), c.sbcAR(highReg), c.sbcAR(lowReg), c.sbcA_Ss_(reg), c.sbcAR('A'),
			c.andR('B'), c.andR('C'), c.andR('D'), c.andR('E'), c.andR(highReg), c.andR(lowReg), c.and_Ss_(reg), c.andR('A'),
			c.xorR('B'), c.xorR('C'), c.xorR('D'), c.xorR('E'), c.xorR(highReg), c.xorR(lowReg), c.xor_Ss_(reg), c.xorR('A'),
			c.orR('B'), c.orR('C'), c.orR('D'), c.orR('E'), c.orR(highReg), c.orR(lowReg), c.or_Ss_(reg), c.orR('A'),
			c.cpR('B'), c.cpR('C'), c.cpR('D'), c.cpR('E'), c.cpR(highReg), c.cpR(lowReg), c.cp_Ss_(reg), c.cpR('A'),
			c.retNz, c.popSs("BC"), c.jpNzNn, c.jpNn, c.callNzNn, c.pushSs("BC"), c.addAN, c.rst(0x00),
			c.retZ, c.ret, c.jpZNn, c.nop, c.callZNn, c.callNn, c.adcAN, c.rst(0x08),
			c.retNc, c.popSs("DE"), c.jpNcNn, c.out_N_A, c.callNcNn, c.pushSs("DE"), c.subN, c.rst(0x10),
			c.retC, c.exx, c.jpCNn, c.inA_N_, c.callCNn, c.nop, c.sbcAN, c.rst(0x18),
			c.retPo, c.popSs(reg), c.jpPoNn, c.ex_Sp_Ss(reg), c.callPoNn, c.pushSs(reg), c.andN, c.rst(0x20),
			c.retPe, c.jp_Ss_(reg), c.jpPeNn, c.exDeSs(reg), c.callPeNn, c.nop, c.xorN, c.rst(0x28),
			c.retP, c.popSs("AF"), c.jpPNn, c.di, c.callPNn, c.pushSs("AF"), c.orN, c.rst(0x30),
			c.retM, c.ldSpSs(reg), c.jpMNn, c.ei, c.callMNn, c.nop, c.cpN, c.rst(0x38),
		}
		bitList := [256]func(){
			c.rlcR('B'), c.rlcR('C'), c.rlcR('D'), c.rlcR('E'), c.rlcR('H'), c.rlcR('L'), c.rlcSs(reg), c.rlcR('A'),
			c.rrcR('B'), c.rrcR('C'), c.rrcR('D'), c.rrcR('E'), c.rrcR('H'), c.rrcR('L'), c.rrcSs(reg), c.rrcR('A'),
			c.rlR('B'), c.rlR('C'), c.rlR('D'), c.rlR('E'), c.rlR('H'), c.rlR('L'), c.rlSs(reg), c.rlR('A'),
			c.rrR('B'), c.rrR('C'), c.rrR('D'), c.rrR('E'), c.rrR('H'), c.rrR('L'), c.rrSs(reg), c.rrR('A'),
			c.slaR('B'), c.slaR('C'), c.slaR('D'), c.slaR('E'), c.slaR('H'), c.slaR('L'), c.slaSs(reg), c.slaR('A'),
			c.sraR('B'), c.sraR('C'), c.sraR('D'), c.sraR('E'), c.sraR('H'), c.sraR('L'), c.sraSs(reg), c.sraR('A'),
			c.sllR('B'), c.sllR('C'), c.sllR('D'), c.sllR('E'), c.sllR('H'), c.sllR('L'), c.sllSs(reg), c.sllR('A'),
			c.srlR('B'), c.srlR('C'), c.srlR('D'), c.srlR('E'), c.srlR('H'), c.srlR('L'), c.srlSs(reg), c.srlR('A'),
			c.bitBR(0, 'B'), c.bitBR(0, 'C'), c.bitBR(0, 'D'), c.bitBR(0, 'E'), c.bitBR(0, highReg), c.bitBR(0, lowReg), c.bitBSs(0, reg), c.bitBR(0, 'A'),
			c.bitBR(1, 'B'), c.bitBR(1, 'C'), c.bitBR(1, 'D'), c.bitBR(1, 'E'), c.bitBR(1, highReg), c.bitBR(1, lowReg), c.bitBSs(1, reg), c.bitBR(1, 'A'),
			c.bitBR(2, 'B'), c.bitBR(2, 'C'), c.bitBR(2, 'D'), c.bitBR(2, 'E'), c.bitBR(2, highReg), c.bitBR(2, lowReg), c.bitBSs(2, reg), c.bitBR(2, 'A'),
			c.bitBR(3, 'B'), c.bitBR(3, 'C'), c.bitBR(3, 'D'), c.bitBR(3, 'E'), c.bitBR(3, highReg), c.bitBR(3, lowReg), c.bitBSs(3, reg), c.bitBR(3, 'A'),
			c.bitBR(4, 'B'), c.bitBR(4, 'C'), c.bitBR(4, 'D'), c.bitBR(4, 'E'), c.bitBR(4, highReg), c.bitBR(4, lowReg), c.bitBSs(4, reg), c.bitBR(4, 'A'),
			c.bitBR(5, 'B'), c.bitBR(5, 'C'), c.bitBR(5, 'D'), c.bitBR(5, 'E'), c.bitBR(5, highReg), c.bitBR(5, lowReg), c.bitBSs(5, reg), c.bitBR(5, 'A'),
			c.bitBR(6, 'B'), c.bitBR(6, 'C'), c.bitBR(6, 'D'), c.bitBR(6, 'E'), c.bitBR(6, highReg), c.bitBR(6, lowReg), c.bitBSs(6, reg), c.bitBR(6, 'A'),
			c.bitBR(7, 'B'), c.bitBR(7, 'C'), c.bitBR(7, 'D'), c.bitBR(7, 'E'), c.bitBR(7, highReg), c.bitBR(7, lowReg), c.bitBSs(7, reg), c.bitBR(7, 'A'),
			c.resBR(0, 'B'), c.resBR(0, 'C'), c.resBR(0, 'D'), c.resBR(0, 'E'), c.resBR(0, 'H'), c.resBR(0, 'L'), c.resBSs(0, reg), c.resBR(0, 'A'),
			c.resBR(1, 'B'), c.resBR(1, 'C'), c.resBR(1, 'D'), c.resBR(1, 'E'), c.resBR(1, 'H'), c.resBR(1, 'L'), c.resBSs(1, reg), c.resBR(1, 'A'),
			c.resBR(2, 'B'), c.resBR(2, 'C'), c.resBR(2, 'D'), c.resBR(2, 'E'), c.resBR(2, 'H'), c.resBR(2, 'L'), c.resBSs(2, reg), c.resBR(2, 'A'),
			c.resBR(3, 'B'), c.resBR(3, 'C'), c.resBR(3, 'D'), c.resBR(3, 'E'), c.resBR(3, 'H'), c.resBR(3, 'L'), c.resBSs(3, reg), c.resBR(3, 'A'),
			c.resBR(4, 'B'), c.resBR(4, 'C'), c.resBR(4, 'D'), c.resBR(4, 'E'), c.resBR(4, 'H'), c.resBR(4, 'L'), c.resBSs(4, reg), c.resBR(4, 'A'),
			c.resBR(5, 'B'), c.resBR(5, 'C'), c.resBR(5, 'D'), c.resBR(5, 'E'), c.resBR(5, 'H'), c.resBR(5, 'L'), c.resBSs(5, reg), c.resBR(5, 'A'),
			c.resBR(6, 'B'), c.resBR(6, 'C'), c.resBR(6, 'D'), c.resBR(6, 'E'), c.resBR(6, 'H'), c.resBR(6, 'L'), c.resBSs(6, reg), c.resBR(6, 'A'),
			c.resBR(7, 'B'), c.resBR(7, 'C'), c.resBR(7, 'D'), c.resBR(7, 'E'), c.resBR(7, 'H'), c.resBR(7, 'L'), c.resBSs(7, reg), c.resBR(7, 'A'),
			c.setBR(0, 'B'), c.setBR(0, 'C'), c.setBR(0, 'D'), c.setBR(0, 'E'), c.setBR(0, 'H'), c.setBR(0, 'L'), c.setBSs(0, reg), c.setBR(0, 'A'),
			c.setBR(1, 'B'), c.setBR(1, 'C'), c.setBR(1, 'D'), c.setBR(1, 'E'), c.setBR(1, 'H'), c.setBR(1, 'L'), c.setBSs(1, reg), c.setBR(1, 'A'),
			c.setBR(2, 'B'), c.setBR(2, 'C'), c.setBR(2, 'D'), c.setBR(2, 'E'), c.setBR(2, 'H'), c.setBR(2, 'L'), c.setBSs(2, reg), c.setBR(2, 'A'),
			c.setBR(3, 'B'), c.setBR(3, 'C'), c.setBR(3, 'D'), c.setBR(3, 'E'), c.setBR(3, 'H'), c.setBR(3, 'L'), c.setBSs(3, reg), c.setBR(3, 'A'),
			c.setBR(4, 'B'), c.setBR(4, 'C'), c.setBR(4, 'D'), c.setBR(4, 'E'), c.setBR(4, 'H'), c.setBR(4, 'L'), c.setBSs(4, reg), c.setBR(4, 'A'),
			c.setBR(5, 'B'), c.setBR(5, 'C'), c.setBR(5, 'D'), c.setBR(5, 'E'), c.setBR(5, 'H'), c.setBR(5, 'L'), c.setBSs(5, reg), c.setBR(5, 'A'),
			c.setBR(6, 'B'), c.setBR(6, 'C'), c.setBR(6, 'D'), c.setBR(6, 'E'), c.setBR(6, 'H'), c.setBR(6, 'L'), c.setBSs(6, reg), c.setBR(6, 'A'),
			c.setBR(7, 'B'), c.setBR(7, 'C'), c.setBR(7, 'D'), c.setBR(7, 'E'), c.setBR(7, 'H'), c.setBR(7, 'L'), c.setBSs(7, reg), c.setBR(7, 'A'),
		}
		switch reg {
		case "HL":
			c.mnemonics.base = baseList
			c.mnemonics.xxBITxx = bitList
		case "IX":
			c.mnemonics.xxIXxx = baseList
			c.mnemonics.xxIXBITxx = bitList
		case "IY":
			c.mnemonics.xxIYxx = baseList
			c.mnemonics.xxIYBITxx = bitList
		}
	}
	c.mnemonics.xx80xx = [256]func(){
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop, c.nop,
		c.inR_C_('B'), c.out_C_R('B'), c.sbcHlRr("BC"), c.ld_Nn_Rr("BC"), c.neg, c.retn, c.im(0), c.ldIA,
		c.inR_C_('C'), c.out_C_R('C'), c.adcHlRr("BC"), c.ldRr_Nn_("BC"), c.neg, c.reti, c.im(0), c.ldRA,
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

func (c *CPU) setF5(value bool) {
	if value {
		c.AF = c.AF | 0x0020
	} else {
		c.AF = c.AF & 0xffdf
	}
}

func (c *CPU) setH(value bool) {
	if value {
		c.AF = c.AF | 0x0010
	} else {
		c.AF = c.AF & 0xffef
	}
}

func (c *CPU) setF3(value bool) {
	if value {
		c.AF = c.AF | 0x0008
	} else {
		c.AF = c.AF & 0xfff7
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
	value = c.readWord(c.SP, 3, 3)
	c.SP += 2

	return
}

func (c *CPU) pushStack(value uint16) {
	c.SP -= 2
	c.writeWord(c.SP, value, 3, 3)
}

func (c *CPU) getPort(highAddrHalf, lowAddrHalf uint8, tstates uint8) uint8 {
	c.tstates += uint64(tstates)
	//if c.States.Tmp > 200 {
	//c.States.Tmp = 0
	//fmt.Printf("PORT address %02x%02x\n", highAddrHalf, lowAddrHalf)
	//if highAddrHalf&4 == 4 {
	//fmt.Printf("%02x-%02x\n", highAddrHalf, lowAddrHalf)
	if highAddrHalf&0x08 == 0 {
		if c.States.Tmp%200 < 10 {
			c.States.Tmp++
			return 0x17
		}
		//fmt.Println(c.States.Tmp)
		c.States.Tmp++
	}
	//return 0x00
	//}
	//}
	//c.States.Tmp++
	return 0xff
}

func (c *CPU) setPort(addr uint8, value uint8, tstates uint8) {
	c.tstates += uint64(tstates)
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

func (c *CPU) shiftedAddress(base uint16, shift uint8) uint16 {
	if shift > 127 {
		c.WZ = base + uint16(shift) - 256
	} else {
		c.WZ = base + uint16(shift)
	}

	return c.WZ
}

func (c *CPU) readByte(address uint16, usedTstates uint8) uint8 {
	value, contended := c.dma.GetMemoryByte(address)

	if contended {
		c.tstates += uint64(c.config.ContentionDelays[c.tstates%c.config.FrameLength])
	}

	c.tstates += uint64(usedTstates)

	return value
}

func (c *CPU) writeByte(address uint16, value uint8, usedTstates uint8) {
	contended := c.dma.SetMemoryByte(address, value)

	if contended && usedTstates > 0 {
		c.tstates += uint64(c.config.ContentionDelays[c.tstates%c.config.FrameLength])
	}

	c.tstates += uint64(usedTstates)
}

// reads word and maintains endianess
// example:
// 0040 34 21
// readWord(0x0040) => 0x1234
func (c *CPU) readWord(address uint16, usedTstates1, usedTstates2 uint8) uint16 {
	return uint16(c.readByte(address+1, usedTstates1))<<8 | uint16(c.readByte(address, usedTstates2))
}

// writes word to given address and address+1 and maintains endianess
// example:
// writeWord(0x1234, 0x5678)
// 1234  78 56
func (c *CPU) writeWord(address uint16, value uint16, usedTstates1, usedTstates2 uint8) {
	c.writeByte(address, uint8(value), usedTstates1)
	c.writeByte(address+1, uint8(value>>8), usedTstates2)
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
	case 'X':
		return uint8(c.IX >> 8)
	case 'x':
		return uint8(c.IX)
	case 'Y':
		return uint8(c.IY >> 8)
	case 'y':
		return uint8(c.IY)
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

// left stores the result
// @todo replace with adc16bit?
func (c *CPU) addRegisters(left *uint16, right uint16) {
	sum := *left + right

	c.setC(sum < *left || sum < right)
	c.setN(false)
	c.setH((*left^right^sum)&0x1000 == 0x1000)
	c.setF5(sum&0x2000 == 0x2000)
	c.setF3(sum&0x0800 == 0x0800)

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
	c.setF5(result&0x20 == 0x20)
	c.setF3(result&0x08 == 0x08)
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
	c.setF5(result&0x2000 == 0x2000)
	c.setF3(result&0x0800 == 0x0800)

	return
}

func (c *CPU) die() uint8 {
	panic("unimplemented mnemonic")
}

func (c *CPU) DebugStep() (tstates uint8) {
	var debugOpcode string

	//pc := c.PC
	debugT := c.tstates % 69888

	opcode := c.readByte(c.PC, 4)
	dbOpcode := opcode

	if dbOpcode == 0xcb || dbOpcode == 0xdd || dbOpcode == 0xed || dbOpcode == 0xfd {
		dbOpcode = c.readByte(c.PC+1, 0)
		if dbOpcode == 0xcb {
			dbOpcode = c.readByte(c.PC+3, 0)
		}
	}

	fmt.Printf(
		"%x: AF=%d BC=%d DE=%d HL=%d AF_=%d BC_=%d DE_=%d HL_=%d IX=%d IY=%d SP=%d PC=%d (HL)=%d t=%d\n",
		dbOpcode, c.AF, c.BC, c.DE, c.HL, c.AF_, c.BC_, c.DE_, c.HL_, c.IX, c.IY, c.SP, c.PC, c.readByte(c.HL, 0), debugT,
	)

	switch opcode {
	case 0xcb:
		opcode = c.readByte(c.PC+1, 4)
		c.mnemonics.xxBITxx[opcode]()
		debugOpcode = fmt.Sprintf("%s (CB %02x)", mnemonicsDebug.xxBITxx[opcode], opcode)
	case 0xdd:
		opcode = c.readByte(c.PC+1, 4)
		switch opcode {
		case 0xcb:
			opcode = c.readByte(c.PC+3, 3)
			c.mnemonics.xxIXBITxx[opcode]()
			debugOpcode = fmt.Sprintf("%s (DD CB %02x)", mnemonicsDebug.xxIXBITxx[opcode], opcode)
		default:
			c.mnemonics.xxIXxx[opcode]()
			debugOpcode = fmt.Sprintf("%s (DD %02x)", mnemonicsDebug.xxIXxx[opcode], opcode)
		}
	case 0xed:
		opcode = c.readByte(c.PC+1, 4)
		c.mnemonics.xx80xx[opcode]()
		debugOpcode = fmt.Sprintf("%s (ED %02x)", mnemonicsDebug.xx80xx[opcode], opcode)
	case 0xfd:
		opcode = c.readByte(c.PC+1, 4)
		switch opcode {
		case 0xcb:
			opcode = c.readByte(c.PC+3, 3)
			c.mnemonics.xxIYBITxx[opcode]()
			debugOpcode = fmt.Sprintf("%s (FD CB %02x)", mnemonicsDebug.xxIYBITxx[opcode], opcode)
		default:
			c.mnemonics.xxIYxx[opcode]()
			debugOpcode = fmt.Sprintf("%s (FD %02x)", mnemonicsDebug.xxIYxx[opcode], opcode)
		}
	default:
		c.mnemonics.base[opcode]()
		debugOpcode = fmt.Sprintf("%s (%02x)", mnemonicsDebug.base[opcode], opcode)
	}

	//47: AF=68 BC=0 DE=65535 HL=0 AF_=0 BC_=0 DE_=0 HL_=0 IX=0 IY=0 SP=0 PC=4555 (HL)=243 t=43

	if debugOpcode == "X" {
		return
	}

	//fmt.Printf("%04x: %s [%d]\n", pc, debugOpcode, c.tstates)
	return
}

func (c *CPU) Step() {
	opcode := c.readByte(c.PC, 4)
	switch opcode {
	case 0xcb:
		opcode = c.readByte(c.PC+1, 4)
		c.mnemonics.xxBITxx[opcode]()
	case 0xdd:
		opcode = c.readByte(c.PC+1, 4)
		switch opcode {
		case 0xcb:
			opcode = c.readByte(c.PC+3, 3)
			c.mnemonics.xxIXBITxx[opcode]()
		default:
			c.mnemonics.xxIXxx[opcode]()
		}
	case 0xed:
		opcode = c.readByte(c.PC+1, 4)
		c.mnemonics.xx80xx[opcode]()
	case 0xfd:
		opcode = c.readByte(c.PC+1, 4)
		switch opcode {
		case 0xcb:
			opcode = c.readByte(c.PC+3, 3)
			c.mnemonics.xxIYBITxx[opcode]()
		default:
			c.mnemonics.xxIYxx[opcode]()
		}
	default:
		c.mnemonics.base[opcode]()
	}

	return
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
	c.tstates = 0
	c.States = CPUStates{IFF1: true, IFF2: true}
}

func NewCPU(dma *dma.DMA, config CPUConfig) *CPU {
	cpu := new(CPU)
	cpu.dma = dma
	cpu.config = config

	cpu.initializeMnemonics()
	cpu.Reset()

	return cpu
}

func (c *CPU) LoadSnapshot(snapshot loader.Snapshot) {
	c.Reset()
	c.PC = snapshot.PC
	c.SP = snapshot.SP
	c.AF = snapshot.AF
	c.AF_ = snapshot.AF_
	c.BC = snapshot.BC
	c.BC_ = snapshot.BC_
	c.DE = snapshot.DE
	c.DE_ = snapshot.DE_
	c.HL = snapshot.HL
	c.HL_ = snapshot.HL_
	c.I = snapshot.I
	c.R = snapshot.R
	c.IX = snapshot.IX
	c.IY = snapshot.IY

	c.States.IM = snapshot.IM
	c.States.IFF1 = snapshot.IFF1
	c.States.IFF2 = snapshot.IFF2
}

func (c *CPU) HandleInterrupt() {
	if c.States.Halt {
		c.States.Halt = false
		c.PC++
	}
	c.States.IFF1, c.States.IFF2 = false, false
	c.pushStack(c.PC)

	switch c.States.IM {
	case 0:
		panic("IM 0")
	case 1:
		c.PC = 0x0038
		c.tstates += 7
	case 2:
		inttemp := uint16((uint16(c.I) << 8) | 0x00ff)
		c.PC = c.readWord(inttemp, 3, 3)
		c.tstates += 7
	}
}
