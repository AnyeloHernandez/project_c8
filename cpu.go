package main

import (
	"fmt"
	"math/rand"
)

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
	drawFlag    bool
}

// CHIP-8 fontset: each character is 4x5 pixels
var chip8_fontset = [80]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

/*
0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
0x200-0xFFF - Program ROM and work RAM
*/

func (c *cpu) init() {
	// Init registers and memory once
	c.PC = 0x200 // Program counter starts at 0x200
	c.I = 0      // Index register starts at 0x000
	c.opcode = 0 // Reset opcode
	c.SP = 0     // Stack pointer starts at 0x00

	// Clear display
	for i := range 2048 {
		c.gfx[i] = 0
	}

	// Clear stack
	for i := range 16 {
		c.stack[i] = 0
	}

	for i := range 16 {
		c.key[i] = 0
		c.V[i] = 0
	}

	// Clear memory
	for i := range 4096 {
		c.memory[i] = 0
	}

	// Load fontset into memory
	for i := 0; i < 80; i++ {
		c.memory[i] = chip8_fontset[i]
	}

	// Reset timers
	c.delay_timer = 0
	c.sound_timer = 0

	c.drawFlag = true

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
	case 0x0000:
		switch c.opcode & 0x000F {
		case 0x0000: // 0x00E0: Clears the screen
			for i := 0; i < 2048; i++ {
				c.gfx[i] = 0x0
			}
			c.drawFlag = true
			c.PC += 2
		case 0x000E: // 0x00EE: Returns from a subroutine
			c.SP--               // Decrement stack pointer
			c.PC = c.stack[c.SP] // Set program counter to the address at the
			// top of the stack
			c.PC += 2 // Increment program counter by 2
		default:
			fmt.Printf("Unknown opcode: 0x%X\n", c.opcode)
		}

	case 0xA000: // ANNN
		c.I = c.opcode & 0x0FFF // Set index register
		c.PC += 2               // Increment program counter by 2
		// If the next opcode should be skiped, increase the PC by 4.
	case 0x1000: // 1NNN: Jumps to address NNN
		// Jumps, so we don't increment SP
		c.PC = c.opcode & 0x0FFF
	case 0x2000: // 2NNN
		// Call subroutine at address NNN
		c.stack[c.SP] = c.PC
		c.SP++
		c.PC = c.opcode & 0x0FFF
	// We are calling a subroutine, so we not increment the PC by 2.
	case 0x3000: // 0x3XNN (Type: conditional)
		// Skips the next instruction if VX equals NN
		if c.V[(c.opcode&0x0F00)>>8] == byte(c.opcode&0x00FF) {
			c.PC += 4 // Skip next instruction
		} else {
			c.PC += 2 // Just increment PC by 2
		}
	case 0x4000: // 0x4XNN (Type: conditional)
		// Skips the next instruction if VX does not equal NN
		if c.V[(c.opcode&0x0F00)>>8] != byte(c.opcode&0x00FF) {
			c.PC += 4 // Skip next instruction
		} else {
			c.PC += 2 // Just increment PC by 2
		}
	case 0x5000: // 0x5XY0 (Type: conditional)
		// Skips the next instruction if VX equals VY
		if c.V[(c.opcode&0x0F00)>>8] == c.V[(c.opcode&0x00F0)>>4] {
			c.PC += 4 // Skip next instruction
		} else {
			c.PC += 2 // Just increment PC by 2
		}
	case 0x6000: // 0x6XNN (Type: const)
		// Sets VX to NN
		c.V[(c.opcode&0x0F00)>>8] = byte(c.opcode & 0x00FF)
		c.PC += 2 // Increment program counter
	case 0x7000: // 0x7XNN (Type: const)
		// Adds NN to VX (Carry flag not changed)
		c.V[(c.opcode&0x0F00)>>8] += byte(c.opcode & 0x00FF)
		c.PC += 2
	case 0x8000: // 8XYn
		switch c.opcode & 0x000F {
		case 0x0000: // 8XY0: Sets Vx to the value of Vy
			c.V[(c.opcode&0x0F00)>>8] = c.V[byte(c.opcode&0x00F0>>4)]
			c.PC += 2
		case 0x0001: // 8XY1: Sets Vx to Vx OR Vy
			c.V[(c.opcode&0x0F00)>>8] |= c.V[byte(c.opcode&0x00F0>>4)]
			c.PC += 2
		case 0x0002: // 8XY2: Sets Vx to Vx AND Vy
			c.V[(c.opcode&0x0F00)>>8] &= c.V[byte(c.opcode&0x00F0>>4)]
			c.PC += 2
		case 0x0003: // 8XY3: Sets Vx to Vx XOR Vy
			c.V[(c.opcode&0x0F00)>>8] ^= c.V[byte(c.opcode&0x00F0>>4)]
			c.PC += 2
		case 0x0004: // 8XY4: Adds Vy to Vx. VF is set to 1 when there's overflow. 0 when not
			// Adds VY to VX. Vf is set to 1 when there's an overflow, and to 0 when there is not.
			if c.V[(c.opcode&0x00F0)>>4] > (0xFF - c.V[(c.opcode&0x0F00)>>8]) {
				c.V[0xF] = 1 // Set carry flag
			} else {
				c.V[0xF] = 0 // Clear carry flag
			}
			c.V[(c.opcode&0x0F00)>>8] += c.V[(c.opcode&0x00F0)>>4]
			c.PC += 2
		case 0x0005: // 8XY5: Vy is subtracted from Vx. VF is set to 0 when there's and underflow. 1 when not
			if c.V[(c.opcode&0x0F00)>>8] >= c.V[(c.opcode&0x00F0)>>4] {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
			c.V[(c.opcode&0x0F00)>>8] -= c.V[(c.opcode&0x00F0)>>4]
		case 0x0006: // 8XY6: Set Vx to Vy and shift Vx one bit to the right, set Vf to the bit shifted out, even if X=F!
			c.V[0xF] = c.V[(c.opcode&0x0F00)>>8] & 0x1
			c.V[(c.opcode&0x0F00)>>8] >>= 1
			c.PC += 2
		case 0x0007: // 8XY7: Set Vx to the result of subtracting Vx from Vy, Vf is set to 0 if an underflow happened, to 1 if not, even if X=F!
			if c.V[(c.opcode&0x0F00)>>8] <= c.V[(c.opcode&0x00F0)>>4] {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
			c.V[(c.opcode&0x0F00)>>8] = c.V[(c.opcode&0x0F00)>>8] - c.V[(c.opcode&0x00F0)>>4]
			c.PC += 2
		case 0x000E: // 8XYE: Set Vx to Vy and shift Vx one bit to the left, set Vf to the bit shifted out, even if X=F!
			c.V[0xF] = c.V[(c.opcode&0x0F00)>>8] >> 7
			c.V[(c.opcode&0x0F00)>>8] <<= 1
			c.PC += 2
		default:
			fmt.Printf("Unknown opcode: 0x%X\n", c.opcode)
		}
	case 0x9000: // 9XY0: Skips the next instruction if VX does not equal VY
		if c.V[(c.opcode&0x0F00)>>8] != c.V[(c.opcode&0x00F0)>>4] {
			c.PC += 4
		} else {
			c.PC += 2
		}
	case 0xB000: // BNNN: Jumps to the addres plus V0. PC = V0 + NNN
		c.PC = uint16(c.V[0x0] + byte(c.opcode&0x0FFF))
	case 0xC000: // 0xCXNN: Sets Vx to the result of a bitwise AND operation on a random number and NN.
		c.V[(c.opcode&0x0F00)>>8] = byte(rand.Intn(256)) & 0x00FF
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

		c.V[0xF] = 0 // Resets the register VF (Collision flag)
		for yline := 0; yline < int(height); yline++ {
			pixel = uint32(c.memory[c.I+uint16(yline)])
			for xline := 0; xline < 8; xline++ {
				if pixel&(0x80>>xline) != 0 {
					if c.gfx[(x+uint32(xline)+((y+uint32(yline))*64))] == 1 {
						// If pixel is set to 1, set VF to 1 (collision)
						c.V[0xF] = 1
					}
					// Set the pixel value by using XOR
					c.gfx[x+uint32(xline)+((y+uint32(yline))*64)] ^= 1
				}
			}
		}
		c.drawFlag = true
		c.PC += 2
	case 0xE000: // 0xE000 is a prefix for key input opcodes
		/*
			0xE09E: Skip next instruction if key with value of Vx is pressed
			0xE0A1: Skip next instruction if key with value of Vx is not pressed

			It doesn't matter which key is pressed, we just check the key state
		*/
		switch c.opcode & 0x00FF {
		case 0x9E: // Skip next instruction if key with value of Vx is pressed
			if c.key[c.V[(c.opcode&0x0F00)>>8]] != 0 {
				c.PC += 4 // Skip next instruction
			} else {
				c.PC += 2 // Just increment PC by 2
			}
		case 0xA1: // Skip next instruction if key with value of Vx is not pressed
			if c.key[c.V[(c.opcode&0x0F00)>>8]] == 0 {
				c.PC += 4 // Skip next instruction
			} else {
				c.PC += 2 // Just increment PC by 2
			}
		}
	case 0xF000: // 0xF000 is a prefix for various opcodes
		switch c.opcode & 0x00FF {
		case 0x0007: // 0xFX07: Sets Vx to the value of the delay timer
			c.V[(c.opcode&0x0F00)>>8] = c.delay_timer
			c.PC += 2
		case 0x000A: // 0xFX0A: A key press is awaited, and then stored in Vx
			// get_key() functionality
			var keypress bool = false
			for i := 0; i < 16; i++ {
				if c.key[i] != 0 {
					c.V[(c.opcode&0x0F00)>>8] = byte(i)
					keypress = true
				}
			}
			// If not key is recieved, skip this cycle
			if !keypress {
				return
			}
			c.PC += 2
		case 0x0015: // 0xFX15: Sets the delay timer to Vx
			c.delay_timer = c.V[(c.opcode&0x0F00)>>8]
			c.PC += 2
		case 0x0018: // 0xFX18: Sets the sound timer to Vx
			c.sound_timer = c.V[(c.opcode&0x0F00)>>8]
			c.PC += 2
		case 0x001E: // 0xFX1E: Adds Vx to I. Vf is not affected.
			// We check if there's an overflow
			if (c.I + uint16(c.V[c.opcode&0x0F00>>8])) > 0x00FF {
				c.V[0xF] = 1 // Set carry flag
			} else {
				c.V[0xF] = 0 // Clear carry flag
			}
			c.I = uint16(c.V[c.opcode&0x0F00>>8])
			c.PC += 2
		case 0x0029: // 0xFX29: Sets I to the location of the sprite for the caracter in Vx(considering the lowest nibble only)
			c.I = uint16(c.V[(c.opcode&0x0F00)>>8]) * 0x5 // Each character is 5 bytes
			c.PC += 2
		case 0x0033: // FX33
			// Stores the binary-coded decimal representation of VX in memory locations I, I+1, and I+2.
			c.memory[c.I] = c.V[(c.opcode&0x0F00)>>8] / 100
			c.memory[c.I+1] = (c.V[(c.opcode&0x0F00)>>8] / 10) % 10
			c.memory[c.I+2] = c.V[(c.opcode&0x0F00)>>8] % 10
			c.PC += 2
		case 0x0055: // 0xFX55: Stores from V0 to Vx in memory, starting at address I. The offset from I is increased by 1 for each
			// value written, but I itself is left unmodified.
			for i := 0; i <= int((c.opcode&0x0F00)>>8); i++ {
				c.memory[c.I+uint16(i)] = c.V[uint16(i)]
			}
			// On the original interpreter, when the operation is done I = I + X +1
			c.I += ((c.opcode & 0x0F00) >> 8) + 1
			c.PC += 2
		case 0x0065: // 0xFX65: Fills from V0 to VX (including VX) with values from memory, starting at address I.
			// The offset from I is increased by 1 for each value read, but I itself is left unmodified.
			for i := 0; i <= int((c.opcode&0x0F00)>>8); i++ {
				c.V[uint16(i)] = c.memory[c.I+uint16(i)]
			}
			// On the original interpreter, when the operation is done I = I + X +1
			c.I += ((c.opcode & 0x0F00) >> 8) + 1
			c.PC += 2
		}

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

func (c *cpu) debugRender() {
	// Draw
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if c.gfx[(y*64)+x] == 0 {
				fmt.Printf("O")
			} else {
				fmt.Printf(" ")
			}
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
	}
}
