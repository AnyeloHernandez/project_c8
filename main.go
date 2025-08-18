// //////////////////////////////////////////////////////////////////////////////////////////////////
// Thanks to Laurence Muller's guide on how to build a Chip8 emulator:
// https://multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/
//
// This is a simple Chip-8 emulator written in Go, using OpenGL for graphics rendering.
//
// Name: project-C8
//
// Author: AnyeloHernandez
// Contact: anyelo.g.hernandez@gmail.com
//
// ///////////////////////////////////////////////////////////////////////////////////////////////////
package main

import (
	"fmt"

	"github.com/go-gl/glfw/v3.2/glfw"
)

var KEY_MAP = map[glfw.Key]byte{
	glfw.Key1: 0x1,
	glfw.Key2: 0x2,
	glfw.Key3: 0x3,
	glfw.Key4: 0xC,
	glfw.KeyQ: 0x4,
	glfw.KeyW: 0x5,
	glfw.KeyE: 0x6,
	glfw.KeyR: 0xD,
	glfw.KeyA: 0x7,
	glfw.KeyS: 0x8,
	glfw.KeyD: 0x9,
	glfw.KeyF: 0xE,
	glfw.KeyZ: 0xA,
	glfw.KeyX: 0x0,
	glfw.KeyC: 0xB,
	glfw.KeyV: 0xF,
}

func main() {

	window := initWindowEmulator()
	defer glfw.Terminate()
	program := initOpenGL()

	// Initialize the CPU
	cpu := cpu{}

	// Initialize the Chip8 system and load the game into the memory
	cpu.initEmulator()
	cpu.loadGame("games/PONG")

	keyboardHandler(window, &cpu)

	windowShouldCloseHandler(&cpu, window, program)

}

func windowShouldCloseHandler(cpu *cpu, window *glfw.Window, program uint32) {
	for !window.ShouldClose() {
		// Main emulation loop
		// Emulate one cycle
		cpu.emulateCycle()
		cpu.codeDebugger()

		if cpu.drawFlag {
			drawSpriteOnWindow(window, program, cpu)
			//cpu.debugRender()
			cpu.drawFlag = false
		}

		// Handle timers
		cpu.handleTimers()

		glfw.PollEvents()

	}
}

func on_keyboard_pressed(c *cpu, key glfw.Key, action glfw.Action) {
	if key == glfw.KeyEscape && action == glfw.Press {
		glfw.Terminate()
		return
	}

	if symbol, ok := KEY_MAP[key]; ok {
		c.key[symbol] = 1
		fmt.Printf("Key pressed: %v, Action: %v\n", key, action)
	}
}

func on_keyboard_released(c *cpu, key glfw.Key, action glfw.Action) {
	if symbol, ok := KEY_MAP[key]; ok {
		c.key[symbol] = 0
		fmt.Printf("Key released: %v, Action: %v\n", key, action)
	}
}

func keyboardHandler(window *glfw.Window, c *cpu) {
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		switch action {
		case glfw.Press:
			on_keyboard_pressed(c, key, action)
		case glfw.Release:
			on_keyboard_released(c, key, action)
		}
	})
}
