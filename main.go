package main

import (
	"fmt"

	"github.com/go-gl/glfw/v3.2/glfw"

	"runtime"
)

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
	c.loadGame("games/TANK")

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			on_keyboard_pressed(&c, window, key, action)
		} else if action == glfw.Release {
			on_keyboard_released(&c, window, key, action)
		}
	})

	for !window.ShouldClose() {

		//setupGraphics()
		// setupInput()

		// Main emulation loop
		// Emulate one cycle
		c.emulateCycle()

		if c.drawFlag {
			draw(window, program, &c)
			c.debugRender()
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

func on_keyboard_pressed(c *cpu, window *glfw.Window, key glfw.Key, action glfw.Action) {
	if key == glfw.KeyEscape && action == glfw.Press {
		glfw.Terminate()
		return
	}

	switch key {
	case glfw.Key1:
		c.key[0x1] = 1
	case glfw.Key2:
		c.key[0x2] = 1
	case glfw.Key3:
		c.key[0x3] = 1
	case glfw.Key4:
		c.key[0xC] = 1

	case glfw.KeyQ:
		c.key[0x4] = 1
	case glfw.KeyW:
		c.key[0x5] = 1
	case glfw.KeyE:
		c.key[0x6] = 1
	case glfw.KeyR:
		c.key[0xD] = 1

	case glfw.KeyA:
		c.key[0x7] = 1
	case glfw.KeyS:
		c.key[0x8] = 1
	case glfw.KeyD:
		c.key[0x9] = 1
	case glfw.KeyF:
		c.key[0xE] = 1

	case glfw.KeyZ:
		c.key[0xA] = 1
	case glfw.KeyX:
		c.key[0x0] = 1
	case glfw.KeyC:
		c.key[0xB] = 1
	case glfw.KeyV:
		c.key[0xF] = 1

		fmt.Printf("Key pressed: %v, Action: %v\n", key, action)
	}
}

func on_keyboard_released(c *cpu, window *glfw.Window, key glfw.Key, action glfw.Action) {
	if action == glfw.Release {
		switch key {
		case glfw.Key1:
			c.key[0x1] = 0
		case glfw.Key2:
			c.key[0x2] = 0
		case glfw.Key3:
			c.key[0x3] = 0
		case glfw.Key4:
			c.key[0xC] = 0

		case glfw.KeyQ:
			c.key[0x4] = 0
		case glfw.KeyW:
			c.key[0x5] = 0
		case glfw.KeyE:
			c.key[0x6] = 0
		case glfw.KeyR:
			c.key[0xD] = 0

		case glfw.KeyA:
			c.key[0x7] = 0
		case glfw.KeyS:
			c.key[0x8] = 0
		case glfw.KeyD:
			c.key[0x9] = 0
		case glfw.KeyF:
			c.key[0xE] = 0

		case glfw.KeyZ:
			c.key[0xA] = 0
		case glfw.KeyX:
			c.key[0x0] = 0
		case glfw.KeyC:
			c.key[0xB] = 0
		case glfw.KeyV:
			c.key[0xF] = 0
		}
	}
}
