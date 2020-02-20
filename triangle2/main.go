package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// 屏幕宽，高
const screenWidth = 600
const screenHeight = 400

// 顶点着色器和片段着色器源码
var vertexShaderSource = `
#version 330
layout (location = 0) in vec3 aPos;
void main() {
    gl_Position = vec4(aPos, 1.0);
}
` + "\x00"

var fragmentShaderSource = `
#version 330
out vec4 FragColor;
void main() {
    FragColor = vec4(1.0, 0.5, 0.2, 1.0f);
}
` + "\x00"

// 三角形的顶点数据
var triangle = []float32{
	// first triangle
	-0.9, -0.5, 0.0, // left
	0.0, -0.5, 0.0, // right
	-0.45, 0.5, 0.0, // top
	// second triangle
	0.0, -0.5, 0.0, // left
	0.9, -0.5, 0.0, // right
	0.45, 0.5, 0.0, // top
}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	// 初始化GLFW
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
	window, err := glfw.CreateWindow(screenWidth, screenHeight, "Triangle2 两个三角形", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}
	// 指定当前视口尺寸(前两个参数为左下角位置，后两个参数是渲染窗口宽、高)
	gl.Viewport(0, 0, screenWidth, screenHeight)

	//生成并绑定VAO和VBO
	var VAO uint32 // VAO
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	var VBO uint32 // VBO
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	// 将顶点数据绑定至当前默认的缓冲中
	gl.BufferData(gl.ARRAY_BUFFER, len(triangle)*4, gl.Ptr(triangle), gl.STATIC_DRAW)
	// 设置顶点属性指针
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	// 解绑VAO和VBO
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	// 生成并编译着色器
	// 顶点着色器
	var vertexShader = gl.CreateShader(gl.VERTEX_SHADER)
	csovert, free := gl.Strs(vertexShaderSource)
	gl.ShaderSource(vertexShader, 1, csovert, nil)
	free()
	gl.CompileShader(vertexShader)
	var success int32
	// 检查着色器是否成功编译，如果编译失败，打印错误信息
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLength1 int32
		gl.GetShaderiv(vertexShader, gl.INFO_LOG_LENGTH, &logLength1)

		log1 := strings.Repeat("\x00", int(logLength1+1))
		gl.GetShaderInfoLog(vertexShader, logLength1, nil, gl.Str(log1))
		fmt.Println("ERROR::SHADER::VERTEX::COMPILATION_FAILED", gl.Str(log1))
	}
	// 片段着色器
	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	csofragm, free := gl.Strs(fragmentShaderSource)
	gl.ShaderSource(fragmentShader, 1, csofragm, nil)
	free()
	gl.CompileShader(fragmentShader)
	// 检查着色器是否成功编译，如果编译失败，打印错误信息
	gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLength2 int32
		gl.GetShaderiv(fragmentShader, gl.INFO_LOG_LENGTH, &logLength2)

		log2 := strings.Repeat("\x00", int(logLength2+1))
		gl.GetShaderInfoLog(fragmentShader, logLength2, nil, gl.Str(log2))
		fmt.Println("ERROR::SHADER::FRAGMENT::COMPILATION_FAILED", gl.Str(log2))
	}
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)
	// 检查着色器是否成功链接，如果链接失败，打印错误信息
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength3 int32
		gl.GetShaderiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength3)

		log3 := strings.Repeat("\x00", int(logLength3+1))
		gl.GetShaderInfoLog(shaderProgram, logLength3, nil, gl.Str(log3))
		fmt.Println("ERROR::SHADER::FRAGMENT::LINKING_FAILED", gl.Str(log3))
	}
	// 删除着色器
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	// 渲染循环
	for !window.ShouldClose() {
		// 清空颜色缓冲
		gl.ClearColor(1.0, 1.0, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// 使用着色器程序
		gl.UseProgram(shaderProgram)
		// 绘制三角形
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		gl.BindVertexArray(0)

		// 交换缓冲并且检查是否有触发事件(比如键盘输入、鼠标移动等）
		glfw.PollEvents()
		window.SwapBuffers()

	}
	// 删除VAO和VBO
	gl.DeleteVertexArrays(1, &VAO)
	gl.DeleteBuffers(1, &VBO)
}
