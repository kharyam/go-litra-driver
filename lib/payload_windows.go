//go:build windows

package lib

// preparePayload prepends a 0x00 Report ID byte required by the Windows HID API.
// On Windows, hid_write() treats the first byte as the Report ID and caps.OutputReportByteLength
// always includes 1 byte for the Report ID (even for single-report devices with ID 0x00).
// Without this prefix, Windows incorrectly uses 0x11 as the Report ID and the device
// receives the wrong data, silently ignoring the command.
func preparePayload(data []byte) []byte {
	result := make([]byte, len(data)+1)
	result[0] = 0x00
	copy(result[1:], data)
	return result
}
