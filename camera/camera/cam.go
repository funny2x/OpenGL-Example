/*
摄像机
*/

package camera

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// 默认的 camera 值
const (
	YAW         = -90.0 //默认的偏航角
	PITCH       = 0.0   //默认的俯仰角
	SPEED       = 2.5   //默认的速度
	SENSITIVITY = 0.1   //默认的灵敏度
	ZOOM        = 45.0  //默认的视角
)

//摄像机移动方向
const (
	FORWARD = iota
	BACKWARD
	LEFT
	RIGHT
)

//Camera 摄像机
type Camera struct {
	Position mgl32.Vec3 //摄像头位置
	Front    mgl32.Vec3 //摄像头方向
	Up       mgl32.Vec3 //摄像头垂直上方向
	Right    mgl32.Vec3 //摄像头右方向
	WorldUp  mgl32.Vec3 //正上方向

	Yaw   float32
	Pitch float32

	MovementSpeed    float32
	MouseSensitivity float32
	Zoom             float32
}

//GetCamera Camera的构造函数
//pos=mgl32.Vec3{0.0,0.0,0.0}
//默认
//up=mgl32.Vec3{0.0,1.0,0.0}
//yaw=YAW,pitch=PITCH float32
func GetCamera(pos mgl32.Vec3) *Camera {
	cam := &Camera{
		Front:            mgl32.Vec3{0.0, 0.0, -1.0},
		MovementSpeed:    SPEED,
		MouseSensitivity: SENSITIVITY,
		Zoom:             ZOOM,
		Position:         pos,
		WorldUp:          mgl32.Vec3{0.0, 1.0, 0.0},
		Yaw:              YAW,
		Pitch:            PITCH,
	}
	cam.updateCameraVectors()
	return cam
}

//GetViewMatrix ->mgl32.LookAtV(c.Position,c.Position.Add(c.Front),c.Up)
func (c *Camera) GetViewMatrix() mgl32.Mat4 {

	return mgl32.LookAtV(c.Position, c.Position.Add(c.Front), c.Up)
}

//ProcessKeyboard 对应键盘移动事件
func (c *Camera) ProcessKeyboard(direction uint32, deltaTime float64) {
	velocity := c.MovementSpeed * float32(deltaTime)
	switch direction {
	case FORWARD:
		c.Position = c.Position.Add(c.Front.Mul(velocity))
	case BACKWARD:
		c.Position = c.Position.Sub(c.Front.Mul(velocity))
	case LEFT:
		c.Position = c.Position.Sub(c.Right.Mul(velocity))
	case RIGHT:
		c.Position = c.Position.Add(c.Right.Mul(velocity))
	}
}

//ProcessMouseMovement 对应鼠标移动事件
func (c *Camera) ProcessMouseMovement(xoffset, yoffset float64, constrainPitch bool) {
	xoffset *= float64(c.MouseSensitivity)
	yoffset *= float64(c.MouseSensitivity)

	c.Yaw += float32(xoffset)
	c.Pitch += float32(yoffset)

	if constrainPitch {
		if c.Pitch > 89.0 {
			c.Pitch = 89.0
		} else if c.Pitch < -89.0 {
			c.Pitch = -89.0
		}
	}
	c.updateCameraVectors()
}

//ProcessMouseScroll 对应鼠标滚轮事件
func (c *Camera) ProcessMouseScroll(yoffset float64) {
	if c.Zoom >= 1.0 && c.Zoom <= 45.0 {
		c.Zoom -= float32(yoffset)
	}
	if c.Zoom <= 1.0 {
		c.Zoom = 1.0
	}
	if c.Zoom >= 45.0 {
		c.Zoom = 45.0
	}
}

func (c *Camera) updateCameraVectors() {
	x := math.Cos(float64(mgl32.DegToRad(c.Yaw))) * math.Cos(float64(mgl32.DegToRad(c.Pitch)))
	y := math.Sin(float64(mgl32.DegToRad(c.Pitch)))
	z := math.Sin(float64(mgl32.DegToRad(c.Yaw))) * math.Cos(float64(mgl32.DegToRad(c.Pitch)))
	var front = mgl32.Vec3{float32(x), float32(y), float32(z)}
	c.Front = front.Normalize()

	c.Right = c.Front.Cross(c.WorldUp).Normalize()
	c.Up = c.Right.Cross(c.Front).Normalize()
}
