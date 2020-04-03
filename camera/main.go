/***
 * 例程  摄像机类
 * 步骤:
 * 加载摄像机类进行操作即可
 */

package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"camera/camera"
	"camera/shader"
	"camera/win"
)

const (
	screenWidth  = 600 //窗口宽度
	screenHeight = 600 //窗口高度
)

//立方体数组
var vertices = []float32{
	-0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
	0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 0.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 0.0, 0.0,
	-0.5, 0.5, -0.5, 1.0, 0.0, 0.0,
	-0.5, -0.5, -0.5, 1.0, 0.0, 0.0,

	-0.5, -0.5, 0.5, 0.0, 1.0, 0.0,
	0.5, -0.5, 0.5, 0.0, 1.0, 0.0,
	0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
	0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
	-0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
	-0.5, -0.5, 0.5, 0.0, 1.0, 0.0,

	-0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
	-0.5, 0.5, -0.5, 0.0, 0.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 0.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 0.0, 1.0,
	-0.5, -0.5, 0.5, 0.0, 0.0, 1.0,
	-0.5, 0.5, 0.5, 0.0, 0.0, 1.0,

	0.5, 0.5, 0.5, 0.5, 0.0, 0.0,
	0.5, 0.5, -0.5, 0.5, 0.0, 0.0,
	0.5, -0.5, -0.5, 0.5, 0.0, 0.0,
	0.5, -0.5, -0.5, 0.5, 0.0, 0.0,
	0.5, -0.5, 0.5, 0.5, 0.0, 0.0,
	0.5, 0.5, 0.5, 0.5, 0.0, 0.0,

	-0.5, -0.5, -0.5, 0.0, 0.5, 0.0,
	0.5, -0.5, -0.5, 0.0, 0.5, 0.0,
	0.5, -0.5, 0.5, 0.0, 0.5, 0.0,
	0.5, -0.5, 0.5, 0.0, 0.5, 0.0,
	-0.5, -0.5, 0.5, 0.0, 0.5, 0.0,
	-0.5, -0.5, -0.5, 0.0, 0.5, 0.0,

	-0.5, 0.5, -0.5, 0.0, 0.0, 0.5,
	0.5, 0.5, -0.5, 0.0, 0.0, 0.5,
	0.5, 0.5, 0.5, 0.0, 0.0, 0.5,
	0.5, 0.5, 0.5, 0.0, 0.0, 0.5,
	-0.5, 0.5, 0.5, 0.0, 0.0, 0.5,
	-0.5, 0.5, -0.5, 0.0, 0.0, 0.5,
}

var cam = camera.GetCamera(mgl32.Vec3{0.0, 0.0, 3.0})

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	window := win.NewWindow(screenWidth, screenHeight, "Camera", cam)
	//-----------------------------------------
	//鼠标设置
	//-----------------------------------------
	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}
	gl.Viewport(0, 0, screenWidth, screenHeight)

	//加载着色器
	camShader, err := shader.NewShader("src/task-camera.vs", "src/task-camera.fs")
	if err != nil {
		log.Panic(err)
	}

	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	gl.GenBuffers(1, &VBO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	// 将顶点数据绑定至当前默认的缓冲中
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	// 设置顶点属性指针
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(12))
	gl.EnableVertexAttribArray(1)

	gl.Enable(gl.DEPTH_TEST)
	for !window.ShouldClose() {
		window.StartProcessInput()
		gl.ClearColor(0.0, 0.34, 0.57, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) //清理颜色缓冲和深度缓冲

		camShader.Use()
		model := mgl32.HomogRotate3D(float32(glfw.GetTime()), mgl32.Vec3{0.5, 1.0, 0.0})
		//-------------------------------------
		// Transform坐标变换矩
		view := cam.GetViewMatrix()

		projection := mgl32.Perspective(45.0, float32(screenWidth/screenHeight), 0.1, 10.0)
		// 向着色器中传入参数
		camShader.SetMat4("model", model)
		camShader.SetMat4("view", view)
		camShader.SetMat4("projection", projection)

		//绘制
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
		gl.BindVertexArray(0)

	}
	//释放VAOVBO
	gl.DeleteVertexArrays(1, &VAO)
	gl.DeleteBuffers(1, &VBO)
}
