FROM debian:bookworm-slim as isolate-builder

RUN apt-get update && apt-get install --no-install-recommends -y \
    build-essential pkg-config libsystemd-dev libcap-dev gcc

WORKDIR /app
COPY third_party/ ./third_party
WORKDIR /app/third_party/isolate
RUN make isolate isolate-check-environment isolate-cg-keeper

FROM golang:1.22.4-bookworm

RUN apt-get update && apt-get install --no-install-recommends -y \
    build-essential bash gcc g++ python3 pypy3

WORKDIR /app/third_party/isolate
COPY --from=isolate-builder /app/third_party /app/third_party
RUN make install

WORKDIR /app
COPY pkg/ ./pkg
COPY go.mod go.sum ./

COPY compose/judge/entrypoint /entrypoint
RUN sed -i 's/\r$//g' /entrypoint
RUN chmod +x /entrypoint

RUN go mod download && go mod verify

ENTRYPOINT ["/entrypoint"]
