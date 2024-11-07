package testcommon

import (
	"errors"

	vmcommon "github.com/kalyan3104/k-chain-vm-common-go"
)

// MockBuiltin defined the functions that can be replaced in order to mock a builtin
type MockBuiltin struct {
	ProcessBuiltinFunctionCall func(acntSnd, _ vmcommon.UserAccountHandler, vmInput *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error)
	setNewGasConfig            func(_ *vmcommon.GasCost)
	isInterfaceNil             func() bool
}

// ProcessBuiltinFunction - see BuiltinFunction.ProcessBuiltInFunction()
func (m *MockBuiltin) ProcessBuiltinFunction(acntSnd, acntRcv vmcommon.UserAccountHandler, vmInput *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	if m.ProcessBuiltinFunctionCall == nil {
		return nil, errors.New("undefined processBuiltinFunction")
	}
	return m.ProcessBuiltinFunctionCall(acntSnd, acntRcv, vmInput)
}

// SetNewGasConfig - see BuiltinFunction.SetNewGasConfig()
func (m *MockBuiltin) SetNewGasConfig(gasCost *vmcommon.GasCost) {
	if m.setNewGasConfig != nil {
		m.setNewGasConfig(gasCost)
	}
}

// IsActive -
func (m *MockBuiltin) IsActive() bool {
	return true
}

// IsInterfaceNil - see BuiltinFunction.IsInterfaceNil()
func (m *MockBuiltin) IsInterfaceNil() bool {
	if m.isInterfaceNil == nil {
		return m == nil
	}
	return m.isInterfaceNil()
}
