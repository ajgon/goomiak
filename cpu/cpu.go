package cpu

type CPU struct {
	PC uint16
}

func (c *CPU) nop() uint8 {
	c.PC++

	return 4
}

func (c *CPU) Reset() {
	c.PC = 0
}

func CPUNew() *CPU {
	cpu := new(CPU)
	return cpu
}
