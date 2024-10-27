# Build the manager binary
FROM golang:1.23 AS builder

WORKDIR /workspace

COPY go.mod go.mod

COPY go.sum go.sum

RUN go mod download

COPY cmd/main.go cmd/main.go

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/main.go


ENTRYPOINT ["./workspace/manager"]
