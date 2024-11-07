package contexts

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/kalyan3104/k-chain-core-go/core"
	"github.com/kalyan3104/k-chain-core-go/data/vm"
	logger "github.com/kalyan3104/k-chain-logger-go"
	vmcommon "github.com/kalyan3104/k-chain-vm-common-go"
	"github.com/kalyan3104/k-chain-vm-v1_4-go/vmhost"
)

var _ vmhost.OutputContext = (*outputContext)(nil)

var logOutput = logger.GetOrCreate("vm/output")

type outputContext struct {
	host        vmhost.VMHost
	outputState *vmcommon.VMOutput
	stateStack  []*vmcommon.VMOutput
	codeUpdates map[string]struct{}
}

// NewOutputContext creates a new outputContext
func NewOutputContext(host vmhost.VMHost) (*outputContext, error) {
	context := &outputContext{
		host:       host,
		stateStack: make([]*vmcommon.VMOutput, 0),
	}

	context.InitState()

	return context, nil
}

// InitState initializes the output state and the code updates.
func (context *outputContext) InitState() {
	context.outputState = newVMOutput()
	context.codeUpdates = make(map[string]struct{})
}

func newVMOutput() *vmcommon.VMOutput {
	return &vmcommon.VMOutput{
		ReturnData:      make([][]byte, 0),
		ReturnCode:      vmcommon.Ok,
		ReturnMessage:   "",
		GasRemaining:    0,
		GasRefund:       big.NewInt(0),
		OutputAccounts:  make(map[string]*vmcommon.OutputAccount),
		DeletedAccounts: make([][]byte, 0),
		TouchedAccounts: make([][]byte, 0),
		Logs:            make([]*vmcommon.LogEntry, 0),
	}
}

// NewVMOutputAccount creates a new output account and sets the given address
func NewVMOutputAccount(address []byte) *vmcommon.OutputAccount {
	return &vmcommon.OutputAccount{
		Address:                 address,
		Nonce:                   0,
		BalanceDelta:            big.NewInt(0),
		Balance:                 nil,
		StorageUpdates:          make(map[string]*vmcommon.StorageUpdate),
		BytesAddedToStorage:     0,
		BytesDeletedFromStorage: 0,
	}
}

// PushState appends the current vmOutput to the state stack
func (context *outputContext) PushState() {
	newState := newVMOutput()
	mergeVMOutputs(newState, context.outputState)
	context.stateStack = append(context.stateStack, newState)
}

// PopSetActiveState removes the latest entry from the state stack and sets it as the current vm output
func (context *outputContext) PopSetActiveState() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	prevState := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]
	context.outputState = prevState
}

// PopMergeActiveState merges the current state into the head of the stateStack,
// then pop the head of the stateStack into the current state.
// Doing this allows the VM to execute a SmartContract into a context on top
// of an existing context (a previous SC) without allowing access to it, but
// later merging the output of the two SCs in chronological order.
func (context *outputContext) PopMergeActiveState() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	prevState := context.stateStack[stateStackLen-1]
	context.stateStack = context.stateStack[:stateStackLen-1]

	mergeVMOutputs(prevState, context.outputState)
	context.outputState = newVMOutput()
	mergeVMOutputs(context.outputState, prevState)
}

// PopDiscard removes the latest entry from the state stack, but maintaining
// all GasUsed values.
func (context *outputContext) PopDiscard() {
	stateStackLen := len(context.stateStack)
	if stateStackLen == 0 {
		return
	}

	context.stateStack = context.stateStack[:stateStackLen-1]
}

// ClearStateStack reinitializes the state stack.
func (context *outputContext) ClearStateStack() {
	context.stateStack = make([]*vmcommon.VMOutput, 0)
}

// CensorVMOutput will cause the next executed SC to appear isolated, as if
// nothing was executed before. Required for ExecuteOnDestContext().
// StorageUpdates are not deleted from context.outputState.OutputAccounts,
// preserving the storage cache.
func (context *outputContext) CensorVMOutput() {
	context.outputState.ReturnData = make([][]byte, 0)
	context.outputState.ReturnCode = vmcommon.Ok
	context.outputState.ReturnMessage = ""
	context.outputState.GasRemaining = 0
	context.outputState.GasRefund = big.NewInt(0)
	context.outputState.Logs = make([]*vmcommon.LogEntry, 0)

	logOutput.Trace("state content censored")
}

// GetOutputAccount returns the output account present at the given address,
// and a bool that is true if the account is new. If no output account is present at that address,
// a new account will be created and added to the output accounts.
func (context *outputContext) GetOutputAccount(address []byte) (*vmcommon.OutputAccount, bool) {
	accountIsNew := false
	account, ok := context.outputState.OutputAccounts[string(address)]
	if !ok {
		account = NewVMOutputAccount(address)
		context.outputState.OutputAccounts[string(address)] = account
		accountIsNew = true
	}

	return account, accountIsNew
}

// GetOutputAccounts returns all the OutputAccounts in the current outputState.
func (context *outputContext) GetOutputAccounts() map[string]*vmcommon.OutputAccount {
	return context.outputState.OutputAccounts
}

// DeleteOutputAccount removes the given address from the output accounts and code updates
func (context *outputContext) DeleteOutputAccount(address []byte) {
	delete(context.outputState.OutputAccounts, string(address))
	delete(context.codeUpdates, string(address))
}

// GetRefund returns the value of the gas refund for the current output state.
func (context *outputContext) GetRefund() uint64 {
	return uint64(context.outputState.GasRefund.Int64())
}

// SetRefund sets the given value as gas refund for the current output state.
func (context *outputContext) SetRefund(refund uint64) {
	context.outputState.GasRefund = big.NewInt(int64(refund))
}

// ReturnData returns the data of the current output state.
func (context *outputContext) ReturnData() [][]byte {
	return context.outputState.ReturnData
}

// ReturnCode returns the code of the current output state
func (context *outputContext) ReturnCode() vmcommon.ReturnCode {
	return context.outputState.ReturnCode
}

// SetReturnCode sets the given return code as the return code for the current output state.
func (context *outputContext) SetReturnCode(returnCode vmcommon.ReturnCode) {
	context.outputState.ReturnCode = returnCode
}

// ReturnMessage returns a string that represents the return message for the current output state.
func (context *outputContext) ReturnMessage() string {
	return context.outputState.ReturnMessage
}

// SetReturnMessage sets the given string as a return message for the current output state.
func (context *outputContext) SetReturnMessage(returnMessage string) {
	context.outputState.ReturnMessage = returnMessage
}

// ClearReturnData reinitializes the return data for the current output state.
func (context *outputContext) ClearReturnData() {
	context.outputState.ReturnData = make([][]byte, 0)
}

// RemoveReturnData removes the return data item located at the specified index
func (context *outputContext) RemoveReturnData(index uint32) {
	returnData := context.outputState.ReturnData
	if index >= uint32(len(returnData)) {
		return
	}
	context.outputState.ReturnData = append(returnData[:index], returnData[index+1:]...)
}

// SelfDestruct does nothing
// TODO change comment when the function is implemented
func (context *outputContext) SelfDestruct(_ []byte, _ []byte) {
}

// Finish appends the given data to the return data of the current output state.
func (context *outputContext) Finish(data []byte) {
	context.outputState.ReturnData = append(context.outputState.ReturnData, data)
}

// PrependFinish appends the given data to the return data of the current output state.
func (context *outputContext) PrependFinish(data []byte) {
	context.outputState.ReturnData = append([][]byte{data}, context.outputState.ReturnData...)
}

// DeleteFirstReturnData deletes the first return data, to be used after prepend
func (context *outputContext) DeleteFirstReturnData() {
	if len(context.outputState.ReturnData) > 0 {
		context.outputState.ReturnData = context.outputState.ReturnData[1:]
	}
}

// WriteLogWithIdentifier creates a new LogEntry and appends it to the logs of the current output state.
func (context *outputContext) WriteLogWithIdentifier(address []byte, topics [][]byte, data []byte, identifier []byte) {
	if context.host.Runtime().ReadOnly() {
		logOutput.Trace("log entry", "error", "cannot write logs in readonly mode")
		return
	}

	newLogEntry := &vmcommon.LogEntry{
		Address:    address,
		Data:       [][]byte{data},
		Identifier: identifier,
	}
	logOutput.Trace("log entry", "address", address, "data", data)

	if len(topics) == 0 {
		context.outputState.Logs = append(context.outputState.Logs, newLogEntry)
		return
	}

	newLogEntry.Topics = topics

	context.outputState.Logs = append(context.outputState.Logs, newLogEntry)
	logOutput.Trace("log entry", "endpoint", newLogEntry.Identifier, "topics", newLogEntry.Topics)
}

// WriteLog creates a new LogEntry and appends it to the logs of the current output state.
func (context *outputContext) WriteLog(address []byte, topics [][]byte, data []byte) {
	context.WriteLogWithIdentifier(address, topics, data, []byte(context.host.Runtime().Function()))
}

// TransferValueOnly will transfer the big.int value and checks if it is possible
func (context *outputContext) TransferValueOnly(destination []byte, sender []byte, value *big.Int, checkPayable bool) error {
	logOutput.Trace("transfer value", "sender", sender, "dest", destination, "value", value)

	if value.Cmp(vmhost.Zero) < 0 {
		logOutput.Trace("transfer value", "error", vmhost.ErrTransferNegativeValue)
		return vmhost.ErrTransferNegativeValue
	}

	if !context.hasSufficientBalance(sender, value) {
		logOutput.Trace("transfer value", "error", vmhost.ErrTransferInsufficientFunds)
		return vmhost.ErrTransferInsufficientFunds
	}

	payable, err := context.host.Blockchain().IsPayable(sender, destination)
	if err != nil {
		logOutput.Trace("transfer value", "error", err)
		return err
	}

	isAsyncCall := context.host.Runtime().GetVMInput().CallType == vm.AsynchronousCall
	hasValue := value.Cmp(vmhost.Zero) > 0
	if checkPayable && !payable && hasValue && !isAsyncCall {
		logOutput.Trace("transfer value", "error", vmhost.ErrAccountNotPayable)
		return vmhost.ErrAccountNotPayable
	}

	senderAcc, _ := context.GetOutputAccount(sender)
	destAcc, _ := context.GetOutputAccount(destination)

	senderAcc.BalanceDelta = big.NewInt(0).Sub(senderAcc.BalanceDelta, value)
	destAcc.BalanceDelta = big.NewInt(0).Add(destAcc.BalanceDelta, value)

	if value.Cmp(vmhost.Zero) > 0 {
		if context.host.Runtime().ReadOnly() && context.host.CheckExecuteReadOnly() {
			return vmhost.ErrInvalidCallOnReadOnlyMode
		}

		context.WriteLogWithIdentifier(
			context.host.Runtime().GetContextAddress(),
			[][]byte{sender, destination, value.Bytes()},
			[]byte{},
			[]byte("transferValueOnly"),
		)
	}

	return nil
}

// Transfer handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
func (context *outputContext) Transfer(destination []byte, sender []byte, gasLimit uint64, gasLocked uint64, value *big.Int, input []byte, callType vm.CallType) error {
	checkPayableIfNotCallback := gasLimit > 0 && callType != vm.AsynchronousCallBack
	err := context.TransferValueOnly(destination, sender, value, checkPayableIfNotCallback)
	if err != nil {
		return err
	}

	destAcc, _ := context.GetOutputAccount(destination)
	outputTransfer := vmcommon.OutputTransfer{
		Value:         big.NewInt(0).Set(value),
		GasLimit:      gasLimit,
		GasLocked:     gasLocked,
		Data:          input,
		CallType:      callType,
		SenderAddress: sender,
	}
	destAcc.OutputTransfers = append(destAcc.OutputTransfers, outputTransfer)

	logOutput.Trace("transfer value added")
	return nil
}

// TransferDCDT makes the dcdt/nft transfer and exports the data if it is cross shard
func (context *outputContext) TransferDCDT(
	destination []byte,
	sender []byte,
	transfers []*vmcommon.DCDTTransfer,
	callInput *vmcommon.ContractCallInput,
) (uint64, error) {
	if len(transfers) == 0 {
		return 0, vmhost.ErrTransferValueOnDCDTCall
	}

	isSmartContract := context.host.Blockchain().IsSmartContract(destination)
	sameShard := context.host.AreInSameShard(sender, destination)
	callType := vm.DirectCall
	isExecution := isSmartContract && callInput != nil
	if isExecution {
		callType = vm.DCDTTransferAndExecute
	}

	vmOutput, gasConsumedByTransfer, err := context.host.ExecuteDCDTTransfer(destination, sender, transfers, callType)
	if err != nil {
		return 0, err
	}

	gasRemaining := uint64(0)

	if callInput != nil && isSmartContract {
		if gasConsumedByTransfer > callInput.GasProvided {
			logOutput.Trace("DCDT post-transfer execution", "error", vmhost.ErrNotEnoughGas)
			return 0, vmhost.ErrNotEnoughGas
		}
		gasRemaining = callInput.GasProvided - gasConsumedByTransfer
	}

	if isExecution {
		if gasRemaining > context.host.Metering().GasLeft() {
			logOutput.Trace("DCDT post-transfer execution", "error", vmhost.ErrNotEnoughGas)
			return 0, vmhost.ErrNotEnoughGas
		}

		if !sameShard {
			context.host.Metering().UseGas(gasRemaining)
		}
	}

	destAcc, _ := context.GetOutputAccount(destination)
	outputTransfer := vmcommon.OutputTransfer{
		Value:         big.NewInt(0),
		GasLimit:      gasRemaining,
		GasLocked:     0,
		Data:          []byte{},
		CallType:      vm.DirectCall,
		SenderAddress: sender,
	}

	outputTransfer.Data = context.getOutputTransferDataFromDCDTTransfer(transfers, vmOutput, sameShard, destination)

	if sameShard {
		outputTransfer.GasLimit = 0
	}

	if callInput != nil {
		scCallData := "@" + hex.EncodeToString([]byte(callInput.Function))
		for _, arg := range callInput.Arguments {
			scCallData += "@" + hex.EncodeToString(arg)
		}
		outputTransfer.Data = append(outputTransfer.Data, []byte(scCallData)...)
	}

	destAcc.OutputTransfers = append(destAcc.OutputTransfers, outputTransfer)

	context.outputState.Logs = append(context.outputState.Logs, vmOutput.Logs...)
	return gasRemaining, nil
}

func (context *outputContext) getOutputTransferDataFromDCDTTransfer(
	transfers []*vmcommon.DCDTTransfer,
	vmOutput *vmcommon.VMOutput,
	sameShard bool,
	destination []byte,
) []byte {

	if len(transfers) == 1 && transfers[0].DCDTTokenNonce == 0 {
		return []byte(core.BuiltInFunctionDCDTTransfer + "@" + hex.EncodeToString(transfers[0].DCDTTokenName) + "@" + hex.EncodeToString(transfers[0].DCDTValue.Bytes()))
	}

	if !sameShard {
		outTransfer, ok := vmOutput.OutputAccounts[string(destination)]
		if ok && len(outTransfer.OutputTransfers) == 1 {
			return outTransfer.OutputTransfers[0].Data
		}
	}

	if len(transfers) == 1 {
		data := []byte(core.BuiltInFunctionDCDTNFTTransfer + "@" +
			hex.EncodeToString(transfers[0].DCDTTokenName) + "@" +
			hex.EncodeToString(big.NewInt(0).SetUint64(transfers[0].DCDTTokenNonce).Bytes()) + "@" +
			hex.EncodeToString(transfers[0].DCDTValue.Bytes()) + "@" +
			hex.EncodeToString(destination))
		return data
	}

	data := core.BuiltInFunctionMultiDCDTNFTTransfer + "@" + hex.EncodeToString(destination) + "@" + hex.EncodeToString(big.NewInt(int64(len(transfers))).Bytes())
	for _, transfer := range transfers {
		data += "@" + hex.EncodeToString(transfer.DCDTTokenName) + "@" + hex.EncodeToString(big.NewInt(0).SetUint64(transfer.DCDTTokenNonce).Bytes()) + "@" + hex.EncodeToString(transfer.DCDTValue.Bytes())
	}

	return []byte(data)
}

func (context *outputContext) hasSufficientBalance(address []byte, value *big.Int) bool {
	senderBalance := context.host.Blockchain().GetBalanceBigInt(address)
	return senderBalance.Cmp(value) >= 0
}

// AddTxValueToAccount adds the given value to the BalanceDelta of the account that is mapped to the given address
func (context *outputContext) AddTxValueToAccount(address []byte, value *big.Int) {
	destAcc, _ := context.GetOutputAccount(address)
	destAcc.BalanceDelta = big.NewInt(0).Add(destAcc.BalanceDelta, value)
}

// RemoveNonUpdatedStorage removes non updated storage from output state
func (context *outputContext) RemoveNonUpdatedStorage() {
	for _, outAcc := range context.outputState.OutputAccounts {
		for _, storageUpdate := range outAcc.StorageUpdates {
			if !storageUpdate.Written {
				delete(outAcc.StorageUpdates, string(storageUpdate.Offset))
			}
		}
	}
}

// GetVMOutput updates the current VMOutput and returns it
func (context *outputContext) GetVMOutput() *vmcommon.VMOutput {
	context.removeNonUpdatedCode()

	metering := context.host.Metering()
	context.outputState.GasRemaining = metering.GasLeft()

	err := metering.UpdateGasStateOnSuccess(context.outputState)
	if err != nil {
		return context.CreateVMOutputInCaseOfError(err)
	}

	return context.outputState
}

// DeployCode sets the given code to a an account, and creates a new codeUpdates entry at the accounts address.
func (context *outputContext) DeployCode(input vmhost.CodeDeployInput) {
	newSCAccount, _ := context.GetOutputAccount(input.ContractAddress)
	newSCAccount.Code = input.ContractCode
	newSCAccount.CodeMetadata = input.ContractCodeMetadata
	newSCAccount.CodeDeployerAddress = input.CodeDeployerAddress

	var empty struct{}
	context.codeUpdates[string(input.ContractAddress)] = empty
}

// CreateVMOutputInCaseOfError creates a new vmOutput with the given error set as return message.
func (context *outputContext) CreateVMOutputInCaseOfError(err error) *vmcommon.VMOutput {
	runtime := context.host.Runtime()
	runtime.AddError(err, runtime.Function())

	returnCode := context.resolveReturnCodeFromError(err)
	returnMessage := context.resolveReturnMessageFromError(err)

	vmOutput := &vmcommon.VMOutput{
		GasRemaining:  0,
		GasRefund:     big.NewInt(0),
		ReturnCode:    returnCode,
		ReturnMessage: returnMessage,
	}

	context.host.Metering().UpdateGasStateOnFailure(vmOutput)

	return vmOutput
}

func (context *outputContext) removeNonUpdatedCode() {
	for address, account := range context.outputState.OutputAccounts {
		_, ok := context.codeUpdates[address]
		if !ok {
			account.Code = nil
			account.CodeMetadata = nil
			account.CodeDeployerAddress = nil
		}
	}
}

func (context *outputContext) resolveReturnMessageFromError(err error) string {
	if errors.Is(err, vmhost.ErrSignalError) {
		return context.ReturnMessage()
	}
	if errors.Is(err, vmhost.ErrMemoryLimit) {
		// ErrMemoryLimit will still produce the 'execution failed' message.
		return vmhost.ErrExecutionFailed.Error()
	}
	if len(context.outputState.ReturnMessage) > 0 {
		// Another return message was already set.
		return context.outputState.ReturnMessage
	}

	return err.Error()
}

func (context *outputContext) resolveReturnCodeFromError(err error) vmcommon.ReturnCode {
	if err == nil {
		return vmcommon.Ok
	}

	if errors.Is(err, vmhost.ErrSignalError) {
		return vmcommon.UserError
	}
	if errors.Is(err, vmhost.ErrFuncNotFound) {
		return vmcommon.FunctionNotFound
	}
	if errors.Is(err, vmhost.ErrFunctionNonvoidSignature) {
		return vmcommon.FunctionWrongSignature
	}
	if errors.Is(err, vmhost.ErrInvalidFunction) {
		return vmcommon.UserError
	}
	if errors.Is(err, vmhost.ErrNotEnoughGas) {
		return vmcommon.OutOfGas
	}
	if errors.Is(err, vmhost.ErrContractNotFound) {
		return vmcommon.ContractNotFound
	}
	if errors.Is(err, vmhost.ErrContractInvalid) {
		return vmcommon.ContractInvalid
	}
	if errors.Is(err, vmhost.ErrUpgradeFailed) {
		return vmcommon.UpgradeFailed
	}
	if errors.Is(err, vmhost.ErrTransferInsufficientFunds) {
		return vmcommon.OutOfFunds
	}

	return vmcommon.ExecutionFailed
}

// AddToActiveState merges the given vmOutput with the outputState.
func (context *outputContext) AddToActiveState(rightOutput *vmcommon.VMOutput) {
	if rightOutput.GasRefund != nil {
		rightOutput.GasRefund.Add(rightOutput.GasRefund, context.outputState.GasRefund)
	}

	for _, rightAccount := range rightOutput.OutputAccounts {
		leftAccount, ok := context.outputState.OutputAccounts[string(rightAccount.Address)]
		if !ok {
			continue
		}

		if rightAccount.BalanceDelta != nil {
			rightAccount.BalanceDelta.Add(rightAccount.BalanceDelta, leftAccount.BalanceDelta)
		}
		if len(rightAccount.OutputTransfers) > 0 {
			leftAccount.OutputTransfers = append(leftAccount.OutputTransfers, rightAccount.OutputTransfers...)
		}
	}

	mergeVMOutputs(context.outputState, rightOutput)
}

func mergeVMOutputs(leftOutput *vmcommon.VMOutput, rightOutput *vmcommon.VMOutput) {
	if leftOutput.OutputAccounts == nil {
		leftOutput.OutputAccounts = make(map[string]*vmcommon.OutputAccount)
	}

	for _, rightAccount := range rightOutput.OutputAccounts {
		leftAccount, ok := leftOutput.OutputAccounts[string(rightAccount.Address)]
		if !ok {
			leftAccount = &vmcommon.OutputAccount{}
			leftOutput.OutputAccounts[string(rightAccount.Address)] = leftAccount
		}
		mergeOutputAccounts(leftAccount, rightAccount)
	}

	leftOutput.Logs = append(leftOutput.Logs, rightOutput.Logs...)
	leftOutput.ReturnData = append(leftOutput.ReturnData, rightOutput.ReturnData...)
	leftOutput.GasRemaining = rightOutput.GasRemaining
	leftOutput.GasRefund = rightOutput.GasRefund
	if leftOutput.GasRefund == nil {
		leftOutput.GasRefund = big.NewInt(0)
	}

	leftOutput.ReturnCode = rightOutput.ReturnCode
	leftOutput.ReturnMessage = rightOutput.ReturnMessage
}

func mergeOutputAccounts(
	leftAccount *vmcommon.OutputAccount,
	rightAccount *vmcommon.OutputAccount,
) {
	if len(rightAccount.Address) != 0 {
		leftAccount.Address = rightAccount.Address
	}

	mergeStorageUpdates(leftAccount, rightAccount)

	if rightAccount.Balance != nil {
		leftAccount.Balance = rightAccount.Balance
	}
	if leftAccount.BalanceDelta == nil {
		leftAccount.BalanceDelta = big.NewInt(0)
	}
	if rightAccount.BalanceDelta != nil {
		leftAccount.BalanceDelta = rightAccount.BalanceDelta
	}
	if len(rightAccount.Code) > 0 {
		leftAccount.Code = rightAccount.Code
	}
	if len(rightAccount.CodeMetadata) > 0 {
		leftAccount.CodeMetadata = rightAccount.CodeMetadata
	}
	if rightAccount.Nonce > leftAccount.Nonce {
		leftAccount.Nonce = rightAccount.Nonce
	}

	lenLeftOutTransfers := len(leftAccount.OutputTransfers)
	lenRightOutTransfers := len(rightAccount.OutputTransfers)
	if lenRightOutTransfers > lenLeftOutTransfers {
		leftAccount.OutputTransfers = append(leftAccount.OutputTransfers, rightAccount.OutputTransfers[lenLeftOutTransfers:]...)
	}

	leftAccount.GasUsed = rightAccount.GasUsed

	if rightAccount.CodeDeployerAddress != nil {
		leftAccount.CodeDeployerAddress = rightAccount.CodeDeployerAddress
	}

	if rightAccount.BytesAddedToStorage > leftAccount.BytesAddedToStorage {
		leftAccount.BytesAddedToStorage = rightAccount.BytesAddedToStorage
	}
	if rightAccount.BytesDeletedFromStorage > leftAccount.BytesDeletedFromStorage {
		leftAccount.BytesDeletedFromStorage = rightAccount.BytesDeletedFromStorage
	}
}

func mergeStorageUpdates(
	leftAccount *vmcommon.OutputAccount,
	rightAccount *vmcommon.OutputAccount,
) {
	if leftAccount.StorageUpdates == nil {
		leftAccount.StorageUpdates = make(map[string]*vmcommon.StorageUpdate)
	}
	for key, update := range rightAccount.StorageUpdates {
		leftAccount.StorageUpdates[key] = update
	}
}
