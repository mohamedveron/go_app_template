FROM golang:1.20-alpine as builder

RUN apk add --no-cache git
RUN apk add --update make
RUN apk add --no-cache openssh

# Move to working directory /app
WORKDIR /app

# Copy the code into the container
COPY . .

# Toggle CGO on your app requirement
RUN CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o /app/appbin ./cmd/main.go

FROM alpine:latest
LABEL MAINTAINER Author mohamed abdelmohaimen
# Add new user 'appuser'. App should be run without root privileges as a security measure
RUN adduser --home "/appuser" --disabled-password appuser \
    --gecos "appuser,-,-,-"
USER appuser

COPY --from=builder /app/internal/server/http/web /home/appuser/app/web
COPY --from=builder /app/appbin /home/appuser/app

WORKDIR /home/appuser/app

# Export necessary port
EXPOSE 9090
# Command to run when starting the container
CMD ["./appbin"]