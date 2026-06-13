FROM golang:1.26-alpine as builder

RUN apk add --no-cache git make openssh

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./internal ./internal
COPY ./cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -extldflags \"-static\" -X main.version=${VERSION}" \
    -o /app/appbin ./cmd/main.go

FROM alpine:latest
LABEL MAINTAINER Author mohamed abdelmohaimen

RUN adduser --home "/appuser" --disabled-password appuser \
    --gecos "appuser,-,-,-"
USER appuser

COPY --from=builder /app/appbin /home/appuser/app/appbin

WORKDIR /home/appuser/app

# Export necessary port
EXPOSE 9090
# Command to run when starting the container
CMD ["./appbin"]