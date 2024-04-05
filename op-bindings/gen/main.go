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

// flags struct stores command line arguments.
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

// data struct is used to pass data to the template.
type data struct {
	Name              string
	StorageLayout     string
	DeployedBin       string
	Package           string
	DeployedSourceMap string
}

func main() {
	var f flags
	// Parse command line flags.
	flag.StringVar(&f.ForgeArtifacts, "forge-artifacts", "", "Forge artifacts directory")
	flag.StringVar(&f.HardhatArtifacts, "hardhat-artifacts", "", "Hardhat artifacts directory")
	flag.StringVar(&f.OutDir, "out", "", "Output directory for generated code")
	flag.StringVar(&f.ForgeContracts, "forge-contracts", "", "File path with list of contracts for forge")
	flag.StringVar(&f.HardhatContracts, "hardhat-contracts", "", "File path with list of contracts for hardhat")
	flag.StringVar(&f.SourceMaps, "source-maps", "", "Comma-separated contracts for source maps")
	flag.StringVar(&f.Package, "package", "bindings", "Package name for generated Go code")
	flag.StringVar(&f.MonorepoBase, "monorepo-base", "", "Monorepo base directory for AST paths")
	flag.Parse()

	// Process artifacts for Forge and Hardhat if specified.
	if f.ForgeArtifacts != "" {
		processForgeArtifacts(f)
	}
	if f.HardhatArtifacts != "" {
		processHardhatArtifacts(f)
	}
}

// processForgeArtifacts processes Forge artifacts to generate Go bindings.
func processForgeArtifacts(f flags) {
	// Read the contracts file.
	contractData, err := os.ReadFile(f.ForgeContracts)
	if err != nil {
		log.Fatalf("error reading contracts file: %v", err)
	}
	var contracts []string
	// Unmarshal the JSON data into the contracts slice.
	if err := json.Unmarshal(contractData, &contracts); err != nil {
		log.Fatalf("error parsing contracts JSON: %v", err)
	}

	// Create a set of source maps.
	sourceMaps := strings.Split(f.SourceMaps, ",")
	sourceMapsSet := make(map[string]struct{})
	for _, contractName := range sourceMaps {
		sourceMapsSet[contractName] = struct{}{}
	}

	// Generate Go bindings for each contract.
	for _, contract := range contracts {
		log.Printf("generating code for : %s", contract)
		artifactPath := filepath.Join(f.ForgeArtifacts, contract+".sol", contract+".json")
		forgeArtifactData, err := os.ReadFile(artifactPath)
		if err != nil {
			log.Printf("cannot find forge-artifact for %s, %v", contract, err)
			continue
		}
		log.Printf("using forge-artifact %s\n", artifactPath)

		var artifact foundry.Artifact
		if err := json.Unmarshal(forgeArtifactData, &artifact); err != nil {
			log.Fatalf("Failed to unmarshal artifact for contract %s: %v", contract, err)
		}

		// Write ABI and BIN files.
		abiPath := filepath.Join(f.OutDir, contract+".abi")
		if err := os.WriteFile(abiPath, artifact.Abi, 0644); err != nil {
			log.Fatalf("Failed to write ABI file for %s: %v", contract, err)
		}

		binPath := filepath.Join(f.OutDir, contract+".bin")
		if err := os.WriteFile(binPath, []byte(artifact.Bytecode.Object.String()), 0644); err != nil {
			log.Fatalf("Failed to write BIN file for %s: %v", contract, err)
		}

		// Use abigen to generate Go bindings.
		outFile := filepath.Join(f.OutDir, contract+".go")
		cmd := exec.Command("abigen", "--abi", abiPath, "--bin", binPath, "--pkg", f.Package, "--type", contract, "--out", outFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatalf("abigen failed for contract %s: %v", contract, err)
		}
		log.Printf("wrote file bindings: %s", contract)
	}
}

// processHardhatArtifacts processes Hardhat artifacts similarly to Forge.
func processHardhatArtifacts(f flags) {
	// Initialize Hardhat with the specified network and artifacts.
	hh, err := hardhat.New("mainnet", []string{f.HardhatArtifacts}, []string{})
	if err != nil {
		log.Fatalf("hardhat initialization failed: %v", err)
	}

	// Read and parse the Hardhat contracts file.
	contractNames, err := parseContracts(f.HardhatContracts)
	if err != nil {
		log.Fatalf("hardhat contract file reading error: %v", err)
	}

	// Template preparation for generating Go files.
	t := template.Must(template.New("artifact").Parse(tmpl))

	// Process each contract to generate Go bindings.
	for _, contractName := range contractNames {
		artifactPath := filepath.Join(f.HardhatArtifacts, contractName+".json")
		fmt.Printf("Processing artifact at path: %s\n", artifactPath)
		log.Printf("generating code for : %s", contractName)
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

		// Prepare data for template execution.
		d := data{
			Name:          contractName,
			StorageLayout: serStr,
			DeployedBin:   art.DeployedBytecode.String(),
			Package:       f.Package,
		}

		// Generate Go file from template.
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
		log.Printf("wrote file bindings: %s", contractName)
	}
}

// parseContracts reads and unmarshals a JSON file listing contracts.
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
