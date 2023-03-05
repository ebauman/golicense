package authority

import (
	"github.com/acorn-io/mink/pkg/stores"
	"github.com/acorn-io/mink/pkg/strategy/remote"
	"k8s.io/apiserver/pkg/registry/rest"
)

func NewStorage() rest.Storage {
	remote.NewRemote()

	stores.NewBuilder()
}
