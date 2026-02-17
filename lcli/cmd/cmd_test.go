package cmd

import (
	"testing"

	"github.com/kharyam/go-litra-driver/lib"
	"github.com/stretchr/testify/mock"
)

// MockLib is a mock implementation of the lib package functions used by the commands.
type MockLib struct {
	mock.Mock
}

func (m *MockLib) ReadCurrentState(deviceIndex int) (brightness int, temperature int, power int) {
	args := m.Called(deviceIndex)
	return args.Int(0), args.Int(1), args.Int(2)
}

func (m *MockLib) LightOn(deviceIndex int) {
	m.Called(deviceIndex)
}

func (m *MockLib) LightOff(deviceIndex int) {
	m.Called(deviceIndex)
}

func (m *MockLib) LightBrightness(deviceIndex int, level int) {
	m.Called(deviceIndex, level)
}

func (m *MockLib) LightBrightDown(deviceIndex int, inc int) {
	m.Called(deviceIndex, inc)
}

func (m *MockLib) LightBrightUp(deviceIndex int, inc int) {
	m.Called(deviceIndex, inc)
}

func (m *MockLib) LightTemperature(deviceIndex int, temp uint16) {
	m.Called(deviceIndex, temp)
}

func (m *MockLib) LightTempDown(deviceIndex int, inc int) {
	m.Called(deviceIndex, inc)
}

func (m *MockLib) LightTempUp(deviceIndex int, inc int) {
	m.Called(deviceIndex, inc)
}

func (m *MockLib) ListDevices() []lib.DiscoveredDevice {
	args := m.Called()
	return args.Get(0).([]lib.DiscoveredDevice)
}

// TestOnCmd_Run tests the Run function of the onCmd.
func TestOnCmd_Run(t *testing.T) {
	mockLib := new(MockLib)
	originalLibImpl := libImpl
	libImpl = mockLib
	originalDeviceIndex := deviceIndex
	deviceIndex = 0
	defer func() {
		libImpl = originalLibImpl
		deviceIndex = originalDeviceIndex
	}()

	mockLib.On("LightOn", 0).Return().Once()
	onCmd.Run(onCmd, []string{})
	mockLib.AssertExpectations(t)
}

// TestOnCmd_RunWithDevice tests the Run function of the onCmd with a specific device.
func TestOnCmd_RunWithDevice(t *testing.T) {
	mockLib := new(MockLib)
	originalLibImpl := libImpl
	libImpl = mockLib
	originalDeviceIndex := deviceIndex
	deviceIndex = 2
	defer func() {
		libImpl = originalLibImpl
		deviceIndex = originalDeviceIndex
	}()

	mockLib.On("LightOn", 2).Return().Once()
	onCmd.Run(onCmd, []string{})
	mockLib.AssertExpectations(t)
}

// TestOffCmd_Run tests the Run function of the offCmd.
func TestOffCmd_Run(t *testing.T) {
	mockLib := new(MockLib)
	originalLibImpl := libImpl
	libImpl = mockLib
	originalDeviceIndex := deviceIndex
	deviceIndex = 0
	defer func() {
		libImpl = originalLibImpl
		deviceIndex = originalDeviceIndex
	}()

	mockLib.On("LightOff", 0).Return().Once()
	offCmd.Run(offCmd, []string{})
	mockLib.AssertExpectations(t)
}

// TestBrightCmd_Run tests the Run function of the brightCmd.
func TestBrightCmd_Run(t *testing.T) {
	t.Run("ValidLevel", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		level := 50
		mockLib.On("LightBrightness", 0, level).Return().Once()
		brightCmd.Run(brightCmd, []string{"50"})
		mockLib.AssertExpectations(t)
	})

	t.Run("InvalidLevel", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		invalidLevel := 150
		mockLib.On("LightBrightness", 0, invalidLevel).Return().Unset()
		brightCmd.Run(brightCmd, []string{"150"})
		mockLib.AssertExpectations(t)
	})

	t.Run("NonIntegerArgument", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		brightCmd.Run(brightCmd, []string{"abc"})
		mockLib.AssertNotCalled(t, "LightBrightness", mock.Anything, mock.Anything)
		mockLib.AssertExpectations(t)
	})
}

// TestBrightDownCmd_Run tests the Run function of the brightdownCmd.
func TestBrightDownCmd_Run(t *testing.T) {
	t.Run("ValidDecrement", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		inc := 5
		mockLib.On("LightBrightDown", 0, inc).Return().Once()
		brightdownCmd.Run(brightdownCmd, []string{"5"})
		mockLib.AssertExpectations(t)
	})

	t.Run("NoArgument", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		mockLib.AssertNotCalled(t, "LightBrightDown", mock.Anything, mock.Anything)
		brightdownCmd.Run(brightdownCmd, []string{})
		mockLib.AssertExpectations(t)
	})

	t.Run("InvalidArgument", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		mockLib.AssertNotCalled(t, "LightBrightDown", mock.Anything, mock.Anything)
		brightdownCmd.Run(brightdownCmd, []string{"abc"})
		mockLib.AssertExpectations(t)
	})
}

// TestBrightUpCmd_Run tests the Run function of the brightupCmd.
func TestBrightUpCmd_Run(t *testing.T) {
	t.Run("ValidIncrement", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		inc := 5
		mockLib.On("LightBrightUp", 0, inc).Return().Once()
		brightupCmd.Run(brightupCmd, []string{"5"})
		mockLib.AssertExpectations(t)
	})

	t.Run("NoArgument", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		brightupCmd.Run(brightupCmd, []string{})
		mockLib.AssertNotCalled(t, "LightBrightUp", mock.Anything, mock.Anything)
		mockLib.AssertExpectations(t)
	})

	t.Run("InvalidArgument", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		brightupCmd.Run(brightupCmd, []string{"abc"})
		mockLib.AssertNotCalled(t, "LightBrightUp", mock.Anything, mock.Anything)
		mockLib.AssertExpectations(t)
	})
}

// TestTempCmd_Run tests the Run function of the tempCmd.
func TestTempCmd_Run(t *testing.T) {
	t.Run("ValidLevel", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		temp := uint16(4000)
		mockLib.On("LightTemperature", 0, temp).Return().Once()
		tempCmd.Run(tempCmd, []string{"4000"})
		mockLib.AssertExpectations(t)
	})

	t.Run("InvalidLevelBelow", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		tempCmd.Run(tempCmd, []string{"2000"})
		mockLib.AssertNotCalled(t, "LightTemperature", mock.Anything, mock.Anything)
		mockLib.AssertExpectations(t)
	})

	t.Run("InvalidLevelAbove", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		tempCmd.Run(tempCmd, []string{"7000"})
		mockLib.AssertNotCalled(t, "LightTemperature", mock.Anything, mock.Anything)
		mockLib.AssertExpectations(t)
	})

	t.Run("NonIntegerArgument", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		tempCmd.Run(tempCmd, []string{"abc"})
		mockLib.AssertNotCalled(t, "LightTemperature", mock.Anything, mock.Anything)
		mockLib.AssertExpectations(t)
	})
}

// TestTempDownCmd_Run tests the Run function of the tempdownCmd.
func TestTempDownCmd_Run(t *testing.T) {
	t.Run("ValidDecrement", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		inc := 100
		mockLib.On("LightTempDown", 0, inc).Return().Once()
		tempdownCmd.Run(tempdownCmd, []string{"100"})
		mockLib.AssertExpectations(t)
	})

	t.Run("InvalidDecrementZero", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		tempdownCmd.Run(tempdownCmd, []string{"0"})
		mockLib.AssertNotCalled(t, "LightTempDown", mock.Anything, mock.Anything)
		mockLib.AssertExpectations(t)
	})

	t.Run("InvalidDecrementNegative", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		tempdownCmd.Run(tempdownCmd, []string{"-50"})
		mockLib.AssertNotCalled(t, "LightTempDown", mock.Anything, mock.Anything)
		mockLib.AssertExpectations(t)
	})

	t.Run("NonIntegerArgument", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		tempdownCmd.Run(tempdownCmd, []string{"abc"})
		mockLib.AssertNotCalled(t, "LightTempDown", mock.Anything, mock.Anything)
		mockLib.AssertExpectations(t)
	})
}

// TestTempUpCmd_Run tests the Run function of the tempupCmd.
func TestTempUpCmd_Run(t *testing.T) {
	t.Run("ValidIncrement", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		inc := 100
		mockLib.On("LightTempUp", 0, inc).Return().Once()
		tempupCmd.Run(tempupCmd, []string{"100"})
		mockLib.AssertExpectations(t)
	})

	t.Run("InvalidIncrementZero", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		tempupCmd.Run(tempupCmd, []string{"0"})
		mockLib.AssertNotCalled(t, "LightTempUp", mock.Anything, mock.Anything)
		mockLib.AssertExpectations(t)
	})

	t.Run("InvalidIncrementNegative", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		tempupCmd.Run(tempupCmd, []string{"-50"})
		mockLib.AssertNotCalled(t, "LightTempUp", mock.Anything, mock.Anything)
		mockLib.AssertExpectations(t)
	})

	t.Run("NonIntegerArgument", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		tempupCmd.Run(tempupCmd, []string{"abc"})
		mockLib.AssertNotCalled(t, "LightTempUp", mock.Anything, mock.Anything)
		mockLib.AssertExpectations(t)
	})
}

// TestToggleCmd_Run tests the Run function of the toggleCmd.
func TestToggleCmd_Run(t *testing.T) {
	t.Run("LightInitiallyOn", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		mockLib.On("ReadCurrentState", 0).Return(0, 0, 1).Once()
		mockLib.On("LightOff", 0).Return().Once()
		toggleCmd.Run(toggleCmd, []string{})
		mockLib.AssertExpectations(t)
	})

	t.Run("LightInitiallyOff", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		originalDeviceIndex := deviceIndex
		deviceIndex = 0
		defer func() {
			libImpl = originalLibImpl
			deviceIndex = originalDeviceIndex
		}()

		mockLib.On("ReadCurrentState", 0).Return(0, 0, 0).Once()
		mockLib.On("LightOn", 0).Return().Once()
		toggleCmd.Run(toggleCmd, []string{})
		mockLib.AssertExpectations(t)
	})
}

// TestDevicesCmd_Run tests the Run function of the devicesCmd.
func TestDevicesCmd_Run(t *testing.T) {
	mockLib := new(MockLib)
	originalLibImpl := libImpl
	libImpl = mockLib
	defer func() { libImpl = originalLibImpl }()

	mockLib.On("ListDevices").Return([]lib.DiscoveredDevice{
		{Index: 1, Name: "Beam", Serial: "ABC123", ProductID: 0xc901},
		{Index: 2, Name: "Glow", Serial: "DEF456", ProductID: 0xc900},
	}).Once()

	devicesCmd.Run(devicesCmd, []string{})
	mockLib.AssertExpectations(t)
}
