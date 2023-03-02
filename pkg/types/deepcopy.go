package types

func (a *Authority) DeepCopyObject() Object {
	var authority = Authority{}

	authority.Id = a.Id
	authority.Name = a.Name
	authority.Products = a.Products
	authority.Metadata = a.Metadata

	return &authority
}

func (l *License) DeepCopyObject() Object {
	var license = License{}

	license.Id = l.Id
	license.Licensee = l.Licensee
	license.Metadata = l.Metadata
	license.Grants = l.Grants
	license.NotBefore = l.NotBefore
	license.NotAfter = l.NotAfter
	license.Key = l.Key
	license.Certificate = l.Certificate

	return &license
}

func (l *Licensee) DeepCopyObject() Object {
	var licensee = Licensee{}

	licensee.Id = l.Id
	licensee.Authority = l.Authority
	licensee.Name = l.Name
	licensee.Metadata = l.Metadata

	return &licensee
}

func (c *Certificate) DeepCopyObject() Object {
	var certificate = Certificate{}

	certificate.Id = c.Id
	certificate.PrivateKey = c.PrivateKey
	certificate.Authority = c.Authority
	certificate.Metadata = c.Metadata

	return &certificate
}

func (p *Product) DeepCopyObject() Object {
	var product = Product{}

	product.Id = p.Id
	product.Name = p.Name
	product.Unit = p.Unit
	product.Metadata = p.Metadata

	return &product
}
