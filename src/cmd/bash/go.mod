module jenkinsci.org/plugins/durabletask/bash

go 1.20

replace jenkinsci.org/plugins/durabletask/common => ../../pkg/common

require (
	golang.org/x/sys v0.7.0
	jenkinsci.org/plugins/durabletask/common v0.0.0-00010101000000-000000000000
)
