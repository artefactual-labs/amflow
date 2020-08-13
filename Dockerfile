FROM node:11 as frontend
WORKDIR /src
COPY ./ui ./ui
RUN yarn --cwd /src/ui install
RUN yarn --cwd /src/ui build

ARG GO_VERSION=1.15
FROM golang:${GO_VERSION}-alpine AS build
RUN apk add --no-cache ca-certificates git make
WORKDIR /src
COPY ./go.mod ./go.sum ./Makefile ./
RUN make deps
COPY ./ ./
RUN make build

FROM alpine:3.11 AS final
RUN apk --no-cache add ca-certificates graphviz
COPY --from=build /tmp/amflow /bin/amflow
ENTRYPOINT ["/amflow"]
