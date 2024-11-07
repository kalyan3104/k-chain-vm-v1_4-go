package mock

import (
	"math/big"

	"github.com/kalyan3104/k-chain-core-go/data/vm"
	vmcommon "github.com/kalyan3104/k-chain-vm-common-go"
	"github.com/kalyan3104/k-chain-vm-v1_4-go/vmhost"
)

var _ vmhost.OutputContext = (*OutputContextStub)(nil)

// OutputContextStub is used in tests to check the OutputContext interface method calls
type OutputContextStub struct {
	InitStateCalled                   func()
	PushStateCalled                   func()
	PopSetActiveStateCalled           func()
	PopMergeActiveStateCalled         func()
	PopDiscardCalled                  func()
	ClearStateStackCalled             func()
	CopyTopOfStackToActiveStateCalled func()
	CensorVMOutputCalled              func()
	GetOutputAccountsCalled           func() map[string]*vmcommon.OutputAccount
	GetOutputAccountCalled            func(address []byte) (*vmcommon.OutputAccount, bool)
	DeleteOutputAccountCalled         func(address []byte)
	WriteLogCalled                    func(address []byte, topics [][]byte, data []byte)
	TransferCalled                    func(destination []byte, sender []byte, gasLimit uint64, gasLocked uint64, value *big.Int, input []byte) error
	TransferDCDTCalled                func(destination []byte, sender []byte, transfers []*vmcommon.DCDTTransfer, input *vmcommon.ContractCallInput) (uint64, error)
	SelfDestructCalled                func(address []byte, beneficiary []byte)
	GetRefundCalled                   func() uint64
	SetRefundCalled                   func(refund uint64)
	ReturnCodeCalled                  func() vmcommon.ReturnCode
	SetReturnCodeCalled               func(returnCode vmcommon.ReturnCode)
	ReturnMessageCalled               func() string
	SetReturnMessageCalled            func(message string)
	ReturnDataCalled                  func() [][]byte
	ClearReturnDataCalled             func()
	RemoveReturnDataCalled            func(index uint32)
	FinishCalled                      func(data []byte)
	PrependFinishCalled               func(data []byte)
	DeleteFirstReturnDataCalled       func()
	GetVMOutputCalled                 func() *vmcommon.VMOutput
	AddTxValueToAccountCalled         func(address []byte, value *big.Int)
	DeployCodeCalled                  func(input vmhost.CodeDeployInput)
	CreateVMOutputInCaseOfErrorCalled func(err error) *vmcommon.VMOutput
	AddToActiveStateCalled            func(vmOutput *vmcommon.VMOutput)
	TransferValueOnlyCalled           func(destination []byte, sender []byte, value *big.Int, checkPayable bool) error
	RemoveNonUpdatedStorageCalled     func()
}

// AddToActiveState mocked method
func (o *OutputContextStub) AddToActiveState(vmOutput *vmcommon.VMOutput) {
	if o.AddToActiveStateCalled != nil {
		o.AddToActiveStateCalled(vmOutput)
	}
}

// InitState mocked method
func (o *OutputContextStub) InitState() {
	if o.InitStateCalled != nil {
		o.InitStateCalled()
	}
}

// PushState mocked method
func (o *OutputContextStub) PushState() {
	if o.PushStateCalled != nil {
		o.PushStateCalled()
	}
}

// PopSetActiveState mocked method
func (o *OutputContextStub) PopSetActiveState() {
	if o.PopSetActiveStateCalled != nil {
		o.PopSetActiveStateCalled()
	}
}

// PopMergeActiveState mocked method
func (o *OutputContextStub) PopMergeActiveState() {
	if o.PopMergeActiveStateCalled != nil {
		o.PopMergeActiveStateCalled()
	}
}

// PopDiscard mocked method
func (o *OutputContextStub) PopDiscard() {
	if o.PopDiscardCalled != nil {
		o.PopDiscardCalled()
	}
}

// ClearStateStack mocked method
func (o *OutputContextStub) ClearStateStack() {
	if o.ClearStateStackCalled != nil {
		o.ClearStateStackCalled()
	}
}

// CopyTopOfStackToActiveState mocked method
func (o *OutputContextStub) CopyTopOfStackToActiveState() {
	if o.CopyTopOfStackToActiveStateCalled != nil {
		o.CopyTopOfStackToActiveStateCalled()
	}
}

// CensorVMOutput mocked method
func (o *OutputContextStub) CensorVMOutput() {
	if o.CensorVMOutputCalled != nil {
		o.CensorVMOutputCalled()
	}
}

// GetOutputAccounts mocked method
func (o *OutputContextStub) GetOutputAccounts() map[string]*vmcommon.OutputAccount {
	if o.GetOutputAccountsCalled != nil {
		return o.GetOutputAccountsCalled()
	}
	return nil
}

// GetOutputAccount mocked method
func (o *OutputContextStub) GetOutputAccount(address []byte) (*vmcommon.OutputAccount, bool) {
	if o.GetOutputAccountCalled != nil {
		return o.GetOutputAccountCalled(address)
	}
	return nil, false
}

// DeleteOutputAccount mocked method
func (o *OutputContextStub) DeleteOutputAccount(address []byte) {
	if o.DeleteOutputAccountCalled != nil {
		o.DeleteOutputAccountCalled(address)
	}
}

// WriteLog mocked method
func (o *OutputContextStub) WriteLog(address []byte, topics [][]byte, data []byte) {
	if o.WriteLogCalled != nil {
		o.WriteLogCalled(address, topics, data)
	}
}

// TransferValueOnly mocked method
func (o *OutputContextStub) TransferValueOnly(destination []byte, sender []byte, value *big.Int, checkPayable bool) error {
	if o.TransferValueOnlyCalled != nil {
		return o.TransferValueOnlyCalled(destination, sender, value, checkPayable)
	}

	return nil
}

// Transfer mocked method
func (o *OutputContextStub) Transfer(destination []byte, sender []byte, gasLimit uint64, gasLocked uint64, value *big.Int, input []byte, _ vm.CallType) error {
	if o.TransferCalled != nil {
		return o.TransferCalled(destination, sender, gasLimit, gasLocked, value, input)
	}

	return nil
}

// TransferDCDT mocked method
func (o *OutputContextStub) TransferDCDT(destination []byte, sender []byte, transfers []*vmcommon.DCDTTransfer, callInput *vmcommon.ContractCallInput) (uint64, error) {
	if o.TransferDCDTCalled != nil {
		return o.TransferDCDTCalled(destination, sender, transfers, callInput)
	}
	return 0, nil
}

// SelfDestruct mocked method
func (o *OutputContextStub) SelfDestruct(address []byte, beneficiary []byte) {
	if o.SelfDestructCalled != nil {
		o.SelfDestructCalled(address, beneficiary)
	}
}

// GetRefund mocked method
func (o *OutputContextStub) GetRefund() uint64 {
	if o.GetRefundCalled != nil {
		return o.GetRefundCalled()
	}
	return 0
}

// SetRefund mocked method
func (o *OutputContextStub) SetRefund(refund uint64) {
	if o.SetRefundCalled != nil {
		o.SetRefundCalled(refund)
	}
}

// ReturnCode mocked method
func (o *OutputContextStub) ReturnCode() vmcommon.ReturnCode {
	if o.ReturnCodeCalled != nil {
		return o.ReturnCodeCalled()
	}
	return vmcommon.Ok
}

// SetReturnCode mocked method
func (o *OutputContextStub) SetReturnCode(returnCode vmcommon.ReturnCode) {
	if o.SetReturnCodeCalled != nil {
		o.SetReturnCodeCalled(returnCode)
	}
}

// ReturnMessage mocked method
func (o *OutputContextStub) ReturnMessage() string {
	if o.ReturnMessageCalled != nil {
		return o.ReturnMessageCalled()
	}
	return ""
}

// SetReturnMessage mocked method
func (o *OutputContextStub) SetReturnMessage(message string) {
	if o.SetReturnMessageCalled != nil {
		o.SetReturnMessageCalled(message)
	}
}

// ReturnData mocked method
func (o *OutputContextStub) ReturnData() [][]byte {
	if o.ReturnDataCalled != nil {
		return o.ReturnDataCalled()
	}
	return [][]byte{}
}

// ClearReturnData mocked method
func (o *OutputContextStub) ClearReturnData() {
	if o.ClearReturnDataCalled != nil {
		o.ClearReturnDataCalled()
	}
}

// RemoveReturnData mocked method
func (o *OutputContextStub) RemoveReturnData(index uint32) {
	if o.RemoveReturnDataCalled != nil {
		o.RemoveReturnDataCalled(index)
	}
}

// Finish mocked method
func (o *OutputContextStub) Finish(data []byte) {
	if o.FinishCalled != nil {
		o.FinishCalled(data)
	}
}

// PrependFinish mocked method
func (o *OutputContextStub) PrependFinish(data []byte) {
	if o.PrependFinishCalled != nil {
		o.PrependFinishCalled(data)
	}
}

// DeleteFirstReturnData mocked method
func (o *OutputContextStub) DeleteFirstReturnData() {
	if o.DeleteFirstReturnDataCalled != nil {
		o.DeleteFirstReturnDataCalled()
	}
}

// GetVMOutput mocked method
func (o *OutputContextStub) GetVMOutput() *vmcommon.VMOutput {
	if o.GetVMOutputCalled != nil {
		return o.GetVMOutputCalled()
	}
	return nil
}

// RemoveNonUpdatedStorage mocked method
func (o *OutputContextStub) RemoveNonUpdatedStorage() {
	if o.RemoveNonUpdatedStorageCalled != nil {
		o.RemoveNonUpdatedStorageCalled()
	}
}

// AddTxValueToAccount mocked method
func (o *OutputContextStub) AddTxValueToAccount(address []byte, value *big.Int) {
	if o.AddTxValueToAccountCalled != nil {
		o.AddTxValueToAccountCalled(address, value)
	}
}

// DeployCode mocked method
func (o *OutputContextStub) DeployCode(input vmhost.CodeDeployInput) {
	if o.DeployCodeCalled != nil {
		o.DeployCodeCalled(input)
	}
}

// CreateVMOutputInCaseOfError mocked method
func (o *OutputContextStub) CreateVMOutputInCaseOfError(err error) *vmcommon.VMOutput {
	if o.CreateVMOutputInCaseOfErrorCalled != nil {
		return o.CreateVMOutputInCaseOfErrorCalled(err)
	}
	return nil
}
