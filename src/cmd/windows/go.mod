module jenkinsci.org/plugins/durabletask/windows

go 1.24.0

replace jenkinsci.org/plugins/durabletask/common => ../../pkg/common

require (
	golang.org/x/sys v0.38.0
	jenkinsci.org/plugins/durabletask/common v0.0.0-00010101000000-000000000000
)
