package contracts

import (
	"fmt"
	"math/big"

	mock "github.com/kalyan3104/k-chain-vm-v1_4-go/mock/context"
	"github.com/kalyan3104/k-chain-vm-v1_4-go/testcommon"
	"github.com/kalyan3104/k-chain-vm-v1_4-go/vmhost/vmhooks"
)

// TransferAndExecuteFuncName -
var TransferAndExecuteFuncName = "transferAndExecute"

// TransferAndExecuteReturnData -
var TransferAndExecuteReturnData = []byte{1, 2, 3}

// TransferAndExecute is an exposed mock contract method
func TransferAndExecute(instanceMock *mock.InstanceMock, config interface{}) {
	testConfig := config.(TransferAndExecuteTestConfig)
	instanceMock.AddMockMethod(TransferAndExecuteFuncName, func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)

		host.Metering().UseGas(testConfig.GasUsedByParent)

		arguments := host.Runtime().Arguments()
		noOfTransfers := int(big.NewInt(0).SetBytes(arguments[0]).Int64())

		for transfer := 0; transfer < noOfTransfers; transfer++ {
			vmhooks.TransferValueExecuteWithTypedArgs(host,
				GetChildAddressForTransfer(transfer),
				big.NewInt(testConfig.TransferFromParentToChild),
				int64(testConfig.GasTransferToChild),
				big.NewInt(int64(transfer)).Bytes(), // transfer data
				[][]byte{},
			)
		}

		host.Output().Finish(TransferAndExecuteReturnData)

		return instance
	})
}

// GetChildAddressForTransfer -
func GetChildAddressForTransfer(transfer int) []byte {
	return testcommon.MakeTestSCAddress(fmt.Sprintf("childSC-%d", transfer))
}
