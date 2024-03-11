# Develpment Notes

## Build Locally

```bash

git clone git@github.com:kharyam/go-litra-driver.git
cd go-litra-driver
go build -v ./config
go build -v ./lib
go build -o lcli -v ./lcli
go build -o lcui -v ./lcui
```

## Cobra Config (for future reference)

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

## Publishing

```bash
export VERSION=v0.0.2

cd config
go mod tidy
cd ../lib
go mod tidy
cd ../lcli
go mod tidy
cd ../lcui
go mod tidy

GOPROXY=proxy.golang.org go list -m github.com/kharyam/go-litra-driver@${VERSION}
```

## Packaging
```bash
podman build -t kharyam/fyne-cross-images:linux build/linux

cd lcli
fyne-cross linux --arch=amd64 --image=kharyam/fyne-cross-images:linux --app-id=net.kharyam.lcli
fyne-cross windows --arch=amd64 --app-id=net.kharyam.lcli
# TODO - Package for OSX
#fyne-cross darwin --arch=amd64 --app-id=net.kharyam.lcli

cd ../lcui
fyne-cross linux --arch=amd64 --image=kharyam/fyne-cross-images:linux --app-id=net.kharyam.lcui
fyne-cross windows --arch=amd64 --app-id=net.kharyam.lcui
# TODO - Package for OSX
#fyne-cross darwin --arch=amd64 --app-id=net.kharyam.lcui
```
