// Package eip1559 provides Holocene/Jovian EIP-1559 parameter encoding for Optimism.
// This is a shim for github.com/ethereum/go-ethereum/consensus/misc/eip1559 which only
// exists in op-geth. The tokamak-thanos-geth fork doesn't include these functions.
package eip1559

import (
	"encoding/binary"
	"errors"
	"fmt"
	gomath "math"
)

const (
	HoloceneExtraDataVersionByte = uint8(0x00)
	JovianExtraDataVersionByte   = uint8(0x01)
)

type ForkChecker interface {
	IsHolocene(time uint64) bool
	IsJovian(time uint64) bool
}

func ValidateOptimismExtraData(fc ForkChecker, time uint64, extraData []byte) error {
	if fc.IsJovian(time) {
		return ValidateJovianExtraData(extraData)
	} else if fc.IsHolocene(time) {
		return ValidateHoloceneExtraData(extraData)
	} else if len(extraData) > 0 {
		return errors.New("extraData must be empty before Holocene")
	}
	return nil
}

func DecodeOptimismExtraData(fc ForkChecker, time uint64, extraData []byte) (uint64, uint64, *uint64) {
	if fc.IsJovian(time) {
		return DecodeJovianExtraData(extraData)
	} else if fc.IsHolocene(time) {
		d, e := DecodeHoloceneExtraData(extraData)
		return d, e, nil
	}
	return 0, 0, nil
}

func EncodeOptimismExtraData(fc ForkChecker, time uint64, denominator, elasticity uint64, minBaseFee *uint64) []byte {
	if fc.IsJovian(time) {
		if minBaseFee == nil {
			panic("minBaseFee cannot be nil since the MinBaseFee feature is enabled")
		}
		return EncodeJovianExtraData(denominator, elasticity, *minBaseFee)
	} else if fc.IsHolocene(time) {
		return EncodeHoloceneExtraData(denominator, elasticity)
	}
	return nil
}

func DecodeHolocene1559Params(params []byte) (uint64, uint64) {
	if len(params) != 8 {
		return 0, 0
	}
	return uint64(binary.BigEndian.Uint32(params[:4])), uint64(binary.BigEndian.Uint32(params[4:]))
}

func DecodeHoloceneExtraData(extra []byte) (uint64, uint64) {
	if len(extra) != 9 {
		return 0, 0
	}
	return DecodeHolocene1559Params(extra[1:])
}

func EncodeHolocene1559Params(denom, elasticity uint64) []byte {
	r := make([]byte, 8)
	if denom > gomath.MaxUint32 || elasticity > gomath.MaxUint32 {
		panic("eip-1559 parameters out of uint32 range")
	}
	binary.BigEndian.PutUint32(r[:4], uint32(denom))
	binary.BigEndian.PutUint32(r[4:], uint32(elasticity))
	return r
}

func EncodeHoloceneExtraData(denom, elasticity uint64) []byte {
	r := make([]byte, 9)
	if denom > gomath.MaxUint32 || elasticity > gomath.MaxUint32 {
		panic("eip-1559 parameters out of uint32 range")
	}
	binary.BigEndian.PutUint32(r[1:5], uint32(denom))
	binary.BigEndian.PutUint32(r[5:], uint32(elasticity))
	return r
}

func ValidateHolocene1559Params(params []byte) error {
	if len(params) != 8 {
		return fmt.Errorf("holocene eip-1559 params should be 8 bytes, got %d", len(params))
	}
	d, e := DecodeHolocene1559Params(params)
	if e != 0 && d == 0 {
		return errors.New("holocene params cannot have a 0 denominator unless elasticity is also 0")
	} else if e == 0 && d != 0 {
		return errors.New("holocene params cannot have a 0 elasticity unless denominator is also 0")
	}
	return nil
}

func ValidateHoloceneExtraData(extra []byte) error {
	if len(extra) != 9 {
		return fmt.Errorf("holocene extraData should be 9 bytes, got %d", len(extra))
	}
	if extra[0] != HoloceneExtraDataVersionByte {
		return fmt.Errorf("holocene extraData version byte should be %d, got %d", HoloceneExtraDataVersionByte, extra[0])
	}
	d, e := DecodeHolocene1559Params(extra[1:])
	if d == 0 {
		return errors.New("holocene extraData must encode a non-zero denominator")
	} else if e == 0 {
		return errors.New("holocene extraData must encode a non-zero elasticity")
	}
	return nil
}

func DecodeJovianExtraData(extra []byte) (uint64, uint64, *uint64) {
	if len(extra) == 9 {
		d, e := DecodeHolocene1559Params(extra[1:9])
		return d, e, nil
	} else if len(extra) == 17 {
		d, e := DecodeHolocene1559Params(extra[1:9])
		minBaseFee := binary.BigEndian.Uint64(extra[9:])
		return d, e, &minBaseFee
	}
	return 0, 0, nil
}

func EncodeJovianExtraData(denom, elasticity, minBaseFee uint64) []byte {
	r := make([]byte, 17)
	if denom > gomath.MaxUint32 || elasticity > gomath.MaxUint32 {
		panic("eip-1559 parameters out of uint32 range")
	}
	r[0] = JovianExtraDataVersionByte
	binary.BigEndian.PutUint32(r[1:5], uint32(denom))
	binary.BigEndian.PutUint32(r[5:9], uint32(elasticity))
	binary.BigEndian.PutUint64(r[9:], minBaseFee)
	return r
}

func ValidateJovianExtraData(extra []byte) error {
	if len(extra) != 17 {
		return fmt.Errorf("Jovian extraData should be 17 bytes, got %d", len(extra))
	}
	if extra[0] != JovianExtraDataVersionByte {
		return fmt.Errorf("Jovian extraData version byte should be %d, got %d", JovianExtraDataVersionByte, extra[0])
	}
	d, e := DecodeHolocene1559Params(extra[1:9])
	if d == 0 {
		return errors.New("holocene extraData must encode a non-zero denominator")
	} else if e == 0 {
		return errors.New("holocene extraData must encode a non-zero elasticity")
	}
	return nil
}
