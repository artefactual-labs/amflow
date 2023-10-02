FROM node:18 as frontend
WORKDIR /src
COPY web /src/web
RUN npm --cwd /src/web ci
RUN npm --cwd /src/web run build

FROM golang:1.21.1-alpine AS build
RUN apk add --no-cache ca-certificates git make
WORKDIR /src
COPY --from=frontend /src/web/dist ./web/dist
COPY . .
RUN make build

FROM debian:12 AS final
RUN apt-get update && apt-get install --yes --no-install-recommends graphviz ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=build /src/dist/amflow /bin/amflow
ENTRYPOINT ["/bin/amflow"]
