package main

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kharyam/go-litra-driver/config"
	"github.com/kharyam/go-litra-driver/lib"
)

var selectedDeviceIndex int = 0

//go:generate fyne bundle -o icons.go Icon.png
func main() {
	application := app.NewWithID("net.khary.lcui")
	application.SetIcon(resourceIconPng)

	mainWindow := application.NewWindow("Litra Controller")

	if desk, ok := application.(desktop.App); ok {
		systrayMenu := fyne.NewMenu("LitraController",
			fyne.NewMenuItem("Show", func() {
				mainWindow.Show()
			}),
			fyne.NewMenuItem("Hide", func() {
				mainWindow.Hide()
			}),
			fyne.NewMenuItem("Off", func() {
				lib.LightOff(0)
			}),
			fyne.NewMenuItem("On", func() {
				lib.LightOn(0)
			}),
		)
		desk.SetSystemTrayMenu(systrayMenu)
	}

	mainWindow.SetCloseIntercept(func() {
		mainWindow.Hide()
	})

	// Power
	powerRadio := widget.NewRadioGroup([]string{"Off", "On"}, func(power string) {
		if power == "Off" {
			lib.LightOff(selectedDeviceIndex)
		} else {
			lib.LightOn(selectedDeviceIndex)
		}
	})
	powerRadio.Horizontal = true
	powerLabel := widget.NewLabel("Power")
	powerGroup := container.New(layout.NewHBoxLayout(), powerLabel, powerRadio)

	// Brightness
	brightnessLabel := widget.NewLabel("Brightness")
	brightnessSlider := widget.NewSlider(1, 100)
	brightnessSlider.Step = 1
	brightnessGroup := container.New(layout.NewVBoxLayout(), brightnessLabel, brightnessSlider)

	// Temperature
	tempLabel := widget.NewLabel("Temperature")
	tempSlider := widget.NewSlider(2700, 6500)
	tempSlider.Step = 100
	tempGroup := container.New(layout.NewVBoxLayout(), tempLabel, tempSlider)

	// Device Selector
	devices := lib.ListDevices()
	deviceOptions := []string{"All Devices"}
	for _, d := range devices {
		deviceOptions = append(deviceOptions, fmt.Sprintf("Device %d: Litra %s", d.Index, d.Name))
	}
	deviceLabel := widget.NewLabel("Device:")
	deviceSelector := widget.NewSelect(deviceOptions, func(selection string) {
		if selection == "All Devices" {
			selectedDeviceIndex = 0
		} else {
			// Parse index from "Device N: ..."
			parts := strings.SplitN(selection, ":", 2)
			if len(parts) > 0 {
				numStr := strings.TrimPrefix(strings.TrimSpace(parts[0]), "Device ")
				if idx, err := strconv.Atoi(numStr); err == nil {
					selectedDeviceIndex = idx
				}
			}
		}
		// Refresh UI from selected device's state
		bright, temp, power := config.ReadCurrentState(selectedDeviceIndex)
		brightnessSlider.SetValue(float64(bright))
		brightnessLabel.SetText(fmt.Sprintf("Brightness %d%%", int(bright)))
		tempSlider.SetValue(float64(temp))
		tempLabel.SetText(fmt.Sprintf("Temperature %dk", uint16(temp)))
		if power == 1 {
			powerRadio.SetSelected("On")
		} else {
			powerRadio.SetSelected("Off")
		}
	})
	deviceSelector.SetSelected("All Devices")
	deviceGroup := container.New(layout.NewHBoxLayout(), deviceLabel, deviceSelector)

	// Profiles
	profileNew := widget.NewButton("New...", func() {
		fmt.Printf("Save As Clicked")
	})
	profileDelete := widget.NewButton("Delete", func() {
		fmt.Printf("Delete Clicked")
	})
	profileDelete.Disable()
	profileNew.Enable()
	profileLabel := widget.NewLabel("Preset:")
	profileSelector := widget.NewSelect(config.GetProfileNames(), func(selection string) {
		if selection == config.CurrentProfileName {
			profileNew.Enable()
			profileDelete.Disable()
		} else {
			profileNew.Disable()
			profileDelete.Enable()
			bright, temp, power := config.ReadProfile(selection)
			brightnessSlider.SetValue(float64(bright))
			brightnessLabel.SetText(fmt.Sprintf("Brightness %d%%", int(bright)))
			tempSlider.SetValue(float64(temp))
			tempLabel.SetText(fmt.Sprintf("Temperature %dk", uint16(temp)))
			config.UpdateCurrentState(selectedDeviceIndex, bright, temp, power)
			lib.LightBrightness(selectedDeviceIndex, bright)
			lib.LightTemperature(selectedDeviceIndex, uint16(temp))
		}
	})
	profileDelete.OnTapped = func() {
		dialog.ShowConfirm("Delete Profile?", fmt.Sprintf("Delete Profile \"%s\"?", profileSelector.Selected), func(delete bool) {
			if delete {
				config.DeleteProfile(profileSelector.Selected)
				profileSelector.SetOptions(config.GetProfileNames())
				profileSelector.SetSelected(config.CurrentProfileName)
			}
		}, mainWindow)
	}

	profileNew.OnTapped = func() {
		dialog.ShowEntryDialog("New Profile", "Name", func(profileName string) {
			_, _, currentPower := config.ReadCurrentState(selectedDeviceIndex)
			config.AddOrUpdateProfile(profileName, int(brightnessSlider.Value), int(tempSlider.Value), currentPower)
			profileSelector.SetOptions(config.GetProfileNames())
			profileSelector.SetSelected(profileName)
		}, mainWindow)
	}
	profileSelector.SetSelected(config.CurrentProfileName)
	profileGroup := container.New(layout.NewHBoxLayout(), profileLabel, profileSelector, profileNew, profileDelete)

	// Exit
	exitButton := widget.NewButton("Exit", func() {
		application.Quit()
	})

	// Callbacks
	brightnessSlider.OnChanged = func(brightness float64) {
		brightnessLabel.SetText(fmt.Sprintf("Brightness %d%%", int(brightness)))
	}
	tempSlider.OnChanged = func(temp float64) {
		tempLabel.SetText(fmt.Sprintf("Temperature %dk", uint16(temp)))
	}

	brightnessSlider.OnChangeEnded = func(brightness float64) {
		lib.LightBrightness(selectedDeviceIndex, int(brightness))
		brightnessLabel.SetText(fmt.Sprintf("Brightness %d%%", int(brightness)))
		_, _, currentPower := config.ReadCurrentState(selectedDeviceIndex)
		config.AddOrUpdateProfile(profileSelector.Selected, int(brightness), -1, currentPower)
	}

	tempSlider.OnChangeEnded = func(temp float64) {
		lib.LightTemperature(selectedDeviceIndex, uint16(temp))
		tempLabel.SetText(fmt.Sprintf("Temperature %dk", uint16(temp)))
		_, _, currentPower := config.ReadCurrentState(selectedDeviceIndex)
		config.AddOrUpdateProfile(profileSelector.Selected, -1, int(temp), currentPower)
	}

	// Set Current Values
	currentBright, currentTemp, currentPower := config.ReadCurrentState(selectedDeviceIndex)
	brightnessSlider.SetValue(float64(currentBright))
	tempSlider.SetValue(float64(currentTemp))
	brightnessLabel.SetText(fmt.Sprintf("Brightness %d%%", int(currentBright)))
	tempLabel.SetText(fmt.Sprintf("Temperature %dk", uint16(currentTemp)))
	if currentPower == 1 {
		powerRadio.SetSelected("On")
	} else {
		powerRadio.SetSelected("Off")
	}

	// Add all widgets to the container
	mainGroup := container.New(layout.NewVBoxLayout(), deviceGroup, powerGroup, profileGroup, brightnessGroup, tempGroup, exitButton)

	mainWindow.SetContent(mainGroup)

	mainWindow.ShowAndRun()
}
