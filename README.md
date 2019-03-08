# amflow

## Installation

### Local installation

Download the [latest binary][0] and simply run it, e.g.:

    $ ./amflow help

amflow needs `dot` installed (part of [GraphViz](https://www.graphviz.org/)).

### Via Docker

You can also use the [Docker image][1], e.g.:

    $ docker run --rm sevein/amflow:latest help

## Usage examples

Serve the latest workflow document found in Archivematica's GitHub repository. It should be accessible at http://127.0.0.1:2323.

    $ docker run -p 2323:2323 --rm sevein/amflow:latest edit --latest

## API docs

See http://petstore.swagger.io/?url=https://raw.githubusercontent.com/artefactual-labs/amflow/main/public/swagger/swagger.yaml.

[0]: https://github.com/artefactual-labs/amflow/releases/latest
[1]: https://hub.docker.com/r/sevein/amflow/tags
