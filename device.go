package simulator

import (
	"time"

	"github.com/joinself/zktf-sdk-go/keypair/signing"
	"github.com/joinself/zktf-sim-go/internal/ffi"
)

// ContentType selects a message content type for MatchContentType.
type ContentType uint32

const (
	ContentUnknown           ContentType = ContentType(ffi.ContentUnknown)
	ContentCustom            ContentType = ContentType(ffi.ContentCustom)
	ContentChat              ContentType = ContentType(ffi.ContentChat)
	ContentReceipt           ContentType = ContentType(ffi.ContentReceipt)
	ContentCredential        ContentType = ContentType(ffi.ContentCredential)
	ContentIntroduction      ContentType = ContentType(ffi.ContentIntroduction)
	ContentDiscoveryRequest  ContentType = ContentType(ffi.ContentDiscoveryRequest)
	ContentDiscoveryResponse ContentType = ContentType(ffi.ContentDiscoveryResponse)
	ContentExchangeRequest   ContentType = ContentType(ffi.ContentExchangeRequest)
	ContentExchangeResponse  ContentType = ContentType(ffi.ContentExchangeResponse)
)

// Match selects which incoming messages a rule applies to. Build one with
// MatchAny, MatchContentType, or MatchRequestID.
type Match struct {
	kind        ffi.MatchKind
	contentType ContentType
	requestID   []byte
}

// MatchAny matches every incoming message.
func MatchAny() Match { return Match{kind: ffi.MatchAny} }

// MatchContentType matches messages of a given content type.
func MatchContentType(ct ContentType) Match {
	return Match{kind: ffi.MatchContentType, contentType: ct}
}

// MatchRequestID matches the message carrying a specific request id.
func MatchRequestID(id []byte) Match {
	return Match{kind: ffi.MatchRequestID, requestID: id}
}

// Behaviour is how the device reacts to a matched message. Build one with
// Accept, Reject, or Ignore, optionally deferred with After.
type Behaviour struct {
	action ffi.Behaviour
	delay  time.Duration
}

// Accept drives the matched workflow to completion with simulated user consent.
func Accept() Behaviour { return Behaviour{action: ffi.BehaveAccept} }

// Reject responds to the matched workflow with a rejection.
func Reject() Behaviour { return Behaviour{action: ffi.BehaveReject} }

// Ignore drops the matched message without responding.
func Ignore() Behaviour { return Behaviour{action: ffi.BehaveIgnore} }

// After defers the behaviour by d before it is applied.
func (b Behaviour) After(d time.Duration) Behaviour {
	b.delay = d
	return b
}

// Device is a simulated mobile client. Register rules with Expect before
// driving a workflow; the device auto-responds to matching messages on a
// background thread.
type Device struct {
	h *ffi.Device
}

// NewDevice creates a device with its own account connected to the network.
func NewDevice(network *Network) *Device {
	return &Device{h: ffi.NewDevice(network.h)}
}

// AttachDevice creates a device attached to a real, test-deployed backend at
// the given endpoints. The device always uses test trust anchors, so it can
// only interoperate with test networks — there is no way to target production.
func AttachDevice(rpcEndpoint, objectEndpoint, messagingEndpoint string) *Device {
	return &Device{h: ffi.DeviceAttach(rpcEndpoint, objectEndpoint, messagingEndpoint)}
}

// Expect registers an auto-response rule. Rules are evaluated in registration
// order; the first match wins.
func (d *Device) Expect(m Match, b Behaviour) {
	d.h.Expect(
		m.kind,
		ffi.ContentType(m.contentType),
		m.requestID,
		b.action,
		uint64(b.delay/time.Millisecond),
	)
}

// Address returns the device's zktf address.
func (d *Device) Address() (*signing.PublicKey, error) {
	b, err := d.h.Address()
	if err != nil {
		return nil, err
	}
	return signing.FromBytes(b)
}

// Inbox returns the device's messaging inbox public key.
func (d *Device) Inbox() (*signing.PublicKey, error) {
	b, err := d.h.Inbox()
	if err != nil {
		return nil, err
	}
	return signing.FromBytes(b)
}

// Register drives the registration workflow against counterparty. It blocks
// until the workflow completes, is rejected, or times out.
func (d *Device) Register(counterparty *signing.PublicKey) error {
	return d.h.Register(counterparty.Bytes())
}

// Connect drives the pairwise connect workflow against counterparty. It blocks
// until the workflow completes, is rejected, or times out.
func (d *Device) Connect(counterparty *signing.PublicKey) error {
	return d.h.Connect(counterparty.Bytes())
}
