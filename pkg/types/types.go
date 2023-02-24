package types

import (
	"crypto/rsa"
	"time"
)

const (
	GrantRegex    = `(?P<domain>[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*)\/(?P<unit>[a-z0-9]([-a-z0-9]*[a-z0-9]?))\=(?P<amount>[0-9]*)`
	MetadataRegex = `(?P<domain>[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*)\/(?P<label>[a-z0-9]([-a-z0-9]*[a-z0-9]?))\=(?P<value>[a-z]([a-z]*))`
)

type License struct {
	Id          string            `json:"id"`
	Licensee    string            `json:"licensee"`
	Metadata    map[string]string `json:"metadata"`
	Grants      map[string]int    `json:"grants"`
	NotBefore   time.Time         `json:"notBefore"`
	NotAfter    time.Time         `json:"notAfter"`
	Key         string            `json:"-"`
	Certificate string            `json:"-"`
}

func (l *License) GetKind() string {
	return "license"
}

func (l *License) GetId() string {
	return l.Id
}

func (l *License) SetMetadata(metadata map[string]string) {
	l.Metadata = metadata
}

func (l *License) GetMetadata() map[string]string {
	return l.Metadata
}

type LicensingResponse struct {
	Success bool
	Error   error
}

func NewLicensingResponse(success bool, err error) LicensingResponse {
	return LicensingResponse{
		success,
		err,
	}
}

type Object interface {
	GetId() string
	GetKind() string
}

type MetaObject interface {
	Object
	SetMetadata(map[string]string)
	GetMetadata() map[string]string
}

type Authority struct {
	Id       string
	Name     string
	Products []string
	Metadata map[string]string
}

func (a *Authority) GetKind() string {
	return "authority"
}

func (a *Authority) GetId() string {
	return a.Id
}

func (a *Authority) SetMetadata(metadata map[string]string) {
	a.Metadata = metadata
}

func (a *Authority) GetMetadata() map[string]string {
	return a.Metadata
}

type Licensee struct {
	Id        string
	Name      string
	Authority string
	Metadata  map[string]string
}

func (l *Licensee) GetKind() string {
	return "licensee"
}

func (l *Licensee) GetId() string {
	return l.Id
}

func (l *Licensee) SetMetadata(metadata map[string]string) {
	l.Metadata = metadata
}

func (l *Licensee) GetMetadata() map[string]string {
	return l.Metadata
}

type Certificate struct {
	Id         string
	Authority  string
	PrivateKey rsa.PrivateKey
	Metadata   map[string]string
}

func (c *Certificate) GetKind() string {
	return "certificate"
}

func (c *Certificate) GetId() string {
	return c.Id
}

func (c *Certificate) SetMetadata(metadata map[string]string) {
	c.Metadata = metadata
}

func (c *Certificate) GetMetadata() map[string]string {
	return c.Metadata
}

type Product struct {
	Id       string
	Name     string
	Unit     string
	Metadata map[string]string
}

func (p *Product) GetKind() string {
	return "product"
}

func (p *Product) GetId() string {
	return p.Id
}

func (p *Product) SetMetadata(metadata map[string]string) {
	p.Metadata = metadata
}

func (p *Product) GetMetadata() map[string]string {
	return p.Metadata
}
