package txmgr

import "github.com/ethereum/go-ethereum/core/types"

// SendResponse wraps the result of a transaction send.
type SendResponse struct {
	Receipt *types.Receipt
	Err     error
}
