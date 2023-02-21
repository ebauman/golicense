package types

import (
	"fmt"
	"reflect"
)

type InvalidLicenseError struct{}

func (i InvalidLicenseError) Error() string {
	return "invalid license"
}

func NewInvalidLicenseError() InvalidLicenseError {
	return InvalidLicenseError{}
}

type InvalidGrantError struct {
	grant string
}

func (ige InvalidGrantError) Error() string {
	return fmt.Sprintf("invalid grant %s. valid grant is of the form {domain}/{unit}={amount}, where {domain} is an RFC 1123 "+
		"domain, {unit} is a minimum 1 character length alphanumeric string, and {amount} is a 32-bit integer (supplied as string). "+
		"regex used for this format is %s", ige.grant, GrantRegex)
}

func NewInvalidGrantError(grant string) InvalidGrantError {
	return InvalidGrantError{
		grant: grant,
	}
}

func IsError[T, K error](err T, is K) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(is)
}

func IsErrorReflect(err error, is error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(is)
}

type InvalidMetadataError struct {
	metadata string
}

func (ime InvalidMetadataError) Error() string {
	return fmt.Sprintf("invalid metadata %s. valid metadata is of the form {domain}/{label}={value}, where {domain} is an RFC 1123 "+
		"domain, {label} is a minimum 1 character length alphanumeric string, and {value} is a minimum 1 character length alphanumeric string. "+
		"regex used for this format is %s", ime.metadata, MetadataRegex)
}

func NewInvalidMetadataError(metadata string) InvalidMetadataError {
	return InvalidMetadataError{
		metadata: metadata,
	}
}

type GrantExceededError struct {
	usage       int
	grant       string
	grantAmount int
}

func (iue GrantExceededError) Error() string {
	return fmt.Sprintf("invalid usage, grant %s exceeded (requested %d, granted %d)", iue.grant, iue.usage, iue.grantAmount)
}

func NewGrantExceededError(grant string, usage int, grantAmount int) GrantExceededError {
	return GrantExceededError{
		grant:       grant,
		usage:       usage,
		grantAmount: grantAmount,
	}
}

type GrantNotFoundError struct {
	usage          int
	requestedGrant string
}

func (gnfe GrantNotFoundError) Error() string {
	return fmt.Sprintf("invalid usage, grant %s not found (requested %d)", gnfe.requestedGrant, gnfe.usage)
}

func NewGrantNotFoundError(requestedGrant string, usage int) GrantNotFoundError {
	return GrantNotFoundError{
		usage:          usage,
		requestedGrant: requestedGrant,
	}
}

type InvalidPublicKeysError struct {
}

func (InvalidPublicKeysError) Error() string {
	return fmt.Sprintf("invalid public keys")
}

func NewInvalidPublicKeysError() InvalidPublicKeysError {
	return InvalidPublicKeysError{}
}
