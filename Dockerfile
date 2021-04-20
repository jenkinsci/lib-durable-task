# NOTE: golang > 1.14 drops darwin 32-bit support
FROM golang:1.14-buster
WORKDIR /durabletask
ENTRYPOINT /bin/sh compile-binaries.sh ${project.version}
