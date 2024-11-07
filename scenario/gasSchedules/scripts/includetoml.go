package main

import (
	"io"
	"os"
	"strings"
)

const suffix = ".toml"

// Reads all .txt files in the current folder
// and encodes them as strings literals in textfiles.go
func main() {
	fs, _ := os.ReadDir(".")
	out, _ := os.Create("gasScheduleEmbedGenerated.go")
	_, _ = out.Write([]byte("package gasschedules \n\n"))
	_, _ = out.Write([]byte("// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n"))
	_, _ = out.Write([]byte("// !!!!!!!!!!!!!!!!!!!!!! AUTO-GENERATED FILE !!!!!!!!!!!!!!!!!!!!!!\n"))
	_, _ = out.Write([]byte("// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n"))
	_, _ = out.Write([]byte("\n"))
	_, _ = out.Write([]byte("// Please do not edit manually!\n"))
	_, _ = out.Write([]byte("// Call `go generate` in `vm-wasm-vm-v1_4/scenarioexec/gasSchedules` to update it.\n"))
	_, _ = out.Write([]byte("\n"))
	_, _ = out.Write([]byte("const (\n"))
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), suffix) {
			_, _ = out.Write([]byte("\t" + strings.TrimSuffix(f.Name(), suffix) + " = `"))
			f, _ := os.Open(f.Name())
			_, _ = io.Copy(out, f)
			_, _ = out.Write([]byte("`\n"))
		}
	}
	_, _ = out.Write([]byte(")\n"))
}
