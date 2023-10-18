package v1

import (
	golicense_1eb100_net "github.com/ebauman/golicense/pkg/apis/golicense.1eb100.net"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const Version = "v1"

var SchemeGroupVersion = schema.GroupVersion{
	Group:   golicense_1eb100_net.Group,
	Version: Version,
}

func AddToScheme(scheme *runtime.Scheme) error {
	return AddToSchemeWithGV(scheme, SchemeGroupVersion)
}

func AddToSchemeWithGV(scheme *runtime.Scheme, schemeGroupVersion schema.GroupVersion) error {

	scheme.AddKnownTypes(schemeGroupVersion,
		&License{},
		&LicenseList{},
		&Authority{},
		&AuthorityList{},
		&Licensee{},
		&LicenseeList{},
		&Certificate{},
		&CertificateList{},
		&Product{},
		&ProductList{},
	)

	scheme.AddKnownTypes(schemeGroupVersion, &metav1.Status{})

	return nil
}
