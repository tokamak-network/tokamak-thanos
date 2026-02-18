package signer

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tokamak-network/tokamak-thanos/op-service/eth"
)

// BlockSigner signs and verifies block payloads for p2p gossip.
type BlockSigner interface {
	SignBlockV1(ctx context.Context, chainID eth.ChainID, payloadHash common.Hash) (sig eth.Bytes65, err error)
}

// OPStackP2PBlockAuthV1 is the authentication context for verifying a signed p2p block.
type OPStackP2PBlockAuthV1 struct {
	Allowed common.Address
	Chain   eth.ChainID
}

// SignedP2PBlock is a raw block payload with its signature.
type SignedP2PBlock struct {
	Raw       []byte
	Signature eth.Bytes65
}

// VerifySignature verifies the block signature against the allowed signer.
func (b *SignedP2PBlock) VerifySignature(auth *OPStackP2PBlockAuthV1) error {
	payloadHash := PayloadHash(b.Raw)
	msg := blockSigningHash(auth.Chain, payloadHash)
	pubKey, err := crypto.SigToPub(msg[:], b.Signature[:])
	if err != nil {
		return fmt.Errorf("failed to recover public key: %w", err)
	}
	addr := crypto.PubkeyToAddress(*pubKey)
	if addr != auth.Allowed {
		return fmt.Errorf("signer %s is not allowed (expected %s)", addr, auth.Allowed)
	}
	return nil
}

// PayloadHash computes the hash of a raw payload for signing.
func PayloadHash(data []byte) common.Hash {
	return crypto.Keccak256Hash(data)
}

// blockSigningHash creates the message to sign for a block.
func blockSigningHash(chainID eth.ChainID, payloadHash common.Hash) common.Hash {
	v, _ := chainID.Uint64()
	domain := crypto.Keccak256Hash([]byte("OPTIMISM_BLOCK_SIGNING"), common.BigToHash(new(big.Int).SetUint64(v)).Bytes())
	return crypto.Keccak256Hash(domain[:], payloadHash[:])
}

// LocalBlockSigner signs blocks with a local private key.
type LocalBlockSigner struct {
	key *ecdsa.PrivateKey
}

var _ BlockSigner = (*LocalBlockSigner)(nil)

// NewLocalBlockSigner creates a new LocalBlockSigner.
func NewLocalBlockSigner(key *ecdsa.PrivateKey) *LocalBlockSigner {
	return &LocalBlockSigner{key: key}
}

func (s *LocalBlockSigner) SignBlockV1(ctx context.Context, chainID eth.ChainID, payloadHash common.Hash) (eth.Bytes65, error) {
	if s.key == nil {
		return eth.Bytes65{}, errors.New("no signing key configured")
	}
	msg := blockSigningHash(chainID, payloadHash)
	sig, err := crypto.Sign(msg[:], s.key)
	if err != nil {
		return eth.Bytes65{}, err
	}
	var out eth.Bytes65
	copy(out[:], sig)
	return out, nil
}

// NewLocalSigner is an alias for NewLocalBlockSigner.
func NewLocalSigner(key *ecdsa.PrivateKey) *LocalBlockSigner {
	return NewLocalBlockSigner(key)
}

// NewRemoteSigner creates a remote signer from CLI config.
// For now this is a stub that returns an error.
func NewRemoteSigner(log interface{}, cfg CLIConfig) (BlockSigner, error) {
	return nil, errors.New("remote signer not yet supported in tokamak-thanos")
}
