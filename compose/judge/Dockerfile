FROM ghcr.io/minhnhatnoe/isolate:latest

RUN apt-get update && apt-get install --no-install-recommends -y \
    gcc g++ python2.7 python3.10 pypy pypy3 \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY pkg/ ./pkg
COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY compose/judge/entrypoint /entrypoint
RUN sed -i 's/\r$//g' /entrypoint
RUN chmod +x /entrypoint

COPY compose/judge/start /start
RUN sed -i 's/\r$//g' /start
RUN chmod +x /start

ENTRYPOINT ["/entrypoint"]
