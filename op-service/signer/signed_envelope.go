package signer

import (
	"fmt"

	"github.com/tokamak-network/tokamak-thanos/op-service/eth"
)

// SignedExecutionPayloadEnvelope wraps an execution payload envelope with a signature.
type SignedExecutionPayloadEnvelope struct {
	Envelope  *eth.ExecutionPayloadEnvelope `json:"envelope"`
	Signature eth.Bytes65                   `json:"signature"`
}

var _ SignedObject = (*SignedExecutionPayloadEnvelope)(nil)

func (s *SignedExecutionPayloadEnvelope) ID() eth.BlockID {
	return s.Envelope.ExecutionPayload.ID()
}

func (s *SignedExecutionPayloadEnvelope) String() string {
	return fmt.Sprintf("signedEnvelope(%s)", s.ID())
}

func (s *SignedExecutionPayloadEnvelope) VerifySignature(auth Authenticator) error {
	return fmt.Errorf("signature verification not implemented")
}
