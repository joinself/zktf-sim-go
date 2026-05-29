package ffi

/*
#include <zktf-sim.h>
#include <stdlib.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// MatchKind mirrors zktf_sim_match_kind.
type MatchKind uint32

const (
	MatchAny         MatchKind = C.ZKTF_SIM_MATCH_ANY
	MatchContentType MatchKind = C.ZKTF_SIM_MATCH_CONTENT_TYPE
	MatchRequestID   MatchKind = C.ZKTF_SIM_MATCH_REQUEST_ID
)

// ContentType mirrors zktf_sim_content_type.
type ContentType uint32

const (
	ContentUnknown           ContentType = C.ZKTF_SIM_CONTENT_UNKNOWN
	ContentCustom            ContentType = C.ZKTF_SIM_CONTENT_CUSTOM
	ContentChat              ContentType = C.ZKTF_SIM_CONTENT_CHAT
	ContentReceipt           ContentType = C.ZKTF_SIM_CONTENT_RECEIPT
	ContentCredential        ContentType = C.ZKTF_SIM_CONTENT_CREDENTIAL
	ContentIntroduction      ContentType = C.ZKTF_SIM_CONTENT_INTRODUCTION
	ContentDiscoveryRequest  ContentType = C.ZKTF_SIM_CONTENT_DISCOVERY_REQUEST
	ContentDiscoveryResponse ContentType = C.ZKTF_SIM_CONTENT_DISCOVERY_RESPONSE
	ContentExchangeRequest   ContentType = C.ZKTF_SIM_CONTENT_EXCHANGE_REQUEST
	ContentExchangeResponse  ContentType = C.ZKTF_SIM_CONTENT_EXCHANGE_RESPONSE
)

// Behaviour mirrors zktf_sim_behaviour.
type Behaviour uint32

const (
	BehaveAccept Behaviour = C.ZKTF_SIM_BEHAVE_ACCEPT
	BehaveReject Behaviour = C.ZKTF_SIM_BEHAVE_REJECT
	BehaveIgnore Behaviour = C.ZKTF_SIM_BEHAVE_IGNORE
)

// Device wraps a zktf_sim_device handle.
type Device struct {
	ptr *C.zktf_sim_device
}

func newDevice(ptr *C.zktf_sim_device) *Device {
	if ptr == nil {
		return nil
	}
	d := &Device{ptr: ptr}
	runtime.AddCleanup(d, func(ptr *C.zktf_sim_device) {
		C.zktf_sim_device_destroy(ptr)
	}, d.ptr)
	return d
}

// NewDevice allocates a simulated mobile device connected to the network.
func NewDevice(network *Network) *Device {
	return newDevice(C.zktf_sim_device_new(network.ptr))
}

// DeviceAttach allocates a simulated device attached to a real, test-deployed
// backend at the given endpoints. The device always uses test trust anchors.
// Returns nil on invalid (non-UTF-8) endpoints.
func DeviceAttach(rpcEndpoint, objectEndpoint, messagingEndpoint string) *Device {
	rpc := C.CString(rpcEndpoint)
	object := C.CString(objectEndpoint)
	messaging := C.CString(messagingEndpoint)
	defer free(unsafe.Pointer(rpc))
	defer free(unsafe.Pointer(object))
	defer free(unsafe.Pointer(messaging))
	return newDevice(C.zktf_sim_device_attach(rpc, object, messaging))
}

// Expect registers an auto-response rule. requestID is used only for
// MatchRequestID; contentType only for MatchContentType; a non-zero delayMs
// wraps the behaviour in a delay.
func (d *Device) Expect(kind MatchKind, contentType ContentType, requestID []byte, behaviour Behaviour, delayMs uint64) {
	buf, length := cbytes(requestID)
	defer free(unsafe.Pointer(buf))
	C.zktf_sim_device_expect(
		d.ptr,
		C.enum_zktf_sim_match_kind(kind),
		C.enum_zktf_sim_content_type(contentType),
		buf, length,
		C.enum_zktf_sim_behaviour(behaviour),
		C.uint64_t(delayMs),
	)
}

func (d *Device) Address() ([]byte, error) {
	return keyBytes(func(buf *C.uint8_t) C.enum_zktf_sim_status {
		return C.zktf_sim_device_address(d.ptr, buf, signingKeyBytesLen)
	})
}

func (d *Device) Inbox() ([]byte, error) {
	return keyBytes(func(buf *C.uint8_t) C.enum_zktf_sim_status {
		return C.zktf_sim_device_inbox(d.ptr, buf, signingKeyBytesLen)
	})
}

// Register drives the registration workflow against counterparty (33-byte
// address). Blocks until the workflow completes, fails, or times out.
func (d *Device) Register(counterparty []byte) error {
	buf, length := cbytes(counterparty)
	defer free(unsafe.Pointer(buf))
	return status(C.zktf_sim_device_register(d.ptr, buf, length))
}

// Connect drives the pairwise connect workflow against counterparty (33-byte
// public key). Blocks until the workflow completes, fails, or times out.
func (d *Device) Connect(counterparty []byte) error {
	buf, length := cbytes(counterparty)
	defer free(unsafe.Pointer(buf))
	return status(C.zktf_sim_device_connect(d.ptr, buf, length))
}

// keyBytes runs a getter that writes a 33-byte signing key into a caller buffer.
func keyBytes(getter func(*C.uint8_t) C.enum_zktf_sim_status) ([]byte, error) {
	buf := C.malloc(signingKeyBytesLen)
	defer C.free(buf)
	if err := status(getter((*C.uint8_t)(buf))); err != nil {
		return nil, err
	}
	return C.GoBytes(buf, signingKeyBytesLen), nil
}
