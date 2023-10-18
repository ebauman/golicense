//go:generate go run github.com/acorn-io/baaah/cmd/deepcopy ./pkg/apis/golicense.1eb100.net/v1/
//go:generate go run k8s.io/kube-openapi/cmd/openapi-gen -i github.com/ebauman/golicense/pkg/apis/golicense.1eb100.net/v1,k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/runtime,k8s.io/apimachinery/pkg/version,k8s.io/apimachinery/pkg/api/resource,k8s.io/api/core/v1,k8s.io/api/rbac/v1,k8s.io/apimachinery/pkg/util/intstr -p github.com/ebauman/golicense/pkg/openapi -h boilerplate/header.txt

package main

import (
	_ "github.com/acorn-io/baaah/pkg/deepcopy"
	_ "github.com/golang/mock/gomock"
	_ "k8s.io/kube-openapi/cmd/openapi-gen/args"
)
