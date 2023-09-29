FROM node:18 as frontend
WORKDIR /src
COPY ./ui ./ui
RUN yarn --cwd /src/ui install
RUN yarn --cwd /src/ui build

FROM golang:1.21.1-alpine AS build
RUN apk add --no-cache ca-certificates git make
WORKDIR /src
COPY ./go.mod ./go.sum ./Makefile ./
COPY ./ ./
RUN make build

FROM alpine:3.18 AS final
RUN apk --no-cache add ca-certificates graphviz
COPY --from=build /src/dist/amflow /bin/amflow
ENTRYPOINT ["/amflow"]
