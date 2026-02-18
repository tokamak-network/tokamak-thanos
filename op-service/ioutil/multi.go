package ioutil

import "io"

// MultiCloser is a collection of closers that can be closed together.
type MultiCloser []io.Closer

// Close closes all closers. Returns the first error encountered.
func (mc MultiCloser) Close() error {
	var firstErr error
	for _, c := range mc {
		if err := c.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
