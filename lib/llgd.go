// Package lib defines a library for accessing the functionality of the
// Logitech Litra Glow and Logitech Litra Beam via USB
package lib

import (
	"encoding/binary"
	"math"
	"os"

	"github.com/karalabe/hid"
	"github.com/kharyam/go-litra-driver/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const VendorId = 0x046d
const LightOffCode = 0x00
const LightOnCode = 0x01
const MinBrightness = 0x14
const MaxBrightness = 0xfa

var firstRun = true

type litraDevice struct {
	name      string
	productId uint
}

var litraProducts = [2]litraDevice{
	{"Glow",
		0xc900},
	{"Beam",
		0xc901},
}

func findDevices() []hid.Device {

	var deviceInfos []hid.DeviceInfo
	var devices []hid.Device

	if !hid.Supported() {
		log.Fatal().Msgf("HID Library reports platform is not supported")
	}

	for i := 0; i < len(litraProducts); i++ {
		deviceInfosList := hid.Enumerate(uint16(VendorId), uint16(litraProducts[i].productId))

		for i := 0; i < len(deviceInfosList); i++ {
			deviceInfos = append(deviceInfos, deviceInfosList[i])
		}
	}

	for _, devInfo := range deviceInfos {
		device, err := devInfo.Open()
		if firstRun {
			// Bug? Serial/Product not populated by hid library
			log.Debug().Msgf("Found device %s", devInfo.Product)
		}
		if err == nil {
			devices = append(devices, *device)
		} else {
			log.Error().Msgf("ERROR %v", err)
		}
	}

	firstRun = false
	return devices
}

func commandDevices(bytes []byte) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	var devices = findDevices()

	for i := 0; i < len(devices); i++ {
		var device = devices[i]
		defer device.Close()
		device.Write(bytes)
	}

}

// LightOn turns on all detected lights
func LightOn() {
	var bytes = []byte{0x11, 0xff, 0x04, 0x1c, LightOnCode, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	commandDevices(bytes)
	config.UpdateCurrentState(-1, -1, 1)
}

// LightOff turns off all detected lights
func LightOff() {
	var bytes = []byte{0x11, 0xff, 0x04, 0x1c, LightOffCode, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	commandDevices(bytes)
	config.UpdateCurrentState(-1, -1, 0)
}

// LightBrightness sets the brightness of all connected lights. Specify a brightness between 0 and 100
func LightBrightness(level int) {
	var adjustedLevel = MinBrightness + math.Floor((float64(level)/float64(100))*(MaxBrightness-MinBrightness))

	var bytes = []byte{0x11, 0xff, 0x04, 0x4c, 0x00, byte(adjustedLevel), 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	commandDevices(bytes)
	config.UpdateCurrentState(level, -1, -1)
}

// LightBrightDown decreases the brightness by the amount specified
func LightBrightDown(inc int) {

	brightness, _, _ := config.ReadCurrentState()
	brightness -= inc

	if brightness < 1 {
		brightness = 0
	}

	LightBrightness(brightness)
}

// LightBrightUp increases the brightness by the amount specified
func LightBrightUp(inc int) {

	brightness, _, _ := config.ReadCurrentState()
	brightness += inc

	if brightness > 100 {
		brightness = 100
	}

	LightBrightness(brightness)
}

// LightTemperature sets a light temperature between 2700 and 6500
func LightTemperature(temp uint16) {

	tempBytes := make([]byte, 2)

	binary.BigEndian.PutUint16(tempBytes, temp)

	var bytes = []byte{0x11, 0xff, 0x04, 0x9c, tempBytes[0], tempBytes[1], 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	commandDevices(bytes)
	config.UpdateCurrentState(-1, int(temp), -1)
}

// LightTempDown decreases the temperature by the amount specified
func LightTempDown(inc int) {

	_, temp, _ := config.ReadCurrentState()
	temp -= inc

	if temp < 2700 {
		temp = 2700
	}

	LightTemperature(uint16(temp))
}

// LightTempUp increases the temperature by the amount specified
func LightTempUp(inc int) {

	_, temp, _ := config.ReadCurrentState()
	temp += inc

	if temp > 6500 {
		temp = 6500
	}

	LightTemperature(uint16(temp))
}
