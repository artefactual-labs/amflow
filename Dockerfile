FROM --platform=$BUILDPLATFORM node:26.12.0-trixie-slim@sha256:e556a144463cc8dc53afe8de6d2fa0db5c132fd39ffa232c17641bd40df2bf11 AS frontend
WORKDIR /src/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web ./
RUN npm run build

FROM --platform=$BUILDPLATFORM golang:1.27.3-trixie@sha256:fbb7ae40b62ff8afe31662df570d39b4f5052ac32b90dae8e60d7657a20a74f4 AS build
ARG TARGETOS
ARG TARGETARCH
ENV GOTOOLCHAIN=local
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /src/web/dist ./web/dist
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH make build

FROM debian:trixie-20261020-slim@sha256:62bd423aeee912639ca8ee681c2b762cb3e25137e0a92b20be455898a9f034ad AS final
# Freeze transitive packages as well as the explicitly installed dependencies.
RUN sed -i \
      -e 's|http://deb.debian.org/debian-security|http://snapshot.debian.org/archive/debian-security/20261026T000000Z|g' \
      -e 's|http://deb.debian.org/debian|http://snapshot.debian.org/archive/debian/20261026T000000Z|g' \
      /etc/apt/sources.list.d/debian.sources \
    && apt-get -o Acquire::Check-Valid-Until=false update \
    && DEBIAN_FRONTEND=noninteractive apt-get install --yes --no-install-recommends \
      graphviz=2.42.4-3 \
      ca-certificates=20250419+deb13u1 \
    && rm -rf /var/lib/apt/lists/* \
    && useradd --uid 10001 --create-home --home-dir=/home/amflow --user-group --shell=/usr/sbin/nologin amflow
COPY --from=build /src/dist/amflow /bin/amflow
USER amflow
ENTRYPOINT ["/bin/amflow"]
