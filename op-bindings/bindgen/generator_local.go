package bindgen

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/ethereum/go-ethereum/log"
	"github.com/tokamak-network/tokamak-thanos/op-bindings/ast"
	"github.com/tokamak-network/tokamak-thanos/op-bindings/foundry"
	"github.com/tokamak-network/tokamak-thanos/op-bindings/hardhat"
)

type BindGenGeneratorLocal struct {
	BindGenGeneratorBase
	SourceMapsList           string
	ForgeArtifactsPath       string
	HardhatArtifactsPath     string
	HardhatContractsListPath string
}

type localContractMetadata struct {
	Name                   string
	StorageLayout          string
	DeployedBin            string
	Package                string
	DeployedSourceMap      string
	HasImmutableReferences bool
}

func (generator *BindGenGeneratorLocal) GenerateBindings() error {
	if generator.ContractsListPath != "" {
		forgeContracts, err := readContractList(generator.Logger, generator.ContractsListPath)
		if err != nil {
			return fmt.Errorf("error reading forge contract list %s: %w", generator.ContractsListPath, err)
		}
		if len(forgeContracts.Local) == 0 {
			return fmt.Errorf("no forge contracts parsed from given contract list: %s", generator.ContractsListPath)
		}
		err = generator.processForgeContracts(forgeContracts.Local)
		if err != nil {
			return err
		}
	}

	if generator.HardhatContractsListPath != "" {
		hardhatContracts, err := readHardhatContractList(generator.Logger, generator.HardhatContractsListPath)
		if err != nil {
			return fmt.Errorf("error reading hardhat contract list %s: %w", generator.HardhatContractsListPath, err)
		}
		if len(hardhatContracts) == 0 {
			return fmt.Errorf("no hardhat contracts parsed from given contract list: %s", generator.HardhatContractsListPath)
		}
		err = generator.processHardhatContracts(hardhatContracts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (generator *BindGenGeneratorLocal) processForgeContracts(forgeContracts []string) error {
	tempArtifactsDir, err := mkTempArtifactsDir(generator.Logger)
	if err != nil {
		return err
	}
	defer func() {
		err := os.RemoveAll(tempArtifactsDir)
		if err != nil {
			generator.Logger.Error("Error removing temporary artifact directory", "path", tempArtifactsDir, "err", err.Error())
		} else {
			generator.Logger.Debug("Successfully removed temporary artifact directory")
		}
	}()

	sourceMapsList := strings.Split(generator.SourceMapsList, ",")
	sourceMapsSet := make(map[string]struct{})
	for _, k := range sourceMapsList {
		sourceMapsSet[k] = struct{}{}
	}

	contractArtifactPaths, err := generator.getContractArtifactPaths()
	if err != nil {
		return err
	}

	contractMetadataFileTemplate := template.Must(template.New("localContractMetadata").Parse(localContractMetadataTemplate))

	for _, contractName := range forgeContracts {
		generator.Logger.Info("Generating bindings and metadata for forge contract", "contract", contractName)

		forgeArtifact, err := generator.readForgeArtifact(contractName, contractArtifactPaths)
		if err != nil {
			return err
		}

		abiFilePath, bytecodeFilePath, err := writeContractArtifacts(generator.Logger, tempArtifactsDir, contractName, forgeArtifact.Abi, []byte(forgeArtifact.Bytecode.Object.String()))
		if err != nil {
			return err
		}

		err = genContractBindings(generator.Logger, generator.MonorepoBasePath, abiFilePath, bytecodeFilePath, generator.BindingsPackageName, contractName)
		if err != nil {
			return err
		}

		deployedSourceMap, canonicalStorageStr, err := generator.canonicalizeStorageLayout(forgeArtifact, sourceMapsSet, contractName)
		if err != nil {
			return err
		}

		re := regexp.MustCompile(`\s+`)
		immutableRefs, err := json.Marshal(re.ReplaceAllString(string(forgeArtifact.DeployedBytecode.ImmutableReferences), ""))
		if err != nil {
			return fmt.Errorf("error marshaling immutable references: %w", err)
		}

		hasImmutables := string(immutableRefs) != `""`

		contractMetaData := localContractMetadata{
			Name:                   contractName,
			StorageLayout:          canonicalStorageStr,
			DeployedBin:            forgeArtifact.DeployedBytecode.Object.String(),
			Package:                generator.BindingsPackageName,
			DeployedSourceMap:      deployedSourceMap,
			HasImmutableReferences: hasImmutables,
		}

		if err := generator.writeContractMetadata(contractMetaData, contractName, contractMetadataFileTemplate); err != nil {
			return err
		}
	}

	return nil
}

func (generator *BindGenGeneratorLocal) processHardhatContracts(hardhatContracts []string) error {
	sourceMapsList := strings.Split(generator.SourceMapsList, ",")
	sourceMapsSet := make(map[string]struct{})
	for _, k := range sourceMapsList {
		sourceMapsSet[k] = struct{}{}
	}

	contractMetadataFileTemplate := template.Must(template.New("localContractMetadata").Parse(localContractMetadataTemplate))

	for _, contractName := range hardhatContracts {
		generator.Logger.Info("Generating bindings and metadata for hardhat contract", "contract", contractName)

		err := generator.processHardhatArtifact(contractName, contractMetadataFileTemplate, sourceMapsSet)
		if err != nil {
			return err
		}
	}

	return nil
}

func (generator *BindGenGeneratorLocal) processHardhatArtifact(contractName string, contractMetadataFileTemplate *template.Template, sourceMapsSet map[string]struct{}) error {
	if generator.HardhatArtifactsPath == "" {
		generator.Logger.Warn("Skipping Hardhat artifact processing as no path is provided")
		return nil
	}

	hh, err := hardhat.New("mainnet", []string{generator.HardhatArtifactsPath}, []string{})
	if err != nil {
		return fmt.Errorf("hardhat initialization failed: %w", err)
	}

	art, err := hh.GetArtifact(contractName)
	if err != nil {
		return fmt.Errorf("error reading artifact %s: %w", contractName, err)
	}

	storage, err := hh.GetStorageLayout(contractName)
	if err != nil {
		return fmt.Errorf("error reading storage layout %s: %w", contractName, err)
	}

	ser, err := json.Marshal(storage)
	if err != nil {
		return fmt.Errorf("error marshaling storage: %w", err)
	}
	serStr := strings.Replace(string(ser), "\"", "\\\"", -1)

	deployedSourceMap := ""
	deployedBin := ""

	// Convert art.DeployedBytecode to DeployedBytecodeObject if it is not a string
	switch v := art.DeployedBytecode.(type) {
	case hardhat.DeployedBytecodeObject:
		deployedSourceMap = v.SourceMap
		deployedBin = v.Object.String()
	case string:
		deployedBin = v
	}

	// GetBuildInfo 함수를 사용하여 immutableReferences를 가져옴
	buildInfo, err := hh.GetBuildInfo(contractName)
	if err != nil {
		return fmt.Errorf("error getting build info for %s: %w", contractName, err)
	}

	hasImmutables := false
	for key, value := range buildInfo.Output.Contracts {
		if strings.Contains(key, fmt.Sprintf("/%s.sol", contractName)) && !strings.Contains(key, "/interfaces/") {
			for _, v := range value {
				if v.Evm.DeployedBytecode.ImmutableReferences != nil && len(v.Evm.DeployedBytecode.ImmutableReferences) > 0 {
					fmt.Printf("key: %s, keyimmutable: %+v\n", key, v.Evm.DeployedBytecode.ImmutableReferences)
					hasImmutables = true
					break
				}
			}
		}
	}

	generator.Logger.Debug("ImmutableReferences found", "hasImmutables", hasImmutables)

	contractMetaData := localContractMetadata{
		Name:                   contractName,
		StorageLayout:          serStr,
		DeployedBin:            deployedBin,
		Package:                generator.BindingsPackageName,
		DeployedSourceMap:      deployedSourceMap,
		HasImmutableReferences: hasImmutables,
	}

	return generator.writeContractMetadata(contractMetaData, contractName, contractMetadataFileTemplate)
}

func (generator *BindGenGeneratorLocal) getContractArtifactPaths() (map[string]string, error) {
	artifactPaths := make(map[string]string)
	if err := filepath.Walk(generator.ForgeArtifactsPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.HasSuffix(path, ".json") {
				base := filepath.Base(path)
				name := strings.TrimSuffix(base, ".json")

				// remove the compiler version from the name
				re := regexp.MustCompile(`\.\d+\.\d+\.\d+`)
				sanitized := re.ReplaceAllString(name, "")
				_, ok := artifactPaths[sanitized]
				if !ok {
					artifactPaths[sanitized] = path
				} else {
					generator.Logger.Warn("Multiple versions of forge artifacts exist, using lesser version", "contract", sanitized)
				}
			}
			return nil
		}); err != nil {
		return artifactPaths, err
	}

	for contract, path := range artifactPaths {
		generator.Logger.Debug("Found artifact", "contract", contract, "path", path)
	}

	return artifactPaths, nil
}

func (generator *BindGenGeneratorLocal) readForgeArtifact(contractName string, contractArtifactPaths map[string]string) (foundry.Artifact, error) {
	var forgeArtifact foundry.Artifact

	contractArtifactPath := path.Join(generator.ForgeArtifactsPath, contractName+".sol", contractName+".json")
	generator.Logger.Debug("Attempting to read forge artifact", "path", contractArtifactPath)
	forgeArtifactRaw, err := os.ReadFile(contractArtifactPath)
	if errors.Is(err, os.ErrNotExist) {
		generator.Logger.Debug("Cannot find forge-artifact at standard path, trying provided path", "contract", contractName, "standardPath", contractArtifactPath, "providedPath", contractArtifactPaths[contractName])
		contractArtifactPath = contractArtifactPaths[contractName]
		forgeArtifactRaw, err = os.ReadFile(contractArtifactPath)
		if errors.Is(err, os.ErrNotExist) {
			return forgeArtifact, fmt.Errorf("cannot find forge-artifact of %q at %q or %q", contractName, contractArtifactPath, contractArtifactPaths[contractName])
		}
	}

	generator.Logger.Debug("Using forge-artifact", "path", contractArtifactPath)
	if err := json.Unmarshal(forgeArtifactRaw, &forgeArtifact); err != nil {
		return forgeArtifact, fmt.Errorf("failed to parse forge artifact of %q: %w", contractName, err)
	}

	return forgeArtifact, nil
}

func (generator *BindGenGeneratorLocal) canonicalizeStorageLayout(forgeArtifact foundry.Artifact, sourceMapsSet map[string]struct{}, contractName string) (string, string, error) {
	artifactStorageStruct := forgeArtifact.StorageLayout
	canonicalStorageStruct := ast.CanonicalizeASTIDs(&artifactStorageStruct, generator.MonorepoBasePath)
	canonicalStorageJson, err := json.Marshal(canonicalStorageStruct)
	if err != nil {
		return "", "", fmt.Errorf("error marshaling canonical storage: %w", err)
	}
	canonicalStorageStr := strings.Replace(string(canonicalStorageJson), "\"", "\\\"", -1)

	deployedSourceMap := ""
	if _, ok := sourceMapsSet[contractName]; ok {
		deployedSourceMap = forgeArtifact.DeployedBytecode.SourceMap
	}

	return deployedSourceMap, canonicalStorageStr, nil
}

func (generator *BindGenGeneratorLocal) writeContractMetadata(contractMetaData localContractMetadata, contractName string, fileTemplate *template.Template) error {
	metadataFilePath := filepath.Join(generator.MetadataOut, strings.ToLower(contractName)+"_more.go")
	metadataFile, err := os.OpenFile(
		metadataFilePath,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		0o600,
	)
	if err != nil {
		return fmt.Errorf("error opening %s's metadata file at %s: %w", contractName, metadataFilePath, err)
	}
	defer metadataFile.Close()

	if err := fileTemplate.Execute(metadataFile, contractMetaData); err != nil {
		return fmt.Errorf("error writing %s's contract metadata at %s: %w", contractName, metadataFilePath, err)
	}

	generator.Logger.Debug("Successfully wrote contract metadata", "contract", contractName, "path", metadataFilePath)
	return nil
}

func readHardhatContractList(logger log.Logger, path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening contract list file %s: %w", path, err)
	}
	defer file.Close()

	var contracts []string
	if err := json.NewDecoder(file).Decode(&contracts); err != nil {
		return nil, fmt.Errorf("error decoding contract list file %s: %w", path, err)
	}

	logger.Debug("Loaded contract list", "path", path, "count", len(contracts))
	return contracts, nil
}

var localContractMetadataTemplate = `// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package {{.Package}}

import (
	"encoding/json"

	"github.com/tokamak-network/tokamak-thanos/op-bindings/solc"
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
	immutableReferences["{{.Name}}"] = {{.HasImmutableReferences}}
}
`
