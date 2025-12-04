ARG BASE_DIR=/durabletask
ARG NAME=durable_task_monitor

FROM golang:1.25.5-bookworm AS builder
ARG BASE_DIR
ARG NAME
COPY cmd $BASE_DIR/cmd
COPY pkg $BASE_DIR/pkg
WORKDIR $BASE_DIR/pkg/common
RUN go mod tidy
RUN go test -v
WORKDIR $BASE_DIR/cmd/bash
RUN go mod tidy
RUN go test -v
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -ldflags="-s -w" -trimpath -o ${NAME}_darwin_amd64
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -ldflags="-s -w" -trimpath -o ${NAME}_darwin_arm64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-s -w" -trimpath -o ${NAME}_linux_64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -a -ldflags="-s -w" -trimpath -o ${NAME}_linux_32
RUN CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le go build -a -ldflags="-s -w" -trimpath -o ${NAME}_linux_ppc64le
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -ldflags="-s -w" -trimpath -o ${NAME}_linux_aarch64
WORKDIR $BASE_DIR/cmd/windows
RUN go mod tidy
# can't test windows on linux
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -ldflags="-s -w" -trimpath -o ${NAME}_win_64.exe
RUN CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -a -ldflags="-s -w" -trimpath -o ${NAME}_win_32.exe

FROM scratch AS export-stage
ARG BASE_DIR
ARG NAME
COPY --from=builder $BASE_DIR/cmd/bash/${NAME}_*_* /
COPY --from=builder $BASE_DIR/cmd/windows/${NAME}_*_* /
