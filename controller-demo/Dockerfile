FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /workspace

COPY . .

RUN go mod download && go build -o controller ./cmd/

FROM alpine

COPY --from=builder /workspace/controller /

ENTRYPOINT ["/controller"]