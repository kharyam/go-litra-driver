// Package lib defines a library for accessing the functionality of the
// Logitech Litra Glow and Logitech Litra Beam via USB
package lib

import (
	"encoding/binary"
	"math"
	"os"

	"github.com/google/gousb"
	"github.com/kharyam/go-litra-driver/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const VendorId = 0x046d
const LightOffCode = 0x00
const LightOnCode = 0x01
const MinBrightness = 0x14
const MaxBrightness = 0xfa
const TimeoutMs = 3000

type litraDevice struct {
	name         string
	productId    uint
	endpoint     uint
	bufferLength uint
}

var litraProducts = [2]litraDevice{
	{"Glow",
		0xc900,
		0x02,
		64},
	{"Beam",
		0xc901,
		0x01,
		32},
}

// var devices []*gousb.Device
// var context *gousb.Context

func findDevices() ([]*gousb.Device, *gousb.Context) {

	// Initialize a new Context.
	var context = gousb.NewContext()
	var devices []*gousb.Device
	log.Debug().Msg("Searching for litra devices...")
	for i := 0; i < len(litraProducts); i++ {

		devs, err := context.OpenDevices(func(desc *gousb.DeviceDesc) bool {
			if desc.Vendor == gousb.ID(VendorId) && desc.Product == gousb.ID(litraProducts[i].productId) {
				return true
			}
			return false
		})

		if err != nil {
			log.Fatal().Msgf("Could not open a device: %v", err)
		}

		for j := 0; j < len(devs); j++ {
			log.Debug().Msgf("Found a %s device!!!", litraProducts[i].name)
			devices = append(devices, devs[j])
		}

	}

	return devices, context
}

func getEndpoint(device *gousb.Device) uint {

	var result uint = 0x0
	for i := 0; i < len(litraProducts); i++ {
		if device.Desc.Product == gousb.ID(litraProducts[i].productId) {
			result = litraProducts[i].endpoint
			break
		}
	}

	return result
}

func getBufferLength(device *gousb.Device) uint {

	var result uint = 0x0
	for i := 0; i < len(litraProducts); i++ {
		if device.Desc.Product == gousb.ID(litraProducts[i].productId) {
			result = litraProducts[i].bufferLength
			break
		}
	}

	return result
}

func commandDevices(bytes []byte) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	var devices, context = findDevices()

	defer context.Close()

	for i := 0; i < len(devices); i++ {
		var device = devices[i]
		device.SetAutoDetach(true)
		defer device.Close()
		var endpoint = getEndpoint(device)

		cfg, err := device.Config(1)
		if err != nil {
			log.Fatal().Msgf("%s.Config(1): %v", device, err)
		}
		defer cfg.Close()

		intf, err := cfg.Interface(0, 0)
		if err != nil {
			log.Fatal().Msgf("%s.Interface(0, 0): %v", cfg, err)
		}
		defer intf.Close()

		epOut, err := intf.OutEndpoint(int(endpoint))
		if err != nil {
			log.Fatal().Msgf("%s.OutEndpoint(%d): %v", intf, endpoint, err)
		}

		epIn, err := intf.InEndpoint(int(endpoint))
		if err != nil {
			log.Fatal().Msgf("%s.InEndpoint(%d): %v", intf, endpoint, err)
		}

		writeBytes, err := epOut.Write(bytes)
		if err != nil {
			log.Fatal().Msgf("Failed writing bytes %v %v", err, writeBytes)
		}

		buffer := make([]byte, getBufferLength(device))
		epIn.Read(buffer)
	}

}

// LightOn turns on all detected lights
func LightOn() {
	var bytes = []byte{0x11, 0xff, 0x04, 0x1c, LightOnCode, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	commandDevices(bytes)
}

// LightOff turns off all detected lights
func LightOff() {
	var bytes = []byte{0x11, 0xff, 0x04, 0x1c, LightOffCode, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	commandDevices(bytes)
}

// LightBrightness sets the brightness of all connected lights. Specify a brightness between 0 and 100
func LightBrightness(level int) {
	var adjustedLevel = MinBrightness + math.Floor((float64(level)/float64(100))*(MaxBrightness-MinBrightness))

	var bytes = []byte{0x11, 0xff, 0x04, 0x4c, 0x00, byte(adjustedLevel), 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	commandDevices(bytes)
	config.UpdateCurrentState(level, -1)
}

// LightTemperature sets a light temperature between 2700 and 6500
func LightTemperature(temp uint16) {

	log.Info().Msgf("%d", temp)
	tempBytes := make([]byte, 2)

	binary.BigEndian.PutUint16(tempBytes, temp)

	log.Info().Msgf("%d %d", tempBytes[0], tempBytes[1])
	var bytes = []byte{0x11, 0xff, 0x04, 0x9c, tempBytes[0], tempBytes[1], 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	commandDevices(bytes)
	config.UpdateCurrentState(-1, int(temp))
}
