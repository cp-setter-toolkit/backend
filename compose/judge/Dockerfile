FROM golang:1.22.4-bookworm

WORKDIR /app

RUN apt-get update && apt-get install --no-install-recommends -y \
    build-essential bash gcc g++ libcap-dev libsystemd-dev

COPY third_party/ ./third_party
WORKDIR /app/third_party/isolate
RUN make install

WORKDIR /app
COPY pkg/ ./pkg
COPY go.mod go.sum ./

RUN go mod download && go mod verify
