package win

import (
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"

	"camera/camera"
)

//Window 窗口结构体
type Window struct {
	width     int
	height    int
	gWin      *glfw.Window //窗口
	winInput  *inputManager
	deltaTime float64
	lastFrame float64
}

//NewWindow 窗口结构体Window构造函数
func NewWindow(width, height int, title string, cam *camera.Camera) *Window {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	gWindow, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		log.Fatalln(err)
	}
	x := float64(width / 2)
	y := float64(height / 2)

	gWindow.MakeContextCurrent()
	gWindow.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	im := &inputManager{
		firstMouse: true,
		lastX:      x,
		lastY:      y,
		cam:        cam,
	}
	// gWindow.SetKeyCallback(im.keyCallback)
	gWindow.SetCursorPosCallback(im.mouseCallback)
	gWindow.SetScrollCallback(im.scrollCallback)

	return &Window{
		width:  width,
		height: height,
		gWin:   gWindow,

		deltaTime: 0.0,
		lastFrame: 0.0,
		winInput:  im,
	}
}

//Width 返回窗口宽
func (w *Window) Width() int {
	return w.width
}

//Height 返回窗口高
func (w *Window) Height() int {
	return w.height
}

//ShouldClose 询问窗口是否需要关闭
func (w *Window) ShouldClose() bool {
	return w.gWin.ShouldClose()
}

//StartProcessInput ->SwapBuffers and glfw.PollEvents processInput键盘输入
func (w *Window) StartProcessInput() {
	// swap in the previous rendered buffer
	w.gWin.SwapBuffers()
	// poll for UI window events
	glfw.PollEvents()
	//检测键盘输入
	w.processInput()
}
