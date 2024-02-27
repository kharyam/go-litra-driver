![Go workflow](https://github.com/kharyam/go-litra-driver/actions/workflows/go.yml/badge.svg)
# Go utility for Litra Glow and Beam

## Introduction

After purchasing a [Logitech Litra Glow](https://www.logitech.com/en-us/products/lighting/litra-glow.946-000001.html) I was unable to find any support for linux. [I created a python project](https://github.com/kharyam/litra-driver) which reverse-engineered the basic functionality of the litra glow so that we can control it via USB without using the physical buttons on the device. It also now supports the [Logitech Litra Beam](https://www.logitech.com/en-us/products/lighting/litra-beam.946-000006.html). This project is the port from python to go.


## Quick Start

### Linux
```bash
# If necessary, create a udev role to grant permission to access the light
sudo tee /etc/udev/rules.d/82-litra-glow.rules <<< 'SUBSYSTEM=="usb", ATTR{idVendor}=="046d", ATTR{idProduct}=="c900",MODE="0666"'
sudo tee /etc/udev/rules.d/82-litra-beam.rules <<< 'SUBSYSTEM=="usb", ATTR{idVendor}=="046d", ATTR{idProduct}=="c901",MODE="0666"'

# For most operating systems, reloading udev rules is enough
sudo udevadm control --reload-rules
sudo udevadm trigger

# For others, a reboot may be required
# sudo reboot

# Command Line Interface
lcli --help

# User Interface
lcui

```

## The UI
A basic UI can be launched to allow control of the light:

```bash
lcui
```

![lcui Screen Shot](images/lcui.png)

## The CLI

```bash
This command line interface allows you to control a litra Glow or Beam 
device using the commands described below.

Usage:
  cli [command]

Available Commands:
  bright      Sets the brightness level (0-100)
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  off         Turn lights off
  on          Turn lights on
  temp        Sets the temperature of the lights (2700 - 6500)

Flags:
  -h, --help     help for cli
  -t, --toggle   Help message for toggle

Use "cli [command] --help" for more information about a command.
```

Sample Usage
```bash
lcli on
lcli bright 10
lcli temp 6500
lcli off
```

## Development

### Linux Pre-requisites
```bash
sudo apt-get install -y libusb-dev libusb-1.0-0-dev
```

### OS X Pre-requisities
```bash
brew install libusb
```

### Windows Pre-requisities (Work in Progress!!)

Follow the steps under `Installing the MinGW-w64 toolchain` [here](https://code.visualstudio.com/docs/cpp/config-mingw)
**Do not follow the step for pkg-config - the link currently goes to a malicious website. pkg-config will be installed as part of the Gcc installation**

Be sure to run the terminal `MSYS32 MINGW64`

```bash
pacman -Sy mingw-w64-x86_64-toolchain
pacman -Sy mingw-w64-x86_64-libusb

```
create /mingw64/lib/pkgconfig/libusb-1.0.pc

```
#lib-usb-1.0.pc
prefix=/mingw64
exec_prefix=${prefix}
libdir=${exec_prefix}/lib
includedir=${prefix}/include

Name: libusb-1.0
Description: C API for USB device access from Linux, Mac OS X, Windows, OpenBSD/NetBSD and Solaris userspace

#update version to match
Version: 1.0.27
Libs: -L${libdir} -llibusb-1.0
Libs.private: 
Cflags: -I${includedir}/libusb-1.0
```



Copy `libusb-1.0.dll` to `c:\Windows\System32`

Use [zadig](https://zadig.akeo.ie/) to replace the USB drivers for the litra glow to use WinUSB...

![Zadig Screen Shot](images/zadig.png)

References
* https://github.com/google/gousb?tab=readme-ov-file#notes-for-installation-on-windows
* https://code.visualstudio.com/docs/cpp/config-mingw



### Build Locally


```bash

# Linux
sudo apt-get install -y libusb-dev libusb-1.0-0-dev

git clone git@github.com:kharyam/go-litra-driver.git
cd go-litra-driver
go build -v ./config
go build -v ./lib
go build -o lcli -v ./cli
go build -o lcui -v ./ui
```

### Cobra Config (for future reference)

```bash

cd cli

# Workaround when using workspaces
GOWORK=off cobra-cli init .

# Create skeleton code for each command
cobra-cli add on
cobra-cli add off
cobra-cli add bright
cobra-cli add temp
```
