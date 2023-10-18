package mysql

import (
	"github.com/acorn-io/mink/pkg/db"
	"github.com/acorn-io/mink/pkg/serializer"
	"github.com/acorn-io/mink/pkg/strategy"
	v1 "github.com/ebauman/golicense/pkg/apis/golicense.1eb100.net/v1"
	golicenseScheme "github.com/ebauman/golicense/pkg/scheme"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

func APIGroups(dsn string) (*genericapiserver.APIGroupInfo, error) {
	scheme := golicenseScheme.Scheme
	dbFactory, err := db.NewFactory(scheme, dsn)
	if err != nil {
		return nil, err
	}

	var authorityStrategy, licenseStrategy, licenseeStrategy, certificateStrategy, productStrategy strategy.CompleteStrategy
	authorityStrategy, err = dbFactory.NewDBStrategy(&v1.Authority{})
	licenseStrategy, err = dbFactory.NewDBStrategy(&v1.License{})
	licenseeStrategy, err = dbFactory.NewDBStrategy(&v1.Licensee{})
	certificateStrategy, err = dbFactory.NewDBStrategy(&v1.Certificate{})
	productStrategy, err = dbFactory.NewDBStrategy(&v1.Product{})

	authorityStorage, err := NewAuthorityStorage(authorityStrategy, productStrategy)
	if err != nil {
		return nil, err
	}

	licenseStorage, err := NewLicenseStorage(licenseStrategy)
	if err != nil {
		return nil, err
	}

	licenseeStorage, err := NewLicenseeStorage(licenseeStrategy, authorityStrategy)
	if err != nil {
		return nil, err
	}

	certificateStorage, err := NewCertificateStorage(certificateStrategy)
	if err != nil {
		return nil, err
	}

	productStorage, err := NewProductStorage(productStrategy)
	if err != nil {
		return nil, err
	}

	stores := map[string]rest.Storage{
		"authorities":  authorityStorage,
		"licenses":     licenseStorage,
		"licensees":    licenseeStorage,
		"certificates": certificateStorage,
		"products":     productStorage,
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
