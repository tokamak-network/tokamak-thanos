package processors

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/tokamak-network/tokamak-thanos/op-supervisor/supervisor/types"
)

// InteropCrossL2InboxAddress is the address of the cross-L2 inbox precompile.
// In geth 1.14+ this lives in params package; we define it locally for old geth compat.
var InteropCrossL2InboxAddress = common.HexToAddress("0x4200000000000000000000000000000000000022")

type EventDecoderFn func(*ethTypes.Log) (*types.ExecutingMessage, error)

func MessageFromLog(l *ethTypes.Log) (*types.Message, error) {
	if l.Address != InteropCrossL2InboxAddress {
		return nil, nil
	}
	if len(l.Topics) != 2 { // topics: event-id and payload-hash
		return nil, nil
	}
	if l.Topics[0] != types.ExecutingMessageEventTopic {
		return nil, nil
	}
	var msg types.Message
	if err := msg.DecodeEvent(l.Topics, l.Data); err != nil {
		return nil, fmt.Errorf("invalid executing message: %w", err)
	}
	return &msg, nil
}

func DecodeExecutingMessageLog(l *ethTypes.Log) (*types.ExecutingMessage, error) {
	msg, err := MessageFromLog(l)
	if err != nil || msg == nil {
		return nil, err
	}
	return &types.ExecutingMessage{
		ChainID:   msg.Identifier.ChainID,
		BlockNum:  msg.Identifier.BlockNumber,
		LogIdx:    msg.Identifier.LogIndex,
		Timestamp: msg.Identifier.Timestamp,
		Checksum:  msg.Checksum(),
	}, nil
}
