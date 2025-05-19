FROM docker.io/fyneio/fyne-cross-images:linux

RUN dpkg --add-architecture arm64 && \
    apt-get update && \
    apt-get install -y libudev-dev udev && \
    apt-get install -y libudev-dev:arm64 udev:arm64
