# amflow

## Installation

### Local installation

Download the [latest binary][0] and simply run it, e.g.:

    ./amflow help

amflow needs `dot` installed (part of [GraphViz](https://www.graphviz.org/)).

### Via Docker

You can also use the [Docker image][1], e.g.:

    $ docker run --rm artefactual/amflow:latest -v warn help
    A tool that facilitates workflow editing for Archivematica.

    Usage:
      amflow [command]

    Available Commands:
      check       Verify workflow integrity
      edit        Edit the workflow
      export      Export the workflow in DOT format
      help        Help about any command
      search      Experimenal solution to do searches in the workflow graph
      version     Print the version information

    Flags:
      -h, --help               help for amflow
      -v, --verbosity string   Log level (debug, info, warn, error, fatal, panic (default "info")

    Use "amflow [command] --help" for more information about a command.

## Usage

The following examples use Docker so you don't have to install amflow locally.

Serve the latest workflow document found in Archivematica's GitHub repository. It should be accessible at http://127.0.0.1:2323.

    docker run -p 2323:2323 --rm artefactual/amflow:latest edit --latest

It is also possible to point to a local or remote workflow document, e.g.:

    docker run -p 2323:2323 --rm artefactual/amflow:latest edit --file=https://raw.githubusercontent.com/artefactual/archivematica/stable/1.10.x/src/MCPServer/lib/assets/workflow.json

Check the integrity of the workflow, e.g.:

    $ docker run --rm artefactual/amflow:latest check --latest
    INFO[0000] amflow (dev)
    INFO[0000] Downloading workfow                           mode=file source="https://raw.githubusercontent.com/artefactual/archivematica/qa/1.x/src/MCPServer/lib/assets/workflow.json"
    WARN[0001] Unhealthy workflow warning                    err="[/system/createAIC/] watched directory is not referenced"
    WARN[0001] Unhealthy workflow warning                    err="[/system/reingestAIP/] watched directory is not referenced"
    WARN[0001] Unhealthy workflow warning                    err="[653b134f-4a37-4578-a286-7f2072e89f9e] link is terminal but has alternative paths [children=1] [refsWD=false]"
    WARN[0001] Unhealthy workflow warning                    err="[16415d2f-5642-496d-a46d-00028ef6eb0a] link is terminal but has alternative paths [children=2] [refsWD=true]"
    WARN[0001] Unhealthy workflow warning                    err="[abd6d60c-d50f-4660-a189-ac1b34fafe85] link is terminal but has alternative paths [children=1] [refsWD=false]"

## Limitations

* Web interface does not provide searching capabilities.

* Web interface hides common links like "Email fail report" in order to speed up the rendering process. Open up `dotviz.go` to know more.
  A workaround is to export the DOT graph with `amflow export --full --format=dot`, then render it locally with GraphViz. E.g.:

  ```amflow export --full --format=dot | dot -v -Tx11```

  It can take a couple of minutes! `-Tlang` determines the output format, `-Tx11` opens an interactive graph viewer based on Xlib canvas supporting panning and zooming functions.

[0]: https://github.com/artefactual-labs/amflow/releases/latest
[1]: https://hub.docker.com/r/artefactual/amflow/tags
