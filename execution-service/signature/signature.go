package signature

import "execution-service/domain"

type Signature struct {
	Value     []byte
	PublicKey domain.PublicKey
}
