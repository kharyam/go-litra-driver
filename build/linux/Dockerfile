FROM docker.io/fyneio/fyne-cross-images:linux

ARG GO_VERSION=1.25.7
RUN rm -rf /usr/local/go && \
    curl -fsSL "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz" | \
    tar -C /usr/local -xzf - && \
    ln -sf /usr/local/go/bin/go /usr/local/bin/go && \
    ln -sf /usr/local/go/bin/gofmt /usr/local/bin/gofmt

RUN dpkg --add-architecture arm64 && \
    apt-get update && \
    apt-get install -y libudev-dev udev && \
    apt-get install -y libudev-dev:arm64 udev:arm64
