package simulator

import "github.com/joinself/zktf-sim-go/internal/ffi"

// Network is an in-process zktf network. Its endpoints are used to configure a
// real zktf-sdk-go account; its fault-injection methods perturb message
// delivery for chaos testing. The underlying listeners are shut down when the
// Network is garbage collected.
type Network struct {
	h *ffi.Network
}

// NewNetwork starts a simulated network on the given ports.
func NewNetwork(apiPort, objectPort, messagingPort, controlPort uint16) *Network {
	return &Network{h: ffi.NewNetwork(apiPort, objectPort, messagingPort, controlPort)}
}

// NewDefaultNetwork starts a simulated network on the default ports
// (3000/3500/4000/9000).
func NewDefaultNetwork() *Network {
	return &Network{h: ffi.NewDefaultNetwork()}
}

// RPCEndpoint returns the rpc endpoint URL.
func (n *Network) RPCEndpoint() string { return n.h.RPCEndpoint() }

// ObjectEndpoint returns the object storage endpoint URL.
func (n *Network) ObjectEndpoint() string { return n.h.ObjectEndpoint() }

// MessagingEndpoint returns the websocket messaging endpoint URL.
func (n *Network) MessagingEndpoint() string { return n.h.MessagingEndpoint() }

// MessageCount returns the number of messages currently queued for address.
func (n *Network) MessageCount(address []byte) int { return n.h.MessageCount(address) }

// MessageBlock blocks delivery to address and closes its subscription.
func (n *Network) MessageBlock(address []byte) { n.h.MessageBlock(address) }

// MessageUnblock resumes delivery to address.
func (n *Network) MessageUnblock(address []byte) { n.h.MessageUnblock(address) }

// FaultReorder shuffles the queued messages in the range [from, until).
func (n *Network) FaultReorder(address []byte, from, until int) {
	n.h.FaultReorder(address, from, until)
}

// FaultRedeliver duplicates and redelivers the messages in [from, until).
func (n *Network) FaultRedeliver(address []byte, from, until int) {
	n.h.FaultRedeliver(address, from, until)
}

// FaultDelete drops the messages in [from, until), optionally preserving MLS
// commits so group state stays recoverable.
func (n *Network) FaultDelete(address []byte, from, until int, preserveCommits bool) {
	n.h.FaultDelete(address, from, until, preserveCommits)
}
