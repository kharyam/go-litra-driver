package cmd

import (
	"github.com/kharyam/go-litra-driver/config"
	"github.com/kharyam/go-litra-driver/lib"
)

// LitraLib defines the interface for the lib package functions used by the commands.
type LitraLib interface {
	LightOn()
	LightOff()
	LightBrightness(level int)
	LightBrightDown(inc int)
	LightBrightUp(inc int)
	LightTemperature(temp uint16)
	LightTempDown(inc int)
	LightTempUp(inc int)
	ReadCurrentState() (brightness int, temperature int, power int)
}

// DefaultLitraLib is the default implementation of the LitraLib interface using the actual lib package.
type DefaultLitraLib struct{}

func (l *DefaultLitraLib) ReadCurrentState() (brightness int, temperature int, power int) {
	return config.ReadCurrentState()
}

func (l *DefaultLitraLib) LightOn() {
	lib.LightOn()
}

func (l *DefaultLitraLib) LightOff() {
	lib.LightOff()
}

func (l *DefaultLitraLib) LightBrightness(level int) {
	lib.LightBrightness(level)
}

func (l *DefaultLitraLib) LightBrightDown(inc int) {
	lib.LightBrightDown(inc)
}

func (l *DefaultLitraLib) LightBrightUp(inc int) {
	lib.LightBrightUp(inc)
}

func (l *DefaultLitraLib) LightTemperature(temp uint16) {
	lib.LightTemperature(temp)
}

func (l *DefaultLitraLib) LightTempDown(inc int) {
	lib.LightTempDown(inc)
}

func (l *DefaultLitraLib) LightTempUp(inc int) {
	lib.LightTempUp(inc)
}

// libImpl is the variable that will hold the implementation of the LitraLib interface.
// It is initialized with the default implementation.
var libImpl LitraLib = &DefaultLitraLib{}
