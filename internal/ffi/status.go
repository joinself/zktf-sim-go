package ffi

/*
#include <zktf-sim.h>
*/
import "C"

// Status wraps a non-OK zktf_sim_status code and implements the error
// interface. The public package returns it as a plain `error`, so the C enum
// never appears in an exported signature.
type Status struct {
	code uint32
}

// Code returns the raw zktf_sim_status code.
func (s *Status) Code() uint32 {
	return s.code
}

// Error returns a human readable message for the status code.
func (s *Status) Error() string {
	switch C.enum_zktf_sim_status(s.code) {
	case C.ZKTF_SIM_NULL_ARGUMENT:
		return "simulator: null argument"
	case C.ZKTF_SIM_BUFFER_INSUFFICIENT:
		return "simulator: output buffer insufficient"
	case C.ZKTF_SIM_TIMEOUT:
		return "simulator: workflow timed out"
	case C.ZKTF_SIM_POLICY_REJECTED:
		return "simulator: workflow rejected by policy"
	case C.ZKTF_SIM_SDK_ERROR:
		return "simulator: sdk error"
	case C.ZKTF_SIM_INVALID_INPUT:
		return "simulator: invalid input"
	default:
		return "simulator: unknown error"
	}
}

// status converts a raw zktf_sim_status result into a Go error. ZKTF_SIM_OK
// returns nil; any other code returns a *Status.
func status(code C.enum_zktf_sim_status) error {
	if code == C.ZKTF_SIM_OK {
		return nil
	}
	return &Status{code: uint32(code)}
}
