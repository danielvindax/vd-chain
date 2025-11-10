# NB: This is a digest for a multi-arch manifest list, you will want to get this by running
# `docker buildx imagetools inspect golang:1.23.1-alpine`
ARG GOLANG_1_23_ALPINE_DIGEST="ac67716dd016429be8d4c2c53a248d7bcdf06d34127d3dc451bda6aa5a87bc06"

# This Dockerfile is a stateless build of the `vindaxd` binary as a Docker container.
# It does not include any configuration, state, or genesis information.

# --------------------------------------------------------
# Builder
# --------------------------------------------------------

FROM golang@sha256:${GOLANG_1_23_ALPINE_DIGEST} as builder
ARG VERSION
ARG COMMIT

RUN set -eux; apk add --no-cache ca-certificates build-base; apk add git linux-headers bash binutils-gold

# Download go dependencies
WORKDIR /vindax
COPY go.* ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download

# Copy the remaining files
COPY . .

# Build vindaxd binary
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go build \
      -mod=readonly \
      -tags "netgo,ledger,muslc" \
      -ldflags "-X github.com/cosmos/cosmos-sdk/version.Name="vindax" \
              -X github.com/cosmos/cosmos-sdk/version.AppName="vindaxd" \
              -X github.com/cosmos/cosmos-sdk/version.Version=$VERSION \
              -X github.com/cosmos/cosmos-sdk/version.Commit=$COMMIT \
              -X github.com/cosmos/cosmos-sdk/version.BuildTags='netgo,ledger,muslc' \
              -w -s -linkmode=external -extldflags '-Wl,-z,muldefs -static'" \
      -trimpath \
      -o /vindax/build/ \
      ./...


# --------------------------------------------------------
# Runner
# --------------------------------------------------------

FROM golang@sha256:${GOLANG_1_23_ALPINE_DIGEST}

RUN apk add --no-cache bash

COPY --from=builder /vindax/build/vindaxd /bin/vindaxd

ENV HOME /vindax
WORKDIR $HOME

# tendermint p2p
EXPOSE 26656
# tendermint rpc
EXPOSE 26657
# rest server
EXPOSE 1317
# grpc
EXPOSE 9090

ENTRYPOINT ["vindaxd"]
