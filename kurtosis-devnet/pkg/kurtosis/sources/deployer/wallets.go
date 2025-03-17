package deployer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	ktfs "github.com/ethereum-optimism/optimism/devnet-sdk/kt/fs"
	"github.com/ethereum-optimism/optimism/op-chain-ops/devkeys"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/yaml.v3"
)

const (
	// TODO: can we figure out how many were actually funded?
	numWallets = 21
)

func getMnemonics(r io.Reader) (string, error) {
	type mnemonicConfig struct {
		Mnemonic string `yaml:"mnemonic"`
		Count    int    `yaml:"count"` // TODO: what does this mean? it seems much larger than the number of wallets
	}

	var config []mnemonicConfig
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(&config); err != nil {
		return "", fmt.Errorf("failed to decode mnemonic config: %w", err)
	}

	// TODO: what does this mean if there are multiple mnemonics in this file?
	return config[0].Mnemonic, nil
}

func (d *Deployer) getL1ValidatorWallets(deployerArtifact *ktfs.Artifact) ([]*Wallet, error) {
	mnemonicsBuffer := bytes.NewBuffer(nil)
	if err := deployerArtifact.ExtractFiles(
		ktfs.NewArtifactFileWriter(d.l1ValidatorMnemonicName, mnemonicsBuffer),
	); err != nil {
		return nil, err
	}

	mnemonic, err := getMnemonics(mnemonicsBuffer)
	if err != nil {
		return nil, err
	}

	m, _ := devkeys.NewMnemonicDevKeys(mnemonic)
	knownWallets := make([]*Wallet, 0)

	var keys []devkeys.Key
	for i := 0; i < numWallets; i++ {
		keys = append(keys, devkeys.UserKey(i))
	}

	for _, key := range keys {
		addr, _ := m.Address(key)
		sec, _ := m.Secret(key)

		knownWallets = append(knownWallets, &Wallet{
			Name:       key.String(),
			Address:    addr,
			PrivateKey: hexutil.Bytes(crypto.FromECDSA(sec)).String(),
		})
	}

	return knownWallets, nil
}

func (d *Deployer) getL1ChainID(genesisArtifact *ktfs.Artifact) (string, error) {
	genesisBuffer := bytes.NewBuffer(nil)
	if err := genesisArtifact.ExtractFiles(
		ktfs.NewArtifactFileWriter(d.l1GenesisName, genesisBuffer),
	); err != nil {
		return "", err
	}

	// Parse the genesis file JSON into a core.Genesis struct
	var genesis core.Genesis
	if err := json.NewDecoder(genesisBuffer).Decode(&genesis); err != nil {
		return "", fmt.Errorf("failed to parse genesis file %s in artifact %s: %w", d.l1GenesisName, d.genesisArtifactName, err)
	}

	return genesis.Config.ChainID.String(), nil
}
