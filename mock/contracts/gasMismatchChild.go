package contracts

import (
	"math/big"

	mock "github.com/kalyan3104/k-chain-vm-v1_4-go/mock/context"
)

// GasMismatchAsyncCallChildMock is an exposed mock contract method
func GasMismatchAsyncCallChildMock(instanceMock *mock.InstanceMock, _ interface{}) {
	instanceMock.AddMockMethod("gasMismatchChild", func() *mock.InstanceMock {
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		host.Output().Finish(big.NewInt(42).Bytes())
		return instance
	})
}
