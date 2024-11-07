package hostCore

import (
	"github.com/kalyan3104/k-chain-vm-v1_4-go/vmhost"
)

func (host *vmHost) handleBreakpointIfAny(executionErr error) error {
	if executionErr == nil {
		return nil
	}

	runtime := host.Runtime()
	breakpointValue := runtime.GetRuntimeBreakpointValue()
	log.Trace("handleBreakpointIfAny", "value", breakpointValue)
	if breakpointValue != vmhost.BreakpointNone {
		err := host.handleBreakpoint(breakpointValue)
		runtime.AddError(err, runtime.Function())
		return err
	}

	log.Trace("wasmer execution error", "err", executionErr)
	runtime.AddError(executionErr, runtime.Function())
	return vmhost.ErrExecutionFailed
}

func (host *vmHost) handleBreakpoint(breakpointValue vmhost.BreakpointValue) error {
	if breakpointValue == vmhost.BreakpointAsyncCall {
		return host.handleAsyncCallBreakpoint()
	}
	if breakpointValue == vmhost.BreakpointExecutionFailed {
		return vmhost.ErrExecutionFailed
	}
	if breakpointValue == vmhost.BreakpointSignalError {
		return vmhost.ErrSignalError
	}
	if breakpointValue == vmhost.BreakpointOutOfGas {
		return vmhost.ErrNotEnoughGas
	}
	if breakpointValue == vmhost.BreakpointMemoryLimit {
		return vmhost.ErrMemoryLimit
	}

	return vmhost.ErrUnhandledRuntimeBreakpoint
}
