package vmserver

import (
	"math/big"

	"github.com/kalyan3104/k-chain-vm-v1_4-go/config"
	"github.com/kalyan3104/k-chain-vm-v1_4-go/vmhost"
	"github.com/kalyan3104/k-chain-vm-v1_4-go/vmhost/hostCore"
	"github.com/kalyan3104/k-chain-vm-v1_4-go/vmhost/mock"

	worldmock "github.com/kalyan3104/k-chain-scenario-go/worldmock"
	vmcommon "github.com/kalyan3104/k-chain-vm-common-go"
	"github.com/kalyan3104/k-chain-vm-common-go/builtInFunctions"
	"github.com/kalyan3104/k-chain-vm-common-go/parsers"
)

type worldDataModel struct {
	ID       string
	Accounts worldmock.AccountMap
}

type world struct {
	id             string
	blockchainHook *worldmock.MockWorld
	vm             vmcommon.VMExecutionHandler
}

func newWorldDataModel(worldID string) *worldDataModel {
	return &worldDataModel{
		ID:       worldID,
		Accounts: worldmock.NewAccountMap(),
	}
}

// newWorld creates a new debugging world
func newWorld(dataModel *worldDataModel) (*world, error) {
	blockchainHook := mock.NewMockWorldVM14()
	blockchainHook.AcctMap = dataModel.Accounts

	vm, err := hostCore.NewVMHost(
		blockchainHook,
		getHostParameters(),
	)
	if err != nil {
		return nil, err
	}

	return &world{
		id:             dataModel.ID,
		blockchainHook: blockchainHook,
		vm:             vm,
	}, nil
}

func getHostParameters() *vmhost.VMHostParameters {
	dcdtTransferParser, _ := parsers.NewDCDTTransferParser(worldmock.WorldMarshalizer)
	return &vmhost.VMHostParameters{
		VMType:                   []byte{5, 0},
		BlockGasLimit:            uint64(10000000),
		GasSchedule:              config.MakeGasMap(1, 1),
		ProtectedKeyPrefix:       []byte("E" + "L" + "R" + "O" + "N" + "D"),
		BuiltInFuncContainer:     builtInFunctions.NewBuiltInFunctionContainer(),
		DCDTTransferParser:       dcdtTransferParser,
		EpochNotifier:            &mock.EpochNotifierStub{},
		EnableEpochsHandler:      &mock.EnableEpochsHandlerStub{},
		WasmerSIGSEGVPassthrough: false,
		Hasher:                   worldmock.DefaultHasher,
	}
}

func (w *world) deploySmartContract(request DeployRequest) *DeployResponse {
	input := w.prepareDeployInput(request)
	log.Trace("w.deploySmartContract()", "input", prettyJson(input))

	vmOutput, err := w.vm.RunSmartContractCreate(input)
	if err == nil {
		_ = w.blockchainHook.UpdateAccounts(vmOutput.OutputAccounts, nil)
	}

	response := &DeployResponse{}
	response.ContractResponseBase = createContractResponseBase(&input.VMInput, vmOutput)
	response.Error = err
	response.ContractAddress = w.blockchainHook.LastCreatedContractAddress
	response.ContractAddressHex = toHex(response.ContractAddress)
	return response
}

func (w *world) upgradeSmartContract(request UpgradeRequest) *UpgradeResponse {
	input := w.prepareUpgradeInput(request)
	log.Trace("w.upgradeSmartContract()", "input", prettyJson(input))

	vmOutput, err := w.vm.RunSmartContractCall(input)
	if err == nil {
		_ = w.blockchainHook.UpdateAccounts(vmOutput.OutputAccounts, nil)
	}

	response := &UpgradeResponse{}
	response.ContractResponseBase = createContractResponseBase(&input.VMInput, vmOutput)
	response.Error = err

	return response
}

func (w *world) runSmartContract(request RunRequest) *RunResponse {
	input := w.prepareCallInput(request)
	log.Trace("w.runSmartContract()", "input", prettyJson(input))

	vmOutput, err := w.vm.RunSmartContractCall(input)
	if err == nil {
		_ = w.blockchainHook.UpdateAccounts(vmOutput.OutputAccounts, nil)
	}

	response := &RunResponse{}
	response.ContractResponseBase = createContractResponseBase(&input.VMInput, vmOutput)
	response.Error = err

	return response
}

func (w *world) querySmartContract(request QueryRequest) *QueryResponse {
	input := w.prepareCallInput(request.RunRequest)
	log.Trace("w.querySmartContract()", "input", prettyJson(input))

	vmOutput, err := w.vm.RunSmartContractCall(input)

	response := &QueryResponse{}
	response.ContractResponseBase = createContractResponseBase(&input.VMInput, vmOutput)
	response.Error = err

	return response
}

func (w *world) createAccount(request CreateAccountRequest) *CreateAccountResponse {
	log.Trace("w.createAccount()", "request", prettyJson(request))

	account := worldmock.Account{
		Address:         request.Address,
		Nonce:           request.Nonce,
		Balance:         request.BalanceAsBigInt,
		BalanceDelta:    big.NewInt(0),
		DeveloperReward: big.NewInt(0),
	}
	w.blockchainHook.AcctMap.PutAccount(&account)
	return &CreateAccountResponse{Account: &account}
}

func (w *world) toDataModel() *worldDataModel {
	accounts := w.blockchainHook.AcctMap.Clone()
	for _, account := range accounts {
		account.MockWorld = nil
	}

	return &worldDataModel{
		ID:       w.id,
		Accounts: accounts,
	}
}
