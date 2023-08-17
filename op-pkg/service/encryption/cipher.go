package encryption

import "context"

type Cipher struct{}

func (c *Cipher) Encrypt(ctx context.Context, payload []byte, secret string) ([]byte, error) {
	return payload, nil
}

type Decipher struct{}

func (d *Decipher) Decrypt(ctx context.Context, payload []byte, secret string) ([]byte, error) {
	return payload, nil
}
