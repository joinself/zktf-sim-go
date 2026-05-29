package ffi

/*
#include <zktf-sim.h>
#include <stdlib.h>
*/
import "C"

import "runtime"

// Verifier wraps a zktf_sim_verifier handle.
type Verifier struct {
	ptr *C.zktf_sim_verifier
}

func newVerifier(ptr *C.zktf_sim_verifier) *Verifier {
	if ptr == nil {
		return nil
	}
	v := &Verifier{ptr: ptr}
	runtime.AddCleanup(v, func(ptr *C.zktf_sim_verifier) {
		C.zktf_sim_verifier_destroy(ptr)
	}, v.ptr)
	return v
}

// NewVerifier allocates a simulated server verifier connected to the network.
func NewVerifier(network *Network) (*Verifier, error) {
	var ptr *C.zktf_sim_verifier
	if err := status(C.zktf_sim_verifier_new(network.ptr, &ptr)); err != nil {
		return nil, err
	}
	return newVerifier(ptr), nil
}

func (v *Verifier) Identifier() ([]byte, error) {
	return keyBytes(func(buf *C.uint8_t) C.enum_zktf_sim_status {
		return C.zktf_sim_verifier_identifier(v.ptr, buf, signingKeyBytesLen)
	})
}

func (v *Verifier) Inbox() ([]byte, error) {
	return keyBytes(func(buf *C.uint8_t) C.enum_zktf_sim_status {
		return C.zktf_sim_verifier_inbox(v.ptr, buf, signingKeyBytesLen)
	})
}

func (v *Verifier) Assertion() ([]byte, error) {
	return keyBytes(func(buf *C.uint8_t) C.enum_zktf_sim_status {
		return C.zktf_sim_verifier_assertion(v.ptr, buf, signingKeyBytesLen)
	})
}
