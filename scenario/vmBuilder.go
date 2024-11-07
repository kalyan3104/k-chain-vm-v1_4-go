package scenario

import (
	"fmt"

	"github.com/kalyan3104/k-chain-vm-v1_4-go/config"
	gasSchedules "github.com/kalyan3104/k-chain-vm-v1_4-go/scenario/gasSchedules"
	"github.com/kalyan3104/k-chain-vm-v1_4-go/vmhost"
	"github.com/kalyan3104/k-chain-vm-v1_4-go/vmhost/hostCore"
	"github.com/kalyan3104/k-chain-vm-v1_4-go/vmhost/mock"

	"github.com/kalyan3104/k-chain-core-go/core"
	scenexec "github.com/kalyan3104/k-chain-scenario-go/scenario/executor"
	scenmodel "github.com/kalyan3104/k-chain-scenario-go/scenario/model"
	"github.com/kalyan3104/k-chain-scenario-go/worldmock"
	"github.com/kalyan3104/k-chain-vm-common-go/parsers"
)

var _ scenexec.VMBuilder = (*ScenarioVMHostBuilder)(nil)

// DefaultVMType is the VM type argument we use in tests.
var DefaultVMType = []byte{5, 0}

// VMTestExecutor parses, interprets and executes both .test.json tests and .scen.json scenarios with VM.
type ScenarioVMHostBuilder struct {
	VMType []byte
}

// NewScenarioVMHostBuilder creates a default ScenarioVMHostBuilder.
func NewScenarioVMHostBuilder() *ScenarioVMHostBuilder {
	return &ScenarioVMHostBuilder{
		VMType: DefaultVMType,
	}
}

// NewMockWorld defines how the MockWorld is initialized.
func (*ScenarioVMHostBuilder) NewMockWorld() *worldmock.MockWorld {
	return mock.NewMockWorldVM14()
}

// GasScheduleMapFromScenarios provides the correct gas schedule for the gas schedule named specified in a scenario.
func (svb *ScenarioVMHostBuilder) GasScheduleMapFromScenarios(scenGasSchedule scenmodel.GasSchedule) (worldmock.GasScheduleMap, error) {
	switch scenGasSchedule {
	case scenmodel.GasScheduleDefault:
		return gasSchedules.LoadGasScheduleConfig(gasSchedules.GetV4())
	case scenmodel.GasScheduleDummy:
		return config.MakeGasMapForTests(), nil
	case scenmodel.GasScheduleV3:
		return gasSchedules.LoadGasScheduleConfig(gasSchedules.GetV3())
	case scenmodel.GasScheduleV4:
		return gasSchedules.LoadGasScheduleConfig(gasSchedules.GetV4())
	default:
		return nil, fmt.Errorf("unknown scenario GasSchedule: %d", scenGasSchedule)
	}
}

// GetVMType returns the configured VM type.
func (svb *ScenarioVMHostBuilder) GetVMType() []byte {
	return svb.VMType
}

// NewVM will create a new VM instance with pointers to a mock world and given gas schedule.
func (svb *ScenarioVMHostBuilder) NewVM(
	world *worldmock.MockWorld,
	gasSchedule map[string]map[string]uint64,
) (scenexec.VMInterface, error) {
	blockGasLimit := uint64(10000000)
	dcdtTransferParser, _ := parsers.NewDCDTTransferParser(worldmock.WorldMarshalizer)

	return hostCore.NewVMHost(
		world,
		&vmhost.VMHostParameters{
			VMType:                   svb.VMType,
			BlockGasLimit:            blockGasLimit,
			GasSchedule:              gasSchedule,
			BuiltInFuncContainer:     world.BuiltinFuncs.Container,
			ProtectedKeyPrefix:       []byte(core.ProtectedKeyPrefix),
			DCDTTransferParser:       dcdtTransferParser,
			EpochNotifier:            &mock.EpochNotifierStub{},
			EnableEpochsHandler:      world.EnableEpochsHandler,
			WasmerSIGSEGVPassthrough: false,
			Hasher:                   worldmock.DefaultHasher,
		})

}

// DefaultScenarioExecutor provides a scenario executor with VM 1.5, default configuration
func DefaultScenarioExecutor() *scenexec.ScenarioExecutor {
	return scenexec.NewScenarioExecutor(NewScenarioVMHostBuilder())
}
