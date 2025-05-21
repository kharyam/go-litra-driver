package lib

import (
	"github.com/kharyam/go-litra-driver/config"
	"github.com/sstallion/go-hid"
)

// HIDDevice is an interface for HID device operations
type HIDDevice interface {
	Write(data []byte) (int, error)
	Close() error
}

// HIDEnumerator is an interface for HID enumeration operations
type HIDEnumerator interface {
	Enumerate(vendorID uint16, productID uint16, enumerationCallback func(*hid.DeviceInfo) error) error
}

// HIDOpener is an interface for opening HID devices
type HIDOpener interface {
	Open(vendorID uint16, productID uint16, serialNumber string) (HIDDevice, error)
}

// ConfigUpdater is an interface for updating config state
type ConfigUpdater interface {
	UpdateCurrentState(brightness int, temperature int, power int)
	ReadCurrentState() (brightness int, temperature int, power int)
}

// Default implementations
type defaultHIDEnumeratorImpl struct{}

func (e *defaultHIDEnumeratorImpl) Enumerate(vendorID uint16, productID uint16, enumerationCallback func(*hid.DeviceInfo) error) error {
	return hid.Enumerate(vendorID, productID, enumerationCallback)
}

type defaultHIDOpenerImpl struct{}

func (o *defaultHIDOpenerImpl) Open(vendorID uint16, productID uint16, serialNumber string) (HIDDevice, error) {
	return hid.Open(vendorID, productID, serialNumber)
}

type defaultConfigUpdaterImpl struct{}

func (c *defaultConfigUpdaterImpl) UpdateCurrentState(brightness int, temperature int, power int) {
	config.UpdateCurrentState(brightness, temperature, power)
}

func (c *defaultConfigUpdaterImpl) ReadCurrentState() (brightness int, temperature int, power int) {
	return config.ReadCurrentState()
}

// Default instances
var defaultHIDEnumerator HIDEnumerator = &defaultHIDEnumeratorImpl{}
var defaultHIDOpener HIDOpener = &defaultHIDOpenerImpl{}
var defaultConfigUpdater ConfigUpdater = &defaultConfigUpdaterImpl{}
