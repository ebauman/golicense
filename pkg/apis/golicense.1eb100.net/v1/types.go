package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type License struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec LicenseSpec `json:"spec"`
}

// +k8s:deepcopy-gen=true

type LicenseSpec struct {
	Id          string            `json:"id"`
	Licensee    string            `json:"licensee"`
	Metadata    map[string]string `json:"metadata"`
	Grants      map[string]int    `json:"grants"`
	NotBefore   metav1.Time       `json:"notBefore"`
	NotAfter    metav1.Time       `json:"notAfter"`
	Key         string            `json:"-"`
	Certificate string            `json:"-"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type LicenseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []License `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Authority struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AuthoritySpec `json:"spec"`
}

// +k8s:deepcopy-gen=true

type AuthoritySpec struct {
	Id       string            `json:"id"`
	Name     string            `json:"name"`
	Products []string          `json:"products"`
	Metadata map[string]string `json:"metadata"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AuthorityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Authority `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Licensee struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec LicenseeSpec `json:"spec"`
}

// +k8s:deepcopy-gen=true

type LicenseeSpec struct {
	Id        string            `json:"id"`
	Name      string            `json:"name"`
	Authority string            `json:"authority"`
	Metadata  map[string]string `json:"metadata"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type LicenseeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Licensee `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Certificate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CertificateSpec `json:"spec"`
}

// +k8s:deepcopy-gen=true

type CertificateSpec struct {
	Id              string            `json:"id"`
	Authority       string            `json:"authority"`
	PrivateKeyBytes []byte            `json:"privateKey"`
	Metadata        map[string]string `json:"metadata"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CertificateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Certificate `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Product struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ProductSpec `json:"spec"`
}

// +k8s:deepcopy-gen=true

type ProductSpec struct {
	Id       string            `json:"id"`
	Name     string            `json:"name"`
	Unit     string            `json:"unit"`
	Metadata map[string]string `json:"metadata"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ProductList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Product `json:"items"`
}
