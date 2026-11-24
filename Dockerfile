FROM --platform=$BUILDPLATFORM node:26.13.1-trixie-slim@sha256:e6bedef217677b276c591795e169840edd9096dbf13f29df8ef8d1856b8820bc AS frontend
WORKDIR /src/web
COPY web/package.json web/package-lock.json ./
RUN npx --yes "$(node -p 'require("./package.json").packageManager')" ci
COPY web ./
RUN npx --yes "$(node -p 'require("./package.json").packageManager')" run build

FROM --platform=$BUILDPLATFORM golang:1.27.4-trixie@sha256:f18c3d73a184fb0582828510722e723b63a93f7188d11dd223d9be6f6880e65d AS build
ARG TARGETOS
ARG TARGETARCH
ENV GOTOOLCHAIN=local
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /src/web/dist ./web/dist
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH make build

FROM debian:trixie-20261123-slim@sha256:8aa5a7bf5a1955bddbe7d74fcdd0b91ab30e74362a8b73af672b0481b35b44a0 AS final
# Freeze transitive packages as well as the explicitly installed dependencies.
RUN sed -i \
      -e 's|http://deb.debian.org/debian-security|http://snapshot.debian.org/archive/debian-security/20261124T000000Z|g' \
      -e 's|http://deb.debian.org/debian|http://snapshot.debian.org/archive/debian/20261124T000000Z|g' \
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
