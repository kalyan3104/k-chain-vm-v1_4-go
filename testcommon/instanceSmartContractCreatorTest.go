package testcommon

import (
	"testing"

	vmcommon "github.com/kalyan3104/k-chain-vm-common-go"
	contextmock "github.com/kalyan3104/k-chain-vm-v1_4-go/mock/context"
	"github.com/kalyan3104/k-chain-vm-v1_4-go/vmhost"
)

// TestCreateTemplateConfig holds the data to build a contract creation test
type TestCreateTemplateConfig struct {
	t                  *testing.T
	address            []byte
	input              *vmcommon.ContractCreateInput
	setup              func(vmhost.VMHost, *contextmock.BlockchainHookStub)
	assertResults      func(*contextmock.BlockchainHookStub, *VMOutputVerifier)
	host               vmhost.VMHost
	blockchainHookStub *contextmock.BlockchainHookStub
}

// BuildInstanceCreatorTest starts the building process for a contract creation test
func BuildInstanceCreatorTest(t *testing.T) *TestCreateTemplateConfig {
	return &TestCreateTemplateConfig{
		t:     t,
		setup: func(vmhost.VMHost, *contextmock.BlockchainHookStub) {},
	}
}

// WithInput provides the ContractCreateInput for a TestCreateTemplateConfig
func (callerTest *TestCreateTemplateConfig) WithInput(input *vmcommon.ContractCreateInput) *TestCreateTemplateConfig {
	callerTest.input = input
	return callerTest
}

// WithAddress provides the address for a TestCreateTemplateConfig
func (callerTest *TestCreateTemplateConfig) WithAddress(address []byte) *TestCreateTemplateConfig {
	callerTest.address = address
	return callerTest
}

// WithSetup provides the setup function for a TestCreateTemplateConfig
func (callerTest *TestCreateTemplateConfig) WithSetup(setup func(vmhost.VMHost, *contextmock.BlockchainHookStub)) *TestCreateTemplateConfig {
	callerTest.setup = setup
	return callerTest
}

// AndAssertResults provides the function that will aserts the results
func (callerTest *TestCreateTemplateConfig) AndAssertResults(assertResults func(*contextmock.BlockchainHookStub, *VMOutputVerifier)) {
	callerTest.assertResults = assertResults
	callerTest.runTest(true)
}

// AndAssertResultsWithoutReset provides the function that will aserts the results
func (callerTest *TestCreateTemplateConfig) AndAssertResultsWithoutReset(assertResults func(*contextmock.BlockchainHookStub, *VMOutputVerifier)) {
	callerTest.assertResults = assertResults
	callerTest.runTest(false)
}

func (callerTest *TestCreateTemplateConfig) runTest(reset bool) {
	if callerTest.host == nil {
		callerTest.host, callerTest.blockchainHookStub = DefaultTestVMForDeployment(callerTest.t, 24, callerTest.address)
		callerTest.setup(callerTest.host, callerTest.blockchainHookStub)
	}
	defer func() {
		if reset {
			callerTest.host.Reset()
		}
	}()

	vmOutput, err := callerTest.host.RunSmartContractCreate(callerTest.input)

	verify := NewVMOutputVerifier(callerTest.t, vmOutput, err)
	callerTest.assertResults(callerTest.blockchainHookStub, verify)
}
