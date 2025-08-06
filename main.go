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

	"runtime"
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

	// Ensure the main thread is the only one running OpenGL
	runtime.LockOSThread()
	// Initialize GLFW
	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()
	//cells := makeCells()

	// Initialize the CPU
	c := cpu{}

	// Initialize the Chip8 system and load the game into the memory
	c.init()
	c.loadGame("games/PONG")

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			on_keyboard_pressed(&c, key, action)
		} else if action == glfw.Release {
			on_keyboard_released(&c, key, action)
		}
	})

	for !window.ShouldClose() {
		// Main emulation loop
		// Emulate one cycle
		c.emulateCycle()
		c.codeDebugger()

		if c.drawFlag {
			draw(window, program, &c)
			//c.debugRender()
			c.drawFlag = false
		}

		// Handle timers
		if c.delay_timer > 0 {
			c.delay_timer--
		}
		if c.sound_timer > 0 {
			if c.sound_timer == 1 {
				fmt.Println("BEEP!")
			}
			c.sound_timer--
		}
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
