package op_e2e

import (
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func EncodeCallData(items ...interface{}) []byte {
	var packed []byte

	for _, item := range items {
		switch v := item.(type) {
		case []byte:
			packed = append(packed, v...)
		case uint32:
			buf := make([]byte, 4)
			binary.BigEndian.PutUint32(buf, v)
			packed = append(packed, buf...)
		case *big.Int:
			bytes := v.Bytes()
			if len(bytes) < 32 {
				padding := make([]byte, 32-len(bytes))
				bytes = append(padding, bytes...)
			}
			packed = append(packed, bytes...)
		case common.Address:
			packed = append(packed, v.Bytes()...)
		}
	}
	return packed
}
