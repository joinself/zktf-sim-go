// Package simulator drives the zktf network simulator from Go. It wraps the
// standalone native simulator library (libzktf_sim) and exposes three handles:
//
//   - Network   — an in-process zktf network (rpc/object/messaging endpoints)
//     with message fault-injection for chaos testing.
//   - Device    — a simulated mobile client that auto-responds to workflows.
//   - Verifier  — a simulated server verifier that issues credentials.
//
// It is intended as a TEST utility: pair it with zktf-sdk-go in your test scope
// to exercise a real account against a simulated counterparty. Identity values
// cross into Go as zktf-sdk-go *signing.PublicKey, so the simulator and the SDK
// interoperate without sharing native symbols.
package simulator
