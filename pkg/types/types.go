package types

import (
	"time"
)

const (
	GrantRegex    = `(?P<domain>[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*)\/(?P<unit>[a-z0-9]([-a-z0-9]*[a-z0-9]?))\=(?P<amount>[0-9]*)`
	MetadataRegex = `(?P<domain>[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*)\/(?P<label>[a-z0-9]([-a-z0-9]*[a-z0-9]?))\=(?P<value>[a-z]([a-z]*))`
)

type License struct {
	Id        string            `json:"id"`
	Licensee  string            `json:"licensee"`
	Metadata  map[string]string `json:"metadata"`
	Grants    map[string]int    `json:"grants"`
	NotBefore time.Time         `json:"notBefore"`
	NotAfter  time.Time         `json:"notAfter"`
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
