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

type authorityValidator struct {
	productGetter strategy.Getter
}

// NewAuthorityStorage receives a factory
func NewAuthorityStorage(authorityStrategy strategy.CompleteStrategy, productGetter strategy.Getter) (rest.Storage, error) {
	authValidator := authorityValidator{
		productGetter: productGetter,
	}
	return stores.NewBuilder(authorityStrategy.Scheme(), &v1.Authority{}).
		WithCompleteCRUD(authorityStrategy).
		WithValidateCreate(authValidator).
		WithValidateUpdate(authValidator).
		Build(), nil
}

func (av authorityValidator) Validate(ctx context.Context, obj runtime.Object) (result field.ErrorList) {
	result = append(result, av.validateProducts(ctx, obj)...)

	return
}

func (av authorityValidator) ValidateUpdate(ctx context.Context, obj runtime.Object, _ runtime.Object) (result field.ErrorList) {
	result = append(result, av.validateProducts(ctx, obj)...)

	return
}

func (av authorityValidator) validateProducts(ctx context.Context, obj runtime.Object) (result field.ErrorList) {
	authority := obj.(*v1.Authority)

	// ensure that all products on the authority exist
	for _, prod := range authority.Spec.Products {
		if _, err := av.productGetter.Get(ctx, authority.GetNamespace(), prod); err != nil {
			result = append(result, field.Invalid(field.NewPath("spec", "products"), prod, err.Error()))
			return
		}
	}

	return
}
