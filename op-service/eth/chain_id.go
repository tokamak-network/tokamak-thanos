package eth

import (
	"fmt"
	"golang.org/x/exp/slog"
	"math/big"
	"sort"
	"strings"

	"github.com/holiman/uint256"
)

type ChainID uint256.Int

var _ slog.LogValuer = ChainID(uint256.Int{})

func ChainIDFromBig(chainID *big.Int) ChainID {
	return ChainID(*uint256.MustFromBig(chainID))
}

func ChainIDFromUInt64(i uint64) ChainID {
	return ChainID(*uint256.NewInt(i))
}

func ChainIDFromString(id string) (ChainID, error) {
	if strings.HasPrefix(id, "0x") {
		return ParseHexChainID(id)
	}
	return ParseDecimalChainID(id)
}

func ChainIDFromBytes32(b [32]byte) ChainID {
	val := new(uint256.Int).SetBytes(b[:])
	return ChainID(*val)
}

func ParseDecimalChainID(chainID string) (ChainID, error) {
	v, err := uint256.FromDecimal(chainID)
	if err != nil {
		return ChainID{}, err
	}
	return ChainID(*v), nil
}

func ParseHexChainID(chainID string) (ChainID, error) {
	v, err := uint256.FromHex(chainID)
	if err != nil {
		return ChainID{}, err
	}
	return ChainID(*v), nil
}

func (id ChainID) LogValue() slog.Value {
	return slog.StringValue(id.String())
}

func (id ChainID) String() string {
	return ((*uint256.Int)(&id)).Dec()
}

func (id ChainID) Bytes32() [32]byte {
	return (*uint256.Int)(&id).Bytes32()
}

func (id ChainID) IsUint64() bool {
	return (*uint256.Int)(&id).IsUint64()
}

func (id ChainID) Uint64() (v uint64, ok bool) {
	return (*uint256.Int)(&id).Uint64(), id.IsUint64()
}

func EvilChainIDToUInt64(id ChainID) uint64 {
	v := (*uint256.Int)(&id)
	if !v.IsUint64() {
		panic(fmt.Errorf("ChainID too large for uint64: %v", id))
	}
	return v.Uint64()
}

func (id ChainID) ToBig() *big.Int {
	return (*uint256.Int)(&id).ToBig()
}

func (id ChainID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

func (id *ChainID) UnmarshalText(data []byte) error {
	var x uint256.Int
	err := x.UnmarshalText(data)
	if err != nil {
		return err
	}
	*id = ChainID(x)
	return nil
}

func (id ChainID) Cmp(other ChainID) int {
	return (*uint256.Int)(&id).Cmp((*uint256.Int)(&other))
}

func SortChainID(ids []ChainID) {
	sort.Slice(ids, func(i, j int) bool {
		return ids[i].Cmp(ids[j]) < 0
	})
}
