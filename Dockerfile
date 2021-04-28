# NOTE: golang > 1.14 drops darwin 32-bit support
#ARG platform_tag
FROM golang:1.14-${image.platform}
WORKDIR /durabletask
ENTRYPOINT /bin/sh ${entrypoint.script} ${project.version}
