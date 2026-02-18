// Package lib defines a library for accessing the functionality of the
// Logitech Litra Glow and Logitech Litra Beam via USB
package lib

import (
	"encoding/binary"
	"math"
	"os"
	"sort"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sstallion/go-hid"
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

// discoveredDeviceInternal pairs an opened HID device handle with its metadata
type discoveredDeviceInternal struct {
	device   HIDDevice
	metadata DiscoveredDevice
}

// findDevicesWithDefaults finds all connected Litra devices using default implementations.
// Devices are sorted by serial number for deterministic ordering and assigned 1-based indices.
func findDevicesWithDefaults() []discoveredDeviceInternal {
	var deviceInfos = make(map[string]*hid.DeviceInfo)
	var productNames = make(map[string]string)

	for i := 0; i < len(litraProducts); i++ {
		productName := litraProducts[i].name
		defaultHIDEnumerator.Enumerate(VendorId, uint16(litraProducts[i].productId), func(info *hid.DeviceInfo) error {
			deviceInfos[info.SerialNbr] = info
			productNames[info.SerialNbr] = productName
			return nil
		})
	}

	// Sort serials for deterministic ordering
	serials := make([]string, 0, len(deviceInfos))
	for serial := range deviceInfos {
		serials = append(serials, serial)
	}
	sort.Strings(serials)

	var devices []discoveredDeviceInternal
	for idx, serial := range serials {
		info := deviceInfos[serial]
		device, err := defaultHIDOpener.Open(info.VendorID, info.ProductID, info.SerialNbr)
		if firstRun {
			log.Debug().Msgf("Found device %s (serial: %s)", productNames[serial], serial)
		}
		if err == nil {
			devices = append(devices, discoveredDeviceInternal{
				device: device,
				metadata: DiscoveredDevice{
					Index:     idx + 1,
					Name:      productNames[serial],
					Serial:    serial,
					ProductID: info.ProductID,
				},
			})
		} else {
			log.Error().Msgf("ERROR %v", err)
		}
	}

	firstRun = false
	return devices
}

// commandDevices sends a command to connected devices.
// deviceIndex 0 writes to all devices, deviceIndex > 0 writes only to the matching device.
func commandDevices(bytes []byte, deviceIndex int) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	var devices = findDevicesWithDefaults()

	for i := 0; i < len(devices); i++ {
		var d = devices[i]
		defer d.device.Close()
		if deviceIndex == 0 || d.metadata.Index == deviceIndex {
			d.device.Write(bytes)
		}
	}
}

// ListDevices returns all connected Litra devices with their metadata
func ListDevices() []DiscoveredDevice {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	devices := findDevicesWithDefaults()
	result := make([]DiscoveredDevice, len(devices))
	for i, d := range devices {
		defer d.device.Close()
		result[i] = d.metadata
	}
	return result
}

// LightOn turns on detected lights. deviceIndex 0 targets all, 1+ targets a specific device.
func LightOn(deviceIndex int) {
	var bytes = []byte{0x11, 0xff, 0x04, 0x1c, LightOnCode, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	commandDevices(bytes, deviceIndex)
	defaultConfigUpdater.UpdateCurrentState(deviceIndex, -1, -1, 1)
}

// LightOff turns off detected lights. deviceIndex 0 targets all, 1+ targets a specific device.
func LightOff(deviceIndex int) {
	var bytes = []byte{0x11, 0xff, 0x04, 0x1c, LightOffCode, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	commandDevices(bytes, deviceIndex)
	defaultConfigUpdater.UpdateCurrentState(deviceIndex, -1, -1, 0)
}

// LightBrightness sets the brightness of connected lights. Specify a brightness between 0 and 100.
// deviceIndex 0 targets all, 1+ targets a specific device.
func LightBrightness(deviceIndex int, level int) {
	var adjustedLevel = MinBrightness + math.Floor((float64(level)/float64(100))*(MaxBrightness-MinBrightness))

	var bytes = []byte{0x11, 0xff, 0x04, 0x4c, 0x00, byte(adjustedLevel), 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	commandDevices(bytes, deviceIndex)
	defaultConfigUpdater.UpdateCurrentState(deviceIndex, level, -1, -1)
}

// Function variables for testing
var lightBrightnessFunc = LightBrightness
var lightTemperatureFunc = LightTemperature

// LightBrightDown decreases the brightness by the amount specified.
// deviceIndex 0 targets all, 1+ targets a specific device.
func LightBrightDown(deviceIndex int, inc int) {
	brightness, _, _ := defaultConfigUpdater.ReadCurrentState(deviceIndex)
	brightness -= inc

	if brightness < 1 {
		brightness = 0
	}

	lightBrightnessFunc(deviceIndex, brightness)
}

// LightBrightUp increases the brightness by the amount specified.
// deviceIndex 0 targets all, 1+ targets a specific device.
func LightBrightUp(deviceIndex int, inc int) {
	brightness, _, _ := defaultConfigUpdater.ReadCurrentState(deviceIndex)
	brightness += inc

	if brightness > 100 {
		brightness = 100
	}

	lightBrightnessFunc(deviceIndex, brightness)
}

// LightTemperature sets a light temperature between 2700 and 6500.
// deviceIndex 0 targets all, 1+ targets a specific device.
func LightTemperature(deviceIndex int, temp uint16) {
	tempBytes := make([]byte, 2)

	binary.BigEndian.PutUint16(tempBytes, temp)

	var bytes = []byte{0x11, 0xff, 0x04, 0x9c, tempBytes[0], tempBytes[1], 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	commandDevices(bytes, deviceIndex)
	defaultConfigUpdater.UpdateCurrentState(deviceIndex, -1, int(temp), -1)
}

// LightTempDown decreases the temperature by the amount specified.
// deviceIndex 0 targets all, 1+ targets a specific device.
func LightTempDown(deviceIndex int, inc int) {
	_, temp, _ := defaultConfigUpdater.ReadCurrentState(deviceIndex)
	temp -= inc

	if temp < 2700 {
		temp = 2700
	}

	lightTemperatureFunc(deviceIndex, uint16(temp))
}

// LightTempUp increases the temperature by the amount specified.
// deviceIndex 0 targets all, 1+ targets a specific device.
func LightTempUp(deviceIndex int, inc int) {
	_, temp, _ := defaultConfigUpdater.ReadCurrentState(deviceIndex)
	temp += inc

	if temp > 6500 {
		temp = 6500
	}

	lightTemperatureFunc(deviceIndex, uint16(temp))
}
