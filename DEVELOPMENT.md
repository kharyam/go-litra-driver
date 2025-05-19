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

## Publishing

```bash
export VERSION=v0.1.3

cd config
go get -u
cd ../lib
go get -u
cd ../lcli
go get -u
cd ../lcui
go get -u
cd ..

cd config
go mod tidy
cd ../lib
go mod tidy
cd ../lcli
go mod tidy
cd ../lcui
go mod tidy
cd ..

# Push to main branch, then tag as the version defined above

GOPROXY=proxy.golang.org go list -m github.com/kharyam/go-litra-driver@${VERSION}
```

## Packaging
```bash
podman build -t kharyam/fyne-cross-images:linux build/linux

cd lcli
fyne-cross linux --arch=amd64 --image=kharyam/fyne-cross-images:linux --app-id=net.kharyam.lcli --name lcli-amd64
fyne-cross linux --arch=arm64 --image=kharyam/fyne-cross-images:linuxcd .. --app-id=net.kharyam.lcli --name lcli-arm64
fyne-cross windows --arch=amd64 --app-id=net.kharyam.lcli
# TODO - Package for OSX
#fyne-cross darwin --arch=amd64 --app-id=net.kharyam.lcli

cd ../lcui
fyne-cross linux --arch=amd64 --image=kharyam/fyne-cross-images:linux --app-id=net.kharyam.lcui --name lcui-amd64
fyne-cross linux --arch=arm64 --image=kharyam/fyne-cross-images:linux --app-id=net.kharyam.lcui --name lcui-arm64
fyne-cross windows --arch=amd64 --app-id=net.kharyam.lcui
# TODO - Package for OSX
#fyne-cross darwin --arch=amd64 --app-id=net.kharyam.lcui

cd ..
find . -name "*.xz" -exec mv {} build \;
find . -name "*.zip" -exec mv {} build \;

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