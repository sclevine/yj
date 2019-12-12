FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /workspace
COPY . .
RUN go build
FROM scratch
COPY --from=builder /workspace/yj /bin/yj
ENTRYPOINT ["/bin/yj"]
