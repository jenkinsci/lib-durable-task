# lib-durable-task

This repository provides the binaries that launch the scripts from [`durable-task`](https://github.com/jenkinsci/durable-task-plugin).

The binary source is written in Golang. See the `src/cmd` folder for the `main` packages.

The code is cross-compiled via docker, and the compiled binaries are packaged into the library for use as a dependency in `durable-task`.