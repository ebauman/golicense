package server

import (
	"github.com/acorn-io/mink/pkg/server"
	"github.com/ebauman/golicense/pkg/scheme"
	"github.com/ebauman/golicense/pkg/storage/providers/mysql"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

func New(httpPort int, httpsPort int, dsn string) (*server.Server, error) {
	apiGroups, err := mysql.APIGroups(dsn)
	if err != nil {
		return nil, err
	}

	return server.New(&server.Config{
		Authenticator:        nil,
		Authorization:        nil,
		HTTPListenPort:       httpPort,
		HTTPSListenPort:      httpsPort,
		LongRunningVerbs:     []string{"watch"},
		LongRunningResources: nil,
		OpenAPIConfig:        nil,
		Scheme:               scheme.Scheme,
		CodecFactory:         &scheme.Codecs,
		APIGroups:            []*genericapiserver.APIGroupInfo{apiGroups},
		Middleware:           nil,
		PostStartFunc:        nil,
	})
}
