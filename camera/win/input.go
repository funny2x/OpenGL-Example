package win

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	"camera/camera"
)

func (w *Window) processInput() {
	//计算每帧的时间差
	currentFrame := glfw.GetTime()
	w.deltaTime = currentFrame - w.lastFrame
	w.lastFrame = currentFrame

	if w.gWin.GetKey(glfw.KeyEscape) == glfw.Press {
		w.gWin.SetShouldClose(true)
	} else if w.gWin.GetKey(glfw.KeyW) == glfw.Press {
		w.winInput.cam.ProcessKeyboard(camera.FORWARD, w.deltaTime)
	} else if w.gWin.GetKey(glfw.KeyS) == glfw.Press {
		w.winInput.cam.ProcessKeyboard(camera.BACKWARD, w.deltaTime)
	} else if w.gWin.GetKey(glfw.KeyA) == glfw.Press {
		w.winInput.cam.ProcessKeyboard(camera.LEFT, w.deltaTime)
	} else if w.gWin.GetKey(glfw.KeyD) == glfw.Press {
		w.winInput.cam.ProcessKeyboard(camera.RIGHT, w.deltaTime)
	}
}

type inputManager struct {
	firstMouse bool
	lastX      float64
	lastY      float64

	xoffset     float64
	yoffset     float64
	cam         *camera.Camera
	keysPressed [glfw.KeyLast]bool
}

func (im *inputManager) mouseCallback(window *glfw.Window, xpos, ypos float64) {

	if im.firstMouse {
		im.lastX = xpos
		im.lastY = ypos
		im.firstMouse = false
	}

	im.xoffset = xpos - im.lastX
	im.yoffset = ypos - im.lastY

	im.lastX = xpos
	im.lastY = ypos
	//TODO
	im.cam.ProcessMouseMovement(im.xoffset, im.yoffset, true)
}

//鼠标滚轮响应
func (im *inputManager) scrollCallback(window *glfw.Window, xoffset, yoffset float64) {
	im.cam.ProcessMouseScroll(yoffset)
}
