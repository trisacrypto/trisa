package crypto

import "crypto/rand"

type Handler interface {
	Encrypt()
	Decrypt()
	Sign()
	Verify()
}

func GenRandom(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
