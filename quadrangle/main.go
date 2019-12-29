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
const screen_width = 600
const screen_height = 400

// 顶点着色器和片段着色器源码
var vertex_shader_source = `
#version 330
layout (location = 0) in vec3 aPos;
void main() {
    gl_Position = vec4(aPos, 1.0);
}
` + "\x00"

var fragment_shader_source = `
#version 330
out vec4 FragColor;
void main() {
    FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
}
` + "\x00"

// 三角形的顶点数据
var triangle = []float32{
	//第一个三角形
	0.5, 0.5, 0.0, //右上
	0.5, -0.5, 0.0, //右下
	-0.5, -0.5, 0.0, //左下

	//第二个三角形
	-0.5, -0.5, 0.0, //左下
	0.5, 0.5, 0.0, //右上
	-0.5, 0.5, 0.0, //左上
}

//索引数据(注意这里是从0开始的)
var indices = []uint32{
	0, 1, 5, //第一个三角形
	1, 2, 5, //第二个三角形
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
	window, err := glfw.CreateWindow(screen_width, screen_height, "Triangle", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}
	// 指定当前视口尺寸(前两个参数为左下角位置，后两个参数是渲染窗口宽、高)
	gl.Viewport(0, 0, screen_width, screen_height)

	//生成并绑定VAO和VBO
	var vertex_array_object uint32 // VAO
	gl.GenVertexArrays(1, &vertex_array_object)
	gl.BindVertexArray(vertex_array_object)

	var vertex_buffer_object uint32 // VBO
	gl.GenBuffers(1, &vertex_buffer_object)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertex_buffer_object)
	// 将顶点数据绑定至当前默认的缓冲中
	gl.BufferData(gl.ARRAY_BUFFER, len(triangle)*4, gl.Ptr(triangle), gl.STATIC_DRAW)

	var element_buffer_object uint32
	gl.GenBuffers(1, &element_buffer_object)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, element_buffer_object)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
	// 设置顶点属性指针
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	// 解绑VAO和VBO
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// 生成并编译着色器
	// 顶点着色器
	var vertex_shader = gl.CreateShader(gl.VERTEX_SHADER)
	csovert, free := gl.Strs(vertex_shader_source)
	gl.ShaderSource(vertex_shader, 1, csovert, nil)
	free()
	gl.CompileShader(vertex_shader)
	var success int32
	// 检查着色器是否成功编译，如果编译失败，打印错误信息
	gl.GetShaderiv(vertex_shader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLength1 int32
		gl.GetShaderiv(vertex_shader, gl.INFO_LOG_LENGTH, &logLength1)

		log1 := strings.Repeat("\x00", int(logLength1+1))
		gl.GetShaderInfoLog(vertex_shader, logLength1, nil, gl.Str(log1))
		fmt.Println("ERROR::SHADER::VERTEX::COMPILATION_FAILED", gl.Str(log1))
	}
	// 片段着色器
	fragment_shader := gl.CreateShader(gl.FRAGMENT_SHADER)
	csofragm, free := gl.Strs(fragment_shader_source)
	gl.ShaderSource(fragment_shader, 1, csofragm, nil)
	free()
	gl.CompileShader(fragment_shader)
	// 检查着色器是否成功编译，如果编译失败，打印错误信息
	gl.GetShaderiv(fragment_shader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLength2 int32
		gl.GetShaderiv(fragment_shader, gl.INFO_LOG_LENGTH, &logLength2)

		log2 := strings.Repeat("\x00", int(logLength2+1))
		gl.GetShaderInfoLog(fragment_shader, logLength2, nil, gl.Str(log2))
		fmt.Println("ERROR::SHADER::FRAGMENT::COMPILATION_FAILED", gl.Str(log2))
	}
	shader_program := gl.CreateProgram()
	gl.AttachShader(shader_program, vertex_shader)
	gl.AttachShader(shader_program, fragment_shader)
	gl.LinkProgram(shader_program)
	// 检查着色器是否成功链接，如果链接失败，打印错误信息
	gl.GetProgramiv(shader_program, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength3 int32
		gl.GetShaderiv(shader_program, gl.INFO_LOG_LENGTH, &logLength3)

		log3 := strings.Repeat("\x00", int(logLength3+1))
		gl.GetShaderInfoLog(shader_program, logLength3, nil, gl.Str(log3))
		fmt.Println("ERROR::SHADER::FRAGMENT::LINKING_FAILED", gl.Str(log3))
	}
	// 删除着色器
	gl.DeleteShader(vertex_shader)
	gl.DeleteShader(fragment_shader)
	// 渲染循环
	for !window.ShouldClose() {
		// 清空颜色缓冲
		gl.ClearColor(1.0, 1.0, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// 使用着色器程序
		gl.UseProgram(shader_program)
		// 绘制四边形
		gl.BindVertexArray(vertex_array_object)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		gl.BindVertexArray(0)

		// 交换缓冲并且检查是否有触发事件(比如键盘输入、鼠标移动等）
		window.SwapBuffers()
		glfw.PollEvents()

	}
	// 删除VAO和VBO
	gl.DeleteVertexArrays(1, &vertex_array_object)
	gl.DeleteBuffers(1, &vertex_buffer_object)
	gl.DeleteBuffers(1, &element_buffer_object)
}
