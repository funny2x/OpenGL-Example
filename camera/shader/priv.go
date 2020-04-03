//私有的对外包不可见的

package shader

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

//getShaderFromFile 从文件中获取shader源码
func getShaderFromFile(file string) (string, error) {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Load Shader File Error!")
		return "", err
	}
	return string(src), nil
}

//generateCompileShader 生成并编译着色器
//shadercode 着色器源码
//sType 编译着色器类型,如:gl.VERTEX_SHADER,gl.FRAGMENT_SHADER
func generateCompileShader(shadercode string, sType uint32) (uint32, error) {
	handle := gl.CreateShader(sType)
	glSrc, freeFn := gl.Strs(shadercode + "\x00")
	gl.ShaderSource(handle, 1, glSrc, nil)
	freeFn()
	gl.CompileShader(handle)
	var failMsg string
	if sType == gl.VERTEX_SHADER {
		failMsg = "ERROR::SHADER::VERTEX::COMPILATION_FAILED"
	} else if sType == gl.FRAGMENT_SHADER {
		failMsg = "ERROR::SHADER::FRAGMENT::COMPILATION_FAILED"
	}
	if err := getGlError(handle, failMsg); err != nil {
		return 0, err
	}
	return handle, nil
}

//getGlError 检查着色器是否成功编译，如果编译失败，打印错误信息
func getGlError(handle uint32, failMsg string) error {
	var success int32
	gl.GetShaderiv(handle, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLength1 int32
		gl.GetShaderiv(handle, gl.INFO_LOG_LENGTH, &logLength1)

		log1 := strings.Repeat("\x00", int(logLength1+1))
		gl.GetShaderInfoLog(handle, logLength1, nil, gl.Str(log1))

		return fmt.Errorf("%s: %v", failMsg, log1)
	}
	return nil
}

//linkShader 链接生成着色器程序
func linkShader(vertexShader, fragmentShader uint32) (uint32, error) {
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)
	err := checkCompileErrors(shaderProgram)
	return shaderProgram, err
}

func checkCompileErrors(shaderProgram uint32) error {
	var success int32
	// 检查着色器是否成功链接，如果链接失败，打印错误信息
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength3 int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength3)

		log3 := strings.Repeat("\x00", int(logLength3+1))
		gl.GetProgramInfoLog(shaderProgram, logLength3, nil, gl.Str(log3))

		return fmt.Errorf("%s: %v", "ERROR::SHADER::FRAGMENT::LINKING_FAILED", log3)
	}
	return nil
}

func (s *Shader) getUniform(name string) int32 {
	position := gl.GetUniformLocation(s.id, gl.Str(name+"\x00"))
	if position == -1 {
		fmt.Println("uniform ", name, " set failed!")
	}
	return position
}
