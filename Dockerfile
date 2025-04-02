# build
FROM golang:1.23.6 AS builder
WORKDIR /app
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x
COPY . /app
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target="/root/.cache/go-build" \
    CGO_ENABLED=0 \
    go build -o todolist-user
# run
FROM alpine:3.19.2
WORKDIR /app
COPY --from=builder /app/todolist-user /app
COPY ./config.yaml /app
ENTRYPOINT ["./todolist-user"]