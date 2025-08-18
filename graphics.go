package main

/*
Especial thanks to:
https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl

For the OpenGL tutorial that helped me understand how to use OpenGL with Go.
*/
import (
	"fmt"
	"strings"

	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	WIDTH  = 64
	HEIGHT = 32

	ROWS = 64
	COLS = 32

	vertexShaderSource = `
		#version 410
		in vec3 vp;
		void main() {
			gl_Position = vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410
		out vec4 frag_colour;
		void main() {
			frag_colour = vec4(1, 1, 1, 1.0);
		}
	` + "\x00"

	vertexShaderSourceV2 = `
		#version 120
		attribute vec3 vp;
		void main() {
			gl_Position = vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentShaderSourceV2 = `
		#version 120
		void main() {
			gl_FragColor = vec4(1, 1, 1, 1.0);
		}
	` + "\x00"

	vertexShaderSourceV3 = `
		#version 330 core
		layout(location = 0) in vec3 vp;
		void main() {
			gl_Position = vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentShaderSourceV3 = `
		#version 330 core
		out vec4 frag_colour;
		void main() {
			frag_colour = vec4(1, 1, 1, 1.0);
		}
	` + "\x00"
)

var (
	square = []float32{
		-0.5, 0.5, 0.0,
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,

		-0.5, 0.5, 0.0,
		0.5, 0.5, 0.0,
		0.5, -0.5, 0.0,
	}
)

type cell struct {
	drawable uint32
	x        int
	y        int
}

// func setupGraphics() {
// 	// Ensure the main thread is the only one running OpenGL
// 	runtime.LockOSThread()
// 	// Initialize GLFW
// 	window := initGlfw()
// 	defer glfw.Terminate()

// 	program := initOpenGL()

// 	for !window.ShouldClose() {
// 		// Do openGl stuff

// 	}
// }

func initWindowEmulator() *glfw.Window {
	runtime.LockOSThread()
	// Initialize GLFW
	return initGlfw()
}

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("failed to initialize GLFW: %v", err))
	}
	// glfw.WindowHint(glfw.ContextVersionMajor, 4)
	// glfw.WindowHint(glfw.ContextVersionMinor, 1)
	// glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	// glfw.WindowHint(glfw.Resizable, glfw.False)
	// glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Create a windowed mode window and its OpenGL context
	window, err := glfw.CreateWindow(WIDTH*10, HEIGHT*10, "Chip8 Emulator", nil, nil)
	if err != nil {
		panic(fmt.Errorf("failed to create window: %v", err))
	}

	window.MakeContextCurrent()
	return window
}

func initOpenGL() uint32 {
	var major int
	major, _, _ = glfw.GetVersion()
	if major < 3 {
		println("OpenGL version 2 found")
		return initOpenGL2()
	} else if major == 3 {
		return initOpenGL3()
	} else {
		println("OpenGL version 3 or higher found")
		return initOpenGL4()
	}

}

func initOpenGL3() uint32 {
	if err := gl.Init(); err != nil {
		panic(fmt.Errorf("failed to initialize OpenGL: %v", err))
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Printf("OpenGL version: %s\n", version)

	vertexShader, err := compileShader(vertexShaderSourceV3, gl.VERTEX_SHADER)
	if err != nil {
		panic(fmt.Errorf("failed to compile vertex shader: %v", err))
	}
	fragmentShader, err := compileShader(fragmentShaderSourceV3, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(fmt.Errorf("failed to compile fragment shader: %v", err))
	}
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func initOpenGL2() uint32 {
	if err := gl.Init(); err != nil {
		panic(fmt.Errorf("failed to initialize OpenGL: %v", err))
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Printf("OpenGL version: %s\n", version)

	vertexShader, err := compileShader(vertexShaderSourceV2, gl.VERTEX_SHADER)
	if err != nil {
		panic(fmt.Errorf("failed to compile vertex shader: %v", err))
	}
	fragmentShader, err := compileShader(fragmentShaderSourceV2, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(fmt.Errorf("failed to compile fragment shader: %v", err))
	}
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func initOpenGL4() uint32 {
	if err := gl.Init(); err != nil {
		panic(fmt.Errorf("failed to initialize OpenGL: %v", err))
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Printf("OpenGL version: %s\n", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(fmt.Errorf("failed to compile vertex shader: %v", err))
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(fmt.Errorf("failed to compile fragment shader: %v", err))
	}
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func drawGraphics(c *cpu, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)
	// Draw the graphics
	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			if c.gfx[y*WIDTH+x] != 0 {
				// Draw a pixel at (x, y)
				drawPixel(x, y, program)
			}
		}
	}

	window.SwapBuffers()
	c.drawFlag = false // Reset the draw flag after drawing
}

func drawPixel(x, y int, program uint32) {
	// Convertir coordenadas CHIP-8 (0-63, 0-31) a OpenGL (-1 a 1)
	pixelWidth := 2.0 / float32(WIDTH)
	pixelHeight := 2.0 / float32(HEIGHT)

	// Ajustamos que Y vaya de arriba hacia abajo como en el Chip8
	adjustedY := HEIGHT - 1 - y

	startX := -1.0 + float32(x)*pixelWidth
	startY := -1.0 + float32(adjustedY)*pixelHeight
	endX := startX + pixelWidth
	endY := startY + pixelHeight

	vertices := []float32{
		startX, endY, 0.0, // Bottom left
		endX, endY, 0.0, // Bottom right
		endX, startY, 0.0, // Top right

		startX, endY, 0.0, // Bottom left
		endX, startY, 0.0, // Top right
		startX, startY, 0.0, // Top left
	}

	vao := makeVao(vertices)
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.DeleteVertexArrays(1, &vao)
}

func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func drawSpriteOnWindow(window *glfw.Window, program uint32, cpu *cpu) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			if cpu.gfx[y*WIDTH+x] != 0 {
				drawPixel(x, y, program)
			}
		}
	}

	glfw.PollEvents()
	window.SwapBuffers()
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csource, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile shader: %s", log)
	}

	return shader, nil
}

func makeCells() [][]*cell {
	cells := make([][]*cell, ROWS, ROWS)
	for x := range ROWS {
		for y := range COLS {
			c := newCell(x, y)
			cells[x] = append(cells[x], c)
		}
	}
	return cells
}

func newCell(x, y int) *cell {
	points := make([]float32, len(square))
	copy(points, square)

	for i := 0; i < len(points); i++ {
		var position float32
		var size float32
		switch i % 3 {
		case 0:
			size = 1.0 / float32(COLS)
			position = float32(x) * size
		case 1:
			size = 1.0 / float32(ROWS)
			position = float32(y) * size
		default:
			continue
		}

		if points[i] < 0 {
			points[i] = (position * 2) - 1
		} else {
			points[i] = ((position + size) * 2) - 1
		}
	}

	return &cell{
		drawable: makeVao(points),

		x: x,
		y: y,
	}
}

func (c *cell) draw() {
	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
}
