# automatically set by buildkit, can be changed with --platform flag
# see https://docs.docker.com/engine/reference/builder/#automatic-platform-args-in-the-global-scope
# TARGETOS
# TARGETARCH
# TARGETPLATFORM
# BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

# All target images use this AS base image, and add the final build results.
# It will default to the target platform.
ARG TARGET_BASE_IMAGE=alpine:3.18

# We may be cross-building for another platform. Specify which platform we need AS builder.
FROM --platform=$BUILDPLATFORM golang:1.21.3-alpine3.18 AS builder

RUN apk add --no-cache make gcc musl-dev linux-headers git jq bash

# We copy the go.mod/sum first, so the `go mod download` does not have to re-run if dependencies do not change.
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app

RUN echo "go mod cache: $(go env GOMODCACHE)"
RUN echo "go build cache: $(go env GOCACHE)"

# warm-up the cache
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build go mod download

# NOTE: the Dockerfile.dockerignore file effectively describes all dependencies
COPY . /app

# We avoid copying the full .git dir into the build for just some metadata.
# Instead, specify:
# --build-arg GIT_COMMIT=$(git rev-parse HEAD)
# --build-arg GIT_DATE=$(git show -s --format='%ct')
ARG GIT_COMMIT
ARG GIT_DATE

# separate docker-builds:
# - proxyd
# - any JS/TS/smart-contract builds

# TODO: if the build involves optional CGO, then:
# 1) we may have to disable CGO explicitly
# 3) apk install cross-build gcc,
#    and specify to go to use that, with the CC flag

ARG TARGETOS
ARG TARGETARCH

# Build the Go services, utilizing caches and share the many common packages.
# The "id" defaults to the value of "target", the cache will thus be reused during this build.
# "sharing" defaults to "shared", the cache will thus be available to other concurrent docker builds.
FROM --platform=$BUILDPLATFORM builder AS cannon-builder
ARG CANNON_VERSION=v0.0.0
RUN --mount=type=cache,target=/root/.cache/go-build cd cannon && make cannon  \
    GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$CANNON_VERSION"

# Build op-program first (needed by asterisc for prestate generation)
FROM --platform=$BUILDPLATFORM builder AS op-program-builder
ARG OP_PROGRAM_VERSION=v0.0.0
# Build host (for running AS server)
RUN --mount=type=cache,target=/root/.cache/go-build cd op-program && make op-program-host  \
    GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_PROGRAM_VERSION"
# Build MIPS client for prestate generation (required for fault proofs)
RUN --mount=type=cache,target=/root/.cache/go-build cd op-program && make op-program-client-mips  \
    GOOS=linux GOARCH=mips GOMIPS=softfloat GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_PROGRAM_VERSION"
# Build RISC-V client for Asterisc prestate generation (GameType 2, 3)
RUN --mount=type=cache,target=/root/.cache/go-build cd op-program && make op-program-client-riscv  \
    GOOS=linux GOARCH=riscv64 GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_PROGRAM_VERSION"
# Generate Cannon prestate (GameType 0, 1)
RUN /app/cannon/bin/cannon load-elf --path /app/op-program/bin/op-program-client.elf --out /app/op-program/bin/prestate.json --meta ""
RUN /app/cannon/bin/cannon run --proof-at '=0' --stop-at '=1' --input /app/op-program/bin/prestate.json --meta "" --proof-fmt '/app/op-program/bin/%d.json' --output "" && \
    mv /app/op-program/bin/0.json /app/op-program/bin/prestate-proof.json

# Build asterisc VM (RISC-V) for GameType 2 and 3 (depends on op-program-builder)
FROM --platform=$BUILDPLATFORM op-program-builder AS asterisc-builder
ARG ASTERISC_VERSION=v0.0.0
# Build asterisc binary and generate prestate (only if source code exists)
# Note: If asterisc source is not present, users should use pull-vm-images.sh to download pre-built binaries
RUN if [ -d "/app/asterisc" ] && [ -f "/app/asterisc/Makefile" ]; then \
        echo "Building asterisc from source..."; \
        cd asterisc && make asterisc \
        GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$ASTERISC_VERSION" && \
        echo "Generating Asterisc prestate..." && \
        /app/asterisc/bin/asterisc load-elf --path /app/op-program/bin/op-program-client-riscv.elf --out /app/asterisc/bin/prestate.json --meta "" && \
        /app/asterisc/bin/asterisc run --proof-at '=0' --stop-at '=1' --input /app/asterisc/bin/prestate.json --meta "" --proof-fmt '/app/asterisc/bin/%d.json' --output "" && \
        mv /app/asterisc/bin/0.json /app/asterisc/bin/prestate-proof.json && \
        echo "Asterisc build and prestate generation complete"; \
    else \
        echo "asterisc source not found, skipping build"; \
        echo "Users should download pre-built binaries via pull-vm-images.sh"; \
        mkdir -p /app/asterisc/bin; \
    fi

FROM --platform=$BUILDPLATFORM builder AS op-heartbeat-builder
ARG OP_HEARTBEAT_VERSION=v0.0.0
RUN --mount=type=cache,target=/root/.cache/go-build cd op-heartbeat && make op-heartbeat  \
    GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_HEARTBEAT_VERSION"

FROM --platform=$BUILDPLATFORM builder AS op-wheel-builder
ARG OP_WHEEL_VERSION=v0.0.0
RUN --mount=type=cache,target=/root/.cache/go-build cd op-wheel && make op-wheel  \
    GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE  VERSION="$OP_WHEEL_VERSION"

FROM --platform=$BUILDPLATFORM builder AS op-node-builder
ARG OP_NODE_VERSION=v0.0.0
RUN echo TARGETOS=$TARGETOS TARGETARCH=$TARGETARCH BUILDPLATFORM=$BUILDPLATFORM TARGETPLATFORM=$TARGETPLATFORM OP_NODE_VERSION=$OP_NODE_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build cd op-node && make op-node  \
    GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE VERSION="$OP_NODE_VERSION"

FROM --platform=$BUILDPLATFORM builder AS op-challenger-builder
ARG OP_CHALLENGER_VERSION=v0.0.0
RUN --mount=type=cache,target=/root/.cache/go-build cd op-challenger && make op-challenger  \
    GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE  VERSION="$OP_CHALLENGER_VERSION"

FROM --platform=$BUILDPLATFORM builder AS op-dispute-mon-builder
ARG OP_DISPUTE_MON_VERSION=v0.0.0
RUN --mount=type=cache,target=/root/.cache/go-build cd op-dispute-mon && make op-dispute-mon  \
    GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE  VERSION="$OP_DISPUTE_MON_VERSION"

FROM --platform=$BUILDPLATFORM builder AS op-batcher-builder
ARG OP_BATCHER_VERSION=v0.0.0
RUN --mount=type=cache,target=/root/.cache/go-build cd op-batcher && make op-batcher  \
    GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE  VERSION="$OP_BATCHER_VERSION"

FROM --platform=$BUILDPLATFORM builder AS op-proposer-builder
ARG OP_PROPOSER_VERSION=v0.0.0
RUN --mount=type=cache,target=/root/.cache/go-build cd op-proposer && make op-proposer  \
    GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE  VERSION="$OP_PROPOSER_VERSION"

FROM --platform=$BUILDPLATFORM builder AS op-conductor-builder
ARG OP_CONDUCTOR_VERSION=v0.0.0
RUN --mount=type=cache,target=/root/.cache/go-build cd op-conductor && make op-conductor  \
    GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE  VERSION="$OP_CONDUCTOR_VERSION"

FROM --platform=$BUILDPLATFORM builder AS da-server-builder
RUN --mount=type=cache,target=/root/.cache/go-build cd op-plasma && make da-server  \
    GOOS=$TARGETOS GOARCH=$TARGETARCH GITCOMMIT=$GIT_COMMIT GITDATE=$GIT_DATE

FROM --platform=$TARGETPLATFORM $TARGET_BASE_IMAGE AS cannon-target
COPY --from=cannon-builder /app/cannon/bin/cannon /usr/local/bin/
CMD ["cannon"]

FROM --platform=$TARGETPLATFORM $TARGET_BASE_IMAGE AS op-program-target
COPY --from=op-program-builder /app/op-program/bin/op-program /usr/local/bin/
CMD ["op-program"]

FROM --platform=$TARGETPLATFORM $TARGET_BASE_IMAGE AS op-heartbeat-target
COPY --from=op-heartbeat-builder /app/op-heartbeat/bin/op-heartbeat /usr/local/bin/
CMD ["op-heartbeat"]

FROM --platform=$TARGETPLATFORM $TARGET_BASE_IMAGE AS op-wheel-target
COPY --from=op-wheel-builder /app/op-wheel/bin/op-wheel /usr/local/bin/
CMD ["op-wheel"]

FROM --platform=$TARGETPLATFORM $TARGET_BASE_IMAGE AS op-node-target
COPY --from=op-node-builder /app/op-node/bin/op-node /usr/local/bin/
CMD ["op-node"]

FROM --platform=$TARGETPLATFORM $TARGET_BASE_IMAGE AS op-challenger-target
COPY --from=op-challenger-builder /app/op-challenger/bin/op-challenger /usr/local/bin/
# Make the bundled op-program the default cannon server
COPY --from=op-program-builder /app/op-program/bin/op-program /usr/local/bin/
ENV OP_CHALLENGER_CANNON_SERVER=/usr/local/bin/op-program
# Make the bundled cannon the default cannon executable
COPY --from=cannon-builder /app/cannon/bin/cannon /usr/local/bin/
ENV OP_CHALLENGER_CANNON_BIN=/usr/local/bin/cannon
# NOTE: Asterisc binary is provided via volume mount (see docker-compose-full.yml)
# COPY --from=asterisc-builder /app/asterisc/bin/asterisc /usr/local/bin/
# ENV OP_CHALLENGER_ASTERISC_BIN=/usr/local/bin/asterisc
# Asterisc uses op-program AS server (same AS cannon)
ENV OP_CHALLENGER_ASTERISC_SERVER=/usr/local/bin/op-program
CMD ["op-challenger"]

FROM --platform=$TARGETPLATFORM $TARGET_BASE_IMAGE AS op-dispute-mon-target
COPY --from=op-dispute-mon-builder /app/op-dispute-mon/bin/op-dispute-mon /usr/local/bin/
CMD ["op-dispute-mon"]

FROM --platform=$TARGETPLATFORM $TARGET_BASE_IMAGE AS op-batcher-target
COPY --from=op-batcher-builder /app/op-batcher/bin/op-batcher /usr/local/bin/
CMD ["op-batcher"]

FROM --platform=$TARGETPLATFORM $TARGET_BASE_IMAGE AS op-proposer-target
COPY --from=op-proposer-builder /app/op-proposer/bin/op-proposer /usr/local/bin/
CMD ["op-proposer"]

FROM --platform=$TARGETPLATFORM $TARGET_BASE_IMAGE AS op-conductor-target
COPY --from=op-conductor-builder /app/op-conductor/bin/op-conductor /usr/local/bin/
CMD ["op-conductor"]

FROM --platform=$TARGETPLATFORM $TARGET_BASE_IMAGE AS da-server-target
COPY --from=da-server-builder /app/op-plasma/bin/da-server /usr/local/bin/
CMD ["da-server"]
