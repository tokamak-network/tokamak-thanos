// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"encoding/json"

	"github.com/tokamak-network/tokamak-thanos/op-bindings/solc"
)

const UnsupportedProtocolStorageLayoutJSON = "{\"storage\":[],\"types\":null}"

var UnsupportedProtocolStorageLayout = new(solc.StorageLayout)

var UnsupportedProtocolDeployedBin = "0x60808060405234603157807fea3559ef0000000000000000000000000000000000000000000000000000000060049252fd5b600080fdfea164736f6c6343000811000a"


func init() {
	if err := json.Unmarshal([]byte(UnsupportedProtocolStorageLayoutJSON), UnsupportedProtocolStorageLayout); err != nil {
		panic(err)
	}

	layouts["UnsupportedProtocol"] = UnsupportedProtocolStorageLayout
	deployedBytecodes["UnsupportedProtocol"] = UnsupportedProtocolDeployedBin
	immutableReferences["UnsupportedProtocol"] = false
}
