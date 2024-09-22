FROM golang:alpine AS builder
LABEL stage=gobuilder

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
COPY etc /rpc/etc
# Build rpc with -s (inlining) -w (dead code elimination)
RUN go build -ldflags="-s -w" -o /rpc/main .


FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata

WORKDIR /rpc
COPY --from=builder /rpc /rpc
COPY --from=builder /rpc/etc /rpc/etc

EXPOSE 9000

CMD ["./main", "-f", "etc/rdportfoliorpc.yaml"]
