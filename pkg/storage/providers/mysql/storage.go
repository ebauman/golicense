package mysql

import (
	"github.com/acorn-io/mink/pkg/db"
	"github.com/acorn-io/mink/pkg/serializer"
	v1 "github.com/ebauman/golicense/pkg/apis/golicense.1eb100.net/v1"
	golicenseScheme "github.com/ebauman/golicense/pkg/scheme"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

func APIGroups(dsn string) (*genericapiserver.APIGroupInfo, error) {
	scheme := golicenseScheme.Scheme
	dbFactory := db.NewFactory(scheme, dsn)

	authorityStorage, err := NewAuthorityStorage(dbFactory)
	if err != nil {
		return nil, err
	}

	stores := map[string]rest.Storage{
		"authorities": authorityStorage,
	}
	err = v1.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}

	err = v1.AddToScheme(scheme) // should this be AddToSchemeWithGV? See github.com/acorn-io/runtime/pkg/server/registry/apigroups/acorn/apigroup.go
	if err != nil {
		return nil, err
	}

	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(v1.SchemeGroupVersion.Group, scheme, golicenseScheme.ParameterCodec, golicenseScheme.Codecs)
	apiGroupInfo.VersionedResourcesStorageMap["v1"] = stores
	apiGroupInfo.NegotiatedSerializer = serializer.NewNoProtobufSerializer(apiGroupInfo.NegotiatedSerializer)

	return &apiGroupInfo, nil
}
