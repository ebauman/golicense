package mysql

import (
	"context"
	"github.com/acorn-io/mink/pkg/stores"
	"github.com/acorn-io/mink/pkg/strategy"
	v1 "github.com/ebauman/golicense/pkg/apis/golicense.1eb100.net/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/rest"
)

type licenseeValidator struct {
	authorityGetter strategy.Getter
}

func NewLicenseeStorage(licenseeStrategy strategy.CompleteStrategy, authorityGetter strategy.Getter) (rest.Storage, error) {
	lv := licenseeValidator{authorityGetter: authorityGetter}

	return stores.NewBuilder(licenseeStrategy.Scheme(), &v1.Licensee{}).
		WithCompleteCRUD(licenseeStrategy).WithValidateCreate(lv).
		Build(), nil
}

func (lv licenseeValidator) Validate(ctx context.Context, obj runtime.Object) (result field.ErrorList) {
	// validation for creation
	licensee := obj.(*v1.Licensee)

	// the licensee must belong to an existing authority
	if _, err := lv.authorityGetter.Get(ctx, licensee.GetNamespace(), licensee.Spec.Authority); err != nil {
		result = append(result, field.Invalid(field.NewPath("spec", "authority"),
			licensee.Spec.Authority,
			err.Error()))
		return
	}

	return
}
