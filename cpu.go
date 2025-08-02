package main

type cpu struct {
	opcode      uint16
	memory      [4096]byte
	V           [16]byte
	I           uint16
	PC          uint16
	SP          byte          // Stack pointer
	gfx         [64 * 32]byte // Graphics buffer (64x32 pixels)
	delay_timer byte          // Delay timer
	sound_timer byte          // Sound timer
	stack       [16]uint16    // Stack for subroutine calls
	key         [16]byte      // Key state (0-15)
	fontset     [80]byte      // Fontset for Chip 8
	drawFlag    bool
}

/*
0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
0x200-0xFFF - Program ROM and work RAM
*/

func (c *cpu) init() {
	// Init registers and memory once
	c.PC = 0x200     // Program counter starts at 0x200
	c.I = 0x000      // Index register starts at 0x000
	c.opcode = 0x000 // Reset opcode
	c.SP = 0x00      // Stack pointer starts at 0x00

	// Load fontset
	for i := 0; i < 80; i++ {
		c.memory[i] = c.fontset[i]
	}
}

/*
Every cycle, the method emulateCycle is called which emulates one cycle of the Chip 8.
During this cycle, the CPU fetches the opcode, decodes it, and executes it.
*/
func (c *cpu) emulateCycle() {
	// Fetch opcode
	c.opcode = uint16(c.memory[c.PC])<<8 | uint16(c.memory[c.PC+1])
	/*
			The system will fetch the opcode from the memory at the location
			specified by the PC.

			So we obtain for example 0xA2F0 by reading PC and 0xF0 on PC+1.
			We have to merge both bytes to form the full opcode.

			First we add 8 bits on PC using a left shift.
			uint16(c.memory[c.PC])<<8
			Then we merge the next 8 bits on PC+1 using a bitwise OR.
			c.opcode = uint16(c.memory[c.PC])<<8 | uint16(c.memory[c.PC+1])


			0xA2       0xA2 << 8 = 0xA200   HEX
			10100010   1010001000000000     BIN

			1010001000000000 | // 0xA200
		       	    11110000 = // 0xF0 (0x00F0)
			------------------
			1010001011110000   // 0xA2F0

		So now we have the full opcode 0xA2F0 in c.opcode.
		We do a switch stmt to verify which opcode we retrieved.

		Since only 12 biits are containing the value we need to store, we use &
		0x0FFF to get rid of the first four bits.
	*/

	switch c.opcode & 0xF000 {
	case 0xA000:
		c.I = c.opcode & 0x0FFF // Set index register
		c.PC += 2               // Increment program counter by 2
		// If the next opcode should be skiped, increase the PC by 4.
	case 0x2000:
		// Call subroutine at address NNN
		c.stack[c.SP] = c.PC
		c.SP++
		c.PC = c.opcode & 0x0FFF
		// We are calling a subroutine, so we not increment the PC by 2.
	case 0x0004: // 8XY4
		// Adds VY to VX. Vf is set to 1 when there's an overflow, and to 0 when there is not.
		if c.V[(c.opcode&0x00F0)>>4] > (0xFF - c.V[(c.opcode&0x0F00)>>8]) {
			c.V[0xF] = 1 // Set carry flag
		} else {
			c.V[0xF] = 0 // Clear carry flag
		}
		c.V[(c.opcode&0x0F00)>>8] += c.V[(c.opcode&0x00F0)>>4]
		c.PC += 2
	case 0x0033: // FX33
		// Stores the binary-coded decimal representation of VX in memory locations I, I+1, and I+2.
		c.memory[c.I] = c.V[(c.opcode&0x0F00)>>8] / 100
		c.memory[c.I+1] = (c.V[(c.opcode&0x0F00)>>8] / 10) % 10
		c.memory[c.I+2] = c.V[(c.opcode&0x0F00)>>8] % 10
		c.PC += 2
	case 0xD000: // DXYN
		/* This opcode is responsible for drawing to the display
		Has a width of 8 pixels and a height of N pixels.

		Each row of 8 pixels is read as a bit-coded starting from memory location.


		*/
		x := uint32(c.V[(c.opcode&0x0F00)>>8])
		y := uint32(c.V[(c.opcode&0x00F0)>>4])
		height := uint32(c.opcode & 0x000F)
		var pixel uint32

		c.V[0xF] = 0
		for yline := 0; yline < int(height); yline++ {
			pixel = uint32(c.memory[c.I+uint16(yline)])
			for xline := 0; xline < 8; xline++ {
				if pixel&(0x80>>xline) != 0 {
					if c.gfx[(x+uint32(xline)+((y+uint32(yline))*64))] == 1 {
						c.V[0xF] = 1
					}
					c.gfx[x+uint32(xline)+((y+uint32(yline))*64)] ^= 1
				}
			}
		}
		c.drawFlag = true
		c.PC += 2

	default:
		println("Unknown opcode:", c.opcode)
		// Decode opcode

		// Execute Opcode

		// Update timers
		if c.delay_timer > 0 {
			c.delay_timer--
		}
		if c.sound_timer > 0 {
			if c.sound_timer == 1 {
				println("BEEP!") // Sound timer reached 0, play sound
			}
			c.sound_timer--
		}
	}
}
