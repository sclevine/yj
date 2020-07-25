FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /workspace
COPY . .
ARG version="0.0.0"
RUN go build -ldflags "-X main.Version=$version"
FROM scratch
COPY --from=builder /workspace/yj /bin/yj
ENTRYPOINT ["/bin/yj"]
