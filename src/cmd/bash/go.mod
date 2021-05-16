module jenkinsci.org/plugins/durabletask/bash

go 1.16

replace jenkinsci.org/plugins/durabletask/common => ../../pkg/common

require (
	golang.org/x/sys v0.0.0-20210507161434-a76c4d0a0096
	jenkinsci.org/plugins/durabletask/common v0.0.0-00010101000000-000000000000
)
