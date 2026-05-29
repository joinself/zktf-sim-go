//go:build integration

// Integration tests exercise the native simulator end to end: a simulated
// device registers against a simulated verifier over an in-process network.
//
//	CGO_CFLAGS="-I/path/to/zktf-sdk/crates/sim-ffi -I/path/to/zktf-sdk/crates/zktf-ffi" \
//	CGO_LDFLAGS=-L/path/to/zktf-sdk/target/debug \
//	LD_LIBRARY_PATH=/path/to/zktf-sdk/target/debug \
//	go test -tags integration -v ./...
package simulator_test

import (
	"testing"

	simulator "github.com/joinself/zktf-sim-go"
)

func TestRegister(t *testing.T) {
	network := simulator.NewDefaultNetwork()

	verifier, err := simulator.NewVerifier(network)
	if err != nil {
		t.Fatalf("new verifier: %v", err)
	}

	device := simulator.NewDevice(network)
	device.Expect(simulator.MatchAny(), simulator.Accept())

	identifier, err := verifier.Identifier()
	if err != nil {
		t.Fatalf("verifier identifier: %v", err)
	}

	if err := device.Register(identifier); err != nil {
		t.Fatalf("register: %v", err)
	}
}
