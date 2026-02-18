# Develpment Notes

## Build Locally

```bash

git clone git@github.com:kharyam/go-litra-driver.git
cd go-litra-driver

go build -v ./config
go build -v ./lib
go build -o lcli -v ./lcli
go build -o lcui -v ./lcui
go generate  -v ./lcui  # to update the icon from the PNG file
```

## Run unit tests

```bash
go test -cover -coverprofile=coverage.out -v ./config ./lib ./lcli/cmd

# View coverage in browser
go tool cover -html=coverage.out
```

## Publishing

```bash
# Update to new version
export VERSION=v0.1.8

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

# Commit and push changes to feature branch

# Push to main branch (e.g., merge pull request for the branch)

# Switch to the main branch
git checkout main
git pull origin main

# Tag as the version defined above
git tag -a ${VERSION} -m "Release version ${VERSION:1}"
git push origin ${VERSION}

# Git hub action will build for all supported platforms, run unit tests, and create the Release

# Go to the release page once the action completes and update the release description (Auto generate)
GOPROXY=proxy.golang.org go list -m github.com/kharyam/go-litra-driver@${VERSION}
```

## Packaging

This is for reference - there is a [GitHub Action](.github/workflows/cross-compile.yml) to build all supported versions.

```bash
podman build -t kharyam/fyne-cross-images:linux build/linux

cd lcli
fyne-cross linux --arch=amd64 --image=kharyam/fyne-cross-images:linux --app-id=net.kharyam.lcli --name lcli-amd64
fyne-cross linux --arch=arm64 --image=kharyam/fyne-cross-images:linux --app-id=net.kharyam.lcli --name lcli-arm64
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
