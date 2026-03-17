//go:build !windows

package lib

// preparePayload returns the data unchanged on Linux and macOS.
// On Linux, hid_linux.c does a raw write() to hidraw, sending bytes as-is to the device.
// Prepending a 0x00 Report ID would be sent literally to the device (wrong).
// On macOS, IOKit handles the single-report case correctly with the current byte sequence.
func preparePayload(data []byte) []byte {
	return data
}
