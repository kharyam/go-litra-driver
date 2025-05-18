// Package config is responsible for parsing the user config file for the light(s).  THe config file
// persists the current state of the light along with any presets.
package config

import (
	"errors"
	"os"

	"path/filepath"

	"strconv"

	"github.com/bigkevmcd/go-configparser"
	"github.com/rs/zerolog/log"
)

const CurrentProfileName = "current"
const Bright = "brightness"
const Temp = "temperature"
const Power = "power"

// getConfig loads the config file
func getConfig() (*configparser.ConfigParser, string) {
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	configFile := ""
	if xdgConfig != "" {
		configFile = filepath.Join(xdgConfig, "llgd")
		pathExists, _ := exists(configFile)
		if !pathExists {
			os.MkdirAll(configFile, os.ModeDir)
		}
		configFile = filepath.Join(configFile, "config")
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Msgf("Failed to find home directory: %v", err)
		}
		configFile = filepath.Join(homeDir, ".llgd_config")
	}

	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		cfile, e := os.Create(configFile)
		if e != nil {
			log.Fatal().Msgf("Failed to create new config file %s %v", configFile, e)
		}
		cfile.Close()
	}

	configParser, err := configparser.NewConfigParserFromFile(configFile)
	if err != nil {
		log.Fatal().Msgf("Failed to load config file %s : %v", configFile, err)
	}

	return configParser, configFile
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// AddOrUpdateProfile will create a new profile or update an existing profile
func AddOrUpdateProfile(profileName string, brightness int, temp int, power int) {
	parser, configFile := getConfig()
	if !parser.HasSection(profileName) {
		parser.AddSection(profileName)
	}
	if brightness != -1 {
		parser.Set(profileName, Bright, strconv.Itoa(brightness))
	}
	if temp != -1 {
		parser.Set(profileName, Temp, strconv.Itoa(temp))
	}
	if power != -1 {
		parser.Set(profileName, Power, strconv.Itoa(power))
	}
	parser.SaveWithDelimiter(configFile, "=")
}

// UpdateCurrentState updates the temperature, brightness, and/or power for current state.
// set any value to -1 to not set it in the section
func UpdateCurrentState(brightness int, temperature int, power int) {
	AddOrUpdateProfile(CurrentProfileName, brightness, temperature, power)
}

// DeleteProfile removes a profile from the configuration file
func DeleteProfile(profileName string) {
	parser, configFile := getConfig()
	if parser.HasSection(profileName) {
		parser.RemoveSection(profileName)
		parser.SaveWithDelimiter(configFile, "=")
	}
}

// ReadProfile will read the brightness, temperature, and power settings from a profile
func ReadProfile(profileName string) (brightness int, temperature int, power int) {
	parser, _ := getConfig()

	brightnessString, err := parser.Get(profileName, Bright)
	if err != nil {
		brightness = -1
	} else {
		brightness, _ = strconv.Atoi(brightnessString)
	}

	temperatureString, err := parser.Get(profileName, Temp)
	if err != nil {
		temperature = -1
	} else {
		temperature, _ = strconv.Atoi(temperatureString)
	}

	powerString, err := parser.Get(profileName, Power)
	if err != nil {
		power = -1
	} else {
		power, _ = strconv.Atoi(powerString)
	}

	return brightness, temperature, power

}

// Read the current state of the lights from the config file
func ReadCurrentState() (brightness int, temperature int, power int) {
	return ReadProfile(CurrentProfileName)
}

// Return the list of profile names with "current" being first
func GetProfileNames() (profiles []string) {
	parser, _ := getConfig()
	allProfiles := parser.Sections()

	profiles = append(profiles, CurrentProfileName)

	for i := 0; i < len(allProfiles); i++ {
		if allProfiles[i] != CurrentProfileName {
			profiles = append(profiles, allProfiles[i])
		}
	}

	return profiles

}
