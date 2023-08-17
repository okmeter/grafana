package encryption

// these interfaces are only used in implementation
type Provider struct{}

func (p *Provider) ProvideCiphers() map[string]Cipher {
	return map[string]Cipher{}
}

func (p *Provider) ProvideDeciphers() map[string]Decipher {
	return map[string]Decipher{}
}
