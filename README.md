# lib-durable-task

This repository provides the binaries that launch the scripts from the [`durable-task-plugin`](https://github.com/jenkinsci/durable-task-plugin).

The binary source is written in Golang and pinned to version 1.14 as that is the last version that supports 32-bit darwin builds.
See the `src/cmd` folder for the `main` packages.

The code is cross-compiled via docker, and the compiled binaries are then packaged into a jar for use as a dependency in `durable-task`.