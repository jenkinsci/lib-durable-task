ARG BASE_DIR=/durabletask
ARG NAME=durable_task_monitor

FROM golang:1.24.3-nanoserver-1809 AS builder
ARG BASE_DIR
ARG NAME
COPY cmd $BASE_DIR/cmd
COPY pkg $BASE_DIR/pkg
WORKDIR $BASE_DIR/pkg/common
RUN go test -v
WORKDIR $BASE_DIR/cmd/bash
# go mod tidy fails in windows CI due to permissions of module cache. Each dependency needs to be fetched with -modcacherw flag
# see golang.org/issue/31481
# RUN go mod tidy
RUN go get -modcacherw jenkinsci.org/plugins/durabletask/common
# TODO find some way to inject specific version of x/sys based on go.mod definitions
RUN go get -modcacherw golang.org/x/sys@v0.38.0
# can't test bash on windows
RUN set CGO_ENABLED=0&& set GOOS=darwin&& set GOARCH=amd64&& go build -a -ldflags="-s -w" -trimpath -o %NAME%_darwin_amd64
RUN set CGO_ENABLED=0&& set GOOS=darwin&& set GOARCH=arm64&& go build -a -ldflags="-s -w" -trimpath -o %NAME%_darwin_arm64
RUN set CGO_ENABLED=0&& set GOOS=linux&& set GOARCH=amd64&& go build -a -ldflags="-s -w" -trimpath -o %NAME%_linux_64
RUN set CGO_ENABLED=0&& set GOOS=linux&& set GOARCH=386&& go build -a -ldflags="-s -w" -trimpath -o %NAME%_linux_32
RUN set CGO_ENABLED=0&& set GOOS=linux&& set GOARCH=ppc64le&& go build -a -ldflags="-s -w" -trimpath -o %NAME%_linux_ppc64le
RUN set CGO_ENABLED=0&& set GOOS=linux&& set GOARCH=arm64&& go build -a -ldflags="-s -w" -trimpath -o %NAME%_linux_aarch64
WORKDIR $BASE_DIR/cmd/windows
# RUN go mod tidy
RUN go get -modcacherw jenkinsci.org/plugins/durabletask/common
RUN go get -modcacherw golang.org/x/sys@v0.38.0
RUN go test -v
RUN set CGO_ENABLED=0& set GOOS=windows& set GOARCH=amd64& go build -a -ldflags="-s -w" -trimpath -o %NAME%_win_64.exe
RUN set CGO_ENABLED=0& set GOOS=windows& set GOARCH=386& go build -a -ldflags="-s -w" -trimpath -o %NAME%_win_32.exe

# TODO: uncomment once docker build --output (i.e. BuildKit) is available for windows containers
# see: https://github.com/microsoft/Windows-Containers/issues/34
# FROM mcr.microsoft.com/windows/nanoserver:20H2 as export-stage
# ARG BASE_DIR
# ARG NAME
# WORKDIR $BASE_DIR
# COPY --from=builder $BASE_DIR/cmd/bash/${NAME}_darwin_amd64 $BASE_DIR
# COPY --from=builder $BASE_DIR/cmd/bash/${NAME}_darwin_arm64 $BASE_DIR
# COPY --from=builder $BASE_DIR/cmd/bash/${NAME}_linux_64 $BASE_DIR
# COPY --from=builder $BASE_DIR/cmd/bash/${NAME}_linux_32 $BASE_DIR
