ARG OP_STACK_GO_BUILDER=tokamaknetwork/thanos-op-stack-go:latest
FROM $OP_STACK_GO_BUILDER as builder
# See "make golang-docker" and /ops/docker/op-stack-go

FROM alpine:3.18

COPY --from=builder /usr/local/bin/op-proposer /usr/local/bin/op-proposer

CMD ["op-proposer"]
