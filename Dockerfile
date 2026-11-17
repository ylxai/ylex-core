FROM golang:1.22-alpine AS builder
RUN apk add --no-cache make gcc g++ musl-dev linux-headers git
WORKDIR /src
COPY . .
RUN make ylexd

FROM alpine:latest
RUN apk add --no-cache libstdc++ ca-certificates
COPY --from=builder /src/build/bin/ylexd /usr/local/bin/ylexd
EXPOSE 30300 30300/udp 8545
ENTRYPOINT ["/usr/local/bin/ylexd"]
