// Package ffi is the single cgo boundary for the zktf Go simulator.
//
// It is the ONLY package in this module that does `import "C"`. The public
// `simulator` package is pure Go and talks to the native simulator library
// exclusively through the Go types and functions exported here. Keeping all cgo
// in one package means C types never escape: each wrapper struct holds its
// `*C.zktf_sim_*` pointer in an UNEXPORTED field, so no exported signature ever
// mentions a C type.
//
// The simulator library is deliberately standalone: it shares no symbols with
// libzktf_sdk, so a test binary may link both. Values cross the boundary as
// raw bytes (33-byte signing keys / addresses); the public package rehydrates
// them into zktf-sdk-go native key types.
//
// Build prerequisites: the native header `zktf-sim.h` must be on the C include
// path and `libzktf_sim` on the linker path. Because this module also imports
// zktf-sdk-go, `zktf-sdk.h` / `libzktf_sdk` must be available too. For local
// development point cgo at the zktf-sdk checkout, e.g.:
//
//	CGO_CFLAGS="-I/path/to/zktf-sdk/crates/sim-ffi -I/path/to/zktf-sdk/crates/zktf-ffi" \
//	CGO_LDFLAGS=-L/path/to/zktf-sdk/target/debug \
//	LD_LIBRARY_PATH=/path/to/zktf-sdk/target/debug \
//	go build ./...
package ffi

/*
#cgo LDFLAGS: -lstdc++ -lm -ldl
#cgo linux LDFLAGS: -lzktf_sim
#cgo darwin LDFLAGS: -lzktf_sim -framework CoreFoundation -framework SystemConfiguration -framework Security
#include <zktf-sim.h>
#include <stdlib.h>
*/
import "C"

import "unsafe"

// signingKeyBytesLen is the length of a raw signing public key / address.
const signingKeyBytesLen = 33

// cbytes copies a Go byte slice into C-allocated memory. The returned pointer
// must be released with free. A nil/empty slice yields a nil pointer.
func cbytes(b []byte) (*C.uint8_t, C.size_t) {
	if len(b) == 0 {
		return nil, 0
	}
	return (*C.uint8_t)(C.CBytes(b)), C.size_t(len(b))
}

// free releases memory allocated by cbytes.
func free(p unsafe.Pointer) {
	C.free(p)
}
