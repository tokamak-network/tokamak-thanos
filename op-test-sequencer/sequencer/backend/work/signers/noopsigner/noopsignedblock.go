package noopsigner

import (
	"github.com/tokamak-network/tokamak-thanos/op-service/signer"
	"github.com/tokamak-network/tokamak-thanos/op-test-sequencer/sequencer/backend/work"
)

type NoopSignedBlock struct {
	work.Block
}

func (s *NoopSignedBlock) VerifySignature(_ signer.Authenticator) error {
	return nil
}
