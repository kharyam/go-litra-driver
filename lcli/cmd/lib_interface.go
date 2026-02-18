package cmd

import (
	"github.com/kharyam/go-litra-driver/config"
	"github.com/kharyam/go-litra-driver/lib"
)

// LitraLib defines the interface for the lib package functions used by the commands.
type LitraLib interface {
	LightOn(deviceIndex int)
	LightOff(deviceIndex int)
	LightBrightness(deviceIndex int, level int)
	LightBrightDown(deviceIndex int, inc int)
	LightBrightUp(deviceIndex int, inc int)
	LightTemperature(deviceIndex int, temp uint16)
	LightTempDown(deviceIndex int, inc int)
	LightTempUp(deviceIndex int, inc int)
	ReadCurrentState(deviceIndex int) (brightness int, temperature int, power int)
	ListDevices() []lib.DiscoveredDevice
}

// DefaultLitraLib is the default implementation of the LitraLib interface using the actual lib package.
type DefaultLitraLib struct{}

func (l *DefaultLitraLib) ReadCurrentState(deviceIndex int) (brightness int, temperature int, power int) {
	return config.ReadCurrentState(deviceIndex)
}

func (l *DefaultLitraLib) LightOn(deviceIndex int) {
	lib.LightOn(deviceIndex)
}

func (l *DefaultLitraLib) LightOff(deviceIndex int) {
	lib.LightOff(deviceIndex)
}

func (l *DefaultLitraLib) LightBrightness(deviceIndex int, level int) {
	lib.LightBrightness(deviceIndex, level)
}

func (l *DefaultLitraLib) LightBrightDown(deviceIndex int, inc int) {
	lib.LightBrightDown(deviceIndex, inc)
}

func (l *DefaultLitraLib) LightBrightUp(deviceIndex int, inc int) {
	lib.LightBrightUp(deviceIndex, inc)
}

func (l *DefaultLitraLib) LightTemperature(deviceIndex int, temp uint16) {
	lib.LightTemperature(deviceIndex, temp)
}

func (l *DefaultLitraLib) LightTempDown(deviceIndex int, inc int) {
	lib.LightTempDown(deviceIndex, inc)
}

func (l *DefaultLitraLib) LightTempUp(deviceIndex int, inc int) {
	lib.LightTempUp(deviceIndex, inc)
}

func (l *DefaultLitraLib) ListDevices() []lib.DiscoveredDevice {
	return lib.ListDevices()
}

// libImpl is the variable that will hold the implementation of the LitraLib interface.
// It is initialized with the default implementation.
var libImpl LitraLib = &DefaultLitraLib{}
