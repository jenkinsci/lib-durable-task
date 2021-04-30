module jenkinsci.org/plugins/durabletask/bash

go 1.14

replace jenkinsci.org/plugins/durabletask/common => ../../pkg/common

require (
	// pin x/sys to 1.14 by manually running: go get golang.org/x/sys@release-branch.go1.14-std
	golang.org/x/sys v0.0.0-20200201011859-915c9c3d4ccf
	jenkinsci.org/plugins/durabletask/common v0.0.0-00010101000000-000000000000
)
