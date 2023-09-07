package mysql

import (
	"github.com/acorn-io/mink/pkg/db"
	"github.com/acorn-io/mink/pkg/stores"
	v1 "github.com/ebauman/golicense/pkg/apis/golicense.1eb100.net/v1"
	"k8s.io/apiserver/pkg/registry/rest"
)

//func blah() error {
//	s := runtime.Scheme{}
//
//	gv := schema.GroupVersion{
//		Group:   "golicense.1eb100.net",
//		Version: "v1",
//	}
//
//	s.AddKnownTypes(gv)
//
//	minkFactory := db.NewFactory(&s, "sqlite://data.db")
//
//	authorityStrategy, err := minkFactory.NewDBStrategy(&v1.Authority{})
//	if err != nil {
//		return err
//	}
//
//	authorityStrategy.
//}

// NewAuthorityStorage receives a factory
func NewAuthorityStorage(factory *db.Factory) (rest.Storage, error) {
	authorityStrategy, err := factory.NewDBStrategy(&v1.Authority{})
	if err != nil {
		return nil, err
	}

	return stores.NewBuilder(factory.Scheme(), &v1.Authority{}).
		WithCompleteCRUD(authorityStrategy).
		Build(), nil
}
