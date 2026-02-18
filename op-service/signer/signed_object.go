package signer

import "github.com/tokamak-network/tokamak-thanos/op-service/eth"

// SignedObject is an interface for signed block payloads.
type SignedObject interface {
	ID() eth.BlockID
	String() string
	VerifySignature(auth Authenticator) error
}

// Authenticator verifies signatures.
type Authenticator interface {
	Verify(data []byte, sig [65]byte) error
}
