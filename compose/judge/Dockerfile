FROM ghcr.io/thepluck/judge-essential:latest

WORKDIR /app

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

RUN go generate ./...

COPY compose/judge/entrypoint /entrypoint
RUN sed -i 's/\r$//g' /entrypoint
RUN chmod +x /entrypoint

COPY compose/judge/start /start
RUN sed -i 's/\r$//g' /start
RUN chmod +x /start

ENTRYPOINT ["/entrypoint"]
