package simulator

import (
	"github.com/joinself/zktf-sdk-go/keypair/signing"
	"github.com/joinself/zktf-sim-go/internal/ffi"
)

// Verifier is a simulated server verifier. It auto-responds to credential
// verification requests on a background thread, issuing credentials signed by
// its assertion key.
type Verifier struct {
	h *ffi.Verifier
}

// NewVerifier creates a verifier with its issuer identity registered on the
// network.
func NewVerifier(network *Network) (*Verifier, error) {
	h, err := ffi.NewVerifier(network.h)
	if err != nil {
		return nil, err
	}
	return &Verifier{h: h}, nil
}

// Identifier returns the verifier's issuer address.
func (v *Verifier) Identifier() (*signing.PublicKey, error) {
	b, err := v.h.Identifier()
	if err != nil {
		return nil, err
	}
	return signing.FromBytes(b)
}

// Inbox returns the verifier's messaging inbox public key.
func (v *Verifier) Inbox() (*signing.PublicKey, error) {
	b, err := v.h.Inbox()
	if err != nil {
		return nil, err
	}
	return signing.FromBytes(b)
}

// Assertion returns the key the verifier signs issued credentials with.
func (v *Verifier) Assertion() (*signing.PublicKey, error) {
	b, err := v.h.Assertion()
	if err != nil {
		return nil, err
	}
	return signing.FromBytes(b)
}
