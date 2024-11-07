package contracts

import (
	"errors"
	"math/big"

	"github.com/kalyan3104/k-chain-core-go/data/vm"
	mock "github.com/kalyan3104/k-chain-vm-v1_4-go/mock/context"
	test "github.com/kalyan3104/k-chain-vm-v1_4-go/testcommon"
	"github.com/kalyan3104/k-chain-vm-v1_4-go/vmhost"
	"github.com/stretchr/testify/require"
)

func TransferToAsyncParentOnCallbackChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*AsyncCallTestConfig)
	instanceMock.AddMockMethod("transferToThirdParty", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		host.Metering().UseGas(testConfig.GasUsedByChild)

		runtime := host.Runtime()
		output := host.Output()

		vmInput := runtime.GetVMInput()
		scAddress := host.Runtime().GetContextAddress()
		arguments := host.Runtime().Arguments()

		valueToTransfer := big.NewInt(0).SetBytes(arguments[0])

		_ = output.Transfer(
			vmInput.CallerAddr,
			scAddress,
			0,
			0,
			valueToTransfer,
			nil,
			vm.DirectCall)
		return instance
	})
}

// TransferToThirdPartyAsyncChildMock is an exposed mock contract method
func TransferToThirdPartyAsyncChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(*AsyncCallTestConfig)
	instanceMock.AddMockMethod("transferToThirdParty", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T

		host.Metering().UseGas(testConfig.GasUsedByChild)

		arguments := host.Runtime().Arguments()
		outputContext := host.Output()

		if len(arguments) != 3 {
			host.Runtime().SignalUserError("wrong num of arguments")
			return instance
		}

		behavior := byte(0)
		if len(arguments[2]) != 0 {
			behavior = arguments[2][0]
		}
		err := handleChildBehaviorArgument(host, behavior)
		if err != nil {
			return instance
		}

		scAddress := host.Runtime().GetContextAddress()
		valueToTransfer := big.NewInt(0).SetBytes(arguments[0])
		err = outputContext.Transfer(
			test.ThirdPartyAddress,
			scAddress,
			0,
			0,
			valueToTransfer,
			arguments[1],
			0)
		require.Nil(t, err)
		outputContext.Finish([]byte("thirdparty"))

		valueToTransfer = big.NewInt(testConfig.TransferToVault)
		err = outputContext.Transfer(
			test.VaultAddress,
			scAddress,
			0,
			0,
			valueToTransfer,
			[]byte{},
			0)
		require.Nil(t, err)
		outputContext.Finish([]byte("vault"))

		_, _ = host.Storage().SetStorage(test.ChildKey, test.ChildData)

		return instance
	})
}

// ExecutedOnSameContextByCallback is an exposed mock contract method
func ExecutedOnSameContextByCallback(instanceMock *mock.InstanceMock, _ interface{}) {
	instanceMock.AddMockMethod("executedOnSameContextByCallback", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		_, _ = host.Storage().SetStorage(test.ParentKeyB, test.ParentDataA)
		return instance
	})
}

func handleChildBehaviorArgument(host vmhost.VMHost, behavior byte) error {
	if behavior == 1 {
		host.Runtime().SignalUserError("child error")
		return errors.New("behavior / child error")
	}
	if behavior == 2 {
		for {
			host.Output().Finish([]byte("loop"))
		}
	}
	host.Output().Finish([]byte{behavior})
	return nil
}
