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
	hh, err := hardhat.New("mainnet", []string{f.HardhatArtifacts}, []string{})
	if err != nil {
		log.Fatalf("hardhat initialization failed: %v", err)
	}

	contractData, err := os.ReadFile(f.HardhatContracts)
	if err != nil {
		log.Fatalf("error reading hardhat contracts file: %v", err)
	}
	var contracts []string
	if err := json.Unmarshal(contractData, &contracts); err != nil {
		log.Fatalf("error parsing hardhat contracts JSON: %v", err)
	}

	for _, contractName := range contracts {
		artifact, err := hh.GetArtifact(contractName)
		if err != nil {
			log.Printf("Failed to get Hardhat artifact for %s, skipping: %v", contractName, err)
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

		// hexutil.Bytes 타입을 직접 파일에 쓰기 위한 처리.
		// 바이트코드를 문자열로 변환하여 파일에 쓰기
		binPath := filepath.Join(f.OutDir, contractName+".bin")
		if err := os.WriteFile(binPath, []byte(artifact.Bytecode), 0600); err != nil {
			log.Fatalf("Failed to write BIN file for %s: %v", contractName, err)
		}

		outFile := filepath.Join(f.OutDir, contractName+".go")
		cmd := exec.Command("abigen", "--abi", abiPath, "--bin", binPath, "--pkg", f.Package, "--type", contractName, "--out", outFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

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
