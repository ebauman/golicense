package storage

import (
	"github.com/ebauman/golicense/pkg/types"
)

type GolicenseStore interface {
	AuthorityStore
	LicenseStore
	LicenseeStore
	CertificateStore
	ProductStore
}

type AuthorityStore interface {
	CreateAuthority(authority *types.Authority) (*types.Authority, error)
	ListAuthorities() ([]*types.Authority, error)
	GetAuthority(id string) (*types.Authority, error)
	UpdateAuthority(authority *types.Authority) (*types.Authority, error)
	DeleteAuthority(id string) error
}

type LicenseeStore interface {
	CreateLicensee(licensee *types.Licensee) (*types.Licensee, error)
	ListLicensees(authority string) ([]*types.Licensee, error)
	GetLicensee(id string) (*types.Licensee, error)
	UpdateLicensee(licensee *types.Licensee) (*types.Licensee, error)
	DeleteLicensee(id string) error
}

type CertificateStore interface {
	CreateCertificate(certificate *types.Certificate) (*types.Certificate, error)
	ListCertificates(authority string) ([]*types.Certificate, error)
	GetCertificate(id string) (*types.Certificate, error)
	UpdateCertificate(certificate *types.Certificate) (*types.Certificate, error)
	DeleteCertificate(id string) error
}

type LicenseStore interface {
	CreateLicense(license *types.License) (*types.License, error)
	ListLicensesForLicensee(licensee string) ([]*types.License, error)
	ListLicensesForCertificate(certificate string) ([]*types.License, error)
	GetLicense(id string) (*types.License, error)
	UpdateLicense(license *types.License) (*types.License, error)
	DeleteLicense(id string) error
}

type ProductStore interface {
	CreateProduct(product *types.Product) (*types.Product, error)
	ListProducts() ([]*types.Product, error)
	GetProduct(id string) (*types.Product, error)
	UpdateProduct(product *types.Product) (*types.Product, error)
	DeleteProduct(id string) error
}
