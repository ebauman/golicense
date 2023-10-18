package mysql

import (
	"github.com/acorn-io/mink/pkg/stores"
	"github.com/acorn-io/mink/pkg/strategy"
	v1 "github.com/ebauman/golicense/pkg/apis/golicense.1eb100.net/v1"
	"k8s.io/apiserver/pkg/registry/rest"
)

func NewLicenseStorage(licenseStrategy strategy.CompleteStrategy) (rest.Storage, error) {
	return stores.NewBuilder(licenseStrategy.Scheme(), &v1.License{}).
		WithCompleteCRUD(licenseStrategy).
		Build(), nil
}
