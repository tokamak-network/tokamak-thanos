package ioutil

// Progressor tracks progress of an operation.
type Progressor func(current, total int64)

// NoopProgressor returns a no-op progress tracker.
var NoopProgressor Progressor = func(_ int64, _ int64) {}
