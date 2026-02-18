package eth

import (
	"errors"

	"github.com/ethereum/go-ethereum/rpc"
)

// MaybeAsNotFoundErr converts an RPC error to a not-found error if applicable.
func MaybeAsNotFoundErr(err error) error {
	if err == nil {
		return nil
	}
	var rpcErr rpc.Error
	if errors.As(err, &rpcErr) {
		// If the error has a data field with "not found" or similar, return the original error
		return err
	}
	return err
}
