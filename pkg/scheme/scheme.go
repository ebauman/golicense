package scheme

import (
	v1 "github.com/ebauman/golicense/pkg/apis/golicense.1eb100.net/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var (
	// Scheme is a new (blank?) scheme to which resources will be added
	Scheme = runtime.NewScheme()
	// Codecs is a new codec factory for the scheme
	Codecs = serializer.NewCodecFactory(Scheme)
	// ParameterCodec is a new parameter codec for the scheme
	ParameterCodec = runtime.NewParameterCodec(Scheme)
)

// AddToScheme adds all Resources to the Scheme
// For instance, we are going to take the types from the v1 package and add them to the Scheme
// This allows our api server to recognize the types
func AddToScheme(scheme *runtime.Scheme) error {
	metav1.AddToGroupVersion(scheme, v1.SchemeGroupVersion)

	if err := v1.AddToScheme(scheme); err != nil {
		return err
	}

	if err := corev1.AddToScheme(scheme); err != nil {
		return err
	}

	return nil
}

// We must complete adding the types to the scheme or fail hard (panic)
func init() {
	utilruntime.Must(AddToScheme(Scheme))
}
