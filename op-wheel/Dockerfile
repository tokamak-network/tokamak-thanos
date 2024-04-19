ARG OP_STACK_GO_BUILDER=onthertech/titan-op-stack-go:latest
FROM $OP_STACK_GO_BUILDER as builder
# See "make golang-docker" and /ops/docker/op-stack-go
FROM alpine:3.18

COPY --from=builder /app/op-wheel/bin/op-wheel /usr/local/bin

CMD ["op-wheel"]
