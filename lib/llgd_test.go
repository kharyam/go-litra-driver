package lib

import (
	"testing"

	"github.com/sstallion/go-hid"
	"github.com/stretchr/testify/mock"
)

// Mock HID device
type MockHIDDevice struct {
	mock.Mock
}

func (m *MockHIDDevice) Write(data []byte) (int, error) {
	args := m.Called(data)
	return args.Int(0), args.Error(1)
}

func (m *MockHIDDevice) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Mock HID enumerator
type MockHIDEnumerator struct {
	mock.Mock
}

func (m *MockHIDEnumerator) Enumerate(vendorID uint16, productID uint16, enumerationCallback func(*hid.DeviceInfo) error) error {
	args := m.Called(vendorID, productID, enumerationCallback)
	return args.Error(0)
}

// Mock HID opener
type MockHIDOpener struct {
	mock.Mock
}

func (m *MockHIDOpener) Open(vendorID uint16, productID uint16, serialNumber string) (HIDDevice, error) {
	args := m.Called(vendorID, productID, serialNumber)
	return args.Get(0).(HIDDevice), args.Error(1)
}

// Mock config updater
type MockConfigUpdater struct {
	mock.Mock
}

func (m *MockConfigUpdater) UpdateCurrentState(deviceIndex int, brightness int, temperature int, power int) {
	m.Called(deviceIndex, brightness, temperature, power)
}

func (m *MockConfigUpdater) ReadCurrentState(deviceIndex int) (brightness int, temperature int, power int) {
	args := m.Called(deviceIndex)
	return args.Int(0), args.Int(1), args.Int(2)
}

// Setup test environment. Returns two mock devices: device1 is Beam (index 1, sorted first),
// device2 is Glow (index 2, sorted second) since serials are sorted alphabetically.
func setupTest() (*MockHIDDevice, *MockHIDDevice, *MockHIDEnumerator, *MockHIDOpener, *MockConfigUpdater, func()) {
	// Save original values
	originalHIDEnumerator := defaultHIDEnumerator
	originalHIDOpener := defaultHIDOpener
	originalConfigUpdater := defaultConfigUpdater
	originalLightBrightnessFunc := lightBrightnessFunc
	originalLightTemperatureFunc := lightTemperatureFunc

	// Create mocks - two separate devices for Beam and Glow
	mockDevice1 := new(MockHIDDevice) // Beam (serial "test-serial-Beam" sorts first)
	mockDevice2 := new(MockHIDDevice) // Glow (serial "test-serial-Glow" sorts second)
	mockEnumerator := new(MockHIDEnumerator)
	mockOpener := new(MockHIDOpener)
	mockConfigUpdater := new(MockConfigUpdater)

	// Replace with mocks
	defaultHIDEnumerator = mockEnumerator
	defaultHIDOpener = mockOpener
	defaultConfigUpdater = mockConfigUpdater

	// Map product names to their mock devices (sorted by serial: Beam first, Glow second)
	mockDevices := map[string]*MockHIDDevice{
		"Beam": mockDevice1,
		"Glow": mockDevice2,
	}

	// Setup mock behavior for findDevices
	for _, product := range litraProducts {
		// Setup the enumerate call to invoke the callback with our device info
		mockEnumerator.On("Enumerate",
			uint16(VendorId),
			uint16(product.productId),
			mock.MatchedBy(func(fn interface{}) bool {
				_, ok := fn.(func(*hid.DeviceInfo) error)
				return ok
			})).Return(nil).
			Run(func(args mock.Arguments) {
				callback := args.Get(2).(func(*hid.DeviceInfo) error)
				deviceInfo := &hid.DeviceInfo{
					VendorID:   uint16(VendorId),
					ProductID:  uint16(product.productId),
					SerialNbr:  "test-serial-" + product.name,
					ProductStr: product.name,
				}
				callback(deviceInfo)
			}).
			Return(nil).Once()

		// Setup the open call for this device
		mockOpener.On("Open",
			uint16(VendorId),
			uint16(product.productId),
			"test-serial-"+product.name).
			Return(mockDevices[product.name], nil).Once()
	}

	// Return cleanup function
	cleanup := func() {
		defaultHIDEnumerator = originalHIDEnumerator
		defaultHIDOpener = originalHIDOpener
		defaultConfigUpdater = originalConfigUpdater
		lightBrightnessFunc = originalLightBrightnessFunc
		lightTemperatureFunc = originalLightTemperatureFunc
	}

	return mockDevice1, mockDevice2, mockEnumerator, mockOpener, mockConfigUpdater, cleanup
}

// Test LightOn function
func TestLightOn(t *testing.T) {
	mockDevice1, mockDevice2, _, _, mockConfigUpdater, cleanup := setupTest()
	defer cleanup()

	// Expected command bytes for turning on the light
	expectedBytes := []byte{0x11, 0xff, 0x04, 0x1c, LightOnCode, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	// Setup expectations - both devices get the write (deviceIndex 0 = all)
	mockDevice1.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice1.On("Close").Return(nil).Once()
	mockDevice2.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice2.On("Close").Return(nil).Once()
	mockConfigUpdater.On("UpdateCurrentState", 0, -1, -1, 1).Once()

	// Call the function
	LightOn(0)

	// Verify expectations
	mockDevice1.AssertExpectations(t)
	mockDevice2.AssertExpectations(t)
	mockConfigUpdater.AssertExpectations(t)
}

// Test LightOff function
func TestLightOff(t *testing.T) {
	mockDevice1, mockDevice2, _, _, mockConfigUpdater, cleanup := setupTest()
	defer cleanup()

	// Expected command bytes for turning off the light
	expectedBytes := []byte{0x11, 0xff, 0x04, 0x1c, LightOffCode, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	// Setup expectations
	mockDevice1.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice1.On("Close").Return(nil).Once()
	mockDevice2.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice2.On("Close").Return(nil).Once()
	mockConfigUpdater.On("UpdateCurrentState", 0, -1, -1, 0).Once()

	// Call the function
	LightOff(0)

	// Verify expectations
	mockDevice1.AssertExpectations(t)
	mockDevice2.AssertExpectations(t)
	mockConfigUpdater.AssertExpectations(t)
}

// Test LightBrightness function
func TestLightBrightness(t *testing.T) {
	mockDevice1, mockDevice2, _, _, mockConfigUpdater, cleanup := setupTest()
	defer cleanup()

	// Test with 50% brightness
	level := 50
	adjustedLevel := MinBrightness + ((MaxBrightness - MinBrightness) * level / 100)

	// Expected command bytes for setting brightness
	expectedBytes := []byte{0x11, 0xff, 0x04, 0x4c, 0x00, byte(adjustedLevel), 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	// Setup expectations
	mockDevice1.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice1.On("Close").Return(nil).Once()
	mockDevice2.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice2.On("Close").Return(nil).Once()
	mockConfigUpdater.On("UpdateCurrentState", 0, level, -1, -1).Once()

	// Call the function
	LightBrightness(0, level)

	// Verify expectations
	mockDevice1.AssertExpectations(t)
	mockDevice2.AssertExpectations(t)
	mockConfigUpdater.AssertExpectations(t)
}

// Test LightBrightDown function
func TestLightBrightDown(t *testing.T) {
	mockDevice1, mockDevice2, _, _, mockConfigUpdater, cleanup := setupTest()
	defer cleanup()

	// Current brightness is 50%
	currentBrightness := 50
	mockConfigUpdater.On("ReadCurrentState", 0).Return(currentBrightness, 4000, 1).Once()

	// Decrease by 10%
	decreaseAmount := 10
	newBrightness := currentBrightness - decreaseAmount
	adjustedLevel := MinBrightness + ((MaxBrightness - MinBrightness) * newBrightness / 100)

	// Expected command bytes for setting brightness
	expectedBytes := []byte{0x11, 0xff, 0x04, 0x4c, 0x00, byte(adjustedLevel), 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	// Setup expectations
	mockDevice1.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice1.On("Close").Return(nil).Once()
	mockDevice2.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice2.On("Close").Return(nil).Once()
	mockConfigUpdater.On("UpdateCurrentState", 0, newBrightness, -1, -1).Once()

	// Call the function
	LightBrightDown(0, decreaseAmount)

	// Verify expectations
	mockDevice1.AssertExpectations(t)
	mockDevice2.AssertExpectations(t)
	mockConfigUpdater.AssertExpectations(t)
}

// Test LightBrightDown function with minimum brightness
func TestLightBrightDownMinimum(t *testing.T) {
	mockDevice1, mockDevice2, _, _, mockConfigUpdater, cleanup := setupTest()
	defer cleanup()

	// Current brightness is 5%
	currentBrightness := 5
	mockConfigUpdater.On("ReadCurrentState", 0).Return(currentBrightness, 4000, 1).Once()

	// Decrease by 10% (should clamp to 0%)
	decreaseAmount := 10
	newBrightness := 0
	adjustedLevel := MinBrightness + ((MaxBrightness - MinBrightness) * newBrightness / 100)

	// Expected command bytes for setting brightness
	expectedBytes := []byte{0x11, 0xff, 0x04, 0x4c, 0x00, byte(adjustedLevel), 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	// Setup expectations
	mockDevice1.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice1.On("Close").Return(nil).Once()
	mockDevice2.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice2.On("Close").Return(nil).Once()
	mockConfigUpdater.On("UpdateCurrentState", 0, newBrightness, -1, -1).Once()

	// Call the function
	LightBrightDown(0, decreaseAmount)

	// Verify expectations
	mockDevice1.AssertExpectations(t)
	mockDevice2.AssertExpectations(t)
	mockConfigUpdater.AssertExpectations(t)
}

// Test LightBrightUp function
func TestLightBrightUp(t *testing.T) {
	mockDevice1, mockDevice2, _, _, mockConfigUpdater, cleanup := setupTest()
	defer cleanup()

	// Current brightness is 50%
	currentBrightness := 50
	mockConfigUpdater.On("ReadCurrentState", 0).Return(currentBrightness, 4000, 1).Once()

	// Increase by 10%
	increaseAmount := 10
	newBrightness := currentBrightness + increaseAmount
	adjustedLevel := MinBrightness + ((MaxBrightness - MinBrightness) * newBrightness / 100)

	// Expected command bytes for setting brightness
	expectedBytes := []byte{0x11, 0xff, 0x04, 0x4c, 0x00, byte(adjustedLevel), 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	// Setup expectations
	mockDevice1.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice1.On("Close").Return(nil).Once()
	mockDevice2.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice2.On("Close").Return(nil).Once()
	mockConfigUpdater.On("UpdateCurrentState", 0, newBrightness, -1, -1).Once()

	// Call the function
	LightBrightUp(0, increaseAmount)

	// Verify expectations
	mockDevice1.AssertExpectations(t)
	mockDevice2.AssertExpectations(t)
	mockConfigUpdater.AssertExpectations(t)
}

// Test LightBrightUp function with maximum brightness
func TestLightBrightUpMaximum(t *testing.T) {
	mockDevice1, mockDevice2, _, _, mockConfigUpdater, cleanup := setupTest()
	defer cleanup()

	// Current brightness is 95%
	currentBrightness := 95
	mockConfigUpdater.On("ReadCurrentState", 0).Return(currentBrightness, 2900, 1).Once()

	// Increase by 10% (should clamp to 100%)
	increaseAmount := 10
	newBrightness := 100

	// Expected command bytes for setting brightness
	expectedBytes := []byte{0x11, 0xff, 0x04, 0x4c, 0x00, byte(MaxBrightness), 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	// Setup expectations
	mockDevice1.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice1.On("Close").Return(nil).Once()
	mockDevice2.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice2.On("Close").Return(nil).Once()
	mockConfigUpdater.On("UpdateCurrentState", 0, newBrightness, -1, -1).Once()

	// Call the function
	LightBrightUp(0, increaseAmount)

	// Verify expectations
	mockDevice1.AssertExpectations(t)
	mockDevice2.AssertExpectations(t)
	mockConfigUpdater.AssertExpectations(t)
}

// Test LightTemperature function
func TestLightTemperature(t *testing.T) {
	mockDevice1, mockDevice2, _, _, mockConfigUpdater, cleanup := setupTest()
	defer cleanup()

	// Test with 4000K temperature
	temp := uint16(4000)
	tempBytes := []byte{0x0f, 0xa0} // 4000 in big-endian bytes

	// Expected command bytes for setting temperature
	expectedBytes := []byte{0x11, 0xff, 0x04, 0x9c, tempBytes[0], tempBytes[1], 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	// Setup expectations
	mockDevice1.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice1.On("Close").Return(nil).Once()
	mockDevice2.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice2.On("Close").Return(nil).Once()
	mockConfigUpdater.On("UpdateCurrentState", 0, -1, int(temp), -1).Once()

	// Call the function
	LightTemperature(0, temp)

	// Verify expectations
	mockDevice1.AssertExpectations(t)
	mockDevice2.AssertExpectations(t)
	mockConfigUpdater.AssertExpectations(t)
}

// Test LightTempDown function
func TestLightTempDown(t *testing.T) {
	mockDevice1, mockDevice2, _, _, mockConfigUpdater, cleanup := setupTest()
	defer cleanup()

	// Current temperature is 4000K
	currentTemp := 4000
	mockConfigUpdater.On("ReadCurrentState", 0).Return(50, currentTemp, 1).Once()

	// Decrease by 200K
	decreaseAmount := 200
	newTemp := currentTemp - decreaseAmount
	tempBytes := []byte{0x0e, 0xd8} // 3800 in big-endian bytes

	// Expected command bytes for setting temperature
	expectedBytes := []byte{0x11, 0xff, 0x04, 0x9c, tempBytes[0], tempBytes[1], 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	// Setup expectations
	mockDevice1.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice1.On("Close").Return(nil).Once()
	mockDevice2.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice2.On("Close").Return(nil).Once()
	mockConfigUpdater.On("UpdateCurrentState", 0, -1, newTemp, -1).Once()

	// Call the function
	LightTempDown(0, decreaseAmount)

	// Verify expectations
	mockDevice1.AssertExpectations(t)
	mockDevice2.AssertExpectations(t)
	mockConfigUpdater.AssertExpectations(t)
}

// Test LightTempDown function with minimum temperature
func TestLightTempDownMinimum(t *testing.T) {
	mockDevice1, mockDevice2, _, _, mockConfigUpdater, cleanup := setupTest()
	defer cleanup()

	// Current temperature is 2800K
	currentTemp := 2800
	mockConfigUpdater.On("ReadCurrentState", 0).Return(50, currentTemp, 1).Once()

	// Decrease by 200K (should clamp to 2700K)
	decreaseAmount := 200
	newTemp := 2700
	tempBytes := []byte{0x0a, 0x8c} // 2700 in big-endian bytes

	// Expected command bytes for setting temperature
	expectedBytes := []byte{0x11, 0xff, 0x04, 0x9c, tempBytes[0], tempBytes[1], 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	// Setup expectations
	mockDevice1.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice1.On("Close").Return(nil).Once()
	mockDevice2.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice2.On("Close").Return(nil).Once()
	mockConfigUpdater.On("UpdateCurrentState", 0, -1, newTemp, -1).Once()

	// Call the function
	LightTempDown(0, decreaseAmount)

	// Verify expectations
	mockDevice1.AssertExpectations(t)
	mockDevice2.AssertExpectations(t)
	mockConfigUpdater.AssertExpectations(t)
}

// Test LightTempUp function
func TestLightTempUp(t *testing.T) {
	mockDevice1, mockDevice2, _, _, mockConfigUpdater, cleanup := setupTest()
	defer cleanup()

	// Current temperature is 4000K
	currentTemp := 4000
	mockConfigUpdater.On("ReadCurrentState", 0).Return(50, currentTemp, 1).Once()

	// Increase by 200K
	increaseAmount := 200
	newTemp := currentTemp + increaseAmount
	tempBytes := []byte{0x10, 0x68} // 4200 in big-endian bytes

	// Expected command bytes for setting temperature
	expectedBytes := []byte{0x11, 0xff, 0x04, 0x9c, tempBytes[0], tempBytes[1], 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	// Setup expectations
	mockDevice1.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice1.On("Close").Return(nil).Once()
	mockDevice2.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice2.On("Close").Return(nil).Once()
	mockConfigUpdater.On("UpdateCurrentState", 0, -1, newTemp, -1).Once()

	// Call the function
	LightTempUp(0, increaseAmount)

	// Verify expectations
	mockDevice1.AssertExpectations(t)
	mockDevice2.AssertExpectations(t)
	mockConfigUpdater.AssertExpectations(t)
}

// Test LightTempUp function with maximum temperature
func TestLightTempUpMaximum(t *testing.T) {
	mockDevice1, mockDevice2, _, _, mockConfigUpdater, cleanup := setupTest()
	defer cleanup()

	// Current temperature is 6400K
	currentTemp := 6400
	mockConfigUpdater.On("ReadCurrentState", 0).Return(50, currentTemp, 1).Once()

	// Increase by 200K (should clamp to 6500K)
	increaseAmount := 200
	newTemp := 6500
	tempBytes := []byte{0x19, 0x64} // 6500 in big-endian bytes

	// Expected command bytes for setting temperature
	expectedBytes := []byte{0x11, 0xff, 0x04, 0x9c, tempBytes[0], tempBytes[1], 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	// Setup expectations
	mockDevice1.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice1.On("Close").Return(nil).Once()
	mockDevice2.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice2.On("Close").Return(nil).Once()
	mockConfigUpdater.On("UpdateCurrentState", 0, -1, newTemp, -1).Once()

	// Call the function
	LightTempUp(0, increaseAmount)

	// Verify expectations
	mockDevice1.AssertExpectations(t)
	mockDevice2.AssertExpectations(t)
	mockConfigUpdater.AssertExpectations(t)
}

// Test LightOn with a single device target
func TestLightOnSingleDevice(t *testing.T) {
	mockDevice1, mockDevice2, _, _, mockConfigUpdater, cleanup := setupTest()
	defer cleanup()

	expectedBytes := []byte{0x11, 0xff, 0x04, 0x1c, LightOnCode, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	// Only device 1 (Beam) should receive the write
	mockDevice1.On("Write", expectedBytes).Return(len(expectedBytes), nil).Once()
	mockDevice1.On("Close").Return(nil).Once()
	// Device 2 (Glow) should only be closed, not written to
	mockDevice2.On("Close").Return(nil).Once()
	mockConfigUpdater.On("UpdateCurrentState", 1, -1, -1, 1).Once()

	LightOn(1)

	mockDevice1.AssertExpectations(t)
	mockDevice2.AssertExpectations(t)
	mockDevice2.AssertNotCalled(t, "Write", mock.Anything)
	mockConfigUpdater.AssertExpectations(t)
}

// Test ListDevices function
func TestListDevices(t *testing.T) {
	mockDevice1, mockDevice2, _, _, _, cleanup := setupTest()
	defer cleanup()

	mockDevice1.On("Close").Return(nil).Once()
	mockDevice2.On("Close").Return(nil).Once()

	devices := ListDevices()

	if len(devices) != 2 {
		t.Fatalf("Expected 2 devices, got %d", len(devices))
	}

	// Sorted by serial: Beam first, Glow second
	if devices[0].Name != "Beam" {
		t.Errorf("Expected first device to be Beam, got %s", devices[0].Name)
	}
	if devices[0].Index != 1 {
		t.Errorf("Expected first device index to be 1, got %d", devices[0].Index)
	}
	if devices[0].Serial != "test-serial-Beam" {
		t.Errorf("Expected first device serial to be test-serial-Beam, got %s", devices[0].Serial)
	}
	if devices[1].Name != "Glow" {
		t.Errorf("Expected second device to be Glow, got %s", devices[1].Name)
	}
	if devices[1].Index != 2 {
		t.Errorf("Expected second device index to be 2, got %d", devices[1].Index)
	}
	if devices[1].Serial != "test-serial-Glow" {
		t.Errorf("Expected second device serial to be test-serial-Glow, got %s", devices[1].Serial)
	}
}
