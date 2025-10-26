package types

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	faultTypes "github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/types"
)

var ErrInvalidPrestate = errors.New("absolute prestate does not match")

type GameStatus uint8

const (
	GameStatusInProgress GameStatus = iota
	GameStatusChallengerWon
	GameStatusDefenderWon
)

// String returns the string representation of the game status.
func (s GameStatus) String() string {
	switch s {
	case GameStatusInProgress:
		return "In Progress"
	case GameStatusChallengerWon:
		return "Challenger Won"
	case GameStatusDefenderWon:
		return "Defender Won"
	default:
		return "Unknown"
	}
}

// GameStatusFromUint8 returns a game status from the uint8 representation.
func GameStatusFromUint8(i uint8) (GameStatus, error) {
	if i > 2 {
		return GameStatus(i), fmt.Errorf("invalid game status: %d", i)
	}
	return GameStatus(i), nil
}

type GameMetadata struct {
	Index     uint64
	GameType  faultTypes.GameType
	Timestamp uint64
	Proxy     common.Address
}
