package cpu

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

var ioIncTemp1 [16]bool = [16]bool{
	false, false, true, false, false, true, false, true,
	true, false, true, true, false, true, true, false,
}

var ioDecTemp1 [16]bool = [16]bool{
	false, true, false, false, true, false, false, true,
	false, false, true, false, false, true, false, true,
}

var ioTemp2 [256]bool = [256]bool{
	/*	      0     1      2     3      4     5      6     7      8     9      A     B      C     D      E     F */
	/* 0 */ false, false, true, true, false, true, false, false, true, true, false, false, true, false, true, true,
	/* 1 */ false, true, false, false, true, false, true, true, false, false, true, true, false, true, false, false,
	/* 2 */ true, true, false, false, true, false, true, true, false, false, true, true, false, true, false, false,
	/* 3 */ true, false, true, true, false, true, false, false, true, true, false, false, true, false, true, true,
	/* 4 */ false, true, false, false, true, false, true, true, false, false, true, true, false, true, false, false,
	/* 5 */ true, false, true, true, false, true, false, false, true, true, false, false, true, false, true, true,
	/* 6 */ false, false, true, true, false, true, false, false, true, true, false, false, true, false, true, true,
	/* 7 */ false, true, false, false, true, false, true, true, false, false, true, true, false, true, false, false,
	/* 8 */ true, true, false, false, true, false, true, true, false, false, true, true, false, true, false, false,
	/* 9 */ true, false, true, true, false, true, false, false, true, true, false, false, true, false, true, true,
	/* A */ false, false, true, true, false, true, false, false, true, true, false, false, true, false, true, true,
	/* B */ false, true, false, false, true, false, true, true, false, false, true, true, false, true, false, false,
	/* C */ true, false, true, true, false, true, false, false, true, true, false, false, true, false, true, true,
	/* D */ false, true, false, false, true, false, true, true, false, false, true, true, false, true, false, false,
	/* E */ true, true, false, false, true, false, true, true, false, false, true, true, false, true, false, false,
	/* F */ true, false, true, true, false, true, false, false, true, true, false, false, true, false, true, true,
}
