package main

import (
	"fmt"
	"os"
)

func (c *cpu) loadGame(path string) {
	// Load the game into memory at 0x200
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error loading game: %v\n", err)
		return
	}
	for i := 0; i < len(data); i++ {
		c.memory[i+512] = data[i]
	}
}
