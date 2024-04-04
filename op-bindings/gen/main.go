package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ethereum-optimism/optimism/op-bindings/foundry"
	"github.com/ethereum-optimism/optimism/op-bindings/hardhat"
)

type flags struct {
	ForgeArtifacts   string
	HardhatArtifacts string
	ForgeContracts   string
	HardhatContracts string
	SourceMaps       string
	OutDir           string
	Package          string
	MonorepoBase     string
}

type data struct {
	Name              string
	StorageLayout     string
	DeployedBin       string
	Package           string
	DeployedSourceMap string
}

func main() {
	var f flags
	flag.StringVar(&f.ForgeArtifacts, "forge-artifacts", "", "Forge artifacts directory")
	flag.StringVar(&f.HardhatArtifacts, "hardhat-artifacts", "", "Hardhat artifacts directory")
	flag.StringVar(&f.OutDir, "out", "", "Output directory to put generated code in")
	flag.StringVar(&f.ForgeContracts, "forge-contracts", "", "Path to file containing list of contracts to generate bindings for from forge-artifacts")
	flag.StringVar(&f.HardhatContracts, "hardhat-contracts", "", "Path to file containing list of contracts to generate bindings for from hardhat-artifacts")
	flag.StringVar(&f.SourceMaps, "source-maps", "", "Comma-separated list of contracts to generate source maps for")
	flag.StringVar(&f.Package, "package", "bindings", "Package name for the generated Go code")
	flag.StringVar(&f.MonorepoBase, "monorepo-base", "", "Base directory of the monorepo for resolving source file paths in ASTs")
	flag.Parse()

	if f.ForgeArtifacts != "" {
		processForgeArtifacts(f)
	}
	if f.HardhatArtifacts != "" {
		processHardhatArtifacts(f)
	}
}

func processForgeArtifacts(f flags) {
	contractData, err := os.ReadFile(f.ForgeContracts)
	if err != nil {
		log.Fatalf("error reading contracts file: %v", err)
	}
	var contracts []string
	if err := json.Unmarshal(contractData, &contracts); err != nil {
		log.Fatalf("error parsing contracts JSON: %v", err)
	}

	sourceMaps := strings.Split(f.SourceMaps, ",")
	sourceMapsSet := make(map[string]struct{})
	for _, contractName := range sourceMaps {
		sourceMapsSet[contractName] = struct{}{}
	}

	for _, contract := range contracts {
		artifactPath := filepath.Join(f.ForgeArtifacts, contract+".sol", contract+".json")
		forgeArtifactData, err := os.ReadFile(artifactPath)
		if err != nil {
			log.Printf("Failed to read artifact for contract %s, skipping: %v", contract, err)
			continue
		}

		var artifact foundry.Artifact
		if err := json.Unmarshal(forgeArtifactData, &artifact); err != nil {
			log.Fatalf("Failed to unmarshal artifact for contract %s: %v", contract, err)
		}

		abiPath := filepath.Join(f.OutDir, contract+".abi")
		if err := os.WriteFile(abiPath, artifact.Abi, 0644); err != nil {
			log.Fatalf("Failed to write ABI file for %s: %v", contract, err)
		}

		binPath := filepath.Join(f.OutDir, contract+".bin")
		if err := os.WriteFile(binPath, []byte(artifact.Bytecode.Object.String()), 0644); err != nil {
			log.Fatalf("Failed to write BIN file for %s: %v", contract, err)
		}

		outFile := filepath.Join(f.OutDir, contract+".go")
		cmd := exec.Command("abigen", "--abi", abiPath, "--bin", binPath, "--pkg", f.Package, "--type", contract, "--out", outFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatalf("abigen failed for contract %s: %v", contract, err)
		}
		log.Printf("Generated Go binding for contract: %s", contract)
	}
}

func processHardhatArtifacts(f flags) {
	fmt.Printf("Path: %s\n", f.HardhatArtifacts)
	hh, err := hardhat.New("mainnet", []string{f.HardhatArtifacts}, []string{})
	if err != nil {
		log.Fatalf("hardhat 초기화 실패: %v", err)
	}

	contractNames, err := parseContracts(f.HardhatContracts)
	if err != nil {
		log.Fatalf("hardhat 계약 파일 읽기 오류: %v", err)
	}

	t := template.Must(template.New("artifact").Parse(tmpl))

	for _, contractName := range contractNames {
		// artifact, err := hh.GetArtifact(contractName)
		// fmt.Printf("Name: %s\n", contractName)
		// if err != nil {
		// 	log.Printf("%s에 대한 Hardhat 아티팩트 가져오기 실패, 건너뜀: %v", contractName, err)
		// 	continue
		// }
		// abii := artifact.Abi
		// fmt.Printf("ABI: %+v\n", abii)

		// // abi, err := json.Marshal(artifact.Abi)
		// // if err != nil {
		// // 	log.Fatalf("%s에 대한 ABI 마샬링 실패: %v", contractName, err)
		// // }

		// // 바이트코드를 문자열로 변환하지 않고 직접 사용
		// bytecode := artifact.Bytecode

		// abiPath := filepath.Join(f.OutDir, contractName+".abi")
		// if err := os.WriteFile(abiPath, artifact.Abi, 0600); err != nil {
		// 	log.Fatalf("%s에 대한 ABI 파일 쓰기 실패: %v", contractName, err)
		// }

		// binPath := filepath.Join(f.OutDir, contractName+".bin")
		// // 바이트코드를 문자열로 변환하여 파일에 쓰기
		// if err := os.WriteFile(binPath, []byte(bytecode), 0600); err != nil {
		// 	log.Fatalf("%s에 대한 BIN 파일 쓰기 실패: %v", contractName, err)
		// }

		// outFile := filepath.Join(f.OutDir, contractName+".go")
		// cmd := exec.Command("abigen", "--abi", abiPath, "--bin", binPath, "--pkg", f.Package, "--type", contractName, "--out", outFile)
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// if err := cmd.Run(); err != nil {
		// 	log.Fatalf("%s에 대한 Go 바인딩 생성 실패: %v", contractName, err)
		// }
		// log.Printf("계약에 대한 Go 바인딩 생성됨: %s", contractName)
		art, err := hh.GetArtifact(contractName)
		if err != nil {
			log.Fatalf("error reading artifact %s: %v\n", contractName, err)
		}

		storage, err := hh.GetStorageLayout(contractName)
		if err != nil {
			log.Fatalf("error reading storage layout %s: %v\n", contractName, err)
		}

		ser, err := json.Marshal(storage)
		if err != nil {
			log.Fatalf("error marshaling storage: %v\n", err)
		}
		serStr := strings.Replace(string(ser), "\"", "\\\"", -1)

		d := data{
			Name:          contractName,
			StorageLayout: serStr,
			DeployedBin:   art.DeployedBytecode.String(),
			Package:       f.Package,
		}

		fname := filepath.Join(f.OutDir, strings.ToLower(contractName)+"_more.go")
		outfile, err := os.OpenFile(
			fname,
			os.O_RDWR|os.O_CREATE|os.O_TRUNC,
			0o600,
		)
		if err != nil {
			log.Fatalf("error opening %s: %v\n", fname, err)
		}

		if err := t.Execute(outfile, d); err != nil {
			log.Fatalf("error writing template %s: %v", outfile.Name(), err)
		}
		outfile.Close()
		log.Printf("wrote file %s\n", outfile.Name())
	}
}

func parseContracts(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var contracts []string
	if err := json.Unmarshal(data, &contracts); err != nil {
		return nil, err
	}
	return contracts, nil
}

var tmpl = `// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package {{.Package}}

import (
	"encoding/json"

	"github.com/ethereum-optimism/optimism/op-bindings/solc"
)

const {{.Name}}StorageLayoutJSON = "{{.StorageLayout}}"

var {{.Name}}StorageLayout = new(solc.StorageLayout)

var {{.Name}}DeployedBin = "{{.DeployedBin}}"
{{if .DeployedSourceMap}}
var {{.Name}}DeployedSourceMap = "{{.DeployedSourceMap}}"
{{end}}
func init() {
	if err := json.Unmarshal([]byte({{.Name}}StorageLayoutJSON), {{.Name}}StorageLayout); err != nil {
		panic(err)
	}

	layouts["{{.Name}}"] = {{.Name}}StorageLayout
	deployedBytecodes["{{.Name}}"] = {{.Name}}DeployedBin
}
`
