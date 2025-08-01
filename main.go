package main

func main() {
	// Initialize the CPU
	c := cpu{}

	setupGraphics()
	setupInput()

	// Initialize the Chip8 system and load the game into the memory
	c.init()
	c.loadGame("path/to/your/game.ch8")

	// Main emulation loop
	for {
		// Emulate one cycle
		c.emulateCycle()

		// If the draw flag is set, update the screen
		if c.drawFlag {
			c.drawGraphics()
		}
		c.setKeys()
	}
}
