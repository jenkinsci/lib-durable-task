module jenkinsci.org/plugins/durabletask/windows

go 1.16

replace jenkinsci.org/plugins/durabletask/common => ../../pkg/common

require (
	golang.org/x/sys v0.1.0
	jenkinsci.org/plugins/durabletask/common v0.0.0-00010101000000-000000000000
)
