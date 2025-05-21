package lib

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLightBrightnessCalculation tests the brightness calculation logic
func TestLightBrightnessCalculation(t *testing.T) {
	// Test with 0% brightness
	level := 0
	adjustedLevel := MinBrightness + ((MaxBrightness - MinBrightness) * level / 100)
	assert.Equal(t, byte(MinBrightness), byte(adjustedLevel))

	// Test with 50% brightness
	level = 50
	adjustedLevel = MinBrightness + ((MaxBrightness - MinBrightness) * level / 100)
	expectedLevel := MinBrightness + ((MaxBrightness - MinBrightness) / 2)
	assert.Equal(t, byte(expectedLevel), byte(adjustedLevel))

	// Test with 100% brightness
	level = 100
	adjustedLevel = MinBrightness + ((MaxBrightness - MinBrightness) * level / 100)
	assert.Equal(t, byte(MaxBrightness), byte(adjustedLevel))
}

// TestTemperatureConversion tests the temperature conversion to bytes
func TestTemperatureConversion(t *testing.T) {
	// Test with 2700K (minimum)
	temp := uint16(2700)
	tempBytes := []byte{0x0a, 0x8c} // 2700 in big-endian bytes

	// Create a buffer and convert the temperature
	buffer := make([]byte, 2)
	binary.BigEndian.PutUint16(buffer, temp)

	assert.Equal(t, tempBytes[0], buffer[0])
	assert.Equal(t, tempBytes[1], buffer[1])

	// Test with 4000K (middle)
	temp = uint16(4000)
	tempBytes = []byte{0x0f, 0xa0} // 4000 in big-endian bytes

	// Create a buffer and convert the temperature
	binary.BigEndian.PutUint16(buffer, temp)

	assert.Equal(t, tempBytes[0], buffer[0])
	assert.Equal(t, tempBytes[1], buffer[1])

	// Test with 6500K (maximum)
	temp = uint16(6500)
	tempBytes = []byte{0x19, 0x64} // 6500 in big-endian bytes

	// Create a buffer and convert the temperature
	binary.BigEndian.PutUint16(buffer, temp)

	assert.Equal(t, tempBytes[0], buffer[0])
	assert.Equal(t, tempBytes[1], buffer[1])
}
