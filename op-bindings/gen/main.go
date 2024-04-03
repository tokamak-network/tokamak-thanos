package main

import (
	"encoding/json"
	"flag"
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
	Contracts        string
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
	flag.StringVar(&f.Contracts, "contracts", "", "Path to file containing list of contracts to generate bindings for")
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
	// contracts.json 파일에서 계약 이름의 목록을 읽어옵니다.
	contractData, err := os.ReadFile(f.Contracts)
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
		log.Printf("Processing Forge artifact for contract: %s\n", contract)
		artifactPath := filepath.Join(f.ForgeArtifacts, contract+".sol", contract+".json")
		forgeArtifactData, err := os.ReadFile(artifactPath)
		if err != nil {
			log.Fatalf("Failed to read artifact for contract %s: %v", contract, err)
		}

		var artifact foundry.Artifact
		if err := json.Unmarshal(forgeArtifactData, &artifact); err != nil {
			log.Fatalf("Failed to unmarshal artifact for contract %s: %v", contract, err)
		}

		// ABI, BIN 파일 및 Go 바인딩 생성 로직은 여기서 처리합니다.
		// 필요한 경우 템플릿을 이용하여 추가 파일을 생성할 수 있습니다.
		// 예시 코드에서는 abigen 커맨드를 사용하는 부분을 유지합니다.

		// ABI 파일 생성
		abiPath := filepath.Join(f.OutDir, contract+".abi")
		if err := os.WriteFile(abiPath, artifact.Abi, 0644); err != nil {
			log.Fatalf("Failed to write ABI file for %s: %v", contract, err)
		}

		// BIN 파일 생성
		binPath := filepath.Join(f.OutDir, contract+".bin")
		if err := os.WriteFile(binPath, []byte(artifact.Bytecode.Object.String()), 0644); err != nil {
			log.Fatalf("Failed to write BIN file for %s: %v", contract, err)
		}

		// Go 바인딩 생성
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
	hh, err := hardhat.New("mainnet", []string{f.HardhatArtifacts}, []string{})
	if err != nil {
		log.Fatalf("hardhat initialization failed: %v", err)
	}

	contracts, err := parseContracts(f.Contracts)
	if err != nil {
		log.Fatalf("error parsing contracts JSON: %v", err)
	}

	for _, contractName := range contracts {
		artifact, err := hh.GetArtifact(contractName)
		if err != nil {
			log.Printf("Failed to get Hardhat artifact for %s: %v", contractName, err)
			continue
		}

		abi, err := json.Marshal(artifact.Abi)
		if err != nil {
			log.Fatalf("Failed to marshal ABI for %s: %v", contractName, err)
		}

		abiPath := filepath.Join(f.OutDir, contractName+".abi")
		if err := os.WriteFile(abiPath, abi, 0600); err != nil {
			log.Fatalf("Failed to write ABI file for %s: %v", contractName, err)
		}

		binPath := filepath.Join(f.OutDir, contractName+".bin")
		if err := os.WriteFile(binPath, artifact.Bytecode, 0600); err != nil {
			log.Fatalf("Failed to write BIN file for %s: %v", contractName, err)
		}

		outFile := filepath.Join(f.OutDir, contractName+".go")
		cmd := exec.Command("abigen", "--abi", abiPath, "--bin", binPath, "--pkg", f.Package, "--type", contractName, "--out", outFile)
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to generate Go binding for %s: %v", contractName, err)
		}
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
