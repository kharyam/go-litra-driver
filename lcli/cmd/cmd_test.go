package cmd

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockLib is a mock implementation of the lib package functions used by the commands.
type MockLib struct {
	mock.Mock
}

func (m *MockLib) ReadCurrentState() (brightness int, temperature int, power int) {
	args := m.Called()
	return args.Int(0), args.Int(1), args.Int(2)
}

// TestOnCmd_Run tests the Run function of the onCmd.
func TestOnCmd_Run(t *testing.T) {
	// Create a new mock lib
	mockLib := new(MockLib)

	// Replace the global libImpl with the mock for testing
	originalLibImpl := libImpl
	libImpl = mockLib
	defer func() { libImpl = originalLibImpl }() // Restore original implementation after test

	// Expect the LightOn function to be called once
	mockLib.On("LightOn").Return().Once()

	// Execute the command's Run function
	onCmd.Run(onCmd, []string{})

	// Assert that all expectations were met
	mockLib.AssertExpectations(t)
}

// TestOffCmd_Run tests the Run function of the offCmd.
func TestOffCmd_Run(t *testing.T) {
	// Create a new mock lib
	mockLib := new(MockLib)

	// Replace the global libImpl with the mock for testing
	originalLibImpl := libImpl
	libImpl = mockLib
	defer func() { libImpl = originalLibImpl }() // Restore original implementation after test

	// Expect the LightOff function to be called once
	mockLib.On("LightOff").Return().Once()

	// Execute the command's Run function
	offCmd.Run(offCmd, []string{})

	// Assert that all expectations were met
	mockLib.AssertExpectations(t)
}

// TestBrightCmd_Run tests the Run function of the brightCmd.
func TestBrightCmd_Run(t *testing.T) {
	// Test case 1: Valid brightness level
	t.Run("ValidLevel", func(t *testing.T) {
		// Create a new mock lib
		mockLib := new(MockLib)

		// Replace the global libImpl with the mock for testing
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }() // Restore original implementation after test

		level := 50
		mockLib.On("LightBrightness", level).Return().Once()
		brightCmd.Run(brightCmd, []string{"50"})
		mockLib.AssertExpectations(t)
	})

	// Test case 2: Invalid brightness level (out of range)
	t.Run("InvalidLevel", func(t *testing.T) {
		// Create a new mock lib
		mockLib := new(MockLib)

		// Replace the global libImpl with the mock for testing
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }() // Restore original implementation after test

		invalidLevel := 150
		// We don't expect LightBrightness to be called for invalid input
		mockLib.On("LightBrightness", invalidLevel).Return().Unset()
		brightCmd.Run(brightCmd, []string{"150"})
		mockLib.AssertExpectations(t)
	})

	// Test case 3: Non-integer argument
	t.Run("NonIntegerArgument", func(t *testing.T) {
		// Create a new mock lib
		mockLib := new(MockLib)

		// Replace the global libImpl with the mock for testing
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }() // Restore original implementation after test

		brightCmd.Run(brightCmd, []string{"abc"})

		// We don't expect LightBrightness to be called for invalid input
		mockLib.AssertNotCalled(t, "LightBrightness", mock.Anything)
		mockLib.AssertExpectations(t)
	})
}

// TestBrightDownCmd_Run tests the Run function of the brightdownCmd.
func TestBrightDownCmd_Run(t *testing.T) {
	// Test case 1: Valid decrement amount
	t.Run("ValidDecrement", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		inc := 5
		mockLib.On("LightBrightDown", inc).Return().Once()
		brightdownCmd.Run(brightdownCmd, []string{"5"})
		mockLib.AssertExpectations(t)
	})

	// Test case 2: No argument (should default to 10)
	t.Run("NoArgument", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		// Verify LightBrightDown not called
		mockLib.AssertNotCalled(t, "LightBrightDown", mock.Anything)
		brightdownCmd.Run(brightdownCmd, []string{})
		mockLib.AssertExpectations(t)
	})

	// Test case 3: Invalid argument (non-integer)
	t.Run("InvalidArgument", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		// Verify LightBrightDown not called
		mockLib.AssertNotCalled(t, "LightBrightDown", mock.Anything)
		brightdownCmd.Run(brightdownCmd, []string{"abc"})
		mockLib.AssertExpectations(t)
	})
}

// TestBrightUpCmd_Run tests the Run function of the brightupCmd.
func TestBrightUpCmd_Run(t *testing.T) {
	// Test case 1: Valid increment amount
	t.Run("ValidIncrement", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		inc := 5
		mockLib.On("LightBrightUp", inc).Return().Once()
		brightupCmd.Run(brightupCmd, []string{"5"})
		mockLib.AssertExpectations(t)
	})

	// Test case 2: No argument (should default to 10)
	t.Run("NoArgument", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		brightupCmd.Run(brightupCmd, []string{})
		mockLib.AssertNotCalled(t, "LightBrightUp", mock.Anything)
		mockLib.AssertExpectations(t)
	})

	// Test case 3: Invalid argument (non-integer)
	t.Run("InvalidArgument", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		brightupCmd.Run(brightupCmd, []string{"abc"})
		mockLib.AssertNotCalled(t, "LightBrightUp", mock.Anything)
		mockLib.AssertExpectations(t)
	})
}

// TestTempCmd_Run tests the Run function of the tempCmd.
func TestTempCmd_Run(t *testing.T) {
	// Test case 1: Valid temperature level
	t.Run("ValidLevel", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		temp := uint16(4000)
		mockLib.On("LightTemperature", temp).Return().Once()
		tempCmd.Run(tempCmd, []string{"4000"})
		mockLib.AssertExpectations(t)
	})

	// Test case 2: Invalid temperature level (below range)
	t.Run("InvalidLevelBelow", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		tempCmd.Run(tempCmd, []string{"2000"})
		mockLib.AssertNotCalled(t, "LightTemperature", mock.Anything)
		mockLib.AssertExpectations(t)
	})

	// Test case 3: Invalid temperature level (above range)
	t.Run("InvalidLevelAbove", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		tempCmd.Run(tempCmd, []string{"7000"})
		mockLib.AssertNotCalled(t, "LightTemperature", mock.Anything)
		mockLib.AssertExpectations(t)
	})

	// Test case 4: Non-integer argument
	t.Run("NonIntegerArgument", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		tempCmd.Run(tempCmd, []string{"abc"})
		mockLib.AssertNotCalled(t, "LightTemperature", mock.Anything)
		mockLib.AssertExpectations(t)
	})
}

// TestTempDownCmd_Run tests the Run function of the tempdownCmd.
func TestTempDownCmd_Run(t *testing.T) {
	// Test case 1: Valid decrement amount
	t.Run("ValidDecrement", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		inc := 100
		mockLib.On("LightTempDown", inc).Return().Once()
		tempdownCmd.Run(tempdownCmd, []string{"100"})
		mockLib.AssertExpectations(t)
	})

	// Test case 2: Invalid decrement amount (zero)
	t.Run("InvalidDecrementZero", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		tempdownCmd.Run(tempdownCmd, []string{"0"})
		mockLib.AssertNotCalled(t, "LightTempDown", mock.Anything)
		mockLib.AssertExpectations(t)
	})

	// Test case 3: Invalid decrement amount (negative)
	t.Run("InvalidDecrementNegative", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		tempdownCmd.Run(tempdownCmd, []string{"-50"})
		mockLib.AssertNotCalled(t, "LightTempDown", mock.Anything)
		mockLib.AssertExpectations(t)
	})

	// Test case 4: Non-integer argument
	t.Run("NonIntegerArgument", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		tempdownCmd.Run(tempdownCmd, []string{"abc"})
		mockLib.AssertNotCalled(t, "LightTempDown", mock.Anything)
		mockLib.AssertExpectations(t)
	})
}

// TestTempUpCmd_Run tests the Run function of the tempupCmd.
func TestTempUpCmd_Run(t *testing.T) {
	// Test case 1: Valid increment amount
	t.Run("ValidIncrement", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		inc := 100
		mockLib.On("LightTempUp", inc).Return().Once()
		tempupCmd.Run(tempupCmd, []string{"100"})
		mockLib.AssertExpectations(t)
	})

	// Test case 2: Invalid increment amount (zero)
	t.Run("InvalidIncrementZero", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		tempupCmd.Run(tempupCmd, []string{"0"})
		mockLib.AssertNotCalled(t, "LightTempUp", mock.Anything)
		mockLib.AssertExpectations(t)
	})

	// Test case 3: Invalid increment amount (negative)
	t.Run("InvalidIncrementNegative", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		tempupCmd.Run(tempupCmd, []string{"-50"})
		mockLib.AssertNotCalled(t, "LightTempUp", mock.Anything)
		mockLib.AssertExpectations(t)
	})

	// Test case 4: Non-integer argument
	t.Run("NonIntegerArgument", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		tempupCmd.Run(tempupCmd, []string{"abc"})
		mockLib.AssertNotCalled(t, "LightTempUp", mock.Anything)
		mockLib.AssertExpectations(t)
	})
}

// TestToggleCmd_Run tests the Run function of the toggleCmd.
func TestToggleCmd_Run(t *testing.T) {
	// Test case 1: Light is initially on
	t.Run("LightInitiallyOn", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		// Mock ReadCurrentState to return power = 1 (on)
		mockLib.On("ReadCurrentState").Return(0, 0, 1).Once()
		// Expect LightOff to be called
		mockLib.On("LightOff").Return().Once()

		toggleCmd.Run(toggleCmd, []string{})
		mockLib.AssertExpectations(t)
	})

	// Test case 2: Light is initially off
	t.Run("LightInitiallyOff", func(t *testing.T) {
		mockLib := new(MockLib)
		originalLibImpl := libImpl
		libImpl = mockLib
		defer func() { libImpl = originalLibImpl }()

		// Mock ReadCurrentState to return power = 0 (off)
		mockLib.On("ReadCurrentState").Return(0, 0, 0).Once()
		// Expect LightOn to be called
		mockLib.On("LightOn").Return().Once()

		toggleCmd.Run(toggleCmd, []string{})
		mockLib.AssertExpectations(t)
	})
}

func (m *MockLib) LightOn() {
	m.Called()
}

func (m *MockLib) LightOff() {
	m.Called()
}

func (m *MockLib) LightBrightness(level int) {
	m.Called(level)
}

func (m *MockLib) LightBrightDown(inc int) {
	m.Called(inc)
}

func (m *MockLib) LightBrightUp(inc int) {
	m.Called(inc)
}

func (m *MockLib) LightTemperature(temp uint16) {
	m.Called(temp)
}

func (m *MockLib) LightTempDown(inc int) {
	m.Called(inc)
}

func (m *MockLib) LightTempUp(inc int) {
	m.Called(inc)
}

// TODO: Implement tests for each command
