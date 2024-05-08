package bindings

import (
	"fmt"
	"log"
	"strings"

	"github.com/ethereum-optimism/optimism/op-bindings/hardhat"
	"github.com/ethereum-optimism/optimism/op-bindings/solc"
	"github.com/ethereum/go-ethereum/common"
)

// layouts respresents the set of storage layouts. It is populated in an init function.
var layouts = make(map[string]*solc.StorageLayout)

// deployedBytecodes represents the set of deployed bytecodes. It is populated
// in an init function.
var deployedBytecodes = make(map[string]string)

func GetStorageLayout(name string) (*solc.StorageLayout, error) {
	layout := layouts[name]
	if len(layout.Storage) == 0 {
		artifactPath := "/Users/aaron/aaron.lee/OR-Ticket/OR-1550/tokamak-thanos/packages/tokamak/contracts-bedrock/hardhat-artifacts/"
		log.Printf("Loading artifacts from: %s", artifactPath)
		hh, err := hardhat.New("DevnetL1", []string{artifactPath}, []string{})
		if err != nil {
			log.Printf("Failed to create hardhat instance for %s: %v", name, err)
			return nil, fmt.Errorf("failed to create hardhat instance: %v", err)
		}

		layout, err = hh.GetStorageLayout(name)
		if err != nil {
			log.Printf("Failed to find storage layout for %s using hardhat: %v", name, err)
			return nil, fmt.Errorf("failed to find storage layout for %s using hardhat: %v", name, err)
		}

		layouts[name] = layout
		log.Printf("Successfully retrieved and cached storage layout for %s", name)
		log.Printf("Layout data: %+v", layout)
	} else {
		log.Printf("Using cached storage layout for %s", name)
		log.Printf("Cached Layout data: %+v", layout)
	}
	return layout, nil
}

// GetDeployedBytecode returns the deployed bytecode of a contract by name.
func GetDeployedBytecode(name string) ([]byte, error) {
	bc := deployedBytecodes[name]
	if bc == "" {
		return nil, fmt.Errorf("%s: deployed bytecode not found", name)
	}

	if !isHex(bc) {
		return nil, fmt.Errorf("%s: invalid deployed bytecode", name)
	}

	return common.FromHex(bc), nil
}

// isHexCharacter returns bool of c being a valid hexadecimal.
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// isHex validates whether each byte is valid hexadecimal string.
func isHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	str = strings.TrimPrefix(str, "0x")

	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}
