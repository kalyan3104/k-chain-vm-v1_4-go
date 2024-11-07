package vmhooks

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
// typedef unsigned long long uint64_t;
//
// extern void	v1_4_managedSCAddress(void *context, int32_t addressHandle);
// extern void	v1_4_managedOwnerAddress(void *context, int32_t addressHandle);
// extern void	v1_4_managedCaller(void *context, int32_t addressHandle);
// extern void	v1_4_managedSignalError(void* context, int32_t errHandle1);
// extern void	v1_4_managedWriteLog(void* context, int32_t topicsHandle, int32_t dataHandle);
//
// extern int32_t	v1_4_managedMultiTransferDCDTNFTExecute(void *context, int32_t dstHandle, int32_t tokenTransfersHandle, long long gasLimit, int32_t functionHandle, int32_t argumentsHandle);
// extern int32_t	v1_4_managedTransferValueExecute(void *context, int32_t dstHandle, int32_t valueHandle, long long gasLimit, int32_t functionHandle, int32_t argumentsHandle);
// extern int32_t	v1_4_managedExecuteOnDestContext(void *context, long long gas, int32_t addressHandle, int32_t valueHandle, int32_t functionHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern int32_t	v1_4_managedExecuteOnDestContextByCaller(void *context, long long gas, int32_t addressHandle, int32_t valueHandle, int32_t functionHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern int32_t	v1_4_managedExecuteOnSameContext(void *context, long long gas, int32_t addressHandle, int32_t valueHandle, int32_t functionHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern int32_t	v1_4_managedExecuteReadOnly(void *context, long long gas, int32_t addressHandle, int32_t functionHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern int32_t	v1_4_managedCreateContract(void *context, long long gas, int32_t valueHandle, int32_t codeHandle, int32_t codeMetadataHandle, int32_t argumentsHandle, int32_t resultAddressHandle, int32_t resultHandle);
// extern int32_t	v1_4_managedDeployFromSourceContract(void *context, long long gas, int32_t valueHandle, int32_t addressHandle, int32_t codeMetadataHandle, int32_t argumentsHandle, int32_t resultAddressHandle, int32_t resultHandle);
// extern void		v1_4_managedUpgradeContract(void *context, int32_t dstHandle, long long gas, int32_t valueHandle, int32_t codeHandle, int32_t codeMetadataHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern void		v1_4_managedUpgradeFromSourceContract(void *context, int32_t dstHandle, long long gas, int32_t valueHandle, int32_t addressHandle, int32_t codeMetadataHandle, int32_t argumentsHandle, int32_t resultHandle);
// extern void		v1_4_managedAsyncCall(void *context, int32_t dstHandle, int32_t valueHandle, int32_t functionHandle, int32_t argumentsHandle);
//
// extern void		v1_4_managedGetMultiDCDTCallValue(void *context, int32_t multiCallValueHandle);
// extern void		v1_4_managedGetDCDTBalance(void *context, int32_t addressHandle, int32_t tokenIDHandle, long long nonce, int32_t valueHandle);
// extern void		v1_4_managedGetDCDTTokenData(void *context, int32_t addressHandle, int32_t tokenIDHandle, long long nonce, int32_t valueHandle, int32_t propertiesHandle, int32_t hashHandle, int32_t nameHandle, int32_t attributesHandle, int32_t creatorHandle, int32_t royaltiesHandle, int32_t urisHandle);
//
// extern void		v1_4_managedGetReturnData(void *context, int32_t resultID, int32_t resultHandle);
// extern void		v1_4_managedGetPrevBlockRandomSeed(void *context, int32_t resultHandle);
// extern void		v1_4_managedGetBlockRandomSeed(void *context, int32_t resultHandle);
// extern void		v1_4_managedGetStateRootHash(void *context, int32_t resultHandle);
// extern void		v1_4_managedGetOriginalTxHash(void *context, int32_t resultHandle);
//
// extern int32_t   v1_4_managedIsDCDTFrozen(void *context, int32_t addressHandle, int32_t tokenIDHandle, long long nonce);
// extern int32_t   v1_4_managedIsDCDTPaused(void *context, int32_t tokenIDHandle);
// extern int32_t   v1_4_managedIsDCDTLimitedTransfer(void *context, int32_t tokenIDHandle);
// extern void      v1_4_managedBufferToHex(void *context, int32_t sourceHandle, int32_t destHandle);
import "C"

import (
	"encoding/hex"
	"errors"
	"unsafe"

	"github.com/kalyan3104/k-chain-vm-common-go/builtInFunctions"

	"github.com/kalyan3104/k-chain-vm-v1_4-go/math"
	"github.com/kalyan3104/k-chain-vm-v1_4-go/vmhost"
	"github.com/kalyan3104/k-chain-vm-v1_4-go/vmhost/vmhooksmeta"
)

const (
	managedSCAddressName                    = "managedSCAddress"
	managedOwnerAddressName                 = "managedOwnerAddress"
	managedCallerName                       = "managedCaller"
	managedSignalErrorName                  = "managedSignalError"
	managedWriteLogName                     = "managedWriteLog"
	managedMultiTransferDCDTNFTExecuteName  = "managedMultiTransferDCDTNFTExecute"
	managedTransferValueExecuteName         = "managedTransferValueExecute"
	managedExecuteOnDestContextName         = "managedExecuteOnDestContext"
	managedExecuteOnDestContextByCallerName = "managedExecuteOnDestContextByCaller"
	managedExecuteOnSameContextName         = "managedExecuteOnSameContext"
	managedExecuteReadOnlyName              = "managedExecuteReadOnly"
	managedCreateContractName               = "managedCreateContract"
	managedDeployFromSourceContractName     = "managedDeployFromSourceContract"
	managedUpgradeContractName              = "managedUpgradeContract"
	managedUpgradeFromSourceContractName    = "managedUpgradeFromSourceContract"
	managedAsyncCallName                    = "managedAsyncCall"
	managedGetMultiDCDTCallValueName        = "managedGetMultiDCDTCallValue"
	managedGetDCDTBalanceName               = "managedGetDCDTBalance"
	managedGetDCDTTokenDataName             = "managedGetDCDTTokenData"
	managedGetReturnDataName                = "managedGetReturnData"
	managedGetPrevBlockRandomSeedName       = "managedGetPrevBlockRandomSeed"
	managedGetBlockRandomSeedName           = "managedGetBlockRandomSeed"
	managedGetStateRootHashName             = "managedGetStateRootHash"
	managedGetOriginalTxHashName            = "managedGetOriginalTxHash"
	managedIsDCDTFrozenName                 = "managedIsDCDTFrozen"
	managedIsDCDTLimitedTransferName        = "managedIsDCDTLimitedTransfer"
	managedIsDCDTPausedName                 = "managedIsDCDTPaused"
	managedBufferToHexName                  = "managedBufferToHex"
)

// ManagedEIImports creates a new wasmer.Imports populated with variants of the API methods that use managed types only.
func ManagedEIImports(imports vmhooksmeta.EIFunctionReceiver) error {
	imports.Namespace("env")

	err := imports.Append("managedSCAddress", v1_4_managedSCAddress, C.v1_4_managedSCAddress)
	if err != nil {
		return err
	}

	err = imports.Append("managedOwnerAddress", v1_4_managedOwnerAddress, C.v1_4_managedOwnerAddress)
	if err != nil {
		return err
	}

	err = imports.Append("managedCaller", v1_4_managedCaller, C.v1_4_managedCaller)
	if err != nil {
		return err
	}

	err = imports.Append("managedSignalError", v1_4_managedSignalError, C.v1_4_managedSignalError)
	if err != nil {
		return err
	}

	err = imports.Append("managedWriteLog", v1_4_managedWriteLog, C.v1_4_managedWriteLog)
	if err != nil {
		return err
	}

	err = imports.Append("managedMultiTransferDCDTNFTExecute", v1_4_managedMultiTransferDCDTNFTExecute, C.v1_4_managedMultiTransferDCDTNFTExecute)
	if err != nil {
		return err
	}

	err = imports.Append("managedTransferValueExecute", v1_4_managedTransferValueExecute, C.v1_4_managedTransferValueExecute)
	if err != nil {
		return err
	}

	err = imports.Append("managedExecuteOnDestContext", v1_4_managedExecuteOnDestContext, C.v1_4_managedExecuteOnDestContext)
	if err != nil {
		return err
	}

	err = imports.Append("managedExecuteOnDestContextByCaller", v1_4_managedExecuteOnDestContextByCaller, C.v1_4_managedExecuteOnDestContextByCaller)
	if err != nil {
		return err
	}

	err = imports.Append("managedExecuteOnSameContext", v1_4_managedExecuteOnSameContext, C.v1_4_managedExecuteOnSameContext)
	if err != nil {
		return err
	}

	err = imports.Append("managedExecuteReadOnly", v1_4_managedExecuteReadOnly, C.v1_4_managedExecuteReadOnly)
	if err != nil {
		return err
	}

	err = imports.Append("managedCreateContract", v1_4_managedCreateContract, C.v1_4_managedCreateContract)
	if err != nil {
		return err
	}

	err = imports.Append("managedDeployFromSourceContract", v1_4_managedDeployFromSourceContract, C.v1_4_managedDeployFromSourceContract)
	if err != nil {
		return err
	}

	err = imports.Append("managedUpgradeContract", v1_4_managedUpgradeContract, C.v1_4_managedUpgradeContract)
	if err != nil {
		return err
	}

	err = imports.Append("managedUpgradeFromSourceContract", v1_4_managedUpgradeFromSourceContract, C.v1_4_managedUpgradeFromSourceContract)
	if err != nil {
		return err
	}

	err = imports.Append("managedAsyncCall", v1_4_managedAsyncCall, C.v1_4_managedAsyncCall)
	if err != nil {
		return err
	}

	err = imports.Append("managedGetMultiDCDTCallValue", v1_4_managedGetMultiDCDTCallValue, C.v1_4_managedGetMultiDCDTCallValue)
	if err != nil {
		return err
	}

	err = imports.Append("managedGetDCDTBalance", v1_4_managedGetDCDTBalance, C.v1_4_managedGetDCDTBalance)
	if err != nil {
		return err
	}

	err = imports.Append("managedGetDCDTTokenData", v1_4_managedGetDCDTTokenData, C.v1_4_managedGetDCDTTokenData)
	if err != nil {
		return err
	}

	err = imports.Append("managedGetReturnData", v1_4_managedGetReturnData, C.v1_4_managedGetReturnData)
	if err != nil {
		return err
	}

	err = imports.Append("managedGetPrevBlockRandomSeed", v1_4_managedGetPrevBlockRandomSeed, C.v1_4_managedGetPrevBlockRandomSeed)
	if err != nil {
		return err
	}

	err = imports.Append("managedGetBlockRandomSeed", v1_4_managedGetBlockRandomSeed, C.v1_4_managedGetBlockRandomSeed)
	if err != nil {
		return err
	}

	err = imports.Append("managedGetStateRootHash", v1_4_managedGetStateRootHash, C.v1_4_managedGetStateRootHash)
	if err != nil {
		return err
	}

	err = imports.Append("managedGetOriginalTxHash", v1_4_managedGetOriginalTxHash, C.v1_4_managedGetOriginalTxHash)
	if err != nil {
		return err
	}

	err = imports.Append("managedIsDCDTFrozen", v1_4_managedIsDCDTFrozen, C.v1_4_managedIsDCDTFrozen)
	if err != nil {
		return err
	}

	err = imports.Append("managedIsDCDTPaused", v1_4_managedIsDCDTPaused, C.v1_4_managedIsDCDTPaused)
	if err != nil {
		return err
	}

	err = imports.Append("managedIsDCDTLimitedTransfer", v1_4_managedIsDCDTLimitedTransfer, C.v1_4_managedIsDCDTLimitedTransfer)
	if err != nil {
		return err
	}

	err = imports.Append("managedBufferToHex", v1_4_managedBufferToHex, C.v1_4_managedBufferToHex)
	if err != nil {
		return err
	}

	return nil
}

//export v1_4_managedSCAddress
func v1_4_managedSCAddress(context unsafe.Pointer, destinationHandle int32) {
	managedType := vmhost.GetManagedTypesContext(context)
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetSCAddress
	metering.UseGasAndAddTracedGas(managedSCAddressName, gasToUse)

	scAddress := runtime.GetContextAddress()

	managedType.SetBytes(destinationHandle, scAddress)
}

//export v1_4_managedOwnerAddress
func v1_4_managedOwnerAddress(context unsafe.Pointer, destinationHandle int32) {
	managedType := vmhost.GetManagedTypesContext(context)
	blockchain := vmhost.GetBlockchainContext(context)
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetOwnerAddress
	metering.UseGasAndAddTracedGas(managedOwnerAddressName, gasToUse)

	owner, err := blockchain.GetOwnerAddress()
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	managedType.SetBytes(destinationHandle, owner)
}

//export v1_4_managedCaller
func v1_4_managedCaller(context unsafe.Pointer, destinationHandle int32) {
	managedType := vmhost.GetManagedTypesContext(context)
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCaller
	metering.UseGasAndAddTracedGas(managedCallerName, gasToUse)

	caller := runtime.GetVMInput().CallerAddr
	managedType.SetBytes(destinationHandle, caller)
}

//export v1_4_managedSignalError
func v1_4_managedSignalError(context unsafe.Pointer, errHandle int32) {
	managedType := vmhost.GetManagedTypesContext(context)
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)
	metering.StartGasTracing(managedSignalErrorName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.SignalError
	metering.UseAndTraceGas(gasToUse)

	errBytes, err := managedType.GetBytes(errHandle)
	if vmhost.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBytes(errBytes)

	gasToUse = metering.GasSchedule().BaseOperationCost.PersistPerByte * uint64(len(errBytes))
	err = metering.UseGasBounded(gasToUse)
	if err != nil {
		_ = vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution())
		return
	}

	runtime.SignalUserError(string(errBytes))
}

//export v1_4_managedWriteLog
func v1_4_managedWriteLog(
	context unsafe.Pointer,
	topicsHandle int32,
	dataHandle int32,
) {
	runtime := vmhost.GetRuntimeContext(context)
	output := vmhost.GetOutputContext(context)
	metering := vmhost.GetMeteringContext(context)
	managedType := vmhost.GetManagedTypesContext(context)
	metering.StartGasTracing(managedWriteLogName)

	topics, sumOfTopicByteLengths, err := managedType.ReadManagedVecOfManagedBuffers(topicsHandle)
	if vmhost.WithFault(err, context, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	dataBytes, err := managedType.GetBytes(dataHandle)
	if vmhost.WithFault(err, context, runtime.ManagedBufferAPIErrorShouldFailExecution()) {
		return
	}
	managedType.ConsumeGasForBytes(dataBytes)
	dataByteLen := uint64(len(dataBytes))

	gasToUse := metering.GasSchedule().BaseOpsAPICost.Log
	gasForData := math.MulUint64(
		metering.GasSchedule().BaseOperationCost.DataCopyPerByte,
		sumOfTopicByteLengths+dataByteLen)
	gasToUse = math.AddUint64(gasToUse, gasForData)
	metering.UseAndTraceGas(gasToUse)

	output.WriteLog(runtime.GetContextAddress(), topics, dataBytes)
}

//export v1_4_managedGetOriginalTxHash
func v1_4_managedGetOriginalTxHash(context unsafe.Pointer, resultHandle int32) {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)
	managedType := vmhost.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetOriginalTxHash
	metering.UseGasAndAddTracedGas(managedGetOriginalTxHashName, gasToUse)

	managedType.SetBytes(resultHandle, runtime.GetOriginalTxHash())
}

//export v1_4_managedGetStateRootHash
func v1_4_managedGetStateRootHash(context unsafe.Pointer, resultHandle int32) {
	blockchain := vmhost.GetBlockchainContext(context)
	metering := vmhost.GetMeteringContext(context)
	managedType := vmhost.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetStateRootHash
	metering.UseGasAndAddTracedGas(managedGetStateRootHashName, gasToUse)

	managedType.SetBytes(resultHandle, blockchain.GetStateRootHash())
}

//export v1_4_managedGetBlockRandomSeed
func v1_4_managedGetBlockRandomSeed(context unsafe.Pointer, resultHandle int32) {
	blockchain := vmhost.GetBlockchainContext(context)
	metering := vmhost.GetMeteringContext(context)
	managedType := vmhost.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRandomSeed
	metering.UseGasAndAddTracedGas(managedGetBlockRandomSeedName, gasToUse)

	managedType.SetBytes(resultHandle, blockchain.CurrentRandomSeed())
}

//export v1_4_managedGetPrevBlockRandomSeed
func v1_4_managedGetPrevBlockRandomSeed(context unsafe.Pointer, resultHandle int32) {
	blockchain := vmhost.GetBlockchainContext(context)
	metering := vmhost.GetMeteringContext(context)
	managedType := vmhost.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetBlockRandomSeed
	metering.UseGasAndAddTracedGas(managedGetPrevBlockRandomSeedName, gasToUse)

	managedType.SetBytes(resultHandle, blockchain.LastRandomSeed())
}

//export v1_4_managedGetReturnData
func v1_4_managedGetReturnData(context unsafe.Pointer, resultID int32, resultHandle int32) {
	runtime := vmhost.GetRuntimeContext(context)
	output := vmhost.GetOutputContext(context)
	metering := vmhost.GetMeteringContext(context)
	managedType := vmhost.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetReturnData
	metering.UseGasAndAddTracedGas(managedGetReturnDataName, gasToUse)

	returnData := output.ReturnData()
	if resultID >= int32(len(returnData)) || resultID < 0 {
		_ = vmhost.WithFault(vmhost.ErrArgOutOfRange, context, runtime.BaseOpsErrorShouldFailExecution())
		return
	}

	managedType.SetBytes(resultHandle, returnData[resultID])
}

//export v1_4_managedGetMultiDCDTCallValue
func v1_4_managedGetMultiDCDTCallValue(context unsafe.Pointer, multiCallValueHandle int32) {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)
	managedType := vmhost.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetCallValue
	metering.UseGasAndAddTracedGas(managedGetMultiDCDTCallValueName, gasToUse)

	dcdtTransfers := runtime.GetVMInput().DCDTTransfers
	multiCallBytes := writeDCDTTransfersToBytes(managedType, dcdtTransfers)
	managedType.ConsumeGasForBytes(multiCallBytes)

	managedType.SetBytes(multiCallValueHandle, multiCallBytes)
}

//export v1_4_managedGetDCDTBalance
func v1_4_managedGetDCDTBalance(context unsafe.Pointer, addressHandle int32, tokenIDHandle int32, nonce int64, valueHandle int32) {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)
	blockchain := vmhost.GetBlockchainContext(context)
	managedType := vmhost.GetManagedTypesContext(context)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	metering.UseGasAndAddTracedGas(managedGetDCDTBalanceName, gasToUse)

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		_ = vmhost.WithFault(vmhost.ErrArgOutOfRange, context, runtime.BaseOpsErrorShouldFailExecution())
		return
	}
	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		_ = vmhost.WithFault(vmhost.ErrArgOutOfRange, context, runtime.BaseOpsErrorShouldFailExecution())
		return
	}

	dcdtToken, err := blockchain.GetDCDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		_ = vmhost.WithFault(vmhost.ErrArgOutOfRange, context, runtime.BaseOpsErrorShouldFailExecution())
		return
	}

	value := managedType.GetBigIntOrCreate(valueHandle)
	value.Set(dcdtToken.Value)
}

//export v1_4_managedGetDCDTTokenData
func v1_4_managedGetDCDTTokenData(context unsafe.Pointer, addressHandle int32, tokenIDHandle int32, nonce int64,
	valueHandle, propertiesHandle, hashHandle, nameHandle, attributesHandle, creatorHandle, royaltiesHandle, urisHandle int32) {
	runtime := vmhost.GetRuntimeContext(context)
	metering := vmhost.GetMeteringContext(context)
	blockchain := vmhost.GetBlockchainContext(context)
	managedType := vmhost.GetManagedTypesContext(context)
	metering.StartGasTracing(managedGetDCDTTokenDataName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	metering.UseAndTraceGas(gasToUse)

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		_ = vmhost.WithFault(vmhost.ErrArgOutOfRange, context, runtime.BaseOpsErrorShouldFailExecution())
		return
	}
	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		_ = vmhost.WithFault(vmhost.ErrArgOutOfRange, context, runtime.BaseOpsErrorShouldFailExecution())
		return
	}

	dcdtToken, err := blockchain.GetDCDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		_ = vmhost.WithFault(vmhost.ErrArgOutOfRange, context, runtime.BaseOpsErrorShouldFailExecution())
		return
	}

	value := managedType.GetBigIntOrCreate(valueHandle)
	value.Set(dcdtToken.Value)

	managedType.SetBytes(propertiesHandle, dcdtToken.Properties)
	if dcdtToken.TokenMetaData != nil {
		managedType.SetBytes(hashHandle, dcdtToken.TokenMetaData.Hash)
		managedType.ConsumeGasForBytes(dcdtToken.TokenMetaData.Hash)
		managedType.SetBytes(nameHandle, dcdtToken.TokenMetaData.Name)
		managedType.ConsumeGasForBytes(dcdtToken.TokenMetaData.Name)
		managedType.SetBytes(attributesHandle, dcdtToken.TokenMetaData.Attributes)
		managedType.ConsumeGasForBytes(dcdtToken.TokenMetaData.Attributes)
		managedType.SetBytes(creatorHandle, dcdtToken.TokenMetaData.Creator)
		managedType.ConsumeGasForBytes(dcdtToken.TokenMetaData.Creator)
		royalties := managedType.GetBigIntOrCreate(royaltiesHandle)
		royalties.SetUint64(uint64(dcdtToken.TokenMetaData.Royalties))

		managedType.WriteManagedVecOfManagedBuffers(dcdtToken.TokenMetaData.URIs, urisHandle)
	}

}

//export v1_4_managedAsyncCall
func v1_4_managedAsyncCall(
	context unsafe.Pointer,
	destHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32) {
	host := vmhost.GetVMHost(context)
	ManagedAsyncCallWithHost(
		host,
		destHandle,
		valueHandle,
		functionHandle,
		argumentsHandle)
}

func ManagedAsyncCallWithHost(
	host vmhost.VMHost,
	destHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32) {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedAsyncCallName)

	gasSchedule := metering.GasSchedule()
	gasToUse := gasSchedule.BaseOpsAPICost.AsyncCallStep
	metering.UseAndTraceGas(gasToUse)

	vmInput, err := readDestinationFunctionArguments(host, destHandle, functionHandle, argumentsHandle)
	if vmhost.WithFaultAndHost(host, err, host.Runtime().BaseOpsErrorShouldFailExecution()) {
		return
	}

	data := makeCrossShardCallFromInput(vmInput.function, vmInput.arguments)

	value, err := managedType.GetBigInt(valueHandle)
	if err != nil {
		_ = vmhost.WithFaultAndHost(host, vmhost.ErrArgOutOfRange, host.Runtime().BaseOpsErrorShouldFailExecution())
		return
	}

	gasToUse = math.MulUint64(gasSchedule.BaseOperationCost.DataCopyPerByte, uint64(len(data)))
	metering.UseAndTraceGas(gasToUse)

	err = runtime.ExecuteAsyncCall(vmInput.destination, []byte(data), value.Bytes())
	if errors.Is(err, vmhost.ErrNotEnoughGas) {
		runtime.SetRuntimeBreakpointValue(vmhost.BreakpointOutOfGas)
		return
	}
	if vmhost.WithFaultAndHost(host, err, host.Runtime().BaseOpsErrorShouldFailExecution()) {
		return
	}
}

//export v1_4_managedUpgradeFromSourceContract
func v1_4_managedUpgradeFromSourceContract(
	context unsafe.Pointer,
	destHandle int32,
	gas int64,
	valueHandle int32,
	addressHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) {
	host := vmhost.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedUpgradeFromSourceContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	metering.UseAndTraceGas(gasToUse)

	vmInput, err := readDestinationValueArguments(host, destHandle, valueHandle, argumentsHandle)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	sourceContractAddress, err := managedType.GetBytes(addressHandle)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	lenReturnData := len(host.Output().ReturnData())

	UpgradeFromSourceContractWithTypedArgs(
		host,
		sourceContractAddress,
		vmInput.destination,
		vmInput.value.Bytes(),
		vmInput.arguments,
		gas,
		codeMetadata,
	)
	setReturnDataIfExists(host, lenReturnData, resultHandle)
}

//export v1_4_managedUpgradeContract
func v1_4_managedUpgradeContract(
	context unsafe.Pointer,
	destHandle int32,
	gas int64,
	valueHandle int32,
	codeHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) {
	host := vmhost.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedUpgradeContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	metering.UseAndTraceGas(gasToUse)

	vmInput, err := readDestinationValueArguments(host, destHandle, valueHandle, argumentsHandle)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	code, err := managedType.GetBytes(codeHandle)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return
	}

	lenReturnData := len(host.Output().ReturnData())

	upgradeContract(host, vmInput.destination, code, codeMetadata, vmInput.value.Bytes(), vmInput.arguments, gas)
	setReturnDataIfExists(host, lenReturnData, resultHandle)
}

//export v1_4_managedDeployFromSourceContract
func v1_4_managedDeployFromSourceContract(
	context unsafe.Pointer,
	gas int64,
	valueHandle int32,
	addressHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultAddressHandle int32,
	resultHandle int32,
) int32 {
	host := vmhost.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedDeployFromSourceContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	metering.UseAndTraceGas(gasToUse)

	vmInput, err := readDestinationValueArguments(host, addressHandle, valueHandle, argumentsHandle)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	lenReturnData := len(host.Output().ReturnData())

	newAddress, err := DeployFromSourceContractWithTypedArgs(
		host,
		vmInput.destination,
		codeMetadata,
		vmInput.value,
		vmInput.arguments,
		gas,
	)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	managedType.SetBytes(resultAddressHandle, newAddress)
	setReturnDataIfExists(host, lenReturnData, resultHandle)

	return 0
}

//export v1_4_managedCreateContract
func v1_4_managedCreateContract(
	context unsafe.Pointer,
	gas int64,
	valueHandle int32,
	codeHandle int32,
	codeMetadataHandle int32,
	argumentsHandle int32,
	resultAddressHandle int32,
	resultHandle int32,
) int32 {
	host := vmhost.GetVMHost(context)
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()
	metering.StartGasTracing(managedCreateContractName)

	gasToUse := metering.GasSchedule().BaseOpsAPICost.CreateContract
	metering.UseAndTraceGas(gasToUse)

	sender := runtime.GetContextAddress()
	value, err := managedType.GetBigInt(valueHandle)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	data, actualLen, err := managedType.ReadManagedVecOfManagedBuffers(argumentsHandle)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	gasToUse = math.MulUint64(metering.GasSchedule().BaseOperationCost.DataCopyPerByte, actualLen)
	metering.UseAndTraceGas(gasToUse)

	codeMetadata, err := managedType.GetBytes(codeMetadataHandle)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	code, err := managedType.GetBytes(codeHandle)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	lenReturnData := len(host.Output().ReturnData())
	newAddress, err := createContract(sender, data, value, metering, gas, code, codeMetadata, host, runtime)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return 1
	}

	managedType.SetBytes(resultAddressHandle, newAddress)
	setReturnDataIfExists(host, lenReturnData, resultHandle)

	return 0
}

func setReturnDataIfExists(
	host vmhost.VMHost,
	oldLen int,
	resultHandle int32,
) {
	returnData := host.Output().ReturnData()
	if len(returnData) > oldLen {
		host.ManagedTypes().WriteManagedVecOfManagedBuffers(returnData[oldLen:], resultHandle)
	} else {
		host.ManagedTypes().SetBytes(resultHandle, make([]byte, 0))
	}
}

//export v1_4_managedExecuteReadOnly
func v1_4_managedExecuteReadOnly(
	context unsafe.Pointer,
	gas int64,
	addressHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) int32 {
	host := vmhost.GetVMHost(context)
	metering := host.Metering()
	metering.StartGasTracing(managedExecuteReadOnlyName)

	vmInput, err := readDestinationFunctionArguments(host, addressHandle, functionHandle, argumentsHandle)
	if vmhost.WithFaultAndHost(host, err, host.Runtime().BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	lenReturnData := len(host.Output().ReturnData())
	returnVal := ExecuteReadOnlyWithTypedArguments(
		host,
		gas,
		[]byte(vmInput.function),
		vmInput.destination,
		vmInput.arguments,
	)
	setReturnDataIfExists(host, lenReturnData, resultHandle)
	return returnVal
}

//export v1_4_managedExecuteOnSameContext
func v1_4_managedExecuteOnSameContext(
	context unsafe.Pointer,
	gas int64,
	addressHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) int32 {
	host := vmhost.GetVMHost(context)
	metering := host.Metering()
	metering.StartGasTracing(managedExecuteOnSameContextName)

	vmInput, err := readDestinationValueFunctionArguments(host, addressHandle, valueHandle, functionHandle, argumentsHandle)
	if vmhost.WithFaultAndHost(host, err, host.Runtime().BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	lenReturnData := len(host.Output().ReturnData())
	returnVal := ExecuteOnSameContextWithTypedArgs(
		host,
		gas,
		vmInput.value,
		[]byte(vmInput.function),
		vmInput.destination,
		vmInput.arguments,
	)
	setReturnDataIfExists(host, lenReturnData, resultHandle)
	return returnVal
}

//export v1_4_managedExecuteOnDestContextByCaller
func v1_4_managedExecuteOnDestContextByCaller(
	context unsafe.Pointer,
	gas int64,
	addressHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) int32 {
	host := vmhost.GetVMHost(context)
	metering := host.Metering()
	metering.StartGasTracing(managedExecuteOnDestContextByCallerName)

	vmInput, err := readDestinationValueFunctionArguments(host, addressHandle, valueHandle, functionHandle, argumentsHandle)
	if vmhost.WithFaultAndHost(host, err, host.Runtime().BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	lenReturnData := len(host.Output().ReturnData())
	returnVal := ExecuteOnDestContextByCallerWithTypedArgs(
		host,
		gas,
		vmInput.value,
		[]byte(vmInput.function),
		vmInput.destination,
		vmInput.arguments,
	)
	setReturnDataIfExists(host, lenReturnData, resultHandle)
	return returnVal
}

//export v1_4_managedExecuteOnDestContext
func v1_4_managedExecuteOnDestContext(
	context unsafe.Pointer,
	gas int64,
	addressHandle int32,
	valueHandle int32,
	functionHandle int32,
	argumentsHandle int32,
	resultHandle int32,
) int32 {
	host := vmhost.GetVMHost(context)
	metering := host.Metering()
	metering.StartGasTracing(managedExecuteOnDestContextName)

	vmInput, err := readDestinationValueFunctionArguments(host, addressHandle, valueHandle, functionHandle, argumentsHandle)
	if vmhost.WithFaultAndHost(host, err, host.Runtime().BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	lenReturnData := len(host.Output().ReturnData())
	returnVal := ExecuteOnDestContextWithTypedArgs(
		host,
		gas,
		vmInput.value,
		[]byte(vmInput.function),
		vmInput.destination,
		vmInput.arguments,
	)
	setReturnDataIfExists(host, lenReturnData, resultHandle)
	return returnVal
}

//export v1_4_managedMultiTransferDCDTNFTExecute
func v1_4_managedMultiTransferDCDTNFTExecute(
	context unsafe.Pointer,
	dstHandle int32,
	tokenTransfersHandle int32,
	gasLimit int64,
	functionHandle int32,
	argumentsHandle int32,
) int32 {
	host := vmhost.GetVMHost(context)
	managedType := host.ManagedTypes()
	runtime := host.Runtime()
	metering := host.Metering()
	metering.StartGasTracing(managedMultiTransferDCDTNFTExecuteName)

	vmInput, err := readDestinationFunctionArguments(host, dstHandle, functionHandle, argumentsHandle)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	transfers, err := readDCDTTransfers(managedType, tokenTransfersHandle)
	if vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return TransferDCDTNFTExecuteWithTypedArgs(
		host,
		vmInput.destination,
		transfers,
		gasLimit,
		[]byte(vmInput.function),
		vmInput.arguments,
	)
}

//export v1_4_managedTransferValueExecute
func v1_4_managedTransferValueExecute(
	context unsafe.Pointer,
	dstHandle int32,
	valueHandle int32,
	gasLimit int64,
	functionHandle int32,
	argumentsHandle int32,
) int32 {
	host := vmhost.GetVMHost(context)
	metering := host.Metering()
	metering.StartGasTracing(managedTransferValueExecuteName)

	vmInput, err := readDestinationValueFunctionArguments(host, dstHandle, valueHandle, functionHandle, argumentsHandle)
	if vmhost.WithFaultAndHost(host, err, host.Runtime().BaseOpsErrorShouldFailExecution()) {
		return -1
	}

	return TransferValueExecuteWithTypedArgs(
		host,
		vmInput.destination,
		vmInput.value,
		gasLimit,
		[]byte(vmInput.function),
		vmInput.arguments,
	)
}

//export v1_4_managedIsDCDTFrozen
func v1_4_managedIsDCDTFrozen(
	context unsafe.Pointer,
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64) int32 {
	host := vmhost.GetVMHost(context)
	return ManagedIsDCDTFrozenWithHost(host, addressHandle, tokenIDHandle, nonce)
}

func ManagedIsDCDTFrozenWithHost(
	host vmhost.VMHost,
	addressHandle int32,
	tokenIDHandle int32,
	nonce int64) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	metering.UseGasAndAddTracedGas(managedIsDCDTFrozenName, gasToUse)

	address, err := managedType.GetBytes(addressHandle)
	if err != nil {
		_ = vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution())
		return -1
	}
	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		_ = vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution())
		return -1
	}

	dcdtToken, err := blockchain.GetDCDTToken(address, tokenID, uint64(nonce))
	if err != nil {
		_ = vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution())
		return -1
	}

	dcdtUserData := builtInFunctions.DCDTUserMetadataFromBytes(dcdtToken.Properties)
	if dcdtUserData.Frozen {
		return 1
	}
	return 0
}

//export v1_4_managedIsDCDTLimitedTransfer
func v1_4_managedIsDCDTLimitedTransfer(context unsafe.Pointer, tokenIDHandle int32) int32 {
	host := vmhost.GetVMHost(context)
	return ManagedIsDCDTLimitedTransferWithHost(host, tokenIDHandle)
}

func ManagedIsDCDTLimitedTransferWithHost(host vmhost.VMHost, tokenIDHandle int32) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	metering.UseGasAndAddTracedGas(managedIsDCDTLimitedTransferName, gasToUse)

	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		_ = vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution())
		return -1
	}

	if blockchain.IsLimitedTransfer(tokenID) {
		return 1
	}

	return 0
}

//export v1_4_managedIsDCDTPaused
func v1_4_managedIsDCDTPaused(context unsafe.Pointer, tokenIDHandle int32) int32 {
	host := vmhost.GetVMHost(context)
	return ManagedIsDCDTPausedWithHost(host, tokenIDHandle)
}

func ManagedIsDCDTPausedWithHost(host vmhost.VMHost, tokenIDHandle int32) int32 {
	runtime := host.Runtime()
	metering := host.Metering()
	blockchain := host.Blockchain()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().BaseOpsAPICost.GetExternalBalance
	metering.UseGasAndAddTracedGas(managedIsDCDTPausedName, gasToUse)

	tokenID, err := managedType.GetBytes(tokenIDHandle)
	if err != nil {
		_ = vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution())
		return -1
	}

	if blockchain.IsPaused(tokenID) {
		return 1
	}

	return 0
}

//export v1_4_managedBufferToHex
func v1_4_managedBufferToHex(context unsafe.Pointer, sourceHandle int32, destHandle int32) {
	host := vmhost.GetVMHost(context)
	ManagedBufferToHexWithHost(host, sourceHandle, destHandle)
}

func ManagedBufferToHexWithHost(host vmhost.VMHost, sourceHandle int32, destHandle int32) {
	runtime := host.Runtime()
	metering := host.Metering()
	managedType := host.ManagedTypes()

	gasToUse := metering.GasSchedule().ManagedBufferAPICost.MBufferSetBytes
	metering.UseGasAndAddTracedGas(managedBufferToHexName, gasToUse)

	mBuff, err := managedType.GetBytes(sourceHandle)
	if err != nil {
		vmhost.WithFaultAndHost(host, err, runtime.BaseOpsErrorShouldFailExecution())
		return
	}

	encoded := hex.EncodeToString(mBuff)
	managedType.SetBytes(destHandle, []byte(encoded))
}
