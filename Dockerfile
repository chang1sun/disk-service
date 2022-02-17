FROM golang:1.17-alpine as build

# Install git.
# Git is required for fetching the dependencies.
RUN apk add --no-cache git bash curl ca-certificates

WORKDIR /go/src

COPY . .

# $GOPATH/bin添加到环境变量中
ENV PATH $GOPATH/bin:$PATH
RUN export GO111MODULE=on
ENV GOPROXY="https://goproxy.cn,direct"

# Fetch dependencies
RUN go mod download

# Build the binary
RUN go build -o /server

EXPOSE 8001
# Run the binary.
CMD [ "/server" ]