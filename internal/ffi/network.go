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

// endpointBufLen is generous for any "ws://127.0.0.1:65535/" style endpoint.
const endpointBufLen = 256

// Network wraps a zktf_sim_network handle.
type Network struct {
	ptr *C.zktf_sim_network
}

func newNetwork(ptr *C.zktf_sim_network) *Network {
	if ptr == nil {
		return nil
	}
	n := &Network{ptr: ptr}
	runtime.AddCleanup(n, func(ptr *C.zktf_sim_network) {
		C.zktf_sim_network_destroy(ptr)
	}, n.ptr)
	return n
}

// NewNetwork allocates a simulated network on the given ports.
func NewNetwork(apiPort, objectPort, messagingPort, controlPort uint16) *Network {
	return newNetwork(C.zktf_sim_network_new(
		C.uint16_t(apiPort),
		C.uint16_t(objectPort),
		C.uint16_t(messagingPort),
		C.uint16_t(controlPort),
	))
}

// NewDefaultNetwork allocates a simulated network on the default ports.
func NewDefaultNetwork() *Network {
	return newNetwork(C.zktf_sim_network_default())
}

func (n *Network) RPCEndpoint() string {
	buf := C.malloc(endpointBufLen)
	defer C.free(buf)
	if status(C.zktf_sim_network_rpc_endpoint(n.ptr, (*C.uint8_t)(buf), endpointBufLen)) != nil {
		return ""
	}
	return C.GoString((*C.char)(buf))
}

func (n *Network) ObjectEndpoint() string {
	buf := C.malloc(endpointBufLen)
	defer C.free(buf)
	if status(C.zktf_sim_network_object_endpoint(n.ptr, (*C.uint8_t)(buf), endpointBufLen)) != nil {
		return ""
	}
	return C.GoString((*C.char)(buf))
}

func (n *Network) MessagingEndpoint() string {
	buf := C.malloc(endpointBufLen)
	defer C.free(buf)
	if status(C.zktf_sim_network_messaging_endpoint(n.ptr, (*C.uint8_t)(buf), endpointBufLen)) != nil {
		return ""
	}
	return C.GoString((*C.char)(buf))
}

func (n *Network) MessageCount(address []byte) int {
	buf, length := cbytes(address)
	defer free(unsafe.Pointer(buf))
	return int(C.zktf_sim_network_message_count(n.ptr, buf, length))
}

func (n *Network) MessageBlock(address []byte) {
	buf, length := cbytes(address)
	defer free(unsafe.Pointer(buf))
	C.zktf_sim_network_message_block(n.ptr, buf, length)
}

func (n *Network) MessageUnblock(address []byte) {
	buf, length := cbytes(address)
	defer free(unsafe.Pointer(buf))
	C.zktf_sim_network_message_unblock(n.ptr, buf, length)
}

func (n *Network) FaultReorder(address []byte, from, until int) {
	buf, length := cbytes(address)
	defer free(unsafe.Pointer(buf))
	C.zktf_sim_network_fault_reorder(n.ptr, buf, length, C.size_t(from), C.size_t(until))
}

func (n *Network) FaultRedeliver(address []byte, from, until int) {
	buf, length := cbytes(address)
	defer free(unsafe.Pointer(buf))
	C.zktf_sim_network_fault_redeliver(n.ptr, buf, length, C.size_t(from), C.size_t(until))
}

func (n *Network) FaultDelete(address []byte, from, until int, preserveCommits bool) {
	buf, length := cbytes(address)
	defer free(unsafe.Pointer(buf))
	C.zktf_sim_network_fault_delete(
		n.ptr, buf, length, C.size_t(from), C.size_t(until), C.bool(preserveCommits),
	)
}
