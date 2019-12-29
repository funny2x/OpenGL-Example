package main

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"runtime"
)

const (
	//屏幕宽高
	Width  = 600
	Height = 400
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()
	// 不可改变窗口大小
	glfw.WindowHint(glfw.Resizable, glfw.False)
	// OpenGL版本为4.1
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// 创建窗口(宽、高、窗口名称)
	window, err := glfw.CreateWindow(Width, Height, "window", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}
	// 指定当前视口尺寸(前两个参数为左下角位置，后两个参数是渲染窗口宽、高)
	gl.Viewport(0, 0, Width, Height)

	//渲染循环
	for !window.ShouldClose() {
		//清空颜色缓存
		gl.ClearColor(0.0, 0.34, 0.57, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		//交换缓冲并且检查是否有触发事件(比如键盘输入、鼠标移动等）
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
