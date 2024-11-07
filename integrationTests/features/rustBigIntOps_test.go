package featuresintegrationtest

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	vmcommon "github.com/kalyan3104/k-chain-vm-common-go"
)

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	vmTestRoot := filepath.Join(exePath, "../../test")
	return vmTestRoot
}

func getFeaturesContractPath() string {
	return filepath.Join(getTestRoot(), "features/basic-features/output/basic-features.wasm")
}

func unsignedInterpreter(bytes []byte) *big.Int {
	return big.NewInt(0).SetBytes(bytes)
}

func appendBinaryOpTestCase(
	testCases []*pureFunctionIO,
	opName string, signed bool,
	arg1, arg2, result []byte,
	expectedStatus vmcommon.ReturnCode, expectedMessage string,
) []*pureFunctionIO {

	var typeName string
	if signed {
		typeName = "big_int"
	} else {
		typeName = "big_uint"
	}

	expectedResults := make([][]byte, 0)
	if expectedStatus == vmcommon.Ok {
		expectedResults = [][]byte{result}
	}

	testCases = append(testCases, &pureFunctionIO{
		functionName:    fmt.Sprintf("%s_%s", opName, typeName),
		arguments:       [][]byte{arg1, arg2},
		expectedStatus:  expectedStatus,
		expectedMessage: expectedMessage,
		expectedResults: expectedResults,
	})

	testCases = append(testCases, &pureFunctionIO{
		functionName:    fmt.Sprintf("%s_%s_ref", opName, typeName),
		arguments:       [][]byte{arg1, arg2},
		expectedStatus:  expectedStatus,
		expectedMessage: expectedMessage,
		expectedResults: expectedResults,
	})

	testCases = append(testCases, &pureFunctionIO{
		functionName:    fmt.Sprintf("%s_assign_%s", opName, typeName),
		arguments:       [][]byte{arg1, arg2},
		expectedStatus:  expectedStatus,
		expectedMessage: expectedMessage,
		expectedResults: expectedResults,
	})

	testCases = append(testCases, &pureFunctionIO{
		functionName:    fmt.Sprintf("%s_assign_%s_ref", opName, typeName),
		arguments:       [][]byte{arg1, arg2},
		expectedStatus:  expectedStatus,
		expectedMessage: expectedMessage,
		expectedResults: expectedResults,
	})

	return testCases
}

// func TestBigIntArith(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("long test")
// 	}

// 	var testCases []*pureFunctionIO

// 	big1, _ := big.NewInt(0).SetString("18446744073709551616", 10)
// 	big2, _ := big.NewInt(0).SetString("-123456789012345678901234567890", 10)
// 	numbers := []*big.Int{
// 		big.NewInt(0),
// 		big.NewInt(1),
// 		big.NewInt(-1),
// 		big.NewInt(12345),
// 		big1,
// 		big2,
// 	}

// 	for _, num1 := range numbers {
// 		for _, num2 := range numbers {
// 			bytes1 := twos.ToBytes(num1)
// 			bytes2 := twos.ToBytes(num2)

// 			// add
// 			sumBytes := twos.ToBytes(big.NewInt(0).Add(num1, num2))
// 			testCases = appendBinaryOpTestCase(testCases,
// 				"add", true,
// 				bytes1, bytes2, sumBytes,
// 				vmcommon.Ok, "")

// 			// sub
// 			diffBytes := twos.ToBytes(big.NewInt(0).Sub(num1, num2))
// 			testCases = appendBinaryOpTestCase(testCases,
// 				"sub", true,
// 				bytes1, bytes2, diffBytes,
// 				vmcommon.Ok, "")

// 			// mul
// 			mulBytes := twos.ToBytes(big.NewInt(0).Mul(num1, num2))
// 			testCases = appendBinaryOpTestCase(testCases,
// 				"mul", true,
// 				bytes1, bytes2, mulBytes,
// 				vmcommon.Ok, "")

// 			// div
// 			if num2.Sign() == 0 {
// 				testCases = appendBinaryOpTestCase(testCases,
// 					"div", true,
// 					bytes1, bytes2, nil,
// 					vmcommon.ExecutionFailed, vmhost.ErrDivZero.Error())
// 			} else {
// 				divBytes := twos.ToBytes(big.NewInt(0).Quo(num1, num2))
// 				testCases = appendBinaryOpTestCase(testCases,
// 					"div", true,
// 					bytes1, bytes2, divBytes,
// 					vmcommon.Ok, "")
// 			}

// 			// mod
// 			if num2.Sign() == 0 {
// 				testCases = appendBinaryOpTestCase(testCases,
// 					"rem", true,
// 					bytes1, bytes2, nil,
// 					vmcommon.ExecutionFailed, vmhost.ErrDivZero.Error())
// 			} else {
// 				remBytes := twos.ToBytes(big.NewInt(0).Rem(num1, num2))
// 				testCases = appendBinaryOpTestCase(testCases,
// 					"rem", true,
// 					bytes1, bytes2, remBytes,
// 					vmcommon.Ok, "")
// 			}
// 		}
// 	}

// 	logFunc := func(testCaseIndex, testCaseCount int) {
// 		if testCaseIndex%100 == 0 {
// 			fmt.Printf("Big int operation test case %d/%d\n", testCaseIndex, len(testCases))
// 		}
// 	}

// 	pfe, err := newPureFunctionExecutor()
// 	require.Nil(t, err)
// 	defer func() {
// 		vmHost := pfe.vm.(vmhost.VMHost)
// 		vmHost.Reset()
// 	}()

// 	pfe.initAccounts(getFeaturesContractPath())
// 	pfe.executePureFunctionTests(t, testCases, unsignedInterpreter, logFunc)
// }

// func TestBigUintArith(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("long test")
// 	}

// 	var testCases []*pureFunctionIO

// 	big1, _ := big.NewInt(0).SetString("18446744073709551615", 10)
// 	big2, _ := big.NewInt(0).SetString("18446744073709551616", 10)
// 	numbers := []*big.Int{
// 		big.NewInt(0),
// 		big.NewInt(1),
// 		big.NewInt(2),
// 		big.NewInt(12345),
// 		big1,
// 		big2,
// 	}

// 	for _, num1 := range numbers {
// 		for _, num2 := range numbers {
// 			bytes1 := num1.Bytes()
// 			bytes2 := num2.Bytes()

// 			// add
// 			sumBytes := big.NewInt(0).Add(num1, num2).Bytes()
// 			testCases = appendBinaryOpTestCase(testCases,
// 				"add", false,
// 				bytes1, bytes2, sumBytes,
// 				vmcommon.Ok, "")

// 			// sub
// 			diff := big.NewInt(0).Sub(num1, num2)
// 			if diff.Sign() < 0 {
// 				testCases = appendBinaryOpTestCase(testCases,
// 					"sub", false,
// 					bytes1, bytes2, nil,
// 					vmcommon.UserError, "cannot subtract because result would be negative")
// 			} else {
// 				testCases = appendBinaryOpTestCase(testCases,
// 					"sub", false,
// 					bytes1, bytes2, diff.Bytes(),
// 					vmcommon.Ok, "")
// 			}

// 			// mul
// 			mulBytes := big.NewInt(0).Mul(num1, num2).Bytes()
// 			testCases = appendBinaryOpTestCase(testCases,
// 				"mul", false,
// 				bytes1, bytes2, mulBytes,
// 				vmcommon.Ok, "")

// 			// div
// 			if num2.Sign() == 0 {
// 				testCases = appendBinaryOpTestCase(testCases,
// 					"div", false,
// 					bytes1, bytes2, nil,
// 					vmcommon.ExecutionFailed, vmhost.ErrDivZero.Error())
// 			} else {
// 				divBytes := big.NewInt(0).Quo(num1, num2).Bytes()
// 				testCases = appendBinaryOpTestCase(testCases,
// 					"div", false,
// 					bytes1, bytes2, divBytes,
// 					vmcommon.Ok, "")
// 			}

// 			// mod
// 			if num2.Sign() == 0 {
// 				testCases = appendBinaryOpTestCase(testCases,
// 					"rem", false,
// 					bytes1, bytes2, nil,
// 					vmcommon.ExecutionFailed, vmhost.ErrDivZero.Error())
// 			} else {
// 				remBytes := big.NewInt(0).Rem(num1, num2).Bytes()
// 				testCases = appendBinaryOpTestCase(testCases,
// 					"rem", false,
// 					bytes1, bytes2, remBytes,
// 					vmcommon.Ok, "")
// 			}
// 		}
// 	}

// 	logFunc := func(testCaseIndex, testCaseCount int) {
// 		if testCaseIndex%100 == 0 {
// 			fmt.Printf("Big uint operation test case %d/%d\n", testCaseIndex, len(testCases))
// 		}
// 	}

// 	pfe, err := newPureFunctionExecutor()
// 	require.Nil(t, err)
// 	defer func() {
// 		vmHost := pfe.vm.(vmhost.VMHost)
// 		vmHost.Reset()
// 	}()

// 	pfe.initAccounts(getFeaturesContractPath())
// 	pfe.executePureFunctionTests(t, testCases, unsignedInterpreter, logFunc)
// }

// func TestBigUintBitwise(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("long test")
// 	}

// 	var testCases []*pureFunctionIO

// 	big1, _ := big.NewInt(0).SetString("18446744073709551615", 10)
// 	big2, _ := big.NewInt(0).SetString("123456789012345678901234567890", 10)
// 	numbers := []*big.Int{
// 		big.NewInt(0),
// 		big.NewInt(1),
// 		big.NewInt(2),
// 		big.NewInt(12345),
// 		big1,
// 		big2,
// 	}

// 	for _, num1 := range numbers {
// 		for _, num2 := range numbers {
// 			bytes1 := num1.Bytes()
// 			bytes2 := num2.Bytes()

// 			// and
// 			sumBytes := big.NewInt(0).And(num1, num2).Bytes()
// 			testCases = appendBinaryOpTestCase(testCases,
// 				"bit_and", false,
// 				bytes1, bytes2, sumBytes,
// 				vmcommon.Ok, "")

// 			// or
// 			orBytes := big.NewInt(0).Or(num1, num2).Bytes()
// 			testCases = appendBinaryOpTestCase(testCases,
// 				"bit_or", false,
// 				bytes1, bytes2, orBytes,
// 				vmcommon.Ok, "")

// 			// xor
// 			xorBytes := big.NewInt(0).Xor(num1, num2).Bytes()
// 			testCases = appendBinaryOpTestCase(testCases,
// 				"bit_xor", false,
// 				bytes1, bytes2, xorBytes,
// 				vmcommon.Ok, "")
// 		}
// 	}

// 	logFunc := func(testCaseIndex, testCaseCount int) {
// 		if testCaseIndex%100 == 0 {
// 			fmt.Printf("Big uint bitwise operation test case %d/%d\n", testCaseIndex, len(testCases))
// 		}
// 	}

// 	pfe, err := newPureFunctionExecutor()
// 	require.Nil(t, err)
// 	defer func() {
// 		vmHost := pfe.vm.(vmhost.VMHost)
// 		vmHost.Reset()
// 	}()

// 	pfe.initAccounts(getFeaturesContractPath())
// 	pfe.executePureFunctionTests(t, testCases, unsignedInterpreter, logFunc)
// }

// func TestBigUintShift(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("long test")
// 	}

// 	var testCases []*pureFunctionIO

// 	big1, _ := big.NewInt(0).SetString("18446744073709551615", 10)
// 	big2, _ := big.NewInt(0).SetString("123456789012345678901234567890", 10)
// 	numbers := []*big.Int{
// 		big.NewInt(0),
// 		big.NewInt(1),
// 		big1,
// 		big2,
// 	}

// 	shiftAmounts := []uint{
// 		0,
// 		1,
// 		10,
// 		100,
// 	}

// 	for _, num := range numbers {
// 		for _, shiftAmount := range shiftAmounts {
// 			bytes1 := num.Bytes()
// 			bytes2 := big.NewInt(int64(shiftAmount)).Bytes()

// 			// shift right
// 			shrBytes := big.NewInt(0).Rsh(num, shiftAmount).Bytes()
// 			testCases = appendBinaryOpTestCase(testCases,
// 				"shr", false,
// 				bytes1, bytes2, shrBytes,
// 				vmcommon.Ok, "")

// 			// shift left
// 			shlBytes := big.NewInt(0).Lsh(num, shiftAmount).Bytes()
// 			testCases = appendBinaryOpTestCase(testCases,
// 				"shl", false,
// 				bytes1, bytes2, shlBytes,
// 				vmcommon.Ok, "")
// 		}
// 	}

// 	logFunc := func(testCaseIndex, testCaseCount int) {
// 		if testCaseIndex%100 == 0 {
// 			fmt.Printf("Big uint bitwise shift test case %d/%d\n", testCaseIndex, len(testCases))
// 		}
// 	}

// 	pfe, err := newPureFunctionExecutor()
// 	require.Nil(t, err)
// 	defer func() {
// 		vmHost := pfe.vm.(vmhost.VMHost)
// 		vmHost.Reset()
// 	}()

// 	pfe.initAccounts(getFeaturesContractPath())
// 	pfe.executePureFunctionTests(t, testCases, unsignedInterpreter, logFunc)
// }
