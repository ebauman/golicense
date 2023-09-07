//go:generate go run github.com/acorn-io/baaah/cmd/deepcopy ./pkg/apis/golicense.1eb100.net/v1/
//go:generate go run k8s.io/kube-openapi/cmd/openapi-gen -i github.com/ebauman/golicense/pkg/apis/golicense.1eb100.net/v1 -p ./pkg/openapi/generate -h boilerplate/header.txt

package main

import (
	_ "github.com/acorn-io/baaah/pkg/deepcopy"
	_ "github.com/golang/mock/gomock"
	_ "k8s.io/kube-openapi/cmd/openapi-gen/args"
)
