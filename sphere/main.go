package main

import (
	"log"
	"math"
	"runtime"
	"sphere/gfx"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

//  将球横纵划分成50X50的网格

const Y_SEGMENTS = 50
const X_SEGMENTS = 50

func init() {
	// GLFW event handling must be run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to inifitialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(600, 600, "sphere", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow (go function bindings)
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// gl.Viewport(0, 0, 600, 600)
	window.SetKeyCallback(keyCallback)

	err = programLoop(window)
	if err != nil {
		log.Fatal(err)
	}
}

/*
 * Creates the Vertex Array Object for a triangle.
 */
func createVAO(vertices []float32, indices []uint32) uint32 {

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	var EBO uint32
	gl.GenBuffers(1, &EBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(VAO)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// copy indices into element buffer
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// 设置顶点属性指针
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 12, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	// gl.BindVertexArray(0)

	return VAO
}

func programLoop(window *glfw.Window) error {

	// the linked shader program determines how the data will be rendered
	vertShader, err := gfx.NewShaderFromFile("shader/task3.vs", gl.VERTEX_SHADER)
	if err != nil {
		return err
	}

	fragShader, err := gfx.NewShaderFromFile("shader/task3.fs", gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}

	shaderProgram, err := gfx.NewProgram(vertShader, fragShader)
	if err != nil {
		return err
	}
	defer shaderProgram.Delete()

	var vertices []float32 //生成球的顶点
	for y := 0; y <= Y_SEGMENTS; y++ {
		for x := 0; x <= X_SEGMENTS; x++ {
			xSegment := float32(x) / float32(X_SEGMENTS)
			ySegment := float32(y) / float32(Y_SEGMENTS)

			xPos := float32(math.Cos(float64(xSegment*math.Pi*2.0)) * math.Sin(float64(ySegment*math.Pi)))
			yPos := float32(math.Cos(float64(ySegment * math.Pi)))
			zPos := float32(math.Sin(float64(xSegment*math.Pi*2.0)) * math.Sin(float64(ySegment*math.Pi)))

			vertices = append(vertices, xPos, yPos, zPos)
		}
	}
	var indices []uint32 //生成球的Indices
	for i := 0; i < Y_SEGMENTS; i++ {
		for j := 0; j < X_SEGMENTS; j++ {
			a1 := uint32(i*(X_SEGMENTS+1) + j)
			a2 := uint32((i+1)*(X_SEGMENTS+1) + j)
			a3 := uint32((i+1)*(X_SEGMENTS+1) + j + 1)
			b1 := uint32(i*(X_SEGMENTS+1) + j)
			b2 := uint32((i+1)*(X_SEGMENTS+1) + j + 1)
			b3 := uint32(i*(X_SEGMENTS+1) + j + 1)
			indices = append(indices, a1, a2, a3, b1, b2, b3)
		}
	}
	VAO := createVAO(vertices, indices)

	for !window.ShouldClose() {
		// poll events and call their registered callbacks
		glfw.PollEvents()

		// background color
		gl.ClearColor(0.0, 0.34, 0.57, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		shaderProgram.Use()
		//开启面剔除(只需要展示一个面，否则会有重合)
		gl.Enable(gl.CULL_FACE)
		gl.CullFace(gl.BACK)

		gl.BindVertexArray(VAO)
		//使用线框模式绘制
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		gl.DrawElements(gl.TRIANGLES, X_SEGMENTS*Y_SEGMENTS*6, gl.UNSIGNED_INT, unsafe.Pointer(nil))
		gl.BindVertexArray(0)

		window.SwapBuffers()
	}

	return nil
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action,
	mods glfw.ModifierKey) {

	// When a user presses the escape key, we set the WindowShouldClose property to true,
	// which closes the application
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}
