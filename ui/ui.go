package main

import (
	"fmt"

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

func main() {
	application := app.NewWithID("net.khary.lcui")
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
				lib.LightOff()
			}),
			fyne.NewMenuItem("On", func() {
				lib.LightOn()
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
			lib.LightOff()
		} else {
			lib.LightOn()
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
			bright, temp := config.ReadProfile(selection)
			brightnessSlider.SetValue(float64(bright))
			brightnessLabel.SetText(fmt.Sprintf("Brightness %d%%", int(bright)))
			tempSlider.SetValue(float64(temp))
			tempLabel.SetText(fmt.Sprintf("Temperature %dk", uint16(temp)))
			config.UpdateCurrentState(bright, temp)
			lib.LightBrightness(bright)
			lib.LightTemperature(uint16(temp))
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
			config.AddOrUpdateProfile(profileName, int(brightnessSlider.Value), int(tempSlider.Value))
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
		lib.LightBrightness(int(brightness))
		brightnessLabel.SetText(fmt.Sprintf("Brightness %d%%", int(brightness)))
		config.AddOrUpdateProfile(profileSelector.Selected, int(brightness), -1)
	}
	tempSlider.OnChanged = func(temp float64) {
		lib.LightTemperature(uint16(temp))
		tempLabel.SetText(fmt.Sprintf("Temperature %dk", uint16(temp)))
		config.AddOrUpdateProfile(profileSelector.Selected, -1, int(temp))
	}

	// Set Current Values
	currentBright, currentTemp := config.ReadCurrentState()
	brightnessSlider.SetValue(float64(currentBright))
	tempSlider.SetValue(float64(currentTemp))
	brightnessLabel.SetText(fmt.Sprintf("Brightness %d%%", int(currentBright)))
	tempLabel.SetText(fmt.Sprintf("Temperature %dk", uint16(currentTemp)))

	// Add all widgets to the container
	mainGroup := container.New(layout.NewVBoxLayout(), powerGroup, profileGroup, brightnessGroup, tempGroup, exitButton)

	mainWindow.SetContent(mainGroup)

	mainWindow.ShowAndRun()
}
